// Copyright 2017 Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package conch

import (
	"fmt"
	"github.com/joyent/conch-shell/pkg/pgtime"
	uuid "gopkg.in/satori/go.uuid.v1"
	"net/http"
	"time"
)

// Device represents what the API docs call a "DetailedDevice"
//
// Instead of having multiple structs representing partial datasets, like the
// API chooses to do, this library will always hand back Devices. In the
// case that the API does not provide all the data, those fields will be null
// or zero values.
type Device struct {
	AssetTag        string             `json:"asset_tag"`
	BootPhase       string             `json:"boot_phase"`
	Created         pgtime.PgTime      `json:"created"`
	Deactivated     pgtime.PgTime      `json:"deactivated"`
	Graduated       pgtime.PgTime      `json:"graduated"`
	HardwareProduct uuid.UUID          `json:"hardware_product"`
	Health          string             `json:"health"`
	ID              string             `json:"id"`
	LastSeen        pgtime.PgTime      `json:"last_seen"`
	Location        DeviceLocation     `json:"location"`
	Nics            []Nic              `json:"nics"`
	Role            uuid.UUID          `json:"role"`
	State           string             `json:"state"`
	SystemUUID      uuid.UUID          `json:"system_uuid"`
	TritonUUID      uuid.UUID          `json:"triton_uuid"`
	TritonSetup     pgtime.PgTime      `json:"triton_setup"`
	Updated         pgtime.PgTime      `json:"updated"`
	UptimeSince     pgtime.PgTime      `json:"uptime_since"`
	Validated       pgtime.PgTime      `json:"validated"`
	Validations     []ValidationReport `json:"validations"`
	LatestReport    interface{}        `json:"latest_report"`
}

// ValidationReport vars provide an abstraction to make sense of the 'status'
// field in ValidationReports
const (
	ValidationReportStatusFail = 0
	ValidationReportStatusOK   = 1
)

// ValidationReport represents the result from the validation engine, comparing
// field data to expectations.
type ValidationReport struct {
	ComponentID   uuid.UUID   `json:"component_id"`
	ComponentName string      `json:"component_name"`
	ComponentType string      `json:"component_type"`
	CriteriaID    uuid.UUID   `json:"criteria_id"`
	Log           string      `json:"log"`
	Metric        interface{} `json:"metric"`
	Status        int         `json:"status"` // Can use the ValidationReportStatus consts to understand status
}

// Datacenter represents a conch datacenter, aka an AZ
type Datacenter struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

// HardwareProfileZpool represents the layout of the target device's ZFS zpools
type HardwareProfileZpool struct {
	ID       uuid.UUID `json:"id"`
	Cache    int       `json:"cache"`
	Log      int       `json:"log"`
	DisksPer int       `json:"disks_per"`
	Spare    int       `json:"spare"`
	VdevN    int       `json:"vdev_n"`
	VdevT    string    `json:"vdev_t"`
}

// HardwareProfile is a detailed accounting of either the actual hardware or
// intended hardware configuration of a Device, depending on the API endpoint
// in question
type HardwareProfile struct {
	ID           uuid.UUID            `json:"id"`
	NumNics      int                  `json:"nics_num"`
	BiosFirmware string               `json:"bios_firmware"`
	NumCPU       int                  `json:"cpu_num"`
	CPUType      string               `json:"cpu_type"`
	NumDimms     int                  `json:"dimms_num"`
	TotalPSU     int                  `json:"psu_total"`
	Purpose      string               `json:"purpose"`
	TotalRAM     int                  `json:"ram_total"`
	SASNum       int                  `json:"sas_num"`
	SizeSAS      int                  `json:"sas_size"`
	SlotsSAS     string               `json:"saas_slots"`
	NumSATA      int                  `json:"sata_num"`
	SizeSATA     int                  `json:"sata_size"`
	SlotsSATA    string               `json:"sata_slots"`
	NumSSD       int                  `json:"ssd_num"`
	SizeSSD      int                  `json:"ssd_size"`
	SlotsSSD     string               `json:"ssd_slots"`
	NumUSB       int                  `json:"usb_num"`
	Zpool        HardwareProfileZpool `json:"zpool"`
}

// HardwareProduct is a type of Device. For instance, "Hallasan C"
type HardwareProduct struct {
	ID                uuid.UUID       `json:"id"`
	Name              string          `json:"name"`
	Alias             string          `json:"alias"`
	Prefix            string          `json:"prefix"`
	Profile           HardwareProfile `json:"profile"`
	Vendor            string          `json:"vendor"`
	Specification     interface{}     `json:"specification"`
	SKU               string          `json:"sku"`
	GenerationName    string          `json:"generation_name"`
	LegacyProductName string          `json:"legacy_product_name"`
}

// HardwareProductTarget represents the HardwareProduct that a device should
// have based on its location
type HardwareProductTarget struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Alias string    `json:"alias"`
}

// Rack represents a physical rack
type Rack struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	Role       string    `json:"role"`
	Unit       int       `json:"unit"`
	Size       int       `json:"size"`
	Datacenter string    `json:"datacenter"`
	// The key of the Slots map is the RU slot number
	Slots        map[int]RackSlot `json:"slots"`
	SerialNumber string           `json:"serial_number"`
	AssetTag     string           `json:"asset_tag"`
}

// RackSlot represents a physical slot in a physical Rack
type RackSlot struct {
	ID       uuid.UUID `json:"id"`
	Size     int       `json:"size"`
	Name     string    `json:"name"`
	Alias    string    `json:"alias"`
	Vendor   string    `json:"vendor"`
	Occupant Device    `json:"occupant"`
}

// Nic is a network interface card, including its peer switch info
type Nic struct {
	MAC         string `json:"mac"`
	IfaceName   string `json:"iface_name"`
	IfaceVendor string `json:"iface_vendor"`
	IfaceType   string `json:"iface_type"`
	PeerMac     string `json:"peer_mac"`
	PeerPort    string `json:"peer_port"`
	PeerSwitch  string `json:"peer_switch"`
}

// DeviceLocation represents the location of a device, including its datacenter
// and rack
type DeviceLocation struct {
	Datacenter            Datacenter            `json:"datacenter"`
	Rack                  Rack                  `json:"rack"`
	TargetHardwareProduct HardwareProductTarget `json:"target_hardware_product"`
}

// DeviceService represents a single device service, usually a piece of
// software like Manta
type DeviceService struct {
	ID      uuid.UUID `json:"id"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
	Name    string    `json:"name"`
}

// DeviceRole represents a device role, a combination of a hardware product and
// a list of services
type DeviceRole struct {
	ID                uuid.UUID   `json:"id"`
	Created           time.Time   `json:"created"`
	Updated           time.Time   `json:"updated"`
	Description       string      `json:"description"`
	HardwareProductID uuid.UUID   `json:"hardware_product_id"`
	Services          []uuid.UUID `json:"services"`
}

// GetWorkspaceDevices retrieves a list of Devices for the given
// workspace.
// Pass true for 'IDsOnly' to get Devices with only the ID field populated
// Pass a string for 'graduated' to filter devices by graduated value, as per https://conch.joyent.us/doc#getdevices
// Pass a string for 'health' to filter devices by health value, as per https://conch.joyent.us/doc#getdevices
func (c *Conch) GetWorkspaceDevices(workspaceUUID fmt.Stringer, idsOnly bool, graduated string, health string) ([]Device, error) {

	devices := make([]Device, 0)

	opts := struct {
		IDsOnly   bool   `url:"ids_only,omitempty"`
		Graduated string `url:"graduated,omitempty"`
		Health    string `url:"health,omitempty"`
	}{
		idsOnly,
		graduated,
		health,
	}

	aerr := &APIError{}

	url := "/workspace/" + workspaceUUID.String() + "/device"
	if idsOnly {
		ids := make([]string, 0)

		res, err := c.sling().New().
			Get(url).
			QueryStruct(opts).
			Receive(&ids, aerr)

		cerr := c.isHTTPResOk(res, err, aerr)

		if cerr != nil {
			return devices, cerr
		}

		for _, v := range ids {
			device := Device{ID: v}
			devices = append(devices, device)
		}
		return devices, cerr
	}

	res, err := c.sling().New().
		Get(url).
		QueryStruct(opts).
		Receive(&devices, aerr)

	return devices, c.isHTTPResOk(res, err, aerr)
}

// GetDevice returns a Device given a specific serial/id
func (c *Conch) GetDevice(serial string) (Device, error) {
	var device Device
	device.ID = serial

	return c.FillInDevice(device)
}

// FillInDevice takes an existing device and fills in its data using "/device"
//
// This exists because the API hands back partial devices in most cases. It's
// likely, though, that any client utility will eventually want all the data
// about a device and not just bits
func (c *Conch) FillInDevice(d Device) (Device, error) {
	aerr := &APIError{}
	res, err := c.sling().New().Get("/device/"+d.ID).Receive(&d, aerr)
	return d, c.isHTTPResOk(res, err, aerr)
}

// GetDeviceLocation fetches the location for a device, via
// /device/:serial/location
func (c *Conch) GetDeviceLocation(serial string) (DeviceLocation, error) {
	var location DeviceLocation

	aerr := &APIError{}
	res, err := c.sling().New().
		Get("/device/"+serial+"/location").
		Receive(&location, aerr)

	return location, c.isHTTPResOk(res, err, aerr)
}

// GetWorkspaceRacks fetchest the list of racks for a workspace, via
// /workspace/:uuid/rack
//
// NOTE: The API currently returns a hash of arrays where teh key is the
// datacenter/az. This routine copies that key into the Datacenter field in the
// Rack struct.
func (c *Conch) GetWorkspaceRacks(workspaceUUID fmt.Stringer) ([]Rack, error) {
	racks := make([]Rack, 0)
	j := make(map[string][]Rack)

	aerr := &APIError{}
	res, err := c.sling().New().
		Get("/workspace/"+workspaceUUID.String()+"/rack").
		Receive(&j, aerr)

	for az, loc := range j {
		for _, rack := range loc {
			rack.Datacenter = az
			racks = append(racks, rack)
		}
	}

	return racks, c.isHTTPResOk(res, err, aerr)
}

// GetWorkspaceRack fetches a single rack for a workspace, via
// /workspace/:uuid/rack/:id
func (c *Conch) GetWorkspaceRack(workspaceUUID fmt.Stringer, rackUUID fmt.Stringer) (Rack, error) {
	var rack Rack

	aerr := &APIError{}
	res, err := c.sling().New().
		Get("/workspace/"+workspaceUUID.String()+"/rack/"+rackUUID.String()).
		Receive(&rack, aerr)

	return rack, c.isHTTPResOk(res, err, aerr)
}

// GetHardwareProduct fetches a single hardware product via
// /hardware/product/:uuid
func (c *Conch) GetHardwareProduct(hardwareProductUUID fmt.Stringer) (HardwareProduct, error) {
	var prod HardwareProduct

	aerr := &APIError{}
	res, err := c.sling().New().
		Get("/hardware_product/"+hardwareProductUUID.String()).
		Receive(&prod, aerr)

	return prod, c.isHTTPResOk(res, err, aerr)
}

// GetHardwareProducts fetches a single hardware product via
// /hardware_product
func (c *Conch) GetHardwareProducts() ([]HardwareProduct, error) {
	prods := make([]HardwareProduct, 0)

	aerr := &APIError{}
	res, err := c.sling().New().
		Get("/hardware_product").
		Receive(&prods, aerr)

	return prods, c.isHTTPResOk(res, err, aerr)
}

// GraduateDevice sets the 'graduated' field for the given device, via
// /device/:serial/graduate
// WARNING: This is a one way operation and cannot currently be undone via the
// API
func (c *Conch) GraduateDevice(serial string) error {
	aerr := &APIError{}
	res, err := c.sling().New().
		Post("/device/"+serial+"/graduate").
		Receive(nil, aerr)

	return c.isHTTPResOk(res, err, aerr)
}

// DeviceTritonReboot sets the 'triton_reboot' field for the given device, via
// /device/:serial/triton_reboot
// WARNING: This is a one way operation and cannot currently be undone via the
// API
func (c *Conch) DeviceTritonReboot(serial string) error {
	aerr := &APIError{}
	res, err := c.sling().New().
		Post("/device/"+serial+"/triton_reboot").
		Receive(nil, aerr)

	return c.isHTTPResOk(res, err, aerr)
}

// SetDeviceTritonUUID sets the triton UUID via /device/:serial/triton_uuid
func (c *Conch) SetDeviceTritonUUID(serial string, id uuid.UUID) error {
	j := struct {
		TritonUUID string `json:"triton_uuid"`
	}{
		id.String(),
	}

	aerr := &APIError{}
	res, err := c.sling().New().
		Post("/device/"+serial+"/triton_uuid").
		BodyJSON(j).
		Receive(nil, aerr)

	return c.isHTTPResOk(res, err, aerr)
}

// MarkDeviceTritonSetup marks the device as setup for Triton
// For this action to succeed, the device must have its Triton UUID set and
// marked as rebooted into Triton. If these conditions are not met, this
// function will return ErrBadInput
func (c *Conch) MarkDeviceTritonSetup(serial string) error {

	aerr := &APIError{}
	res, err := c.sling().New().
		Post("/device/"+serial+"/triton_setup").
		Receive(nil, aerr)

	return c.isHTTPResOk(res, err, aerr)
}

// SetDeviceAssetTag sets the asset tag for the provided serial
func (c *Conch) SetDeviceAssetTag(serial string, tag string) error {
	j := struct {
		AssetTag string `json:"asset_tag"`
	}{
		tag,
	}

	aerr := &APIError{}
	res, err := c.sling().New().
		Post("/device/"+serial+"/asset_tag").
		BodyJSON(j).
		Receive(nil, aerr)
	return c.isHTTPResOk(res, err, aerr)
}

// GetDeviceServices gets a list of all active device services
func (c *Conch) GetDeviceServices() ([]DeviceService, error) {
	s := make([]DeviceService, 0)

	aerr := &APIError{}
	res, err := c.sling().New().
		Get("/device/service").
		Receive(&s, aerr)

	return s, c.isHTTPResOk(res, err, aerr)
}

// GetDeviceService gets a single device service
func (c *Conch) GetDeviceService(id fmt.Stringer) (DeviceService, error) {
	var service DeviceService

	aerr := &APIError{}
	res, err := c.sling().New().
		Get("/device/service/"+id.String()).
		Receive(&service, aerr)

	return service, c.isHTTPResOk(res, err, aerr)
}

// GetDeviceServiceByName gets a single device service by its name
func (c *Conch) GetDeviceServiceByName(name string) (DeviceService, error) {
	var service DeviceService

	aerr := &APIError{}
	res, err := c.sling().New().
		Get("/device/service/name="+name).
		Receive(&service, aerr)

	return service, c.isHTTPResOk(res, err, aerr)
}

// SaveDeviceService creates or updates a device service. A new service is
// created if the structure lacks an ID. Otherwise, an update is attempted.
//
func (c *Conch) SaveDeviceService(s *DeviceService) error {
	if s.Name == "" {
		return ErrBadInput
	}
	var err error
	var res *http.Response
	aerr := &APIError{}

	if uuid.Equal(s.ID, uuid.UUID{}) {
		j := struct {
			Name string `json:"name"`
		}{
			s.Name,
		}
		res, err = c.sling().New().
			Post("/device/service").
			BodyJSON(j).
			Receive(&s, aerr)
	} else {
		j := struct {
			Name string `json:"name"`
		}{
			s.Name,
		}
		res, err = c.sling().New().
			Post("/device/service/"+s.ID.String()).
			BodyJSON(j).
			Receive(&s, aerr)
	}

	return c.isHTTPResOk(res, err, aerr)
}

// DeleteDeviceService deletes a device service, given a UUID
func (c *Conch) DeleteDeviceService(id fmt.Stringer) error {
	aerr := &APIError{}

	res, err := c.sling().New().
		Delete("/device/service/"+id.String()).Receive(nil, aerr)

	return c.isHTTPResOk(res, err, aerr)
}

// GetDeviceRoles returns a list of all active device roles
func (c *Conch) GetDeviceRoles() ([]DeviceRole, error) {
	d := make([]DeviceRole, 0)

	aerr := &APIError{}
	res, err := c.sling().New().
		Get("/device/role").
		Receive(&d, aerr)

	return d, c.isHTTPResOk(res, err, aerr)
}

// GetDeviceRole gets a single device role, given its uuid
func (c *Conch) GetDeviceRole(id fmt.Stringer) (DeviceRole, error) {
	var role DeviceRole

	aerr := &APIError{}
	res, err := c.sling().New().
		Get("/device/role/"+id.String()).
		Receive(&role, aerr)
	return role, c.isHTTPResOk(res, err, aerr)
}

// SaveDeviceRole creates or updates a device role. A new role is
// created if the structure lacks an ID. Otherwise, an update is attempted.
func (c *Conch) SaveDeviceRole(r *DeviceRole) error {
	if uuid.Equal(r.HardwareProductID, uuid.UUID{}) {
		return ErrBadInput
	}

	var err error
	var res *http.Response
	aerr := &APIError{}

	j := struct {
		Description       string    `json:"description"`
		HardwareProductID uuid.UUID `json:"hardware_product_id"`
	}{
		r.Description,
		r.HardwareProductID,
	}

	if uuid.Equal(r.ID, uuid.UUID{}) {
		res, err = c.sling().New().
			Post("/device/role").
			BodyJSON(j).
			Receive(&r, aerr)
	} else {
		res, err = c.sling().New().
			Post("/device/role/"+r.ID.String()).
			BodyJSON(j).
			Receive(&r, aerr)
	}

	return c.isHTTPResOk(res, err, aerr)
}

//DeleteDeviceRole deletes the device role for the given uuid
func (c *Conch) DeleteDeviceRole(id fmt.Stringer) error {

	aerr := &APIError{}
	res, err := c.sling().New().
		Delete("/device/role/"+id.String()).Receive(nil, aerr)

	return c.isHTTPResOk(res, err, aerr)
}

// AddServiceToDeviceRole removes the given service from the given device role
func (c *Conch) AddServiceToDeviceRole(d DeviceRole, ds DeviceService) error {
	if uuid.Equal(d.ID, uuid.UUID{}) {
		return ErrBadInput
	}
	if uuid.Equal(ds.ID, uuid.UUID{}) {
		return ErrBadInput
	}

	var err error
	var res *http.Response
	aerr := &APIError{}

	j := struct {
		Service uuid.UUID `json:"service"`
	}{
		ds.ID,
	}

	res, err = c.sling().New().
		Post("/device/role/"+d.ID.String()+"/add_service").
		BodyJSON(j).
		Receive(nil, aerr)

	return c.isHTTPResOk(res, err, aerr)
}

// RemoveServiceFromDeviceRole removes the given service from the given role
func (c *Conch) RemoveServiceFromDeviceRole(d DeviceRole, ds DeviceService) error {
	if uuid.Equal(d.ID, uuid.UUID{}) {
		return ErrBadInput
	}
	if uuid.Equal(ds.ID, uuid.UUID{}) {
		return ErrBadInput
	}

	var err error
	var res *http.Response
	aerr := &APIError{}

	j := struct {
		Service uuid.UUID `json:"service"`
	}{
		ds.ID,
	}

	res, err = c.sling().New().
		Post("/device/role/"+d.ID.String()+"/remove_service").
		BodyJSON(j).
		Receive(nil, aerr)
	return c.isHTTPResOk(res, err, aerr)
}

// GetDeviceIPMI retrieves "/device/:serial/interface/impi1/ipaddr"
func (c *Conch) GetDeviceIPMI(serial string) (string, error) {
	j := make(map[string]string)

	aerr := &APIError{}
	res, err := c.sling().New().
		Get("/device/"+serial+"/interface/ipmi1/ipaddr").Receive(&j, aerr)

	if herr := c.isHTTPResOk(res, err, aerr); herr != nil {
		return "", herr
	}

	return j["ipaddr"], nil
}
