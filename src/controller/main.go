package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/ldej/go-acapy-client"
)

func main() {
	//issuer creates an invitation
	issuer := acapy.NewClient("http://192.168.229.8:8002/")
	holder := acapy.NewClient("http://192.168.229.8:9002/")

	fmt.Println("Create Invitation")
	invitation, err := issuer.CreateInvitation("createdByCodeTest", true, false, false)
	if err != nil {
		log.Fatal(err)
	}

	//holder receive the invitation
	parsedInvitation, _ := json.Marshal(invitation.Invitation)

	var invitationTyped acapy.Invitation

	json.Unmarshal(parsedInvitation, &invitationTyped)

	fmt.Println("Receive Invitation")
	connection, err := holder.ReceiveInvitation(invitationTyped, true)
	if err != nil {
		log.Fatal(err)
	}
	_ = connection

	//issuer register a schema
	schemaName := fmt.Sprintf("schema-elton-%v", time.Now().Unix())

	fmt.Println("Register Schema")
	schema, err := issuer.RegisterSchema(schemaName, "0.1", []string{"name", "age"})
	if err != nil {
		log.Default().Println(err)
	}

	fmt.Println(schema.ID)

	time.Sleep(1 * time.Second)

	// issuer creates a credential definition
	fmt.Println("Create Cred Definition")
	credentialDefinition, err := issuer.CreateCredentialDefinition("default", false, 0, schema.ID)
	if err != nil {
		log.Default().Println(err)
	}

	fmt.Println(credentialDefinition)

	// get public DID
	fmt.Println("Create Cred Definition")
	issuerDID, err := issuer.GetPublicDID()
	if err != nil {
		log.Fatal(err)
	}

	// issuer issues a credential
	attributes := []acapy.CredentialPreviewAttributeV2{
		{
			MimeType: "text/plain",
			Name:     "name",
			Value:    "Elton",
		},
		{
			MimeType: "text/plain",
			Name:     "age",
			Value:    "30",
		},
	}

	credentialPreview := acapy.NewCredentialPreviewV2(attributes)

	fmt.Println("Proposal Credential")
	holder.ProposeCredentialV2(connection.ConnectionID, credentialPreview, "first comment", credentialDefinition, issuerDID.DID, schema.ID)

	// holder sends a presentation

	// issuer verifies
}
