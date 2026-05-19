package models

// Transaction ตรงกับ schema ของจันทร์พงศ์ใน migrations/schema.sql
type Transaction struct {
	ID         int     `json:"id"`
	UserID     int     `json:"user_id"`
	CategoryID int     `json:"category_id"`
	Type       string  `json:"type"`
	Amount     float64 `json:"amount"`
	Note       string  `json:"note"`
	Date       string  `json:"date"`
	CreatedAt  string  `json:"created_at"`
	UpdatedAt  string  `json:"updated_at"`
}

// UpdateTransactionRequest ใช้ pointer เพื่อแยกแยะ "ไม่ส่งมา" vs "ส่ง null"
type UpdateTransactionRequest struct {
	CategoryID *int     `json:"category_id"`
	Type       *string  `json:"type"`
	Amount     *float64 `json:"amount"`
	Note       *string  `json:"note"`
	Date       *string  `json:"date"`
}
