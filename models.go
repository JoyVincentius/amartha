package main

type MakePaymentRequest struct {
	LoanID int     `json:"loan_id"` 
	Week   int     `json:"week"`    
	Amount float64 `json:"amount"`  
}

type GetOutstandingRequest struct {
	LoanID int `json:"loan_id"`
}

type IsDelinquentRequest struct {
	LoanID      int `json:"loan_id"`
	CurrentWeek int `json:"current_week"` 
}

type ApiResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}
