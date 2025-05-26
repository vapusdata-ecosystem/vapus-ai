package nemo

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/rs/zerolog"
)

// Instead of UserId, I need doamin Id
// Taking File as an Input and saving it
func NewNemoGuardrail(ctx context.Context, fileName string, userId string, logger zerolog.Logger) error {
	// Open the source file
	srcFile, err := os.Open(fileName)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to open input file")
		return err
	}
	defer srcFile.Close()

	// Create target directory if it doesn't exist
	targetDir := filepath.Join("configfiles", userId)
	err = os.MkdirAll(targetDir, os.ModePerm)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to create target directory")
		return err
	}

	// Create the target file path
	targetFile := filepath.Join(targetDir, fmt.Sprintf("config_%s.yaml", userId))
	dstFile, err := os.Create(targetFile)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to create target file")
		return err
	}
	defer dstFile.Close()

	// Copy content from source to destination
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to copy file content")
		return err
	}

	logger.Info().Str("file", targetFile).Msg("Configuration file saved successfully")
	return nil
}
