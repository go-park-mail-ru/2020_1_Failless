package security

import "golang.org/x/crypto/bcrypt"

func ComparePasswords(hash []byte, p string) bool {
	err := bcrypt.CompareHashAndPassword(hash, []byte(p))
	return err == nil
}
