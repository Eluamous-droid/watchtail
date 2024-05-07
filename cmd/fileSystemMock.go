package cmd

import (
	"os"
	"time"
)

type MockFileInfo struct {
    FileName    string
    IsDirectory bool
		LastModTime			time.Time
}

func (mfi MockFileInfo) Name() string       { return mfi.FileName }
func (mfi MockFileInfo) Size() int64        { return int64(8) }
func (mfi MockFileInfo) Mode() os.FileMode  { return os.ModePerm }
func (mfi MockFileInfo) ModTime() time.Time { return mfi.LastModTime }
func (mfi MockFileInfo) IsDir() bool        { return mfi.IsDirectory }
func (mfi MockFileInfo) Sys() interface{}   { return nil }
