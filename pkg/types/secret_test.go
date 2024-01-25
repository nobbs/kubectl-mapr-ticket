package types_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/nobbs/kubectl-mapr-ticket/pkg/ticket"
	. "github.com/nobbs/kubectl-mapr-ticket/pkg/types"

	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	testTicketsRaw = [][]byte{
		[]byte("demo.mapr.com +Cze+qwYCbAXGbz56OO7UF+lGqL3WPXrNkO1SLawEEDmSbgNl019xBeBY3kvh+R13iz/mCnwpzsLQw4Y5jEnv5GtuIWbeoC95ha8VKwX8MKcE6Kn9nZ2AF0QminkHwNVBx6TDriGZffyJCfZzivBwBSdKoQEWhBOPFCIMAi7w2zV/SX5Ut7u4qIKvEpr0JHV7sLMWYLhYncM6CKMd7iECGvECsBvEZRVj+dpbEY0BaRN/W54/7wNWaSVELUF6JWHQ8dmsqty4cZlI0/MV10HZzIbl9sMLFQ="),
		[]byte("demo.mapr.com cj1FDarNNKh7f+hL5ho1m32RzYyHPKuGIPJzE/CkUqEfcTGEP4YJuFlTsBmHuifI5LvNob/Y4xmDsrz9OxrBnhly/0g9xAs5ApZWNY8Rcab8q70IBYIbpu7xsBBTAiVRyLJkAtGFXNn104BB0AsS55GbQFUN9NAiWLzZY3/X1ITfGfDEGaYbWWTb1LGx6C0Jjgnr7TzXv1GqwiASbcUQCXOx4inguwMneYt9KhOp89smw6GBKP064DfIMHHR6lgv0XhBP6d9FVJ1QWKvcccvi2F3LReBtqA="),
		[]byte("demo.mapr.com IGem6fUksZ1pd4iut978SKElS4ktecRsAkrl+qwPYc7xhfMg4wkwALKDmFmpc8Xvrm1L9Et0jVBoyhCWMDCjhToZ8b6FsfCn8wdCOB0MWm9CRobGv7MDsoEO2TQ5Bnh8i/VfuthKFxd3Om9iZPVCI4I1S9h4p/77Al1GzTGcfFFf1g9fq1HXftT9TEDyLdABIyATJbzv8zD10IDT8P1f8nxl7lgT/7ZhGz7N24vSz6jBxHE7oHmvHzjW22xJwt7TJgvrP21boH9HTsTPiKZOpQMZ4zFo6JA4aNVlQQ0="),
	}
)

func TestNewMaprSecret(t *testing.T) {
	tests := []struct {
		name string
		s    *Secret
		want *MaprSecret
	}{
		{
			name: "nil",
			s:    nil,
			want: &MaprSecret{},
		},
		{
			name: "empty",
			s:    &Secret{},
			want: &MaprSecret{
				Secret: &Secret{},
			},
		},
		{
			name: "ticket",
			s: &Secret{
				Data: map[string][]byte{
					ticket.SecretMaprTicketKey: testTicketsRaw[0],
				},
			},
			want: &MaprSecret{
				Ticket: &ticket.Ticket{
					Cluster: "demo.mapr.com",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewMaprSecret(tt.s)

			if tt.want.Secret != nil {
				assert.Equal(t, tt.want.Secret, got.Secret)
			}

			if tt.want.Ticket != nil {
				assert.NotNil(t, got.Ticket)
			}
		})
	}
}

func TestMaprSecret_GetSecretName(t *testing.T) {
	tests := []struct {
		name string
		t    *MaprSecret
		want string
	}{
		{
			name: "nil",
			t:    nil,
			want: "",
		},
		{
			name: "empty",
			t:    &MaprSecret{},
			want: "",
		},
		{
			name: "name",
			t: &MaprSecret{
				Secret: &Secret{
					ObjectMeta: metaV1.ObjectMeta{
						Name: "test",
					},
				},
			},
			want: "test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.t.GetSecretName()

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestMaprSecret_GetSecretNamespace(t *testing.T) {
	tests := []struct {
		name string
		t    *MaprSecret
		want string
	}{
		{
			name: "nil",
			t:    nil,
			want: "",
		},
		{
			name: "empty",
			t:    &MaprSecret{},
			want: "",
		},
		{
			name: "namespace",
			t: &MaprSecret{
				Secret: &Secret{
					ObjectMeta: metaV1.ObjectMeta{
						Namespace: "test",
					},
				},
			},
			want: "test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.t.GetSecretNamespace()

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestMaprSecret_GetCluster(t *testing.T) {
	tests := []struct {
		name string
		t    *MaprSecret
		want string
	}{
		{
			name: "nil",
			t:    nil,
			want: "",
		},
		{
			name: "empty",
			t:    &MaprSecret{},
			want: "",
		},
		{
			name: "cluster",
			t: &MaprSecret{
				Ticket: &ticket.Ticket{
					Cluster: "test",
				},
			},
			want: "test",
		},
		{
			name: "demo.mapr.com",
			t: NewMaprSecret(
				&Secret{
					Data: map[string][]byte{
						ticket.SecretMaprTicketKey: testTicketsRaw[0],
					},
				},
			),
			want: "demo.mapr.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.t.GetCluster()

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestMaprSecret_GetUser(t *testing.T) {
	tests := []struct {
		name string
		t    *MaprSecret
		want string
	}{
		{
			name: "nil",
			t:    nil,
			want: "",
		},
		{
			name: "empty",
			t:    &MaprSecret{},
			want: "",
		},
		{
			name: "demo.mapr.com",
			t: NewMaprSecret(
				&Secret{
					Data: map[string][]byte{
						ticket.SecretMaprTicketKey: testTicketsRaw[0],
					},
				},
			),
			want: "mapr",
		},
		{
			name: "demo.mapr.com",
			t: NewMaprSecret(
				&Secret{
					Data: map[string][]byte{
						ticket.SecretMaprTicketKey: testTicketsRaw[1],
					},
				},
			),
			want: "mapr",
		},
		{
			name: "demo.mapr.com",
			t: NewMaprSecret(
				&Secret{
					Data: map[string][]byte{
						ticket.SecretMaprTicketKey: testTicketsRaw[2],
					},
				},
			),
			want: "mapr",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.t.GetUser()

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestMaprSecret_GetExpirationTime(t *testing.T) {
	tests := []struct {
		name string
		t    *MaprSecret
		want time.Time
	}{
		{
			name: "nil",
			t:    nil,
			want: time.Time{},
		},
		{
			name: "empty",
			t:    &MaprSecret{},
			want: time.Time{},
		},
		{
			name: "demo.mapr.com",
			t: NewMaprSecret(
				&Secret{
					Data: map[string][]byte{
						ticket.SecretMaprTicketKey: testTicketsRaw[0],
					},
				},
			),
			want: time.Date(29229672, time.June, 17, 19, 31, 17, 0, time.Local),
		},
		{
			name: "demo.mapr.com",
			t: NewMaprSecret(
				&Secret{
					Data: map[string][]byte{
						ticket.SecretMaprTicketKey: testTicketsRaw[1],
					},
				},
			),
			want: time.Date(2019, time.February, 19, 13, 13, 49, 0, time.Local),
		},
		{
			name: "demo.mapr.com",
			t: NewMaprSecret(
				&Secret{
					Data: map[string][]byte{
						ticket.SecretMaprTicketKey: testTicketsRaw[2],
					},
				},
			),
			want: time.Date(2021, time.April, 30, 0, 32, 46, 0, time.Local),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.t.GetExpirationTime()

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestMaprSecret_GetCreationTime(t *testing.T) {
	tests := []struct {
		name string
		t    *MaprSecret
		want time.Time
	}{
		{
			name: "nil",
			t:    nil,
			want: time.Time{},
		},
		{
			name: "empty",
			t:    &MaprSecret{},
			want: time.Time{},
		},
		{
			name: "demo.mapr.com",
			t: NewMaprSecret(
				&Secret{
					Data: map[string][]byte{
						ticket.SecretMaprTicketKey: testTicketsRaw[0],
					},
				},
			),
			want: time.Date(2018, time.April, 4, 16, 31, 37, 0, time.Local),
		},
		{
			name: "demo.mapr.com",
			t: NewMaprSecret(
				&Secret{
					Data: map[string][]byte{
						ticket.SecretMaprTicketKey: testTicketsRaw[1],
					},
				},
			),
			want: time.Date(2019, time.February, 5, 13, 13, 49, 0, time.Local),
		},
		{
			name: "demo.mapr.com",
			t: NewMaprSecret(
				&Secret{
					Data: map[string][]byte{
						ticket.SecretMaprTicketKey: testTicketsRaw[2],
					},
				},
			),
			want: time.Date(2021, time.April, 16, 0, 32, 46, 0, time.Local),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.t.GetCreationTime()

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestMaprSecret_GetStatusString(t *testing.T) {
	tests := []struct {
		name          string
		t             *MaprSecret
		shouldContain string
	}{
		{
			name:          "nil",
			t:             nil,
			shouldContain: "Not found / Invalid",
		},
		{
			name:          "empty",
			t:             &MaprSecret{},
			shouldContain: "No secret found",
		},
		{
			name: "empty secret",
			t: &MaprSecret{
				Secret: &Secret{},
			},
			shouldContain: "No ticket found",
		},
		{
			name: "demo.mapr.com",
			t: NewMaprSecret(
				&Secret{
					Data: map[string][]byte{
						ticket.SecretMaprTicketKey: testTicketsRaw[0],
					},
				},
			),
			shouldContain: "Valid",
		},
		{
			name: "demo.mapr.com",
			t: NewMaprSecret(
				&Secret{
					Data: map[string][]byte{
						ticket.SecretMaprTicketKey: testTicketsRaw[1],
					},
				},
			),
			shouldContain: "Expired",
		},
		{
			name: "demo.mapr.com",
			t: NewMaprSecret(
				&Secret{
					Data: map[string][]byte{
						ticket.SecretMaprTicketKey: testTicketsRaw[2],
					},
				},
			),
			shouldContain: "Expired",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.t.GetStatusString()

			assert.Contains(t, got, tt.shouldContain)
		})
	}
}
