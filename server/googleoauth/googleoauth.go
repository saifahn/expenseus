package googleoauth

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/saifahn/expenseus"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	defaultConfig = &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}
	//TODO: randomize
	oauthStateString = ""
)

type googleUserInfo struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Verified bool   `json:"verified_email"`
}

type GoogleOauthConfig struct {
	config oauth2.Config
}

func New() *GoogleOauthConfig {
	return &GoogleOauthConfig{config: *defaultConfig}
}

func (g *GoogleOauthConfig) AuthCodeURL(state string, opts ...oauth2.AuthCodeOption) string {
	return g.config.AuthCodeURL(state)
}

func (g *GoogleOauthConfig) Exchange(ctx context.Context, code string, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error) {
	token, err := g.config.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("code exchange failed: %s", err.Error())
	}

	return token, nil
}

func (g *GoogleOauthConfig) GetInfoAndGenerateUser(state string, code string) (expenseus.User, error) {
	// TODO: use pointer to User so I can return nil instead of empty struct literal
	if state != oauthStateString {
		return expenseus.User{}, fmt.Errorf("invalid oauth state")
	}

	token, err := g.Exchange(context.Background(), code)
	if err != nil {
		return expenseus.User{}, fmt.Errorf("code exchange failed: %s", err.Error())
	}

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return expenseus.User{}, fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer response.Body.Close()

	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return expenseus.User{}, fmt.Errorf("failed reading response body: %s", err.Error())
	}

	// convert the contents into a google user struct
	var googleUser googleUserInfo
	err = json.Unmarshal(contents, &googleUser)
	if err != nil {
		panic(err)
	}

	if !googleUser.Verified {
		// TODO: define an error type
		return expenseus.User{}, fmt.Errorf("user is not verified")
	}
	// use that information to create an expenseus.User
	return expenseus.User{
		ID:       googleUser.ID,
		Name:     googleUser.Name,
		Username: googleUser.Email,
	}, nil
}
