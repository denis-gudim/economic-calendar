package data

type InvestingCountry struct {
	Id         int
	Title      string
	LanguageId int
}

func (country *InvestingCountry) GetId() int {
	return country.Id
}

func (country *InvestingCountry) GetLanguageId() int {
	return country.LanguageId
}
