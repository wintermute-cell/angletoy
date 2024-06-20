package tool

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var (
	entityTemplateFile string
	sceneTemplateFile  string
	entitiesOutputDir  string
	scenesOutputDir    string
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create [some_name_in_snake_case]",
	Short: "Create a new entity or scene from a template",
	Long:  `Create a new file based on a template and provided name, placing it in the appropriate directory.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		someName := args[0]
		if cmd.Flags().Changed("entity") {
			processTemplate(someName, entityTemplateFile, entitiesOutputDir)
		} else if cmd.Flags().Changed("scene") {
			processTemplate(someName, sceneTemplateFile, scenesOutputDir)
		} else {
			fmt.Println("Please specify either --entity, --entity2d, or --scene")
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	// Define flags for template files and output directories
	createCmd.Flags().StringVar(&entityTemplateFile, "entity-template",
		"fw/core/entities/entity.template", "Template file for entity")
	createCmd.Flags().StringVar(&sceneTemplateFile, "scene-template",
		"fw/modules/scenes/scene.template", "Template file for scene")
	createCmd.Flags().StringVar(&entitiesOutputDir, "entities-output-dir",
		"game/entities", "Output directory for entities")
	createCmd.Flags().StringVar(&scenesOutputDir, "scenes-output-dir",
		"game/scenes", "Output directory for scenes")

	// Define flags for the type of file to create
	createCmd.Flags().BoolP("entity", "e", false, "Create an entity with the given name")
	createCmd.Flags().BoolP("scene", "s", false, "Create a scene with the given name")
}

func processTemplate(someName, templateFilePath, outputDir string) {

	// check if someName is valid
	if someName == "" {
		fmt.Println("Please provide a valid name")
		os.Exit(1)
	}
	if strings.Contains(someName, " ") {
		fmt.Println("Please provide a valid name without spaces")
		os.Exit(1)
	}
	if strings.Contains(someName, "-") {
		fmt.Println("Please provide a valid name without dashes")
		os.Exit(1)
	}
	if strings.Contains(someName, ".") {
		fmt.Println("Please provide a valid name without dots")
		os.Exit(1)
	}
	if strings.ContainsAny(someName, "ABCDEFGHIJKLMNOPQRSTUVWXYZ") {
		fmt.Println("Please provide a valid name without uppercase letters")
		os.Exit(1)
	}
	// if someName == "test" { // NOTE: this is not needed anymore, since we postfixed the files with _entity.go and _scene.go
	// 	fmt.Println("Please provide a valid name that is not 'test' (file name reserved by go)")
	// 	os.Exit(1)
	// }

	// Convert some_name to SomeName format
	someNameFormatted := toCamelCase(someName)

	// Read the template file
	templateContent, err := os.ReadFile(templateFilePath)
	if err != nil {
		fmt.Printf("Error reading template file: %v\n", err)
		os.Exit(1)
	}

	// Replace "Template" with the formatted name
	newContent := strings.ReplaceAll(string(templateContent), "Template", someNameFormatted)

	// Define the new file path
	newFilePath := filepath.Join(outputDir, fmt.Sprintf("%s_entity.go", someName))
	if strings.Contains(templateFilePath, "scene.template") {
		newFilePath = filepath.Join(outputDir, fmt.Sprintf("%s_scene.go", someName))
	}

	// Check if the file already exists
	if _, err := os.Stat(newFilePath); !os.IsNotExist(err) {
		fmt.Printf("File %s already exists, stopping\n", newFilePath)
		os.Exit(1)
	}

	// Write the new content to the new file
	err = os.WriteFile(newFilePath, []byte(newContent), 0644)
	if err != nil {
		fmt.Printf("Error writing new file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Processed template for %s and saved to %s\n", someNameFormatted, newFilePath)
}

func toCamelCase(input string) string {
	parts := strings.Split(input, "_")
	for i := range parts {
		parts[i] = strings.Title(parts[i])
	}
	return strings.Join(parts, "")
}
