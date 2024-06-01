package donor

import "time"

type Donor struct {
	Name     string
	Amount   int64
	CCNumber string
	CVV      string
	ExpMonth time.Month
	ExpYear  int
}
