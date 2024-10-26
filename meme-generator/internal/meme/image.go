package meme

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
)

var DefaultImageName = "meme.png"

var allowedImageTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
}

var (
	allowedColors = []string{"white", "black"}
	defaultColor  = "white"
)

func CheckImageType(fileHeader *multipart.FileHeader) (bool, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return false, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	buffer := make([]byte, 512)
	if _, err = file.Read(buffer); err != nil {
		return false, fmt.Errorf("failed to read file: %w", err)
	}

	if _, err := file.Seek(0, 0); err != nil {
		return false, fmt.Errorf("failed to seek file: %w", err)
	}

	contentType := http.DetectContentType(buffer)
	if _, ok := allowedImageTypes[contentType]; ok {
		return true, nil
	}

	return false, nil
}

func getImageDimensions(filePath string) (width, height int, err error) {
	cmd := exec.Command("bash", "-c", "identify "+"-format "+"%wx%h "+"'"+filePath+"'")
	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return 0, 0, fmt.Errorf("failed to get image dimensions: %v", err)
	}

	dimensions := strings.TrimSpace(out.String())
	parts := strings.Split(dimensions, "x")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("unexpected dimensions format: %s", dimensions)
	}

	width, err = strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, fmt.Errorf("failed to parse width: %v", err)
	}

	height, err = strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, fmt.Errorf("failed to parse height: %v", err)
	}

	return width, height, nil
}

func calculateFontSize(text string, imageWidth int) int {
	minFontSize := imageWidth / 20
	maxFontSize := imageWidth / 6

	estimatedFontSize := imageWidth / (len(text)/5 + 1)

	if estimatedFontSize < minFontSize {
		estimatedFontSize = minFontSize
	} else if estimatedFontSize > maxFontSize {
		estimatedFontSize = maxFontSize
	}

	return estimatedFontSize
}

func getTextAllowedTextColor(color string) string {
	for _, allowedColor := range allowedColors {
		if color == allowedColor {
			return color
		}
	}

	return defaultColor
}
