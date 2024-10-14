package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
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

func getPassphraseFromStdin(confirm bool) string {
	fmt.Print("Enter passphrase: ")
	key0, err := term.ReadPassword(0)
	if err != nil {
		panic(err)
	}
	fmt.Println()

	if confirm {
		fmt.Print("Confirm passphrase: ")
		key1, err := term.ReadPassword(0)
		if err != nil {
			panic(err)
		}
		fmt.Println()
		if string(key0) != string(key1) {
			panic(fmt.Errorf("passphrases do not match"))
		}
	}
	return string(key0)
}

func GetPassphrase(confirm bool) string {
	// passphrase from cli arguments
	if passphrase != nil && *passphrase != "" {
		return *passphrase
	}

	// passphrase from file
	if passfile != nil && *passfile != "" {
		src, err := os.ReadFile(*passfile)
		if err != nil {
			panic(err)
		}
		return string(src)
	}

	return getPassphraseFromStdin(confirm)
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
			ModTime:       info.ModTime().UTC(),
		}
		files = append(files, meta)
		return nil
	})
	if err != nil {
		panic(err)
	}
	return files
}

func CheckGPG() {
	_, err := exec.LookPath(gpg)
	if err != nil {
		panic(err)
	}
}
