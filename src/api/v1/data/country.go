package data

type Country struct {
	Id            int    `json:"id"`
	Code          string `json:"code"`
	ContinentCode string `db:"continent_code" json:"continentCode"`
	Name          string `json:"name"`
	Currency      string `json:"currency"`
}
