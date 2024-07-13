package utils

import (
	"moba-converter-go/internal/config"
)

// #region Value Handling

// applyTemplate applies template data to the session data.
func ApplyTemplate(sessionData, templateData map[string]string) map[string]string {
	for key, value := range templateData {
		sessionData[key] = value
	}
	return sessionData
}

// setDefaultValues sets default values for missing fields in session data.
func SetDefaultValues(sessionData map[string]string, optionsMap config.OptionsMap) map[string]string {
	for key, valueSpec := range optionsMap {
		if _, exists := sessionData[key]; !exists {
			sessionData[key] = valueSpec.Default
		}
	}
	// Custom default for folder
	if _, exists := sessionData["folder"]; !exists {
		sessionData["folder"] = "/"
	}
	return sessionData
}

// applyValueReplacements replaces option values in session data.
func ApplyValueReplacements(sessionData map[string]string, optionsMap config.OptionsMap) map[string]string {
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

// #region Other

func GroupByFolder(slice []map[string]string) map[string][]map[string]string {
	groupKey := "folder"
	grouped := make(map[string][]map[string]string)
	for _, item := range slice {
		if value, exists := item[groupKey]; exists {
			grouped[value] = append(grouped[value], item)
		}
	}
	return grouped
}

func Check(err error) {
	if err != nil {
		panic(err)
	}
}
