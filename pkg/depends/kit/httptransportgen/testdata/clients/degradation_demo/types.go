// This is a generated source file. DO NOT EDIT
// Source: degradation_demo/types.go

package degradation_demo

import jwt "github.com/golang-jwt/jwt/v4"

type DemoApiResp struct {
	Info GithubComGolangJwtJwtV4RegisteredClaims `json:"info"`
}

type GithubComGolangJwtJwtV4RegisteredClaims = jwt.RegisteredClaims
