package data

type Country struct {
	Id               int
	Code             string
	ContinentCode    string
	Name             string
	Currency         string
	NameTranslations Translations
}
