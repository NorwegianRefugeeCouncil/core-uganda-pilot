package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"gopkg.in/square/go-jose.v2"
	jwt2 "gopkg.in/square/go-jose.v2/jwt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type Claimer interface {
	GetClaims() *jwt2.Claims
}

func TokenReader(ctx context.Context) func(req *http.Request, into Claimer) error {

	var jwks jose.JSONWebKeySet

	issuer := "http://localhost:4444"

	req, err := http.NewRequest("GET", issuer+"/.well-known/jwks.json", nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Accept", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	if res.StatusCode != http.StatusOK {
		panic("unexpected status code")
	}
	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(bodyBytes, &jwks); err != nil {
		panic(err)
	}

	return func(req *http.Request, into Claimer) error {
		authorization := req.Header.Get("Authorization")
		if len(authorization) == 0 {
			return fmt.Errorf("authorization header not present")
		}

		parts := strings.Split(authorization, " ")

		if parts[0] != "Bearer" {
			return fmt.Errorf("malformed authorization header")
		}

		token, err := jwt2.ParseSigned(parts[0])
		if err != nil {
			return fmt.Errorf("could not parse token: %v", err)
		}

		if err := token.Claims(jwks.Keys[0].Key, into); err != nil {
			return fmt.Errorf("could not retrieve token claims: %v", err)
		}

		if err := into.GetClaims().Validate(jwt2.Expected{
			Issuer:   issuer,
			Audience: []string{},
			Time:     time.Now(),
		}); err != nil {
			return fmt.Errorf("invalid token: %v", err)
		}

		return nil

	}
}
