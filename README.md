This repository contains a simple **Billing Engine** written in Go that manages loan repayment schedules, stores data in MySQL, and exposes a RESTful HTTP API.

## Features

- Calculate weekly payment amount based on flat interest.
- Record weekly payments in a MySQL `Payments` table.
- Detect delinquent loans (3 or more consecutive missed weeks).
- HTTP endpoints for creating loans, making payments, checking outstanding balance and delinquency status.
- Git‑initialized project with a ready‑to‑run `main.go`.

## Prerequisites

- Go 1.22+ installed
- MySQL server (local or remote) reachable from the application
- `git` (optional, the project is already initialized)

## Installation

1. **Clone the repository** (or copy the project folder) and navigate into it:
   ```bash
   git clone <your-repo-url>
   cd project

2. **Run the project**
   ```bash
   go run .

## Postman Collection

Please import collection to your postman.
Notes:
- /make-payment : save payment in week (?)
- /outstanding : outstanding amount
- /delinquent : check if delinquent