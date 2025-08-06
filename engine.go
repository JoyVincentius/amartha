package main

import (
	"database/sql"
	"fmt"
)

type BillingEngine struct {
	db            *sql.DB
	loanID        int
	principal     float64
	interestRate  float64
	durationWeeks int
	weeklyPay     float64
}

// Hitung cicilan mingguan (principal + bunga flat per minggu)
func (be *BillingEngine) computeWeeklyPayment() {
	interestPerWeek := (be.principal * be.interestRate) / float64(be.durationWeeks)
	be.weeklyPay = (be.principal / float64(be.durationWeeks)) + interestPerWeek
}

// Konstruktor
func NewBillingEngine(db *sql.DB, loanID int) (*BillingEngine, error) {
	row := db.QueryRow(`
        SELECT Principal, InterestRate, DurationWeeks 
        FROM Loans WHERE LoanID = ?`, loanID)

	var p, r float64
	var w int
	if err := row.Scan(&p, &r, &w); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("loan dengan ID %d tidak ditemukan", loanID)
		}
		return nil, fmt.Errorf("gagal membaca data loan: %w", err)
	}

	be := &BillingEngine{
		db:            db,
		loanID:        loanID,
		principal:     p,
		interestRate:  r,
		durationWeeks: w,
	}
	be.computeWeeklyPayment()
	return be, nil
}

// GetOutstanding menghitung sisa hutang (principal + bunga – total pembayaran)
func (be *BillingEngine) GetOutstanding() (float64, error) {
	var paidSum sql.NullFloat64
	err := be.db.QueryRow(`
        SELECT SUM(Amount) FROM Payments 
        WHERE LoanID = ? AND Paid = 1`, be.loanID).Scan(&paidSum)
	if err != nil && err != sql.ErrNoRows {
		return 0, fmt.Errorf("gagal menghitung total bayar: %w", err)
	}

	totalPaid := 0.0
	if paidSum.Valid {
		totalPaid = paidSum.Float64
	}
	outstanding := be.principal + (be.principal * be.interestRate) - totalPaid
	if outstanding < 0 {
		outstanding = 0
	}
	return outstanding, nil
}

// MakePayment menyimpan pembayaran satu minggu
func (be *BillingEngine) MakePayment(week int, amount float64) error {
	var alreadyPaid bool
	err := be.db.QueryRow(`
        SELECT Paid FROM Payments 
        WHERE LoanID = ? AND Week = ?`, be.loanID, week).Scan(&alreadyPaid)

	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("gagal cek status pembayaran: %w", err)
	}
	if alreadyPaid {
		return fmt.Errorf("minggu ke-%d sudah dibayar", week)
	}

	_, err = be.db.Exec(`
        INSERT INTO Payments (LoanID, Week, Amount, Paid) 
        VALUES (?, ?, ?, 1) 
        ON DUPLICATE KEY UPDATE Amount = ?, Paid = 1, PaidAt = CURRENT_TIMESTAMP`,
		be.loanID, week, amount, amount)
	return err
}

// IsDelinquent memeriksa apakah terdapat >2 minggu berturut‑tahun tidak bayar sampai currentWeek
func (be *BillingEngine) IsDelinquent(currentWeek int) (bool, error) {
	rows, err := be.db.Query(`
        SELECT Week FROM Payments 
        WHERE LoanID = ? AND Paid = 1 AND Week <= ? 
        ORDER BY Week ASC`, be.loanID, currentWeek)
	if err != nil {
		return false, fmt.Errorf("gagal query pembayaran: %w", err)
	}
	defer rows.Close()

	paidWeeks := make(map[int]bool)
	for rows.Next() {
		var w int
		if err := rows.Scan(&w); err != nil {
			return false, err
		}
		paidWeeks[w] = true
	}

	consecMiss := 0
	for w := 1; w <= currentWeek; w++ {
		if !paidWeeks[w] {
			consecMiss++
			if consecMiss > 2 {
				return true, nil
			}
		} else {
			consecMiss = 0
		}
	}
	return false, nil
}
