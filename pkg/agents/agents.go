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

func ChangeInvitations(issuer, holder *acapy.Client) error {
	invitation, err := issuer.CreateInvitation("createdByCode", true, false, false)
	if err != nil {
		log.Println(err)
		return err
	}

	parsedInvitation := parseInvitation(invitation.Invitation)

	time.Sleep(time.Second)
	holder.ReceiveInvitation(parsedInvitation, true)

	return nil
}
