package main

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/Instituto-Atlantico/janus/src/controller/holder"
	"github.com/Instituto-Atlantico/janus/src/controller/issuer"
	"github.com/ldej/go-acapy-client"
)

var (
	validAttributes = []acapy.CredentialPreviewAttributeV2{
		{
			MimeType: "text/plain",
			Name:     "temperature",
			Value:    "false",
		},
		{
			MimeType: "text/plain",
			Name:     "humidity",
			Value:    "true",
		},
	}
)

func TryUtilNoError[R any](fn func() (R, error)) (R, error) {
	cResponse := make(chan R)
	cTimeout := make(chan string)

	go func() {
		time.Sleep(20 * time.Second)
		cTimeout <- ""
	}()

	go func() {
		for {
			response, err := fn()
			if err == nil {
				cResponse <- response
				return
			}
			time.Sleep(time.Second)
		}
	}()

	select {
	case data := <-cResponse:
		return data, nil
	case <-cTimeout:
		return *new(R), errors.New("Timeout")
	}
}

func main() {
	//schema
	schema := "EZpfyRHcXuohyTvbgsrg7S:2:janus-sensors:1.0"
	fmt.Println("Schema attributes: ", issuer.GetSchemaAttributes(schema))

	//cred definition
	fmt.Println("\nCreating a credential definition for the schema")
	credDef, err := issuer.GetCredDef(schema)
	if err != nil && err.Error() == "empty" {
		credDef, _ = issuer.CreateCredentialDefinition("default", false, 0, schema)
	}
	fmt.Println("Cred Def ID: ", credDef)

	//invitations
	fmt.Println("\nGetting connections and making invitations")

	issuerConnection, err := issuer.GetConnection()
	if err != nil && err.Error() == "empty" {
		invitation, _ := issuer.CreateInvitation("createdByCode", true, false, false)
		holder.ReceiveInvitation(invitation, true)

		issuerConnection, err = TryUtilNoError(issuer.GetConnection)
		if err != nil {
			log.Fatal("timeout on issuer.GetConnection")
		}
	}

	holderConnection, _ := holder.GetConnection()
	fmt.Println("issuer connection: ", issuerConnection.ConnectionID)
	fmt.Println("holder connection: ", holderConnection.ConnectionID)

	time.Sleep(2 * time.Second)

	// issuing credentials
	fmt.Println("\nIssuing Credentials")
	cred, err := holder.GetCredential("schema_id", schema)
	if err != nil {
		issuer.OfferCredentialV2(issuerConnection.ConnectionID, credDef, "good credential", validAttributes)
		cred, err = TryUtilNoError(func() (acapy.Credential, error) { return holder.GetCredential("schema_id", schema) })
		if err != nil {
			log.Fatal("timeout on holder.GetCredential")
		}
	}

	fmt.Println("cred: ", cred.Referent)

	//Presentation

	//ask for presentation
	fmt.Println("\nAsking for presentation")
	presentationIssuerHumidity, _ := issuer.PresentationRequestRequest(credDef, issuerConnection, "humidity")
	presentationIssuerTemp, _ := issuer.PresentationRequestRequest(credDef, issuerConnection, "temperature")
	time.Sleep(2 * time.Second)
	fmt.Println("presentation request IDs:", presentationIssuerHumidity.ThreadID, presentationIssuerTemp.ThreadID)

	//send presentation
	holder.SendPresentationByID(presentationIssuerHumidity, cred)
	holder.SendPresentationByID(presentationIssuerTemp, cred)

	_, err = TryUtilNoError(func() ([]acapy.PresentationExchangeRecord, error) {
		return issuer.IsPresentationDone(presentationIssuerHumidity.ThreadID)
	})
	if err != nil {
		log.Fatal("timeout issuer.IsPresentationDone for humidity")
	}

	_, err = TryUtilNoError(func() ([]acapy.PresentationExchangeRecord, error) {
		return issuer.IsPresentationDone(presentationIssuerTemp.ThreadID)
	})
	if err != nil {
		log.Fatal("timeout issuer.IsPresentationDone for temperature")
	}

	//verify
	fmt.Println("Verifing presentation")
	_, err = issuer.VerifyPresentationByID(presentationIssuerHumidity)
	if err != nil {
		log.Fatal("verification failed: ", err)
	}

	_, err = issuer.VerifyPresentationByID(presentationIssuerTemp)
	if err != nil {
		log.Fatal("verification failed: ", err)
	}

	LogMessageIfPresentationIsValid(presentationIssuerHumidity.ThreadID, "YEEEEE VALID PRESENTATION for humidity🥳")
	LogMessageIfPresentationIsValid(presentationIssuerTemp.ThreadID, "YEEEEE VALID PRESENTATION for temperature🥳")
}

func LogMessageIfPresentationIsValid(threadID, message string) {
	presentation, err := issuer.GetPresentationExchangeByThreadId(threadID)
	if err != nil {
		log.Fatal("Presentation validation failed, ", err)
	}

	if presentation.Verified == "true" {
		fmt.Println("Message from holder: ", message) //this can be changed for other behaviors
	} else {
		log.Fatal("Presentation validation failed, presentation is not valid")
	}
}
