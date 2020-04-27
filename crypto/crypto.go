package crypto

import (
	"crypto/md5"
	"crypto/sha512"
	"fmt"
)

var salt = "ringle"

func HashPassword(password string) string {
	h := sha512.New()
	h.Write([]byte(salt))
	hashBytes := h.Sum([]byte(password))

	return Md5Data(hashBytes)
}

func Md5Data(data []byte) string {
	has := md5.Sum(data)
	md5str1 := fmt.Sprintf("%x", has) //将[]byte转成16进制
	return md5str1
}

func GenerateOpenIdForImiFromUserId(userId uint64) string {
	id := fmt.Sprintf("%v", userId)

	h := sha512.New()
	h.Write([]byte(salt))
	hashBytes := h.Sum([]byte(id))

	return "wanzi" + Md5Data(hashBytes)[:27]
}

func HashDeviceId(uuid string) string {

	h := sha512.New()
	h.Write([]byte(salt))
	hashBytes := h.Sum([]byte(uuid))

	return Md5Data(hashBytes)
}
