package email

import "net/http"

type Delivery interface {
	SendReminder(w http.ResponseWriter, r *http.Request, _ map[string]string)
}