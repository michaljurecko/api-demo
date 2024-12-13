package entitygen

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
)

type fsWriter struct {
	logger    *slog.Logger
	targetDir string
}

func newFsWriter(logger *slog.Logger, targetDir string) *fsWriter {
	return &fsWriter{logger: logger, targetDir: targetDir}
}

func (w *fsWriter) WriteFiles(ctx context.Context, files []*genFile) (err error) {
	if err := os.RemoveAll(w.targetDir); err != nil {
		return fmt.Errorf("cannot clear directory: %w", err)
	}

	if err := os.MkdirAll(w.targetDir, dirPerm); err != nil {
		return fmt.Errorf("cannot create directory: %w", err)
	}

	var errs []error
	for _, file := range files {
		if err := w.writeEntityFile(ctx, file); err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}

func (w *fsWriter) writeEntityFile(ctx context.Context, spec *genFile) (err error) {
	w.logger.InfoContext(ctx, "writing file", slog.String("file", spec.FileName))

	osFile, err := os.OpenFile(filepath.Join(w.targetDir, spec.FileName), fileMode, filePerm)
	if err != nil {
		return fmt.Errorf(`cannot write file "%s": %w`, spec.FileName, err)
	}

	defer func() {
		if closeErr := osFile.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}()

	// spec.File.NoFormat = true

	if err := spec.File.Render(osFile); err != nil {
		return fmt.Errorf(`cannot render file "%s": %w`, spec.FileName, err)
	}

	return nil
}
