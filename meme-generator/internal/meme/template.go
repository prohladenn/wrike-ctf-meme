package meme

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"meme-generator/internal/model"
	"meme-generator/internal/utils"

	"github.com/gosimple/slug"
)

const (
	templatesDir     = "templates"
	templateInfoFile = "info.txt"
)

func sanitizeFileName(fileName string) error {
	if strings.Contains(fileName, ";") {
		return fmt.Errorf("invalid file name: contains ';' character")
	}
	return nil
}

func (m *MemeManager) NewTemplate(name, comment string, fileHeader *multipart.FileHeader, userID uint) (*model.Template, error) {
	// Prepare template struct
	template := &model.Template{
		Name:   name,
		UserID: userID,
	}

	// Create template directory path in format "data/user_{user_id}/templates/slug({name})"
	templateDir := m.dataDir + "/user_" + strconv.Itoa(int(userID)) + "/" + templatesDir + "/" + slug.Make(name)

	// Check if the directory already exists
	info, err := os.Stat(templateDir)
	if err == nil && info.IsDir() {
		randomSuffix, err := utils.GenerateRandomString(4)
		if err != nil {
			return nil, fmt.Errorf("failed to generate random string: %w", err)
		}

		// Add a random suffix to the directory name to avoid conflicts
		templateDir = templateDir + "-" + randomSuffix
	} else if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to stat directory: %w", err)
	}

	// Create the directory
	if err := os.MkdirAll(templateDir, os.ModePerm); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	template.DirPath = templateDir

	filename := fileHeader.Filename
	if filename == "" {
		filename = DefaultImageName
	}

	if err := sanitizeFileName(filename); err != nil {
		return nil, err
	}

	template.FileName = filename

	file, err := fileHeader.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	dstFile, err := os.OpenFile(filepath.Join(templateDir, filename), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer dstFile.Close()

	hash := sha256.New()

	// Create a multi writer that writes to both the destination file and the hash
	multiWriter := io.MultiWriter(dstFile, hash)

	if _, err := io.Copy(multiWriter, file); err != nil {
		return nil, fmt.Errorf("failed to copy file: %w", err)
	}

	hashValue := hash.Sum(nil)
	hashHex := hex.EncodeToString(hashValue)

	// Write the hash to a file in the meme directory with a .info extension
	hashFilePath := filepath.Join(templateDir, templateInfoFile)
	hashFile, err := os.OpenFile(hashFilePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer hashFile.Close()

	// Write the hash to the file
	if _, err := hashFile.WriteString("SHA256:" + hashHex + "\n"); err != nil {
		return nil, fmt.Errorf("failed to write hash to file: %w", err)
	}

	// Write the comment to the file
	if comment != "" {
		if _, err := hashFile.WriteString(comment + "\n"); err != nil {
			return nil, fmt.Errorf("failed to write comment to file: %w", err)
		}
	}

	return template, nil
}

func (m *MemeManager) GetTemplatePrivateInfo(tempalte *model.Template) (string, error) {
	filePath := filepath.Join(tempalte.DirPath, templateInfoFile)

	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	return string(content), nil
}
