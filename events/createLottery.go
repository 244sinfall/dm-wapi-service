package events

import (
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

type lotteryCreatorError string

func (l lotteryCreatorError) Error() string {
	return string(l)
}

type LotteryObject struct {
	Epic    int `json:"epic"`
	Rare    int `json:"rare"`
	Unusual int `json:"unusual"`
	Usual   int `json:"usual"`
	Low     int `json:"low"`
}

type RespondLotteryObject struct {
	ParticipantsCount  int           `json:"participantsCount"`
	BankPerParticipant float64       `json:"bankPerParticipant"`
	Bank               float64       `json:"bank"`
	BankRemain         float64       `json:"bankRemain"`
	PotentialWinners   int           `json:"potentialWinners"`
	Lottery            LotteryObject `json:"lottery"`
}

type LotteryCreator struct {
	ParticipantsCount       int  `json:"participantsCount"  binding:"required"`
	Rate                    int  `json:"rate"  binding:"required"`
	QualityOverQuantityMode bool `json:"qualityOverQuantityMode"  binding:"required"`
}

func (r *RespondLotteryObject) generateLotteryObject(qualityOverQuantityMode bool) {
	if qualityOverQuantityMode {
		for r.BankRemain >= 1 {
			switch {
			case r.BankRemain >= float64(epicItem):
				r.Lottery.Epic++
				r.BankRemain -= float64(epicItem)
			case r.BankRemain >= float64(rareItem):
				r.Lottery.Rare++
				r.BankRemain -= float64(rareItem)
			case r.BankRemain >= float64(unusualItem):
				r.Lottery.Unusual++
				r.BankRemain -= float64(unusualItem)
			case r.BankRemain >= float64(usualItem):
				r.Lottery.Usual++
				r.BankRemain -= float64(usualItem)
			case r.BankRemain >= float64(lowQualityItem):
				r.Lottery.Low++
				r.BankRemain -= float64(lowQualityItem)
			default:
				break
			}
		}
	} else {
		potentialWinners := r.PotentialWinners
	potentialWinnersLoop:
		for r.PotentialWinners > 0 {
			bankForParticipantLeft := r.BankRemain / float64(potentialWinners)
			switch {
			case r.BankRemain > float64(epicItem) && bankForParticipantLeft >= 9:
				r.Lottery.Epic++
				r.BankRemain -= float64(epicItem)
				potentialWinners--
			case r.BankRemain > float64(rareItem) && bankForParticipantLeft >= 4.5:
				r.Lottery.Rare++
				r.BankRemain -= float64(rareItem)
				potentialWinners--
			case r.BankRemain > float64(unusualItem) && bankForParticipantLeft >= 2:
				r.Lottery.Unusual++
				r.BankRemain -= float64(unusualItem)
				potentialWinners--
			case r.BankRemain > float64(usualItem) && bankForParticipantLeft >= 1.2:
				r.Lottery.Usual++
				r.BankRemain -= float64(usualItem)
				potentialWinners--
			default:
				break potentialWinnersLoop
			}
		}
		if r.BankRemain >= 1 {
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
	respond.Lottery = LotteryObject{0, 0, 0, 0, 0}
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if lotteryCreator.Rate < 7 {
		c.JSON(http.StatusBadRequest, gin.H{"error": lotteryCreatorError("Lottery could be created for rate 7 and more")})
		return
	}
	if lotteryCreator.ParticipantsCount < 10 {
		c.JSON(http.StatusBadRequest, gin.H{"error": lotteryCreatorError("Lottery could be created for more than 10 participants.")})
		return
	}
	respond := lotteryCreator.generateLottery()
	c.JSON(http.StatusOK, respond)
}
