package validation

import (
	"k8s.io/apimachinery/pkg/util/validation/field"
	"code.alipay.com/ant-iac/karbour/pkg/apis/cluster"
)

func ValidateClusterExtension(f *cluster.ClusterExtension) field.ErrorList {
	allErrs := field.ErrorList{}
	return allErrs
}
