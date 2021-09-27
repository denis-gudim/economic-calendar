package data

type Language struct {
	Code       string
	Name       string
	NativeName string `db:"native_name"`
}
