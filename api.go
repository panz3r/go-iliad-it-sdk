package iliad

import (
	"errors"
	"log"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var baseURL = "https://www.iliad.it/account/"

// Client allows to interact with Iliad IT services
type Client struct {
	userToken string
}

// NewClient initializes a new Iliad IT client
func NewClient() Client {
	return Client{}
}

// NewClientWithToken initializes a new Iliad IT client with a previously obtained user token
func NewClientWithToken(token string) Client {
	clt := Client{}
	clt.userToken = token
	return clt
}

// Login is used to get an auth token for the current user to access Iliad IT services
func (*Client) Login(username string, password string) (string, error) {
	page, err := postForm(baseURL, url.Values{
		"login-ident": {username},
		"login-pwd":   {password},
	})
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	defer page.Body.Close()

	// Create a goquery document from the HTTP response
	document, gqErr := goquery.NewDocumentFromReader(page.Body)
	if gqErr != nil {
		log.Fatal("Error loading HTTP response body. ", gqErr)
		return "", gqErr
	}

	iErr := getError(document)
	if iErr != nil {
		log.Fatal(iErr)
		return "", iErr
	}

	cookies := page.Cookies()
	// log.Printf("Cookies: %+v", cookies)
	tknCookie := getCookieByName(cookies, "ACCOUNT_SESSID")
	if tknCookie == "" {
		log.Fatal("Cookie token not found")
		return "", errors.New("Token not available")
	}

	return tknCookie, nil
}

// Private Methods

func getError(doc *goquery.Document) error {
	iliadErr := strings.Split(strings.TrimSpace(doc.Find("div.flash.flash-error").Text()), "\n")[0]
	if iliadErr != "" {
		return errors.New(iliadErr)
	}

	return nil
}
