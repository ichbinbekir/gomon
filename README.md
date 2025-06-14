# Gomon

You can capture file or folder changes, but that's not all.

Gomon saves the last modification dates of the watched files, so it can detect any changes when it is started at a later time.

## ⚡️ Quickstart

```go
package main

import (
  "log"

  "github.com/ichbinbekir/gomon"
)

func main() {
  watcher, err := gomon.NewWatcher(Config{Save: "dates.json"})
  if err != nil {
    log.Fatal(err)
  }
  defer watcher.Close()

  if _, err := watcher.Add("test.txt"); err != nil {
    log.Fatal(err)
  }

  for event := range watcher.Events {
    // Print file modifications
    log.Println(event)
  }
}
```

## ⚙️ Installation

```bash
go get -u github.com/ichbinbekir/keyboard
```
