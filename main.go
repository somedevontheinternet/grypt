package main

import (
	"flag"
	"fmt"
	"os"
)

const MetaDecrypted = "META"
const MetaEncrypted = "META.enc"

var SrcRoot = flag.String("src", "src", "source directory")
var BlobsRoot = flag.String("blobs", "blobs", "blob directory")
var OutRoot = flag.String("out", "out", "output directory")

func main() {
	GenerateExtras()
	flag.Parse()
	args := flag.Args()
	if len(args) != 1 {
		fmt.Println("Usage: grypt encrypt|decrypt")
		os.Exit(1)
	}
	switch args[0] {
	case "encrypt":
		PerformEncryption()
	case "decrypt":
		PerformDecryption()
	default:
		fmt.Println("Usage: grypt encrypt|decrypt")
		os.Exit(1)
	}
}
