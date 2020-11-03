package drivers

import (
	"errors"
	"fmt"

	gonanoid "github.com/matoous/go-nanoid"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type K8SOrakkiDriver struct {
	kubeConfigPath string

	namespace         string
	baseAppName       string
	nodeSelectorKey   string
	nodeSelectorValue string

	orakkiImage string
	gipanImage  string

	mqRpcUri string
	mqRpcNs  string

	turnUri      string
	turnUsername string
	turnPassword string

	kubeConfig *restclient.Config
	kubeOpSet  *kubernetes.Clientset
}

func (d *K8SOrakkiDriver) RunInstance(id string) (string, error) {
	podName := d.newOrakkiPodName(id)
	podObj := d.createOrakkiPod(podName)

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

func (d *K8SOrakkiDriver) newOrakkiPodName(id string) string {
	if id != "" {
		return id
	}

	seedId, _ := gonanoid.Nanoid(7)
	return fmt.Sprintf("%s-%s", d.baseAppName, seedId)
}

func (d *K8SOrakkiDriver) createOrakkiPod(podName string) *core.Pod {
	pod := &core.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      podName,
			Namespace: d.namespace,
			Labels: map[string]string{
				"app": d.baseAppName,
			},
		},
		Spec: core.PodSpec{
			Containers: []core.Container{
				{
					Name:            "orakki",
					Image:           d.orakkiImage,
					ImagePullPolicy: core.PullIfNotPresent,
					Env: []core.EnvVar{
						{
							Name:  "MQRPC_URI",
							Value: d.mqRpcUri,
						},
						{
							Name:  "MQRPC_NAMESPACE",
							Value: d.mqRpcNs,
						},
						{
							Name:  "MQRPC_IDENTIFIER",
							Value: podName,
						},
						{
							Name:  "IPC_IMAGE_FRAMES",
							Value: "tcp://127.0.0.1:8765",
						},
						{
							Name:  "IPC_SOUND_FRAMES",
							Value: "tcp://127.0.0.1:8766",
						},
						{
							Name:  "IPC_KEY_INPUTS",
							Value: "tcp://127.0.0.1:8767",
						},
						{
							Name:  "TURN_URI",
							Value: d.turnUri,
						},
						{
							Name:  "TURN_USERNAME",
							Value: d.turnUsername,
						},
						{
							Name:  "TURN_PASSWORD",
							Value: d.turnPassword,
						},
					},
				},
				{
					Name:            "gipan",
					Image:           d.gipanImage,
					ImagePullPolicy: core.PullIfNotPresent,
					Env: []core.EnvVar{
						{
							Name:  "GAME",
							Value: "dino",
						},
						{
							Name:  "IPC_IMAGE_FRAMES",
							Value: "tcp://127.0.0.1:8765",
						},
						{
							Name:  "IPC_SOUND_FRAMES",
							Value: "tcp://127.0.0.1:8766",
						},
						{
							Name:  "IPC_KEY_INPUTS",
							Value: "tcp://127.0.0.1:8767",
						},
					},
				},
			},
		},
	}

	if d.nodeSelectorKey != "" && d.nodeSelectorValue != "" {
		pod.Spec.NodeSelector = map[string]string{
			d.nodeSelectorKey: d.nodeSelectorValue,
		}
	}

	return pod
}

func NewK8SOrakkiDriver(
	kubeConfigPath,
	namespace,
	nodeSelectorKey,
	nodeSelectorValue,
	orakkiImage,
	gipanImage,
	mqRpcUri,
	mqRpcNs,
	turnUri,
	turnUsername,
	turnPassword string) (*K8SOrakkiDriver, error) {

	if kubeConfigPath == "" {
		return nil, errors.New("invalid kube config path")
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
		kubeConfigPath:    kubeConfigPath,
		namespace:         namespace,
		nodeSelectorKey:   nodeSelectorKey,
		nodeSelectorValue: nodeSelectorValue,
		baseAppName:       "orakki",
		kubeConfig:        config,
		kubeOpSet:         clientset,
		orakkiImage:       orakkiImage,
		gipanImage:        gipanImage,
		mqRpcUri:          mqRpcUri,
		mqRpcNs:           mqRpcNs,
		turnUri:           turnUri,
		turnUsername:      turnUsername,
		turnPassword:      turnPassword,
	}, nil
}
