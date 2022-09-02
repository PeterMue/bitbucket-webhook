package handler

import (
	"bytes"
	"encoding/json"
	"log"
	"os/exec"
	"text/template"

	"github.com/PeterMue/bitbucket-webhook/header"
)

type Handler struct {
	event   string
	command string
	args    []string
	asnyc   bool
}

func New(event, command string, args []string, async bool) *Handler {
	return &Handler{
		event, command, args, async,
	}
}

func (h *Handler) templateArgs(message []byte) ([]string, error) {
	m := map[string]interface{}{}
	if err := json.Unmarshal(message, &m); err != nil {
		return nil, err
	}

	args := make([]string, len(h.args))
	for i, arg := range h.args {
		template, err := template.New("").Parse(arg)
		if err != nil {
			return nil, err
		}
		var b bytes.Buffer
		if err := template.Execute(&b, m); err != nil {
			return nil, err
		}
		args[i] = b.String()
	}

	return args, nil
}

func (h *Handler) run(headers *header.Headers, message []byte) error {
	args, err := h.templateArgs(message)
	if err != nil {
		return err
	}

	cmd := exec.Command(h.command, args...)
	
	if out, err := cmd.Output(); err != nil {
		log.Printf("Error[requestId=%s, event=%s, err=%s, output=%s]", headers.EventKey, headers.RequestId, err, out)
	} else {
		log.Printf("Done[requestId=%s, event=%s, output=%s]", headers.EventKey, headers.RequestId, out)
	}

	return nil
}

func (h *Handler) Run(headers *header.Headers, message []byte) error {
	if h.asnyc {
		log.Printf("Accepted[requestId=%s, event=%s]", headers.EventKey, headers.RequestId)
		go func() {
			h.run(headers, message)
		}()
		return nil
	}
	return h.run(headers, message)
}
