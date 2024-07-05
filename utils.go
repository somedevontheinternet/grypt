package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/term"
)

func prepareDir(dir string) {
	info, err := os.Stat(dir)
	if errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(dir, 0755)
		if err != nil {
			panic(err)
		}
		return
	}
	if err != nil {
		panic(err)
	}
	if !info.IsDir() {
		panic(fmt.Errorf("%s is not a directory", dir))
	}
}

func GetEncryptionKey(confirm bool) string {
	fmt.Print("Enter encryption key: ")
	key0, err := term.ReadPassword(0)
	if err != nil {
		panic(err)
	}
	fmt.Println()

	if confirm {
		fmt.Print("Confirm encryption key: ")
		key1, err := term.ReadPassword(0)
		if err != nil {
			panic(err)
		}
		fmt.Println()
		if string(key0) != string(key1) {
			panic(fmt.Errorf("encryption keys do not match"))
		}
	}
	return string(key0)
}

func ListFiles(dir string) []FileMeta {
	dirAbs, err := filepath.Abs(dir)
	if err != nil {
		panic(err)
	}
	var files []FileMeta
	err = filepath.Walk(dirAbs, func(path string, info os.FileInfo, _ error) error {
		if info.IsDir() {
			return nil // skip
		}
		abs, err := filepath.Abs(path)
		if err != nil {
			return err
		}
		meta := FileMeta{
			DecryptedName: abs[len(dirAbs)+len(string(os.PathSeparator)):],
			ModTime:       info.ModTime(),
		}
		files = append(files, meta)
		return nil
	})
	if err != nil {
		panic(err)
	}
	return files
}
