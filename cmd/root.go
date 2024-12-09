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
	Long:  `Generate greentext memes from a given text.`,
	Run: func(cmd *cobra.Command, args []string) {
		lineCount, err := cmd.Flags().GetInt("lines")
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Generating greentext with", lineCount, "lines")

		lines, err := pkg.ReadInputLines(lineCount)
		if err != nil {
			fmt.Println(err)
			return
		}

		dest, err := cmd.Flags().GetString("output")
		if err != nil {
			fmt.Println(err)
			return
		}

		thumbnail, err := cmd.Flags().GetString("thumbnail")
		if err != nil {
			fmt.Println(err)
			return
		}

		if err := pkg.WriteToImage(dest, lines, thumbnail); err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("Greentext generated and saved to", dest)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().IntP("lines", "l", 5, "Number of lines to include in the greentext. Default is 5")
	rootCmd.Flags().StringP("output", "o", "greentext-output.png", "Output file to write the greentext to. Default is greentext-output.png")
	rootCmd.Flags().StringP("thumbnail", "t", "", "Thumbnail to use for the greentext. Default is no thumbnail. Only supports png files at the moment")
}
