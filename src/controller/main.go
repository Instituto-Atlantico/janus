package main

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/Instituto-Atlantico/janus/src/controller/dojot"
	"github.com/Instituto-Atlantico/janus/src/controller/holder"
	"github.com/Instituto-Atlantico/janus/src/controller/issuer"
	"github.com/ldej/go-acapy-client"
)

var (
	validAttributes = []acapy.CredentialPreviewAttributeV2{
		{
			MimeType: "text/plain",
			Name:     "name",
			Value:    "Vitor",
		},
		{
			MimeType: "text/plain",
			Name:     "age",
			Value:    "20",
		},
	}

	invalidAttributes = []acapy.CredentialPreviewAttributeV2{
		{
			MimeType: "text/plain",
			Name:     "name",
			Value:    "Vitor",
		},
		{
			MimeType: "text/plain",
			Name:     "age",
			Value:    "10",
		},
	}
)

type Try interface {
	acapy.Connection | acapy.Credential
}

func TryUtilNoError[R Try](fn func() (R, error)) (R, error) {
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
	//schemas and cred defs
	schema := "9RgfwfrRcTjESbVVGaSQa:2:schema-janus-0104:0.1"
	fmt.Println("Schema attributes: ", issuer.GetSchemaAttributes(schema))

	fmt.Println("\nCreating a credential definition for the schema")
	credDef, err := issuer.GetCredDef(schema)
	if err != nil && err.Error() == "empty" {
		credDef, _ = issuer.CreateCredentialDefinition("default", false, 0, schema)
	}
	fmt.Println("Cred Def ID: ", credDef)

	// Invitations
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

	goodCred, err := holder.GetCredential("age", "20")
	if err != nil && err.Error() == "empty" {
		issuer.OfferCredentialV2(issuerConnection.ConnectionID, credDef, "good credential", validAttributes)

		goodCred, err = TryUtilNoError(func() (acapy.Credential, error) { return holder.GetCredential("age", "20") })
		if err != nil {
			log.Fatal("timeout on holder.GetCredential('age', '20')")
		}
	}
	fmt.Println("Good cred: ", goodCred.Referent)

	badCred, err := holder.GetCredential("age", "10")
	if err != nil && err.Error() == "empty" {
		issuer.OfferCredentialV2(issuerConnection.ConnectionID, credDef, "bad credential", invalidAttributes)

		badCred, err = TryUtilNoError(func() (acapy.Credential, error) { return holder.GetCredential("age", "10") })
		if err != nil {
			log.Fatal("timeout on holder.GetCredential('age', '10')")
		}
	}
	fmt.Println("Bad cred: ", badCred.Referent)

	time.Sleep(2 * time.Second)
	fmt.Println("\nAsking for presentation (good)")
	presentationIssuer, _ := issuer.PresentationRequestRequest(credDef, issuerConnection)

	time.Sleep(1 * time.Second)

	holder.SendPresentationByID(presentationIssuer.ThreadID, goodCred)

	time.Sleep(1 * time.Second)

	_, err = issuer.VerifyPresentationByID(presentationIssuer)
	if err != nil {
		log.Fatal("verification failed: ", err)
	}

	LogMessageIfPresentationIsValid(presentationIssuer.ThreadID, "hello world")

	// fmt.Println("\nAsking for presentation (bad)")
	// presentationIssuer, _ := issuer.PresentationRequestRequest(credDef, issuerConnection)

	// time.Sleep(1 * time.Second)

	// holder.SendPresentationByID(presentationIssuer.ThreadID, badCred)

	// time.Sleep(1 * time.Second)

	// _, err = issuer.VerifyPresentationByID(presentationIssuer)
	// if err != nil {
	// 	log.Fatal("verification failed: ", err)
	// }

	// LogMessageIfPresentationIsValid(presentationIssuer.ThreadID, "hello world")
}

func LogMessageIfPresentationIsValid(threadID, message string) {
	presentation, err := issuer.GetPresentationExchangeByThreadId(threadID)
	if err != nil {
		log.Fatal("Presentation validation failed, ", err)
	}

	if presentation.Verified == "true" {
		fmt.Println("Publishing message to Dojot...")
		dojot.PublishMessage("localhost", "admin:9bda75", "admin", "http://192.168.0.4", "admin:9bda75/attrs")
	} else {
		log.Fatal("Presentation validation failed.")
	}
}
