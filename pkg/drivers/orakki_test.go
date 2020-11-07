package drivers

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrakkiPodObject(t *testing.T) {
	orakkiImage := "registry.gitlab.com/oraksil/orakki:latest"
	gipanImage := "registry.gitlab.com/oraksil/gipan:latest"
	drv, err := NewK8SOrakkiDriver("../../configs/kube/config.local", "", "", "", orakkiImage, gipanImage, "", "", "", "", "")
	assert.Nil(t, err)
	assert.NotNil(t, drv)
	assert.Equal(t, "orakki", drv.baseAppName)

	po := drv.createOrakkiPod(drv.baseAppName, "dino")
	assert.True(t, strings.HasPrefix(po.ObjectMeta.Name, drv.baseAppName))
	assert.Equal(t, drv.namespace, po.Namespace)

	// orakkiId, err := drv.RunInstance("abcd")
	// assert.Nil(t, err)
	// assert.True(t, strings.HasPrefix(orakkiId, "orakki-"))

	// err = drv.DeleteInstance(orakkiId)
	// assert.Nil(t, err)
}
