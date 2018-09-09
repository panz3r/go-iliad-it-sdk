package iliad

import (
	"log"
	"net/http"
	"net/url"
)

func getPageWithToken(url string, userToken string) (*http.Response, error) {
	// Create HTTP client with timeout
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

	response, reqErr := http.PostForm(
		url,
		formData,
	)
	if reqErr != nil {
		log.Fatal(reqErr)
		return nil, reqErr
	}

	return response, nil
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
