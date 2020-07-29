package drivers

import (
	core "k8s.io/api/core/v1"
)

type K8SOrakkiDriver struct {
}

func (d *K8SOrakkiDriver) RunInstance(peerName string) (string, error) {
	return "", nil
}

func (d *K8SOrakkiDriver) DeleteInstance(id string) error {
	return nil
}

func createPodObject() *core.Pod {
	return nil
}
