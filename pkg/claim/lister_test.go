package claim_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	. "github.com/nobbs/kubectl-mapr-ticket/pkg/claim"
	"github.com/nobbs/kubectl-mapr-ticket/pkg/types"
	"github.com/nobbs/kubectl-mapr-ticket/pkg/util"

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
		name     string
		fields   listerFields
		expected []expecetedClaim
		wantErr  bool
	}{
		{
			name: "no claims",
			fields: listerFields{
				client:    fake.NewSimpleClientset(),
				namespace: "default",
			},
			expected: []expecetedClaim{},
			wantErr:  false,
		},
		{
			name: "one claim without volume",
			fields: listerFields{
				client: fake.NewSimpleClientset(
					newClaim("default", "claim-1"),
				),
				namespace: "default",
			},
			expected: []expecetedClaim{},
			wantErr:  false,
		},
		{
			name: "one claim with volume",
			fields: listerFields{
				client: fake.NewSimpleClientset(
					newClaim("default", "claim-1", withVolumeName("volume-1"), withPhase(coreV1.ClaimBound)),
					newCSIVolume("volume-1", CSIProvisionerMapr, withClaimRef("default", "claim-1")),
				),
				namespace: "default",
			},
			expected: []expecetedClaim{
				expectClaim("claim-1", "default"),
			},
			wantErr: false,
		},
		{
			name: "three claims, one with mapr volume, one with random volume, one unbound",
			fields: listerFields{
				client: fake.NewSimpleClientset(
					newClaim("default", "claim-1", withVolumeName("volume-1"), withPhase(coreV1.ClaimBound)),
					newClaim("default", "claim-2", withVolumeName("volume-2"), withPhase(coreV1.ClaimBound)),
					newClaim("default", "claim-3", withVolumeName("volume-3"), withPhase(coreV1.ClaimPending)),
					newCSIVolume("volume-1", CSIProvisionerMapr, withClaimRef("default", "claim-1")),
					newCSIVolume("volume-2", "random-provisioner", withClaimRef("default", "claim-2")),
				),
				namespace: "default",
			},
			expected: []expecetedClaim{
				expectClaim("claim-1", "default"),
			},
			wantErr: false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			l := NewLister(test.fields.client, test.fields.namespace, test.fields.opts...)

			actual, err := l.List()

			assertClaims(t, test.expected, actual)
			assert.Equal(t, test.wantErr, err != nil)
		})
	}
}

func TestLister_WithAllNamespaces(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		fields   listerFields
		expected []expecetedClaim
		wantErr  bool
	}{
		{
			name: "no claims",
			fields: listerFields{
				client:    fake.NewSimpleClientset(),
				namespace: util.NamespaceAll,
				opts:      []ListerOption{},
			},
			expected: []expecetedClaim{},
			wantErr:  false,
		},
		{
			name: "three claims, two with mapr volumes, one other",
			fields: listerFields{
				client: fake.NewSimpleClientset(
					newClaim("default", "claim-1", withVolumeName("volume-1"), withPhase(coreV1.ClaimBound)),
					newClaim("kube-public", "claim-2", withVolumeName("volume-2"), withPhase(coreV1.ClaimBound)),
					newClaim("kube-system", "claim-3", withVolumeName("volume-3"), withPhase(coreV1.ClaimBound)),
					newCSIVolume("volume-1", CSIProvisionerMapr, withClaimRef("default", "claim-1")),
					newCSIVolume("volume-2", CSIProvisionerMapr, withClaimRef("kube-public", "claim-2")),
					newCSIVolume("volume-3", "random-provisioner", withClaimRef("kube-system", "claim-3")),
				),
				namespace: util.NamespaceAll,
			},
			expected: []expecetedClaim{
				expectClaim("claim-1", "default"),
				expectClaim("claim-2", "kube-public"),
			},
			wantErr: false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			l := NewLister(test.fields.client, test.fields.namespace, test.fields.opts...)

			actual, err := l.List()

			assertClaims(t, test.expected, actual)
			assert.Equal(t, test.wantErr, err != nil)
		})
	}
}

type listerFields struct {
	client    kubernetes.Interface
	namespace string
	opts      []ListerOption
}

type expecetedClaim struct {
	name      string
	namespace string
}

func expectClaim(name, namespace string) expecetedClaim {
	return expecetedClaim{
		name:      name,
		namespace: namespace,
	}
}

func assertClaims(t *testing.T, expected []expecetedClaim, actual []types.MaprVolumeClaim) {
	t.Helper()

	assert.Len(t, actual, len(expected))

	for i, claim := range actual {
		assert.Equal(t, expected[i].name, claim.Claim.Name)
		assert.Equal(t, expected[i].namespace, claim.Claim.Namespace)
	}
}

type claimOptions struct {
	name       string
	namespace  string
	volumeName string
	phase      coreV1.PersistentVolumeClaimPhase
}

type claimOption func(*claimOptions)

func withVolumeName(volumeName string) claimOption {
	return func(c *claimOptions) {
		c.volumeName = volumeName
	}
}

func withPhase(phase coreV1.PersistentVolumeClaimPhase) claimOption {
	return func(c *claimOptions) {
		c.phase = phase
	}
}

func newClaim(namespace, name string, opts ...claimOption) *coreV1.PersistentVolumeClaim {
	c := &claimOptions{
		name:      name,
		namespace: namespace,
	}

	for _, opt := range opts {
		opt(c)
	}

	return &coreV1.PersistentVolumeClaim{
		ObjectMeta: metaV1.ObjectMeta{
			Name:      c.name,
			Namespace: c.namespace,
		},
		Spec: coreV1.PersistentVolumeClaimSpec{
			VolumeName: c.volumeName,
		},
		Status: coreV1.PersistentVolumeClaimStatus{
			Phase: c.phase,
		},
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
