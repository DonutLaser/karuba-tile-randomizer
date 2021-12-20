package main

import (
	"io/fs"
	"os"
)

func ReadDirectory(fullPath string) []fs.DirEntry {
	dir, err := os.Open(fullPath)
	checkError(err)
	defer dir.Close()

	items, err := dir.ReadDir(-1)
	checkError(err)

	return items
}
