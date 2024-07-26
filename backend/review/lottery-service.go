package review

import (
	"math"
	"math/rand/v2"
)

func (r *lotteryResponse) generateLotteryObject(qualityOverQuantityMode bool) {
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

func (l lotteryOptions) generateLottery() lotteryResponse {
	var respond lotteryResponse
	respond.ParticipantsCount = l.ParticipantsCount
	respond.BankPerParticipant = getBankPerParticipant(l.Rate)
	respond.Bank = float64(l.ParticipantsCount) * respond.BankPerParticipant
	respond.BankRemain = respond.Bank
	respond.PotentialWinners = int(math.Round(
		float64(l.ParticipantsCount) * (maxPotentialWinnersPart - (potentialWinnerStep * rand.Float64()))))
	respond.Lottery = lotteryResult{0, 0, 0, 0, 0}
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
