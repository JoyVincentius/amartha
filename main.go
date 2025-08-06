package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func main() {
	dsn := "root:root@tcp(127.0.0.1:3306)/amartha?parseTime=true"
	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("gagal connect DB: %v", err)
	}
	defer db.Close()

	if err = ensureTables(db); err != nil {
		log.Fatalf("gagal buat tabel: %v", err)
	}
	if err = seedExampleLoan(db); err != nil {
		log.Fatalf("gagal seed loan: %v", err)
	}

	http.HandleFunc("/make-payment", makePaymentHandler)
	http.HandleFunc("/outstanding", outstandingHandler)
	http.HandleFunc("/delinquent", delinquentHandler)

	log.Println("Server berjalan di http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func writeJSON(w http.ResponseWriter, status int, resp ApiResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(resp)
}

func makePaymentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, ApiResponse{Success: false, Message: "Method not allowed"})
		return
	}

	var req MakePaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, ApiResponse{Success: false, Message: "Invalid JSON"})
		return
	}

	be, err := NewBillingEngine(db, req.LoanID)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, ApiResponse{Success: false, Message: err.Error()})
		return
	}

	if err := be.MakePayment(req.Week, req.Amount); err != nil {
		writeJSON(w, http.StatusBadRequest, ApiResponse{Success: false, Message: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, ApiResponse{Success: true, Message: "Payment recorded"})
}

func outstandingHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, ApiResponse{Success: false, Message: "Method not allowed"})
		return
	}

	var req GetOutstandingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, ApiResponse{Success: false, Message: "Invalid JSON"})
		return
	}

	be, err := NewBillingEngine(db, req.LoanID)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, ApiResponse{Success: false, Message: err.Error()})
		return
	}

	out, err := be.GetOutstanding()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, ApiResponse{Success: false, Message: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, ApiResponse{Success: true, Data: map[string]float64{"outstanding": out}})
}

func delinquentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, ApiResponse{Success: false, Message: "Method not allowed"})
		return
	}

	var req IsDelinquentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, ApiResponse{Success: false, Message: "Invalid JSON"})
		return
	}

	be, err := NewBillingEngine(db, req.LoanID)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, ApiResponse{Success: false, Message: err.Error()})
		return
	}

	delinquent, err := be.IsDelinquent(req.CurrentWeek)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, ApiResponse{Success: false, Message: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, ApiResponse{Success: true, Data: map[string]bool{"delinquent": delinquent}})
}
