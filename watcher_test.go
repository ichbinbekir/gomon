package gomon

import (
	"os"
	"sync"
	"testing"
	"time"
)

func TestMain(t *testing.T) {
	if err := os.MkdirAll(".test", os.ModePerm); err != nil {
		t.Fatal(err)
	}
	file, err := os.Create(".test/test.txt")
	if err != nil {
		t.Fatal(err)
	}
	file.Close()

	watcher, err := NewWatcher(Config{Save: ".test/dates.json"})
	if err != nil {
		t.Fatal(err)
	}

	var wg sync.WaitGroup
	defer wg.Wait()
	wg.Add(1)
	go func() {
		time.Sleep(time.Second * 5)
		t.Log("closing...")
		if err := watcher.Close(); err != nil {
			t.Error(err)
		}
		wg.Done()
	}()

	go func() {
		for {
			time.Sleep(time.Second)
			os.WriteFile(".test/test.txt", []byte("asdasd"), os.ModePerm)
		}
	}()

	op, err := watcher.Add(".test/test.txt")
	if err != nil {
		t.Error(err)
	}
	t.Log(op)

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				break
			}
			t.Log(event)
		case err := <-watcher.Errors:
			t.Error(err)
		}
	}
}
