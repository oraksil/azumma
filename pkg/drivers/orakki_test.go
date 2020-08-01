package drivers

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrakkiPodObject(t *testing.T) {
	drv, err := NewK8SOrakkiDriver("", "busybox")
	assert.Nil(t, err)
	assert.NotNil(t, drv)
	assert.Equal(t, "orakki", drv.baseAppName)

	peerName := "test-peer-name"
	po := drv.createOrakkiPod(drv.baseAppName, peerName)
	assert.True(t, strings.HasPrefix(po.ObjectMeta.Name, drv.baseAppName))
	assert.Equal(t, drv.namespace, po.Namespace)
	assert.Equal(t, "PEER_NAME", po.Spec.Containers[0].Env[0].Name)
	assert.Equal(t, peerName, po.Spec.Containers[0].Env[0].Value)

	// orakkiId, err := drv.RunInstance("abcd")
	// assert.Nil(t, err)
	// assert.True(t, strings.HasPrefix(orakkiId, "orakki-"))

	// err = drv.DeleteInstance(orakkiId)
	// assert.Nil(t, err)
}
