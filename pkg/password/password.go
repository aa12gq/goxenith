package password

import (
	"crypto/md5"
	"encoding/hex"
	"golang.org/x/crypto/bcrypt"
)

func EncryptPassword(password string) (encrypted string) {
	sum := md5.Sum([]byte(password))
	return hex.EncodeToString(sum[:])
}

// BcryptPassword 使用 bcrypt 对密码进行加密 随机盐值和多次迭代来增强密码的安全性。
func BcryptPassword(password string) string {
	// GenerateFromPassword 的第二个参数是 cost 值。建议大于 12，数值越大耗费时间越长
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 4)
	if err != nil {
		return ""
	}
	return string(bytes)
}

// BcryptPasswordMatch 对比明文密码和数据库的哈希值
func BcryptPasswordMatch(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// BcryptIsHashed 判断字符串是否是哈希过的数据
func BcryptIsHashed(str string) bool {
	// bcrypt 加密后的长度等于 60
	return len(str) == 60
}
