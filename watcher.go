package gomon

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"time"
	"unsafe"

	"github.com/fsnotify/fsnotify"
)

type Watcher struct {
	base   *fsnotify.Watcher
	config Config
	Events chan Event
	Errors chan error
	save   *os.File
	dates  map[string]time.Time
}

func NewWatcher(configs ...Config) (*Watcher, error) {
	var err error
	w := &Watcher{config: mergeConfigs(configs...)}

	w.base, err = fsnotify.NewBufferedWatcher(w.config.BufferSize)
	if err != nil {
		return nil, err
	}

	w.Events = *(*chan Event)(unsafe.Pointer(&w.base.Events))
	w.Errors = w.base.Errors

	if w.config.Save != "" {
		w.save, err = os.OpenFile(w.config.Save, os.O_CREATE|os.O_RDWR, os.ModePerm)
		if err != nil {
			// You can continue without save system.
			return w, err
		}

		if err := json.NewDecoder(w.save).Decode(&w.dates); err != nil && !errors.Is(err, io.EOF) {
			// You can continue without save system.
			return w, err
		}
	}
	return w, nil
}

func (w Watcher) Add(path string) (Op, error) {
	var op Op
	if w.dates != nil {
		info, err := os.Stat(path)
		if err != nil {
			if os.IsNotExist(err) {
				if w.dates[path].IsZero() {
					return Non, err
				}
				op = Remove | Rename
			} else {
				return Non, err
			}
		} else {
			if !w.dates[path].IsZero() && !w.dates[path].Equal(info.ModTime()) {
				op = Write | Chmod
			}
		}
	}
	return op, w.base.Add(path)
}

func (w Watcher) Close() error {
	if w.save != nil {
		names := w.base.WatchList()
		if err := w.base.Close(); err != nil {
			return err
		}
		dates := make(map[string]time.Time, len(names))
		for _, name := range names {
			info, err := os.Stat(name)
			// TODO: Think renaming and dates size
			if err != nil {
				continue
			}
			dates[name] = info.ModTime()
		}

		if err := w.save.Truncate(0); err != nil {
			// TODO: It is very likely to get an error, return should not be made directly.
			// We can continue without save.
			return err
		}
		if _, err := w.save.Seek(0, 0); err != nil {
			// TODO: It is very likely to get an error, return should not be made directly.
			// We can continue without save.
			return err
		}

		if err := json.NewEncoder(w.save).Encode(dates); err != nil {
			// TODO: It is very likely to get an error, return should not be made directly.
			// We can continue without save.
			return err
		}
		return w.save.Close()
	}
	return w.base.Close()
}

func (w Watcher) Config() Config {
	return w.config
}
