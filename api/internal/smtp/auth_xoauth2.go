package smtp

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/emersion/go-sasl"
)

const XOAUTH2 = "XOAUTH2"
const XOAuth2SplitChar = 0x01 // ^A
const XOAuth2MaxParts = 3
const XOAuthUserKey = "user="
const XOAuthTokenKey = "auth="

type XOAuth2Error struct {
	Status  string `json:"status"`
	Schemes string `json:"schemes"`
	Scopes  string `json:"scopes"`
}

type XOAuth2Options struct {
	Username string
	Token    string
	Host     string
	Port     int
}

func (err *XOAuth2Error) Error() string {
	return fmt.Sprintf("XOAuth2 error (%v)", err.Status)
}

type XOAuth2Authenticator func(opts XOAuth2Options) error

type XOAuth2Server struct {
	done         bool
	failErr      error
	authenticate XOAuth2Authenticator
}

func (a *XOAuth2Server) fail(descr string) ([]byte, bool, error) {
	blob, err := json.Marshal(XOAuth2Error{
		Status:  "invalid_request",
		Schemes: "bearer",
	})
	if err != nil {
		panic(err) // wtf
	}
	a.failErr = errors.New("sasl: client error: " + descr)
	return blob, false, nil
}

func (a *XOAuth2Server) Next(response []byte) (challenge []byte, done bool, err error) {
	if a.done {
		err = sasl.ErrUnexpectedClientResponse
		return
	}

	if response == nil {
		return []byte{}, false, nil
	}

	// incoming format: base64("user=test@contoso.onmicrosoft.com^Aauth=Bearer EwBAAl3BAAUFFpUAo7J3Ve0bjLBWZWCclRC3EoAA^A^A")
	parts := bytes.SplitN(response, []byte{XOAuth2SplitChar}, XOAuth2MaxParts)
	if len(parts) != 3 {
		return a.fail("Invalid response")
	}
	user, found := bytes.CutPrefix(parts[0], []byte(XOAuthUserKey))
	if !found {
		return a.fail("Invalid response")
	}
	token, found := bytes.CutPrefix(parts[1], []byte(XOAuthTokenKey))

	if !found {
		return a.fail("Invalid response")
	}

	opts := XOAuth2Options{
		Username: string(user),
		Token:    string(token),
	}

	if len(opts.Username) == 0 || len(opts.Token) == 0 {
		return a.fail("Invalid response")
	}

	return nil, true, a.authenticate(opts)
}

func NewXOAuth2Server(auth XOAuth2Authenticator) sasl.Server {
	return &XOAuth2Server{authenticate: auth}
}
