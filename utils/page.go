package utils

import "strconv"

// PageStringStart return the start index for this page, default page is 1 if page
// string is empty or isn't a legal number or less than 1
func PageStringStart(page string, countPerPage int) int {
	if page == "" {
		return 0
	}

	val, err := strconv.Atoi(page)
	if err != nil {
		return 0
	}

	return PageStart(val, countPerPage)
}

// PageStart return the start index for this page, default page is 1 if the
// page less than 1
func PageStart(page, countPerPage int) int {
	if page <= 0 {
		return 0
	}

	return (page - 1) * countPerPage
}
