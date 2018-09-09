/*
 * Copyright (c) 2018 Mattia Panzeri <mattia.panzeri93@gmail.com>
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/.
 */

package iliad

import (
	"errors"
	"log"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var baseURL = "https://www.iliad.it/account/"

var creditPageEndpoint = "consumi-e-credito"

// Client allows to interact with Iliad IT services
type Client struct {
	userToken string
}

// UserInfo contains general info about a User
type UserInfo struct {
	ID          string
	Name        string
	PhoneNumber string
}

// UserCreditInfo contains CreditInfo (such as call time, messages count, internet traffic, etc.) for a User
type UserCreditInfo struct {
	CallTime                string
	MessagesCount           string
	MultimediaMessagesCount string
	InternetTraffic         string
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

// Login can be used to get an auth token for the current user to access Iliad IT services
func (clt *Client) Login(username string, password string) (string, error) {
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

	clt.userToken = tknCookie
	return tknCookie, nil
}

// GetUserInfo can be used to retrieve basic info about the current User
func (clt *Client) GetUserInfo() (UserInfo, error) {
	usrInfo := UserInfo{}

	page, err := getPageWithToken(baseURL+creditPageEndpoint, clt.userToken)
	if err != nil {
		return usrInfo, err
	}
	defer page.Body.Close()

	// Create a goquery document from the HTTP response
	document, gqErr := goquery.NewDocumentFromReader(page.Body)
	if gqErr != nil {
		log.Fatal("Error loading HTTP response body. ", gqErr)
		return usrInfo, gqErr
	}

	usrInfo.Name = document.Find("div.current-user .bold").First().Text()
	usrInfo.ID = strings.TrimSpace(strings.Split(document.Find("div.current-user .smaller").First().Text(), ":")[1])
	usrInfo.PhoneNumber = strings.TrimSpace(strings.Split(document.Find("div.current-user .smaller").Last().Text(), ":")[1])

	return usrInfo, nil
}

// GetUserCreditInfo can be used to retrieve Credit info for the current User
func (clt *Client) GetUserCreditInfo() (UserCreditInfo, error) {
	crdInfo := UserCreditInfo{}

	page, err := getPageWithToken(baseURL+creditPageEndpoint, clt.userToken)
	if err != nil {
		return crdInfo, err
	}
	defer page.Body.Close()

	// Create a goquery document from the HTTP response
	document, gqErr := goquery.NewDocumentFromReader(page.Body)
	if gqErr != nil {
		log.Fatal("Error loading HTTP response body. ", gqErr)
		return crdInfo, gqErr
	}

	document.Find("div.conso-infos .conso__text .red").Each(func(i int, s *goquery.Selection) {
		switch i {
		case 0:
			crdInfo.CallTime = s.Text()
			break
		case 2:
			crdInfo.MessagesCount = s.Text()
			break
		case 4:
			crdInfo.InternetTraffic = s.Text()
			break
		case 6:
			crdInfo.MultimediaMessagesCount = s.Text()
		}
	})

	return crdInfo, nil
}

// Private Methods

func getError(doc *goquery.Document) error {
	iliadErr := strings.Split(strings.TrimSpace(doc.Find("div.flash.flash-error").Text()), "\n")[0]
	if iliadErr != "" {
		return errors.New(iliadErr)
	}

	return nil
}
