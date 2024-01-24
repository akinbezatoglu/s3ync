package watcher

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/fsnotify/fsnotify"
)

var (
	done chan bool
	wg   sync.WaitGroup
)

type Watcher struct {
	*fsnotify.Watcher
}

func NewWatcher() (*Watcher, error) {
	done = make(chan bool)
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	return &Watcher{w}, nil
}

func (w *Watcher) AddRecursive(path string) {
	filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			w.Add(path)
		}
		return nil
	})
}

func (w *Watcher) Watch() {
	go func() {
		for {
			select {
			case event, ok := <-w.Events:
				if !ok {
					return
				}
				wg.Add(1)
				handleEvent(event)
			case err, ok := <-w.Errors:
				if !ok {
					return
				}
				fmt.Println("Error:", err)
			case <-done:
				return
			}
		}
	}()
	wg.Wait()
}

func Stop() {
	close(done)
}

func handleEvent(e fsnotify.Event) {
	go func() {
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
	}()
}
