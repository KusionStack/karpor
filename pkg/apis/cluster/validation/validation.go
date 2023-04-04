package validation

import (
	"github.com/KusionStack/karbour/pkg/apis/cluster"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func ValidateCluster(f *cluster.Cluster) field.ErrorList {
	allErrs := field.ErrorList{}
	return allErrs
}
