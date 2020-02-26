package routes

import (
	htmux "github.com/dimfeld/httptreemux"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

func ComparePasswords(hash []byte, p string) bool {
	err := bcrypt.CompareHashAndPassword(hash, []byte(p))
	if err != nil {
		return false
	}
	return true
}

func createNewUser(writer http.ResponseWriter, request *http.Request, ps map[string]string) {
}

func getUserProfile(writer http.ResponseWriter, request *http.Request, ps map[string]string) {
}

func updateUserProfile(writer http.ResponseWriter, request *http.Request, ps map[string]string) {
}

func UserHandler(router *htmux.TreeMux) {
	router.POST("/api/user/create", createNewUser)
	router.GET("/api/user/profile/:id", getUserProfile)
	router.POST("/api/user/profile", updateUserProfile)
}
