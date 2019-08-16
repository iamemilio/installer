package validation

import (
	"errors"
	"fmt"

	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/openstack"
)

// ValidatePlatform checks that the specified platform is valid.
func ValidatePlatform(p *openstack.Platform, n *types.Networking, fldPath *field.Path, fetcher ValidValuesFetcher) field.ErrorList {
	allErrs := field.ErrorList{}
	validClouds, err := fetcher.GetCloudNames()
	if err != nil {
		allErrs = append(allErrs, field.InternalError(fldPath.Child("cloud"), errors.New("could not retrieve valid clouds")))
	} else if !isValidValue(p.Cloud, validClouds) {
		allErrs = append(allErrs, field.NotSupported(fldPath.Child("cloud"), p.Cloud, validClouds))
	} else {
		validRegions, err := fetcher.GetRegionNames(p.Cloud)
		if err != nil {
			allErrs = append(allErrs, field.InternalError(fldPath.Child("region"), errors.New("could not retrieve valid regions")))
		} else if !isValidValue(p.Region, validRegions) {
			allErrs = append(allErrs, field.NotSupported(fldPath.Child("region"), p.Region, validRegions))
		}
		validNetworkIDs, validNetworkNames, err := fetcher.GetNetworks(p.Cloud)
		fmt.Printf("Name: %s\nID: %s\n", p.ExternalNetwork, p.ExternalNetworkID)
		if err != nil {
			allErrs = append(allErrs, field.InternalError(fldPath.Child("externalNetwork"), errors.New("could not retrieve valid networks")))
		} else {
			if p.ExternalNetwork == "" && p.ExternalNetworkID == "" {
				allErrs = append(allErrs, field.NotFound(fldPath.Child("ExternalNetwork"), errors.New("No external network provided")))
			}
			if p.ExternalNetwork != "" {
				if !isValidValue(p.ExternalNetwork, validNetworkNames) {
					allErrs = append(allErrs, field.NotSupported(fldPath.Child("externalNetwork"), p.ExternalNetwork, validNetworkNames))
				} else if p.ExternalNetworkID == "" {
					p.ExternalNetworkID = IDFromName(p.ExternalNetwork, validNetworkIDs, validNetworkNames)
				}
			}
			if p.ExternalNetworkID != "" {
				if !isValidValue(p.ExternalNetworkID, validNetworkIDs) {
					allErrs = append(allErrs, field.NotSupported(fldPath.Child("externalNetworkID"), p.ExternalNetworkID, validNetworkIDs))
				} else if p.ExternalNetwork == "" {
					p.ExternalNetwork = nameFromID(p.ExternalNetworkID, validNetworkIDs, validNetworkNames)
				}
			}
			// If both are provided, then both have already been individually validated
			// just need to test to make sure they match
			if p.ExternalNetworkID != "" && p.ExternalNetwork != "" {
				if nameFromID(p.ExternalNetworkID, validNetworkIDs, validNetworkNames) != p.ExternalNetwork {
					allErrs = append(allErrs, field.NotFound(fldPath.Child("ExternalNetwork"), errors.New("External network name and ID mismatch")))
				}
			}
		}

		fmt.Printf("Name: %s\nID: %s\n", p.ExternalNetwork, p.ExternalNetworkID)

		validFlavors, err := fetcher.GetFlavorNames(p.Cloud)
		if err != nil {
			allErrs = append(allErrs, field.InternalError(fldPath.Child("computeFlavor"), errors.New("could not retrieve valid flavors")))
		} else if !isValidValue(p.FlavorName, validFlavors) {
			allErrs = append(allErrs, field.NotSupported(fldPath.Child("computeFlavor"), p.FlavorName, validFlavors))
		}
		netExts, err := fetcher.GetNetworkExtensionsAliases(p.Cloud)
		if err != nil {
			allErrs = append(allErrs, field.InternalError(fldPath.Child("trunkSupport"), errors.New("could not retrieve networking extension aliases")))
		} else {
			if isValidValue("trunk", netExts) {
				p.TrunkSupport = "1"
			} else {
				p.TrunkSupport = "0"
			}
		}
		serviceCatalog, err := fetcher.GetServiceCatalog(p.Cloud)
		if err != nil {
			allErrs = append(allErrs, field.InternalError(fldPath.Child("octaviaSupport"), errors.New("could not retrieve service catalog")))
		} else {
			if isValidValue("octavia", serviceCatalog) {
				p.OctaviaSupport = "1"
			} else {
				p.OctaviaSupport = "0"
			}
		}
	}
	if p.DefaultMachinePlatform != nil {
		allErrs = append(allErrs, ValidateMachinePool(p.DefaultMachinePlatform, fldPath.Child("defaultMachinePlatform"))...)
	}

	return allErrs
}

func isValidValue(s string, validValues []string) bool {
	for _, v := range validValues {
		if s == v {
			return true
		}
	}
	return false
}

func IDFromName(name string, ids []string, names []string) string {
	for k, v := range names {
		if name == v {
			return ids[k]
		}
	}
	return ""
}

func nameFromID(id string, ids []string, names []string) string {
	for k, v := range ids {
		if id == v {
			return names[k]
		}
	}
	return ""
}
