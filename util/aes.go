package util

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"errors"
)

// 对明文进行填充
func pkcs5Padding(plainText []byte, blockSize int) []byte {
	//计算要填充的长度
	n := blockSize - len(plainText)%blockSize
	//对原来的明文填充n个n
	temp := bytes.Repeat([]byte{byte(n)}, n)
	plainText = append(plainText, temp...)
	return plainText
}

// 对明文进行填充
func zeroPadding(plainText []byte, blockSize int) []byte {
	//计算要填充的长度
	n := blockSize - len(plainText)%blockSize
	//对原来的明文填充n个n
	temp := bytes.Repeat([]byte{byte(0)}, n)
	plainText = append(plainText, temp...)
	return plainText
}

// 对密文删除填充
func pkcs5UnPadding(cipherText []byte) ([]byte, error) {
	//取出密文最后一个字节end
	end := cipherText[len(cipherText)-1]
	if int(end) < 0 || int(end) > len(cipherText) {
		return cipherText, errors.New("invalid end")
	}
	//删除填充
	cipherText = cipherText[:len(cipherText)-int(end)]
	return cipherText, nil
}

// AEC加密（CBC模式） iv 16 字节,key 可以 16 24 32字节
func AES_CBC_Encrypt(plainText []byte, key, iv []byte) []byte {
	//指定加密算法，返回一个AES算法的Block接口对象
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	//进行PKCS5填充
	plainText = pkcs5Padding(plainText, block.BlockSize())
	// fmt.Println(len(plainText), string(plainText))
	//指定初始向量iv,长度和block的块尺寸一致
	//指定分组模式，返回一个BlockMode接口对象
	blockMode := cipher.NewCBCEncrypter(block, iv)
	//加密连续数据库
	cipherText := make([]byte, len(plainText))
	blockMode.CryptBlocks(cipherText, plainText)
	//返回密文
	return cipherText
}

// AEC解密（CBC模式）
func AES_CBC_Decrypt(cipherText []byte, key, iv []byte) ([]byte, error) {
	//校验
	if len(cipherText) == 0 {
		return nil, errors.New("src is nil")
	}

	//指定解密算法，返回一个AES算法的Block接口对象
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	//校验
	if len(cipherText)%block.BlockSize() != 0 {
		return nil, errors.New("crypto/cipher: input not full blocks")
	}

	//指定分组模式，返回一个BlockMode接口对象
	blockMode := cipher.NewCBCDecrypter(block, iv)

	plainText := make([]byte, len(cipherText))
	//解密
	blockMode.CryptBlocks(plainText, cipherText)

	//删除填充
	plainText, err = pkcs5UnPadding(plainText)
	return plainText, err
}
