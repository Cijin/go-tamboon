package main

import (
	"fmt"
	"log"
	"os"

	"go-tamboon/cipher"
	"go-tamboon/internal/donor"
	"go-tamboon/internal/transaction"

	"github.com/joho/godotenv"
	"github.com/omise/omise-go"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("warning: assuming default configuration. .env unreadable: %v", err)
	}

	filePath := os.Args[1]
	publicKey := os.Getenv("OMISE_PUBLIC_KEY")
	privateKey := os.Getenv("OMISE_PRIVATE_KEY")

	if len(publicKey) == 0 || len(privateKey) == 0 {
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

	donorCSV, err := donor.NewDonorCSV(cipherReader, donorChan)
	if err != nil {
		log.Println("Csv might be corrupted")
		os.Exit(1)
	}

	donors := donorCSV.Read()

	client, err := omise.NewClient(publicKey, privateKey)
	if err != nil {
		log.Println("There was an error initializing omise client:", err)
		os.Exit(1)
	}

	fmt.Println("performing donations...")

	summary := transaction.ProcessDonations(client, donors)
	fmt.Println(summary)
}
