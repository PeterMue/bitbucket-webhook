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
	event string
	command string
	args []string
}

func New(event, command string, args []string) *Handler {
	return &Handler{
		event, command, args,
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

func (h *Handler) Run(headers *header.Headers, message []byte) error {
	
    args, err := h.templateArgs(message)
	if err != nil {
		return err
	}

	log.Printf("Running event=%s, requestId=%s\n", headers.EventKey, headers.RequestId )

	cmd := exec.Command(h.command, args...)
	if _, err := cmd.Output(); err != nil {
		return err
	}
	
	return nil
}