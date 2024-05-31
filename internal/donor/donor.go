package donor

import "time"

type Donor struct {
	Name     string
	Amount   int
	CCNumber string
	CVV      string
	ExpMonth time.Month
	ExpYear  int
}
