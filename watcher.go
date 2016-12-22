/**
@author: iswarezwp
Created on 2016-12-16 09:23
**/

package dconf

import (
	"github.com/fsnotify/fsnotify"
)

func eventProcessor(watcher *fsnotify.Watcher, callback func()) {
	defer watcher.Close()

	fired := false

	for {
		select {
		case event := <-watcher.Events:
			if event.Op&fsnotify.Write == fsnotify.Write {
				if !fired {
					// Avoid multiple reload at the same time
					fired = true
					go func() {
						callback()
						fired = false
					}()
				}
			}
		case <-watcher.Errors:
			return
		}
	}
}

func watchFile(filename string, callback func()) (*fsnotify.Watcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	err = watcher.Add(filename)
	if err != nil {
		return nil, err
	}

	go eventProcessor(watcher, callback)
	return watcher, nil
}
