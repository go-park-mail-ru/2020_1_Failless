package repository

import (
	"failless/internal/pkg/models"
	"time"
)

type ProfileInfo struct {
	About     *string               `json:"about"`
	Photos    *[]string             `json:"photos"`
	Rating    *float32              `json:"rating"`
	Birthday  *time.Time            `json:"birthday"`
	Gender    *int                  `json:"gender"`
	LoginDate *time.Time            `json:"login_date"`
	Location  *models.LocationPoint `json:"location"`
}
