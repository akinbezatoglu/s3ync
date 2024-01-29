package watcher

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/akinbezatoglu/s3ync/internal/service/config"
	"github.com/akinbezatoglu/s3ync/internal/service/ops"
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

var syncs *config.Syncs
var bucketbasics *ops.BucketBasics

func InitWatcher() (*Watcher, error) {
	done = make(chan struct{})
	addPath = make(chan string)
	rmPath = make(chan string)

	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, &WatcherFailedInitError{Err: err}
	}

	// Gets all pre-configured paths and bucket infos from config file
	syncs = config.GetAllSyncList()

	bucketbasics, err = ops.NewBucketBasics()
	if err != nil {
		return nil, &ops.S3ClientFailedError{Err: err}
	}

	return &Watcher{w}, nil
}

// AddPathsAlreadyConfigured adds pre-configured paths to the watcher
func (w *Watcher) AddPathsAlreadyConfigured() {
	if len(syncs.All) != 0 {
		// map[string][]string
		// Key: filesytem path, Value: [aws-profile, bucket-name]
		for path := range syncs.All {
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

func (w *Watcher) AddPathRecursiveAndUpload(root, bucketname, rootDirObjectKey, profile string) {
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			// If it is a directory, add it to the watcher.
			// Only directories will be added to the watcher.
			w.Add(path)
		} else {
			// If added is a file, It does not need to be added to the watcher.
			// Just upload it to the bucket.
			relativepath := strings.TrimPrefix(filepath.ToSlash(path), filepath.ToSlash(root))
			go bucketbasics.UploadFile(bucketname, rootDirObjectKey+relativepath, path, profile)
		}
		return nil
	})
}

// (*fsnotify.Watcher).Remove() function do not remove recursively.
// RemovePathRecursive recursively removes all sub-directories of the root path from the watcher.
func (w *Watcher) RemovePathRecursive(root string) {
	watches := w.WatchList()
	for _, watch := range watches {
		if strings.Contains(watch, root) {
			// watch is a subdir of the root dir
			err := w.Remove(filepath.ToSlash(watch))
			if err != nil {
				fmt.Println(err)
			}
			//fmt.Println(w.WatchList())
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
			go w.handleEvent(event)
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
	w.Close()
}

func (w *Watcher) handleEvent(e fsnotify.Event) {
	defer wg.Done()
	var profile, bucketname, relativepath string
	for path := range syncs.All {
		if strings.Contains(filepath.ToSlash(e.Name), path) {
			profile = syncs.All[path][0]
			bucketname = syncs.All[path][1]

			//ex. /path/to/watch -> event (Create): /path/to/watch/file1.txt
			// 	  relativepath: file1.txt
			relativepath = strings.TrimPrefix(filepath.ToSlash(e.Name), filepath.ToSlash(path))[1:]
			break
		}
	}

	if e.Has(fsnotify.Create) {
		if fileInfo, err := os.Stat(e.Name); err == nil && fileInfo.IsDir() {
			if empty, err := isDirEmpty(e.Name); err == nil && empty {
				fmt.Printf("Created empty directory: %q", e.Name)
				// created an empty directory
				w.Add(e.Name)
			} else {
				fmt.Printf("Created non-empty directory: %q\n", e.Name)
				// moved from a directory to the watched-directory
				w.AddPathRecursiveAndUpload(e.Name, bucketname, relativepath, profile)
			}
		} else {
			fmt.Printf("Created file: %q\n", e.Name)
			bucketbasics.UploadFile(bucketname, relativepath, e.Name, profile)
		}
		return
	}

	if e.Has(fsnotify.Write) {
		if fileInfo, err := os.Stat(e.Name); err == nil && !fileInfo.IsDir() {
			fmt.Printf("Modified file: %q\n", e.Name)
			// All directories are watched recursively.
			// Receiving a Write event from a directory is redundant.
			// File updates are necessary only in the presence of a Write event specific to a file.
			bucketbasics.UploadFile(bucketname, relativepath, e.Name, profile)
		}
		return
	}

	// When a file/directory is renamed, the Watcher catches two events.
	// Rename and  Create. Rename has the old name, and Create has the new name.
	// Rename logic needs to be handled in two events.
	// Respectively, Rename removes first the old named file/directory.
	// So, Rename shares the same logic with the Remove event.
	// After, Create will upload the newly named file/directory.
	if e.Has(fsnotify.Remove) || e.Has(fsnotify.Rename) {
		fmt.Printf("Removed directory: %q\n", e.Name)
		w.RemovePathRecursive(e.Name)
		bucketbasics.DeleteDirectory(bucketname, relativepath, profile)
		return
	}
}

func isDirEmpty(name string) (bool, error) {
	f, err := os.Open(name)
	if err != nil {
		return false, err
	}
	defer f.Close()

	// read in ONLY one file
	_, err = f.Readdir(1)

	// and if the file is EOF... well, the dir is empty.
	if err == io.EOF {
		return true, nil
	}
	return false, err
}
