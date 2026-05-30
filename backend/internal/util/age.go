package util

import "time"

func CalculateAge(birthDate time.Time) int {
	today := time.Now()

	age := today.Year() - birthDate.Year()

	todayMonthDay := int(today.Month())*100 + today.Day()
	birthMonthDay := int(birthDate.Month())*100 + birthDate.Day()

	if todayMonthDay < birthMonthDay {
		age--
	}

	return age
}

func isLeapYear(year int) bool {
	return (year%4 == 0 && year%100 != 0) || year%400 == 0
}

func CalculateDateRangeByAge(minAge, maxAge int) (minDate, maxDate time.Time) {
	today := time.Now()

	minDate = time.Date(
		today.Year()-maxAge,
		today.Month(),
		today.Day(),
		0, 0, 0, 0, today.Location(),
	)

	maxDate = time.Date(
		today.Year()-minAge,
		today.Month(),
		today.Day(),
		0, 0, 0, 0, today.Location(),
	)

	return minDate, maxDate
}
