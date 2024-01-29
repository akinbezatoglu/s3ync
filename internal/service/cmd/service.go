package main

import (
	"fmt"

	"github.com/akinbezatoglu/s3ync/internal/service/watcher"
)

func main() {
	w, err := watcher.InitWatcher()
	if err != nil {
		fmt.Println(err)
	}
	defer w.Close()

	w.AddPathsAlreadyConfigured()
	w.Watch()
}
