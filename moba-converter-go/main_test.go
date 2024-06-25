package main

import (
	"os"
	"testing"
)

func TestLoadConfigurations(t *testing.T) {
	// Prepare a sample JSON configuration
	configJSON := `
	{
		"session_types": {
			"example_session": {
				"session_type": "example",
				"tmplString": "{{.Example}}",
				"allowed_options": ["example_option"],
				"options": {
					"example_option": "example_value"
				}
			}
		},
		"options": {
			"example_option": {
				"default_value": "default",
				"section": "example_section",
				"help": "This is an example option",
				"options": {
					"old_value": "new_value"
				}
			}
		},
		"_meta": {
    		"version": "0.0.1",
   			 "changed_when": "1970-01-01"
 		 }
	}`

	// Write the JSON to a temporary file
	tempFile := "temp_config.json"
	if err := os.WriteFile(tempFile, []byte(configJSON), 0644); err != nil {
		t.Fatalf("Failed to write temp config file: %v", err)
	}
	defer os.Remove(tempFile)

	// Load configurations
	optionsMap, sessionMap, meta, err := LoadConfigurations(tempFile)
	if err != nil {
		t.Fatalf("Error loading configurations: %v", err)
	}

	// Check if configurations are loaded correctly
	if _, exists := optionsMap["example_option"]; !exists {
		t.Errorf("Expected example_option in optionsMap")
	}
	if _, exists := sessionMap["example_session"]; !exists {
		t.Errorf("Expected example_session in sessionMap")
	}
	if _, exists := meta["version"]; !exists {
		t.Errorf("Expected version in _meta")
	}
}

func TestApplyTemplate(t *testing.T) {
	sessionData := map[string]string{"key1": "value1"}
	templateData := map[string]string{"key2": "value2"}

	result := applyTemplate(sessionData, templateData)

	if result["key1"] != "value1" || result["key2"] != "value2" {
		t.Errorf("Template was not applied correctly")
	}
}

func TestSetDefaultValues(t *testing.T) {
	sessionData := map[string]string{"key1": "value1"}
	optionsMap := OptionsMap{
		"key1": {Default: "default1"},
		"key2": {Default: "default2"},
	}

	result := setDefaultValues(sessionData, optionsMap)

	if result["key1"] != "value1" {
		t.Errorf("Expected key1 to be value1, got %s", result["key1"])
	}
	if result["key2"] != "default2" {
		t.Errorf("Expected key2 to be default2, got %s", result["key2"])
	}
}

func TestApplyValueReplacements(t *testing.T) {
	sessionData := map[string]string{"key1": "old_value"}
	optionsMap := OptionsMap{
		"key1": {Options: Options{"old_value": "new_value"}},
	}

	result := applyValueReplacements(sessionData, optionsMap)

	if result["key1"] != "new_value" {
		t.Errorf("Expected key1 to be new_value, got %s", result["key1"])
	}
}

func TestParseTemplates(t *testing.T) {
	sessionMap := SessionMap{
		"example_session": {
			TmplString: "{{.Example}}",
		},
	}

	templates := parseTemplates(sessionMap)

	if templates["example_session"] == nil {
		t.Errorf("Expected template for example_session, got nil")
	}
}