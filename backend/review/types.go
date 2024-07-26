package review

// Review
type reviewRate struct {
	RateName  string `json:"rateName"`
	RateValue int    `json:"rateValue"`
}

type review struct {
	Rates           []reviewRate `json:"rates" binding:"required"`
	TotalRate       int          `json:"totalRate"`
	CharName        string       `json:"charName" binding:"required"`
	ReviewerProfile string       `json:"reviewerProfile" binding:"required"`
	ReviewerDiscord string       `json:"reviewerDiscord" binding:"required"`
}

type reviewOutput struct {
	Review string `json:"review"`
}

// Lottery
const (
	epicItem       int = 72
	rareItem       int = 24
	unusualItem    int = 6
	usualItem      int = 2
	lowQualityItem int = 1

	minPotentialWinnersPart     = 0.15
	maxPotentialWinnersPart     = 0.35
	potentialWinnerStep         = maxPotentialWinnersPart - minPotentialWinnersPart
	minRate                 int = 7
	minParticipants         int = 10
)

type lotteryResult struct {
	Epic    int `json:"epic"`
	Rare    int `json:"rare"`
	Unusual int `json:"unusual"`
	Usual   int `json:"usual"`
	Low     int `json:"low"`
}

type lotteryResponse struct {
	ParticipantsCount  int           `json:"participantsCount"`
	BankPerParticipant float64       `json:"bankPerParticipant"`
	Bank               float64       `json:"bank"`
	BankRemain         float64       `json:"bankRemain"`
	PotentialWinners   int           `json:"potentialWinners"`
	Lottery            lotteryResult `json:"lottery"`
}

type lotteryOptions struct {
	ParticipantsCount       int  `json:"participantsCount"  binding:"required"`
	Rate                    int  `json:"rate"  binding:"required"`
	QualityOverQuantityMode bool `json:"qualityOverQuantityMode"`
}
