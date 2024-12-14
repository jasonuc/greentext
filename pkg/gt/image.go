package gt

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/png"
	"os"
	"path/filepath"

	"github.com/chromedp/chromedp"
	"github.com/nfnt/resize"
)

const (
	defaultWidth  uint = 512
	defaultHeight uint = 0 // If one of the parameters width or height is set to 0, its size will be calculated so that the aspect ratio is that of the originating image.
)

func CaptureElementScreenshot(htmlFile, outputImageFile string, w, h int) error {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	absPath, err := filepath.Abs(htmlFile)
	if err != nil {
		return fmt.Errorf("failed to get absolute path for HTML file: %w", err)
	}

	fileURL := "file://" + absPath

	var buf []byte //

	err = chromedp.Run(ctx,
		chromedp.Navigate(fileURL),
		chromedp.WaitVisible("#greentext"),
		chromedp.Screenshot("#greentext", &buf, chromedp.NodeVisible),
	)
	if err != nil {
		return fmt.Errorf("failed to capture screenshot: %w", err)
	}

	img, _, err := image.Decode(bytes.NewReader(buf))
	if err != nil {
		return fmt.Errorf("failed to decode screenshot: %w", err)
	}

	var width, height uint = uint(w), uint(h)
	if width+height == 0 {
		width, height = defaultWidth, defaultHeight
	}

	resizedImg := resize.Resize(width, height, img, resize.Lanczos3)

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
