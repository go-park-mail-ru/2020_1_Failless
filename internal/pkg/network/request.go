package network

import (
	"net/http"
	"strconv"
)

func GetIdFromRequest(w http.ResponseWriter, r *http.Request, ps map[string]string) int64 {
	uid, err := strconv.ParseInt(ps["id"], 10, 64)
	if err != nil {
		GenErrorCode(w, r, MessageInvalidID, http.StatusBadRequest)
		return -1
	}
	return uid
}

func GetEIdFromRequest(w http.ResponseWriter, r *http.Request, ps map[string]string) int64 {
	uid, err := strconv.ParseInt(ps["eid"], 10, 64)
	if err != nil {
		GenErrorCode(w, r, MessageErrorRetrievingEidFromUrl, http.StatusBadRequest)
		return -1
	}
	return uid
}

func GetPageFromRequest(w http.ResponseWriter, r *http.Request, ps *map[string]string) int {
	page, err := strconv.Atoi((*ps)["page"])
	if err != nil {
		return 1
	}
	return page
}
