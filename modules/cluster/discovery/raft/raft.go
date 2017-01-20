// Package store provides a simple distributed key-value store. The keys and
// associated values are changed via distributed consensus, meaning that the
// values are changed only when a majority of nodes in the cluster agree on
// the new value.
//
// Distributed consensus is provided via the Raft algorithm.
package raft

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	log "github.com/cihub/seelog"
	"github.com/hashicorp/raft"
	"github.com/hashicorp/raft-boltdb"
	apihandler "github.com/medcl/gopa/core/api"
	"github.com/medcl/gopa/core/global"
	"github.com/medcl/gopa/core/util"
	. "github.com/medcl/gopa/modules/cluster/discovery/config"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const (
	retainSnapshotCount = 2
	raftTimeout         = 10 * time.Second
)

type RaftModule struct {
	cfg  RaftConfig
	raft *raft.Raft
	mu   sync.Mutex
	m    map[string]string // The key-value store for the system.

}

func New() *RaftModule {
	cfg := RaftConfig{}
	cfg.Init()
	return &RaftModule{
		m:   make(map[string]string),
		cfg: cfg,
	}
}

type NullWriter int

func (NullWriter) Write([]byte) (int, error) { return 0, nil }

// handle cluster join function
func (s *RaftModule) handleJoin(w http.ResponseWriter, r *http.Request) {
	log.Debug("receive join request")

	m := map[string]string{}
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(m) != 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	remoteAddr, ok := m["addr"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := s.Join(remoteAddr); err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write([]byte(s.cfg.Bind))
}

func (s *RaftModule) clusterInfo(w http.ResponseWriter, r *http.Request) {

	stats := map[string]interface{}{}
	stats["leader"] = s.raft.Leader()
	stats["stats"] = s.raft.Stats()
	stats["addr"] = s.raft.String()
	b, _ := json.MarshalIndent(stats, "", "\t")
	w.Write(b)
	w.WriteHeader(http.StatusOK)
	return
}

// Open opens the store. If enableSingle is set, and there are no existing peers,
// then this node becomes the first node, and therefore leader, of the cluster.
func (s *RaftModule) Open() error {
	// Setup Raft configuration.
	config := raft.DefaultConfig()
	// Check for any existing peers.
	peers, err := readPeersJSON(filepath.Join(s.cfg.DataDir, "peers.json"))
	if err != nil {
		log.Error(err)
		return err
	}

	enableSingle := len(s.cfg.Seeds) == 0

	if !global.Env().IsDebug {
		//disable raft logging
		config.LogOutput = new(NullWriter)
	}

	log.Debug("cluster previous persisted seed peers: ", len(peers), ", ", strings.Join(peers, ","))

	// Allow the node to entry single-mode, potentially electing itself, if
	// explicitly enabled and there is only 1 node in the cluster already.
	if enableSingle && len(peers) <= 1 {
		log.Debug("raft enabling single-node mode")
		config.EnableSingleNode = true
		config.DisableBootstrapAfterElect = false
	}

	// Setup Raft communication.
	address := util.AutoGetAddress(s.cfg.Bind)
	addr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		log.Error(err)
		return err
	}
	transport, err := raft.NewTCPTransport(address, addr, 3, 10*time.Second, os.Stderr)
	if err != nil {
		log.Error(err)
		return err
	}

	log.Debug("raft listen on: ", address)

	// Create peer storage.
	peerStore := raft.NewJSONPeers(s.cfg.DataDir, transport)

	// Create the snapshot store. This allows the Raft to truncate the log.
	snapshots, err := raft.NewFileSnapshotStore(s.cfg.DataDir, retainSnapshotCount, os.Stderr)
	if err != nil {
		return fmt.Errorf("file snapshot store: %s", err)
	}

	// Create the log store and stable store.
	logStore, err := raftboltdb.NewBoltStore(filepath.Join(s.cfg.DataDir, "raft.db"))
	if err != nil {
		return fmt.Errorf("new bolt store: %s", err)
	}

	// Instantiate the Raft systems.
	ra, err := raft.NewRaft(config, (*fsm)(s), logStore, logStore, snapshots, peerStore, transport)
	if err != nil {
		return fmt.Errorf("new raft: %s", err)
	}
	s.raft = ra

	// If join was specified, make the join request.
	if len(s.cfg.Seeds) > 0 {
		for _, v := range s.cfg.Seeds {
			if err := join(v, address); err != nil {
				log.Errorf("failed to join node at %s: %s", v, err.Error())
			}
		}
	}

	if global.Env().IsDebug {
		go func() {
			t := time.NewTicker(time.Duration(1) * time.Second)

			for {
				select {
				case <-t.C:
					log.Debug("raft leader: ", ra.Leader())
				}
			}
		}()
	}

	apihandler.HandleFunc("/cluster/status", s.clusterInfo)
	apihandler.HandleFunc("/cluster/node/_join", s.handleJoin)

	apihandler.HandleFunc("/cache", s.handleKeyRequest)
	apihandler.HandleFunc("/cache/", s.handleKeyRequest)

	return nil
}

// sent join requests to seed host
func join(joinAddr, raftAddr string) error {

	log.Debug("start join address, ", joinAddr, ",", raftAddr)
	raftAddr = util.GetValidAddress(raftAddr)

	b, err := json.Marshal(map[string]string{"addr": raftAddr})
	if err != nil {
		log.Error(err)
		return err
	}

	joinAddr = util.GetValidAddress(joinAddr)

	if len(global.Env().SystemConfig.CertPath) > 0 {
		url := fmt.Sprintf("https://%s/cluster/node/_join", joinAddr)

		log.Info("try to join the cluster, ", url, ", ", string(b))

		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client := &http.Client{Transport: tr}
		resp, err := client.Post(url, "application-type/json", bytes.NewReader(b))
		if err != nil {
			log.Error("Get error:", err)
			return err
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Error(url, err)
			return err
		}
		log.Debug(string(body))
		return nil
	}

	url := fmt.Sprintf("http://%s/cluster/node/_join", joinAddr)

	log.Info("try to join the cluster, ", url, ", ", string(b))

	resp, err := http.Post(url, "application-type/json", bytes.NewReader(b))
	if err != nil {
		log.Error(err)
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(url, err)
		return err
	}

	log.Debug(string(body))
	defer resp.Body.Close()
	return nil
}

// Join joins a node, located at addr, to this store. The node must be ready to
// respond to Raft communications at that address.
func (s *RaftModule) Join(addr string) error {
	log.Debug("received join request for remote node as :", addr)

	f := s.raft.AddPeer(addr)
	if f.Error() != nil {
		log.Error(f.Error())
		return f.Error()
	}
	log.Info("node at %s joined successfully", addr)
	return nil
}

func readPeersJSON(path string) ([]string, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	if len(b) == 0 {
		return nil, nil
	}

	var peers []string
	dec := json.NewDecoder(bytes.NewReader(b))
	if err := dec.Decode(&peers); err != nil {
		return nil, err
	}

	return peers, nil
}

// handle cache function
func (s *RaftModule) handleKeyRequest(w http.ResponseWriter, r *http.Request) {

	getKey := func() string {
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) != 3 {
			return ""
		}
		return parts[2]
	}

	switch r.Method {
	case "GET":
		k := getKey()
		if k == "" {
			w.WriteHeader(http.StatusBadRequest)
		}
		v, err := s.Get(k)
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		b, err := json.Marshal(map[string]string{k: v})
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		io.WriteString(w, string(b))

	case "POST":
		// Read the value from the POST body.
		m := map[string]string{}
		if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		for k, v := range m {
			if err := s.Set(k, v); err != nil {
				log.Error(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

	case "DELETE":
		k := getKey()
		if k == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if err := s.Delete(k); err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		s.Delete(k)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
	return
}

// Get returns the value for the given key.
func (s *RaftModule) Get(key string) (string, error) {
	log.Trace("gettting ,", key)

	s.mu.Lock()
	defer s.mu.Unlock()
	return s.m[key], nil
}

// Set sets the value for the given key.
func (s *RaftModule) Set(key, value string) error {

	log.Trace("setting ,", key, ",", value)

	log.Error(s.raft)
	if s.raft.State() != raft.Leader {
		return fmt.Errorf("not leader")
	}

	c := &Command{
		Op:    "set",
		Key:   key,
		Value: value,
	}
	b, err := json.Marshal(c)
	if err != nil {
		return err
	}

	f := s.raft.Apply(b, raftTimeout)
	return f.Error()
}

// Delete deletes the given key.
func (s *RaftModule) Delete(key string) error {
	if s.raft.State() != raft.Leader {
		return fmt.Errorf("not leader")
	}

	c := &Command{
		Op:  "delete",
		Key: key,
	}
	b, err := json.Marshal(c)
	if err != nil {
		return err
	}

	f := s.raft.Apply(b, raftTimeout)
	return f.Error()
}

type fsm RaftModule

// Apply applies a Raft log entry to the key-value store.
func (f *fsm) Apply(l *raft.Log) interface{} {
	var c Command
	if err := json.Unmarshal(l.Data, &c); err != nil {
		panic(fmt.Sprintf("failed to unmarshal command: %s", err.Error()))
	}

	switch c.Op {
	case "set":
		return f.applySet(c.Key, c.Value)
	case "delete":
		return f.applyDelete(c.Key)
	default:
		panic(fmt.Sprintf("unrecognized command op: %s", c.Op))
	}
}

// Snapshot returns a snapshot of the key-value store.
func (f *fsm) Snapshot() (raft.FSMSnapshot, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	// Clone the map.
	o := make(map[string]string)
	for k, v := range f.m {
		o[k] = v
	}
	return &fsmSnapshot{store: o}, nil
}

// Restore stores the key-value store to a previous state.
func (f *fsm) Restore(rc io.ReadCloser) error {
	o := make(map[string]string)
	if err := json.NewDecoder(rc).Decode(&o); err != nil {
		return err
	}

	// Set the state from the snapshot, no lock required according to
	// Hashicorp docs.
	f.m = o
	return nil
}

func (f *fsm) applySet(key, value string) interface{} {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.m[key] = value
	return nil
}

func (f *fsm) applyDelete(key string) interface{} {
	f.mu.Lock()
	defer f.mu.Unlock()
	delete(f.m, key)
	return nil
}

type fsmSnapshot struct {
	store map[string]string
}

func (f *fsmSnapshot) Persist(sink raft.SnapshotSink) error {
	err := func() error {
		// Encode data.
		b, err := json.Marshal(f.store)
		if err != nil {
			return err
		}

		// Write data to sink.
		if _, err := sink.Write(b); err != nil {
			return err
		}

		// Close the sink.
		if err := sink.Close(); err != nil {
			return err
		}

		return nil
	}()

	if err != nil {
		sink.Cancel()
		return err
	}

	return nil
}

func (f *fsmSnapshot) Release() {}
