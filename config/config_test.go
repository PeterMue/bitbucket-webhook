package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseConfig(t *testing.T) {
	cases := []struct {
		args   []string
		expect *Config
	}{
		{[]string{"program", "-config", "config_test.yaml"}, &Config{
			Listen: ":3000",
			Secret: "config-secret",
			Hooks: []struct {
				EventType string "yaml:\"event\""
				Command   string "yaml:\"command\""
				Args []string "yaml:\"args\""
			}{
				{"dummy:event", "bash", nil},
			},
		}},
		{[]string{"program", "-config", "config_test.yaml", "-listen", ":1337", "-secret", "flag-secret"}, &Config{
			Listen: ":1337",
			Secret: "flag-secret",
			Hooks: []struct {
				EventType string "yaml:\"event\""
				Command   string "yaml:\"command\""
				Args []string "yaml:\"args\""
			}{
				{"dummy:event", "bash", nil},
			},
		}},
	}

	for _, c := range cases {

		config, err := ParseFlags(c.args)
		if err != nil {
			t.Errorf("Parse flags failed: %s", err)
		}

		assert.Equal(t, config, c.expect)
	}

}
