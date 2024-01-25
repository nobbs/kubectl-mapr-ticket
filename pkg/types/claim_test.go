package types_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/nobbs/kubectl-mapr-ticket/pkg/types"

	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestPersistentVolumeClaim_GetNamespace(t *testing.T) {
	tests := []struct {
		name string
		c    *PersistentVolumeClaim
		want string
	}{
		{
			name: "nil",
			c:    nil,
			want: "",
		},
		{
			name: "empty",
			c:    &PersistentVolumeClaim{},
			want: "",
		},
		{
			name: "namespace",
			c: &PersistentVolumeClaim{
				ObjectMeta: metaV1.ObjectMeta{
					Namespace: "test",
				},
			},
			want: "test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.c.GetNamespace()

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestPersistentVolumeClaim_GetName(t *testing.T) {
	tests := []struct {
		name string
		c    *PersistentVolumeClaim
		want string
	}{
		{
			name: "nil",
			c:    nil,
			want: "",
		},
		{
			name: "empty",
			c:    &PersistentVolumeClaim{},
			want: "",
		},
		{
			name: "name",
			c: &PersistentVolumeClaim{
				ObjectMeta: metaV1.ObjectMeta{
					Name: "test",
				},
			},
			want: "test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.c.GetName()

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestPersistentVolumeClaim_IsBound(t *testing.T) {
	tests := []struct {
		name string
		c    *PersistentVolumeClaim
		want bool
	}{
		{
			name: "nil",
			c:    nil,
			want: false,
		},
		{
			name: "empty",
			c:    &PersistentVolumeClaim{},
			want: false,
		},
		{
			name: "unbound",
			c: &PersistentVolumeClaim{
				Status: coreV1.PersistentVolumeClaimStatus{
					Phase: coreV1.ClaimPending,
				},
			},
			want: false,
		},
		{
			name: "bound",
			c: &PersistentVolumeClaim{
				Status: coreV1.PersistentVolumeClaimStatus{
					Phase: coreV1.ClaimBound,
				},
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.c.IsBound()

			assert.Equal(t, tt.want, got)
		})
	}
}
