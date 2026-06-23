package platform

import (
	"context"
	"os"
)

type FS interface {
	MkdirAll(ctx context.Context, path string, perm os.FileMode) error
	ReadFile(ctx context.Context, path string) ([]byte, error)
	WriteFile(ctx context.Context, path string, data []byte, perm os.FileMode) error
	Rename(ctx context.Context, oldPath, newPath string) error
	Remove(ctx context.Context, path string) error
}

type OSFS struct{}

func (OSFS) MkdirAll(ctx context.Context, path string, perm os.FileMode) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return os.MkdirAll(path, perm)
}

func (OSFS) ReadFile(ctx context.Context, path string) ([]byte, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	return os.ReadFile(path)
}

func (OSFS) WriteFile(ctx context.Context, path string, data []byte, perm os.FileMode) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return os.WriteFile(path, data, perm)
}

func (OSFS) Rename(ctx context.Context, oldPath, newPath string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return os.Rename(oldPath, newPath)
}

func (OSFS) Remove(ctx context.Context, path string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return os.Remove(path)
}
