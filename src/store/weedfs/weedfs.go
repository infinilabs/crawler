package main

import (
	"encoding/json"
	"flag"
	"fmt"
	log "github.com/cihub/seelog"
	"math/rand"
	"mime"
	"net/http"
	"store/weedfs/directory"
	"store/weedfs/storage"
	"strconv"
	"strings"
	"time"
)

var (
	//master
	clusterName       = flag.String("cluster", "gopa", "cluster name")
	port              = flag.Int("port", 9333, "http listen port")
	metaFolder        = flag.String("mdir", "/tmp", "data directory to store mappings")
	capacity          = flag.Int("capacity", 100, "maximum number of volumes to hold")
	mapper            *directory.Mapper
	IsDebug           = flag.Bool("debug", false, "verbose debug information")
	volumeSizeLimitMB = flag.Uint("volumeSizeLimitMB", 32*1024, "Default Volume Size in MegaBytes")

	//volume
	storePort   = flag.Int("storePort", 8080, "http listen port")
	chunkFolder = flag.String("dir", "/tmp", "data directory to store files")
	volumes     = flag.String("volumes", "0,1-3,4", "comma-separated list of volume ids or range of ids")
	publicUrl   = flag.String("publicUrl", "localhost:8080", "public url to serve data read")
	metaServer  = flag.String("mserver", "localhost:9333", "master directory server to store mappings")

	pulse = flag.Int("pulseSeconds", 5, "number of seconds between heartbeats")

	store *storage.Store
)

//weedfs master
func dirLookupHandler(w http.ResponseWriter, r *http.Request) {
	vid := r.FormValue("volumeId")
	commaSep := strings.Index(vid, ",")
	if commaSep > 0 {
		vid = vid[0:commaSep]
	}
	volumeId, _ := strconv.ParseUint(vid, 10, 64)
	machine, e := mapper.Get(uint32(volumeId))
	if e == nil {
		writeJson(w, r, machine.Server)
	} else {
		log.Error("Invalid volume id", volumeId)
		writeJson(w, r, map[string]string{"error": "volume id " + strconv.FormatUint(volumeId, 10) + " not found"})
	}
}
func dirAssignHandler(w http.ResponseWriter, r *http.Request) {
	c := r.FormValue("count")
	fid, count, machine, err := mapper.PickForWrite(c)
	if err == nil {
		writeJson(w, r, map[string]string{"fid": fid, "url": machine.Url, "publicUrl": machine.PublicUrl, "count": strconv.Itoa(count)})
	} else {
		log.Error(err)
		writeJson(w, r, map[string]string{"error": err.Error()})
	}
}
func dirJoinHandler(w http.ResponseWriter, r *http.Request) {
	s := r.RemoteAddr[0:strings.Index(r.RemoteAddr, ":")+1] + r.FormValue("port")
	publicUrl := r.FormValue("publicUrl")
	volumes := new([]storage.VolumeInfo)
	json.Unmarshal([]byte(r.FormValue("volumes")), volumes)
	if *IsDebug {
		log.Info(s, "volumes", r.FormValue("volumes"))
	}
	mapper.Add(*directory.NewMachine(s, publicUrl, *volumes))
}
func dirStatusHandler(w http.ResponseWriter, r *http.Request) {
	writeJson(w, r, mapper)
}

func writeJson(w http.ResponseWriter, r *http.Request, obj interface{}) {
	w.Header().Set("Content-Type", "application/javascript")
	bytes, _ := json.Marshal(obj)
	callback := r.FormValue("callback")
	if callback == "" {
		w.Write(bytes)
	} else {
		w.Write([]uint8(callback))
		w.Write([]uint8("("))
		fmt.Fprint(w, string(bytes))
		w.Write([]uint8(")"))
	}
}

//weedfs volume
func statusHandler(w http.ResponseWriter, r *http.Request) {
	writeJson(w, r, store.Status())
}
func addVolumeHandler(w http.ResponseWriter, r *http.Request) {
	store.AddVolume(r.FormValue("volume"))
	writeJson(w, r, store.Status())
}
func storeHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		GetHandler(w, r)
	case "DELETE":
		DeleteHandler(w, r)
	case "POST":
		PostHandler(w, r)
	}
}
func GetHandler(w http.ResponseWriter, r *http.Request) {
	n := new(storage.Needle)
	vid, fid, ext := parseURLPath(r.URL.Path)
	volumeId, _ := strconv.ParseUint(vid, 10, 64)
	n.ParsePath(fid)

	if *IsDebug {
		log.Info("volume", volumeId, "reading", n)
	}
	cookie := n.Cookie
	count, e := store.Read(volumeId, n)
	if *IsDebug {
		log.Info("read bytes", count, "error", e)
	}
	if n.Cookie != cookie {
		log.Info("request with unmaching cookie from ", r.RemoteAddr, "agent", r.UserAgent())
		return
	}
	if ext != "" {
		w.Header().Set("Content-Type", mime.TypeByExtension(ext))
	}
	w.Write(n.Data)
}
func PostHandler(w http.ResponseWriter, r *http.Request) {
	vid, _, _ := parseURLPath(r.URL.Path)
	volumeId, e := strconv.ParseUint(vid, 10, 64)
	if e != nil {
		writeJson(w, r, e)
	} else {
		needle, ne := storage.NewNeedle(r)
		if ne != nil {
			writeJson(w, r, ne)
		} else {
			ret := store.Write(volumeId, needle)
			m := make(map[string]uint32)
			m["size"] = ret
			writeJson(w, r, m)
		}
	}
}
func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	n := new(storage.Needle)
	vid, fid, _ := parseURLPath(r.URL.Path)
	volumeId, _ := strconv.ParseUint(vid, 10, 64)
	n.ParsePath(fid)

	if *IsDebug {
		log.Info("deleting", n)
	}

	cookie := n.Cookie
	count, ok := store.Read(volumeId, n)

	if ok != nil {
		m := make(map[string]uint32)
		m["size"] = 0
		writeJson(w, r, m)
		return
	}

	if n.Cookie != cookie {
		log.Info("delete with unmaching cookie from ", r.RemoteAddr, "agent", r.UserAgent())
		return
	}

	n.Size = 0
	store.Delete(volumeId, n)
	m := make(map[string]uint32)
	m["size"] = uint32(count)
	writeJson(w, r, m)
}

func parseURLPath(path string) (vid, fid, ext string) {
	sepIndex := strings.LastIndex(path, "/")
	commaIndex := strings.LastIndex(path[sepIndex:], ",")
	if commaIndex <= 0 {
		if "favicon.ico" != path[sepIndex+1:] {
			log.Info("unknown file id", path[sepIndex+1:])
		}
		return
	}
	dotIndex := strings.LastIndex(path[sepIndex:], ".")
	vid = path[sepIndex+1 : commaIndex]
	fid = path[commaIndex+1:]
	ext = ""
	if dotIndex > 0 {
		fid = path[commaIndex+1 : dotIndex]
		ext = path[dotIndex+1:]
	}
	return
}

func main() {
	flag.Parse()

	//volume master block
	log.Info("Volume Size Limit is", *volumeSizeLimitMB, "MB")
	mapper = directory.NewMapper(*metaFolder, "directory", uint64(*volumeSizeLimitMB)*1024*1024)

	//weedfs master handler
	http.HandleFunc("/dir/assign", dirAssignHandler)
	http.HandleFunc("/dir/lookup", dirLookupHandler)
	http.HandleFunc("/dir/join", dirJoinHandler)
	http.HandleFunc("/dir/status", dirStatusHandler)

	//start weedfs master
	go func() {
		log.Info("Start directory service at http://127.0.0.1:" + strconv.Itoa(*port))
		e := http.ListenAndServe(":"+strconv.Itoa(*port), nil)
		if e != nil {
			log.Error("Fail to start:", e)
		}
	}()

	//volume block
	//TODO: now default to 1G, this value should come from server?
	store = storage.NewStore(*storePort, *publicUrl, *chunkFolder, *volumes)
	defer store.Close()

	//weedfs volume handler
	http.HandleFunc("/", storeHandler)
	http.HandleFunc("/status", statusHandler)
	http.HandleFunc("/add_volume", addVolumeHandler)

	go func() {
		for {
			store.Join(*metaServer)
			time.Sleep(time.Duration(float32(*pulse*1e3)*(1+rand.Float32())) * time.Millisecond)
		}
	}()
	log.Info("store joined at", *metaServer)

	//start weedfs volume
	go func() {
		log.Info("Start storage service at http://127.0.0.1:"+strconv.Itoa(*storePort), "public url", *publicUrl)
		e := http.ListenAndServe(":"+strconv.Itoa(*storePort), nil)
		if e != nil {
			log.Error("Fail to start:", e)
		}
	}()

}
