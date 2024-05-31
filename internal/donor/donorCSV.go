package donor

import (
	"encoding/csv"
	"fmt"
	"io"
)

type DonorCSV struct {
	Reader    *csv.Reader
	DonorChan chan<- *Donor
}

func NewDonorCSV(r io.Reader, donorChan chan<- *Donor) (*DonorCSV, error) {
	csvReader := csv.NewReader(r)

	// header
	header, err := csvReader.Read()
	if err != nil {
		return nil, err
	}

	fmt.Println(header)

	return &DonorCSV{
		csvReader, donorChan,
	}, nil
}
