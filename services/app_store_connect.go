package services

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type AppStoreConnect struct {
	cfg AppStoreConnectConfig
}

func NewAppStoreConnect(cfg AppStoreConnectConfig) *AppStoreConnect {
	return &AppStoreConnect{
		cfg: cfg,
	}
}

type AppStoreConnectConfig struct {
	KeyID   string `mapstructure:"key_id"`
	PrivKey string `mapstructure:"priv_key"`
	Issuer  string
}

type AppStoreConnectIdentity struct {
}

// Implement Service for AppStoreConnect.

func (g AppStoreConnect) GroupMembers(group string) ([]User, error) {
	token, err := g.getToken()
	if err != nil {
		return []User{}, err
	}

	fmt.Println(token)

	client := &http.Client{}

	req, err := http.NewRequest(
		"GET",
		"https://api.appstoreconnect.apple.com/v1/users",
		nil,
	)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	fmt.Printf("%s\n", body)

	return []User{}, nil
}

func (g AppStoreConnect) getToken() (string, error) {
	// ACS requires tokens to be encrypted strictly with ES256
	//token := jwt.New(jwt.SigningMethodES256)
	token := jwt.NewWithClaims(
		jwt.SigningMethodES256,
		&jwt.StandardClaims{
			ExpiresAt: time.Now().
				Add(time.Minute * 20).
				UTC().
				Unix(),
			Issuer:   g.cfg.Issuer,
			Audience: "appstoreconnect-v1",
		},
	)

	var block *pem.Block
	if block, _ = pem.Decode([]byte(g.cfg.PrivKey)); block == nil {
		panic("expected pem block")
	}

	key, err := x509.ParsePKCS8PrivateKey([]byte(block.Bytes))
	if err != nil {
		return "", err
	}

	token.Header["kid"] = g.cfg.KeyID

	return token.SignedString(key)
}
