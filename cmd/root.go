package cmd

import (
	"fmt"
	"os"

	"github.com/jasonuc/greentext/pkg"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "greentext",
	Short: "Generate greentext memes",
	Long: `Generate greentext memes.
	
Created by github.com/jasonuc.
Visit https://github.com/jasonuc/greentext for more information.`,
	Run: func(cmd *cobra.Command, args []string) {
		lineCount, err := cmd.Flags().GetInt("lines")
		if err != nil {
			fmt.Println("Error reading line count:", err)
			return
		}
		fmt.Println("Generating greentext with", lineCount, "lines")

		lines, err := pkg.ReadInputLines(lineCount)
		if err != nil {
			fmt.Println("Error reading input lines:", err)
			return
		}

		dest, err := cmd.Flags().GetString("output")
		if err != nil {
			fmt.Println("Error reading output flag:", err)
			return
		}

		thumbnail, err := cmd.Flags().GetString("thumbnail")
		if err != nil {
			fmt.Println("Error reading thumbnail flag:", err)
			return
		}

		templatePath := "templates/greentext_template.html"

		fontSize, err := cmd.Flags().GetInt("font-size")

		if fontSize < 8 || fontSize > 100 {
			fmt.Println("Error: Font size must be between 8 and 100.")
			return
		}

		if err != nil {
			fmt.Println("Error reading font size flag:", err)
			return
		}

		font, err := cmd.Flags().GetString("font")
		if err != nil {
			fmt.Println("Error reading font flag:", err)
			return
		}

		previewOnly, err := cmd.Flags().GetBool("preview-only")
		if err != nil {
			fmt.Println("Error reading preview flag:", err)
			return
		}

		err = pkg.WriteToMemeImage(dest, lines, thumbnail, templatePath, font, fontSize, previewOnly)
		if err != nil {
			fmt.Println("Error generating greentext meme:", err)
			return
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().IntP("lines", "l", 5, "Number of lines to include in the greentext")
	rootCmd.Flags().StringP("output", "o", "greentext.png", "Output file for the greentext. Supports PNG (default) and other formats based on the file extension.")
	rootCmd.Flags().StringP("thumbnail", "t", "", "Thumbnail to use for the greentext. Default is no thumbnail. Supports image file paths or URLs. Example: /path/to/image.png or https://example.com/image.png.")
	rootCmd.Flags().IntP("font-size", "s", 12, "Font size for the greentext lines.")
	rootCmd.Flags().StringP("font", "f", "Roboto Mono", "Font family to use for the entire greentext meme. Only supports built-in web-safe fonts (e.g., 'Courier New', 'Comic Sans MS', 'Monaco')")
	rootCmd.Flags().BoolP("preview-only", "P", false, "Preview the greentext in the browser without generating an image.")
}
