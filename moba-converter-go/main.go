package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"text/template"
)

// #region Global var definition

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
	SessionType    string   `json:"session_type"`
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
	SessionTypes SessionMap        `json:"session_types"`
	Options      OptionsMap        `json:"options"`
	Meta         map[string]string `json:"_meta"`
}

// JSONInput represents the structure of the new input JSON file.
type JSONInput struct {
	Meta      map[string]interface{}       `json:"_meta"`
	Sessions  []map[string]string          `json:"sessions"`
	Templates map[string]map[string]string `json:"templates"`
}

// #endregion

// LoadConfigurations read and unmarshal the JSON configuration file.
func LoadConfigurations(filename string) (OptionsMap, SessionMap, map[string]string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, nil, nil, err
	}

	var conf Config
	if err := json.Unmarshal(data, &conf); err != nil {
		return nil, nil, nil, err
	}

	return conf.Options, conf.SessionTypes, conf.Meta, nil
}

func printVersionInfo(meta map[string]string) {
	fmt.Println("Version of Value-Database")
	if version, ok := meta["version"]; ok {
		fmt.Printf("Version: %s\n", version)
	} else {
		fmt.Println("Version information not found.")
	}
	if changedWhen, ok := meta["changed_when"]; ok {
		fmt.Printf("Changed When: %s\n", changedWhen)
	} else {
		fmt.Println("Changed When information not found.")
	}
}

func printValueInfo(ops OptionsMap) {
	for key, option := range ops {
		fmt.Printf("Option: %s\n", key)
		fmt.Printf("  Section: %s\n", option.Section)
		fmt.Printf("  Default Value: %s\n", option.Default)
		fmt.Printf("  Help: %s\n", option.Help)
		if len(option.Options) > 0 {
			fmt.Println("  Possible Values:")
			for optKey := range option.Options {
				fmt.Printf("    - %s\n", optKey)
			}
		}
		fmt.Println()
	}
}

// #region Value handling

// applyTemplate applies template data to the session data.
func applyTemplate(sessionData, templateData map[string]string) map[string]string {
	for key, value := range templateData {
		sessionData[key] = value
	}
	return sessionData
}

// setDefaultValues sets default values for missing fields in session data.
func setDefaultValues(sessionData map[string]string, optionsMap OptionsMap) map[string]string {
	for key, valueSpec := range optionsMap {
		if _, exists := sessionData[key]; !exists {
			sessionData[key] = valueSpec.Default
		}
	}
	return sessionData
}

// applyValueReplacements replaces option values in session data.
func applyValueReplacements(sessionData map[string]string, optionsMap OptionsMap) map[string]string {
	for key, valueSpec := range optionsMap {
		if valueSpec.Options != nil {
			if val, exists := sessionData[key]; exists {
				if replacement, found := valueSpec.Options[val]; found {
					sessionData[key] = replacement
				}
			}
		}
	}
	return sessionData
}

// #endregion

// #region Template
// parseTemplates parses the templates from the session map.
func parseTemplates(sessionMap SessionMap) map[string]*template.Template {
	parsedTemplates := make(map[string]*template.Template)
	for key, value := range sessionMap {
		parsedTemplates[key] = template.Must(template.New(key).Parse(value.TmplString))
	}
	return parsedTemplates
}

// renderSession renders the session using the appropriate template.
func renderSession(session map[string]string, templates map[string]*template.Template) {
	tmpl, ok := templates[session["session_type"]]
	if !ok {
		if session["session_type"] == "" {
			fmt.Printf("Session type not supported: <NO SESSION TYPE SET>\n")
			return
		}
		fmt.Printf("Session type not supported: %s\n", session["session_type"])
		return
	}

	var rendered bytes.Buffer
	if err := tmpl.Execute(&rendered, session); err != nil {
		fmt.Printf("Error rendering template for type: %s, error: %v\n", session["session_type"], err)
		return
	}

	fmt.Println(rendered.String())
}

// #endregion

// #region Main Function
// main function
func main() {
	infoFlag := flag.Bool("info", false, "Prints the version and changed_when information from the config file.")
	valueInfo := flag.Bool("value-info", false, "List all possible options in a formatted manner")
	flag.Parse()

	optionsMap, sessionMap, meta, err := LoadConfigurations("config.json")
	if err != nil {
		log.Fatalf("Error loading configurations: %v", err)
	}

	if *infoFlag {
		printVersionInfo(meta)
		os.Exit(0)
	}

	if *valueInfo {
		printValueInfo(optionsMap)
		os.Exit(0)
	}

	file, err := os.ReadFile("input.json")
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}

	var input JSONInput
	if err := json.Unmarshal(file, &input); err != nil {
		log.Fatalf("Error parsing input file: %v", err)
	}

	sessionTemplates := input.Templates
	sessions := input.Sessions

	for i, session := range sessions {
		if templateName, hasTemplate := session["template"]; hasTemplate {
			fmt.Printf("Session %s uses template %s\n", session["sessionName"], templateName)
			session = applyTemplate(session, sessionTemplates[templateName])
		}

		session = setDefaultValues(session, optionsMap)
		session = applyValueReplacements(session, optionsMap)

		sessions[i] = session
	}

	parsedTemplates := parseTemplates(sessionMap)

	for _, session := range sessions {
		renderSession(session, parsedTemplates)
	}
}

// #endregion
