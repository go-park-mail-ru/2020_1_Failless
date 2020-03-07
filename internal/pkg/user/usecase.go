package user


// Get gender string by int id
func GenderById(genderId int) string {
	switch genderId {
	case Male:
		return "male"
	case Female:
		return "female"
	}
	return "other"
}

// Get gender id by string name
func GenderByStr(gender string) int {
	switch gender {
	case "male":
		return Male
	case "female":
		return Female
	}

	// Other gender
	return Other
}
