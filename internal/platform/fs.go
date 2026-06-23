package platform

import (
	"context"
	"os"
)

// FS abstracts filesystem operations used by stores and source readers.
type FS interface {
	MkdirAll(ctx context.Context, path string, perm os.FileMode) error
	ReadFile(ctx context.Context, path string) ([]byte, error)
	WriteFile(ctx context.Context, path string, data []byte, perm os.FileMode) error
	Rename(ctx context.Context, oldPath, newPath string) error
	Remove(ctx context.Context, path string) error
}

// OSFS implements FS with the host operating system filesystem.
type OSFS struct{}

// MkdirAll creates a directory tree on the host filesystem.
func (OSFS) MkdirAll(ctx context.Context, path string, perm os.FileMode) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return os.MkdirAll(path, perm)
}

// ReadFile reads a file from the host filesystem.
func (OSFS) ReadFile(ctx context.Context, path string) ([]byte, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	// #nosec G304 -- platform.FS is the repository's explicit host filesystem
	// capability; callers own path confinement at the store/source boundary.
	return os.ReadFile(path)
}

// WriteFile writes a file to the host filesystem.
func (OSFS) WriteFile(ctx context.Context, path string, data []byte, perm os.FileMode) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return os.WriteFile(path, data, perm)
}

// Rename renames a file or directory on the host filesystem.
func (OSFS) Rename(ctx context.Context, oldPath, newPath string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return os.Rename(oldPath, newPath)
}

// Remove removes a file or empty directory from the host filesystem.
func (OSFS) Remove(ctx context.Context, path string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return os.Remove(path)
}
