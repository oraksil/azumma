package drivers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrakkiPodObject(t *testing.T) {
	// drv := drivers.K8SOrakkiDriver{}
	po := createPodObject()
	assert.Equal(t, nil, po)

	// assert.Equal(t, "Pod", po.Kind)
}
