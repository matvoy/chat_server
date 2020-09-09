package main

type webhookReqBody struct {
	Message struct {
		MessageID int64  `json:"message_id"`
		Text      string `json:"text"`
		From      struct {
			Username  string `json:"username"`
			ID        int64  `json:"id"`
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
		} `json:"from"`
		Chat struct {
			ID int64 `json:"id"`
		} `json:"chat"`
	} `json:"message"`
}
