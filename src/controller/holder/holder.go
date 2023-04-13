package holder

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/ldej/go-acapy-client"
)

var holder = acapy.NewClient("http://raspberrypi:8002/")

func GetConnection() (acapy.Connection, error) {
	conns, err := holder.QueryConnections(&acapy.QueryConnectionsParams{State: "active"})
	if err != nil {
		log.Fatal(err)
	}
	if len(conns) == 0 {
		return acapy.Connection{}, errors.New("empty")
	}

	return conns[0], nil
}

func GetCredential(attr, value string) (acapy.Credential, error) {
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

// holder receive the invitation
func ReceiveInvitation(invitation acapy.CreateInvitationResponse, autoAccept bool) (acapy.Connection, error) {
	parsedInvitation, _ := json.Marshal(invitation.Invitation)

	var invitationTyped acapy.Invitation

	json.Unmarshal(parsedInvitation, &invitationTyped)

	connection, err := holder.ReceiveInvitation(invitationTyped, autoAccept)
	if err != nil {
		log.Fatal(err)
	}

	return connection, err
}

func GeneratePresentationByRequest(request acapy.PresentationExchangeRecord, credentialID string) acapy.PresentationProof {
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

// holder send presentation
func SendPresentationByID(request acapy.PresentationExchangeRecord, credential acapy.Credential) error {
	proof := GeneratePresentationByRequest(request, credential.Referent)

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