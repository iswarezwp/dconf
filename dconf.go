/**
@author: iswarezwp
Created on 2016-12-15 18:20
**/

package dconf

import (
	"bufio"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"io"
	"os"
	"strings"
	"sync"
)

type DConf struct {
	filename string
	conf     map[string]map[string]string

	// If true, reload configurations when file changed.
	reload bool

	// Need a mutex when support reload on change
	watcher *fsnotify.Watcher
	mutex   *sync.Mutex
	loadOK  bool
}

func NewDConf(filename string, reload bool) (*DConf, error) {
	d := &DConf{
		filename: filename,
		conf:     make(map[string]map[string]string),
		reload:   reload,
		mutex:    &sync.Mutex{},
		loadOK:   true,
	}
	if err := d.Load(); err != nil {
		if os.IsNotExist(err) {
			d.loadOK = false
		} else {
			return nil, err
		}
	}

	return d, nil
}

func (d *DConf) lock() {
	if d.reload {
		d.mutex.Lock()
	}
}

func (d *DConf) unlock() {
	if d.reload {
		d.mutex.Unlock()
	}
}

func (d *DConf) Get(key, defaultValue string) string {
	return d.GetValue("Default", key, defaultValue)
}

func (d *DConf) GetValue(sec, key, defaultValue string) string {
	// If load configuration failed on start, try to load it first.
	// This will ensure that all changes being tracked even if the
	// fsnotify backend not work.
	if !d.loadOK {
		if err := d.Load(); err == nil {
			d.loadOK = true
		}
	}

	d.lock()
	defer d.unlock()
	if s, ok := d.conf[sec]; ok {
		if v, ok := s[key]; ok {
			return v
		}
	}

	return defaultValue
}

func (d *DConf) IsLoaded() bool {
	return d.loadOK
}

func (d *DConf) setValue(sec, key, value string) {
	key = strings.TrimSpace(key)
	if key == "" {
		return
	}

	d.lock()
	defer d.unlock()

	if s, ok := d.conf[sec]; ok {
		s[key] = value
	} else {
		s := make(map[string]string)
		s[key] = value
		d.conf[sec] = s
	}
}

func (d *DConf) OnFileChange() {
	d.Load()
}

func (d *DConf) Load() error {
	fr, err := os.Open(d.filename)
	if err != nil {
		return err
	}
	defer fr.Close()

	d.conf = make(map[string]map[string]string)

	buf := bufio.NewReader(fr)

	// Handle BOM-UTF8.
	mask, err := buf.Peek(3)
	if err == nil && len(mask) >= 3 &&
		mask[0] == 239 && mask[1] == 187 && mask[2] == 191 {
		buf.Read(mask)
	}

	var (
		i               int
		sec, key, value string
		valQuote        string
	)

	// Parse line-by-line
	for {
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)
		lineLengh := len(line)

		// Skip the empty line or comment line
		if lineLengh != 0 && line[0] != '#' {
			if line[0] == '[' && line[len(line)-1] == ']' {
				// The section line
				sec = strings.TrimSpace(line[1 : len(line)-1])
			} else {
				// The `key=value` line must behind the section line
				if sec == "" {
					// support non-section configuration
					sec = "Default"
				}

				i = strings.IndexAny(line, "=")
				if i <= 0 {
					// Ignore error config iterm
					fmt.Println("Config error, ignore: ", line)
					continue
				}
				key = strings.TrimSpace(line[0:i])

				// Support double quote value like: `key = "value"`
				lineRight := strings.TrimSpace(line[i+1:])
				lineRightLength := len(lineRight)
				if lineRightLength >= 2 {
					valQuote = lineRight[0:1]
				}

				if valQuote == `"` {
					qLen := len(valQuote)
					pos := strings.LastIndex(lineRight[qLen:], valQuote)
					if pos == -1 {
						fmt.Println("Config error, ignore: ", line)
					}
					pos = pos + qLen
					value = lineRight[qLen:pos]
				} else {
					value = strings.TrimSpace(lineRight[0:])
				}

				d.setValue(sec, key, value)
			}
		}

		if err != nil {
			// Reached end of file.
			if err == io.EOF {
				break
			}
			return err
		}
	}

	if d.reload {
		d.watcher, err = watchFile(d.filename, d.OnFileChange)
		if err != nil {
			return err
		}
	}

	return nil
}

func (d *DConf) Close() {
	if d.reload {
		d.watcher.Close()
		d.loadOK = false
	}
}
