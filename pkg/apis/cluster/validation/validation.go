package validation

import (
	"code.alipay.com/multi-cluster/karbour/pkg/apis/cluster"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func ValidateClusterExtension(f *cluster.ClusterExtension) field.ErrorList {
	allErrs := field.ErrorList{}
	return allErrs
}
