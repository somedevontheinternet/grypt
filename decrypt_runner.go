package main

import (
	"fmt"
	"os"
	"path/filepath"
)

type DecryptRunner struct {
	Key          string
	DecryptedDir string
	EncryptDir   string
	Meta         Meta
}

func NewDecryptRunner(key, encryptDir, decryptedDir string, meta Meta) *DecryptRunner {
	return &DecryptRunner{
		Key:          key,
		DecryptedDir: decryptedDir,
		EncryptDir:   encryptDir,
		Meta:         meta,
	}
}

func (r *DecryptRunner) Run(meta Meta) {
	prepareDir(*OutRoot)
	for _, file := range meta {
		fmt.Printf("DECRYPTING: %s\n", file.DecryptedName)
		inPath := filepath.Join(r.EncryptDir, file.EncryptedName)
		outPath := filepath.Join(r.DecryptedDir, file.DecryptedName)
		dir := filepath.Dir(outPath)
		os.MkdirAll(dir, 0755)
		Decrypt(r.Key, inPath, outPath)
		os.Chtimes(outPath, file.ModTime, file.ModTime)
	}
}

func PerformDecryption() {
	key := GetEncryptionKey(false)
	meta := GetMeta(key)
	r := NewDecryptRunner(key, *BlobsRoot, *OutRoot, meta)
	r.Run(meta)
}
