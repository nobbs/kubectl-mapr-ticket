// Copyright (c) 2024 Alexej Disterhoft
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
//
// SPDX-License-Identifier: MIT

package secret_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	. "github.com/nobbs/kubectl-mapr-ticket/pkg/secret"
	"github.com/nobbs/kubectl-mapr-ticket/pkg/ticket"
	"github.com/nobbs/kubectl-mapr-ticket/pkg/types"
	"github.com/nobbs/mapr-ticket-parser/pkg/parse"

	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
)

func TestLister_Default(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		fields  listerFields
		want    []expectedSecret
		wantErr bool
	}{
		{
			name: "no secrets in default namespace",
			fields: listerFields{
				client:    fake.NewSimpleClientset(),
				namespace: "default",
			},
			want:    []expectedSecret{},
			wantErr: false,
		},
		{
			name: "one secret without ticket in default namespace",
			fields: listerFields{
				client: fake.NewSimpleClientset(&coreV1.Secret{
					ObjectMeta: metaV1.ObjectMeta{
						Name:      "test-secret",
						Namespace: "default",
					},
				}),
				namespace: "default",
			},
			want:    []expectedSecret{},
			wantErr: false,
		},
		{
			name: "secret with invalid ticket data",
			fields: listerFields{
				client: fake.NewSimpleClientset(&coreV1.Secret{
					ObjectMeta: metaV1.ObjectMeta{
						Name:      "test-secret",
						Namespace: "default",
					},
					Data: map[string][]byte{
						ticket.SecretMaprTicketKey: []byte("invalid ticket data"),
					},
				}),
				namespace: "default",
			},
			want:    []expectedSecret{},
			wantErr: false,
		},
		{
			name: "one secret with ticket in default namespace",
			fields: listerFields{
				client: fake.NewSimpleClientset(
					secretFromTicketJSON(
						t,
						"default",
						"test-secret",
						[]byte(`{"ticket":{"cluster":"test-cluster"}}`),
					),
				),
				namespace: "default",
			},
			want: []expectedSecret{
				newExpectedSecret("default", "test-secret"),
			},
			wantErr: false,
		},
		{
			name: "two secrets with ticket in kube-system namespace",
			fields: listerFields{
				client: fake.NewSimpleClientset(
					secretFromTicketJSON(
						t,
						"kube-system",
						"test-secret-1",
						[]byte(`{"ticket":{"cluster":"test-cluster-1"}}`),
					),
					secretFromTicketJSON(
						t,
						"kube-system",
						"test-secret-2",
						[]byte(`{"ticket":{"cluster":"test-cluster-2"}}`),
					),
				),
				namespace: "kube-system",
			},
			want: []expectedSecret{
				newExpectedSecret("kube-system", "test-secret-1"),
				newExpectedSecret("kube-system", "test-secret-2"),
			},
			wantErr: false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			l := NewLister(test.fields.client, test.fields.namespace, test.fields.opts...)

			got, err := l.List()

			assertTicketSecret(t, got, test.want)
			assert.Equal(t, test.wantErr, err != nil)
		})
	}
}

func TestLister_WithFilterByMaprCluster(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		fields  listerFields
		want    []expectedSecret
		wantErr bool
	}{
		{
			name: "one secret with ticket in default namespace",
			fields: listerFields{
				client: fake.NewSimpleClientset(
					secretFromTicketJSON(
						t,
						"default",
						"test-secret",
						[]byte(`{"cluster":"test-cluster"}`),
					),
				),
				namespace: "default",
				opts: []ListerOption{
					WithFilterByMaprCluster("test-cluster"),
				},
			},
			want: []expectedSecret{
				newExpectedSecret("default", "test-secret"),
			},
			wantErr: false,
		},
		{
			name: "one secret with ticket in default namespace, filter by different cluster",
			fields: listerFields{
				client: fake.NewSimpleClientset(
					secretFromTicketJSON(
						t,
						"default",
						"test-secret",
						[]byte(`{"cluster":"test-cluster"}`),
					),
				),
				namespace: "default",
				opts: []ListerOption{
					WithFilterByMaprCluster("test-cluster-2"),
				},
			},
			want:    []expectedSecret{},
			wantErr: false,
		},
		{
			name: "two secrets with ticket in kube-system namespace, filter by one cluster",
			fields: listerFields{
				client: fake.NewSimpleClientset(
					secretFromTicketJSON(
						t,
						"kube-system",
						"test-secret-1",
						[]byte(`{"cluster":"test-cluster-1"}`),
					),
					secretFromTicketJSON(
						t,
						"kube-system",
						"test-secret-2",
						[]byte(`{"cluster":"test-cluster-2"}`),
					),
				),
				namespace: "kube-system",
				opts: []ListerOption{
					WithFilterByMaprCluster("test-cluster-1"),
				},
			},
			want: []expectedSecret{
				newExpectedSecret("kube-system", "test-secret-1"),
			},
			wantErr: false,
		},
		{
			name: "two secrets with ticket in kube-system namespace, filter by different cluster",
			fields: listerFields{
				client: fake.NewSimpleClientset(
					secretFromTicketJSON(
						t,
						"kube-system",
						"test-secret-1",
						[]byte(`{"cluster":"test-cluster-1"}`),
					),
					secretFromTicketJSON(
						t,
						"kube-system",
						"test-secret-2",
						[]byte(`{"cluster":"test-cluster-2"}`),
					),
				),
				namespace: "kube-system",
				opts: []ListerOption{
					WithFilterByMaprCluster("test-cluster-3"),
				},
			},
			want:    []expectedSecret{},
			wantErr: false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			l := NewLister(test.fields.client, test.fields.namespace, test.fields.opts...)

			got, err := l.List()

			assertTicketSecret(t, got, test.want)
			assert.Equal(t, test.wantErr, err != nil)
		})
	}
}

func TestLister_WithFilterByMaprUser(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		fields  listerFields
		want    []expectedSecret
		wantErr bool
	}{
		{
			name: "one secret with ticket in default namespace",
			fields: listerFields{
				client: fake.NewSimpleClientset(
					secretFromTicketJSON(
						t,
						"default",
						"test-secret",
						[]byte(`{"ticket":{"userCreds":{"userName":"test-user"}}}`),
					),
				),
				namespace: "default",
				opts: []ListerOption{
					WithFilterByMaprUser("test-user"),
				},
			},
			want: []expectedSecret{
				newExpectedSecret("default", "test-secret"),
			},
			wantErr: false,
		},
		{
			name: "one secret with ticket in default namespace, filter by different user",
			fields: listerFields{
				client: fake.NewSimpleClientset(
					secretFromTicketJSON(
						t,
						"default",
						"test-secret",
						[]byte(`{"ticket":{"userCreds":{"userName":"test-user"}}}`),
					),
				),
				namespace: "default",
				opts: []ListerOption{
					WithFilterByMaprUser("test-user-2"),
				},
			},
			want:    []expectedSecret{},
			wantErr: false,
		},
		{
			name: "two secrets with ticket in kube-system namespace, filter by one user",
			fields: listerFields{
				client: fake.NewSimpleClientset(
					secretFromTicketJSON(
						t,
						"kube-system",
						"test-secret-1",
						[]byte(`{"ticket":{"userCreds":{"userName":"test-user-1"}}}`),
					),
					secretFromTicketJSON(
						t,
						"kube-system",
						"test-secret-2",
						[]byte(`{"ticket":{"userCreds":{"userName":"test-user-2"}}}`),
					),
				),
				namespace: "kube-system",
				opts: []ListerOption{
					WithFilterByMaprUser("test-user-1"),
				},
			},
			want: []expectedSecret{
				newExpectedSecret("kube-system", "test-secret-1"),
			},
			wantErr: false,
		},
		{
			name: "two secrets with ticket in kube-system namespace, filter by different user",
			fields: listerFields{
				client: fake.NewSimpleClientset(
					secretFromTicketJSON(
						t,
						"kube-system",
						"test-secret-1",
						[]byte(`{"ticket":{"userCreds":{"userName":"test-user-1"}}}`),
					),
					secretFromTicketJSON(
						t,
						"kube-system",
						"test-secret-2",
						[]byte(`{"ticket":{"userCreds":{"userName":"test-user-2"}}}`),
					),
				),
				namespace: "kube-system",
				opts: []ListerOption{
					WithFilterByMaprUser("test-user-3"),
				},
			},
			want:    []expectedSecret{},
			wantErr: false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			l := NewLister(test.fields.client, test.fields.namespace, test.fields.opts...)

			got, err := l.List()

			assertTicketSecret(t, got, test.want)
			assert.Equal(t, test.wantErr, err != nil)
		})
	}
}

func TestLister_WithFilterByUID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		fields  listerFields
		want    []expectedSecret
		wantErr bool
	}{
		{
			name: "one secret with ticket in default namespace",
			fields: listerFields{
				client: fake.NewSimpleClientset(
					secretFromTicketJSON(
						t,
						"default",
						"test-secret",
						[]byte(`{"ticket":{"userCreds":{"uid":1000}}}`),
					),
				),
				namespace: "default",
				opts: []ListerOption{
					WithFilterByUID(1000),
				},
			},
			want: []expectedSecret{
				newExpectedSecret("default", "test-secret"),
			},
			wantErr: false,
		},
		{
			name: "one secret with ticket in default namespace, filter by different uid",
			fields: listerFields{
				client: fake.NewSimpleClientset(
					secretFromTicketJSON(
						t,
						"default",
						"test-secret",
						[]byte(`{"ticket":{"userCreds":{"uid":1000}}}`),
					),
				),
				namespace: "default",
				opts: []ListerOption{
					WithFilterByUID(2000),
				},
			},
			want:    []expectedSecret{},
			wantErr: false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			l := NewLister(test.fields.client, test.fields.namespace, test.fields.opts...)

			got, err := l.List()

			assertTicketSecret(t, got, test.want)
			assert.Equal(t, test.wantErr, err != nil)
		})
	}
}

func TestLister_WithFilterByGID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		fields  listerFields
		want    []expectedSecret
		wantErr bool
	}{
		{
			name: "one secret with ticket in default namespace",
			fields: listerFields{
				client: fake.NewSimpleClientset(
					secretFromTicketJSON(
						t,
						"default",
						"test-secret",
						[]byte(`{"ticket":{"userCreds":{"gids":[1000]}}}`),
					),
				),
				namespace: "default",
				opts: []ListerOption{
					WithFilterByGID(1000),
				},
			},
			want: []expectedSecret{
				newExpectedSecret("default", "test-secret"),
			},
			wantErr: false,
		},
		{
			name: "one secret with ticket in default namespace, filter by different gid",
			fields: listerFields{
				client: fake.NewSimpleClientset(
					secretFromTicketJSON(
						t,
						"default",
						"test-secret",
						[]byte(`{"ticket":{"userCreds":{"gids":[1000]}}}`),
					),
				),
				namespace: "default",
				opts: []ListerOption{
					WithFilterByGID(2000),
				},
			},
			want:    []expectedSecret{},
			wantErr: false,
		},
		{
			name: "two secrets with ticket in kube-system namespace, filter by one common gid",
			fields: listerFields{
				client: fake.NewSimpleClientset(
					secretFromTicketJSON(
						t,
						"kube-system",
						"test-secret-1",
						[]byte(`{"ticket":{"userCreds":{"gids":[1000,2000]}}}`),
					),
					secretFromTicketJSON(
						t,
						"kube-system",
						"test-secret-2",
						[]byte(`{"ticket":{"userCreds":{"gids":[2000,3000]}}}`),
					),
				),
				namespace: "kube-system",
				opts: []ListerOption{
					WithFilterByGID(2000),
				},
			},
			want: []expectedSecret{
				newExpectedSecret("kube-system", "test-secret-1"),
				newExpectedSecret("kube-system", "test-secret-2"),
			},
			wantErr: false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			l := NewLister(test.fields.client, test.fields.namespace, test.fields.opts...)

			got, err := l.List()

			assertTicketSecret(t, got, test.want)
			assert.Equal(t, test.wantErr, err != nil)
		})
	}
}

func TestLister_WithFilterOnlyExpired(t *testing.T) {
	ticketWithExpiryTime := func(t *testing.T, expiryTime time.Time) []byte {
		unix := uint64(expiryTime.Unix())
		return []byte(fmt.Sprintf(`{"ticket":{"expiryTime":%d}}`, unix))
	}

	t.Parallel()

	tests := []struct {
		name    string
		fields  listerFields
		want    []expectedSecret
		wantErr bool
	}{
		{
			name: "one unexpired secret with ticket in default namespace",
			fields: listerFields{
				client: fake.NewSimpleClientset(
					secretFromTicketJSON(
						t,
						"default",
						"test-secret-1",
						ticketWithExpiryTime(t, time.Now().Add(1*time.Hour)),
					),
				),
				namespace: "default",
				opts: []ListerOption{
					WithFilterOnlyExpired(),
				},
			},
			want:    []expectedSecret{},
			wantErr: false,
		},
		{
			name: "one expired secret with ticket in default namespace",
			fields: listerFields{
				client: fake.NewSimpleClientset(
					secretFromTicketJSON(
						t,
						"default",
						"expired-secret-1",
						ticketWithExpiryTime(t, time.Now().Add(-1*time.Hour)),
					),
				),
				namespace: "default",
				opts: []ListerOption{
					WithFilterOnlyExpired(),
				},
			},
			want: []expectedSecret{
				newExpectedSecret("default", "expired-secret-1"),
			},
			wantErr: false,
		},
		{
			name: "two secrets with ticket in kube-system namespace, one expired, one unexpired",
			fields: listerFields{
				client: fake.NewSimpleClientset(
					secretFromTicketJSON(
						t,
						"kube-system",
						"expired-secret-1",
						ticketWithExpiryTime(t, time.Now().Add(-1*time.Hour)),
					),
					secretFromTicketJSON(
						t,
						"kube-system",
						"test-secret-1",
						ticketWithExpiryTime(t, time.Now().Add(1*time.Hour)),
					),
				),
				namespace: "kube-system",
				opts: []ListerOption{
					WithFilterOnlyExpired(),
				},
			},
			want: []expectedSecret{
				newExpectedSecret("kube-system", "expired-secret-1"),
			},
			wantErr: false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			l := NewLister(test.fields.client, test.fields.namespace, test.fields.opts...)

			got, err := l.List()

			assertTicketSecret(t, got, test.want)
			assert.Equal(t, test.wantErr, err != nil)
		})
	}
}

func TestLister_WithFilterOnlyUnexpired(t *testing.T) {
	ticketWithExpiryTime := func(t *testing.T, expiryTime time.Time) []byte {
		unix := uint64(expiryTime.Unix())
		return []byte(fmt.Sprintf(`{"ticket":{"expiryTime":%d}}`, unix))
	}

	t.Parallel()

	tests := []struct {
		name    string
		fields  listerFields
		want    []expectedSecret
		wantErr bool
	}{
		{
			name: "one expired secret with ticket in default namespace",
			fields: listerFields{
				client: fake.NewSimpleClientset(
					secretFromTicketJSON(
						t,
						"default",
						"test-secret-1",
						ticketWithExpiryTime(t, time.Now().Add(-1*time.Hour)),
					),
				),
				namespace: "default",
				opts: []ListerOption{
					WithFilterOnlyUnexpired(),
				},
			},
			want:    []expectedSecret{},
			wantErr: false,
		},
		{
			name: "one unexpired secret with ticket in default namespace",
			fields: listerFields{
				client: fake.NewSimpleClientset(
					secretFromTicketJSON(
						t,
						"default",
						"expired-secret-1",
						ticketWithExpiryTime(t, time.Now().Add(1*time.Hour)),
					),
				),
				namespace: "default",
				opts: []ListerOption{
					WithFilterOnlyUnexpired(),
				},
			},
			want: []expectedSecret{
				newExpectedSecret("default", "expired-secret-1"),
			},
			wantErr: false,
		},
		{
			name: "two secrets with ticket in kube-system namespace, one expired, one unexpired",
			fields: listerFields{
				client: fake.NewSimpleClientset(
					secretFromTicketJSON(
						t,
						"kube-system",
						"expired-secret-1",
						ticketWithExpiryTime(t, time.Now().Add(1*time.Hour)),
					),
					secretFromTicketJSON(
						t,
						"kube-system",
						"test-secret-1",
						ticketWithExpiryTime(t, time.Now().Add(-1*time.Hour)),
					),
				),
				namespace: "kube-system",
				opts: []ListerOption{
					WithFilterOnlyUnexpired(),
				},
			},
			want: []expectedSecret{
				newExpectedSecret("kube-system", "expired-secret-1"),
			},
			wantErr: false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			l := NewLister(test.fields.client, test.fields.namespace, test.fields.opts...)

			got, err := l.List()

			assertTicketSecret(t, got, test.want)
			assert.Equal(t, test.wantErr, err != nil)
		})
	}
}

func TestLister_WithFilterExpiresBefore(t *testing.T) {
	ticketWithExpiryTime := func(t *testing.T, expiryTime time.Time) []byte {
		unix := uint64(expiryTime.Unix())
		return []byte(fmt.Sprintf(`{"ticket":{"expiryTime":%d}}`, unix))
	}

	t.Parallel()

	tests := []struct {
		name    string
		fields  listerFields
		want    []expectedSecret
		wantErr bool
	}{
		{
			name: "one secret with ticket in default namespace, expires in 24 hours, filter expires before 12 hours",
			fields: listerFields{
				client: fake.NewSimpleClientset(
					secretFromTicketJSON(
						t,
						"default",
						"test-secret-1",
						ticketWithExpiryTime(t, time.Now().Add(24*time.Hour)),
					),
				),
				namespace: "default",
				opts: []ListerOption{
					WithFilterExpiresBefore(12 * time.Hour),
				},
			},
			want:    []expectedSecret{},
			wantErr: false,
		},
		{
			name: "one secret with ticket in default namespace, expires in 12 hours, filter expires before 24 hours",
			fields: listerFields{
				client: fake.NewSimpleClientset(
					secretFromTicketJSON(
						t,
						"default",
						"test-secret-1",
						ticketWithExpiryTime(t, time.Now().Add(12*time.Hour)),
					),
				),
				namespace: "default",
				opts: []ListerOption{
					WithFilterExpiresBefore(24 * time.Hour),
				},
			},
			want: []expectedSecret{
				newExpectedSecret("default", "test-secret-1"),
			},
			wantErr: false,
		},
		{
			name: "two secrets with ticket in kube-system namespace, one expires in 24 hours, one expires in 6 hours, filter expires before 12 hours",
			fields: listerFields{
				client: fake.NewSimpleClientset(
					secretFromTicketJSON(
						t,
						"kube-system",
						"test-secret-1",
						ticketWithExpiryTime(t, time.Now().Add(6*time.Hour)),
					),
					secretFromTicketJSON(
						t,
						"kube-system",
						"test-secret-2",
						ticketWithExpiryTime(t, time.Now().Add(24*time.Hour)),
					),
				),
				namespace: "kube-system",
				opts: []ListerOption{
					WithFilterExpiresBefore(12 * time.Hour),
				},
			},
			want: []expectedSecret{
				newExpectedSecret("kube-system", "test-secret-1"),
			},
			wantErr: false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			l := NewLister(test.fields.client, test.fields.namespace, test.fields.opts...)

			got, err := l.List()

			assertTicketSecret(t, got, test.want)
			assert.Equal(t, test.wantErr, err != nil)
		})
	}
}

func TestLister_WithMultipleFilters(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		fields  listerFields
		want    []expectedSecret
		wantErr bool
	}{
		{
			name: "multiple secrets with ticket in kube-system namespace, filter by cluster and user",
			fields: listerFields{
				client: fake.NewSimpleClientset(
					secretFromTicketJSON(
						t,
						"kube-system",
						"test-secret-1",
						[]byte(`{"cluster":"test-cluster-1","ticket":{"userCreds":{"userName":"test-user-1"}}}`),
					),
					secretFromTicketJSON(
						t,
						"kube-system",
						"test-secret-2",
						[]byte(`{"cluster":"test-cluster-2","ticket":{"userCreds":{"userName":"test-user-2"}}}`),
					),
					secretFromTicketJSON(
						t,
						"kube-system",
						"test-secret-3",
						[]byte(`{"cluster":"test-cluster-1","ticket":{"userCreds":{"userName":"test-user-2"}}}`),
					),
				),
				namespace: "kube-system",
				opts: []ListerOption{
					WithFilterByMaprCluster("test-cluster-2"),
					WithFilterByMaprUser("test-user-2"),
				},
			},
			want: []expectedSecret{
				newExpectedSecret("kube-system", "test-secret-2"),
			},
			wantErr: false,
		},
		{
			name: "multiple secrets with ticket in kube-system namespace, filter by uid and gid",
			fields: listerFields{
				client: fake.NewSimpleClientset(
					secretFromTicketJSON(
						t,
						"kube-system",
						"test-secret-1",
						[]byte(`{"ticket":{"userCreds":{"uid":1000,"gids":[1000,2000]}}}`),
					),
					secretFromTicketJSON(
						t,
						"kube-system",
						"test-secret-2",
						[]byte(`{"ticket":{"userCreds":{"uid":2000,"gids":[1000,2000]}}}`),
					),
					secretFromTicketJSON(
						t,
						"kube-system",
						"test-secret-3",
						[]byte(`{"ticket":{"userCreds":{"uid":1000,"gids":[1000,3000]}}}`),
					),
				),
				namespace: "kube-system",
				opts: []ListerOption{
					WithFilterByUID(1000),
					WithFilterByGID(1000),
				},
			},
			want: []expectedSecret{
				newExpectedSecret("kube-system", "test-secret-1"),
				newExpectedSecret("kube-system", "test-secret-3"),
			},
			wantErr: false,
		},
		{
			name: "multiple secrets with ticket in kube-system namespace, filter by cluster, user, uid and gid",
			fields: listerFields{
				client: fake.NewSimpleClientset(
					secretFromTicketJSON(
						t,
						"kube-system",
						"test-secret-1",
						[]byte(`{"cluster":"test-cluster-1","ticket":{"userCreds":{"userName":"test-user-2","uid":1000,"gids":[1000,2000]}}}`),
					),
					secretFromTicketJSON(
						t,
						"kube-system",
						"test-secret-2",
						[]byte(`{"cluster":"test-cluster-2","ticket":{"userCreds":{"userName":"test-user-2","uid":2000,"gids":[1000,2000]}}}`),
					),
				),

				namespace: "kube-system",
				opts: []ListerOption{
					WithFilterByMaprCluster("test-cluster-1"),
					WithFilterByMaprUser("test-user-1"),
					WithFilterByUID(1000),
					WithFilterByGID(1000),
				},
			},
			want:    []expectedSecret{},
			wantErr: false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			l := NewLister(test.fields.client, test.fields.namespace, test.fields.opts...)

			got, err := l.List()

			assertTicketSecret(t, got, test.want)
			assert.Equal(t, test.wantErr, err != nil)
		})
	}
}

func TestLister_WithSortByName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		fields  listerFields
		want    []expectedSecret
		wantErr bool
	}{
		{
			name: "two secrets with ticket in kube-system namespace, sort by name",
			fields: listerFields{
				client: fake.NewSimpleClientset(
					secretFromTicketJSON(
						t,
						"kube-system",
						"test-secret-2",
						[]byte(`{"ticket":{"cluster":"test-cluster-2"}}`),
					),
					secretFromTicketJSON(
						t,
						"kube-system",
						"test-secret-1",
						[]byte(`{"ticket":{"cluster":"test-cluster-1"}}`),
					),
				),
				namespace: "kube-system",
				opts: []ListerOption{
					WithSortBy([]SortOption{SortByName}),
				},
			},
			want: []expectedSecret{
				newExpectedSecret("kube-system", "test-secret-1"),
				newExpectedSecret("kube-system", "test-secret-2"),
			},
			wantErr: false,
		},
		{
			name: "three secrets, different namespaces, sort all by name",
			fields: listerFields{
				client: fake.NewSimpleClientset(
					secretFromTicketJSON(
						t,
						"default",
						"test-secret-2",
						[]byte(`{"ticket":{"cluster":"test-cluster-2"}}`),
					),
					secretFromTicketJSON(
						t,
						"kube-system",
						"test-secret-1",
						[]byte(`{"ticket":{"cluster":"test-cluster-1"}}`),
					),
					secretFromTicketJSON(
						t,
						"kube-public",
						"test-secret-3",
						[]byte(`{"ticket":{"cluster":"test-cluster-3"}}`),
					),
				),
				namespace: metaV1.NamespaceAll,
				opts: []ListerOption{
					WithSortBy([]SortOption{SortByName}),
				},
			},
			want: []expectedSecret{
				newExpectedSecret("kube-system", "test-secret-1"),
				newExpectedSecret("default", "test-secret-2"),
				newExpectedSecret("kube-public", "test-secret-3"),
			},
			wantErr: false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			l := NewLister(test.fields.client, test.fields.namespace, test.fields.opts...)

			got, err := l.List()

			assertTicketSecret(t, got, test.want)
			assert.Equal(t, test.wantErr, err != nil)
		})
	}
}

func TestLister_WithSortByNamespace(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		fields  listerFields
		want    []expectedSecret
		wantErr bool
	}{
		{
			name: "two secrets with ticket in kube-system namespace, sort by namespace",
			fields: listerFields{
				client: fake.NewSimpleClientset(
					secretFromTicketJSON(
						t,
						"default",
						"test-secret-2",
						[]byte(`{"ticket":{"cluster":"test-cluster-2"}}`),
					),
					secretFromTicketJSON(
						t,
						"kube-system",
						"test-secret-1",
						[]byte(`{"ticket":{"cluster":"test-cluster-1"}}`),
					),
				),
				namespace: metaV1.NamespaceAll,
				opts: []ListerOption{
					WithSortBy([]SortOption{SortByNamespace}),
				},
			},
			want: []expectedSecret{
				newExpectedSecret("default", "test-secret-2"),
				newExpectedSecret("kube-system", "test-secret-1"),
			},
			wantErr: false,
		},
		{
			name: "three secrets, different namespaces, sort all by namespace",
			fields: listerFields{
				client: fake.NewSimpleClientset(
					secretFromTicketJSON(
						t,
						"default",
						"test-secret-2",
						[]byte(`{"ticket":{"cluster":"test-cluster-2"}}`),
					),
					secretFromTicketJSON(
						t,
						"kube-system",
						"test-secret-1",
						[]byte(`{"ticket":{"cluster":"test-cluster-1"}}`),
					),
					secretFromTicketJSON(
						t,
						"kube-public",
						"test-secret-3",
						[]byte(`{"ticket":{"cluster":"test-cluster-3"}}`),
					),
				),
				namespace: metaV1.NamespaceAll,
				opts: []ListerOption{
					WithSortBy([]SortOption{SortByNamespace}),
				},
			},
			want: []expectedSecret{
				newExpectedSecret("default", "test-secret-2"),
				newExpectedSecret("kube-public", "test-secret-3"),
				newExpectedSecret("kube-system", "test-secret-1"),
			},
			wantErr: false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			l := NewLister(test.fields.client, test.fields.namespace, test.fields.opts...)

			got, err := l.List()

			assertTicketSecret(t, got, test.want)
			assert.Equal(t, test.wantErr, err != nil)
		})
	}
}

func TestLister_WithSortByMaprCluster(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		fields  listerFields
		want    []expectedSecret
		wantErr bool
	}{
		{
			name: "two secrets with ticket in kube-system namespace, sort by mapr cluster",
			fields: listerFields{
				client: fake.NewSimpleClientset(
					secretFromTicketJSON(
						t,
						"default",
						"test-secret-2",
						[]byte(`{"cluster":"test-cluster-2"}`),
					),
					secretFromTicketJSON(
						t,
						"kube-system",
						"test-secret-1",
						[]byte(`{"cluster":"test-cluster-1"}`),
					),
				),
				namespace: metaV1.NamespaceAll,
				opts: []ListerOption{
					WithSortBy([]SortOption{SortByMaprCluster}),
				},
			},
			want: []expectedSecret{
				newExpectedSecret("kube-system", "test-secret-1"),
				newExpectedSecret("default", "test-secret-2"),
			},
			wantErr: false,
		},
		{
			name: "three secrets, different namespaces, sort all by mapr cluster",
			fields: listerFields{
				client: fake.NewSimpleClientset(
					secretFromTicketJSON(
						t,
						"default",
						"test-secret-2",
						[]byte(`{"cluster":"test-cluster-2"}`),
					),
					secretFromTicketJSON(
						t,
						"kube-system",
						"test-secret-1",
						[]byte(`{"cluster":"test-cluster-1"}`),
					),
					secretFromTicketJSON(
						t,
						"kube-public",
						"test-secret-3",
						[]byte(`{"cluster":"test-cluster-3"}`),
					),
				),
				namespace: metaV1.NamespaceAll,
				opts: []ListerOption{
					WithSortBy([]SortOption{SortByMaprCluster}),
				},
			},
			want: []expectedSecret{
				newExpectedSecret("kube-system", "test-secret-1"),
				newExpectedSecret("default", "test-secret-2"),
				newExpectedSecret("kube-public", "test-secret-3"),
			},
			wantErr: false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			l := NewLister(test.fields.client, test.fields.namespace, test.fields.opts...)

			got, err := l.List()

			assertTicketSecret(t, got, test.want)
			assert.Equal(t, test.wantErr, err != nil)
		})
	}
}

func TestLister_WithSortByMaprUser(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		fields  listerFields
		want    []expectedSecret
		wantErr bool
	}{
		{
			name: "two secrets with ticket in kube-system namespace, sort by mapr user",
			fields: listerFields{
				client: fake.NewSimpleClientset(
					secretFromTicketJSON(
						t,
						"default",
						"test-secret-2",
						[]byte(`{"ticket":{"userCreds":{"userName":"test-user-2"}}}`),
					),
					secretFromTicketJSON(
						t,
						"kube-system",
						"test-secret-1",
						[]byte(`{"ticket":{"userCreds":{"userName":"test-user-1"}}}`),
					),
				),
				namespace: metaV1.NamespaceAll,
				opts: []ListerOption{
					WithSortBy([]SortOption{SortByMaprUser}),
				},
			},
			want: []expectedSecret{
				newExpectedSecret("kube-system", "test-secret-1"),
				newExpectedSecret("default", "test-secret-2"),
			},
			wantErr: false,
		},
		{
			name: "three secrets, different namespaces, sort all by mapr user",
			fields: listerFields{
				client: fake.NewSimpleClientset(
					secretFromTicketJSON(
						t,
						"default",
						"test-secret-2",
						[]byte(`{"ticket":{"userCreds":{"userName":"test-user-2"}}}`),
					),
					secretFromTicketJSON(
						t,
						"kube-system",
						"test-secret-1",
						[]byte(`{"ticket":{"userCreds":{"userName":"test-user-1"}}}`),
					),
					secretFromTicketJSON(
						t,
						"kube-public",
						"test-secret-3",
						[]byte(`{"ticket":{"userCreds":{"userName":"test-user-3"}}}`),
					),
				),
				namespace: metaV1.NamespaceAll,
				opts: []ListerOption{
					WithSortBy([]SortOption{SortByMaprUser}),
				},
			},
			want: []expectedSecret{
				newExpectedSecret("kube-system", "test-secret-1"),
				newExpectedSecret("default", "test-secret-2"),
				newExpectedSecret("kube-public", "test-secret-3"),
			},
			wantErr: false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			l := NewLister(test.fields.client, test.fields.namespace, test.fields.opts...)

			got, err := l.List()

			assertTicketSecret(t, got, test.want)
			assert.Equal(t, test.wantErr, err != nil)
		})
	}
}

func TestLister_WithSortByCreationTime(t *testing.T) {
	ticketWithCreationTime := func(t *testing.T, creationTime time.Time) []byte {
		unix := uint64(creationTime.Unix())
		return []byte(fmt.Sprintf(`{"ticket":{"creationTimeSec":%d}}`, unix))
	}

	t.Parallel()

	tests := []struct {
		name    string
		fields  listerFields
		want    []expectedSecret
		wantErr bool
	}{

		{
			name: "two secrets with ticket in kube-system namespace, sort by creation time",
			fields: listerFields{
				client: fake.NewSimpleClientset(
					secretFromTicketJSON(
						t,
						"kube-system",
						"test-secret-2",
						ticketWithCreationTime(t, time.Now().Add(-1*time.Hour)),
					),
					secretFromTicketJSON(
						t,
						"kube-system",
						"test-secret-1",
						ticketWithCreationTime(t, time.Now().Add(-2*time.Hour)),
					),
				),
				namespace: metaV1.NamespaceAll,
				opts: []ListerOption{
					WithSortBy([]SortOption{SortByAge}),
				},
			},
			want: []expectedSecret{
				newExpectedSecret("kube-system", "test-secret-1"),
				newExpectedSecret("kube-system", "test-secret-2"),
			},
			wantErr: false,
		},
		{
			name: "three secrets, different namespaces, sort all by creation time",
			fields: listerFields{
				client: fake.NewSimpleClientset(
					secretFromTicketJSON(
						t,
						"default",
						"test-secret-2",
						ticketWithCreationTime(t, time.Now().Add(-1*time.Hour)),
					),
					secretFromTicketJSON(
						t,
						"kube-system",
						"test-secret-1",
						ticketWithCreationTime(t, time.Now().Add(-2*time.Hour)),
					),
					secretFromTicketJSON(
						t,
						"kube-public",
						"test-secret-3",
						ticketWithCreationTime(t, time.Now().Add(-3*time.Hour)),
					),
				),
				namespace: metaV1.NamespaceAll,
				opts: []ListerOption{
					WithSortBy([]SortOption{SortByAge}),
				},
			},
			want: []expectedSecret{
				newExpectedSecret("kube-public", "test-secret-3"),
				newExpectedSecret("kube-system", "test-secret-1"),
				newExpectedSecret("default", "test-secret-2"),
			},
			wantErr: false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			l := NewLister(test.fields.client, test.fields.namespace, test.fields.opts...)

			got, err := l.List()

			assertTicketSecret(t, got, test.want)
			assert.Equal(t, test.wantErr, err != nil)
		})
	}
}

func TestLister_WithSortByExpiryTime(t *testing.T) {
	ticketWithExpiryTime := func(t *testing.T, expiryTime time.Time) []byte {
		unix := uint64(expiryTime.Unix())
		return []byte(fmt.Sprintf(`{"ticket":{"expiryTime":%d}}`, unix))
	}

	t.Parallel()

	tests := []struct {
		name    string
		fields  listerFields
		want    []expectedSecret
		wantErr bool
	}{
		{
			name: "two secrets with ticket in kube-system namespace, sort by expiry time",
			fields: listerFields{
				client: fake.NewSimpleClientset(
					secretFromTicketJSON(
						t,
						"kube-system",
						"test-secret-2",
						ticketWithExpiryTime(t, time.Now().Add(1*time.Hour)),
					),
					secretFromTicketJSON(
						t,
						"kube-system",
						"test-secret-1",
						ticketWithExpiryTime(t, time.Now().Add(2*time.Hour)),
					),
				),
				namespace: metaV1.NamespaceAll,
				opts: []ListerOption{
					WithSortBy([]SortOption{SortByExpiration}),
				},
			},
			want: []expectedSecret{
				newExpectedSecret("kube-system", "test-secret-2"),
				newExpectedSecret("kube-system", "test-secret-1"),
			},
			wantErr: false,
		},
		{
			name: "three secrets, different namespaces, sort all by expiry time",
			fields: listerFields{
				client: fake.NewSimpleClientset(
					secretFromTicketJSON(
						t,
						"default",
						"test-secret-2",
						ticketWithExpiryTime(t, time.Now().Add(1*time.Hour)),
					),
					secretFromTicketJSON(
						t,
						"kube-system",
						"test-secret-1",
						ticketWithExpiryTime(t, time.Now().Add(-2*time.Hour)),
					),
					secretFromTicketJSON(
						t,
						"kube-public",
						"test-secret-3",
						ticketWithExpiryTime(t, time.Now().Add(3*time.Hour)),
					),
				),
				namespace: metaV1.NamespaceAll,
				opts: []ListerOption{
					WithSortBy([]SortOption{SortByExpiration}),
				},
			},
			want: []expectedSecret{
				newExpectedSecret("kube-system", "test-secret-1"),
				newExpectedSecret("default", "test-secret-2"),
				newExpectedSecret("kube-public", "test-secret-3"),
			},
			wantErr: false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			l := NewLister(test.fields.client, test.fields.namespace, test.fields.opts...)

			got, err := l.List()

			assertTicketSecret(t, got, test.want)
			assert.Equal(t, test.wantErr, err != nil)
		})
	}
}

type listerFields struct {
	client    kubernetes.Interface
	namespace string
	opts      []ListerOption
}

type expectedSecret struct {
	name      string
	namespace string
}

func newExpectedSecret(namespace, name string) expectedSecret {
	return expectedSecret{
		name:      name,
		namespace: namespace,
	}
}

func assertTicketSecret(t *testing.T, secrets []types.MaprSecret, expected []expectedSecret) {
	t.Helper()

	assert.Len(t, secrets, len(expected))

	for i := range secrets {
		assert.Equal(t, expected[i].name, secrets[i].Secret.Name)
		assert.Equal(t, expected[i].namespace, secrets[i].Secret.Namespace)
	}
}

func secretFromTicketJSON(t *testing.T, namespace, name string, in []byte) *coreV1.Secret {
	t.Helper()

	obj := ticket.NewMaprTicket()
	err := json.Unmarshal(in, &obj)
	if err != nil {
		t.Fatal(err)
	}

	out, err := parse.Marshal(obj.AsMaprTicket())
	if err != nil {
		t.Fatal(err)
	}

	return &coreV1.Secret{
		ObjectMeta: metaV1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Data: map[string][]byte{
			ticket.SecretMaprTicketKey: out,
		},
	}
}
