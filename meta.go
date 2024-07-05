package main

import (
	"encoding/json"
	"errors"
	"os"
	"time"
)

type FileMeta struct {
	EncryptedName string    `json:"encryptedName"`
	DecryptedName string    `json:"decryptedName"`
	ModTime       time.Time `json:"modTime"`
}

type Meta []FileMeta

func CreateMeta(key string) {
	f, err := os.OpenFile(MetaDecrypted, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	defer os.Remove(MetaDecrypted)
	err = json.NewEncoder(f).Encode(Meta{})
	if err != nil {
		panic(err)
	}
	Encrypt(key, MetaDecrypted, MetaEncrypted)
}

func ReadMeta(key string) Meta {
	Decrypt(key, MetaEncrypted, MetaDecrypted)
	f, err := os.Open(MetaDecrypted)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	defer os.Remove(MetaDecrypted)
	var meta Meta
	err = json.NewDecoder(f).Decode(&meta)
	if err != nil {
		panic(err)
	}
	return meta
}

func GetMeta(key string) Meta {
	_, err := os.Stat(MetaEncrypted)
	if errors.Is(err, os.ErrNotExist) {
		CreateMeta(key)
		return ReadMeta(key)
	}
	if err != nil {
		panic(err)
	}
	return ReadMeta(key)
}

func SaveMeta(key string, meta Meta) {
	f, err := os.OpenFile(MetaDecrypted, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	defer os.Remove(MetaDecrypted)
	err = json.NewEncoder(f).Encode(meta)
	if err != nil {
		panic(err)
	}

	Encrypt(key, MetaDecrypted, MetaEncrypted)
}
