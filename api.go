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
var optionsPageEndpoint = "le-mie-opzioni"
var servicesPageEndpoint = "i-miei-servizi"

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
	// Format user credentials as FormData
	authForm := url.Values{
		"login-ident": {username},
		"login-pwd":   {password},
	}

	// First request to authenticated user
	page, err := postForm(baseURL, authForm)
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	defer page.Body.Close()

	document, gqErr := goquery.NewDocumentFromReader(page.Body)
	if gqErr != nil {
		log.Fatal("Error loading HTTP response body. ", gqErr)
		return "", gqErr
	}

	// Check for submitted data errors
	iErr := getError(document)
	if iErr != nil {
		log.Fatal(iErr)
		return "", iErr
	}

	// Retrieve user token for current session
	cookies := page.Cookies()
	tknCookie := getCookieByName(cookies, "ACCOUNT_SESSID")
	if tknCookie == "" {
		log.Fatal("Cookie token not found")
		return "", errors.New("Token not available")
	}

	// Second request to activate session token (otherwise token won't be valid)
	aPage, err := postFormWithToken(baseURL, authForm, tknCookie)
	if err != nil {
		log.Fatal("Confirm token request failed", err)
		return "", err
	}

	aDoc, err := goquery.NewDocumentFromReader(aPage.Body)
	if err != nil {
		log.Fatal("Error loading HTTP response body. ", err)
		return "", err
	}

	iaErr := getError(aDoc)
	if iaErr != nil {
		log.Fatal(iaErr)
		return "", iaErr
	}

	// Set Client user token and return it
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

	document, gqErr := goquery.NewDocumentFromReader(page.Body)
	if gqErr != nil {
		log.Fatal("Error loading HTTP response body. ", gqErr)
		return crdInfo, gqErr
	}

	crdInfo.AvailableCredit = document.Find("div.page.p-conso h2 b.red").First().Text()

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

// GetUserOptions can be used to retrieve Options status for the current User
func (clt *Client) GetUserOptions() (UserOptionsStatus, error) {
	optsSts := UserOptionsStatus{}

	page, err := getPageWithToken(baseURL+optionsPageEndpoint, clt.userToken)
	if err != nil {
		return optsSts, err
	}
	defer page.Body.Close()

	document, gqErr := goquery.NewDocumentFromReader(page.Body)
	if gqErr != nil {
		log.Fatal("Error loading HTTP response body. ", gqErr)
		return optsSts, gqErr
	}

	document.Find("div.array-status div.grid-l.as__item").Each(func(i int, o *goquery.Selection) {
		status := o.Find("div.as__status.as__status--on.as__status--active").Length() > 0
		switch i {
		case 0:
			optsSts.LTE = status
			break
		case 1:
			optsSts.PremiumNumbers = status
			break
		case 2:
			optsSts.OverThresholdInternetTraffic = status
			break
		case 3:
			optsSts.LocalOverThresholdInternetTraffic = status
			break
		case 4:
			optsSts.RoamingOverThresholdInternetTraffic = status
			break
		case 5:
			optsSts.ShowLast3PhoneNumberDigits = status
			break
		}
	})

	return optsSts, nil
}

// GetUserServices can be used to retrieve Services status for the current User
func (clt *Client) GetUserServices() (UserServicesStatus, error) {
	svcsSts := UserServicesStatus{}

	page, err := getPageWithToken(baseURL+servicesPageEndpoint, clt.userToken)
	if err != nil {
		return svcsSts, err
	}
	defer page.Body.Close()

	document, gqErr := goquery.NewDocumentFromReader(page.Body)
	if gqErr != nil {
		log.Fatal("Error loading HTTP response body. ", gqErr)
		return svcsSts, gqErr
	}

	document.Find("div.array-status div.grid-l.as__item").Each(func(i int, o *goquery.Selection) {
		status := o.Find("div.as__status.as__status--on.as__status--active").Length() > 0

		switch i {
		case 0:
			svcsSts.BlockHiddenNumbers = status
			break
		case 1:
			svcsSts.RoamingVoicemail = status
			break
		case 2:
			svcsSts.BlockRedirect = status
			break
		case 3:
			svcsSts.AppearAsAbsent = status
			break
		case 4:
			svcsSts.QuickNumbers = status
			break
		case 5:
			svcsSts.CallsMessagesFilter = status
			break
		}
	})

	return svcsSts, nil
}

// Private Methods

func getError(doc *goquery.Document) error {
	iliadErr := strings.Split(strings.TrimSpace(doc.Find("div.flash.flash-error").Text()), "\n")[0]
	if iliadErr != "" {
		return errors.New(iliadErr)
	}

	return nil
}
