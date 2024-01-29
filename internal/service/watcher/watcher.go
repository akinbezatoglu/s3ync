package watcher

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/akinbezatoglu/s3ync/internal/service/config"
	"github.com/fsnotify/fsnotify"
)

type Watcher struct {
	*fsnotify.Watcher
}

var (
	done    chan struct{}  // waits for a signal to stop
	addPath chan string    // receives a path to add to the watcher
	rmPath  chan string    // receives a path to remove from the watcher
	wg      sync.WaitGroup // waits for all events to be handled.
)

var Syncs *config.Syncs

func InitWatcher() (*Watcher, error) {
	done = make(chan struct{})
	addPath = make(chan string)
	rmPath = make(chan string)

	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, &WatcherFailedInitError{err}
	}

	// Gets all pre-configured paths and bucket infos from config file
	Syncs = config.GetAllSyncList()

	return &Watcher{w}, nil
}

// AddPathsAlreadyConfigured adds pre-configured paths to the watcher
func (w *Watcher) AddPathsAlreadyConfigured() {
	if len(Syncs.All) != 0 {
		// map[string][]string
		// Key: filesytem path, Value: [aws-profile, bucket-name]
		for path := range Syncs.All {
			home, _ := os.UserHomeDir()
			// Watcher runs in a container. So event.Name/s is related to container paths.
			// Removes home/userprofile from path to set container's home as container path
			w.AddPathRecursive(strings.Replace(path, home, "", 1))
		}
	}
}

// (*fsnotify.Watcher).Add() function do not add recursively.
// AddPathRecursive recursively adds all directories inside the root path to the watcher.
func (w *Watcher) AddPathRecursive(root string) {
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			w.Add(path)
		}
		return nil
	})
}

// (*fsnotify.Watcher).Remove() function do not remove recursively.
// RemovePathRecursive recursively removes all sub-directories of the root path from the watcher.
func (w *Watcher) RemovePathRecursivee(root string) {
	watches := w.WatchList()
	for _, watch := range watches {
		if strings.Contains(watch, root) {
			// watch is a subdir of the root dir
			w.Remove(watch)
		}
	}
}

// Watch watches and handles events async
func (w *Watcher) Watch() error {
	for {
		select {
		case event, ok := <-w.Events:
			if !ok {
				return nil
			}
			wg.Add(1)
			go handleEvent(event)
		case err, ok := <-w.Errors:
			if !ok {
				return &EventError{err}
			}
		case <-done:
			wg.Wait()
			return nil
		}
	}
}

// Close path channels and lastly done channel to send stop signal to the watcher
func (w *Watcher) Stop() {
	close(addPath)
	close(rmPath)
	close(done)
}

func handleEvent(e fsnotify.Event) {
	defer wg.Done()
	if e.Has(fsnotify.Create) {
		// Handle file/directory creation
		fmt.Println("File/Direcory created:", e.Name)
	}
	if e.Has(fsnotify.Write) {
		// Handle file modification
		fmt.Println("File modified:", e.Name)
	}
	if e.Has(fsnotify.Remove) {
		// Handle file/directory deletion
		fmt.Println("File/Direcory removed:", e.Name)
	}
	if e.Has(fsnotify.Rename) {
		// Handle file/directory renaming
		fmt.Println("File/Direcory renamed:", e.Name)
	}
}
