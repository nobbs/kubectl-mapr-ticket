package secret_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	. "github.com/nobbs/kubectl-mapr-ticket/internal/secret"
	"github.com/nobbs/kubectl-mapr-ticket/internal/ticket"
	"github.com/nobbs/mapr-ticket-parser/pkg/parse"

	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
)

func TestLister_List(t *testing.T) {
	tests := []struct {
		name    string
		fields  fields
		want    []testSecret
		wantErr bool
	}{
		{
			name: "no secrets in default namespace",
			fields: fields{
				client:    fake.NewSimpleClientset(),
				namespace: "default",
			},
			want:    []testSecret{},
			wantErr: false,
		},
		{
			name: "one secret without ticket in default namespace",
			fields: fields{
				client: fake.NewSimpleClientset(&coreV1.Secret{
					ObjectMeta: metaV1.ObjectMeta{
						Name:      "test-secret",
						Namespace: "default",
					},
				}),
				namespace: "default",
			},
			want:    []testSecret{},
			wantErr: false,
		},
		{
			name: "secret with invalid ticket data",
			fields: fields{
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
			want:    []testSecret{},
			wantErr: false,
		},
		{
			name: "one secret with ticket in default namespace",
			fields: fields{
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
			want: []testSecret{
				newTestSecret("default", "test-secret"),
			},
			wantErr: false,
		},
		{
			name: "two secrets with ticket in kube-system namespace",
			fields: fields{
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
			want: []testSecret{
				newTestSecret("kube-system", "test-secret-1"),
				newTestSecret("kube-system", "test-secret-2"),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLister(tt.fields.client, tt.fields.namespace, tt.fields.opts...)

			got, err := l.List()

			checkTicketSecret(t, got, tt.want)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestLister_ListWithFilterByMaprCluster(t *testing.T) {
	tests := []struct {
		name    string
		fields  fields
		want    []testSecret
		wantErr bool
	}{
		{
			name: "one secret with ticket in default namespace",
			fields: fields{
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
			want: []testSecret{
				newTestSecret("default", "test-secret"),
			},
			wantErr: false,
		},
		{
			name: "one secret with ticket in default namespace, filter by different cluster",
			fields: fields{
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
			want:    []testSecret{},
			wantErr: false,
		},
		{
			name: "two secrets with ticket in kube-system namespace, filter by one cluster",
			fields: fields{
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
			want: []testSecret{
				newTestSecret("kube-system", "test-secret-1"),
			},
			wantErr: false,
		},
		{
			name: "two secrets with ticket in kube-system namespace, filter by different cluster",
			fields: fields{
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
			want:    []testSecret{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLister(tt.fields.client, tt.fields.namespace, tt.fields.opts...)

			got, err := l.List()

			checkTicketSecret(t, got, tt.want)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestLister_ListWithFilterByMaprUser(t *testing.T) {
	tests := []struct {
		name    string
		fields  fields
		want    []testSecret
		wantErr bool
	}{
		{
			name: "one secret with ticket in default namespace",
			fields: fields{
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
			want: []testSecret{
				newTestSecret("default", "test-secret"),
			},
			wantErr: false,
		},
		{
			name: "one secret with ticket in default namespace, filter by different user",
			fields: fields{
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
			want:    []testSecret{},
			wantErr: false,
		},
		{
			name: "two secrets with ticket in kube-system namespace, filter by one user",
			fields: fields{
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
			want: []testSecret{
				newTestSecret("kube-system", "test-secret-1"),
			},
			wantErr: false,
		},
		{
			name: "two secrets with ticket in kube-system namespace, filter by different user",
			fields: fields{
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
			want:    []testSecret{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLister(tt.fields.client, tt.fields.namespace, tt.fields.opts...)

			got, err := l.List()

			checkTicketSecret(t, got, tt.want)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestLister_ListWithFilterByUID(t *testing.T) {
	tests := []struct {
		name    string
		fields  fields
		want    []testSecret
		wantErr bool
	}{
		{
			name: "one secret with ticket in default namespace",
			fields: fields{
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
			want: []testSecret{
				newTestSecret("default", "test-secret"),
			},
			wantErr: false,
		},
		{
			name: "one secret with ticket in default namespace, filter by different uid",
			fields: fields{
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
			want:    []testSecret{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLister(tt.fields.client, tt.fields.namespace, tt.fields.opts...)

			got, err := l.List()

			checkTicketSecret(t, got, tt.want)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestLister_ListWithFilterByGID(t *testing.T) {
	tests := []struct {
		name    string
		fields  fields
		want    []testSecret
		wantErr bool
	}{
		{
			name: "one secret with ticket in default namespace",
			fields: fields{
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
			want: []testSecret{
				newTestSecret("default", "test-secret"),
			},
			wantErr: false,
		},
		{
			name: "one secret with ticket in default namespace, filter by different gid",
			fields: fields{
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
			want:    []testSecret{},
			wantErr: false,
		},
		{
			name: "two secrets with ticket in kube-system namespace, filter by one common gid",
			fields: fields{
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
			want: []testSecret{
				newTestSecret("kube-system", "test-secret-1"),
				newTestSecret("kube-system", "test-secret-2"),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLister(tt.fields.client, tt.fields.namespace, tt.fields.opts...)

			got, err := l.List()

			checkTicketSecret(t, got, tt.want)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestLister_ListWithFilterOnlyExpired(t *testing.T) {
	ticketWithExpiryTime := func(t *testing.T, expiryTime time.Time) []byte {
		var unix uint64 = uint64(expiryTime.Unix())
		return []byte(fmt.Sprintf(`{"ticket":{"expiryTime":%d}}`, unix))
	}

	tests := []struct {
		name    string
		fields  fields
		want    []testSecret
		wantErr bool
	}{
		{
			name: "one unexpired secret with ticket in default namespace",
			fields: fields{
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
			want:    []testSecret{},
			wantErr: false,
		},
		{
			name: "one expired secret with ticket in default namespace",
			fields: fields{
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
			want: []testSecret{
				newTestSecret("default", "expired-secret-1"),
			},
			wantErr: false,
		},
		{
			name: "two secrets with ticket in kube-system namespace, one expired, one unexpired",
			fields: fields{
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
			want: []testSecret{
				newTestSecret("kube-system", "expired-secret-1"),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLister(tt.fields.client, tt.fields.namespace, tt.fields.opts...)

			got, err := l.List()

			checkTicketSecret(t, got, tt.want)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestLister_ListWithFilterOnlyUnexpired(t *testing.T) {
	ticketWithExpiryTime := func(t *testing.T, expiryTime time.Time) []byte {
		var unix uint64 = uint64(expiryTime.Unix())
		return []byte(fmt.Sprintf(`{"ticket":{"expiryTime":%d}}`, unix))
	}

	tests := []struct {
		name    string
		fields  fields
		want    []testSecret
		wantErr bool
	}{
		{
			name: "one expired secret with ticket in default namespace",
			fields: fields{
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
			want:    []testSecret{},
			wantErr: false,
		},
		{
			name: "one unexpired secret with ticket in default namespace",
			fields: fields{
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
			want: []testSecret{
				newTestSecret("default", "expired-secret-1"),
			},
			wantErr: false,
		},
		{
			name: "two secrets with ticket in kube-system namespace, one expired, one unexpired",
			fields: fields{
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
			want: []testSecret{
				newTestSecret("kube-system", "expired-secret-1"),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLister(tt.fields.client, tt.fields.namespace, tt.fields.opts...)

			got, err := l.List()

			checkTicketSecret(t, got, tt.want)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestLister_ListWithFilterExpiresBefore(t *testing.T) {
	ticketWithExpiryTime := func(t *testing.T, expiryTime time.Time) []byte {
		var unix uint64 = uint64(expiryTime.Unix())
		return []byte(fmt.Sprintf(`{"ticket":{"expiryTime":%d}}`, unix))
	}

	tests := []struct {
		name    string
		fields  fields
		want    []testSecret
		wantErr bool
	}{
		{
			name: "one secret with ticket in default namespace, expires in 24 hours, filter expires before 12 hours",
			fields: fields{
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
			want:    []testSecret{},
			wantErr: false,
		},
		{
			name: "one secret with ticket in default namespace, expires in 12 hours, filter expires before 24 hours",
			fields: fields{
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
			want: []testSecret{
				newTestSecret("default", "test-secret-1"),
			},
			wantErr: false,
		},
		{
			name: "two secrets with ticket in kube-system namespace, one expires in 24 hours, one expires in 6 hours, filter expires before 12 hours",
			fields: fields{
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
			want: []testSecret{
				newTestSecret("kube-system", "test-secret-1"),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLister(tt.fields.client, tt.fields.namespace, tt.fields.opts...)

			got, err := l.List()

			checkTicketSecret(t, got, tt.want)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestList_WithMultipleFilters(t *testing.T) {
	tests := []struct {
		name    string
		fields  fields
		want    []testSecret
		wantErr bool
	}{
		{
			name: "multiple secrets with ticket in kube-system namespace, filter by cluster and user",
			fields: fields{
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
			want: []testSecret{
				newTestSecret("kube-system", "test-secret-2"),
			},
			wantErr: false,
		},
		{
			name: "multiple secrets with ticket in kube-system namespace, filter by uid and gid",
			fields: fields{
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
			want: []testSecret{
				newTestSecret("kube-system", "test-secret-1"),
				newTestSecret("kube-system", "test-secret-3"),
			},
			wantErr: false,
		},
		{
			name: "multiple secrets with ticket in kube-system namespace, filter by cluster, user, uid and gid",
			fields: fields{
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
			want:    []testSecret{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLister(tt.fields.client, tt.fields.namespace, tt.fields.opts...)

			got, err := l.List()

			checkTicketSecret(t, got, tt.want)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestList_WithSortByName(t *testing.T) {
	tests := []struct {
		name    string
		fields  fields
		want    []testSecret
		wantErr bool
	}{
		{
			name: "two secrets with ticket in kube-system namespace, sort by name",
			fields: fields{
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
					WithSortBy([]SortOptions{SortByName}),
				},
			},
			want: []testSecret{
				newTestSecret("kube-system", "test-secret-1"),
				newTestSecret("kube-system", "test-secret-2"),
			},
			wantErr: false,
		},
		{
			name: "three secrets, different namespaces, sort all by name",
			fields: fields{
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
					WithSortBy([]SortOptions{SortByName}),
				},
			},
			want: []testSecret{
				newTestSecret("kube-system", "test-secret-1"),
				newTestSecret("default", "test-secret-2"),
				newTestSecret("kube-public", "test-secret-3"),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLister(tt.fields.client, tt.fields.namespace, tt.fields.opts...)

			got, err := l.List()

			checkTicketSecret(t, got, tt.want)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestList_WithSortByNamespace(t *testing.T) {
	tests := []struct {
		name    string
		fields  fields
		want    []testSecret
		wantErr bool
	}{
		{
			name: "two secrets with ticket in kube-system namespace, sort by namespace",
			fields: fields{
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
					WithSortBy([]SortOptions{SortByNamespace}),
				},
			},
			want: []testSecret{
				newTestSecret("default", "test-secret-2"),
				newTestSecret("kube-system", "test-secret-1"),
			},
			wantErr: false,
		},
		{
			name: "three secrets, different namespaces, sort all by namespace",
			fields: fields{
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
					WithSortBy([]SortOptions{SortByNamespace}),
				},
			},
			want: []testSecret{
				newTestSecret("default", "test-secret-2"),
				newTestSecret("kube-public", "test-secret-3"),
				newTestSecret("kube-system", "test-secret-1"),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLister(tt.fields.client, tt.fields.namespace, tt.fields.opts...)

			got, err := l.List()

			checkTicketSecret(t, got, tt.want)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestList_WithSortByMaprCluster(t *testing.T) {
	tests := []struct {
		name    string
		fields  fields
		want    []testSecret
		wantErr bool
	}{
		{
			name: "two secrets with ticket in kube-system namespace, sort by mapr cluster",
			fields: fields{
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
					WithSortBy([]SortOptions{SortByMaprCluster}),
				},
			},
			want: []testSecret{
				newTestSecret("kube-system", "test-secret-1"),
				newTestSecret("default", "test-secret-2"),
			},
			wantErr: false,
		},
		{
			name: "three secrets, different namespaces, sort all by mapr cluster",
			fields: fields{
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
					WithSortBy([]SortOptions{SortByMaprCluster}),
				},
			},
			want: []testSecret{
				newTestSecret("kube-system", "test-secret-1"),
				newTestSecret("default", "test-secret-2"),
				newTestSecret("kube-public", "test-secret-3"),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLister(tt.fields.client, tt.fields.namespace, tt.fields.opts...)

			got, err := l.List()

			checkTicketSecret(t, got, tt.want)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestList_WithSortByMaprUser(t *testing.T) {
	tests := []struct {
		name    string
		fields  fields
		want    []testSecret
		wantErr bool
	}{
		{
			name: "two secrets with ticket in kube-system namespace, sort by mapr user",
			fields: fields{
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
					WithSortBy([]SortOptions{SortByMaprUser}),
				},
			},
			want: []testSecret{
				newTestSecret("kube-system", "test-secret-1"),
				newTestSecret("default", "test-secret-2"),
			},
			wantErr: false,
		},
		{
			name: "three secrets, different namespaces, sort all by mapr user",
			fields: fields{
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
					WithSortBy([]SortOptions{SortByMaprUser}),
				},
			},
			want: []testSecret{
				newTestSecret("kube-system", "test-secret-1"),
				newTestSecret("default", "test-secret-2"),
				newTestSecret("kube-public", "test-secret-3"),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLister(tt.fields.client, tt.fields.namespace, tt.fields.opts...)

			got, err := l.List()

			checkTicketSecret(t, got, tt.want)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestList_WithSortByCreationTime(t *testing.T) {
	ticketWithCreationTime := func(t *testing.T, creationTime time.Time) []byte {
		var unix uint64 = uint64(creationTime.Unix())
		return []byte(fmt.Sprintf(`{"ticket":{"creationTimeSec":%d}}`, unix))
	}

	tests := []struct {
		name    string
		fields  fields
		want    []testSecret
		wantErr bool
	}{

		{
			name: "two secrets with ticket in kube-system namespace, sort by creation time",
			fields: fields{
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
					WithSortBy([]SortOptions{SortByCreationTimestamp}),
				},
			},
			want: []testSecret{
				newTestSecret("kube-system", "test-secret-1"),
				newTestSecret("kube-system", "test-secret-2"),
			},
			wantErr: false,
		},
		{
			name: "three secrets, different namespaces, sort all by creation time",
			fields: fields{
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
					WithSortBy([]SortOptions{SortByCreationTimestamp}),
				},
			},
			want: []testSecret{
				newTestSecret("kube-public", "test-secret-3"),
				newTestSecret("kube-system", "test-secret-1"),
				newTestSecret("default", "test-secret-2"),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLister(tt.fields.client, tt.fields.namespace, tt.fields.opts...)

			got, err := l.List()

			checkTicketSecret(t, got, tt.want)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestList_WithSortByExpiryTime(t *testing.T) {
	ticketWithExpiryTime := func(t *testing.T, expiryTime time.Time) []byte {
		var unix uint64 = uint64(expiryTime.Unix())
		return []byte(fmt.Sprintf(`{"ticket":{"expiryTime":%d}}`, unix))
	}

	tests := []struct {
		name    string
		fields  fields
		want    []testSecret
		wantErr bool
	}{
		{
			name: "two secrets with ticket in kube-system namespace, sort by expiry time",
			fields: fields{
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
					WithSortBy([]SortOptions{SortByExpiryTime}),
				},
			},
			want: []testSecret{
				newTestSecret("kube-system", "test-secret-2"),
				newTestSecret("kube-system", "test-secret-1"),
			},
			wantErr: false,
		},
		{
			name: "three secrets, different namespaces, sort all by expiry time",
			fields: fields{
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
					WithSortBy([]SortOptions{SortByExpiryTime}),
				},
			},
			want: []testSecret{
				newTestSecret("kube-system", "test-secret-1"),
				newTestSecret("default", "test-secret-2"),
				newTestSecret("kube-public", "test-secret-3"),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLister(tt.fields.client, tt.fields.namespace, tt.fields.opts...)

			got, err := l.List()

			checkTicketSecret(t, got, tt.want)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

type fields struct {
	client    kubernetes.Interface
	namespace string
	opts      []ListerOption
}

type testSecret struct {
	name      string
	namespace string
}

func newTestSecret(namespace, name string) testSecret {
	return testSecret{
		name:      name,
		namespace: namespace,
	}
}

func checkTicketSecret(t *testing.T, secrets []TicketSecret, tickets []testSecret) {
	t.Helper()

	assert.Equal(t, len(tickets), len(secrets))

	for i := range secrets {
		assert.Equal(t, tickets[i].name, secrets[i].Secret.Name)
		assert.Equal(t, tickets[i].namespace, secrets[i].Secret.Namespace)
	}
}

func secretFromTicketJSON(t *testing.T, namespace, name string, in []byte) *coreV1.Secret {
	t.Helper()

	v := ticket.NewMaprTicket()
	err := json.Unmarshal(in, &v)
	if err != nil {
		t.Fatal(err)
	}

	out := marshalTicket(v)

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

func marshalTicket(t *ticket.MaprTicket) []byte {
	b, _ := parse.Marshal((*parse.MaprTicket)(t))
	return b
}
