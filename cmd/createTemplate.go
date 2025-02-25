package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var createTemplateCmd = &cobra.Command{
	Use:   "create-template",
	Short: "Initializes a new template file using the default template, that you can work on",
	Long: `Initializes a new template file, that you can work on. This is for those who want to create a new HTML template from scratch.
The supported variables are:
- {{.StyleBlock}}: The dynamic style block
- {{.Timestamp}}: The current timestamp
- {{.UniqueID}}: A unique identifier derived from the timestamp
- {{.Thumbnail}}: The thumbnail image
- {{.ThumbnailSize}}: The size of the thumbnail image
- {{.Lines}}: The array of lines

If there are any other variables you would like to see, please open an issue or a pull request the GitHub repository: jasonuc/greentext.
If you want your template to be included in the default templates, please open a pull request on the GitHub repository: jasonuc/greentext.`,
	Run: func(cmd *cobra.Command, args []string) {
		templateFileName, _ := cmd.Flags().GetString("output")
		if templateFileName == "" {
			templateFileName = "template.html"
		}

		// Check if file already exists
		if _, err := os.Stat(templateFileName); err == nil {
			fmt.Printf("File %s already exists. Do you want to overwrite it? (y/n): ", templateFileName)
			var response string
			fmt.Scanln(&response)
			if response != "y" && response != "Y" {
				fmt.Println("Operation cancelled")
				return
			}
		}

		f, err := os.Create(templateFileName)
		if err != nil {
			fmt.Println("Error creating template file:", err)
			return
		}

		defer f.Close()

		defaultTemplate, ok := cmd.Context().Value(defaultTemplateKey).([]byte)
		if !ok {
			fmt.Println("Invalid template passed")
			return
		}

		if _, err = f.Write(defaultTemplate); err != nil {
			fmt.Println("Error writing to template file:", err)
			return
		}

		absPath, _ := filepath.Abs(templateFileName)
		fmt.Println("Template file created successfully at", absPath)
		fmt.Println("You can now edit the template file to your liking. Enjoy :)")
		fmt.Println("\nTo use this template with the greentext command, run:")
		fmt.Printf("  greentext --tmpl %s [other options]\n", templateFileName)
	},
}

func init() {
	createTemplateCmd.Flags().StringP("output", "o", "template.html", "Name of the template file to create")
	rootCmd.AddCommand(createTemplateCmd)
}
