package mxtsession

import (
	"bufio"
	"bytes"
	"fmt"
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
func RenderSession(session map[string]string, sessionConfig config.SessionMap, wr *bufio.Writer) {

	tmpls := parseTmpl(sessionConfig)

	// FIXME:THis is ugly and should be dynamic (like it was before)

	tmpl, ok := tmpls["ssh"]

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
