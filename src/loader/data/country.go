package data

type Country struct {
	Code          string
	ContinentCode string `db:"continent_code"`
	Name          string
	Currency      string
	InvestingId   *int `db:"investing_id"`
	Translations  map[string]string
}
