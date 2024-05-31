package main

import (
	"log"
	"os"

	"go-tamboon/cipher"
	"go-tamboon/internal/donor"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("warning: assuming default configuration. .env unreadable: %v", err)
	}

	filePath := os.Args[1]
	publicKey := os.Getenv("OMISE_PUBLIC_KEY")
	secretKey := os.Getenv("OMISE_SECRET_KEY")

	if len(publicKey) == 0 || len(secretKey) == 0 {
		log.Println("Missing omise public and secret enviornment keys")
		os.Exit(1)
	}

	file, err := os.Open(filePath)
	if err != nil {
		log.Println("Unable to open csv file:", err)
		os.Exit(1)
	}

	defer file.Close()

	// err always nil
	cipherReader, _ := cipher.NewRot128Reader(file)

	donorChan := make(chan *donor.Donor)
	_, err = donor.NewDonorCSV(cipherReader, donorChan)
	if err != nil {
		log.Println("Csv might be corrupted")
		os.Exit(1)
	}
}
