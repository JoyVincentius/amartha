package main

import "database/sql"

func ensureTables(db *sql.DB) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS Loans (
            LoanID INT PRIMARY KEY,
            Principal DECIMAL(15,2) NOT NULL,
            InterestRate DECIMAL(5,4) NOT NULL,
            DurationWeeks INT NOT NULL,
            CreatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            UpdatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
        )`,
		`CREATE TABLE IF NOT EXISTS Payments (
            PaymentID INT AUTO_INCREMENT PRIMARY KEY,
            LoanID INT NOT NULL,
            Week INT NOT NULL,
            Amount DECIMAL(15,2) NOT NULL,
            Paid BOOL NOT NULL DEFAULT 1,
            PaidAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            UNIQUE KEY uq_loan_week (LoanID, Week),
            FOREIGN KEY (LoanID) REFERENCES Loans(LoanID) ON DELETE CASCADE
        )`,
	}
	for _, q := range queries {
		if _, err := db.Exec(q); err != nil {
			return err
		}
	}
	return nil
}

func seedExampleLoan(db *sql.DB) error {
	_, err := db.Exec(`
        INSERT IGNORE INTO Loans (LoanID, Principal, InterestRate, DurationWeeks) 
        VALUES (100, 5500000, 0.10, 500)`)
	return err
}
