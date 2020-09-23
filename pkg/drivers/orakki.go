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
	namespace      string
	baseAppName    string

	orakkiImage string
	gipanImage  string
	mqRpcUri    string
	mqRpcNs     string

	kubeConfig *restclient.Config
	kubeOpSet  *kubernetes.Clientset
}

func (d *K8SOrakkiDriver) RunInstance() (string, error) {
	podName := d.newOrakkiPodName()
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

func (d *K8SOrakkiDriver) newOrakkiPodName() string {
	id, _ := gonanoid.Generate("abcdef", 7)
	return fmt.Sprintf("%s-%s", d.baseAppName, id)
}

func (d *K8SOrakkiDriver) createOrakkiPod(podName string) *core.Pod {
	return &core.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      podName,
			Namespace: d.namespace,
			Labels: map[string]string{
				"app": d.baseAppName,
			},
		},
		Spec: core.PodSpec{
			Volumes: []core.Volume{
				{
					Name: "shared-vol-for-ipc",
					VolumeSource: core.VolumeSource{
						EmptyDir: &core.EmptyDirVolumeSource{
							Medium: "",
						},
					},
				},
			},
			Containers: []core.Container{
				{
					Name:            "orakki",
					Image:           d.orakkiImage,
					ImagePullPolicy: core.PullIfNotPresent,
					VolumeMounts: []core.VolumeMount{
						{
							Name:      "shared-vol-for-ipc",
							MountPath: "/var/oraksil/ipc",
							ReadOnly:  true,
						},
					},
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
							Value: "/var/oraksil/ipc/images.ipc",
						},
						{
							Name:  "IPC_SOUND_FRAMES",
							Value: "/var/oraksil/ipc/sounds.ipc",
						},
						{
							Name:  "IPC_KEY_INPUTS",
							Value: "/var/oraksil/ipc/keys.ipc",
						},
					},
				},
				{
					Name:            "gipan",
					Image:           d.gipanImage,
					ImagePullPolicy: core.PullIfNotPresent,
					VolumeMounts: []core.VolumeMount{
						{
							Name:      "shared-vol-for-ipc",
							MountPath: "/var/oraksil/ipc",
						},
					},
					Env: []core.EnvVar{
						{
							Name:  "GAME",
							Value: "dino",
						},
						{
							Name:  "IPC_IMAGE_FRAMES",
							Value: "/var/oraksil/ipc/images.ipc",
						},
						{
							Name:  "IPC_SOUND_FRAMES",
							Value: "/var/oraksil/ipc/sounds.ipc",
						},
						{
							Name:  "IPC_KEY_INPUTS",
							Value: "/var/oraksil/ipc/keys.ipc",
						},
					},
				},
			},
		},
	}
}

func NewK8SOrakkiDriver(kubeConfigPath, orakkiImage, gipanImage, mqRpcUri, mqRpcNs string) (*K8SOrakkiDriver, error) {
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
		namespace:      "oraksil-dev",
		baseAppName:    "orakki",
		kubeConfig:     config,
		kubeOpSet:      clientset,
		orakkiImage:    orakkiImage,
		gipanImage:     gipanImage,
		mqRpcUri:       mqRpcUri,
		mqRpcNs:        mqRpcNs,
	}, nil
}
