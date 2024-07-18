package convert

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"moba-converter-go/internal/config"
	"moba-converter-go/internal/mxtsession"
	"moba-converter-go/internal/utils"
	"os"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

var outputPathM2J string
var inputPathM2J string

var jsonOutput config.JSONInput

var configPath string

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

	configPath, _ = cmd.Flags().GetString("configPath")
	optionsMap, sessionMap, _ := config.LoadConfigurations(configPath)

	reduce, _ := cmd.Flags().GetBool("reduce")

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

	regex_sessionType := regexp.MustCompile(`[^%]+#(?P<sessionType>\d)%`)

	regex_bookmark := regexp.MustCompile(`\[Bookmarks(_\d+)?\]`)
	regex_SubRep := regexp.MustCompile("SubRep=(.*)")
	regex_ImgNum := regexp.MustCompile(`ImgNum=(\d+)`)

	var isInBookmarkHeader bool // Defines if we are in a bookmark header. We the expect SubRep and ImgNum as the next lines
	var currentFolder string
	var currentImgNum string

	var sessionSlice []map[string]string

	// iterate over mxt line by line

	for i, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			// Skip empty lines
			continue
		}

		// Parse Bookmark lines
		if isInBookmarkHeader {
			if m := regex_SubRep.FindStringSubmatch(line); len(m) > 0 {
				currentFolder = strings.ReplaceAll(strings.Trim(m[1], "\r"), "\\", "/")
			} else if m := regex_ImgNum.FindStringSubmatch(line); len(m) > 0 {
				currentImgNum = m[1]
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
			c_session := evalSession(regex_sessionType.FindStringSubmatch(line)[1], line, i, currentImgNum, optionsMap, sessionMap)
			// Add folder to session
			if currentFolder != "" {
				// FIXME: This is ugly. There should always a folder exported . This currently doesn't work, since j2m doesn't accept "/" or ""
				c_session["folder"] = currentFolder
			}

			// Reduce vars  => only show vars which are non-default (if --reduce is set)
			if reduce {
				c_session = utils.ReduceOptions(c_session, optionsMap)
			}

			sessionSlice = append(sessionSlice, c_session)
		}
	}

	jsonOutput.Meta["description"] = "This file was created using moba-converter-go"
	jsonOutput.Sessions = sessionSlice

	xx, err := json.MarshalIndent(jsonOutput, "", "  ")
	if err != nil {
		log.Fatalln(err)
	}

	err = os.WriteFile(outputPathM2J, xx, 0644)
	if err != nil {
		log.Fatalln("Error writing to file")
	}

}

func evalSession(sessionType string, line string, lineNumber int, ImgNum string, optionsMap config.OptionsMap, sessionMap config.SessionMap) map[string]string {
	// TODO: Add error handling with lineNumber
	// TODO: Add Folder img support

	vars := make(map[string]string)
	tmpl := mxtsession.GetTmplBySessionTypeId(sessionType, optionsMap, sessionMap)

	// Remove the SessionName and add it directly to the vars.
	// This is done to remove "=" from the string to parse

	vars["SessionName"] = strings.Split(line, "=")[0]
	line = strings.Split(line, "=")[1]
	tmpl = strings.Split(tmpl, "=")[1]

	tmpl_slice := strings.Split(tmpl, "#")

	inColorScheme := false
	var optionName string
	// Split Parameter list by "#" to get Sections
	for i_section, section := range strings.Split(line, "#") {
		// Get tmpl for current section and split it
		tmpl_section := strings.Split(tmpl_slice[i_section], "%")
		// Split section by % to get Option
		// The offset is used in combination with Parameters which change the number of fields (%)(ColorScheme)
		offset := 0

		for i_option, option := range strings.Split(section, "%") {
			// fmt.Printf("\n\n")
			optionName = strings.Trim(tmpl_section[i_option+offset], "{. }")

			if i_option == 16 && option != "_Std_Colors_0_" && optionName == "ColorScheme" {
				inColorScheme = true
			}
			if inColorScheme {

				if i_option != 16 {
					option = "%" + option
					offset -= 1
				}
				vars[optionName] = vars[optionName] + option
				if i_option == 31 {
					// exit ColorScheme
					inColorScheme = false
				}
				continue
			}
			vars[optionName] = option
		}

	}

	vars = utils.ReverseValueReplacements(vars, optionsMap)

	return vars
}
