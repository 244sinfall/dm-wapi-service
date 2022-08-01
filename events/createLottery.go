package events

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"math"
	"math/rand"
	"net/http"
)

type itemPrice int

const (
	epicItem       itemPrice = 72
	rareItem       itemPrice = 24
	unusualItem    itemPrice = 6
	usualItem      itemPrice = 2
	lowQualityItem itemPrice = 1

	minPotentialWinnersPart = 0.15
	maxPotentialWinnersPart = 0.35
	potentialWinnerStep     = maxPotentialWinnersPart - minPotentialWinnersPart
)

func (i itemPrice) getItemName() string {
	switch i {
	case 72:
		return "epic"
	case 24:
		return "rare"
	case 6:
		return "unusual"
	case 2:
		return "usual"
	case 1:
		return "low"
	default:
		return "undefined"
	}
}

type lotteryCreatorError string

func (l lotteryCreatorError) Error() string {
	return string(l)
}

type RespondLotteryObject struct {
	ParticipantsCount  int            `json:"participantsCount"`
	BankPerParticipant float64        `json:"bankPerParticipant"`
	Bank               float64        `json:"bank"`
	BankRemain         float64        `json:"bankRemain"`
	PotentialWinners   int            `json:"potentialWinners"`
	Lottery            map[string]int `json:"lottery"`
}

type LotteryCreator struct {
	ParticipantsCount       int  `json:"participantsCount"`
	Rate                    int  `json:"rate"`
	QualityOverQuantityMode bool `json:"qualityOverQuantityMode"`
}

func (r *RespondLotteryObject) generateLotteryObject(qualityOverQuantityMode bool) {
	if qualityOverQuantityMode {
		fmt.Println(r.BankRemain)
		for r.BankRemain > 1 {
			fmt.Println(r.BankRemain, r.Lottery, "quality on")
			switch {
			case r.BankRemain >= float64(epicItem):
				r.Lottery[epicItem.getItemName()]++
				r.BankRemain -= float64(epicItem)
			case r.BankRemain >= float64(rareItem):
				r.Lottery[rareItem.getItemName()]++
				r.BankRemain -= float64(rareItem)
			case r.BankRemain >= float64(unusualItem):
				r.Lottery[unusualItem.getItemName()]++
				r.BankRemain -= float64(unusualItem)
			case r.BankRemain >= float64(usualItem):
				r.Lottery[usualItem.getItemName()]++
				r.BankRemain -= float64(usualItem)
			case r.BankRemain >= float64(lowQualityItem):
				r.Lottery[lowQualityItem.getItemName()]++
				r.BankRemain -= float64(lowQualityItem)
			default:
				break
			}
		}
	} else {
		potentialWinners := r.PotentialWinners
		fmt.Println(r.BankRemain, r.PotentialWinners, "quality off")
	potentialWinnersLoop:
		for r.PotentialWinners > 0 {
			bankForParticipantLeft := r.BankRemain / float64(potentialWinners)
			fmt.Println(r.BankRemain, r.Lottery)
			switch {
			case r.BankRemain > float64(epicItem) && bankForParticipantLeft >= 9:
				r.Lottery[epicItem.getItemName()]++
				r.BankRemain -= float64(epicItem)
				potentialWinners--
			case r.BankRemain > float64(rareItem) && bankForParticipantLeft >= 4.5:
				r.Lottery[rareItem.getItemName()]++
				r.BankRemain -= float64(rareItem)
				potentialWinners--
			case r.BankRemain > float64(unusualItem) && bankForParticipantLeft >= 2:
				r.Lottery[unusualItem.getItemName()]++
				r.BankRemain -= float64(unusualItem)
				potentialWinners--
			case r.BankRemain > float64(usualItem) && bankForParticipantLeft >= 1.2:
				r.Lottery[usualItem.getItemName()]++
				r.BankRemain -= float64(usualItem)
				potentialWinners--
			default:
				break potentialWinnersLoop
			}
		}
		if r.BankRemain > 1 {
			r.generateLotteryObject(true)
		}
	}
}

func (l LotteryCreator) generateLottery() RespondLotteryObject {
	var respond RespondLotteryObject
	respond.ParticipantsCount = l.ParticipantsCount
	respond.BankPerParticipant = getBankPerParticipant(respond.ParticipantsCount)
	respond.Bank = float64(l.ParticipantsCount) * respond.BankPerParticipant
	respond.BankRemain = respond.Bank
	respond.PotentialWinners = int(math.Round(
		float64(respond.ParticipantsCount) * (maxPotentialWinnersPart - (potentialWinnerStep * rand.Float64()))))
	respond.Lottery = map[string]int{}
	respond.generateLotteryObject(l.QualityOverQuantityMode)
	return respond
}

func getBankPerParticipant(rate int) float64 {
	switch {
	case rate >= 13:
		return 4.5
	case rate >= 10 && rate <= 12:
		return 3
	case rate >= 7 && rate <= 9:
		return 1.5
	default:
		return 0
	}

}

func CreateLottery(c *gin.Context) {
	var lotteryCreator LotteryCreator
	if err := c.BindJSON(&lotteryCreator); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}
	if lotteryCreator.Rate < 7 {
		c.AbortWithError(http.StatusBadRequest, lotteryCreatorError("Lottery could be created for rate > 7"))
		return
	}
	if lotteryCreator.ParticipantsCount < 10 {
		c.AbortWithError(http.StatusBadRequest, lotteryCreatorError("Lottery could be created for >= 10 participants."))
	}
	respond := lotteryCreator.generateLottery()
	c.JSON(http.StatusOK, respond)
}
