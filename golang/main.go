package main

import (
	"github.com/google/uuid"
	"log"
	"net/http"
)

var (
	AppId  = "REPLACE_WITH_YOUR_APP_ID"
	AppKey = "REPLACE_WITH_YOUR_APP_KEY"
)

var client = NewHumanCodeClient(http.Client{}, &ClientConfig{
	BaseUrl: "https://humancodeai.com",
	AppId:   AppId,
	AppKey:  AppKey,
})

func GetSessionIdExample() {
	if result, err := client.GetSessionId("123123"); err != nil {
		panic(err)
	} else {
		log.Println("SessionId:", result.SessionId)
	}
}

func GenRegistrationUrlExample() {
	if result, err := client.GenRegistrationUrl("", ""); err != nil {
		panic(err)
	} else {
		log.Println("RegistrationUrl:", result)
	}
}

func GenGenVerificationUrlExample() {
	if result, err := client.GenVerificationUrl("", "", ""); err != nil {
		panic(err)
	} else {
		log.Println("GenVerificationUrl:", result)
	}
}

func VerifyVCodeExample() {
	if result, err := client.Verify("", "", uuid.NewString()); err != nil {
		panic(err)
	} else {
		log.Println("HumanId:", result.HumanId)
	}
}

func main() {
	GetSessionIdExample()
	//GenRegistrationUrlExample()
	//GenGenVerificationUrlExample()
	//VerifyVCodeExample()
}
