package service_chrome

import (
	"fmt"
	log "github.com/cihub/seelog"
	. "github.com/infinitbyte/gopa/core/config"
	"os/exec"
)

type ChromePlugin struct {
}

func (plugin ChromePlugin) Name() string {
	return "Chrome"
}

var cmd *exec.Cmd

func (plugin ChromePlugin) Start(cfg *Config) {

	config := struct {
		Command   string `config:"command"`
		DebugPort string `config:"debug_port"`
	}{DebugPort: "9223", Command: "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"}

	cfg.Unpack(&config)

	go func() {
		cmd = exec.Command(config.Command, "--headless", "-disable-gpu", fmt.Sprintf("--remote-debugging-port=%v", config.DebugPort), "--no-sandbox")
		err := cmd.Start()
		if err != nil {
			panic(err)
		}
		err = cmd.Wait()
		if err != nil {
			log.Debug(err)
		}
		log.Debug("chrome service normal exit")

	}()
}

func (plugin ChromePlugin) Stop() error {
	if cmd != nil {
		if cmd.ProcessState != nil {
			if !cmd.ProcessState.Exited() {
				return cmd.Process.Kill()
			}
		}
	}
	return nil
}
