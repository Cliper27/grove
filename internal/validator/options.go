package validator

import (
	"io/fs"
)

func checkMaxSize(entry fs.DirEntry, maxSize uint64) (bool, error) {
	info, err := entry.Info()
	if err != nil {
		return false, err
	}
	return uint64(info.Size()) <= maxSize, nil
}

func dirSize(fsys fs.FS, dir string) (uint64, error) {
	var total uint64

	err := fs.WalkDir(fsys, dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		info, err := d.Info()
		if err != nil {
			return err
		}
		total += uint64(info.Size())
		return nil
	})

	return total, err
}
