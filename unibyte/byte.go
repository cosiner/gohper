package unibyte

import "unicode"

// IsLower check letter is lower case or not
func IsLower(b byte) bool {
	return b >= 'a' && b <= 'z'
}

// IsUpper check letter is upper case or not
func IsUpper(b byte) bool {
	return b >= 'A' && b <= 'Z'
}

// IsLetter check character is a letter or not
func IsLetter(b byte) bool {
	return IsLower(b) || IsUpper(b)
}

// IsSpaceQuote return wehter a byte is space or quote characters
func IsSpaceQuote(b byte) bool {
	return IsSpace(b) || b == '"' || b == '\''
}

// IsSpace only call unicode.IsSpace
func IsSpace(b byte) bool {
	return unicode.IsSpace(rune(b))
}

// ToLower convert a byte to lower case
func ToLower(b byte) byte {
	if IsUpper(b) {
		b = b - 'A' + 'a'
	}

	return b
}

// ToUpper convert a byte to upper case
func ToUpper(b byte) byte {
	if IsLower(b) {
		b = b - 'a' + 'A'
	}

	return b
}

// ToLowerString convert a byte to lower case string
func ToLowerString(b byte) string {
	if IsUpper(b) {
		b = b - 'A' + 'a'
	}

	return string(b)
}

// ToUpperString convert a byte to upper case string
func ToUpperString(b byte) string {
	if IsLower(b) {
		b = b - 'a' + 'A'
	}

	return string(b)
}
