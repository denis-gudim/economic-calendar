package data

type Country struct {
	Id            int    `json:"id" example:"56"`
	Code          string `json:"code" example:"RU"`
	ContinentCode string `db:"continent_code" json:"continentCode" example:"EU"`
	Name          string `json:"name" example:"Russian Federation"`
	Currency      string `json:"currency" example:"RUB"`
}
