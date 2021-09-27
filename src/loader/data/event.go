package data

type Event struct {
	Id           int
	CountryCode  string
	ImpactLevel  int
	Unit         string
	Source       string
	SourceUrl    string
	Translations map[string]string
}
