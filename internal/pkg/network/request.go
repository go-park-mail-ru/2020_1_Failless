package network

import (
	"net/http"
	"strconv"
)

func GetIdFromRequest(w http.ResponseWriter, r *http.Request, ps *map[string]string) int {
	uid, err := strconv.Atoi((*ps)["id"])
	if err != nil {
		GenErrorCode(w, r, "Incorrect id", http.StatusBadRequest)
		return -1
	}
	return uid
}
