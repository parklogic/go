package value

var validTypes = map[string]bool{
	"2048": true,
	"P256": true,
}

func IsValidPrivateKeyType(t string) bool {
	return validTypes[t]
}
