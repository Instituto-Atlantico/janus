package agents

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/Instituto-Atlantico/go-acapy-client"
)

func parseInvitation(invitation any) acapy.Invitation {
	parsedInvitation, _ := json.Marshal(invitation)

	var invitationTyped acapy.Invitation

	json.Unmarshal(parsedInvitation, &invitationTyped)

	return invitationTyped
}

func ChangeInvitations(issuer, holder *acapy.Client) (string, string, error) {
	invitationResponse, err := issuer.CreateInvitation("createdByCode", true, false, false)
	if err != nil {
		return "", "", err
	}

	parsedInvitation := parseInvitation(invitationResponse.Invitation)

	time.Sleep(time.Second)
	connection, err := holder.ReceiveInvitation(parsedInvitation, true)
	if err != nil {
		return "", "", err
	}

	return invitationResponse.ConnectionID, connection.ConnectionID, nil
}

func CreateCredDef(issuer *acapy.Client, schemaId string) (string, error) {
	credentialDefinition, err := issuer.CreateCredentialDefinition("default", false, 0, schemaId)
	if err != nil {
		log.Println(err)
		return "", err
	}

	return credentialDefinition, nil
}

func GetCredDef(issuer *acapy.Client, schemaId string) (string, error) {
	credentialDefinitions, err := issuer.QueryCredentialDefinitions(acapy.QueryCredentialDefinitionsParams{
		SchemaID: schemaId,
	})
	if err != nil {
		log.Println(err)
		return "", err
	}

	if len(credentialDefinitions) == 0 {
		return "", errors.New("empty")
	}

	return credentialDefinitions[0], nil
}

func IssueCredential(issuer *acapy.Client, credentialDefinition string, connectionId string, attribute []acapy.CredentialPreviewAttributeV2) error {
	credentialPreview := acapy.NewCredentialPreviewV2(attribute)

	_, err := issuer.OfferCredentialV2(connectionId, credentialPreview, credentialDefinition, "Janus Credential")
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func GetCredential(holder *acapy.Client, attr, value string) (acapy.Credential, error) {
	query := fmt.Sprintf(`{"%s":"%s"}`, attr, value)

	credentials, err := holder.GetCredentials(1, 0, query)
	if err != nil {
		log.Fatal(err)
	}
	if len(credentials) == 0 {
		return acapy.Credential{}, errors.New("empty")
	}

	return credentials[0], nil
}

// issuer sends a presentation
func CreateRequestPresentationForSensor(issuer *acapy.Client, credentialDefinition, connectionId, sensorName string) (acapy.PresentationExchangeRecord, error) {
	requestedPredicates := map[string]acapy.RequestedPredicate{}

	requestedAttributes := map[string]acapy.RequestedAttribute{
		sensorName: acapy.RequestedAttribute{
			Restrictions: []map[string]string{
				{
					"cred_def_id": "credentialDefinition",
				},
				{
					fmt.Sprintf("attr::%s::value", sensorName): "true",
				},
			},
			Name: sensorName,
			NonRevoked: acapy.NonRevoked{
				From: time.Now().Add(-time.Hour * 24 * 7).Unix(),
				To:   time.Now().Unix(),
			},
		},
	}

	presentationRequestRequest := acapy.PresentationRequestRequest{
		Trace:        false,
		Comment:      "Presentation request test",
		ConnectionID: connectionId,
		ProofRequest: acapy.NewProofRequest(
			"SensorMeasurementAuthorization",
			"12515784361",
			requestedPredicates,
			requestedAttributes,
			"1.0",
			&acapy.NonRevoked{
				From: time.Now().Add(-time.Hour * 24 * 7).Unix(),
				To:   time.Now().Unix(),
			},
		),
	}

	sendPresentation, err := issuer.SendPresentationRequest(presentationRequestRequest)
	if err != nil {
		log.Fatal(err)
	}

	return sendPresentation, err
}

func IsPresentationDone(issuer *acapy.Client, ThreadID string) ([]acapy.PresentationExchangeRecord, error) {
	params := acapy.QueryPresentationExchangeParams{ThreadID: ThreadID, State: "presentation_received"}
	records, _ := issuer.QueryPresentationExchange(params)

	if len(records) == 0 {
		params = acapy.QueryPresentationExchangeParams{ThreadID: ThreadID, State: "abandoned"}
		records, _ = issuer.QueryPresentationExchange(params)
		if len(records) == 0 {
			return []acapy.PresentationExchangeRecord{}, errors.New("Empty")
		}
	}

	return records, nil
}

// holder send presentation
func SendPresentationByID(holder *acapy.Client, request acapy.PresentationExchangeRecord, credential acapy.Credential) error {
	proof := generatePresentationByRequest(request, credential.Referent)

	//query presentation ID
	param := acapy.QueryPresentationExchangeParams{
		ThreadID: request.ThreadID,
	}

	presentations, err := holder.QueryPresentationExchange(param)
	if err != nil {
		log.Fatal(err)
	}

	_, err = holder.SendPresentationByID(presentations[0].PresentationExchangeID, proof)
	if err != nil {
		log.Fatal(err)
	}

	return err
}

func generatePresentationByRequest(request acapy.PresentationExchangeRecord, credentialID string) acapy.PresentationProof {
	fields := request.PresentationRequest.RequestedAttributes

	attributesProof := map[string]acapy.PresentationProofAttribute{}
	for _, field := range fields {
		attributesProof[field.Name] = acapy.PresentationProofAttribute{CredentialID: credentialID,
			Revealed: true,
		}
	}

	predicatesProof := map[string]acapy.PresentationProofPredicate{}

	proof := acapy.PresentationProof{
		SelfAttestedAttributes: map[string]string{},
		RequestedAttributes:    attributesProof,
		RequestedPredicates:    predicatesProof,
	}

	return proof
}

// issuer verifies
func VerifyPresentationByID(issuer *acapy.Client, sendPresentation acapy.PresentationExchangeRecord) (acapy.PresentationExchangeRecord, error) {
	resp, err := issuer.VerifyPresentationByID(sendPresentation.PresentationExchangeID)
	if err != nil {
		log.Fatal(err)
	}

	return resp, err
}
