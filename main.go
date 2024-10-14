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
var passphrase = flag.String("passphrase", "", "passphrase to use for encryption")
var passfile = flag.String("passfile", "", "passfile to use for encryption")

const usage = "Usage: grypt encrypt|decrypt"

func main() {
	CheckGPG()
	GenerateExtras()
	flag.Parse()
	args := flag.Args()
	if len(args) != 1 {
		fmt.Println(usage)
		os.Exit(1)
	}
	switch args[0] {
	case "encrypt":
		PerformEncryption()
	case "decrypt":
		PerformDecryption()
	default:
		fmt.Println(usage)
		os.Exit(1)
	}
}
