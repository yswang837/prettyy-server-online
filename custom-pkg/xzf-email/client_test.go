package xzf_email

import (
	"fmt"
	"log"
	"testing"
)

func TestSendEmail(t *testing.T) {
	code, err := SendEmail("yswang837@gmail.com")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("code:", code)
}
