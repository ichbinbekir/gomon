package gomon

import "testing"

func TestConfig(t *testing.T) {
	config := mergeConfigs(Config{BufferSize: 20, Save: "old"}, Config{Save: "new"})
	if config.BufferSize != 20 {
		t.Error("buffer size not working")
	}
	if config.Save != "new" {
		if config.Save == "old" {
			t.Error("save not updated")
		} else {
			t.Error("save not working")
		}
	}
}
