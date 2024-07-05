package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
)

func Encrypt(key, in, out string) {
	cmd := exec.Command("gpg", "--batch", "--yes", "--passphrase", key, "--output", out, "--symmetric", "--cipher-algo", "AES256", in)
	b := new(bytes.Buffer)
	cmd.Stderr = b
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error encrypting %s\n", in)
		fmt.Println(cmd.String())
		io.Copy(os.Stderr, b)
		panic(err)
	}
}

func Decrypt(key, in, out string) {
	cmd := exec.Command("gpg", "--batch", "--yes", "--passphrase", key, "--output", out, "--decrypt", in)
	b := new(bytes.Buffer)
	cmd.Stderr = b
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error decrypting %s\n", in)
		fmt.Println(cmd.String())
		io.Copy(os.Stderr, b)
		panic(err)
	}
}