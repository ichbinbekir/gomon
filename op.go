package gomon

import "github.com/fsnotify/fsnotify"

type Op fsnotify.Op

const (
	// Default fsnotify
	Create Op = 1 << iota
	Write
	Remove
	Rename
	Chmod

	// Special for gomon
	Non Op = 0
)

func (o Op) Has(h Op) bool {
	return fsnotify.Op(h).Has(fsnotify.Op(h))
}

func (o Op) String() string {
	if o == Non {
		return "NON"
	}
	return fsnotify.Op(o).String()
}
