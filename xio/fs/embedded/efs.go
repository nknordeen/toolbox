package embedded

import (
	"bytes"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/richardwilkes/toolbox/log/jot"
)

type efs struct {
	files      map[string]*File
	dirModTime time.Time
}

func (f *efs) IsLive() bool {
	return false
}

func (f *efs) actualPath(path string) string {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	return filepath.ToSlash(filepath.Clean(path))
}

func (f *efs) Open(path string) (http.File, error) {
	one, ok := f.files[f.actualPath(path)]
	if !ok {
		return nil, os.ErrNotExist
	}
	if one.isDir {
		return one, nil
	}
	if err := one.uncompressData(); err != nil {
		return nil, err
	}
	return &File{
		Reader:  bytes.NewReader(one.data),
		name:    one.name,
		size:    one.size,
		modTime: one.modTime,
		data:    one.data,
	}, nil
}

func (f *efs) ContentAsBytes(path string) ([]byte, bool) {
	if one, ok := f.files[f.actualPath(path)]; ok {
		if err := one.uncompressData(); err != nil {
			return nil, false
		}
		return one.data, true
	}
	return nil, false
}

func (f *efs) MustContentAsBytes(path string) []byte {
	if d, ok := f.ContentAsBytes(path); ok {
		return d
	}
	jot.Fatal(1, path+" does not exist")
	return nil
}

func (f *efs) ContentAsString(path string) (string, bool) {
	if d, ok := f.ContentAsBytes(path); ok {
		return string(d), true
	}
	return "", false
}

func (f *efs) MustContentAsString(path string) string {
	if s, ok := f.ContentAsString(path); ok {
		return s
	}
	jot.Fatal(1, path+" does not exist")
	return ""
}
