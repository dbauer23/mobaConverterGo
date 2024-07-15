package mxtsession

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"moba-converter-go/internal/config"
	"os"
	"text/template"
)

// parseTmpl parses the tmpl strings from the session map.
func parseTmpl(sessionMap config.SessionMap) map[string]*template.Template {
	parsedTmpl := make(map[string]*template.Template)
	for key, value := range sessionMap {
		parsedTmpl[key] = template.Must(template.New(key).Parse(value.TmplString))
	}
	return parsedTmpl
}

// renderSession renders the session using the appropriate tmpl string.
func RenderSession(session map[string]string, sessionConfig config.SessionMap, optionsMap config.OptionsMap, wr *bufio.Writer) {

	tmpls := parseTmpl(sessionConfig)
	tmpl, ok := tmpls[getSessionTypeById(session["SessionType"], optionsMap)]

	if !ok {
		if session["SessionType"] == "" {
			fmt.Fprintf(os.Stderr, "Session type not supported: <NO SESSION TYPE SET> in session '%s'\n", session["SessionName"])
			return
		}
		fmt.Fprintf(os.Stderr, "Session type not supported: %s in session '%s'\n", session["SessionType"], session["SessionName"])
		return
	}

	var rendered bytes.Buffer
	if err := tmpl.Execute(&rendered, session); err != nil {
		fmt.Fprintf(os.Stderr, "Error rendering tmpl for type: %s, error: %v\n", session["SessionType"], err)
		return
	}

	fmt.Fprintf(wr, "%s\r\n", rendered.String())
}

func getSessionTypeById(id string, optionsMap config.OptionsMap) string {
	// Return the sessionType string for a specific id (0 => ssh)
	sessionTypes := optionsMap["SessionType"].Options

	for k, v := range sessionTypes {
		if v == id {
			return k
		}
	}

	log.Fatalln("Invalid SessionType")
	return ""

}

func getTmplBySessionTypeName(sessionType string, sessionMap config.SessionMap) string {
	for key, v := range sessionMap {
		if key == sessionType {
			return v.TmplString
		}
	}
	// TODO: Check if this should be fatal

	log.Fatalln("unknown session")
	return ""
}

func GetTmplBySessionTypeId(sessionTypeId string, optionsMap config.OptionsMap, sessionMap config.SessionMap) string {
	return getTmplBySessionTypeName(getSessionTypeById(sessionTypeId, optionsMap), sessionMap)
}
