/**
@author: iswarezwp
Created on 2016-12-16 09:23
**/

package dconf

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
)

func eventProcessor(watcher *fsnotify.Watcher, callback func()) {
	defer watcher.Close()

	fired := false

	for {
		select {
		case event := <-watcher.Events:
			if event.Op&fsnotify.Write == fsnotify.Write {
				fmt.Println("file modified")
				if !fired {
					// Avoid multiple reload at the same time
					fired = true
					go func() {
						callback()
						fired = false
					}()
				}
			}
		case err := <-watcher.Errors:
			fmt.Println("error:", err)
			return
		}
	}
}

func watchFile(filename string, callback func()) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	err = watcher.Add(filename)
	if err != nil {
		return err
	}

	go eventProcessor(watcher, callback)
	return nil
}
