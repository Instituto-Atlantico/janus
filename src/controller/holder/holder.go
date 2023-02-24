package holder

import (
	"encoding/json"
	"log"

	"github.com/ldej/go-acapy-client"
)

var holder = acapy.NewClient("http://172.24.188.202:9002/")

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
func SendPresentationByID(connection acapy.Connection) error {
	//query presentation ID
	param := acapy.QueryPresentationExchangeParams{
		ConnectionID: connection.ConnectionID,
	}

	presentations, err := holder.QueryPresentationExchange(param)
	if err != nil {
		log.Fatal(err)
	}

	// query Credential ID
	credentials, err := holder.GetCredentials(1, 0, "")
	if err != nil {
		log.Fatal(err)
	}

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

	_, err = holder.SendPresentationByID(presentations[0].PresentationExchangeID, proof)
	if err != nil {
		log.Fatal(err)
	}

	return err
}
