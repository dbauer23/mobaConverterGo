package convert

import (
	"bufio"
	"fmt"
	"moba-converter-go/internal/config"
	"moba-converter-go/internal/io"
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
	// Load configuration
	configPath, _ := cmd.Flags().GetString("configPath")
	optionsMap, sessionMap, _ := config.LoadConfigurations(configPath)
	// Load Input File
	input := io.LoadJsonInput(&inputPath)

	// Create a new sessions variable to store the final session is.
	// After applying defaults, templates etc.
	var sessions []map[string]string

	// Iterate over input sessions
	for _, session := range input.Sessions {
		//  - Apply templates if there are any.
		if templateName, hasTemplate := session["template"]; hasTemplate {
			fmt.Fprintf(os.Stderr, "Session %s uses template %s\n", session["SessionName"], templateName)
			session = utils.ApplyTemplate(session, input.Templates[templateName])
		}

		//	- Set defaults
		session = utils.SetDefaultValues(session, optionsMap)
		// Apply replacements. Until now the session would contain "true / false " instead of the required "0 / -1"
		session = utils.ApplyValueReplacements(session, optionsMap)

		// Store the session in the final session var
		sessions = append(sessions, session)
	}

	// Group sessions by folder
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
			if f, exists := input.Folders[currentFolder]; exists {
				imageNum = f["Icon"]
			}

			// Print new folder heading
			// Also replace "/" with "\" since moba needs backslashes...
			clearedFolderPath := strings.TrimLeft(strings.ReplaceAll(currentFolder, "/", "\\"), "\\")
			fmt.Fprintf(writer, "[Bookmarks_%d]\r\nSubRep=%s\r\nImgNum=%s\r\n", idx, clearedFolderPath, imageNum)
			utils.Check(err)
		}

		for _, session := range sessions {
			mxtsession.RenderSession(session, sessionMap, optionsMap, writer)
		}

	}
	writer.Flush()
}
