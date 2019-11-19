package utils

import (
    "bytes"
    "crypto/aes"
    "crypto/cipher"
    "log"
)

func init() {
    log.SetFlags(log.Lshortfile | log.LstdFlags)
}

// AesEncrypt 加密函数
func AesEncrypt(plaintext []byte, key []byte) []byte {
    // 分组密钥
    block, err := aes.NewCipher(key)
    if err != nil {
        log.Fatal(" [AES - Encrypt Error]: ", err)
    }
    // 获取秘钥块的长度
    blockSize := block.BlockSize()
    // 填充码
    paddingData := pKCS7Padding(plaintext, blockSize)
    // 加密模式
    blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
    // 创建加密数组
    cyphertext := make([]byte, len(paddingData))
    // 执行加密
    blockMode.CryptBlocks(cyphertext, paddingData)
    return cyphertext
}

// AesDecrypt 解密函数
func AesDecrypt(cryted []byte, key []byte) []byte {
    // 分组秘钥
    block, err := aes.NewCipher(key)
    if err != nil {
        log.Fatal(" [AES - Decrypt Error]: ", err)
    }
    // 获取秘钥块的长度
    blockSize := block.BlockSize()
    // 加密模式
    blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
    // 创建解密数组
    plaintext := make([]byte, len(cryted))
    // 解密
    blockMode.CryptBlocks(plaintext, cryted)
    // 去码
    plaintext = pKCS7UnPadding(plaintext)
    return plaintext
}

// PKCS7Padding 填充码
func pKCS7Padding(ciphertext []byte, blocksize int) []byte {
    padding := blocksize - len(ciphertext)%blocksize
    padtext := bytes.Repeat([]byte{byte(padding)}, padding)
    return append(ciphertext, padtext...)
}

// pKCS7UnPadding 去码
func pKCS7UnPadding(plain []byte) []byte {
    length := len(plain)
    unpadding := int(plain[length-1])
    return plain[:(length - unpadding)]
}
