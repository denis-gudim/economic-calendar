package data

type Event struct {
	Id                   int
	CountryId            int
	ImpactLevel          int
	Unit                 string
	Source               string
	SourceUrl            string
	TitleTranslations    Translations
	OverviewTranslations Translations
}
