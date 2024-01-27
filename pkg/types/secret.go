// Copyright (c) 2024 Alexej Disterhoft
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
//
// SPDX-License-Identifier: MIT

package types

import (
	"fmt"
	"time"

	"github.com/nobbs/kubectl-mapr-ticket/pkg/ticket"
	"github.com/nobbs/kubectl-mapr-ticket/pkg/util"

	coreV1 "k8s.io/api/core/v1"
)

// Secret is a wrapper around coreV1.Secret that provides additional functionality.
type Secret coreV1.Secret

// MaprSecret is used to store a secret and its corresponding ticket parsed from that secret
type MaprSecret struct {
	Secret *Secret        `json:"secret"`
	Ticket *ticket.Ticket `json:"ticket"`
	NumPVC uint32         `json:"-"`
}

// NewMaprSecret creates a new MaprSecret from a Secret
func NewMaprSecret(s *Secret) *MaprSecret {
	if s == nil {
		return &MaprSecret{}
	}

	v := &MaprSecret{
		Secret: s,
	}

	ticket, err := ticket.NewMaprTicketFromSecret((*coreV1.Secret)(s))
	if err != nil {
		return v
	}

	v.Ticket = ticket

	return v
}

// GetSecretName returns the name of the secret
func (t *MaprSecret) GetSecretName() string {
	if t == nil || t.Secret == nil {
		return ""
	}

	return t.Secret.GetName()
}

// GetSecretNamespace returns the namespace of the secret
func (t *MaprSecret) GetSecretNamespace() string {
	if t == nil || t.Secret == nil {
		return ""
	}

	return t.Secret.GetNamespace()
}

// GetCluster returns the cluster of the ticket
func (t *MaprSecret) GetCluster() string {
	if t == nil || t.Ticket == nil {
		return ""
	}

	return t.Ticket.GetCluster()
}

// GetUser returns the user of the ticket
func (t *MaprSecret) GetUser() string {
	if t == nil || t.Ticket == nil {
		return ""
	}

	return t.Ticket.GetUser()
}

// GetExpirationTime returns the expiration time of the ticket
func (t *MaprSecret) GetExpirationTime() time.Time {
	if t == nil || t.Ticket == nil {
		return time.Time{}
	}

	return t.Ticket.ExpirationTime()
}

// GetCreationTime returns the creation time of the ticket
func (t *MaprSecret) GetCreationTime() time.Time {
	if t == nil || t.Ticket == nil {
		return time.Time{}
	}

	return t.Ticket.CreationTime()
}

// GetStatusString returns a human readable string describing the status of the ticket
func (t *MaprSecret) GetStatusString() string {
	if t == nil {
		return "Not found / Invalid"
	}

	if t.Secret == nil {
		return "No secret found"
	}

	if t == nil || t.Ticket == nil {
		return "No ticket found"
	}

	if t.Ticket.IsExpired() {
		return fmt.Sprintf("Expired (%s ago)", util.ShortHumanDurationComparedToNow(t.Ticket.ExpirationTime()))
	}

	return fmt.Sprintf("Valid (%s left)", util.ShortHumanDurationComparedToNow(t.Ticket.ExpirationTime()))
}
