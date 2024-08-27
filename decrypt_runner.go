package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

type DecryptRunner struct {
	Key          string
	DecryptedDir string
	EncryptDir   string
}

func NewDecryptRunner(key, encryptDir, decryptedDir string, jobs chan FileMeta, wg *sync.WaitGroup) {
	for job := range jobs {
		fmt.Printf("DECRYPTING: %s\n", job.DecryptedName)
		inPath := filepath.Join(encryptDir, job.EncryptedName)
		outPath := filepath.Join(decryptedDir, job.DecryptedName)
		dir := filepath.Dir(outPath)
		os.MkdirAll(dir, 0755)
		Decrypt(key, inPath, outPath)
		os.Chtimes(outPath, job.ModTime, job.ModTime)
		wg.Done()
	}
}

func PerformDecryption() {
	key := GetEncryptionKey(false)
	meta := GetMeta(key)
	prepareDir(*OutRoot)
	runners := runtime.NumCPU()
	var wg sync.WaitGroup
	jobs := make(chan FileMeta)
	for i := 0; i < runners; i++ {
		go NewDecryptRunner(key, *BlobsRoot, *OutRoot, jobs, &wg)
	}
	before := time.Now()
	for _, file := range meta {
		wg.Add(1)
		jobs <- file
	}
	wg.Wait()
	after := time.Now()
	fmt.Printf("Decryption completed in %s\n", (after.Sub(before)))
}
