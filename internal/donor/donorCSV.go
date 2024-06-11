package donor

import (
	"encoding/csv"
	"errors"
	"io"
	"log"
	"strconv"
	"time"
)

type DonorCSV struct {
	Reader *csv.Reader
}

func NewDonorCSV(r io.Reader) (*DonorCSV, error) {
	csvReader := csv.NewReader(r)

	// header
	header, err := csvReader.Read()
	if err != nil {
		return nil, err
	}

	// minimum validation of csv
	if len(header) < 6 {
		return nil, errors.New("CSV has missing fields")
	}

	return &DonorCSV{
		csvReader,
	}, nil
}

func (d *DonorCSV) Read() []*Donor {
	var donors []*Donor

	for {
		record, err := d.Reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}

			log.Println("CSV Read error:", err)
		}

		donors = append(donors, parseDonorRecord(record))
	}

	return donors
}

func parseDonorRecord(r []string) *Donor {
	// fit to int64
	amount, err := strconv.ParseInt(r[1], 10, 64)
	if err != nil {
		log.Println("Malformed amount in record:", err)
	}

	month, err := strconv.Atoi(r[4])
	if err != nil {
		log.Println("Malformed month in record:", err)
	}

	year, err := strconv.Atoi(r[4])
	if err != nil {
		log.Println("Malformed year in record:", err)
	}

	return &Donor{
		Name:     r[0],
		Amount:   amount,
		CCNumber: r[2],
		CVV:      r[3],
		ExpMonth: time.Month(month),
		ExpYear:  year,
	}
}
