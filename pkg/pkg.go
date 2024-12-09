package pkg

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/fogleman/gg"
	"github.com/golang/freetype"
	"github.com/nfnt/resize"

	_ "embed"
)

//go:embed assets/RobotoMonoForPowerline.ttf
var fontBytes []byte

func ReadInputLines(linesCount int) ([]string, error) {
	var lines []string

	reader := bufio.NewReader(os.Stdin)
	for i := 0; i < linesCount; i++ {
		fmt.Print("> ")
		line, err := reader.ReadString('\n')
		line = strings.TrimSpace(line)
		if err != nil {
			return nil, err
		}
		lines = append(lines, line)
	}

	return lines, nil
}

func WriteToImage(dest string, lines []string, thumbnailPath string) error {

	fontParsed, err := freetype.ParseFont(fontBytes)
	if err != nil {
		return err
	}

	unixTime := time.Now().Unix()
	timestamp := time.Unix(unixTime, 0).Format("02/01/2006, 15:04:05")
	noField := strconv.FormatInt(unixTime, 10)[:8]

	var thumb image.Image
	var thumbHeight int

	if thumbnailPath == "" {

		maxThumbWidth := 120
		thumbHeight = 120
		placeholder := image.NewRGBA(image.Rect(0, 0, maxThumbWidth, thumbHeight))
		placeholderBg := color.RGBA{200, 200, 200, 255}
		draw.Draw(placeholder, placeholder.Bounds(), &image.Uniform{C: placeholderBg}, image.Point{}, draw.Src)

		ftContext := freetype.NewContext()
		ftContext.SetDPI(100)
		ftContext.SetFont(fontParsed)
		ftContext.SetFontSize(10)
		ftContext.SetDst(placeholder)
		ftContext.SetClip(placeholder.Bounds())
		ftContext.SetSrc(&image.Uniform{C: color.Black})
		_, err := ftContext.DrawString("tfw.png", freetype.Pt(20, 60))
		if err != nil {
			return err
		}
		thumb = placeholder
	} else {

		thumbImgFile, err := os.Open(thumbnailPath)
		if err != nil {
			return err
		}
		defer thumbImgFile.Close()
		thumbImg, _, err := image.Decode(thumbImgFile)
		if err != nil {
			return err
		}

		maxThumbWidth := 120
		scaleFactor := float64(maxThumbWidth) / float64(thumbImg.Bounds().Dx())
		thumbHeight = int(float64(thumbImg.Bounds().Dy()) * scaleFactor)
		thumb = resize.Resize(uint(maxThumbWidth), uint(thumbHeight), thumbImg, resize.NearestNeighbor)
	}

	textWidth := 512 - 120 - 50
	lineHeight := 20
	textLines := wrapText(lines, textWidth, 10)
	textHeight := len(textLines) * lineHeight
	imgWidth := 512
	imgHeight := 100 + textHeight
	if len(lines) < 5 {
		imgHeight = 180 + textHeight
	}

	img := image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight))
	bgColor := color.RGBA{240, 224, 214, 255}
	draw.Draw(img, img.Bounds(), &image.Uniform{C: bgColor}, image.Point{}, draw.Src)

	drawRoundedThumbnail(img, thumb, 20, 40, 120, thumbHeight, 5)

	ftContext := freetype.NewContext()
	ftContext.SetDPI(100)
	ftContext.SetFont(fontParsed)
	ftContext.SetClip(img.Bounds())
	ftContext.SetDst(img)

	headerHeight := 30
	headerBgColor := color.RGBA{226, 209, 198, 255}
	headerRect := image.Rect(0, 0, imgWidth, headerHeight)
	draw.Draw(img, headerRect, &image.Uniform{C: headerBgColor}, image.Point{}, draw.Src)

	ftContext.SetFontSize(10)

	ftContext.SetSrc(&image.Uniform{C: color.RGBA{0, 114, 50, 255}})
	_, err = ftContext.DrawString("Anonymous", freetype.Pt(10, 20))
	if err != nil {
		return err
	}

	ftContext.SetSrc(&image.Uniform{C: color.RGBA{138, 0, 0, 255}})
	_, err = ftContext.DrawString(timestamp, freetype.Pt(150, 20))
	if err != nil {
		return err
	}

	ftContext.SetSrc(&image.Uniform{C: color.RGBA{138, 0, 0, 255}})
	_, err = ftContext.DrawString("No."+noField, freetype.Pt(imgWidth-100, 20))
	if err != nil {
		return err
	}

	lineColor := color.RGBA{229, 207, 199, 255}
	headerLineHeight := 1
	lineRect := image.Rect(0, headerHeight, imgWidth, headerHeight+headerLineHeight)
	draw.Draw(img, lineRect, &image.Uniform{C: lineColor}, image.Point{}, draw.Src)

	fileSizeText := fmt.Sprintf("%d KB PNG", 120*thumbHeight/1024)
	ftContext.SetSrc(&image.Uniform{C: color.RGBA{136, 136, 136, 255}})
	_, err = ftContext.DrawString(fileSizeText, freetype.Pt(20, 55+thumbHeight))
	if err != nil {
		return err
	}

	for y := 10; y < imgHeight-10; y++ {
		img.Set(20+120+10, y, lineColor)
	}

	ftContext.SetFontSize(12)
	ftContext.SetSrc(&image.Uniform{C: color.RGBA{120, 153, 34, 255}})
	yOffset := 50
	for _, line := range textLines {
		_, err := ftContext.DrawString(line, freetype.Pt(20+120+20, yOffset))
		if err != nil {
			return err
		}
		yOffset += lineHeight
	}

	outFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer outFile.Close()

	err = png.Encode(outFile, img)
	if err != nil {
		return err
	}

	return nil
}

func wrapText(lines []string, maxWidth int, fontSize int) []string {
	wrappedLines := []string{}
	for _, line := range lines {
		current := "> "
		for _, word := range strings.Split(line, " ") {
			testLine := current + " " + word
			if len(testLine)*fontSize > maxWidth {
				wrappedLines = append(wrappedLines, current)
				current = word
			} else {
				current = testLine
			}
		}
		wrappedLines = append(wrappedLines, current)
	}
	return wrappedLines
}

func drawRoundedThumbnail(img draw.Image, thumb image.Image, x, y, width, height, radius int) {
	mask := image.NewRGBA(image.Rect(0, 0, width, height))
	dc := gg.NewContextForRGBA(mask)
	dc.DrawRoundedRectangle(0, 0, float64(width), float64(height), float64(radius))
	dc.SetRGB(1, 1, 1)
	dc.Fill()

	thumbRect := image.Rect(x, y, x+width, y+height)

	draw.DrawMask(img, thumbRect, thumb, image.Point{}, mask, image.Point{}, draw.Over)
}
