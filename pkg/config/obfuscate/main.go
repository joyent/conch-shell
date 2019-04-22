// Copyright Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package obfuscate

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
)

func Obfuscate(data string, key string) (string, error) {

	// We need a 32 byte key which sha256 is happy to give us
	sum := sha256.Sum256([]byte(key))
	c, err := aes.NewCipher(sum[:])

	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())

	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	b := gcm.Seal(nonce, nonce, []byte(data), nil)
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func Deobfuscate(text string, key string) (string, error) {

	data, err := base64.RawURLEncoding.DecodeString(text)
	if err != nil {
		return "", err
	}

	sum := sha256.Sum256([]byte(key))
	c, err := aes.NewCipher(sum[:])
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("ciphertext smaller than nonce")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]

	plain, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}
	return string(plain), nil
}
