package pkg

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/nfnt/resize"
)

type MemeData struct {
	Timestamp  string
	UniqueID   string
	Lines      []string
	Thumbnail  string
	FontSize   int
	StyleBlock template.HTML // holds dynamic CSS styles
}

func GenerateHTMLFile(outputFile string, data MemeData, templatePath string) error {

	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	outFile, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	err = tmpl.Execute(outFile, data)
	if err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return nil
}

func CaptureElementScreenshot(htmlFile, outputImageFile string) error {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	absPath, err := filepath.Abs(htmlFile)
	if err != nil {
		return fmt.Errorf("failed to get absolute path for HTML file: %w", err)
	}

	fileURL := "file://" + absPath

	var buf []byte

	err = chromedp.Run(ctx,
		chromedp.Navigate(fileURL),
		chromedp.WaitVisible("#meme-container"),
		chromedp.Screenshot("#meme-container", &buf, chromedp.NodeVisible),
	)
	if err != nil {
		return fmt.Errorf("failed to capture screenshot: %w", err)
	}

	img, _, err := image.Decode(bytes.NewReader(buf))
	if err != nil {
		return fmt.Errorf("failed to decode screenshot: %w", err)
	}

	resizedImg := resize.Resize(512, 0, img, resize.Lanczos3)

	outFile, err := os.Create(outputImageFile)
	if err != nil {
		return fmt.Errorf("failed to create output image file: %w", err)
	}
	defer outFile.Close()

	err = png.Encode(outFile, resizedImg)
	if err != nil {
		return fmt.Errorf("failed to save resized image: %w", err)
	}

	return nil
}

func WriteToMemeImage(dest string, lines []string, thumbnailPath, templatePath string, font string, fontSize int, previewOnly bool, bgColor, textColor string) error {

	unixTime := time.Now().Unix()
	timestamp := time.Unix(unixTime, 0).Format("02/01/2006, 15:04:05")
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
    </style>`, font, bgColor, fontSize, textColor))

	data := MemeData{
		Timestamp:  timestamp,
		UniqueID:   uniqueID,
		Lines:      lines,
		Thumbnail:  thumbnailPath,
		FontSize:   fontSize,
		StyleBlock: styleBlock,
	}

	htmlFile := "temp_meme.html"
	err := GenerateHTMLFile(htmlFile, data, templatePath)
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
