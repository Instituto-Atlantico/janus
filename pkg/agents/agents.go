package agents

import (
	"encoding/json"
	"log"
	"time"

	"github.com/ldej/go-acapy-client"
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
		log.Println(err)
		return "", "", err
	}

	parsedInvitation := parseInvitation(invitationResponse.Invitation)

	time.Sleep(time.Second)
	connection, err := holder.ReceiveInvitation(parsedInvitation, true)
	if err != nil {
		log.Println(err)
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

func IssueCredential(issuer *acapy.Client, credentialDefinition string, connectionId string, attribute []acapy.CredentialPreviewAttributeV2) error {
	credentialPreview := acapy.NewCredentialPreviewV2(attribute)

	_, err := issuer.OfferCredentialV2(connectionId, credentialPreview, credentialDefinition, "Janus Credential")
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
