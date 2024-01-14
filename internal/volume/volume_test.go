package volume_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/nobbs/kubectl-mapr-ticket/internal/volume"

	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
)

const (
	CSIProvisionerMapr    = "com.mapr.csi-kdf"
	CSIProvisionerMaprNFS = "com.mapr.csi-nfskdf"
)

func TestLister_List(t *testing.T) {
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
					newCSIVolume("csi-volume", CSIProvisionerMapr, "mapr", "mapr-secret"),
				),
				secretName: "mapr-secret",
				namespace:  "mapr",
			},
			want: []expectedVolume{
				newExpectedSecret("csi-volume"),
			},
			wantErr: false,
		},
		{
			name: "one other csi volume for one secret",
			fields: listerFields{
				client: fake.NewSimpleClientset(
					newCSIVolume("csi-volume", "some-other-csi-driver", "mapr", "mapr-secret"),
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
					newCSIVolume("csi-volume-1", CSIProvisionerMapr, "mapr", "mapr-secret"),
					newCSIVolume("csi-volume-2", CSIProvisionerMapr, "mapr", "mapr-secret"),
				),
				secretName: "mapr-secret",
				namespace:  "mapr",
			},
			want: []expectedVolume{
				newExpectedSecret("csi-volume-1"),
				newExpectedSecret("csi-volume-2"),
			},
			wantErr: false,
		},
		{
			name: "three mapr volumes for two secrets, one secret specified",
			fields: listerFields{
				client: fake.NewSimpleClientset(
					newCSIVolume("csi-volume-1", CSIProvisionerMapr, "mapr", "mapr-secret-1"),
					newCSIVolume("csi-volume-2", CSIProvisionerMapr, "mapr", "mapr-secret-2"),
					newCSIVolume("csi-volume-3", CSIProvisionerMaprNFS, "mapr", "mapr-secret-2"),
				),
				secretName: "mapr-secret-2",
				namespace:  "mapr",
			},
			want: []expectedVolume{
				newExpectedSecret("csi-volume-2"),
				newExpectedSecret("csi-volume-3"),
			},
			wantErr: false,
		},
		{
			name: "one mapr volume for one secret in another namespace",
			fields: listerFields{
				client: fake.NewSimpleClientset(
					newCSIVolume("csi-volume", CSIProvisionerMapr, "mapr", "mapr-secret"),
				),
				secretName: "mapr-secret",
				namespace:  "mapr-2",
			},
			want:    []expectedVolume{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLister(tt.fields.client, tt.fields.secretName, tt.fields.namespace)

			got, err := l.List()

			assertVolumes(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestLister_ListWithAllSecrets(t *testing.T) {
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
				secretName: SecretAll,
				namespace:  "mapr",
			},
			want:    []expectedVolume{},
			wantErr: false,
		},
		{
			name: "one mapr volume for one secret",
			fields: listerFields{
				client: fake.NewSimpleClientset(
					newCSIVolume("csi-volume", CSIProvisionerMapr, "mapr", "mapr-secret"),
				),
				secretName: SecretAll,
				namespace:  "mapr",
			},
			want: []expectedVolume{
				newExpectedSecret("csi-volume"),
			},
			wantErr: false,
		},
		{
			name: "one other csi volume for one secret",
			fields: listerFields{
				client: fake.NewSimpleClientset(
					newCSIVolume("csi-volume", "some-other-csi-driver", "mapr", "mapr-secret"),
				),
				secretName: SecretAll,
				namespace:  "mapr",
			},
			want:    []expectedVolume{},
			wantErr: false,
		},
		{
			name: "two mapr volumes for one secret",
			fields: listerFields{
				client: fake.NewSimpleClientset(
					newCSIVolume("csi-volume-1", CSIProvisionerMapr, "mapr", "mapr-secret"),
					newCSIVolume("csi-volume-2", CSIProvisionerMapr, "mapr", "mapr-secret"),
				),
				secretName: SecretAll,
				namespace:  "mapr",
			},
			want: []expectedVolume{
				newExpectedSecret("csi-volume-1"),
				newExpectedSecret("csi-volume-2"),
			},
			wantErr: false,
		},
		{
			name: "three mapr volumes for two secrets, one secret specified",
			fields: listerFields{
				client: fake.NewSimpleClientset(
					newCSIVolume("csi-volume-1", CSIProvisionerMapr, "mapr", "mapr-secret-1"),
					newCSIVolume("csi-volume-2", CSIProvisionerMapr, "mapr", "mapr-secret-2"),
					newCSIVolume("csi-volume-3", CSIProvisionerMaprNFS, "mapr", "mapr-secret-2"),
				),
				secretName: SecretAll,
				namespace:  "mapr",
			},
			want: []expectedVolume{
				newExpectedSecret("csi-volume-1"),
				newExpectedSecret("csi-volume-2"),
				newExpectedSecret("csi-volume-3"),
			},
			wantErr: false,
		},
		{
			name: "one mapr volume for one secret in another namespace",
			fields: listerFields{
				client: fake.NewSimpleClientset(
					newCSIVolume("csi-volume", CSIProvisionerMapr, "mapr", "mapr-secret"),
				),
				secretName: SecretAll,
				namespace:  "mapr-2",
			},
			want:    []expectedVolume{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLister(tt.fields.client, tt.fields.secretName, tt.fields.namespace)

			got, err := l.List()

			assertVolumes(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestLister_ListWithAllNamespaces(t *testing.T) {
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
				namespace:  metaV1.NamespaceAll,
			},
			want:    []expectedVolume{},
			wantErr: false,
		},
		{
			name: "one mapr volume for one secret",
			fields: listerFields{
				client: fake.NewSimpleClientset(
					newCSIVolume("csi-volume", CSIProvisionerMapr, "mapr", "mapr-secret"),
				),
				secretName: "we-dont-care-about-this-when-namespace-is-all",
				namespace:  metaV1.NamespaceAll,
			},
			want: []expectedVolume{
				newExpectedSecret("csi-volume"),
			},
			wantErr: false,
		},
		{
			name: "one other csi volume for one secret",
			fields: listerFields{
				client: fake.NewSimpleClientset(
					newCSIVolume("csi-volume", "some-other-csi-driver", "mapr", "mapr-secret"),
				),
				secretName: SecretAll,
				namespace:  metaV1.NamespaceAll,
			},
			want:    []expectedVolume{},
			wantErr: false,
		},
		{
			name: "two mapr volumes for two secrets in different namespaces",
			fields: listerFields{
				client: fake.NewSimpleClientset(
					newCSIVolume("csi-volume-1", CSIProvisionerMapr, "mapr-1", "mapr-secret"),
					newCSIVolume("csi-volume-2", CSIProvisionerMapr, "mapr-2", "mapr-secret"),
				),
				secretName: "",
				namespace:  metaV1.NamespaceAll,
			},
			want: []expectedVolume{
				newExpectedSecret("csi-volume-1"),
				newExpectedSecret("csi-volume-2"),
			},
			wantErr: false,
		},
		{
			name: "three mapr volumes for three secrets in different namespaces",
			fields: listerFields{
				client: fake.NewSimpleClientset(
					newCSIVolume("csi-volume-1", CSIProvisionerMapr, "default", "mapr-secret-1"),
					newCSIVolume("csi-volume-2", CSIProvisionerMapr, "kube-system", "mapr-secret-2"),
					newCSIVolume("csi-volume-3", CSIProvisionerMaprNFS, "kube-public", "mapr-secret-2"),
				),
				secretName: "mapr-secret-2",
				namespace:  metaV1.NamespaceAll,
			},
			want: []expectedVolume{
				newExpectedSecret("csi-volume-1"),
				newExpectedSecret("csi-volume-2"),
				newExpectedSecret("csi-volume-3"),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLister(tt.fields.client, tt.fields.secretName, tt.fields.namespace)

			got, err := l.List()

			assertVolumes(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

type listerFields struct {
	client     kubernetes.Interface
	secretName string
	namespace  string
}

type expectedVolume struct {
	name string
}

func newExpectedSecret(name string) expectedVolume {
	return expectedVolume{
		name: name,
	}
}

func assertVolumes(t *testing.T, expected []expectedVolume, actual []coreV1.PersistentVolume) {
	t.Helper()

	assert := assert.New(t)

	assert.Len(actual, len(expected))

	for i, e := range expected {
		assert.Equal(e.name, actual[i].Name)
	}
}

func newCSIVolume(name, driver, secretNamespace, secretName string) *coreV1.PersistentVolume {
	return &coreV1.PersistentVolume{
		ObjectMeta: metaV1.ObjectMeta{
			Name: name,
		},
		Spec: coreV1.PersistentVolumeSpec{
			PersistentVolumeSource: coreV1.PersistentVolumeSource{
				CSI: &coreV1.CSIPersistentVolumeSource{
					Driver: driver,
					NodePublishSecretRef: &coreV1.SecretReference{
						Namespace: secretNamespace,
						Name:      secretName,
					},
				},
			},
		},
	}
}
