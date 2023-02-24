package main

import (
	"fmt"
	"time"

	"github.com/Instituto-Atlantico/janus/src/controller/holder"
	"github.com/Instituto-Atlantico/janus/src/controller/issuer"
)

func main() {
	fmt.Println("CreateInvitation")
	invitation, _ := issuer.CreateInvitation("createdByCodeTestWFunc", true, false, false)

	fmt.Println("ReceiveInvitation")
	connection, _ := holder.ReceiveInvitation(invitation, true)

	fmt.Println("RegisterSchema")
	schema, _ := issuer.RegisterSchema("schema-00124", "0.1", []string{"name", "age"})

	fmt.Println("CreateCredentialDefinition")
	credDefinition, _ := issuer.CreateCredentialDefinition("default", false, 0, schema.ID)

	time.Sleep(2 * time.Second)

	fmt.Println("OfferCredentialV2")
	issuer.OfferCredentialV2(invitation.ConnectionID, credDefinition, "first comment bellow")

	fmt.Println("PresentationRequestRequest")
	presentationRequest, _ := issuer.PresentationRequestRequest(credDefinition, invitation)

	fmt.Println("SendPresentationByID")
	time.Sleep(3 * time.Second)
	holder.SendPresentationByID(connection)

	fmt.Println("VerifyPresentationByID")
	time.Sleep(3 * time.Second)
	issuer.VerifyPresentationByID(presentationRequest)

	fmt.Println("GetPresentationExchangeByID")
	proofValidation, _ := issuer.GetPresentationExchangeByID(presentationRequest)
	fmt.Println(proofValidation)
}
