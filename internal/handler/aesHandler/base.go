package aesHandler

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"fmt"
)

// PKCS7Padding 填充
func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// PKCS7UnPadding 去除填充
func PKCS7UnPadding(origData []byte) ([]byte, error) {
	length := len(origData)
	if length <= 0 {
		// 处理错误或返回一个错误提示
		return nil, fmt.Errorf("invalid data length")
	}
	unpadding := int(origData[length-1])
	if unpadding > length {
		return nil, fmt.Errorf("unpadding error")
	}
	return origData[:(length - unpadding)], nil
}

// ZeroPadding 填充
func ZeroPadding(origData []byte, blockSize int) []byte {
	padding := blockSize - len(origData)%blockSize
	padtext := bytes.Repeat([]byte{0}, padding)
	return append(origData, padtext...)
}

// ZeroUnPadding 去除填充
func ZeroUnPadding(origData []byte) ([]byte, error) {
	return bytes.TrimRightFunc(origData, func(r rune) bool {
		return r == rune(0)
	}), nil
}

// // AesEncryptECB AES-ECB加密
// func AesEncrypt(origData, key []byte) ([]byte, error) {
// 	if len(key) != 16 {
// 		return nil, fmt.Errorf("invalid key size for AES-128: %d bytes", len(key))
// 	}
// 	block, err := aes.NewCipher(key)
// 	if err != nil {
// 		return nil, err
// 	}
// 	blockSize := block.BlockSize()
// 	origData = PKCS7Padding(origData, blockSize)

// 	crypted := make([]byte, len(origData))
// 	for bs, be := 0, blockSize; bs < len(origData); bs, be = bs+blockSize, be+blockSize {
// 		block.Encrypt(crypted[bs:be], origData[bs:be])
// 	}

// 	return crypted, nil
// }

// // AesDecryptECB AES-ECB解密
// func AesDecrypt(crypted, key []byte) ([]byte, error) {
// 	if len(key) != 16 {
// 		return nil, fmt.Errorf("invalid key size for AES-128: %d bytes", len(key))
// 	}
// 	block, err := aes.NewCipher(key)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if len(crypted)%aes.BlockSize != 0 {
// 		return nil, fmt.Errorf("crypted data is not a multiple of the block size")
// 	}

// 	origData := make([]byte, len(crypted))
// 	for bs, be := 0, aes.BlockSize; bs < len(crypted); bs, be = bs+aes.BlockSize, be+aes.BlockSize {
// 		block.Decrypt(origData[bs:be], crypted[bs:be])
// 	}

// 	return PKCS7UnPadding(origData)
// }

// // AesEncrypt AES加密，CBC-128 不使用IV，应用ZeroPadding
// func AesEncrypt(origData, key []byte) ([]byte, error) {
// 	if len(key) != 16 {
// 		return nil, fmt.Errorf("invalid key size for AES-128: %d bytes", len(key))
// 	}
// 	block, err := aes.NewCipher(key)
// 	if err != nil {
// 		return nil, err
// 	}
// 	blockSize := block.BlockSize()
// 	origData = ZeroPadding(origData, blockSize)

// 	// 使用全零向量
// 	iv := bytes.Repeat([]byte{0}, blockSize)

// 	blockMode := cipher.NewCBCEncrypter(block, iv)
// 	crypted := make([]byte, len(origData))
// 	blockMode.CryptBlocks(crypted, origData)
// 	return crypted, nil
// }

// // AesDecrypt AES解密，CBC-128 不使用IV，应用ZeroPadding
// func AesDecrypt(crypted, key []byte) ([]byte, error) {
// 	if len(key) != 16 {
// 		return nil, fmt.Errorf("invalid key size for AES-128: %d bytes", len(key))
// 	}
// 	block, err := aes.NewCipher(key)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if len(crypted)%aes.BlockSize != 0 {
// 		return nil, fmt.Errorf("crypted data is not a multiple of the block size")
// 	}

// 	// 使用全零向量
// 	iv := bytes.Repeat([]byte{0}, aes.BlockSize)

// 	blockMode := cipher.NewCBCDecrypter(block, iv)
// 	origData := make([]byte, len(crypted))
// 	blockMode.CryptBlocks(origData, crypted)

// 	// 去除填充
// 	return ZeroUnPadding(origData)
// }

// AesEncrypt AES加密 CBC-128
func AesEncrypt(origData, key []byte, iv []byte) ([]byte, error) {
	if len(key) != 16 {
		return nil, fmt.Errorf("invalid key size for AES-128: %d bytes", len(key))
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	origData = PKCS7Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, iv)
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

// AesDecrypt AES解密 CBC-128
func AesDecrypt(crypted, key []byte, iv []byte) ([]byte, error) {
	if len(key) != 16 {
		return nil, fmt.Errorf("invalid key size for AES-128: %d bytes", len(key))
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, iv)
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData, err = PKCS7UnPadding(origData)
	if err != nil {
		return nil, err
	}
	return origData, nil
}
