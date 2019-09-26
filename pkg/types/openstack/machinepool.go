package openstack

// MachinePool stores the configuration for a machine pool installed
// on OpenStack.
type MachinePool struct {
	// FlavorName defines the OpenStack Nova flavor.
	// eg. m1.large
	FlavorName string `json:"type"`

	// RootVolume defines the root volume for instances in the machine pool.
	RootVolume *RootVolume `json:"rootVolume,omitempty"`
}

// Set sets the values from `required` to `a`.
func (o *MachinePool) Set(required *MachinePool) {
	if required == nil || o == nil {
		return
	}

	if required.FlavorName != "" {
		o.FlavorName = required.FlavorName
	}

	if required.RootVolume.Size != 0 {
		o.RootVolume.Size = required.RootVolume.Size
	}
	if required.RootVolume.Type != "" {
		o.RootVolume.Type = required.RootVolume.Type
	}
}

// RootVolume defines the storage for an instance.
type RootVolume struct {
	// Size defines the size of the volume in gibibytes (GiB).
	Size int `json:"size"`
	// Type defines the type of the volume.
	Type string `json:"type"`
}
