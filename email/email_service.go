package email

import (
	"crypto/rand"
	"fmt"
	"github.com/joho/godotenv"
	"gopkg.in/mail.v2"
	"log"
	"math/big"
	"os"
	"strconv"
)

func generateSecureToken() string {
	maxNumber := big.NewInt(100000)
	n, err := rand.Int(rand.Reader, maxNumber)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%05d", n) // Ensure it's always 5 digits
}

func SendVerificationEmail(email string) {
	secureToken := generateSecureToken()
	err := godotenv.Load("resources/.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	host := os.Getenv("EMAIL_HOST")
	port, _ := strconv.Atoi(os.Getenv("EMAIL_PORT"))
	username := os.Getenv("EMAIL_USERNAME")
	password := os.Getenv("EMAIL_PASSWORD")

	d := mail.NewDialer(host, port, username, password)

	m := mail.NewMessage()

	m.SetHeader("From", "email@gmail.com")
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Hello!")
	m.SetBody("text/plain", "Your secure token is: "+secureToken)

	// Send the email
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}

	fmt.Printf("Sending verification email to %s. Verify at /verify-email?token=%s\n", email, secureToken)
}
