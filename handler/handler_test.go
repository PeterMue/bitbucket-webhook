package handler

import (
	"testing"

	"github.com/PeterMue/bitbucket-webhook/header"
)

func TestRun(t *testing.T) {
	json := []byte(`{ "actor" : { "name" : "Johnny" } }`)

	h := New("dummy", "/bin/bash", []string{"-c", "echo \"Halo {{ .actor.name }}\""}, false)
	if err := h.Run(&header.Headers{}, json); err != nil {
		t.Errorf("Run failed: %s", err)
	}

}
