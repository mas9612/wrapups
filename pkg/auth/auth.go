package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	pb "github.com/mas9612/authserver/pkg/authserver"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

type config struct {
	AuthserverURL string `json:"authserver_url"`
}

var (
	// ErrNotConfigured is the error which indicates the user credential has not configured yet.
	ErrNotConfigured = errors.New("credential not configured")
)

type credential struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

// Token returns the token issued with configured credentials.
func Token() (string, error) {
	userHome, err := os.UserHomeDir()
	if err != nil {
		panic("failed to get user home")
	}

	wrapupsDir := path.Join(userHome, ".wrapups")
	if _, err := os.Stat(wrapupsDir); err != nil {
		if os.IsNotExist(err) {
			os.Mkdir(wrapupsDir, 0700)
		}
	}

	configPath := path.Join(wrapupsDir, "config")
	var wrapupsConfig config
	if _, err = os.Stat(configPath); err == nil {
		b, err := ioutil.ReadFile(configPath)
		if err != nil {
			panic("failed to read config")
		}
		if err := json.Unmarshal(b, &wrapupsConfig); err != nil {
			panic("failed to parse config file")
		}
	}
	if wrapupsConfig.AuthserverURL == "" {
		wrapupsConfig.AuthserverURL = "localhost:10000"
	}

	// check token
	tokenPath := path.Join(wrapupsDir, "token")
	_, err = os.Stat(tokenPath)
	if err == nil { // token file exists
		b, err := ioutil.ReadFile(tokenPath)
		if err != nil {
			return "", errors.Wrap(err, "failed to read token file")
		}
		return string(b), nil
	}

	credentialPath := path.Join(wrapupsDir, "credential")
	if _, err := os.Stat(credentialPath); err != nil {
		if os.IsNotExist(err) {
			return "", ErrNotConfigured
		}
	}
	b, err := ioutil.ReadFile(credentialPath)
	if err != nil {
		return "", errors.Wrap(err, "failed to read credential file")
	}
	var wrapupCredential credential
	if err := json.Unmarshal(b, &wrapupCredential); err != nil {
		return "", errors.Wrap(err, "failed to parse credential config")
	}

	conn, err := grpc.Dial(wrapupsConfig.AuthserverURL, grpc.WithInsecure())
	if err != nil {
		return "", errors.Wrap(err, "failed to create gRPC client")
	}
	client := pb.NewAuthserverClient(conn)
	req := &pb.CreateTokenRequest{
		User:     wrapupCredential.User,
		Password: wrapupCredential.Password,
		OrigHost: "wrapups",
	}
	token, err := client.CreateToken(context.Background(), req)
	if err != nil {
		return "", errors.Wrap(err, "failed to issue new token")
	}

	// save token to ~/.wrapups/token
	tokenFile, err := os.OpenFile(tokenPath, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return "", errors.Wrap(err, "failed to save token")
	}
	fmt.Fprintf(tokenFile, token.Token)
	return token.Token, nil
}
