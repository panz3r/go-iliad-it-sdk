/*
 * Copyright (c) 2018 Mattia Panzeri <mattia.panzeri93@gmail.com>
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/.
 */

package iliad

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
	AvailableCredit         string
	CallTime                string
	MessagesCount           string
	MultimediaMessagesCount string
	InternetTraffic         string
}

// UserOptionsStatus contains Options activation statuses for a User
type UserOptionsStatus struct {
	LTE                                 bool
	PremiumNumbers                      bool
	OverThresholdInternetTraffic        bool
	LocalOverThresholdInternetTraffic   bool
	RoamingOverThresholdInternetTraffic bool
	ShowLast3PhoneNumberDigits          bool
}

// UserServicesStatus contains Services activation statuses for a User
type UserServicesStatus struct {
	BlockHiddenNumbers  bool
	RoamingVoicemail    bool
	BlockRedirect       bool
	AppearAsAbsent      bool
	QuickNumbers        bool
	CallsMessagesFilter bool
}
