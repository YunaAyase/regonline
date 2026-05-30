package util

// NormalizeIDNumber removes whitespace and converts lowercase x to uppercase X.
func NormalizeIDNumber(id string) string {
	result := make([]byte, 0, len(id))
	for i := 0; i < len(id); i++ {
		c := id[i]
		if c == ' ' || c == '\t' || c == '\n' || c == '\r' {
			continue
		}
		if i == len(id)-1 && c == 'x' {
			result = append(result, 'X')
		} else {
			result = append(result, c)
		}
	}
	return string(result)
}

// ValidateIDNumberLength checks the ID number is exactly 18 characters.
func ValidateIDNumberLength(id string) bool {
	return len(id) == 18
}

// ValidateBirthDateMatchesID checks if the birth date (positions 7-14) matches the given date.
func ValidateBirthDateMatchesID(birthDate string, idNumber string) bool {
	if len(idNumber) < 14 {
		return false
	}
	birthFromID := idNumber[6:14]
	return birthFromID == birthDate
}

// ValidateGenderMatchesID checks if the gender matches the 17th digit of the ID number.
// Odd digit (1,3,5,7,9) = male, even digit (0,2,4,6,8) = female.
func ValidateGenderMatchesID(gender string, idNumber string) bool {
	if len(idNumber) < 17 {
		return false
	}
	genderDigit := idNumber[16] - '0'
	isMale := genderDigit%2 == 1

	switch gender {
	case "男":
		return isMale
	case "女":
		return !isMale
	default:
		return false
	}
}
