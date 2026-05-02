package models

type Transaction struct {
	ID          int     `json:"id"`
	UserID      int     `json:"user_id"`
	Type        string  `json:"type"`
	Amount      float64 `json:"amount"`
	Category    string  `json:"category"`
	Description string  `json:"description"`
	Date        string  `json:"date"`
	CreatedAt   string  `json:"created_at"`
}

type UpdateTransactionRequest struct {
	Type        *string  `json:"type"`
	Amount      *float64 `json:"amount"`
	Category    *string  `json:"category"`
	Description *string  `json:"description"`
	Date        *string  `json:"date"`
}
