package gt

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/png"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

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

func checkIfThumnailIsURL(thumbnail string) bool {
	return strings.HasPrefix(thumbnail, "http://") || strings.HasPrefix(thumbnail, "https://")
}

func getImageSizeFromURL(url string) (int64, error) {
	// Send a HEAD request
	resp, err := http.Head(url)
	if err != nil {
		return 0, fmt.Errorf("error sending HEAD request: %v", err)
	}
	defer resp.Body.Close()

	// Check for the Content-Length header
	contentLength := resp.Header.Get("Content-Length")
	if contentLength != "" {
		size, err := strconv.ParseInt(contentLength, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("error parsing Content-Length: %v", err)
		}
		return size, nil
	}

	// If Content-Length is not available, download the file to measure its size
	resp, err = http.Get(url)
	if err != nil {
		return 0, fmt.Errorf("error sending GET request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body to count bytes
	var totalBytes int64
	buf := make([]byte, 1024)
	for {
		n, err := resp.Body.Read(buf)
		totalBytes += int64(n)
		if err == io.EOF {
			break
		}
		if err != nil {
			return 0, fmt.Errorf("error reading response body: %v", err)
		}
	}

	// convert totalBytes to KB
	totalBytes /= 1024
	return totalBytes, nil
}

func getThumbnailSize(thumbnail string) int {
	if checkIfThumnailIsURL(thumbnail) {
		size, err := getImageSizeFromURL(thumbnail)
		if err != nil {
			fmt.Println("Error getting thumbnail size:", err)
			return 0
		}
		return int(size)
	}

	info, err := os.Stat(thumbnail)
	if err != nil {
		fmt.Println("Error getting thumbnail size:", err)
		return 0
	}
	return int(info.Size() / 1024) // convert bytes to KB
}
