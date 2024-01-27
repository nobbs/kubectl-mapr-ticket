// Copyright (c) 2024 Alexej Disterhoft
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
//
// SPDX-License-Identifier: MIT

package volume_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/nobbs/kubectl-mapr-ticket/pkg/secret"
	"github.com/nobbs/kubectl-mapr-ticket/pkg/ticket"
	"github.com/nobbs/kubectl-mapr-ticket/pkg/types"
	"github.com/nobbs/kubectl-mapr-ticket/pkg/util"
	"github.com/nobbs/kubectl-mapr-ticket/pkg/volume"
	. "github.com/nobbs/kubectl-mapr-ticket/pkg/volume"
	"github.com/nobbs/mapr-ticket-parser/pkg/parse"

	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
)

const (
	CSIProvisionerMapr    = "com.mapr.csi-kdf"
	CSIProvisionerMaprNFS = "com.mapr.csi-nfskdf"
)

func TestLister_Default(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		fields  listerFields
		want    []expectedVolume
		wantErr bool
	}{
		{
			name: "no volumes",
			fields: listerFields{
				client:     fake.NewSimpleClientset(),
				secretName: "mapr-secret",
				namespace:  "mapr",
			},
			want:    []expectedVolume{},
			wantErr: false,
		},
		{
			name: "one mapr volume for one secret",
			fields: listerFields{
				client: fake.NewSimpleClientset(
					newCSIVolume("csi-volume", CSIProvisionerMapr, withSecretRef("mapr", "mapr-secret")),
				),
				secretName: "mapr-secret",
				namespace:  "mapr",
			},
			want: []expectedVolume{
				expectVolume("csi-volume"),
			},
			wantErr: false,
		},
		{
			name: "one other csi volume for one secret",
			fields: listerFields{
				client: fake.NewSimpleClientset(
					newCSIVolume("csi-volume", "some-other-csi-driver", withSecretRef("mapr", "mapr-secret")),
				),
				secretName: "mapr-secret",
				namespace:  "mapr",
			},
			want:    []expectedVolume{},
			wantErr: false,
		},
		{
			name: "two mapr volumes for one secret",
			fields: listerFields{
				client: fake.NewSimpleClientset(
					newCSIVolume("csi-volume-1", CSIProvisionerMapr, withSecretRef("mapr", "mapr-secret")),
					newCSIVolume("csi-volume-2", CSIProvisionerMapr, withSecretRef("mapr", "mapr-secret")),
				),
				secretName: "mapr-secret",
				namespace:  "mapr",
			},
			want: []expectedVolume{
				expectVolume("csi-volume-1"),
				expectVolume("csi-volume-2"),
			},
			wantErr: false,
		},
		{
			name: "three mapr volumes for two secrets, one secret specified",
			fields: listerFields{
				client: fake.NewSimpleClientset(
					newCSIVolume("csi-volume-1", CSIProvisionerMapr, withSecretRef("mapr", "mapr-secret-1")),
					newCSIVolume("csi-volume-2", CSIProvisionerMapr, withSecretRef("mapr", "mapr-secret-2")),
					newCSIVolume("csi-volume-3", CSIProvisionerMaprNFS, withSecretRef("mapr", "mapr-secret-2")),
				),
				secretName: "mapr-secret-2",
				namespace:  "mapr",
			},
			want: []expectedVolume{
				expectVolume("csi-volume-2"),
				expectVolume("csi-volume-3"),
			},
			wantErr: false,
		},
		{
			name: "one mapr volume for one secret in another namespace",
			fields: listerFields{
				client: fake.NewSimpleClientset(
					newCSIVolume("csi-volume", CSIProvisionerMapr, withSecretRef("mapr", "mapr-secret")),
				),
				secretName: "mapr-secret",
				namespace:  "mapr-2",
			},
			want:    []expectedVolume{},
			wantErr: false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			l := NewLister(test.fields.client, test.fields.secretName, test.fields.namespace)

			got, err := l.List()

			assertVolumes(t, test.want, got)
			assert.Equal(t, test.wantErr, err != nil)
		})
	}
}

func TestLister_WithAllSecrets(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		fields  listerFields
		want    []expectedVolume
		wantErr bool
	}{
		{
			name: "no volumes",
			fields: listerFields{
				client:     fake.NewSimpleClientset(),
				secretName: util.SecretAll,
				namespace:  "mapr",
			},
			want:    []expectedVolume{},
			wantErr: false,
		},
		{
			name: "one mapr volume for one secret",
			fields: listerFields{
				client: fake.NewSimpleClientset(
					newCSIVolume("csi-volume", CSIProvisionerMapr, withSecretRef("mapr", "mapr-secret")),
				),
				secretName: util.SecretAll,
				namespace:  "mapr",
			},
			want: []expectedVolume{
				expectVolume("csi-volume"),
			},
			wantErr: false,
		},
		{
			name: "one other csi volume for one secret",
			fields: listerFields{
				client: fake.NewSimpleClientset(
					newCSIVolume("csi-volume", "some-other-csi-driver", withSecretRef("mapr", "mapr-secret")),
				),
				secretName: util.SecretAll,
				namespace:  "mapr",
			},
			want:    []expectedVolume{},
			wantErr: false,
		},
		{
			name: "two mapr volumes for one secret",
			fields: listerFields{
				client: fake.NewSimpleClientset(
					newCSIVolume("csi-volume-1", CSIProvisionerMapr, withSecretRef("mapr", "mapr-secret")),
					newCSIVolume("csi-volume-2", CSIProvisionerMapr, withSecretRef("mapr", "mapr-secret")),
				),
				secretName: util.SecretAll,
				namespace:  "mapr",
			},
			want: []expectedVolume{
				expectVolume("csi-volume-1"),
				expectVolume("csi-volume-2"),
			},
			wantErr: false,
		},
		{
			name: "three mapr volumes for two secrets, one secret specified",
			fields: listerFields{
				client: fake.NewSimpleClientset(
					newCSIVolume("csi-volume-1", CSIProvisionerMapr, withSecretRef("mapr", "mapr-secret-1")),
					newCSIVolume("csi-volume-2", CSIProvisionerMapr, withSecretRef("mapr", "mapr-secret-2")),
					newCSIVolume("csi-volume-3", CSIProvisionerMaprNFS, withSecretRef("mapr", "mapr-secret-2")),
				),
				secretName: util.SecretAll,
				namespace:  "mapr",
			},
			want: []expectedVolume{
				expectVolume("csi-volume-1"),
				expectVolume("csi-volume-2"),
				expectVolume("csi-volume-3"),
			},
			wantErr: false,
		},
		{
			name: "one mapr volume for one secret in another namespace",
			fields: listerFields{
				client: fake.NewSimpleClientset(
					newCSIVolume("csi-volume", CSIProvisionerMapr, withSecretRef("mapr", "mapr-secret")),
				),
				secretName: util.SecretAll,
				namespace:  "mapr-2",
			},
			want:    []expectedVolume{},
			wantErr: false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			l := NewLister(test.fields.client, test.fields.secretName, test.fields.namespace)

			got, err := l.List()

			assertVolumes(t, test.want, got)
			assert.Equal(t, test.wantErr, err != nil)
		})
	}
}

func TestLister_WithAllNamespaces(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		fields  listerFields
		want    []expectedVolume
		wantErr bool
	}{
		{
			name: "no volumes",
			fields: listerFields{
				client:     fake.NewSimpleClientset(),
				secretName: "mapr-secret",
				namespace:  util.NamespaceAll,
			},
			want:    []expectedVolume{},
			wantErr: false,
		},
		{
			name: "one mapr volume for one secret",
			fields: listerFields{
				client: fake.NewSimpleClientset(
					newCSIVolume("csi-volume", CSIProvisionerMapr, withSecretRef("mapr", "mapr-secret")),
				),
				secretName: "we-dont-care-about-this-when-namespace-is-all",
				namespace:  util.NamespaceAll,
			},
			want: []expectedVolume{
				expectVolume("csi-volume"),
			},
			wantErr: false,
		},
		{
			name: "one other csi volume for one secret",
			fields: listerFields{
				client: fake.NewSimpleClientset(
					newCSIVolume("csi-volume", "some-other-csi-driver", withSecretRef("mapr", "mapr-secret")),
				),
				secretName: util.SecretAll,
				namespace:  util.NamespaceAll,
			},
			want:    []expectedVolume{},
			wantErr: false,
		},
		{
			name: "two mapr volumes for two secrets in different namespaces",
			fields: listerFields{
				client: fake.NewSimpleClientset(
					newCSIVolume("csi-volume-1", CSIProvisionerMapr, withSecretRef("mapr-1", "mapr-secret")),
					newCSIVolume("csi-volume-2", CSIProvisionerMapr, withSecretRef("mapr-2", "mapr-secret")),
				),
				secretName: "",
				namespace:  util.NamespaceAll,
			},
			want: []expectedVolume{
				expectVolume("csi-volume-1"),
				expectVolume("csi-volume-2"),
			},
			wantErr: false,
		},
		{
			name: "three mapr volumes for three secrets in different namespaces",
			fields: listerFields{
				client: fake.NewSimpleClientset(
					newCSIVolume("csi-volume-1", CSIProvisionerMapr, withSecretRef("default", "mapr-secret-1")),
					newCSIVolume("csi-volume-2", CSIProvisionerMapr, withSecretRef("kube-public", "mapr-secret-2")),
					newCSIVolume("csi-volume-3", CSIProvisionerMaprNFS, withSecretRef("kube-system", "mapr-secret-3")),
				),
				secretName: "mapr-secret-2",
				namespace:  util.NamespaceAll,
			},
			want: []expectedVolume{
				expectVolume("csi-volume-1"),
				expectVolume("csi-volume-2"),
				expectVolume("csi-volume-3"),
			},
			wantErr: false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			l := NewLister(test.fields.client, test.fields.secretName, test.fields.namespace)

			got, err := l.List()

			assertVolumes(t, test.want, got)
			assert.Equal(t, test.wantErr, err != nil)
		})
	}
}

func TestLister_WithSortByName(t *testing.T) {
	opts := []ListerOption{
		volume.WithSortBy([]SortOption{SortByName}),
	}

	t.Parallel()

	tests := []struct {
		name    string
		fields  listerFields
		want    []expectedVolume
		wantErr bool
	}{
		{
			name: "no volumes",
			fields: listerFields{
				client:     fake.NewSimpleClientset(),
				secretName: util.SecretAll,
				namespace:  util.NamespaceAll,
				opts:       opts,
			},
			want:    []expectedVolume{},
			wantErr: false,
		},
		{
			name: "three mapr volumes for three secrets in different namespaces",
			fields: listerFields{
				client: fake.NewSimpleClientset(
					newCSIVolume("csi-volume-3", CSIProvisionerMaprNFS, withSecretRef("kube-system", "mapr-secret-3")),
					newCSIVolume("csi-volume-1", CSIProvisionerMapr, withSecretRef("default", "mapr-secret-1")),
					newCSIVolume("csi-volume-2", CSIProvisionerMapr, withSecretRef("kube-public", "mapr-secret-2")),
				),
				secretName: util.SecretAll,
				namespace:  util.NamespaceAll,
				opts:       opts,
			},
			want: []expectedVolume{
				expectVolume("csi-volume-1"),
				expectVolume("csi-volume-2"),
				expectVolume("csi-volume-3"),
			},
			wantErr: false,
		},
		{
			name: "three mapr volumes for three secrets in different namespaces, one secret specified",
			fields: listerFields{
				client: fake.NewSimpleClientset(
					newCSIVolume("csi-volume-3", CSIProvisionerMaprNFS, withSecretRef("kube-public", "mapr-secret-2")),
					newCSIVolume("csi-volume-1", CSIProvisionerMapr, withSecretRef("default", "mapr-secret-1")),
					newCSIVolume("csi-volume-2", CSIProvisionerMapr, withSecretRef("kube-public", "mapr-secret-2")),
				),
				secretName: "mapr-secret-2",
				namespace:  "kube-public",
				opts:       opts,
			},
			want: []expectedVolume{
				expectVolume("csi-volume-2"),
				expectVolume("csi-volume-3"),
			},
			wantErr: false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			l := NewLister(test.fields.client, test.fields.secretName, test.fields.namespace, test.fields.opts...)

			got, err := l.List()

			assertVolumes(t, test.want, got)
			assert.Equal(t, test.wantErr, err != nil)
		})
	}
}

func TestLister_WithSortBySecretNamespace(t *testing.T) {
	opts := []ListerOption{
		volume.WithSortBy([]SortOption{SortBySecretNamespace}),
	}

	t.Parallel()

	tests := []struct {
		name    string
		fields  listerFields
		want    []expectedVolume
		wantErr bool
	}{
		{
			name: "no volumes",
			fields: listerFields{
				client:     fake.NewSimpleClientset(),
				secretName: util.SecretAll,
				namespace:  util.NamespaceAll,
				opts:       opts,
			},
			want:    []expectedVolume{},
			wantErr: false,
		},
		{
			name: "three mapr volumes for three secrets in different namespaces",
			fields: listerFields{
				client: fake.NewSimpleClientset(
					newCSIVolume("csi-volume-3", CSIProvisionerMaprNFS, withSecretRef("kube-system", "mapr-secret-3")),
					newCSIVolume("csi-volume-1", CSIProvisionerMapr, withSecretRef("default", "mapr-secret-1")),
					newCSIVolume("csi-volume-2", CSIProvisionerMapr, withSecretRef("kube-public", "mapr-secret-2")),
				),
				secretName: util.SecretAll,
				namespace:  util.NamespaceAll,
				opts:       opts,
			},
			want: []expectedVolume{
				expectVolume("csi-volume-1"),
				expectVolume("csi-volume-2"),
				expectVolume("csi-volume-3"),
			},
			wantErr: false,
		},
		{
			name: "three mapr volumes for three secrets in different namespaces, one secret namespace specified",
			fields: listerFields{
				client: fake.NewSimpleClientset(
					newCSIVolume("csi-volume-3", CSIProvisionerMaprNFS, withSecretRef("kube-public", "mapr-secret-3")),
					newCSIVolume("csi-volume-1", CSIProvisionerMapr, withSecretRef("default", "mapr-secret-1")),
					newCSIVolume("csi-volume-2", CSIProvisionerMapr, withSecretRef("kube-public", "mapr-secret-2")),
				),
				secretName: util.SecretAll,
				namespace:  "kube-public",
				opts:       opts,
			},
			want: []expectedVolume{
				expectVolume("csi-volume-2"),
				expectVolume("csi-volume-3"),
			},
			wantErr: false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			l := NewLister(test.fields.client, test.fields.secretName, test.fields.namespace, test.fields.opts...)

			got, err := l.List()

			assertVolumes(t, test.want, got)
			assert.Equal(t, test.wantErr, err != nil)
		})
	}
}

func TestLister_WithSortBySecretName(t *testing.T) {
	opts := []ListerOption{
		volume.WithSortBy([]SortOption{SortBySecretName}),
	}

	t.Parallel()

	tests := []struct {
		name    string
		fields  listerFields
		want    []expectedVolume
		wantErr bool
	}{
		{
			name: "no volumes",
			fields: listerFields{
				client:     fake.NewSimpleClientset(),
				secretName: util.SecretAll,
				namespace:  util.NamespaceAll,
				opts:       opts,
			},
			want:    []expectedVolume{},
			wantErr: false,
		},
		{
			name: "three mapr volumes for three secrets in different namespaces",
			fields: listerFields{
				client: fake.NewSimpleClientset(
					newCSIVolume("csi-volume-3", CSIProvisionerMaprNFS, withSecretRef("kube-system", "mapr-secret-3")),
					newCSIVolume("csi-volume-2", CSIProvisionerMapr, withSecretRef("kube-public", "mapr-secret-2")),
					newCSIVolume("csi-volume-1", CSIProvisionerMapr, withSecretRef("default", "mapr-secret-1")),
				),
				secretName: util.SecretAll,
				namespace:  util.NamespaceAll,
				opts:       opts,
			},
			want: []expectedVolume{
				expectVolume("csi-volume-1"),
				expectVolume("csi-volume-2"),
				expectVolume("csi-volume-3"),
			},
			wantErr: false,
		},
		{
			name: "three mapr volumes for three secrets in different namespaces, one secret namespace specified",
			fields: listerFields{
				client: fake.NewSimpleClientset(
					newCSIVolume("csi-volume-1", CSIProvisionerMapr, withSecretRef("default", "mapr-secret-1")),
					newCSIVolume("csi-volume-3", CSIProvisionerMaprNFS, withSecretRef("kube-public", "mapr-secret-3")),
					newCSIVolume("csi-volume-2", CSIProvisionerMapr, withSecretRef("kube-public", "mapr-secret-2")),
				),
				secretName: util.SecretAll,
				namespace:  "kube-public",
				opts:       opts,
			},
			want: []expectedVolume{
				expectVolume("csi-volume-2"),
				expectVolume("csi-volume-3"),
			},
			wantErr: false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			l := NewLister(test.fields.client, test.fields.secretName, test.fields.namespace, test.fields.opts...)

			got, err := l.List()

			assertVolumes(t, test.want, got)
			assert.Equal(t, test.wantErr, err != nil)
		})
	}
}

func TestLister_WithSortByClaimNamespace(t *testing.T) {
	opts := []ListerOption{
		volume.WithSortBy([]SortOption{SortByClaimNamespace}),
	}

	t.Parallel()

	tests := []struct {
		name    string
		fields  listerFields
		want    []expectedVolume
		wantErr bool
	}{
		{
			name: "no volumes",
			fields: listerFields{
				client:     fake.NewSimpleClientset(),
				secretName: util.SecretAll,
				namespace:  util.NamespaceAll,
				opts:       opts,
			},
			want:    []expectedVolume{},
			wantErr: false,
		},
		{
			name: "three mapr volumes for three secrets in different namespaces",
			fields: listerFields{
				client: fake.NewSimpleClientset(
					newCSIVolume(
						"csi-volume-3", CSIProvisionerMaprNFS,
						withSecretRef("kube-system", "mapr-secret-3"),
						withClaimRef("kube-system", "mapr-claim-3"),
					),
					newCSIVolume(
						"csi-volume-1", CSIProvisionerMapr,
						withSecretRef("default", "mapr-secret-1"),
						withClaimRef("default", "mapr-claim-1"),
					),
					newCSIVolume(
						"csi-volume-2", CSIProvisionerMapr,
						withSecretRef("kube-public", "mapr-secret-2"),
						withClaimRef("kube-public", "mapr-claim-2"),
					),
				),
				secretName: util.SecretAll,
				namespace:  util.NamespaceAll,
				opts:       opts,
			},
			want: []expectedVolume{
				expectVolume("csi-volume-1"),
				expectVolume("csi-volume-2"),
				expectVolume("csi-volume-3"),
			},
			wantErr: false,
		},
		{
			name: "three mapr volumes for three secrets in different namespaces, one secret namespace specified",
			fields: listerFields{
				client: fake.NewSimpleClientset(
					newCSIVolume(
						"csi-volume-1", CSIProvisionerMapr,
						withSecretRef("default", "mapr-secret-1"),
						withClaimRef("default", "mapr-claim-1"),
					),
					newCSIVolume(
						"csi-volume-3", CSIProvisionerMaprNFS,
						withSecretRef("kube-public", "mapr-secret-3"),
						withClaimRef("kube-public", "mapr-claim-3"),
					),
					newCSIVolume(
						"csi-volume-2", CSIProvisionerMapr,
						withSecretRef("kube-public", "mapr-secret-2"),
						withClaimRef("kube-public", "mapr-claim-2"),
					),
				),
				secretName: util.SecretAll,
				namespace:  "kube-public",
				opts:       opts,
			},
			want: []expectedVolume{
				expectVolume("csi-volume-2"),
				expectVolume("csi-volume-3"),
			},
			wantErr: false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			l := NewLister(test.fields.client, test.fields.secretName, test.fields.namespace, test.fields.opts...)

			got, err := l.List()

			assertVolumes(t, test.want, got)
			assert.Equal(t, test.wantErr, err != nil)
		})
	}
}

func TestLister_WithSortByClaimName(t *testing.T) {
	opts := []ListerOption{
		volume.WithSortBy([]SortOption{SortByClaimName}),
	}

	t.Parallel()

	tests := []struct {
		name    string
		fields  listerFields
		want    []expectedVolume
		wantErr bool
	}{
		{
			name: "no volumes",
			fields: listerFields{
				client:     fake.NewSimpleClientset(),
				secretName: util.SecretAll,
				namespace:  util.NamespaceAll,
				opts:       opts,
			},
			want:    []expectedVolume{},
			wantErr: false,
		},
		{
			name: "three mapr volumes for three secrets in different namespaces",
			fields: listerFields{
				client: fake.NewSimpleClientset(
					newCSIVolume(
						"csi-volume-3", CSIProvisionerMaprNFS,
						withSecretRef("kube-system", "mapr-secret-3"),
						withClaimRef("kube-system", "mapr-claim-3"),
					),
					newCSIVolume(
						"csi-volume-1", CSIProvisionerMapr,
						withSecretRef("default", "mapr-secret-1"),
						withClaimRef("default", "mapr-claim-1"),
					),
					newCSIVolume(
						"csi-volume-2", CSIProvisionerMapr,
						withSecretRef("kube-public", "mapr-secret-2"),
						withClaimRef("kube-public", "mapr-claim-2"),
					),
				),
				secretName: util.SecretAll,
				namespace:  util.NamespaceAll,
				opts:       opts,
			},
			want: []expectedVolume{
				expectVolume("csi-volume-1"),
				expectVolume("csi-volume-2"),
				expectVolume("csi-volume-3"),
			},
			wantErr: false,
		},
		{
			name: "three mapr volumes for three secrets in different namespaces, one secret namespace specified",
			fields: listerFields{
				client: fake.NewSimpleClientset(
					newCSIVolume(
						"csi-volume-1", CSIProvisionerMapr,
						withSecretRef("default", "mapr-secret-1"),
						withClaimRef("default", "mapr-claim-1"),
					),
					newCSIVolume(
						"csi-volume-3", CSIProvisionerMaprNFS,
						withSecretRef("kube-public", "mapr-secret-3"),
						withClaimRef("kube-public", "mapr-claim-3"),
					),
					newCSIVolume(
						"csi-volume-2", CSIProvisionerMapr,
						withSecretRef("kube-public", "mapr-secret-2"),
						withClaimRef("kube-public", "mapr-claim-2"),
					),
				),
				secretName: util.SecretAll,
				namespace:  "kube-public",
				opts:       opts,
			},
			want: []expectedVolume{
				expectVolume("csi-volume-2"),
				expectVolume("csi-volume-3"),
			},
			wantErr: false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			l := NewLister(test.fields.client, test.fields.secretName, test.fields.namespace, test.fields.opts...)

			got, err := l.List()

			assertVolumes(t, test.want, got)
			assert.Equal(t, test.wantErr, err != nil)
		})
	}
}

func TestLister_WithSortByVolumePath(t *testing.T) {
	opts := []ListerOption{
		volume.WithSortBy([]SortOption{SortByVolumePath}),
	}

	t.Parallel()

	tests := []struct {
		name    string
		fields  listerFields
		want    []expectedVolume
		wantErr bool
	}{
		{
			name: "no volumes",
			fields: listerFields{
				client:     fake.NewSimpleClientset(),
				secretName: util.SecretAll,
				namespace:  util.NamespaceAll,
				opts:       opts,
			},
			want:    []expectedVolume{},
			wantErr: false,
		},
		{
			name: "three mapr volumes for three secrets in different namespaces",
			fields: listerFields{
				client: fake.NewSimpleClientset(
					newCSIVolume(
						"csi-volume-3", CSIProvisionerMaprNFS,
						withVolumePath("/mapr"),
						withSecretRef("kube-system", "mapr-secret-3"),
					),
					newCSIVolume(
						"csi-volume-2", CSIProvisionerMapr,
						withVolumePath("/mapr-2"),
						withSecretRef("kube-public", "mapr-secret-2"),
					),
					newCSIVolume(
						"csi-volume-1", CSIProvisionerMapr,
						withVolumePath("/mapr-1"),
						withSecretRef("default", "mapr-secret-1"),
					),
				),
				secretName: util.SecretAll,
				namespace:  util.NamespaceAll,
				opts:       opts,
			},
			want: []expectedVolume{
				expectVolume("csi-volume-3"),
				expectVolume("csi-volume-1"),
				expectVolume("csi-volume-2"),
			},
			wantErr: false,
		},
		{
			name: "three mapr volumes for three secrets in different namespaces, one secret namespace specified",
			fields: listerFields{
				client: fake.NewSimpleClientset(
					newCSIVolume(
						"csi-volume-1", CSIProvisionerMapr,
						withVolumePath("/mapr-1"),
						withSecretRef("default", "mapr-secret-1"),
					),
					newCSIVolume(
						"csi-volume-3", CSIProvisionerMaprNFS,
						withVolumePath("/mapr-3"),
						withSecretRef("kube-public", "mapr-secret-3"),
					),
					newCSIVolume(
						"csi-volume-2", CSIProvisionerMapr,
						withVolumePath("/mapr-2"),
						withSecretRef("kube-public", "mapr-secret-2"),
					),
				),
				secretName: util.SecretAll,
				namespace:  "kube-public",
				opts:       opts,
			},
			want: []expectedVolume{
				expectVolume("csi-volume-2"),
				expectVolume("csi-volume-3"),
			},
			wantErr: false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			l := NewLister(test.fields.client, test.fields.secretName, test.fields.namespace, test.fields.opts...)

			got, err := l.List()

			assertVolumes(t, test.want, got)
			assert.Equal(t, test.wantErr, err != nil)
		})
	}
}

func TestLister_WithSortByVolumeHandle(t *testing.T) {
	opts := []ListerOption{
		volume.WithSortBy([]SortOption{SortByVolumeHandle}),
	}

	t.Parallel()

	tests := []struct {
		name    string
		fields  listerFields
		want    []expectedVolume
		wantErr bool
	}{
		{
			name: "no volumes",
			fields: listerFields{
				client:     fake.NewSimpleClientset(),
				secretName: util.SecretAll,
				namespace:  util.NamespaceAll,
				opts:       opts,
			},
			want:    []expectedVolume{},
			wantErr: false,
		},
		{
			name: "three mapr volumes for three secrets in different namespaces",
			fields: listerFields{
				client: fake.NewSimpleClientset(
					newCSIVolume(
						"csi-volume-3", CSIProvisionerMaprNFS,
						withVolumeHandle("csi.mapr"),
						withSecretRef("kube-system", "mapr-secret-3"),
					),
					newCSIVolume(
						"csi-volume-2", CSIProvisionerMapr,
						withVolumeHandle("csi.mapr-2"),
						withSecretRef("kube-public", "mapr-secret-2"),
					),
					newCSIVolume(
						"csi-volume-1", CSIProvisionerMapr,
						withVolumeHandle("csi.mapr-1"),
						withSecretRef("default", "mapr-secret-1"),
					),
				),
				secretName: util.SecretAll,
				namespace:  util.NamespaceAll,
				opts:       opts,
			},
			want: []expectedVolume{
				expectVolume("csi-volume-3"),
				expectVolume("csi-volume-1"),
				expectVolume("csi-volume-2"),
			},
			wantErr: false,
		},
		{
			name: "three mapr volumes for three secrets in different namespaces, one secret namespace specified",
			fields: listerFields{
				client: fake.NewSimpleClientset(
					newCSIVolume(
						"csi-volume-1", CSIProvisionerMapr,
						withVolumeHandle("csi.mapr-1"),
						withSecretRef("default", "mapr-secret-1"),
					),
					newCSIVolume(
						"csi-volume-3", CSIProvisionerMaprNFS,
						withVolumeHandle("csi.mapr-3"),
						withSecretRef("kube-public", "mapr-secret-3"),
					),
					newCSIVolume(
						"csi-volume-2", CSIProvisionerMapr,
						withVolumeHandle("csi.mapr-2"),
						withSecretRef("kube-public", "mapr-secret-2"),
					),
				),
				secretName: util.SecretAll,
				namespace:  "kube-public",
				opts:       opts,
			},
			want: []expectedVolume{
				expectVolume("csi-volume-2"),
				expectVolume("csi-volume-3"),
			},
			wantErr: false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			l := NewLister(test.fields.client, test.fields.secretName, test.fields.namespace, test.fields.opts...)

			got, err := l.List()

			assertVolumes(t, test.want, got)
			assert.Equal(t, test.wantErr, err != nil)
		})
	}
}

func TestLister_WithSortByExpiryTime(t *testing.T) {
	opts := []ListerOption{
		volume.WithSortBy([]SortOption{SortByExpiration}),
	}

	t.Parallel()

	tests := []struct {
		name    string
		fields  listerFields
		want    []expectedVolume
		wantErr bool
	}{
		{
			name: "no volumes",
			fields: listerFields{
				client:     fake.NewSimpleClientset(),
				secretName: util.SecretAll,
				namespace:  util.NamespaceAll,
				opts:       opts,
			},
			want:    []expectedVolume{},
			wantErr: false,
		},
		{
			name: "three mapr volumes for three secrets in different namespaces",
			fields: listerFields{
				client: fake.NewSimpleClientset(
					newCSIVolume(
						"csi-volume-1", CSIProvisionerMapr,
						withSecretRef("default", "mapr-secret"),
					),
					secretFromTicketJSON(
						t,
						"default",
						"mapr-secret",
						ticketWithExpiryTime(t, time.Now().Add(1*time.Hour)),
					),
					newCSIVolume(
						"csi-volume-2", CSIProvisionerMaprNFS,
						withSecretRef("kube-system", "mapr-secret"),
					),
					secretFromTicketJSON(
						t,
						"kube-system",
						"mapr-secret",
						ticketWithExpiryTime(t, time.Now().Add(5*time.Hour)),
					),
					newCSIVolume(
						"csi-volume-3", CSIProvisionerMapr,
						withSecretRef("kube-public", "mapr-secret"),
					),
					secretFromTicketJSON(
						t,
						"kube-public",
						"mapr-secret",
						ticketWithExpiryTime(t, time.Now().Add(-1*time.Hour)),
					),
				),
				secretName: util.SecretAll,
				namespace:  util.NamespaceAll,
				opts:       opts,
			},
			want: []expectedVolume{
				expectVolume("csi-volume-3"),
				expectVolume("csi-volume-1"),
				expectVolume("csi-volume-2"),
			},
			wantErr: false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			l := NewLister(test.fields.client, test.fields.secretName, test.fields.namespace, append(
				test.fields.opts,
				WithSecretLister(secret.NewLister(test.fields.client, util.NamespaceAll)),
			)...)

			got, err := l.List()

			assertVolumes(t, test.want, got)
			assert.Equal(t, test.wantErr, err != nil)
		})
	}
}

func TestLister_WithSortByAge(t *testing.T) {
	opts := []ListerOption{
		volume.WithSortBy([]SortOption{SortByAge}),
	}

	t.Parallel()

	tests := []struct {
		name    string
		fields  listerFields
		want    []expectedVolume
		wantErr bool
	}{
		{
			name: "no volumes",
			fields: listerFields{
				client:     fake.NewSimpleClientset(),
				secretName: util.SecretAll,
				namespace:  util.NamespaceAll,
				opts:       opts,
			},
			want:    []expectedVolume{},
			wantErr: false,
		},
		{
			name: "three mapr volumes for three secrets in different namespaces",
			fields: listerFields{
				client: fake.NewSimpleClientset(
					newCSIVolume(
						"csi-volume-1", CSIProvisionerMapr,
						withSecretRef("default", "mapr-secret"),
						withCreateTimeRelative(-1*time.Hour),
					),
					newCSIVolume(
						"csi-volume-2", CSIProvisionerMaprNFS,
						withSecretRef("kube-system", "mapr-secret"),
						withCreateTimeRelative(5*time.Hour),
					),
					newCSIVolume(
						"csi-volume-3", CSIProvisionerMapr,
						withSecretRef("kube-public", "mapr-secret"),
						withCreateTimeRelative(1*time.Hour),
					),
				),
				secretName: util.SecretAll,
				namespace:  util.NamespaceAll,
				opts:       opts,
			},
			want: []expectedVolume{
				expectVolume("csi-volume-1"),
				expectVolume("csi-volume-3"),
				expectVolume("csi-volume-2"),
			},
			wantErr: false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			l := NewLister(test.fields.client, test.fields.secretName, test.fields.namespace, append(
				test.fields.opts,
				WithSecretLister(secret.NewLister(test.fields.client, util.NamespaceAll)),
			)...)

			got, err := l.List()

			assertVolumes(t, test.want, got)
			assert.Equal(t, test.wantErr, err != nil)
		})
	}
}

type listerFields struct {
	client     kubernetes.Interface
	secretName string
	namespace  string
	opts       []ListerOption
}

type expectedVolume struct {
	name string
}

func expectVolume(name string) expectedVolume {
	return expectedVolume{
		name: name,
	}
}

func assertVolumes(t *testing.T, expected []expectedVolume, actual []types.MaprVolume) {
	t.Helper()

	assert.Len(t, actual, len(expected))

	for i, e := range expected {
		assert.Equal(t, e.name, actual[i].Volume.Name)
	}
}

type csiVolume struct {
	name            string
	driver          string
	secretNamespace string
	secretName      string
	claimNamespace  string
	claimName       string
	volumePath      string
	volumeHandle    string
	createTime      time.Time
}

type csiVolumeOption func(*csiVolume)

func withSecretRef(secretNamespace, secretName string) csiVolumeOption {
	return func(v *csiVolume) {
		v.secretNamespace = secretNamespace
		v.secretName = secretName
	}
}

func withClaimRef(claimNamespace, claimName string) csiVolumeOption {
	return func(v *csiVolume) {
		v.claimNamespace = claimNamespace
		v.claimName = claimName
	}
}

func withVolumePath(volumePath string) csiVolumeOption {
	return func(v *csiVolume) {
		v.volumePath = volumePath
	}
}

func withVolumeHandle(volumeHandle string) csiVolumeOption {
	return func(v *csiVolume) {
		v.volumeHandle = volumeHandle
	}
}

func withCreateTimeRelative(duration time.Duration) csiVolumeOption {
	return func(v *csiVolume) {
		v.createTime = time.Now().Add(duration)
	}
}

func newCSIVolume(name, driver string, opts ...csiVolumeOption) *coreV1.PersistentVolume {
	v := &csiVolume{
		name:   name,
		driver: driver,
	}

	for _, opt := range opts {
		opt(v)
	}

	return &coreV1.PersistentVolume{
		ObjectMeta: metaV1.ObjectMeta{
			Name:              v.name,
			CreationTimestamp: metaV1.NewTime(v.createTime),
		},
		Spec: coreV1.PersistentVolumeSpec{
			ClaimRef: &coreV1.ObjectReference{
				Namespace: v.claimNamespace,
				Name:      v.claimName,
			},
			PersistentVolumeSource: coreV1.PersistentVolumeSource{
				CSI: &coreV1.CSIPersistentVolumeSource{
					Driver: v.driver,
					NodePublishSecretRef: &coreV1.SecretReference{
						Namespace: v.secretNamespace,
						Name:      v.secretName,
					},
					VolumeHandle: v.volumeHandle,
					VolumeAttributes: map[string]string{
						"volumePath": v.volumePath,
					},
				},
			},
		},
	}
}

func ticketWithExpiryTime(t *testing.T, expiryTime time.Time) []byte {
	unix := uint64(expiryTime.Unix())
	return []byte(fmt.Sprintf(`{"ticket":{"expiryTime":%d}}`, unix))
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
