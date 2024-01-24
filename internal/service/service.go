package service

import (
	"log"
	"runtime"

	"github.com/akinbezatoglu/s3ync/internal/watcher"
	"github.com/takama/daemon"
)

const (
	name        = "s3ync"
	description = "Watcher service to catch file system events."
	startCmd    = "start"
	stopCmd     = "stop"
	upStatus    = "up"
	downStatus  = "down"
)

var (
	stdlog, errlog *log.Logger
	command        string
	status         string
)

func SetStartCmd() {
	command = startCmd
}

func SetStopCmd() {
	command = stopCmd
}

type Service struct {
	daemon.Daemon
}

func InitService() *Service {
	// to init service, first set the status up
	status = upStatus

	var d daemon.Daemon
	if runtime.GOOS == "darwin" {
		// for MacOS
		d, _ = daemon.New(name, description, daemon.GlobalDaemon)
	} else {
		d, _ = daemon.New(name, description, daemon.SystemDaemon)
	}
	return &Service{d}
}

func (s *Service) Manage() {
	// TODO: When the service runs, it must listen for commands from the CLI to receive a stop command. This needs to be async. Meanwhile, the watcher will continue listening to the file system.
	go func() {
		for {
			switch status {
			case downStatus:
			case upStatus:
			}
		}
	}()
	initWatcher()
	s.Stop()
}

func initWatcher() *watcher.Watcher {
	w, err := watcher.NewWatcher()
	if err != nil {
		errlog.Println("Error creating watcher:", err)
	}
	defer w.Close()
	w.AddRecursive("")
	w.Watch()
	return w
}
