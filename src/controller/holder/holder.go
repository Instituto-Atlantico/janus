package holder

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/ldej/go-acapy-client"
)

var holder = acapy.NewClient("http://raspberrypi.local:8002/")

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
	query := fmt.Sprintf(`{"attr::%s::value": "%s"}`, attr, value)

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

// holder send presentation
func SendPresentationByID(threadID string, credential acapy.Credential) error {
	//query presentation ID
	param := acapy.QueryPresentationExchangeParams{
		ThreadID: threadID,
	}

	presentations, err := holder.QueryPresentationExchange(param)
	if err != nil {
		log.Fatal(err)
	}

	attributesProof := map[string]acapy.PresentationProofAttribute{
		"name": {
			CredentialID: credential.Referent,
			Revealed:     true,
		},
	}

	predicatesProof := map[string]acapy.PresentationProofPredicate{
		"age": {
			CredentialID: credential.Referent,
		},
	}

	proof := acapy.PresentationProof{
		SelfAttestedAttributes: map[string]string{},
		RequestedAttributes:    attributesProof,
		RequestedPredicates:    predicatesProof,
	}

	_, err = holder.SendPresentationByID(presentations[0].PresentationExchangeID, proof)
	if err != nil {
		log.Fatal(err)
	}

	return err
}
