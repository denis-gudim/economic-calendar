package data

type Country struct {
	Id            int
	Code          string
	ContinentCode string `db:"continent_code"`
	Name          string
	Currency      string
	Translations  Translations
}
