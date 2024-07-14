package convert

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"moba-converter-go/internal/config"
	"os"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

var outputPathM2J string
var inputPathM2J string

var jsonOutput config.JSONInput

func init() {
	convertCmd.AddCommand(moba2jsonCmd)
	moba2jsonCmd.Flags().StringVar(&outputPathM2J, "output", "converted.json", "Path to write the output json file.")
	moba2jsonCmd.Flags().StringVar(&inputPathM2J, "input", "", "Path to input mxtsessions file. If not set, reads from stdin.")
	moba2jsonCmd.Flags().BoolP("reduce", "r", true, "Only export non-default Parameters to the json file")

	jsonOutput.Meta = make(map[string]interface{})
	jsonOutput.Templates = make(map[string]map[string]string)
}

var moba2jsonCmd = &cobra.Command{
	Use:     "moba2json",
	Short:   "Convert mobaxterm session file into json",
	Long:    "",
	Aliases: []string{"m2j"},
	Run:     convertMoba2Json,
}

func convertMoba2Json(cmd *cobra.Command, args []string) {

	// Read moba file
	var data []byte
	var err error
	if inputPathM2J != "" {
		data, err = os.ReadFile(inputPathM2J)
		if err != nil {
			log.Fatalf("Error opening file: %v", err)
		}
	} else {
		// Read from stdin
		fmt.Fprintln(os.Stderr, "<<Reading from stdin.>>")
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			data = append(data, scanner.Bytes()...)
		}
		if err := scanner.Err(); err != nil {
			log.Fatalf("Error reading from stdin: %v", err)
		}
	}

	regex_sessionType := regexp.MustCompile(`[^%]+#(\d)%`)

	regex_bookmark := regexp.MustCompile(`\[Bookmarks(_\d+)?\]`)
	regex_SubRep := regexp.MustCompile("SubRep=(.*)")
	regex_ImgNum := regexp.MustCompile(`ImgNum=\d+`)

	var isInBookmarkHeader bool // defines if we are in a bookmark header. We the expect SubRep and ImgNum as the next lines
	var currentFolder string
	var currentImgNum string

	var sessionSlice []map[string]string

	// iterate over mxt line by line

	for i, line := range strings.Split(string(data), "\n") {
		fmt.Println("Working on line ", i)
		if line == "" {
			// Skip empty lines
			continue
		}

		// Parse Bookmark lines
		if isInBookmarkHeader {
			// TODO: Use the folder info and add it to the json output
			if m := regex_SubRep.FindString(line); m != "" {
				currentFolder = m
			} else if m := regex_ImgNum.FindString(line); m != "" {
				currentImgNum = m
				isInBookmarkHeader = false
			} else {
				log.Fatalf("Expected SubRep or ImgNum line. Got %s in line %d", line, i)
			}
			continue
		}

		if regex_bookmark.MatchString(line) {
			// Set Folder
			isInBookmarkHeader = true
		} else {
			// Now this should be a session line
			sessionSlice = append(sessionSlice, evalSession(regex_sessionType.FindString(line), line, i, currentFolder, currentImgNum))
		}
	}
	// Finished session section
	fmt.Print(sessionSlice)

	// Build json

	fmt.Println("")

	jsonOutput.Meta["description"] = "This file was created using moba-converter-go"
	jsonOutput.Sessions = sessionSlice

	xx, err := json.MarshalIndent(jsonOutput, "", "  ")
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(string(xx))
	fmt.Println(outputPathM2J)

	err = os.WriteFile(outputPathM2J, xx, 0644)
	if err != nil {
		log.Fatalln("Error writing to file")
	}

}

func evalSession(sessionType string, line string, lineNumber int, folder string, ImgNum string) map[string]string {

	// TODO: use the session type and set tmpl
	tmpl := "{{.SessionName}}={{ .Logout }}#{{ .IconNumber }}#{{ .SessionType }}%{{ .RemoteHost }}%{{ .Port }}%{{ .Username }}%%{{ .X11Forwarding }}%{{ .Compression }}%{{ .ExecuteCommand }}%{{ .SSHGWHostList }}%{{ .SSHGWPortList }}%{{ .SSHGWUserList }}%{{ .DoNotExitAfterLoginCommand }}%{{ .DontSpecifyUsername }}%{{ .RemoteEnvironment }}%{{ .PrivateKeyPath }}%{{ .SSHGWPrivateKeyList }}%{{ .SSHBrowserType }}%{{ .FollowSSHPath }}%0%{{ .ProxyType }}%{{ .ProxyHost }}%{{ .ProxyPort }}%{{ .ProxyLogin }}%{{ .AdaptLocales }}%{{ .FileBrowserSCPOverSFTP }}%{{ .FileBrowserProtocol }}%{{ .LocalProxyCommand }}%{{ .SSHProtocolVersion }}%{{ .KeyExchangeAlgos }}%{{ .HostKeyTypes }}%{{ .Ciphers }}%{{ .DisconnectIfAuthSucceedsTrivially }}%{{ .PreferHostKeyAlgorithms }}%{{ .AttemptAuthUsingSSHAgent }}%{{ .AllowAgentForwarding }}#{{ .TerminalFont }}%{{ .FontSize }}%{{ .TerminalFontBold }}%%{{ .AppendWindowsPath }}%{{ .TerminalCharset }}%{{ .Foreground }}%{{ .Background }}%{{ .CursorColor }}%{{ .CursorType }}%{{ .BackspaceSendsCtrlH }}%{{ .LogOutput }}%{{ .LogFolderPath }}%{{ .TerminalType }}%{{ .LockTerminalTitle }}%%{{ .ColorScheme }}%{{ .TerminalRows }}%{{ .TerminalColumns }}%{{ .ForceFixedRowsCols }}%{{ .SyntaxHighlighting }}%{{ .ShowBoldFontAsBrighter }}%{{ .CustomMacroType }}%{{ .CustomMacroText }}%{{ .PasteDelay }}%{{ .FontCharset }}%{{ .FontAntialiasing }}%{{ .FontLigatures }}#{{ .StartSessionIn }}#{{ .Comments }}#{{ .CustomTabColor }}"

	vars, err := extractVariables(tmpl, line)
	if err != nil {
		log.Fatalf("Error in Line %d: %s", lineNumber, err)
	}

	// TODO: convert session type to key

	vars["SessionType"] = "ssh"

	// TODO: Convert all vars to their text form if possible

	// TODO: Reduce vars  => only show vars which are non-default (if --reduce is set)

	return vars
}

func createRegexFromTemplate(template string) (string, error) {
	// The placeholder regex pattern to match {{.VariableName}}
	placeholderRegex := regexp.MustCompile(`\\{\\{\s*\\.([a-zA-Z0-9_]+)\s*\\}\\}`)
	// Replace placeholders with named capturing groups
	regexStr := placeholderRegex.ReplaceAllString(template, `(?P<$1>.*?)`)

	// Append $ to make sure the regex is read until EOL
	regexStr += "$"

	return regexStr, nil
}

func extractVariables(template, input string) (map[string]string, error) {
	escapedTemplate := regexp.QuoteMeta(template)
	regexStr, err := createRegexFromTemplate(escapedTemplate)
	if err != nil {
		return nil, err
	}

	regex, err := regexp.Compile(regexStr)
	if err != nil {
		return nil, err
	}

	// fmt.Println(input)
	// fmt.Println(regexStr)

	match := regex.FindStringSubmatch(input)
	if match == nil {
		return nil, fmt.Errorf("input does not match template")
	}

	vars := make(map[string]string)
	for i, name := range regex.SubexpNames() {
		if i != 0 && name != "" {
			vars[name] = match[i]
		}
	}

	return vars, nil
}
