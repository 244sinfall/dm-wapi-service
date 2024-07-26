package participants

type participantsRequest struct {
	RawText string `json:"rawText"  binding:"required"`
}

type participantsResponse struct {
	CleanedText string `json:"cleanedText"`
	Count       int    `json:"cleanedCount"`
	EditedLines string `json:"editedLines"`
}
