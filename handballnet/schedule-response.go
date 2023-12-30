package handballnet

type scheduleResponse struct {
	PageProps leaguePageProps `json:"pageProps"`
}

type leaguePageProps struct {
	Schedule schedule `json:"schedule"`
}

type schedule struct {
	Data []game `json:"data"`
}

type game struct {
	ID            string     `json:"id"`
	Tournament    tournament `json:"tournament"`
	Phase         phase      `json:"phase"`
	Round         round      `json:"round"`
	Field         field      `json:"field"`
	HomeTeam      team       `json:"homeTeam"`
	AwayTeam      team       `json:"awayTeam"`
	StartsAt      int64      `json:"startsAt"`
	ShowTime      bool       `json:"showTime"`
	State         string     `json:"state"`
	HomeGoals     int        `json:"homeGoals"`
	AwayGoals     int        `json:"awayGoals"`
	HomeGoalsHalf int        `json:"homeGoalsHalf"`
	AwayGoalsHalf int        `json:"awayGoalsHalf"`
}

type tournament struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Acronym        string `json:"acronym"`
	TournamentType string `json:"tournamentType"`
	AgeGroup       string `json:"ageGroup"`
}

type phase struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Acronym string `json:"acronym"`
	Type    string `json:"type"`
}

type round struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Acronym  string `json:"acronym"`
	Number   int    `json:"number"`
	StartsAt int64  `json:"startsAt"`
	EndsAt   int64  `json:"endsAt"`
	Type     string `json:"type"`
}

type field struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Acronym string `json:"acronym"`
	City    string `json:"city"`
	Type    string `json:"type"`
}

type team struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Acronym     string `json:"acronym"`
	Logo        string `json:"logo"`
	TeamGroupId int    `json:"teamGroupId"`
	Type        string `json:"type"`
}
