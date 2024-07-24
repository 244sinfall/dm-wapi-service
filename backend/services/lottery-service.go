package services

import (
	"math"
	"math/rand/v2"
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

type LotteryResult struct {
	Epic    int `json:"epic"`
	Rare    int `json:"rare"`
	Unusual int `json:"unusual"`
	Usual   int `json:"usual"`
	Low     int `json:"low"`
}

type LotteryResponse struct {
	ParticipantsCount  int           `json:"participantsCount"`
	BankPerParticipant float64       `json:"bankPerParticipant"`
	Bank               float64       `json:"bank"`
	BankRemain         float64       `json:"bankRemain"`
	PotentialWinners   int           `json:"potentialWinners"`
	Lottery            LotteryResult `json:"lottery"`
}

type LotteryOptions struct {
	ParticipantsCount       int  `json:"participantsCount"  binding:"required"`
	Rate                    int  `json:"rate"  binding:"required"`
	QualityOverQuantityMode bool `json:"qualityOverQuantityMode"`
}

func (r *LotteryResponse) generateLotteryObject(qualityOverQuantityMode bool) {
	if qualityOverQuantityMode {
	BankLoop:
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
				break BankLoop
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

func (l LotteryOptions) GenerateLottery() LotteryResponse {
	var respond LotteryResponse
	respond.ParticipantsCount = l.ParticipantsCount
	respond.BankPerParticipant = getBankPerParticipant(l.Rate)
	respond.Bank = float64(l.ParticipantsCount) * respond.BankPerParticipant
	respond.BankRemain = respond.Bank
	respond.PotentialWinners = int(math.Round(
		float64(l.ParticipantsCount) * (maxPotentialWinnersPart - (potentialWinnerStep * rand.Float64()))))
	respond.Lottery = LotteryResult{0, 0, 0, 0, 0}
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
