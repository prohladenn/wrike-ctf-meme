package meme

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	"meme-generator/internal/model"
	"meme-generator/internal/utils"

	"github.com/gosimple/slug"
)

const (
	memesDir = "memes"
)

func (m *MemeManager) NewMeme(template *model.Template, templateOwnerUsername string, authorID uint, name, caption, textColor string) (*model.Meme, error) {
	// Prepare template struct
	meme := &model.Meme{
		Name:           name,
		UserID:         authorID,
		MemeTemplateID: template.ID,
	}

	// Create template directory path in format "data/user_{user_id}_memes/slug({name})"
	memeDir := filepath.Join(m.dataDir, "/user_"+strconv.Itoa(int(authorID)), memesDir, slug.Make(name))

	// Check if the directory already exists
	info, err := os.Stat(memeDir)
	if err == nil && info.IsDir() {
		randomSuffix, err := utils.GenerateRandomString(4)
		if err != nil {
			return nil, fmt.Errorf("failed to generate random string: %w", err)
		}

		// Add a random suffix to the directory name to avoid conflicts
		memeDir = memeDir + "-" + randomSuffix
	} else if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to stat directory: %w", err)
	}

	// Create the directory
	if err := os.MkdirAll(memeDir, os.ModePerm); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	meme.DirPath = memeDir
	meme.FileName = slug.Make(name) + filepath.Ext(template.FileName)

	// Get full path to the template file and new meme file
	templateFilePath := filepath.Join(template.DirPath, template.FileName)
	memeFilePath := filepath.Join(memeDir, meme.FileName)

	// Generate the meme image with caption and watermark using ImageMagick
	if err := m.generateMemeImage(templateFilePath, memeFilePath, caption, "@"+templateOwnerUsername, textColor); err != nil {
		return nil, fmt.Errorf("failed to generate meme image: %w", err)
	}

	return meme, nil
}

func (m *MemeManager) generateMemeImage(templateFilePath, memeFilePath, caption, watermark, textColor string) error {
	width, height, err := getImageDimensions(templateFilePath)
	if err != nil {
		return err
	}

	captionFontSize := calculateFontSize(caption, width)
	watermarkFontSize := width / 20

	captionWidth := width - (width / 10)

	watermarkPadding := width / 30
	topPadding := height / 30

	cmd := exec.Command("convert",
		templateFilePath, // Input template image

		"(",
		"-background", "none",
		"-fill", getTextAllowedTextColor(textColor),
		"-font", "DejaVu-Sans",
		"-pointsize", fmt.Sprintf("%d", captionFontSize),
		"-gravity", "center",
		"-size", fmt.Sprintf("%d", captionWidth),
		"caption:"+caption,
		")",

		"-gravity", "north",
		"-geometry", fmt.Sprintf("+0+%d", topPadding),
		"-composite", // Overlay the caption onto the template image

		// Add watermark
		"-gravity", "southwest",
		"-fill", getTextAllowedTextColor(textColor),
		"-font", "DejaVu-Sans",
		"-pointsize", fmt.Sprintf("%d", watermarkFontSize),
		"-annotate", fmt.Sprintf("+%d+%d", watermarkPadding, watermarkPadding), watermark,

		memeFilePath, // Output meme file
	)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("ImageMagick convert command failed: %v, stderr: %s", err, stderr.String())
	}

	return nil
}
