package investing

type InvestingCalendarEvent struct {
	Id         int
	Title      string
	Overview   string
	Source     string
	SourceUrl  string
	Unit       string
	LanguageId int
}

func (event *InvestingCalendarEvent) GetId() int {
	return event.Id
}

func (event *InvestingCalendarEvent) GetLanguageId() int {
	return event.LanguageId
}
