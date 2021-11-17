package login

import (
	"github.com/looplab/fsm"
	"github.com/nrc-no/core/pkg/server/login/authrequest"
)

func handleAwaitingIDPCallback() func(authRequest *authrequest.AuthRequest, evt *fsm.Event) error {
	return func(authRequest *authrequest.AuthRequest, evt *fsm.Event) error {
		return nil
	}
}
