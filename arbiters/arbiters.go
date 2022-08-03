package arbiters

import (
	"bufio"
	"darkmoonWebApi/events"
	"fmt"
	"github.com/gin-gonic/gin"
	"math"
	"net/http"
	"strings"
)

type ArbiterWorkMode string

const (
	giveXP   ArbiterWorkMode = "givexp"
	takeXP   ArbiterWorkMode = "takexp"
	giveGold ArbiterWorkMode = "givegold"

	giveXPCommand   string = ".exp game"
	takeXPCommand   string = ".exp oth"
	giveGoldCommand string = ".send mo"

	xpToRateModifier   float64 = 1000
	goldToRateModifier float64 = 6000

	writerModifier           float64 = 1.2
	masterModifier           float64 = 0.3
	masterAndWriterModifier  float64 = 0.5
	crafterModifier          float64 = 0.5
	crafterAndWriterModifier float64 = 0.7
)

type ArbiterWorkRequest struct {
	ParticipantsCleanedText string          `json:"participantsCleanedText" binding:"required"`
	Mode                    ArbiterWorkMode `json:"mode" binding:"required"` // "givexp", "takexp", "givegold"
	Rate                    int             `json:"rate" binding:"required"`
	EventLink               string          `json:"eventLink" binding:"required"`
}

type ArbiterWorkResponse struct {
	Commands             string `json:"commands"`
	ParticipantsModified string `json:"participantsModified"`
}

func (r ArbiterWorkRequest) generateResponse() ArbiterWorkResponse {
	response := ArbiterWorkResponse{"", ""}
	scanner := bufio.NewScanner(strings.NewReader(r.ParticipantsCleanedText))
	participantsSlice := make([]string, 0)
	participantsAmount := 0
	for scanner.Scan() {
		line := scanner.Text()
		isLegitSuffix, suffix := events.CheckForLegitSuffixes(line)
		if strings.Count(line, " ") == 0 ||
			(strings.Count(line, " ") == 1 && isLegitSuffix) {
			participantsSlice = append(participantsSlice, line)
			if suffix != " M" && suffix != " WM" && suffix != "MW" {
				participantsAmount += 1
			}
		} else {
			response.ParticipantsModified += line + " - не прошел проверку на корректность\n"
		}
	}
	for _, participant := range participantsSlice {
		var defaultValueToManipulate float64 = float64(r.Rate)
		participantName := participant
		var valueToManipulate float64 = 0
		isLegitSuffix, suffix := events.CheckForLegitSuffixes(participant)
		if isLegitSuffix {
			participantName, _, _ = strings.Cut(participant, " ")
			if suffix == " W" {
				if participantsAmount >= 5 {
					valueToManipulate = defaultValueToManipulate * writerModifier
				} else {
					valueToManipulate = defaultValueToManipulate
					response.ParticipantsModified += participant + " - нет оснований для бонуса писателя\n"
				}
			} else if suffix == " WD" || suffix == " DW" {
				if participantsAmount >= 5 {
					valueToManipulate = defaultValueToManipulate * crafterAndWriterModifier
				} else {
					valueToManipulate = defaultValueToManipulate * crafterModifier
					response.ParticipantsModified += participant + " - нет оснований для бонуса писателя\n"
				}
			} else if suffix == " D" {
				valueToManipulate = defaultValueToManipulate * crafterModifier
			} else if suffix == " M" {
				if participantsAmount >= 5 {
					valueToManipulate = defaultValueToManipulate * masterModifier
				}
			} else if suffix == " WM" || suffix == " MW" {
				if participantsAmount >= 5 {
					valueToManipulate = defaultValueToManipulate * masterAndWriterModifier
				}
			}
		} else {
			valueToManipulate = defaultValueToManipulate
		}
		if valueToManipulate == 0 {
			response.ParticipantsModified += participant + " - нет оснований для выдачи награды\n"
		} else {
			switch r.Mode {
			case giveXP:
				response.Commands += fmt.Sprintf("%v %v %v %v\n", giveXPCommand, participantName, math.Round(valueToManipulate*xpToRateModifier), r.EventLink)
			case takeXP:
				response.Commands += fmt.Sprintf("%v %v -%v %v\n", takeXPCommand, participantName, math.Round(valueToManipulate*xpToRateModifier), r.EventLink)
			case giveGold:
				response.Commands += fmt.Sprintf("%v %v \"%v\" \"\" %v\n", giveGoldCommand, participantName, r.EventLink, math.Round(valueToManipulate*goldToRateModifier))
			}
		}
	}
	return response
}

func ArbiterWork(c *gin.Context) {
	var request ArbiterWorkRequest
	err := c.BindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	switch request.Mode {
	case giveXP, giveGold, takeXP:
		response := request.generateResponse()
		c.JSON(http.StatusOK, response)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unknown arbiter work mode"})
		return
	}
}
