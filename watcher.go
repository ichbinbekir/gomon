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
			return w, &SaveError{Op: "open file", Err: err}
		}

		if err := json.NewDecoder(w.save).Decode(&w.dates); err != nil && !errors.Is(err, io.EOF) {
			return w, &SaveError{Op: "decode", Err: err}
		}
	}
	return w, nil
}

func (w *Watcher) Add(path string) (Op, error) {
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

func (w *Watcher) Close() error {
	if w.save != nil {
		names := w.base.WatchList()
		if err := w.base.Close(); err != nil {
			return err
		}
		return saveDates(w.save, names)
	}
	return w.base.Close()
}

func (w *Watcher) Config() Config {
	return w.config
}

func saveDates(file *os.File, names []string) error {
	var statErr error
	dates := make(map[string]time.Time, len(names))
	for _, name := range names {
		info, err := os.Stat(name)
		if err != nil {
			statErr = errors.Join(statErr, err)
			continue
		}
		dates[name] = info.ModTime()
	}

	if err := file.Truncate(0); err != nil {
		return &SaveError{Op: "truncate", Err: errors.Join(statErr, err)}
	}
	if _, err := file.Seek(0, 0); err != nil {
		return &SaveError{Op: "seek", Err: errors.Join(statErr, err)}
	}
	if err := json.NewEncoder(file).Encode(dates); err != nil {
		return &SaveError{Op: "encode", Err: errors.Join(statErr, err)}
	}
	if err := file.Close(); err != nil {
		return &SaveError{Op: "close file", Err: errors.Join(statErr, err)}
	}
	return statErr
}
