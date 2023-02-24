package issuer

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/ldej/go-acapy-client"
)

var issuer = acapy.NewClient("http://172.24.188.202:8002/")

// issuer creates an invitation
func CreateInvitation(alias string, autoAccept bool, multiUse bool, public bool) (acapy.CreateInvitationResponse, error) {
	fmt.Println("Create Invitation")

	invitation, err := issuer.CreateInvitation(alias, autoAccept, multiUse, public)
	if err != nil {
		log.Fatal(err)
	}

	return invitation, err
}

// issuer register a schema
func RegisterSchema(name string, version string, attributes []string) (acapy.Schema, error) {
	//schemaName := fmt.Sprintf("schema-elton-%v", time.Now().Unix())

	schema, err := issuer.RegisterSchema(name, version, attributes)
	if err != nil {
		log.Default().Println(err)
	}

	return schema, err
}

// issuer creates a credential definition
func CreateCredentialDefinition(tag string, supportRevocation bool, revocationRegistrySize int, schemaID string) (string, error) {
	credentialDefinition, err := issuer.CreateCredentialDefinition(tag, supportRevocation, revocationRegistrySize, schemaID)
	if err != nil {
		log.Default().Println(err)
	}

	return credentialDefinition, err
}

// issuer issues a credential
func OfferCredentialV2(connectionID, credentialDefinition, comment string) (acapy.CredentialExchangeRecordResult, error) {
	attributes := []acapy.CredentialPreviewAttributeV2{
		{
			MimeType: "text/plain",
			Name:     "name",
			Value:    "Elton",
		},
		{
			MimeType: "text/plain",
			Name:     "age",
			Value:    "22",
		},
	}

	credentialPreview := acapy.NewCredentialPreviewV2(attributes)

	offerCredential, err := issuer.OfferCredentialV2(connectionID, credentialPreview, credentialDefinition, "first comment bellow")
	if err != nil {
		log.Fatal(err)
	}

	return offerCredential, err
}

// issuer sends a presentation
func PresentationRequestRequest(credentialDefinition string, invitation acapy.CreateInvitationResponse) (acapy.PresentationExchangeRecord, error) {
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

	return sendPresentation, err
}

// issuer verifies
func VerifyPresentationByID(sendPresentation acapy.PresentationExchangeRecord) ([]byte, error) {
	resp, err := issuer.VerifyPresentationByID(sendPresentation.PresentationExchangeID)
	if err != nil {
		log.Fatal(err)
	}

	parsed, _ := json.Marshal(resp)

	return parsed, err
}

// issuer query presentation ID
func GetPresentationExchangeByID(sendPresentation acapy.PresentationExchangeRecord) (string, error) {
	presentationByID, err := issuer.GetPresentationExchangeByID(sendPresentation.PresentationExchangeID)
	if err != nil {
		log.Fatal(err)
	}

	getPresentationExchangeParsed, _ := json.Marshal(presentationByID.State)

	stringParsed := string(getPresentationExchangeParsed)
	//fmt.Println(string(getPresentationExchangeParsed))

	return stringParsed, err
}
