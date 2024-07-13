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
	"strings"

	"github.com/spf13/cobra"
)

var outputPath string
var inputPath string

func init() {
	convertCmd.AddCommand(json2mobaCmd)
	json2mobaCmd.Flags().StringVar(&outputPath, "output", "converted.mxtsessions", "Path to write the output mxtsessions file ")
	json2mobaCmd.Flags().StringVar(&inputPath, "input", "", "Path to input JSON file. If not set, reads from stdin.")
}

var json2mobaCmd = &cobra.Command{
	Use:     "json2moba",
	Short:   "Convert a json file into a mobaxterm session file",
	Long:    "",
	Aliases: []string{"j2m"},
	Run:     convertJson2Moba,
}

func convertJson2Moba(cmd *cobra.Command, args []string) {
	configPath, _ := cmd.Flags().GetString("configPath")
	optionsMap, sessionMap, _, err := config.LoadConfigurations(configPath)
	if err != nil {
		log.Fatalf("Error loading configurations: %v", err)
	}

	var data []byte

	// Read from file if inputPath is provided
	if inputPath != "" {
		data, err = os.ReadFile(inputPath)
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

	var input config.JSONInput
	if err := json.Unmarshal(data, &input); err != nil {
		log.Fatalf("Error parsing input file: %v", err)
	}

	sessionTemplates := input.Templates
	sessions := input.Sessions
	folders := input.Folders

	// Iterate over input sessions and
	//  - Apply templates
	//	- Set defaults
	//  - replace values

	for i, session := range sessions {
		if templateName, hasTemplate := session["template"]; hasTemplate {
			fmt.Fprintf(os.Stderr, "Session %s uses template %s\n", session["SessionName"], templateName)
			session = utils.ApplyTemplate(session, sessionTemplates[templateName])
		}

		session = utils.SetDefaultValues(session, optionsMap)
		session = utils.ApplyValueReplacements(session, optionsMap)

		sessions[i] = session
	}

	// Group list by folder
	groupedSessions := utils.GroupByFolder(sessions)

	// Open output file
	f, err := os.Create(outputPath)
	utils.Check(err)
	defer f.Close()

	writer := bufio.NewWriter(f)

	fmt.Fprintf(writer, "[Bookmarks]\r\nSubRep=\r\nImgNum=42\r\n")
	idx := 0
	for currentFolder, sessions := range groupedSessions {
		if currentFolder != "/" {
			idx++
			// Empty line after each Folder Block
			fmt.Fprintln(writer, "")
			utils.Check(err)
			// Default Image Number
			imageNum := "41"
			//  Set new Image number if folder is specified
			if _, exists := folders[currentFolder]; exists {
				imageNum = folders[currentFolder]["Icon"]
			}

			// Print new folder heading
			// Also replace "/" with "\" since moba needs backslashes...
			clearedFolderPath := strings.TrimLeft(strings.ReplaceAll(currentFolder, "/", "\\"), "\\")
			fmt.Fprintf(writer, "[Bookmarks_%d]\r\nSubRep=%s\r\nImgNum=%s\r\n", idx, clearedFolderPath, imageNum)
			utils.Check(err)
		}

		for _, session := range sessions {
			mxtsession.RenderSession(session, mxtsession.ParseTmpl(sessionMap), writer)
		}

	}
	writer.Flush()
}
