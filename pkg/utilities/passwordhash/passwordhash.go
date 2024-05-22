package passwordhash

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
)

func HashPassword(salt string, password string) string {
	md5Hash := md5.Sum([]byte(password))
	md5Password := hex.EncodeToString(md5Hash[:])

	saltedPassword := md5Password + salt
	sha256Hash := sha256.Sum256([]byte(saltedPassword))
	sha256Password := hex.EncodeToString(sha256Hash[:])

	return sha256Password
}
