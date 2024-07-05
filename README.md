# Grypt

Grypt is a utility tool to encrypt / decrypt entire repo, focusing on simplicity and longetivity. It is intended for personal. Small uses. Such as "safely" storing a copy of your driver's license or tax report.

# Encryption commands

Every file is encrypted using the following command:

```
$ gpg --batch --yes --passphrase $KEY --output $OUTPUT_FILE --symmetric --cipher-algo AES256 $INPUT_FILE
```

and every file is decrypted using the following command:

```
$ gpg --batch --yes --passphrase $KEY --output $OUTPUT_FILE --decrypt $INPUT_FILE
```

and whatever version of `gpg` is installed on your system. I used this one:

```
$ gpg version
gpg (GnuPG/MacGPG2) 2.2.41
libgcrypt 1.8.10
Copyright (C) 2022 g10 Code GmbH
License GNU GPL-3.0-or-later <https://gnu.org/licenses/gpl.html>
This is free software: you are free to change and redistribute it.
There is NO WARRANTY, to the extent permitted by law.

Home: /Users/USERNAME/.gnupg
Supported algorit ms:
Pubkey: RSA, ELG, DSA, ECDH, ECDSA, EDDSA
Cipher: IDEA, 3DES, CAST5, BLOWFISH, AES, AES192, AES256, TWOFISH,
        CAMELLIA128, CAMELLIA192, CAMELLIA256
Hash: SHA1, RIPEMD160, SHA256, SHA384, SHA512, SHA224
Compression: Uncompressed, ZIP, ZLIB, BZIP2
```

## Encryption

Start with a folder called `src` and put all the files that you want to encrypt in there.

```
└── src
    ├── DRIVER_LICENSE.jpg
    └── TAXES
        └── 2022_FINAL.txt
```

Then go on the folder above `src` and run `grypt`:

```

$ grypt encrypt
Enter encryption key: 🔒
Confirm encryption key: 🔒
NEW FILE: DRIVER_LICENSE.jpg
NEW FILE: TAXES/2022_FINAL.txt

```

Once that's done several new files will be created:

- `blobs/*`

  These are the encrypted files. Unfortunately an attacker can still see the relative size of the files.

- `META.enc`

  This file contains meta information about the mapping between the original files and the blobs. It is encrypted using the same key as the files. It is a json array using the following format:

```

    [
        {
            "encryptedName": "string", // name of the blob
            "decryptedName": "string", // original name of the file
            "modTime": number, // last time the original file was modified, this is useful to keep track of changes
        }
    ]

```

- `.gitignore`

Automatically ignore `src/`, `META`, `out/`, but also other normal filesystem files.

- `README.md`

A readme is created in the repo in case you forget how these files were encrypted. It outlines the manual process of recovering the files. In the event `grypt` becomes unavailable.

## Decryption

With a valid `META.enc` and `blobs` folder:

```
$ grypt decrypt
Enter encryption key: 🔒
DECRYPTING: DRIVER_LICENSE.jpg
DECRYPTING: TAXES/2022_FINAL.txt
```

This will decrypt all the blobs in a `out` folder. The last modified time will be set in order to keep version control under control.

## Intended usage

The intended usage is to allow you to "safely" store private information on public servers.

1. Create a private git repo.
1. Add all your files in `src/`
1. Run `grypt encrypt` to encrypt your files.
1. Commit, verify that no private file is uploaded as-is.
1. Delete your `src/` folder so that these file almost never exist.
1. When you need to read/write those files: `grypt decrypt && mv out src`
1. Make more modification to the files in your `src/` folder.
1. Repeat `grypt encrypt` / `commit` / `rm -rf src/`
