package drivers

import (
	"fmt"
	"path/filepath"

	gonanoid "github.com/matoous/go-nanoid"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type K8SOrakkiDriver struct {
	kubeConfigPath string
	orakkiImage    string
	namespace      string
	baseAppName    string

	kubeConfig *restclient.Config
	kubeOpSet  *kubernetes.Clientset
}

func (d *K8SOrakkiDriver) RunInstance(peerName string) (string, error) {
	podName := d.newOrakkiPodName()
	podObj := d.createOrakkiPod(podName, peerName)

	_, err := d.kubeOpSet.CoreV1().Pods(d.namespace).Create(podObj)
	if err != nil {
		return "", err
	}

	return podName, nil
}

func (d *K8SOrakkiDriver) DeleteInstance(id string) error {
	podName := id

	var gracePeriod int64 = 0
	delOpts := metav1.DeleteOptions{GracePeriodSeconds: &gracePeriod}

	err := d.kubeOpSet.CoreV1().Pods(d.namespace).Delete(podName, &delOpts)
	if err != nil {
		return err
	}

	return nil
}

func (d *K8SOrakkiDriver) newOrakkiPodName() string {
	id, _ := gonanoid.Generate("abcdef", 7)
	return fmt.Sprintf("%s-%s", d.baseAppName, id)
}

func (d *K8SOrakkiDriver) createOrakkiPod(podName, peerName string) *core.Pod {
	return &core.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      podName,
			Namespace: d.namespace,
		},
		Spec: core.PodSpec{
			Containers: []core.Container{
				{
					Name:            podName,
					Image:           d.orakkiImage,
					ImagePullPolicy: core.PullIfNotPresent,
					Env: []core.EnvVar{
						{
							Name:  "PEER_NAME",
							Value: peerName,
						},
					},
				},
			},
		},
	}
}

func NewK8SOrakkiDriver(kubeConfigPath, orakkiImage string) (*K8SOrakkiDriver, error) {
	if kubeConfigPath == "" {
		kubeConfigPath = filepath.Join(homedir.HomeDir(), ".kube", "config")
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &K8SOrakkiDriver{
		kubeConfigPath: kubeConfigPath,
		orakkiImage:    orakkiImage,
		namespace:      "oraksil-dev",
		baseAppName:    "orakki",
		kubeConfig:     config,
		kubeOpSet:      clientset,
	}, nil
}
