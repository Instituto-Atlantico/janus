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

func ChangeInvitations(issuer, holder *acapy.Client) (acapy.Invitation, acapy.Connection, error) {
	invitation, err := issuer.CreateInvitation("createdByCode", true, false, false)
	if err != nil {
		log.Println(err)
		return acapy.Invitation{}, acapy.Connection{}, err
	}

	parsedInvitation := parseInvitation(invitation.Invitation)

	time.Sleep(time.Second)
	connection, err := holder.ReceiveInvitation(parsedInvitation, true)
	if err != nil {
		log.Println(err)
		return acapy.Invitation{}, connection, err
	}

	return parsedInvitation, connection, nil
}

func CreateCredDef(issuer *acapy.Client, schemaId string) (string, error) {
	credentialDefinition, err := issuer.CreateCredentialDefinition("default", false, 0, schemaId)
	if err != nil {
		log.Println(err)
		return "", err
	}

	return credentialDefinition, nil
}

func issuerCredential(issuer *acapy.Client, credentialDefinition string, connectionId string, attribute []acapy.CredentialPreviewAttributeV2) error {
	credentialPreview := acapy.NewCredentialPreviewV2(attribute)

	_, err := issuer.OfferCredentialV2(connectionId, credentialPreview, credentialDefinition, "Janus Credential")
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
