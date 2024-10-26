package meme

import (
	"fmt"
	"os/exec"

	"meme-generator/internal/storage"
)

type MemeManager struct {
	dataDir string

	store storage.Store
}

func NewMemeManager(dataDir string, store storage.Store) (*MemeManager, error) {
	if err := checkImageMagickAvailable(); err != nil {
		return nil, err
	}

	return &MemeManager{
		dataDir: dataDir,
		store:   store,
	}, nil
}

func checkImageMagickAvailable() error {
	_, err := exec.LookPath("convert")
	if err != nil {
		return fmt.Errorf("ImageMagick 'convert' tool not found in PATH: %w", err)
	}
	return nil
}
