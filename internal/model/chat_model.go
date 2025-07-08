package model

import "encoding/json"

type CreateChatRequest struct {
	PrID        int             `json:"pr_id"`
	ChatHistory json.RawMessage `json:"chat_history"`
}

type ChatReponse struct {
	ID          int             `json:"id"`
	PrID        int             `json:"pr_id"`
	ChatHistory json.RawMessage `json:"chat_history"`
}
