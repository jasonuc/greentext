package cmd

import (
	"context"
	"fmt"

	_ "embed"

	"github.com/jasonuc/greentext/pkg/gt"
	"github.com/jasonuc/greentext/pkg/version"
	"github.com/spf13/cobra"
)

type ctxKey string

const (
	defaultTemplateKey ctxKey = "DEFAULT"
)

var rootCmd = &cobra.Command{
	Use:   "greentext",
	Short: "Generate greentext memes",
	Long: `Generate greentext memes.
	
Created by github.com/jasonuc.
Visit https://github.com/jasonuc/greentext for more information.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Helper functions
		getStringFlag := func(name string) (string, error) {
			return cmd.Flags().GetString(name)
		}
		getIntFlag := func(name string) (int, error) {
			return cmd.Flags().GetInt(name)
		}

		defaultTemplate, ok := cmd.Context().Value(defaultTemplateKey).([]byte)
		if !ok {
			fmt.Println("Invalid template passed")
			return
		}

		// Read flags
		lineCount, err := getIntFlag("lines")
		if err != nil {
			fmt.Println("Error reading line count:", err)
			return
		}

		textFile, err := getStringFlag("input-file")
		if err != nil {
			fmt.Println("Error reading file flag:", err)
			return
		}

		var lines []string
		if textFile != "" {
			// Read lines from file using the file from the -i flag
			lines, err = gt.ReadLinesFromFile(textFile)
			if err != nil {
				fmt.Println("Error reading lines from file:", err)
				return
			}
		} else {
			// Read lines interactively using the no. of lines from the -l flag
			fmt.Println("Generating greentext with", lineCount, "lines")
			lines, err = gt.ReadInputLines(lineCount)
			if err != nil {
				fmt.Println("Error reading input lines:", err)
				return
			}
		}

		dest, err := getStringFlag("output")
		if err != nil {
			fmt.Println("Error reading output flag:", err)
			return
		}

		thumbnail, err := getStringFlag("thumbnail")
		if err != nil {
			fmt.Println("Error reading thumbnail flag:", err)
			return
		}

		fontSize, err := getIntFlag("font-size")
		if err != nil {
			fmt.Println("Error reading font size flag:", err)
			return
		}

		// Validate font size
		if fontSize < 8 || fontSize > 100 {
			fmt.Println("Error: Font size must be between 8 and 100.")
			return
		}

		font, err := getStringFlag("font")
		if err != nil {
			fmt.Println("Error reading font flag:", err)
			return
		}

		previewOnly, err := cmd.Flags().GetBool("preview-only")
		if err != nil {
			fmt.Println("Error reading preview flag:", err)
			return
		}

		bgColor, err := getStringFlag("background-color")
		if err != nil {
			fmt.Println("Error reading background color flag:", err)
			return
		}

		textColor, err := getStringFlag("text-color")
		if err != nil {
			fmt.Println("Error reading text color flag:", err)
			return
		}

		width, err := getIntFlag("width")
		if err != nil {
			fmt.Println("Error reading text width flag:", err)
			return
		}

		height, err := getIntFlag("height")
		if err != nil {
			fmt.Println("Error reading text height flag:", err)
			return
		}

		err = gt.WriteToGreentext(dest, defaultTemplate, lines, thumbnail, font, fontSize, previewOnly, bgColor, textColor, width, height)
		if err != nil {
			fmt.Println("Error generating greentext:", err)
			return
		}
	},
}

func Execute(currentVersion string, defaultTemplate []byte) error {
	rootCmd.Version = currentVersion
	info := version.FetchUpdateInfo(rootCmd.Version)
	defer info.PromptUpdateIfAvailable()
	ctx := version.WithContext(context.Background(), &info)

	ctx = context.WithValue(ctx, defaultTemplateKey, defaultTemplate)
	return rootCmd.ExecuteContext(ctx)
}

func init() {
	rootCmd.Flags().IntP("lines", "l", 5, "Number of lines to include in the greentext")
	rootCmd.Flags().StringP("output", "o", "greentext.png", "Output file for the greentext. Supports PNG (default) and other formats based on the file extension")
	rootCmd.Flags().StringP("thumbnail", "t", "", "Thumbnail to use for the greentext. Default is no thumbnail. Supports image file paths or URLs")
	rootCmd.Flags().IntP("font-size", "s", 12, "Font size for the greentext lines")
	rootCmd.Flags().StringP("font", "f", "Roboto Mono", "Font family to use for the entire greentext meme. Only supports built-in web-safe fonts")
	rootCmd.Flags().BoolP("preview-only", "P", false, "Preview the greentext in the browser without generating an image")
	rootCmd.Flags().StringP("background-color", "b", "#f0e0d6", "Background color for the greentext meme in HEX format")
	rootCmd.Flags().StringP("text-color", "c", "#819f32", "Text color for the greentext lines in HEX format")
	rootCmd.Flags().StringP("input-file", "i", "", "Path to a text file containing the greentext lines. Overrides the --lines flag")
	rootCmd.Flags().String("tmpl", "", "Path to a custom template file for the greentext")
	rootCmd.Flags().Int("width", 0, "Width of the generated greentext image in pixels. Default is auto")
	rootCmd.Flags().Int("height", 0, "Height of the generated greentext image in pixels. Default is auto")
}
