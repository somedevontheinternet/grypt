package main

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type EncryptRunner struct {
	Key          string
	DecryptedDir string
	EncryptedDir string
	Meta         Meta
}

func NewEncryptRunner(key, decryptedDir, encryptedDir string, meta Meta) *EncryptRunner {
	return &EncryptRunner{
		Key:          key,
		DecryptedDir: decryptedDir,
		EncryptedDir: encryptedDir,
		Meta:         meta,
	}
}

func (r *EncryptRunner) findMeta(file FileMeta) (int, bool) {
	for i, m := range r.Meta {
		if m.DecryptedName == file.DecryptedName {
			return i, true
		}
	}
	return -1, false
}

func (r *EncryptRunner) processFile(file FileMeta) {
	lastMetaI, found := r.findMeta(file)
	if !found { // brand new file
		fmt.Println("NEW FILE: ", file.DecryptedName)
		file.EncryptedName = base64.URLEncoding.EncodeToString(sha256.New().Sum([]byte(file.DecryptedName)))
		Encrypt(r.Key, filepath.Join(r.DecryptedDir, file.DecryptedName), filepath.Join(r.EncryptedDir, file.EncryptedName))
		r.Meta = append(r.Meta, file)
		return
	}
	if r.Meta[lastMetaI].ModTime.UTC() == file.ModTime.UTC() { // no change
		return
	}

	fmt.Println("MODIFIED: ", file.DecryptedName)
	r.Meta[lastMetaI].ModTime = file.ModTime
	Encrypt(r.Key, filepath.Join(r.DecryptedDir, file.DecryptedName), filepath.Join(r.EncryptedDir, r.Meta[lastMetaI].EncryptedName))
}

func (r *EncryptRunner) Run(files []FileMeta) {
	prepareDir(r.EncryptedDir)

	// Encrypt every new / changed files
	for _, file := range files {
		r.processFile(file)
	}

	// Check every old file for deletion

	for i := 0; i < len(r.Meta); i++ {
		_, err := os.Stat(filepath.Join(r.DecryptedDir, r.Meta[i].DecryptedName))
		if errors.Is(err, os.ErrNotExist) {
			fmt.Println("DELETED: ", r.Meta[i].DecryptedName)
			// Delete the file FIRST and then remove the meta
			os.Remove(filepath.Join(r.EncryptedDir, r.Meta[i].EncryptedName))

			r.Meta[i] = r.Meta[len(r.Meta)-1]
			r.Meta = r.Meta[:len(r.Meta)-1]
			i--
			continue
		}
	}

	SaveMeta(r.Key, r.Meta)
}

func PerformEncryption() {
	key := GetEncryptionKey(true)
	meta := GetMeta(key)
	files := ListFiles(*SrcRoot)
	r := NewEncryptRunner(key, *SrcRoot, *BlobsRoot, meta)
	r.Run(files)
}
