package gt

import (
	"fmt"
	"html/template"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	_ "embed"
)

// greentext data struct
type GTData struct {
	Timestamp     string
	UniqueID      string
	Lines         []string
	Thumbnail     string
	ThumbnailSize int
	FontSize      int
	StyleBlock    template.HTML // holds dynamic CSS styles
	Template      []byte
}

func WriteToGreentext(dest string, tmpl []byte, lines []string, thumbnailPath, font string, fontSize int, previewOnly bool, bgColor, textColor string, width, height int, customDateTime string) error {

	var unixTime int64
	var timestamp string

	if customDateTime != "" {
		// Check input is just a date (DD/MM/YYYY)
		var parsedTime time.Time
		var err error

		if len(customDateTime) <= 10 && !strings.Contains(customDateTime, ":") {
			// It's just a date, generate a random time
			r := rand.New(rand.NewSource(time.Now().UnixNano()))
			randomHour := r.Intn(24)   // 0-23 hours
			randomMinute := r.Intn(60) // 0-59 minutes
			randomSecond := r.Intn(60) // 0-59 seconds

			// Format date with random time
			fullDateTime := fmt.Sprintf("%s, %02d:%02d:%02d",
				customDateTime,
				randomHour,
				randomMinute,
				randomSecond)

			parsedTime, err = time.Parse("02/01/2006, 15:04:05", fullDateTime)
			if err != nil {
				return fmt.Errorf("invalid date format: %v. Expected DD/MM/YYYY", err)
			}

			// Update timestamp to include random time
			timestamp = fullDateTime
		} else {
			parsedTime, err = time.Parse("02/01/2006, 15:04:05", customDateTime)
			if err != nil {
				return fmt.Errorf("invalid datetime format: %v. Expected DD/MM/YYYY, HH:MM:SS", err)
			}
			timestamp = customDateTime
		}

		unixTime = parsedTime.Unix()
	} else {
		unixTime = time.Now().Unix()
		timestamp = time.Unix(unixTime, 0).Format("02/01/2006, 15:04:05")
	}

	uniqueID := strconv.FormatInt(unixTime, 10)[:8]

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

	data.ThumbnailSize = getThumbnailSize(data.Thumbnail)

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

	fmt.Println("Generating greentext...")
	err = CaptureElementScreenshot(htmlFile, dest, width, height)
	if err != nil {
		return err
	}

	fmt.Println("Greentext generated and saved to", dest)
	return nil
}
