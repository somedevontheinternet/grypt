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
		(func(job FileMeta) {
			defer wg.Done()
			fmt.Printf("DECRYPTING: %s\n", job.DecryptedName)
			inPath := filepath.Join(encryptDir, job.EncryptedName)
			outPath := filepath.Join(decryptedDir, job.DecryptedName)
			dir := filepath.Dir(outPath)
			err := os.MkdirAll(dir, 0755)
			if err != nil {
				fmt.Printf("Skipping %s: %s\n", job.DecryptedName, err)
				return
			}
			Decrypt(key, inPath, outPath)
			err = os.Chtimes(outPath, job.ModTime.UTC(), job.ModTime.UTC())
			if err != nil {
				fmt.Printf("Could not set back the modtime for %s\n", job.DecryptedName)

			}
		})(job)
	}
}

func PerformDecryption() {
	key := GetPassphrase(false)
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
