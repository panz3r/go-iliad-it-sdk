/*
 * Copyright (c) 2018 Mattia Panzeri <mattia.panzeri93@gmail.com>
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/.
 */

package iliad

import (
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func getPageWithToken(url string, userToken string) (*http.Response, error) {
	// Create HTTP client
	client := &http.Client{}

	// Create and modify HTTP request before sending
	request, reqErr := http.NewRequest("GET", url, nil)
	if reqErr != nil {
		log.Fatal(reqErr)
		return nil, reqErr
	}

	request.Header.Set("cookie", "ACCOUNT_SESSID="+userToken)

	// Make request
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return response, nil
}

func postForm(url string, formData url.Values) (*http.Response, error) {
	return postFormWithToken(url, formData, "")
}

func postFormWithToken(url string, formData url.Values, userToken string) (*http.Response, error) {
	// Create HTTP client
	client := &http.Client{}

	// Create and modify HTTP request before sending
	req, err := http.NewRequest("POST", url, strings.NewReader(formData.Encode()))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(formData.Encode())))

	if userToken != "" {
		req.Header.Set("cookie", "ACCOUNT_SESSID="+userToken)
	}

	// Make request
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return res, nil
}

func getCookieByName(cookie []*http.Cookie, name string) string {
	cookieLen := len(cookie)
	result := ""
	for i := 0; i < cookieLen; i++ {
		if cookie[i].Name == name {
			result = cookie[i].Value
		}
	}
	return result
}
