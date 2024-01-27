// Copyright (c) 2024 Alexej Disterhoft
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
//
// SPDX-License-Identifier: MIT

// Package ticket provides functionality to work with MapR tickets, including parsing tickets either
// from their raw string representation or from Kubernetes secrets.
//
// The package relies on https://pkg.go.dev/github.com/nobbs/mapr-ticket-parser for the actual
// ticket parsing. Most of the functionality in this package is just a wrapper around the parser to
// add some convenience methods.
package ticket

import (
	"errors"
	"fmt"
	"time"

	"github.com/nobbs/kubectl-mapr-ticket/pkg/util"
	"github.com/nobbs/mapr-ticket-parser/pkg/parse"

	coreV1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/json"
	"sigs.k8s.io/yaml"
)

const (
	// SecretMaprTicketKey is the key used for MapR tickets in secrets
	SecretMaprTicketKey = "CONTAINER_TICKET"

	// DefaultTimeFormat is the default time format used for human readable time strings
	DefaultTimeFormat = time.RFC3339
)

// SecretContainsMaprTicket returns true if the secret contains the key typically
// used for MapR tickets
type ErrSecretDoesNotContainMaprTicket struct {
	Name      string
	Namespace string
}

// NewErrSecretDoesNotContainMaprTicket returns a new ErrSecretDoesNotContainMaprTicket
func NewErrSecretDoesNotContainMaprTicket(namespace, name string) ErrSecretDoesNotContainMaprTicket {
	return ErrSecretDoesNotContainMaprTicket{
		Name:      name,
		Namespace: namespace,
	}
}

// Error returns the error message for ErrSecretDoesNotContainMaprTicket
func (err ErrSecretDoesNotContainMaprTicket) Error() string {
	return fmt.Sprintf("secret %s/%s does not contain a MapR ticket", err.Namespace, err.Name)
}

// SecretContainsMaprTicket returns true if the secret contains the key typically
// used for MapR tickets
func SecretContainsMaprTicket(secret *coreV1.Secret) bool {
	_, ok := secret.Data[SecretMaprTicketKey]
	return ok
}

// Wrapper around parse.Ticket to add methods
type Ticket parse.MaprTicket

// NewMaprTicket returns a new empty MaprTicket
func NewMaprTicket() *Ticket {
	return (*Ticket)(parse.NewMaprTicket())
}

// NewMaprTicketFromSecret parses the ticket from the secret and returns it
func NewMaprTicketFromSecret(secret *coreV1.Secret) (*Ticket, error) {
	// get ticket from secret
	ticketBytes, ok := secret.Data[SecretMaprTicketKey]
	if !ok {
		return nil, NewErrSecretDoesNotContainMaprTicket(secret.Namespace, secret.Name)
	}

	// parse ticket
	ticket, err := parse.Unmarshal(ticketBytes)
	if err != nil {
		return nil, err
	}

	return (*Ticket)(ticket), nil
}

// NewMaprTicketFromBytes parses the ticket from the given bytes and returns it
func NewMaprTicketFromBytes(ticketBytes []byte) (*Ticket, error) {
	// try to parse ticket directly
	ticket, errTicket := parseTicket(ticketBytes)
	if errTicket == nil {
		return ticket, nil
	}

	// try to parse as secret
	ticket, errSecret := parseSecret(ticketBytes)
	if errSecret == nil {
		return ticket, nil
	}

	// if we get here, we couldn't parse the ticket
	return nil, errors.Join(errTicket, errSecret)
}

// GetCluster returns the cluster that the ticket is for
func (ticket *Ticket) GetCluster() string {
	if ticket == nil {
		return ""
	}

	return ticket.Cluster
}

// GetUser returns the user that the ticket is for
func (ticket *Ticket) GetUser() string {
	if ticket == nil || ticket.UserCreds == nil || ticket.UserCreds.UserName == nil {
		return ""
	}

	return *ticket.UserCreds.UserName
}

// isExpired returns true if the ticket is expired
func (ticket *Ticket) IsExpired() bool {
	return time.Now().After(ticket.ExpirationTime())
}

// ExpirationTime returns the expiry time of the ticket as a time.Time object
func (ticket *Ticket) ExpirationTime() time.Time {
	return time.Unix(int64(ticket.GetExpiryTime()), 0)
}

// CreationTime returns the creation time of the ticket as a time.Time object
func (ticket *Ticket) CreationTime() time.Time {
	return time.Unix(int64(ticket.GetCreationTimeSec()), 0)
}

// ExpiresBefore returns true if the ticket expires before the given duration
func (ticket *Ticket) ExpiresBefore(duration time.Duration) bool {
	return ticket.ExpirationTime().Before(time.Now().Add(duration))
}

// AsMaprTicket returns the ticket as a parse.MaprTicket object
func (ticket *Ticket) AsMaprTicket() *parse.MaprTicket {
	return (*parse.MaprTicket)(ticket)
}

// String returns a string representation of the ticket
func parseTicket(ticketBytes []byte) (*Ticket, error) {
	// try to parse ticket directly
	ticket, errPlain := parse.Unmarshal(ticketBytes)
	if errPlain == nil {
		return (*Ticket)(ticket), nil
	}

	// try to parse ticket as base64 encoded
	ticketBytes, errDecode := util.DecodeBase64(string(ticketBytes))
	ticket, errBase64 := parse.Unmarshal(ticketBytes)
	if errBase64 == nil {
		return (*Ticket)(ticket), nil
	}

	// if we get here, we couldn't parse the ticket
	return nil, errors.Join(errPlain, errDecode, errBase64)
}

// parseSecret parses the secret and returns the ticket if it contains one
func parseSecret(secretBytes []byte) (*Ticket, error) {
	// try to parse as YAML into a secret
	var secret coreV1.Secret
	var errYAML error
	var errJSON error

	if errYAML = yaml.Unmarshal(secretBytes, &secret); errYAML == nil {
		return NewMaprTicketFromSecret(&secret)
	}

	// try to parse as JSON into a secret
	if errJSON = json.Unmarshal(secretBytes, &secret); errJSON == nil {
		return NewMaprTicketFromSecret(&secret)
	}

	// if we get here, we couldn't parse the secret
	return nil, errors.Join(errYAML, errJSON)
}
