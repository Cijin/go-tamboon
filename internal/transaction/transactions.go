package transaction

import (
	"fmt"
	"log"
	"math/rand"
	"sort"
	"strings"
	"time"

	"go-tamboon/internal/donor"

	"github.com/omise/omise-go"
	"github.com/omise/omise-go/operations"
)

const currency = "thb"

func ProcessDonations(client *omise.Client, donors []*donor.Donor) string {
	var successfullDonors []*donor.Donor
	var successfullDonations int64
	var failedDonations int64

	chanSuccess := make(chan *donor.Donor)
	chanFail := make(chan *donor.Donor)

	defer func() {
		close(chanSuccess)
		close(chanFail)
	}()

	for _, donor := range donors {
		go processTransaction(client, donor, chanSuccess, chanFail)

		// basic rate limiter
		source := rand.NewSource(time.Now().UnixNano())
		rng := rand.New(source)
		time.Sleep(time.Duration(rng.Intn(2)) * time.Second)
	}

	for i := 0; i < len(donors); i++ {
		select {
		case donor := <-chanSuccess:
			successfullDonations += donor.Amount
			successfullDonors = append(successfullDonors, donor)

		case donor := <-chanFail:
			failedDonations += donor.Amount
		}
	}

	return summary(successfullDonations, failedDonations, successfullDonors)
}

func processTransaction(client *omise.Client, d *donor.Donor, chanSuccess, chanFail chan *donor.Donor) {
	card, createToken := &omise.Card{}, &operations.CreateToken{
		Name:            d.Name,
		Number:          d.CCNumber,
		ExpirationMonth: d.ExpMonth,
		ExpirationYear:  d.ExpYear,
		SecurityCode:    d.CVV,
	}

	// Requests fail at times	with error: looking for value with response body: <html>...
	if err := client.Do(card, createToken); err != nil {
		log.Println("Create omise token error:", err)

		chanFail <- d
		return
	}

	charge, createCharge := &omise.Charge{}, &operations.CreateCharge{
		Amount:   d.Amount,
		Currency: currency,
		Card:     card.ID,
	}

	if err := client.Do(charge, createCharge); err != nil {
		log.Println("Create omise charge err:", err)

		chanFail <- d
		return
	}

	if charge.Paid {
		chanSuccess <- d
		return
	}

	chanFail <- d
}

func summary(sucessAmount, failedAmount int64, sucessfullDonors []*donor.Donor) string {
	var summary strings.Builder

	// Use the Builder to efficiently construct the string
	summary.WriteString(fmt.Sprintf("total received:\t THB\t %d\n", sucessAmount+failedAmount))
	summary.WriteString(fmt.Sprintf("successfully donated:\t THB\t %d\n", sucessAmount))
	summary.WriteString(fmt.Sprintf("faulty donation:\t THB\t %d\n", failedAmount))
	summary.WriteString("\n")

	if len(sucessfullDonors) != 0 {
		summary.WriteString(fmt.Sprintf("average per person:\t THB\t %d\n", sucessAmount/int64(len(sucessfullDonors))))
	}

	sort.Slice(sucessfullDonors, func(i, j int) bool {
		return sucessfullDonors[i].Amount > sucessfullDonors[j].Amount
	})

	summary.WriteString("top donors:\n")
	for i, donor := range sucessfullDonors {
		if i > 2 {
			break
		}

		summary.WriteString(fmt.Sprintf("\t%s\n", donor.Name))
	}

	// Return the built string
	return summary.String()
}
