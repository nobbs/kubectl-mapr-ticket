package ticket_test

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	. "github.com/nobbs/kubectl-mapr-ticket/pkg/ticket"
	"github.com/nobbs/mapr-ticket-parser/pkg/parse"

	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
)

func TestNewTicketFromSecret(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		name   string
		secret *coreV1.Secret
		err    error
	}{
		{
			name: "secret contains MaprTicket",
			secret: func() *coreV1.Secret {
				secret := &coreV1.Secret{
					ObjectMeta: metaV1.ObjectMeta{
						Name:      "secret",
						Namespace: "default",
					},
					Data: map[string][]byte{
						SecretMaprTicketKey: []byte("demo.mapr.com +Cze+qwYCbAXGbz56OO7UF+lGqL3WPXrNkO1SLawEEDmSbgNl019xBeBY3kvh+R13iz/mCnwpzsLQw4Y5jEnv5GtuIWbeoC95ha8VKwX8MKcE6Kn9nZ2AF0QminkHwNVBx6TDriGZffyJCfZzivBwBSdKoQEWhBOPFCIMAi7w2zV/SX5Ut7u4qIKvEpr0JHV7sLMWYLhYncM6CKMd7iECGvECsBvEZRVj+dpbEY0BaRN/W54/7wNWaSVELUF6JWHQ8dmsqty4cZlI0/MV10HZzIbl9sMLFQ="),
					},
				}

				return secret
			}(),
			err: nil,
		},
		{
			name: "secret does not contain MaprTicket",
			secret: func() *coreV1.Secret {
				secret := &coreV1.Secret{
					ObjectMeta: metaV1.ObjectMeta{
						Name:      "secret",
						Namespace: "default",
					},
					Data: make(map[string][]byte),
				}

				return secret
			}(),
			err: NewErrSecretDoesNotContainMaprTicket("default", "secret"),
		},
		{
			name: "secret contains invalid MaprTicket",
			secret: func() *coreV1.Secret {
				secret := &coreV1.Secret{
					ObjectMeta: metaV1.ObjectMeta{
						Name:      "secret",
						Namespace: "default",
					},
					Data: map[string][]byte{
						SecretMaprTicketKey: []byte("invalid ticket"),
					},
				}

				return secret
			}(),
			err: errors.New("invalid mapr ticket: illegal base64 data at input byte 4"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := NewMaprTicketFromSecret(test.secret)

			if test.err == nil {
				assert.NoError(err)
			} else {
				assert.EqualError(err, test.err.Error())
			}
		})
	}
}

func TestSecretContainsMaprTicket(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		name     string
		secret   *coreV1.Secret
		expected bool
	}{
		{
			name: "secret contains MaprTicket",
			secret: func() *coreV1.Secret {
				secret := &coreV1.Secret{}
				secret.Data = make(map[string][]byte)
				secret.Data[SecretMaprTicketKey] = []byte("dummy ticket")
				return secret
			}(),
			expected: true,
		},
		{
			name: "secret does not contain MaprTicket",
			secret: func() *coreV1.Secret {
				secret := &coreV1.Secret{}
				secret.Data = make(map[string][]byte)
				return secret
			}(),
			expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := SecretContainsMaprTicket(test.secret)

			assert.Equal(test.expected, result)
		})
	}
}

func TestExpirationTime(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		name     string
		ticket   *Ticket
		expected time.Time
	}{
		{
			name: "ticket has expiry time",
			ticket: func() *Ticket {
				ticket := NewMaprTicket()
				ticket.TicketAndKey.ExpiryTime = ptr.To[uint64](1234567890)
				return ticket
			}(),
			expected: time.Unix(1234567890, 0),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.ticket.ExpirationTime()

			assert.Equal(test.expected, result)
		})
	}
}

func TestCreationTime(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		name     string
		ticket   *Ticket
		expected time.Time
	}{
		{
			name: "ticket has creation time",
			ticket: func() *Ticket {
				ticket := NewMaprTicket()
				ticket.TicketAndKey.CreationTimeSec = ptr.To[uint64](1234567890)
				return ticket
			}(),
			expected: time.Unix(1234567890, 0),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.ticket.CreationTime()

			assert.Equal(test.expected, result)
		})
	}
}

func TestIsExpired(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		name     string
		ticket   *Ticket
		expected bool
	}{
		{
			name: "ticket is not expired",
			ticket: func() *Ticket {
				ticket := NewMaprTicket()
				expiresInOneHour := time.Now().Add(1 * time.Hour).Unix()
				ticket.TicketAndKey.ExpiryTime = ptr.To[uint64](uint64(expiresInOneHour))
				return ticket
			}(),
			expected: false,
		},
		{
			name: "ticket is expired",
			ticket: func() *Ticket {
				ticket := NewMaprTicket()
				expiredOneHourAgo := time.Now().Add(-1 * time.Hour).Unix()
				ticket.TicketAndKey.ExpiryTime = ptr.To[uint64](uint64(expiredOneHourAgo))
				return ticket
			}(),
			expected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.ticket.IsExpired()

			assert.Equal(test.expected, result)
		})
	}
}

func TestExpiresBefore(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		name     string
		ticket   *Ticket
		time     time.Duration
		expected bool
	}{
		{
			name: "ticket expires before the provided time",
			ticket: func() *Ticket {
				ticket := NewMaprTicket()
				expiresInOneHour := time.Now().Add(1 * time.Hour).Unix()
				ticket.TicketAndKey.ExpiryTime = ptr.To[uint64](uint64(expiresInOneHour))
				return ticket
			}(),
			time:     2 * time.Hour,
			expected: true,
		},
		{
			name: "ticket does not expire before the provided time",
			ticket: func() *Ticket {
				ticket := NewMaprTicket()
				expiresInTwoHours := time.Now().Add(2 * time.Hour).Unix()
				ticket.TicketAndKey.ExpiryTime = ptr.To[uint64](uint64(expiresInTwoHours))
				return ticket
			}(),
			time:     1 * time.Hour,
			expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.ticket.ExpiresBefore(test.time)

			assert.Equal(test.expected, result)
		})
	}
}

func TestAsMaprTicket(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		name     string
		ticket   *Ticket
		expected *parse.MaprTicket
	}{
		{
			name:     "ticket is nil",
			ticket:   nil,
			expected: nil,
		},
		{
			name:     "ticket is not nil",
			ticket:   NewMaprTicket(),
			expected: parse.NewMaprTicket(),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.ticket.AsMaprTicket()

			assert.Equal(test.expected, result)
		})
	}
}
