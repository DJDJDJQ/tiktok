package utils

import (
	"encoding/base64"
	"log"

	"golang.org/x/crypto/bcrypt"
)

//base64加密方式
func base64Encode(src []byte) []byte {
	return []byte(base64.StdEncoding.EncodeToString(src))
}

//base64解密方式
func base64Decode(src []byte) ([]byte, error) {
	return base64.StdEncoding.DecodeString(string(src))
}

// StrEncrypt 对传入字符串进行加密
func StrEncrypt(str string, salt string) string {
	str = str + salt
	hash, err := bcrypt.GenerateFromPassword([]byte(str), bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}

// StrMatch 对传入的加密字符串进行比对,str2为明文
func StrMatch(str1 string, str2 string, salt string) bool {
	str2 = str2 + salt
	err := bcrypt.CompareHashAndPassword([]byte(str1), []byte(str2))
	return err == nil
}
