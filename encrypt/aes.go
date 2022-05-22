package encrypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"errors"
	"fmt"
	"log"
)

var key []byte
var block cipher.Block
var encrypter cipher.BlockMode
var decrypter cipher.BlockMode
var BlockSize int

func RestKey(token string) {
	sum256 := sha256.Sum256([]byte(token))
	key = sum256[:]
	log.Printf("init key length(%d)\n", len(key))
	tmp, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	block = tmp
	BlockSize = block.BlockSize()
	encrypter = cipher.NewCBCEncrypter(block, key[:BlockSize])
	decrypter = cipher.NewCBCDecrypter(block, key[:BlockSize])
}

func init() {
	RestKey("1234567890123456")
}

func AesEncryptData(data []byte) []byte {
	encryptBytes := pkcs7Padding(data, BlockSize)
	encrypted := make([]byte, len(encryptBytes))
	encrypter.CryptBlocks(encrypted, encryptBytes)
	return encrypted
}

func AesDecryptData(encrypted []byte) ([]byte, error) {
	length := len(encrypted)
	if length%BlockSize != 0 {
		return make([]byte, 0), errors.New(fmt.Sprintf("Decrypt length(%d) error", length))
	}
	paddingBytes := make([]byte, length)
	decrypter.CryptBlocks(paddingBytes, encrypted)
	data, err := pkcs7UnPadding(paddingBytes)
	if err != nil {
		return make([]byte, 0), err
	}
	return data, nil
}

func pkcs7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

func pkcs7UnPadding(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return make([]byte, 0), errors.New("UnPadding error length(0)")
	}
	unPadding := int(data[length-1])
	index := length - unPadding
	if index < 0 {
		return make([]byte, 0), errors.New(fmt.Sprintf("UnPadding error length(%d),unPadding(%d)", length, unPadding))
	}
	return data[:index], nil
}
