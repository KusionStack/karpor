package scheme

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"

	clusterinstall "code.alipay.com/multi-cluster/karbour/pkg/apis/cluster/install"
	clusterv1beta1 "code.alipay.com/multi-cluster/karbour/pkg/apis/cluster/v1beta1"
	searchinstall "code.alipay.com/multi-cluster/karbour/pkg/apis/search/install"
	searchv1beta1 "code.alipay.com/multi-cluster/karbour/pkg/apis/search/v1beta1"
)

var (
	// Scheme defines methods for serializing and deserializing API objects.
	Scheme = runtime.NewScheme()
	// Codecs provides methods for retrieving codecs and serializers for specific
	// versions and content types.
	Codecs = serializer.NewCodecFactory(Scheme)
	// ParameterCodec handles versioning of objects that are converted to query parameters.
	ParameterCodec = runtime.NewParameterCodec(Scheme)

	Versions = []schema.GroupVersion{
		clusterv1beta1.SchemeGroupVersion,
		searchv1beta1.SchemeGroupVersion,
	}
)

func init() {
	clusterinstall.Install(Scheme)
	searchinstall.Install(Scheme)

	// we need to add the options to empty v1
	// TODO fix the server code to avoid this
	metav1.AddToGroupVersion(Scheme, schema.GroupVersion{Version: "v1"})

	// TODO: keep the generic API server from wanting this
	unversioned := schema.GroupVersion{Group: "", Version: "v1"}
	Scheme.AddUnversionedTypes(unversioned,
		&metav1.Status{},
		&metav1.APIVersions{},
		&metav1.APIGroupList{},
		&metav1.APIGroup{},
		&metav1.APIResourceList{},
	)
}
