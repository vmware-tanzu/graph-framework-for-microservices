package common

var Mode string

func IsModeAdmin() bool {
	return Mode == "admin"
}
