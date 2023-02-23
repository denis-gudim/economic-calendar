package data

type Event struct {
	EventRow
	Type        int    `json:"type"`
	ImpactLevel int    `db:"impact_level" json:"impactLevel"`
	Code        string `json:"code"`
	Unit        string `json:"unit"`
	Title       string `json:"title"`
}
