package util

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
)

func Encode(plain []byte) string {
	ranStr := GetRanStr(16)

	tmp1 := md5.Sum([]byte(ranStr))
	tmp2 := md5.Sum(tmp1[:])

	var iv, key []byte
	iv = tmp1[:]
	key = append(key, tmp1[:]...)
	key = append(key, tmp2[:]...)

	encryptData := AES_CBC_Encrypt(plain, key, iv)
	x := hex.EncodeToString(encryptData)

	return x + ranStr
}

// mixStr = hex(encrypted)+ ranStr(16)
func Decoded(mixStr string) ([]byte, error) {
	if len(mixStr) <= 16 {
		return nil, errors.New("mixstr is invalid")
	}
	encryptedStr := mixStr[:len(mixStr)-16]
	ranStr := mixStr[len(mixStr)-16:]

	var key, iv []byte
	tmp1 := md5.Sum([]byte(ranStr))
	tmp2 := md5.Sum(tmp1[:])
	iv = tmp1[:]
	key = append(key, tmp1[:]...)
	key = append(key, tmp2[:]...)

	encrypted, err := hex.DecodeString(encryptedStr)
	if err != nil {
		return nil, err
	}
	plain, err := AES_CBC_Decrypt(encrypted, key, iv)
	return plain, err
}
