package handler

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	listerv1 "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	db        = make(map[string]string)
	podlister listerv1.PodLister
)

func init() {
	var err error
	podlister, err = newPodLister()
	if err != nil {
		panic(err)
	}
}

func PingHandler(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}

func GetUserHandler(c *gin.Context) {
	user := c.Params.ByName("name")
	if value, ok := db[user]; ok {
		c.JSON(http.StatusOK, gin.H{"user": user, "value": value})
	} else {
		c.JSON(http.StatusOK, gin.H{"user": user, "status": "no value"})
	}
}

/* example curl for /admin with basicauth header
   Zm9vOmJhcg== is base64("foo:bar")
	curl -X POST \
  	http://localhost:8080/admin \
  	-H 'authorization: Basic Zm9vOmJhcg==' \
  	-H 'content-type: application/json' \
  	-d '{"value":"bar"}'
*/
func AdminHandler(c *gin.Context) {
	user := c.MustGet(gin.AuthUserKey).(string)

	// Parse JSON
	var json struct {
		Value string `json:"value" binding:"required"`
	}

	if c.Bind(&json) == nil {
		db[user] = json.Value
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	}
}

func GetPods(c *gin.Context) {
	ns := c.Params.ByName("ns")
	pods, err := podlister.Pods(ns).List(labels.Set{}.AsSelector())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// We only want pod names.
	var podNames []string
	for _, pod := range pods {
		podNames = append(podNames, pod.Name)
	}

	c.JSON(http.StatusOK, gin.H{"pods": podNames})
}

func newPodLister() (listerv1.PodLister, error) {
	defaultKubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")
	if _, err := os.Stat(defaultKubeconfig); err != nil {
		defaultKubeconfig = ""
	}
	log.Printf("using kubeconfig: %s", defaultKubeconfig)

	// Create a kubernetes config.
	config, err := clientcmd.BuildConfigFromFlags("", defaultKubeconfig)
	if err != nil {
		return nil, fmt.Errorf("get cluster K8s config failed: %v", err)
	}

	// Create the clientset.
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("make new k8s clientset failed: %v", err)
	}

	// Create informer.
	kubeInformers := informers.NewSharedInformerFactoryWithOptions(
		clientset, 5*time.Minute,
		informers.WithTweakListOptions(func(options *metav1.ListOptions) {
		}),
	)

	podInformer := kubeInformers.Core().V1().Pods()

	// Register pod informer
	inf := podInformer.Informer()

	log.Println("start pod informers")
	stopCh := make(chan struct{})
	kubeInformers.Start(stopCh)

	log.Println("wait for pod informer cache sync ready...")
	if !cache.WaitForCacheSync(stopCh, inf.HasSynced) {
		return nil, fmt.Errorf("failed to sync K8s pod cache")
	}

	log.Println("pod informer cache is synced done")
	return podInformer.Lister(), nil
}
