package security

import "golang.org/x/crypto/bcrypt"

func ComparePasswords(hash []byte, p string) bool {
	err := bcrypt.CompareHashAndPassword(hash, []byte(p))
	return err == nil
}

func EncryptPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}
// TODO: add to header CSRF-token when create request for up vote and down vote and all post|put|delete requests
// login | register pages not needed in csrf token
// and if we tape reload in the browser our SPA have to still work fine
// we get token without session, using separate method for this (get token and after this go to handler)
// we may save token to local storage in browser but not into page memory

