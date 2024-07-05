# Grypt encrypted repo

This data was encrypted using Grypt. In the event that `grypt` is not available, follow these steps to recover the data:

## GPG

Get gpg with the following (or compatible) version:

```
$ gpg version
gpg (GnuPG/MacGPG2) 2.2.41
libgcrypt 1.8.10
Copyright (C) 2022 g10 Code GmbH
License GNU GPL-3.0-or-later <https://gnu.org/licenses/gpl.html>
This is free software: you are free to change and redistribute it.
There is NO WARRANTY, to the extent permitted by law.

Home: /Users/USERNAME/.gnupg
Supported algorithms:
Pubkey: RSA, ELG, DSA, ECDH, ECDSA, EDDSA
Cipher: IDEA, 3DES, CAST5, BLOWFISH, AES, AES192, AES256, TWOFISH,
 CAMELLIA128, CAMELLIA192, CAMELLIA256
Hash: SHA1, RIPEMD160, SHA256, SHA384, SHA512, SHA224
Compression: Uncompressed, ZIP, ZLIB, BZIP2
```

## META.enc

Decrypt this file to retrieve the metadata for the encrypted files. The key is the same as the one used to encrypt the files.

The structure is the following:

```
[
    {
        "encryptedName": "string", // sha256 hash of the decrypted name
        "decryptedName": "string", // original name of the file
        "modTime": number // last modification time of the file, for change detection, not important in recovery
    }
]
```

It should be decrypted using the following command:

```
$ gpg --batch --yes --passphrase $KEY --output META --decrypt META.enc
```

## blobs/\*

Once you have recovered `META`. You can start decrypting the blobs.

For every entry in `META` run the following command:

```
$ gpg --batch --yes --passphrase $KEY --output $ENTRY_DECRYPTED_NAME --decrypt $ENTRY_ENCRYPTED_NAME
```
