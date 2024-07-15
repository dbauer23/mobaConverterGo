package config

import (
	"encoding/json"
	"log"
	"os"
)

// Options represents a map of string key-value pairs for options.
type Options map[string]string

// OptionsConfiguration holds the configuration for each option.
type OptionsConfiguration struct {
	Required *bool   `json:"required,omitempty"`
	Default  string  `json:"default_value"`
	Section  string  `json:"section"`
	Help     string  `json:"help"`
	Options  Options `json:"options,omitempty"`
}

// SessionConfiguration holds the configuration for each session.
type SessionConfiguration struct {
	TmplString     string   `json:"tmplString"`
	AllowedOptions []string `json:"allowed_options"`
	Options        Options  `json:"options,omitempty"`
}

// OptionsMap maps option names to their configurations.
type OptionsMap map[string]OptionsConfiguration

// SessionMap maps session types to their configurations.
type SessionMap map[string]SessionConfiguration

// Config represents the overall configuration with session types and options.
type Config struct {
	SessionTypes SessionMap        `json:"sessionTypes"`
	Options      OptionsMap        `json:"options"`
	Meta         map[string]string `json:"_meta"`
}

// JSONInput represents the structure of the new input JSON file.
type JSONInput struct {
	Meta      map[string]interface{}       `json:"_meta"`
	Sessions  []map[string]string          `json:"sessions"`
	Templates map[string]map[string]string `json:"templates"`
	Folders   map[string]map[string]string `json:"folders,omitempty"`
}

// LoadConfigurations read and unmarshal the JSON configuration file.
func LoadConfigurations(filename string) (OptionsMap, SessionMap, map[string]string) {

	var data []byte
	var err error

	if filename == "" {
		data, err = Asset("data/config.json")
		if err != nil {
			log.Fatalf("Error loading configurations: %v", err)
		}
	} else {
		data, err = os.ReadFile(filename)
		if err != nil {
			log.Fatalf("Error loading configurations: %v", err)
		}

	}

	var conf Config
	if err := json.Unmarshal(data, &conf); err != nil {
		log.Fatalf("Error loading configurations: %v", err)
	}

	return conf.Options, conf.SessionTypes, conf.Meta
}
