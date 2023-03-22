package validation

import (
	"code.alipay.com/multi-cluster/karbour/pkg/apis/search"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func ValidateSearchExtension(f *search.SearchExtension) field.ErrorList {
	allErrs := field.ErrorList{}
	return allErrs
}
