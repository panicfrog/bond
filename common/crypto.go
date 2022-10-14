package common

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"errors"
)

func pkcs5Padding(cipherText []byte, blockSize int) []byte {
	padding := blockSize - len(cipherText)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(cipherText, padText...)
}

func pkcs5UnPadding(origData []byte) []byte {
	length := len(origData)
	unPadding := int(origData[length-1])
	return origData[:(length - unPadding)]
}

// AesEncryptCBC  key 的长度为 16 24 32
func AesEncryptCBC(originData []byte, key []byte) (encrypted []byte, err error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}
	blockSize := block.BlockSize()
	origData := pkcs5Padding(originData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	encrypted = make([]byte, len(origData))
	blockMode.CryptBlocks(encrypted, origData)
	return
}

// AesDecryptCBC  key 的长度为 16 24 32
func AesDecryptCBC(encrypted []byte, key []byte) (decrypted []byte, err error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	decrypted = make([]byte, len(encrypted))
	if len(decrypted) < blockSize {
		err = errors.New("encrypted is not full blocks")
		return
	}
	blockMode.CryptBlocks(decrypted, encrypted)
	decrypted = pkcs5UnPadding(decrypted)
	return
}
