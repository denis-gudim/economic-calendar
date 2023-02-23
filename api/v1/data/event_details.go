package data

type EventDetails struct {
	Event
	Overview  string `json:"overview"`
	Source    string `json:"source"`
	SourceUrl string `db:"source_url" json:"sourceUrl"`
}
