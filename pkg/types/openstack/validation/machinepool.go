package validation

import (
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types/openstack"
)

// ValidateMachinePool checks that the specified machine pool is valid.
func ValidateMachinePool(p *openstack.MachinePool, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	// Validate Volumes
	if p.Type == "" && p.Size > 0 {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("Type"), p.Type, "Volume type must be specified to use root volumes"))
	}
	if p.Type != "" && p.Size <= 0 {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("Size"), p.Type, "Volume size must be greater than zero to use root volumes"))
	}

	return allErrs
}
