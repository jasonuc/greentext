package gt

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strconv"
	"time"

	_ "embed"
)

// greentext data struct
type GTData struct {
	Timestamp  string
	UniqueID   string
	Lines      []string
	Thumbnail  string
	FontSize   int
	StyleBlock template.HTML // holds dynamic CSS styles
	Template   []byte
}

func WriteToGreentext(dest string, tmpl []byte, lines []string, thumbnailPath, font string, fontSize int, previewOnly bool, bgColor, textColor string) error {

	unixTime := time.Now().Unix()
	timestamp := time.Unix(unixTime, 0).Format("02/01/2006, 15:04:05")
	uniqueID := strconv.FormatInt(unixTime, 10)[:8] // For the part that says no.xxxxxxxx
	styleBlock := template.HTML(fmt.Sprintf(`
    <style>
        body {
            font-family: %s, 'Roboto Mono', monospace;
        }

        .container {
            background-color: %s;
        }

        .text small {
            font-size: %dpx;
            color: %s;
        }
    </style>`, font, bgColor, fontSize, textColor)) // extra style block for dynamic styles like font, background colour, text colour...

	data := GTData{
		Timestamp:  timestamp,
		UniqueID:   uniqueID,
		Lines:      lines,
		Thumbnail:  thumbnailPath,
		FontSize:   fontSize,
		StyleBlock: styleBlock,
		Template:   tmpl,
	}

	htmlFile := "temp_meme.html"
	err := GenerateHTMLFile(htmlFile, data)
	if err != nil {
		return err
	}
	defer os.Remove(htmlFile)

	if previewOnly {
		fmt.Println("Preview mode enabled. Opening browser...")
		var absHtmlFile string
		if absHtmlFile, err = filepath.Abs(htmlFile); err != nil {
			return fmt.Errorf("failed to get absolute path for HTML file: %w", err)
		}

		if err := openBrowser("file://" + absHtmlFile); err != nil {
			return fmt.Errorf("failed to open browser: %w", err)
		}

		fmt.Println("Press Enter to exit after viewing the preview...")
		_, _ = fmt.Scanln() // Wait for user input before exiting
		fmt.Println("Exiting...")
		return nil
	}

	err = CaptureElementScreenshot(htmlFile, dest)
	if err != nil {
		return err
	}

	fmt.Println("Greentext generated and saved to", dest)
	return nil
}
