package cmd

import (
	"fmt"
	"os"

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
		f, err := os.Create("template.html")
		if err != nil {
			fmt.Println("Error creating template file:", err)
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

		fmt.Println("Template file created successfully at template.html")
		fmt.Println("You can now edit the template file to your liking. Enjoy :)")
	},
}

func init() {
	rootCmd.AddCommand(createTemplateCmd)
}
