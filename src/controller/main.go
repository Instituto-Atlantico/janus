package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Instituto-Atlantico/go-acapy-client"
)

func main() {
	//issuer creates an invitation
	issuer := acapy.NewClient("http://192.168.182.51:8002/")
	holder := acapy.NewClient("http://192.168.182.51:9002/")

	fmt.Println("Create Invitation")
	invitation, err := issuer.CreateInvitation("createdByCodeTest", true, false, false)
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println(invitation.ConnectionID)

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

	// fmt.Println(connection)

	//issuer register a schema
	schemaName := fmt.Sprintf("schema-elton-%v", time.Now().Unix())

	fmt.Println("Register Schema")
	schema, err := issuer.RegisterSchema(schemaName, "0.1", []string{"name", "age"})
	if err != nil {
		log.Default().Println(err)
	}

	// fmt.Println(schema.ID)

	time.Sleep(1 * time.Second)

	// issuer creates a credential definition
	fmt.Println("Create Cred Definition")
	credentialDefinition, err := issuer.CreateCredentialDefinition("default", false, 0, schema.ID)
	if err != nil {
		log.Default().Println(err)
	}

	// fmt.Println(credentialDefinition)

	// get public DID
	fmt.Println("Issuer DID")
	issuerDID, err := issuer.GetPublicDID()
	if err != nil {
		log.Fatal(err)
	}

	_ = issuerDID
	// fmt.Println(issuerDID.DID)

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

	fmt.Println("Offer Credential")
	offerCredential, err := issuer.OfferCredentialV2(invitation.ConnectionID, credentialPreview, credentialDefinition, "first comment bellow")
	if err != nil {
		log.Fatal(err)
	}

	_ = offerCredential

	// fmt.Println(offerCredential)

	// holder sends a presentation
	requestedPredicates := map[string]acapy.RequestedPredicate{
		"age": acapy.RequestedPredicate{
			Restrictions: []acapy.Restrictions{{ // Required in case of Names
				CredentialDefinitionID: credentialDefinition,
			}},
			Name:   "age", // XOR with Names
			PType:  acapy.PredicateGT,
			PValue: 18,
			NonRevoked: acapy.NonRevoked{
				From: time.Now().Add(-time.Hour * 24 * 7).Unix(),
				To:   time.Now().Unix(),
			},
		},
	}

	requestedAttributes := map[string]acapy.RequestedAttribute{
		"name": acapy.RequestedAttribute{
			Restrictions: []acapy.Restrictions{{ // Required in case of Names
				CredentialDefinitionID: credentialDefinition,
			}},
			Name: "name", // XOR with Names
			NonRevoked: acapy.NonRevoked{
				From: time.Now().Add(-time.Hour * 24 * 7).Unix(),
				To:   time.Now().Unix(),
			},
		},
	}

	presentationRequestRequest := acapy.PresentationRequestRequest{
		Trace:        false,
		Comment:      "Presentation request test",
		ConnectionID: invitation.ConnectionID,
		ProofRequest: acapy.NewProofRequest(
			"Elton",
			"1234567890",
			requestedPredicates,
			requestedAttributes,
			"1.0",
			&acapy.NonRevoked{
				From: time.Now().Add(-time.Hour * 24 * 7).Unix(),
				To:   time.Now().Unix(),
			},
		),
	}

	fmt.Println("Send Presentation Request")
	sendPresentation, err := issuer.SendPresentationRequest(presentationRequestRequest)
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println(sendPresentation)

	// Query presentation ID
	fmt.Println("Query presentation ID")
	param := acapy.QueryPresentationExchangeParams{
		ConnectionID: connection.ConnectionID,
	}

	time.Sleep(1 * time.Second)

	presentations, err := holder.QueryPresentationExchange(param)
	if err != nil {
		log.Fatal(err)
	}

	//Query Credential ID

	fmt.Println("Query credentials")
	credentials, err := holder.GetCredentials(1, 0, "")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(credentials[0])

	// holder send presentation
	attributesProof := map[string]acapy.PresentationProofAttribute{
		"name": {
			CredentialID: credentials[0].Referent,
			Revealed:     true,
		},
	}

	predicatesProof := map[string]acapy.PresentationProofPredicate{
		"age": {
			CredentialID: credentials[0].Referent,
		},
	}

	proof := acapy.PresentationProof{
		SelfAttestedAttributes: map[string]string{},
		RequestedAttributes:    attributesProof,
		RequestedPredicates:    predicatesProof,
	}

	fmt.Println("sending presentation")
	_, err = holder.SendPresentationByID(presentations[0].PresentationExchangeID, proof)
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(1 * time.Second)

	// issuer verifies
	fmt.Println("verifying presentation🙏")
	resp, err := issuer.VerifyPresentationByID(sendPresentation.PresentationExchangeID)
	if err != nil {
		log.Fatal(err)
	}

	parsed, _ := json.Marshal(resp)
	fmt.Println(string(parsed))
}
