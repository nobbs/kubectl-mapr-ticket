// Copyright (c) 2024 Alexej Disterhoft
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
//
// SPDX-License-Identifier: MIT

package types_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/nobbs/kubectl-mapr-ticket/pkg/types"
	"github.com/nobbs/kubectl-mapr-ticket/pkg/util"

	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestPersistentVolume_GetName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		v    *PersistentVolume
		want string
	}{
		{
			name: "nil",
			v:    nil,
			want: "",
		},
		{
			name: "empty",
			v:    &PersistentVolume{},
			want: "",
		},
		{
			name: "name",
			v: &PersistentVolume{

				ObjectMeta: metaV1.ObjectMeta{
					Name: "test",
				},
			},
			want: "test",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got := test.v.GetName()

			assert.Equal(t, test.want, got)
		})
	}
}

func TestPersistentVolume_GetClaimName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		v    *PersistentVolume
		want string
	}{
		{
			name: "nil",
			v:    nil,
			want: "",
		},
		{
			name: "empty",
			v:    &PersistentVolume{},
			want: "",
		},
		{
			name: "claim name",
			v: &PersistentVolume{
				Spec: coreV1.PersistentVolumeSpec{
					ClaimRef: &coreV1.ObjectReference{
						Name: "test",
					},
				},
			},
			want: "test",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got := test.v.GetClaimName()

			assert.Equal(t, test.want, got)
		})
	}
}

func TestPersistentVolume_GetClaimNamespace(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		v    *PersistentVolume
		want string
	}{
		{
			name: "nil",
			v:    nil,
			want: "",
		},
		{
			name: "empty",
			v:    &PersistentVolume{},
			want: "",
		},
		{
			name: "claim namespace",
			v: &PersistentVolume{
				Spec: coreV1.PersistentVolumeSpec{
					ClaimRef: &coreV1.ObjectReference{
						Namespace: "test",
					},
				},
			},
			want: "test",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got := test.v.GetClaimNamespace()

			assert.Equal(t, test.want, got)
		})
	}
}

func TestPersistentVolume_GetVolumePath(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		v    *PersistentVolume
		want string
	}{
		{
			name: "nil",
			v:    nil,
			want: "",
		},
		{
			name: "empty",
			v:    &PersistentVolume{},
			want: "",
		},
		{
			name: "volume path",
			v: &PersistentVolume{
				Spec: coreV1.PersistentVolumeSpec{
					PersistentVolumeSource: coreV1.PersistentVolumeSource{
						CSI: &coreV1.CSIPersistentVolumeSource{
							VolumeAttributes: map[string]string{
								"volumePath": "test",
							},
						},
					},
				},
			},
			want: "test",
		},
		{
			name: "no volume path",
			v: &PersistentVolume{
				Spec: coreV1.PersistentVolumeSpec{
					PersistentVolumeSource: coreV1.PersistentVolumeSource{
						CSI: &coreV1.CSIPersistentVolumeSource{
							VolumeAttributes: map[string]string{},
						},
					},
				},
			},
			want: "",
		},
		{
			name: "nil volume attributes",
			v: &PersistentVolume{
				Spec: coreV1.PersistentVolumeSpec{
					PersistentVolumeSource: coreV1.PersistentVolumeSource{
						CSI: &coreV1.CSIPersistentVolumeSource{},
					},
				},
			},
			want: "",
		},
		{
			name: "nil CSI",
			v: &PersistentVolume{
				Spec: coreV1.PersistentVolumeSpec{},
			},
			want: "",
		},
		{
			name: "nil spec",
			v:    &PersistentVolume{},
			want: "",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got := test.v.GetVolumePath()

			assert.Equal(t, test.want, got)
		})
	}
}

func TestPersistentVolume_GetVolumeHandle(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		v    *PersistentVolume
		want string
	}{
		{
			name: "nil",
			v:    nil,
			want: "",
		},
		{
			name: "empty",
			v:    &PersistentVolume{},
			want: "",
		},
		{
			name: "volume handle",
			v: &PersistentVolume{
				Spec: coreV1.PersistentVolumeSpec{
					PersistentVolumeSource: coreV1.PersistentVolumeSource{
						CSI: &coreV1.CSIPersistentVolumeSource{
							VolumeHandle: "test",
						},
					},
				},
			},
			want: "test",
		},
		{
			name: "no volume handle",
			v: &PersistentVolume{
				Spec: coreV1.PersistentVolumeSpec{
					PersistentVolumeSource: coreV1.PersistentVolumeSource{
						CSI: &coreV1.CSIPersistentVolumeSource{},
					},
				},
			},
			want: "",
		},
		{
			name: "nil CSI",
			v: &PersistentVolume{
				Spec: coreV1.PersistentVolumeSpec{},
			},
			want: "",
		},
		{
			name: "nil spec",
			v:    &PersistentVolume{},
			want: "",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got := test.v.GetVolumeHandle()

			assert.Equal(t, test.want, got)
		})
	}
}

func TestPersistentVolume_GetSecretName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		v    *PersistentVolume
		want string
	}{
		{
			name: "nil",
			v:    nil,
			want: "",
		},
		{
			name: "empty",
			v:    &PersistentVolume{},
			want: "",
		},
		{
			name: "secret name",
			v: &PersistentVolume{
				Spec: coreV1.PersistentVolumeSpec{
					PersistentVolumeSource: coreV1.PersistentVolumeSource{
						CSI: &coreV1.CSIPersistentVolumeSource{
							NodePublishSecretRef: &coreV1.SecretReference{
								Name: "test",
							},
						},
					},
				},
			},
			want: "test",
		},
		{
			name: "no secret name",
			v: &PersistentVolume{
				Spec: coreV1.PersistentVolumeSpec{
					PersistentVolumeSource: coreV1.PersistentVolumeSource{
						CSI: &coreV1.CSIPersistentVolumeSource{},
					},
				},
			},
			want: "",
		},
		{
			name: "nil CSI",
			v: &PersistentVolume{
				Spec: coreV1.PersistentVolumeSpec{},
			},
			want: "",
		},
		{
			name: "nil spec",
			v:    &PersistentVolume{},
			want: "",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got := test.v.GetSecretName()

			assert.Equal(t, test.want, got)
		})
	}
}

func TestPersistentVolume_GetSecretNamespace(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		v    *PersistentVolume
		want string
	}{
		{
			name: "nil",
			v:    nil,
			want: "",
		},
		{
			name: "empty",
			v:    &PersistentVolume{},
			want: "",
		},
		{
			name: "secret name",
			v: &PersistentVolume{
				Spec: coreV1.PersistentVolumeSpec{
					PersistentVolumeSource: coreV1.PersistentVolumeSource{
						CSI: &coreV1.CSIPersistentVolumeSource{
							NodePublishSecretRef: &coreV1.SecretReference{
								Namespace: "test",
							},
						},
					},
				},
			},
			want: "test",
		},
		{
			name: "no secret name",
			v: &PersistentVolume{
				Spec: coreV1.PersistentVolumeSpec{
					PersistentVolumeSource: coreV1.PersistentVolumeSource{
						CSI: &coreV1.CSIPersistentVolumeSource{},
					},
				},
			},
			want: "",
		},
		{
			name: "nil CSI",
			v: &PersistentVolume{
				Spec: coreV1.PersistentVolumeSpec{},
			},
			want: "",
		},
		{
			name: "nil spec",
			v:    &PersistentVolume{},
			want: "",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got := test.v.GetSecretNamespace()

			assert.Equal(t, test.want, got)
		})
	}
}

func TestPersistentVolume_IsMaprCSIBased(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		v    *PersistentVolume
		want bool
	}{
		{
			name: "nil",
			v:    nil,
			want: false,
		},
		{
			name: "empty",
			v:    &PersistentVolume{},
			want: false,
		},
		{
			name: "not MapR CSI based",
			v: &PersistentVolume{
				Spec: coreV1.PersistentVolumeSpec{
					PersistentVolumeSource: coreV1.PersistentVolumeSource{
						CSI: &coreV1.CSIPersistentVolumeSource{
							Driver: "test",
						},
					},
				},
			},
			want: false,
		},
		{
			name: "MapR CSI based",
			v: &PersistentVolume{
				Spec: coreV1.PersistentVolumeSpec{
					PersistentVolumeSource: coreV1.PersistentVolumeSource{
						CSI: &coreV1.CSIPersistentVolumeSource{
							Driver: MaprCSIProvisionerKDF,
						},
					},
				},
			},
			want: true,
		},
		{
			name: "MapR CSI nfs based",
			v: &PersistentVolume{
				Spec: coreV1.PersistentVolumeSpec{
					PersistentVolumeSource: coreV1.PersistentVolumeSource{
						CSI: &coreV1.CSIPersistentVolumeSource{
							Driver: MaprCSIProvisionerNFSKDF,
						},
					},
				},
			},
			want: true,
		},
		{
			name: "nil CSI",
			v: &PersistentVolume{
				Spec: coreV1.PersistentVolumeSpec{},
			},
			want: false,
		},
		{
			name: "nil spec",
			v:    &PersistentVolume{},
			want: false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got := test.v.IsMaprCSIBased()

			assert.Equal(t, test.want, got)
		})
	}
}

func TestPersistentVolume_UsesSecret(t *testing.T) {
	type args struct {
		namespace string
		name      string
	}

	t.Parallel()

	tests := []struct {
		name string
		v    *PersistentVolume
		args args
		want bool
	}{
		{
			name: "nil",
			v:    nil,
			args: args{
				namespace: "test",
				name:      "test",
			},
			want: false,
		},
		{
			name: "empty",
			v:    &PersistentVolume{},
			args: args{
				namespace: "test",
				name:      "test",
			},
			want: false,
		},
		{
			name: "uses secret",
			v: &PersistentVolume{
				Spec: coreV1.PersistentVolumeSpec{
					PersistentVolumeSource: coreV1.PersistentVolumeSource{
						CSI: &coreV1.CSIPersistentVolumeSource{
							NodePublishSecretRef: &coreV1.SecretReference{
								Namespace: "test",
								Name:      "test",
							},
						},
					},
				},
			},
			args: args{
				namespace: "test",
				name:      "test",
			},
			want: true,
		},
		{
			name: "does not use secret",
			v: &PersistentVolume{
				Spec: coreV1.PersistentVolumeSpec{
					PersistentVolumeSource: coreV1.PersistentVolumeSource{
						CSI: &coreV1.CSIPersistentVolumeSource{},
					},
				},
			},
			args: args{
				namespace: "test",
				name:      "test",
			},
			want: false,
		},
		{
			name: "nil CSI",
			v: &PersistentVolume{
				Spec: coreV1.PersistentVolumeSpec{},
			},
			args: args{
				namespace: "test",
				name:      "test",
			},
			want: false,
		},
		{
			name: "nil spec",
			v:    &PersistentVolume{},
			args: args{
				namespace: "test",
				name:      "test",
			},
			want: false,
		},
		{
			name: "wrong secret name",
			v: &PersistentVolume{
				Spec: coreV1.PersistentVolumeSpec{
					PersistentVolumeSource: coreV1.PersistentVolumeSource{
						CSI: &coreV1.CSIPersistentVolumeSource{
							NodePublishSecretRef: &coreV1.SecretReference{
								Namespace: "test",
								Name:      "test",
							},
						},
					},
				},
			},
			args: args{
				namespace: "test",
				name:      "wrong",
			},
			want: false,
		},
		{
			name: "wrong namespace",
			v: &PersistentVolume{
				Spec: coreV1.PersistentVolumeSpec{
					PersistentVolumeSource: coreV1.PersistentVolumeSource{
						CSI: &coreV1.CSIPersistentVolumeSource{
							NodePublishSecretRef: &coreV1.SecretReference{
								Namespace: "test",
								Name:      "test",
							},
						},
					},
				},
			},
			args: args{
				namespace: "wrong",
				name:      "test",
			},
			want: false,
		},
		{
			name: "all secrets in namespace",
			v: &PersistentVolume{
				Spec: coreV1.PersistentVolumeSpec{
					PersistentVolumeSource: coreV1.PersistentVolumeSource{
						CSI: &coreV1.CSIPersistentVolumeSource{
							NodePublishSecretRef: &coreV1.SecretReference{
								Namespace: "test",
								Name:      "test",
							},
						},
					},
				},
			},
			args: args{
				namespace: "test",
				name:      util.SecretAll,
			},
			want: true,
		},
		{
			name: "all secrets in namespace, wrong namespace",
			v: &PersistentVolume{
				Spec: coreV1.PersistentVolumeSpec{
					PersistentVolumeSource: coreV1.PersistentVolumeSource{
						CSI: &coreV1.CSIPersistentVolumeSource{
							NodePublishSecretRef: &coreV1.SecretReference{
								Namespace: "test",
								Name:      "test",
							},
						},
					},
				},
			},
			args: args{
				namespace: "wrong",
				name:      util.SecretAll,
			},
			want: false,
		},
		{
			name: "all secrets in all namespaces",
			v: &PersistentVolume{
				Spec: coreV1.PersistentVolumeSpec{
					PersistentVolumeSource: coreV1.PersistentVolumeSource{
						CSI: &coreV1.CSIPersistentVolumeSource{
							NodePublishSecretRef: &coreV1.SecretReference{
								Namespace: "test",
								Name:      "test",
							},
						},
					},
				},
			},
			args: args{
				namespace: util.NamespaceAll,
				name:      "does not matter",
			},
			want: true,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got := test.v.UsesSecret(test.args.namespace, test.args.name)

			assert.Equal(t, test.want, got)
		})
	}
}
