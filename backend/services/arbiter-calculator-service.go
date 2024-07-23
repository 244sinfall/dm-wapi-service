package services

import (
	"bufio"
	"fmt"
	"math"
	"strings"
)

type ArbiterCalculationMode string

const (
	giveXP   ArbiterCalculationMode = "givexp"
	takeXP   ArbiterCalculationMode = "takexp"
	giveGold ArbiterCalculationMode = "givegold"

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

type ArbiterCalculationRequest struct {
	ParticipantsCleanedText string          `json:"participantsCleanedText" binding:"required"`
	Mode                    ArbiterCalculationMode `json:"mode" binding:"required"` // "givexp", "takexp", "givegold"
	Rate                    int             `json:"rate" binding:"required"`
	EventLink               string          `json:"eventLink" binding:"required"`
}

type ArbiterCalculationResponse struct {
	Commands             string `json:"commands"`
	ParticipantsModified string `json:"participantsModified"`
}

func (r ArbiterCalculationRequest) GenerateResponse() ArbiterCalculationResponse {
	response := ArbiterCalculationResponse{"", ""}
	scanner := bufio.NewScanner(strings.NewReader(r.ParticipantsCleanedText))
	participantsSlice := make([]string, 0)
	participantsAmount := 0
	for scanner.Scan() {
		line := scanner.Text()
		if len(strings.TrimSpace(line)) > 1 {
			isLegitSuffix, suffix := CheckForLegitSuffixes(line)
			if strings.Count(line, " ") == 0 ||
				(strings.Count(line, " ") == 1 && isLegitSuffix) {
				participantsSlice = append(participantsSlice, line)
				if suffix != " M" && suffix != " WM" && suffix != " MW" {
					participantsAmount += 1
				}
			} else {
				response.ParticipantsModified += line + " - не прошел проверку на корректность\n"
			}
		}
	}
	for _, participant := range participantsSlice {
		var defaultValueToManipulate float64 = float64(r.Rate)
		participantName := participant
		var valueToManipulate float64 = 0
		isLegitSuffix, suffix := CheckForLegitSuffixes(participant)
		if isLegitSuffix {
			participantName, _, _ = strings.Cut(participant, " ")
			if suffix == " W" {
				//if participantsAmount >= 5 {
				valueToManipulate = defaultValueToManipulate * writerModifier
				//} else {
				//	valueToManipulate = defaultValueToManipulate
				//	response.ParticipantsModified += participant + " - нет оснований для бонуса писателя\n"
				//}
			} else if suffix == " WD" || suffix == " DW" {
				//if participantsAmount >= 5 {
				valueToManipulate = defaultValueToManipulate * crafterAndWriterModifier
				//} else {
				//	valueToManipulate = defaultValueToManipulate * crafterModifier
				//	response.ParticipantsModified += participant + " - нет оснований для бонуса писателя\n"
				//}
			} else if suffix == " D" {
				valueToManipulate = defaultValueToManipulate * crafterModifier
			} else if suffix == " M" {
				//if participantsAmount >= 5 {
				valueToManipulate = defaultValueToManipulate * masterModifier
				//}
			} else if suffix == " WM" || suffix == " MW" {
				//if participantsAmount >= 5 {
				valueToManipulate = defaultValueToManipulate * masterAndWriterModifier
				//}
			}
		} else {
			valueToManipulate = defaultValueToManipulate
		}
		//if valueToManipulate == 0 {
		//	response.ParticipantsModified += participant + " - нет оснований для выдачи награды\n"
		//} else {
		switch r.Mode {
		case giveXP:
			response.Commands += fmt.Sprintf("%v %v %v %v\n", giveXPCommand, participantName, math.Round(valueToManipulate*xpToRateModifier), r.EventLink)
		case takeXP:
			response.Commands += fmt.Sprintf("%v %v -%v %v\n", takeXPCommand, participantName, math.Round(valueToManipulate*xpToRateModifier), r.EventLink)
		case giveGold:
			response.Commands += fmt.Sprintf("%v %v \"%v\" \"\" %v\n", giveGoldCommand, participantName, r.EventLink, math.Round(valueToManipulate*goldToRateModifier))
		}
		//}
	}
	return response
}