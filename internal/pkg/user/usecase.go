package user

import "failless/internal/pkg/models"

// Get gender string by int id
func GenderById(genderId int) string {
	switch genderId {
	case models.Male:
		return "male"
	case models.Female:
		return "female"
	}
	return "other"
}

// Get gender id by string name
func GenderByStr(gender string) int {
	switch gender {
	case "male":
		return models.Male
	case "female":
		return models.Female
	}

	// Other gender
	return models.Other
}
