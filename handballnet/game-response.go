package handballnet

type gameResponse struct {
	PageProps gamePageProps `json:"pageProps"`
}

type gamePageProps struct {
	Game   game    `json:"game"`
	Lineup lineup  `json:"lineup"`
	Events []event `json:"events"`
}

type lineup struct {
	Home          []player      `json:"home"`
	HomeOfficials []staffMember `json:"homeOfficials"`
	Away          []player      `json:"away"`
	AwayOfficials []staffMember `json:"awayOfficials"`
	Referees      []referee     `json:"referees"`
}

type player struct {
	ID            string `json:"id"`
	Firstname     string `json:"firstname"`
	Lastname      string `json:"lastname"`
	Position      string `json:"position"`
	Number        int    `json:"number"`
	Goals         uint   `json:"goals"`
	PenaltyGoals  uint   `json:"penaltyGoals"`
	PenaltyMissed uint   `json:"penaltyMissed"`
	Penalties     uint   `json:"penalties"`
	YellowCards   uint   `json:"yellowCards"`
	RedCards      uint   `json:"redCards"`
	BlueCards     uint   `json:"blueCards"`
}

type staffMember struct {
	ID                            string   `json:"id"`
	Firstname                     string   `json:"firstname"`
	Lastname                      string   `json:"lastname"`
	Position                      string   `json:"position"`
	TimePenalties                 []uint   `json:"timePenalties"`
	Warnings                      []string `json:"warnings"`
	Disqualifications             []string `json:"disqualifications"`
	DisqualificationWithBlueCards []string `json:"disqualificationWithBlueCards"`
}

type referee struct {
	ID        string `json:"id"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Position  string `json:"position"`
}

type eventType string

const (
	eventTypeStopPeriod                   eventType = "StopPeriod"
	eventTypeGoal                         eventType = "Goal"
	eventTypeSevenMeterGoal               eventType = "SevenMeterGoal"
	eventTypeTwoMinutePenalty             eventType = "TwoMinutePenalty"
	eventTypeTimeout                      eventType = "Timeout"
	eventTypeSevenMeterMissed             eventType = "SevenMeterMissed"
	eventTypeWarning                      eventType = "Warning"
	eventTypeDisqualification             eventType = "Disqualification"
	eventTypeDisqualificationWithBlueCard eventType = "DisqualificationWithBlueCard"
)

type event struct {
	ID        uint      `json:"id"`
	Type      eventType `json:"type"`
	Time      string    `json:"time"`
	Score     string    `json:"score"`
	Timestamp int64     `json:"timestamp"`
	Team      string    `json:"team"`
	Message   string    `json:"message"`
}
