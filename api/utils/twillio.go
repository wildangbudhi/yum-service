package utils

import "github.com/twilio/twilio-go"

func NewTwillioConnection(username, password string) *twilio.RestClient {

	var client *twilio.RestClient

	client = twilio.NewRestClientWithParams(twilio.RestClientParams{
		Username: username,
		Password: password,
	})

	return client
}
