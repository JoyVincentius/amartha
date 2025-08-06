package main

type MakePaymentRequest struct {
	LoanID int     `json:"loan_id"` // ID pinjaman
	Week   int     `json:"week"`    // minggu ke‑berapa yang dibayar
	Amount float64 `json:"amount"`  // jumlah pembayaran
}

type GetOutstandingRequest struct {
	LoanID int `json:"loan_id"`
}

type IsDelinquentRequest struct {
	LoanID      int `json:"loan_id"`
	CurrentWeek int `json:"current_week"` // minggu yang sedang dicek
}

type ApiResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}
