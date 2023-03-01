package issuer

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/ldej/go-acapy-client"
)

var issuer = acapy.NewClient("http://localhost:8002/")

func GetSchema(schemaName string) (string, error) {
	schemas, err := issuer.QuerySchemas(acapy.QuerySchemasParams{SchemaName: schemaName})
	if err != nil {
		log.Fatal(err)
	}
	if len(schemas) == 0 {
		return "", errors.New("empty")
	}

	return schemas[0], nil
}

func GetCredDef(schemaID string) (string, error) {
	credDefs, err := issuer.QueryCredentialDefinitions(acapy.QueryCredentialDefinitionsParams{SchemaID: schemaID})
	if err != nil {
		log.Fatal(err)
	}
	if len(credDefs) == 0 {
		return "", errors.New("empty")
	}

	return credDefs[0], nil
}

func GetConnection() (acapy.Connection, error) {
	conns, err := issuer.QueryConnections(&acapy.QueryConnectionsParams{})
	if err != nil {
		log.Fatal(err)
	}
	if len(conns) == 0 {
		return acapy.Connection{}, errors.New("empty")
	}

	return conns[0], nil
}

// issuer creates an invitation
func CreateInvitation(alias string, autoAccept bool, multiUse bool, public bool) (acapy.CreateInvitationResponse, error) {
	invitation, err := issuer.CreateInvitation(alias, autoAccept, multiUse, public)
	if err != nil {
		log.Fatal(err)
	}

	return invitation, err
}

// issuer register a schema
func RegisterSchema(name string, version string, attributes []string) (acapy.Schema, error) {
	schema, err := issuer.RegisterSchema(name, version, attributes)
	if err != nil {
		log.Default().Println(err)
	}

	return schema, err
}

func GetSchemaAttributes(id string) []string {
	schema, _ := issuer.GetSchema(id)

	return schema.AttributeNames
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
func OfferCredentialV2(connectionID, credentialDefinition, comment string, attributes []acapy.CredentialPreviewAttributeV2) (acapy.CredentialExchangeRecordResult, error) {
	credentialPreview := acapy.NewCredentialPreviewV2(attributes)

	offerCredential, err := issuer.OfferCredentialV2(connectionID, credentialPreview, credentialDefinition, comment)
	if err != nil {
		log.Fatal(err)
	}

	return offerCredential, err
}

// issuer sends a presentation
func PresentationRequestRequest(credentialDefinition string, invitation acapy.Connection) (acapy.PresentationExchangeRecord, error) {
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

	getPresentationExchangeParsed, _ := json.Marshal(presentationByID)

	stringParsed := string(getPresentationExchangeParsed)

	return stringParsed, err
}

func GetPresentationExchangeByThreadId(id string) (acapy.PresentationExchangeRecord, error) {
	presentations, _ := issuer.QueryPresentationExchange(acapy.QueryPresentationExchangeParams{ThreadID: id, State: "verified"})

	if len(presentations) == 0 {
		return acapy.PresentationExchangeRecord{}, errors.New("no verified presentation found")
	}

	return presentations[0], nil
}
