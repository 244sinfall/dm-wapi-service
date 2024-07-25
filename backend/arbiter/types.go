package arbiter

const (
	giveXP   string = "givexp"
	takeXP   string = "takexp"
	giveGold string = "givegold"

	giveXPCommand   string = ".exp game"
	takeXPCommand   string = ".exp oth"
	giveGoldCommand string = ".send mo"

	xpToRateModifier   float64 = 1000
	goldToRateModifier float64 = 6000

	writerModifier           float64 = 1.5
	masterModifier           float64 = 1
	masterAndWriterModifier  float64 = 1.5
	crafterModifier          float64 = 0.5
	crafterAndWriterModifier float64 = 1
)

type arbiterCommandsRequest struct {
	ParticipantsCleanedText string `json:"participantsCleanedText" binding:"required"`
	Mode                    string `json:"mode" binding:"required"` // "givexp", "takexp", "givegold"
	Rate                    int    `json:"rate" binding:"required"`
	EventLink               string `json:"eventLink" binding:"required"`
}

type arbiterCommandsResponse struct {
	Commands             string `json:"commands"`
	ParticipantsModified string `json:"participantsModified"`
}
