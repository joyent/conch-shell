// Copyright Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package conch

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/joyent/conch-shell/pkg/conch/uuid"
)

// ValidationReport vars provide an abstraction to make sense of the 'status'
// field in ValidationReports
const (
	ValidationReportStatusFail = 0
	ValidationReportStatusOK   = 1
)

// Conch contains auth and configuration data
type Conch struct {
	BaseURL   string
	UserAgent map[string]string
	Debug     bool
	Trace     bool
	Token     string

	HTTPClient *http.Client
}

type Datacenter struct {
	ID         uuid.UUID `json:"id"`
	Vendor     string    `json:"vendor"`
	VendorName string    `json:"vendor_name"`
	Region     string    `json:"region"`
	Location   string    `json:"location"`
	Created    time.Time `json:"created"`
	Updated    time.Time `json:"updated"`
}

type DatacenterDetailedRoom struct {
	ID           uuid.UUID `json:"id"`
	AZ           string    `json:"az"`
	Alias        string    `json:"alias"`
	VendorName   string    `json:"vendor_name,omitempty"`
	DatacenterID uuid.UUID `json:"datacenter'`
	Created      time.Time `json:"created"`
	Updated      time.Time `json:"updated"`
}

// Device represents what the API docs call a "DetailedDevice"
//
// Instead of having multiple structs representing partial datasets, like the
// API chooses to do, this library will always hand back Devices. In the
// case that the API does not provide all the data, those fields will be null
// or zero values.
type Device struct {
	AssetTag              string             `json:"asset_tag"`
	Created               time.Time          `json:"created"`
	Deactivated           time.Time          `json:"deactivated"`
	Graduated             time.Time          `json:"graduated"`
	HardwareProduct       uuid.UUID          `json:"hardware_product"`
	Health                string             `json:"health"`
	Hostname              string             `json:"hostname"`
	ID                    string             `json:"id"`
	LastSeen              time.Time          `json:"last_seen"`
	Location              DeviceLocation     `json:"location"`
	Nics                  []Nic              `json:"nics"`
	State                 string             `json:"state"`
	SystemUUID            uuid.UUID          `json:"system_uuid"`
	TritonUUID            uuid.UUID          `json:"triton_uuid"`
	TritonSetup           time.Time          `json:"triton_setup"`
	Updated               time.Time          `json:"updated"`
	UptimeSince           time.Time          `json:"uptime_since"`
	Validated             time.Time          `json:"validated"`
	Validations           []ValidationReport `json:"validations"`
	LatestReport          interface{}        `json:"latest_report"`
	LatestReportIsInvalid bool               `json:"latest_report_is_invalid"`
	InvalidReport         string             `json:"invalid_report"`
	Disks                 []Disk             `json:"disks"`
	RackUnitStart         int                `json:"rack_unit_start`
	RackID                uuid.UUID          `json:"rack_id"`
	Phase                 string             `json:"phase"`
}

type Devices []Device

func (d Devices) Len() int {
	return len(d)
}

func (d Devices) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}

func (d Devices) Less(i, j int) bool {
	return d[i].ID < d[j].ID
}

// DeviceDisk ...
type Disk struct {
	ID           uuid.UUID `json:"id"`
	Created      time.Time `json:"created"`
	Updated      time.Time `json:"updated"`
	DriveType    string    `json:"drive_type"`
	Enclosure    string    `json:"enclosure"`
	Firmware     string    `json:"firmware"`
	HBA          string    `json:"hba"`
	Health       string    `json:"health"`
	Model        string    `json:"model"`
	SerialNumber string    `json:"serial_number"`
	Size         int       `json:"size"`
	Slot         int       `json:"slot"`
	Temp         int       `json:"temp"`
	Transport    string    `json:"transport"`
	Vendor       string    `json:"vendor"`
}

// DeviceLocation represents the location of a device, including its
// datacenter, room and rack
type DeviceLocation struct {
	Datacenter            Datacenter             `json:"datacenter"`
	Room                  DatacenterDetailedRoom `json:"datacenter_room"`
	Rack                  Rack                   `json:"rack"`
	TargetHardwareProduct HardwareProductTarget  `json:"target_hardware_product"`
	RackUnitStart         int                    `json:"rack_unit_start"`
}

type ExtendedDevice struct {
	Device
	IPMI          string                    `json:"ipmi"`
	HardwareName  string                    `json:"hardware_name"`
	RackRole      RackRole                  `json:"rack_role"`
	SKU           string                    `json:"sku"`
	Enclosures    map[string]map[int]Disk   `json:"enclosures"`
	IsGraduated   bool                      `json:"is_graduated"`
	IsTritonSetup bool                      `json:"is_triton_setup"`
	IsValidated   bool                      `json:"is_validated"`
	Validations   []ValidationPlanExecution `json:"validations"`
}

type Rack struct {
	ID               uuid.UUID `json:"id"`
	Created          time.Time `json:"created"`
	Updated          time.Time `json:"updated"`
	DatacenterRoomID uuid.UUID `json:"datacenter_room_id"`
	Name             string    `json:"name"`
	RoleID           uuid.UUID `json:"role"`
	SerialNumber     string    `json:"serial_number"`
	AssetTag         string    `json:"asset_tag"`
	Phase            string    `json:"phase"`
}

type RackLayoutSlot struct {
	ID        uuid.UUID `json:"id"`
	Created   time.Time `json:"created"`
	Updated   time.Time `json:"updated"`
	RackID    uuid.UUID `json:"rack_id"`
	ProductID uuid.UUID `json:"product_id"`
	RUStart   int       `json:"ru_start"`
}

type RackLayoutSlots []RackLayoutSlot

func (g RackLayoutSlots) Len() int {
	return len(g)
}
func (g RackLayoutSlots) Swap(i, j int) {
	g[i], g[j] = g[j], g[i]
}

func (g RackLayoutSlots) Less(i, j int) bool {
	return g[i].RUStart > g[j].RUStart
}

type RackRole struct {
	ID       uuid.UUID `json:"id"`
	Created  time.Time `json:"created"`
	Updated  time.Time `json:"updated"`
	Name     string    `json:"name"`
	RackSize int       `json:"rack_size"`
}

type Room struct {
	ID           uuid.UUID `json:"id"`
	Created      time.Time `json:"created"`
	Updated      time.Time `json:"updated"`
	DatacenterID uuid.UUID `json:"datacenter"`
	AZ           string    `json:"az"`
	Alias        string    `json:"alias"`
	VendorName   string    `json:"vendor_name"`
}

// HardwareProduct is a type of Device. For instance, "Hallasan C"
type HardwareProduct struct {
	ID                uuid.UUID       `json:"id"`
	Name              string          `json:"name"`
	Alias             string          `json:"alias"`
	Prefix            string          `json:"prefix,omitempty"`
	HardwareVendorID  uuid.UUID       `json:"hardware_vendor_id"`
	GenerationName    string          `json:"generation_name,omitempty"`
	LegacyProductName string          `json:"legacy_product_name,omitempty"`
	SKU               string          `json:"sku,omitempty"`
	Specification     interface{}     `json:"specification"`
	Profile           HardwareProfile `json:"hardware_product_profile"`
	Created           time.Time       `json:"created"`
	Updated           time.Time       `json:"updated"`
}

func (h *HardwareProduct) UnmarshalJSON(data []byte) error {
	r := struct {
		ID                uuid.UUID       `json:"id"`
		Name              string          `json:"name"`
		Alias             string          `json:"alias"`
		Prefix            string          `json:"prefix"`
		HardwareVendorID  uuid.UUID       `json:"hardware_vendor_id"`
		GenerationName    string          `json:"generation_name"`
		LegacyProductName string          `json:"legacy_product_name"`
		SKU               string          `json:"sku"`
		Specification     string          `json:"specification"`
		Profile           HardwareProfile `json:"hardware_product_profile"`
		Created           time.Time       `json:"created"`
		Updated           time.Time       `json:"updated"`
	}{}

	if err := json.Unmarshal(data, &r); err != nil {
		return err
	}

	h.ID = r.ID
	h.Name = r.Name
	h.Alias = r.Alias
	h.Prefix = r.Prefix
	h.HardwareVendorID = r.HardwareVendorID
	h.GenerationName = r.GenerationName
	h.LegacyProductName = r.LegacyProductName
	h.SKU = r.SKU
	h.Profile = r.Profile
	h.Created = r.Created
	h.Updated = r.Updated

	if r.Specification == "" {
		h.Specification = make(map[string]interface{})
	} else {
		var s interface{}
		if err := json.Unmarshal([]byte(r.Specification), &s); err != nil {
			return err
		}
		h.Specification = s
	}
	return nil
}

// HardwareProductTarget represents the HardwareProduct that a device should
// have based on its location
type HardwareProductTarget struct {
	ID     uuid.UUID `json:"id"`
	Name   string    `json:"name"`
	Alias  string    `json:"alias"`
	Vendor string    `json:"vendor"`
}

// HardwareProfile is a detailed accounting of either the actual hardware or
// intended hardware configuration of a Device, depending on the API endpoint
// in question
type HardwareProfile struct {
	ID           uuid.UUID `json:"id"`
	BiosFirmware string    `json:"bios_firmware"`
	CPUType      string    `json:"cpu_type"`
	HbaFirmware  string    `json:"hba_firmware,omitempty"`
	NumCPU       int       `json:"cpu_num"`
	NumDimms     int       `json:"dimms_num"`
	NumNics      int       `json:"nics_num"`
	NumUSB       int       `json:"usb_num"`
	Purpose      string    `json:"purpose"`

	SasHddNum   int    `json:"sas_hdd_num,omitempty"`
	SasHddSize  int    `json:"sas_hdd_size,omitempty"`
	SasHddSlots string `json:"sas_hdd_slots,omitempty"`

	SataHddNum   int    `json:"sata_hdd_num,omitempty"`
	SataHddSize  int    `json:"sata_hdd_size,omitempty"`
	SataHddSlots string `json:"sata_hdd_slots,omitempty"`

	SataSsdNum   int    `json:"sata_ssd_num,omitempty"`
	SataSsdSize  int    `json:"sata_ssd_size,omitempty"`
	SataSsdSlots string `json:"sata_ssd_slots,omitempty"`

	NvmeSsdNum   int    `json:"nvme_ssd_num,omitempty"`
	NvmeSsdSize  int    `json:"nvme_ssd_size,omitempty"`
	NvmeSsdSlots string `json:"nvme_ssd_slots,omitempty"`
	RaidLunNum   int    `json:"raid_lun_num,omitempty"`
	TotalPSU     int    `json:"psu_total,omitempty"`
	TotalRAM     int    `json:"ram_total"`
	RackUnit     int    `json:"rack_unit"`
}

// HardwareVendor ...
type HardwareVendor struct {
	ID      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
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

type WorkspaceRack struct {
	ID           uuid.UUID          `json:"id"`
	Name         string             `json:"name"`
	Role         string             `json:"role"`
	Unit         int                `json:"unit"` // BUG(sungo): This exists because device locations provide rack info, but also slot info. This is a sloppy combination of data streams
	Size         int                `json:"size"`
	Datacenter   string             `json:"datacenter"`
	Slots        WorkspaceRackSlots `json:"slots,omitempty"`
	SerialNumber string             `json:"serial_number"`
	AssetTag     string             `json:"asset_tag"`
	Phase        string             `json:"phase"`
}

type WorkspaceRackSlot struct {
	ID            uuid.UUID `json:"id"`
	Size          int       `json:"size"`
	Name          string    `json:"name"`
	Alias         string    `json:"alias"`
	Vendor        string    `json:"vendor"`
	Occupant      Device    `json:"occupant"`
	RackUnitStart int       `json:"rack_unit_start"`
}

type WorkspaceRackSlots []WorkspaceRackSlot

func (r WorkspaceRackSlots) Len() int {
	return len(r)
}
func (r WorkspaceRackSlots) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func (r WorkspaceRackSlots) Less(i, j int) bool {
	return r[i].RackUnitStart > r[j].RackUnitStart
}

// User represents a person able to access the Conch API or UI
type User struct {
	ID      string    `json:"id,omitempty"`
	Email   string    `json:"email"`
	Name    string    `json:"name"`
	Role    string    `json:"role"`
	RoleVia uuid.UUID `json:"role_via,omitempty"`
}

type Users []User

func (u Users) Len() int {
	return len(u)
}
func (u Users) Swap(i, j int) {
	u[i], u[j] = u[j], u[i]
}

func (u Users) Less(i, j int) bool {
	return strings.ToLower(u[i].Name) < strings.ToLower(u[j].Name)
}

// UserDetailed ...
type UserDetailed struct {
	ID                  uuid.UUID          `json:"id"`
	Name                string             `json:"name"`
	Email               string             `json:"email"`
	Created             time.Time          `json:"created"`
	LastLogin           time.Time          `json:"last_login"`
	RefuseSessionAuth   bool               `json:"refuse_session_auth"`
	ForcePasswordChange bool               `json:"force_password_change"`
	Workspaces          WorkspacesAndRoles `json:"workspaces,omitempty"`
	IsAdmin             bool               `json:"is_admin"`
}

type UsersDetailed []UserDetailed

func (u UsersDetailed) Len() int {
	return len(u)
}
func (u UsersDetailed) Swap(i, j int) {
	u[i], u[j] = u[j], u[i]
}

func (u UsersDetailed) Less(i, j int) bool {
	return strings.ToLower(u[i].Name) < strings.ToLower(u[j].Name)
}

type UserProfile struct {
	Created             time.Time          `json:"created"`
	Email               string             `json:"email"`
	ForcePasswordChange bool               `json:"force_password_change"`
	ID                  uuid.UUID          `json:"id"`
	LastLogin           time.Time          `json:"last_login"`
	Name                string             `json:"name"`
	RefuseSessionAuth   bool               `json:"refuse_session_auth"`
	Workspaces          WorkspacesAndRoles `json:"workspaces"`
}

// Validation represents device validations loaded into Conch
type Validation struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Version     int       `json:"version"`
	Description string    `json:"description"`
	Created     time.Time `json:"created"`
	Updated     time.Time `json:"updated"`
	Deactivated time.Time `json:"deactivated"`
}

type Validations []Validation

func (v Validations) Len() int {
	return len(v)
}
func (v Validations) Swap(i, j int) {
	v[i], v[j] = v[j], v[i]
}

func (v Validations) Less(i, j int) bool {
	return strings.ToLower(v[i].Name) < strings.ToLower(v[j].Name)
}

// ValidationPlan represents an organized association of Validations
type ValidationPlan struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Created     time.Time `json:"created"`
}

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

type ValidationPlanExecution struct {
	ID          uuid.UUID      `json:"id"`
	Name        string         `json:"name"`
	Validations ValidationRuns `json:"validations"`
}

// ValidationResult is a result of running a validation on a device
type ValidationResult struct {
	ID              uuid.UUID `json:"id"`
	Category        string    `json:"category"`
	ComponentID     string    `json:"component_id"`
	DeviceID        string    `json:"device_id"`
	HardwareProduct uuid.UUID `json:"hardware_product_id"`
	Hint            string    `json:"hint"`
	Message         string    `json:"message"`
	Status          string    `json:"status"`
	ValidationID    uuid.UUID `json:"validation_id"`
}

type ValidationRun struct {
	ID      uuid.UUID          `json:"id"`
	Name    string             `json:"name"`
	Passed  bool               `json:"passed"`
	Results []ValidationResult `json:"results"`
}

type ValidationRuns []ValidationRun

func (v ValidationRuns) Len() int {
	return len(v)
}

func (v ValidationRuns) Swap(i, j int) {
	v[i], v[j] = v[j], v[i]
}

func (v ValidationRuns) Less(i, j int) bool {
	return v[i].Name < v[j].Name
}

// ValidationState is the result of running a validation plan on a device
type ValidationState struct {
	ID               uuid.UUID          `json:"id"`
	Created          time.Time          `json:"created"`
	Completed        time.Time          `json:"completed"`
	DeviceID         string             `json:"device_id"`
	Results          []ValidationResult `json:"results"`
	Status           string             `json:"status"`
	ValidationPlanID uuid.UUID          `json:"validation_plan_id"`
}

// WorkspaceRelays ...
type WorkspaceRelays []WorkspaceRelay

func (w WorkspaceRelays) Len() int {
	return len(w)
}

func (w WorkspaceRelays) Swap(i, j int) {
	w[i], w[j] = w[j], w[i]
}

func (w WorkspaceRelays) Less(i, j int) bool {
	return w[i].Updated.Before(w[j].Updated)
}

// WorkspaceRelay represents a Conch Relay unit, a physical piece of hardware that
// mediates Livesys interactions in the field
type WorkspaceRelay struct {
	ID         string                 `json:"id"` // *not* a UUID
	Created    time.Time              `json:"created"`
	Updated    time.Time              `json:"updated"`
	Alias      string                 `json:"alias"`
	IPAddr     string                 `json:"ipaddr"`
	SSHPort    int                    `json:"ssh_port"`
	Version    string                 `json:"version"`
	LastSeen   time.Time              `json:"last_seen"`
	NumDevices int                    `json:"num_devices"`
	Location   WorkspaceRelayLocation `json:"location"`
}

// WorkspaceRelayLocation ...
type WorkspaceRelayLocation struct {
	Az            string    `json:"az"`
	RackID        uuid.UUID `json:"rack_id"`
	RackName      string    `json:"rack_name"`
	RackUnitStart int       `json:"rack_unit_start"`
	RoleName      string    `json:"role_name"`
}

// Workspace represents a Conch data partition which allows users to create
// custom lists of hardware
type Workspace struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Role        string    `json:"role"`
	ParentID    uuid.UUID `json:"parent_id,omitempty"`
}

type Workspaces []Workspace

func (w Workspaces) Len() int {
	return len(w)
}

func (w Workspaces) Swap(i, j int) {
	w[i], w[j] = w[j], w[i]
}

func (w Workspaces) Less(i, j int) bool {
	return w[i].Name < w[j].Name
}

// WorkspaceAndRole ...
type WorkspaceAndRole struct {
	Workspace
	RoleVia uuid.UUID `json:"role_via"`
}

type WorkspacesAndRoles []WorkspaceAndRole

func (w WorkspacesAndRoles) Len() int {
	return len(w)
}

func (w WorkspacesAndRoles) Swap(i, j int) {
	w[i], w[j] = w[j], w[i]
}

func (w WorkspacesAndRoles) Less(i, j int) bool {
	return w[i].Name < w[j].Name
}

// WorkspaceUser ...
type WorkspaceUser struct {
	User
	RoleVia uuid.UUID `json:"role_via,omitempty"`
}

/* This is a piece of fun from /workspace/:id/rack/:id/layout
The payload looks like:
{ "my-device-id": 47 }

Where '47' is the rack unit start for the device
*/
type WorkspaceRackLayoutAssignments map[string]int

// corresponds to conch.git/json-schema/input.yaml;NewUserToken
type CreateNewUserToken struct {
	Name string `json:"name"`
}

// corresponds to conch.git/json-schema/response.yaml;UserToken
type UserToken struct {
	Name     string    `json:"name"`
	Created  time.Time `json:"created"`
	LastUsed time.Time `json:"last_used,omitempty"`
	Expires  time.Time `json:"expires"`
}

type UserTokens []UserToken

func (u UserTokens) Len() int {
	return len(u)
}

func (u UserTokens) Swap(i, j int) {
	u[i], u[j] = u[j], u[i]
}

func (u UserTokens) Less(i, j int) bool {
	return u[i].Name < u[j].Name
}

// corresponds to conch.git/json-schema/response.yaml;NewUserToken
type NewUserToken struct {
	UserToken
	Token string `json:"token"`
}

/**/

type RequestRackAssignmentUpdate struct {
	DeviceID       string `json:"device_id"`
	RackUnitStart  int    `json:"rack_unit_start"`
	DeviceAssetTag string `json:"device_asset_tag,omitempty"`
}

type RequestRackAssignmentUpdates []RequestRackAssignmentUpdate

func (r RequestRackAssignmentUpdates) Len() int {
	return len(r)
}

func (r RequestRackAssignmentUpdates) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func (r RequestRackAssignmentUpdates) Less(i, j int) bool {
	return r[i].RackUnitStart < r[j].RackUnitStart
}

/**/
type RequestRackAssignmentDelete struct {
	DeviceID      string `json:"device_id"`
	RackUnitStart int    `json:"rack_unit_start"`
}

type RequestRackAssignmentDeletes []RequestRackAssignmentDelete

func (r RequestRackAssignmentDeletes) Len() int {
	return len(r)
}

func (r RequestRackAssignmentDeletes) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func (r RequestRackAssignmentDeletes) Less(i, j int) bool {
	return r[i].RackUnitStart < r[j].RackUnitStart
}

/**/

type ResponseRackAssignment struct {
	DeviceID        string `json:"device_id,omitempty"`
	DeviceAssetTag  string `json:"device_asset_tag,omitempty"`
	HardwareProduct string `json:"hardware_product"`
	RackUnitStart   int    `json:"rack_unit_start"`
	RackUnitSize    int    `json:"rack_unit_size"`
}
type ResponseRackAssignments []ResponseRackAssignment

func (r ResponseRackAssignments) Len() int {
	return len(r)
}

func (r ResponseRackAssignments) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func (r ResponseRackAssignments) Less(i, j int) bool {
	return r[i].RackUnitStart < r[j].RackUnitStart
}
