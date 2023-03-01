package main

import (
	"fmt"
	"log"
	"time"

	"github.com/Instituto-Atlantico/janus/src/controller/holder"
	"github.com/Instituto-Atlantico/janus/src/controller/issuer"
	"github.com/ldej/go-acapy-client"
)

func LogHelloForVerifiedDevices(presentation any) {

}

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

func main() {
	// Invitations

	fmt.Println("\nGetting connections and making invitations")

	issuerConnection, err := issuer.GetConnection()
	if err != nil && err.Error() == "empty" {
		invitation, _ := issuer.CreateInvitation("createdByCode", true, false, false)
		holder.ReceiveInvitation(invitation, true)
		time.Sleep(5 * time.Second)

		issuerConnection, _ = issuer.GetConnection()
	}

	holderConnection, _ := holder.GetConnection()
	fmt.Println("issuer connection: ", issuerConnection.ConnectionID)
	fmt.Println("holder connection: ", holderConnection.ConnectionID)

	fmt.Println("\nGetting and Registering schemas and cred defs")

	schema, err := issuer.GetSchema("schema-janus-0104")
	if err != nil && err.Error() == "empty" {
		resp, _ := issuer.RegisterSchema("schema-janus-0104", "0.1", []string{"name", "age"})

		schema = resp.ID
	}

	credDef, err := issuer.GetCredDef(schema)
	if err != nil && err.Error() == "empty" {
		credDef, _ = issuer.CreateCredentialDefinition("default", false, 0, schema)
	}

	fmt.Println("Schema ID: ", schema)
	fmt.Println("Cred Def ID: ", credDef)

	time.Sleep(3 * time.Second)

	fmt.Println("\nIssuing Credentials")

	goodCred, err := holder.GetCredential("age", "20")
	if err != nil && err.Error() == "empty" {
		issuer.OfferCredentialV2(issuerConnection.ConnectionID, credDef, "good credential", validAttributes)

		time.Sleep(10 * time.Second)
		goodCred, _ = holder.GetCredential("age", "20")
	}

	badCred, err := holder.GetCredential("age", "10")
	if err != nil && err.Error() == "empty" {
		issuer.OfferCredentialV2(issuerConnection.ConnectionID, credDef, "bad credential", invalidAttributes)

		time.Sleep(10 * time.Second)
		badCred, _ = holder.GetCredential("age", "10")
	}

	fmt.Println("Good cred: ", goodCred.Referent)
	fmt.Println("Bad cred: ", badCred.Referent)

	fmt.Println("\nAsking for presentation (good)")

	presentationIssuer, _ := issuer.PresentationRequestRequest(credDef, issuerConnection)

	time.Sleep(1 * time.Second)

	holder.SendPresentationByID(presentationIssuer.ThreadID, badCred)

	time.Sleep(1 * time.Second)

	_, err = issuer.VerifyPresentationByID(presentationIssuer)
	if err != nil {
		log.Fatal("verification failed: ", err)
	}

	LogMessageIfPresentationIsValid(presentationIssuer.ThreadID, "hello world")
}

func LogMessageIfPresentationIsValid(threadID, message string) {
	presentation, err := issuer.GetPresentationExchangeByThreadId(threadID)
	if err != nil {
		log.Fatal("Presentation validation failed, ", err)
	}

	if presentation.Verified == "true" {
		fmt.Println("Message from holder: ", message) //this can be changed for other behaviors
	} else {
		log.Fatal("Presentation validation failed.")
	}
}
