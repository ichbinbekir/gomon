package gomon

import "github.com/fsnotify/fsnotify"

type Event struct {
	Name string
	Op   Op
}

func (e Event) Has(op Op) bool {
	return gomon2fsnotify(e).Has(fsnotify.Op(op))
}

func (e Event) String() string {
	return gomon2fsnotify(e).String()
}

func gomon2fsnotify(e Event) fsnotify.Event {
	return fsnotify.Event{Name: e.Name, Op: fsnotify.Op(e.Op)}
}
