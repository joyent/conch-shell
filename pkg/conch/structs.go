// Copyright Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package conch

import (
	"net/http"
	"net/http/cookiejar"
	"time"

	"github.com/joyent/conch-shell/pkg/pgtime"
	uuid "gopkg.in/satori/go.uuid.v1"
)

// ValidationReport vars provide an abstraction to make sense of the 'status'
// field in ValidationReports
const (
	ValidationReportStatusFail = 0
	ValidationReportStatusOK   = 1
)

// Conch contains auth and configuration data
type Conch struct {
	Session string // DEPRECATED
	BaseURL string
	UA      string
	JWToken string
	Expires int // This will be overwritten by JWT claims

	HTTPClient *http.Client
	CookieJar  *cookiejar.Jar
}

// Datacenter represents a conch datacenter, aka an AZ
type Datacenter struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	VendorName string    `json:"vendor_name"`
}

// Device represents what the API docs call a "DetailedDevice"
//
// Instead of having multiple structs representing partial datasets, like the
// API chooses to do, this library will always hand back Devices. In the
// case that the API does not provide all the data, those fields will be null
// or zero values.
type Device struct {
	AssetTag              string             `json:"asset_tag"`
	Created               pgtime.PgTime      `json:"created"`
	Deactivated           pgtime.PgTime      `json:"deactivated"`
	Graduated             pgtime.PgTime      `json:"graduated"`
	HardwareProduct       uuid.UUID          `json:"hardware_product"`
	Health                string             `json:"health"`
	Hostname              string             `json:"hostname"`
	ID                    string             `json:"id"`
	LastSeen              pgtime.PgTime      `json:"last_seen"`
	Location              DeviceLocation     `json:"location"`
	Nics                  []Nic              `json:"nics"`
	State                 string             `json:"state"`
	SystemUUID            uuid.UUID          `json:"system_uuid"`
	TritonUUID            uuid.UUID          `json:"triton_uuid"`
	TritonSetup           pgtime.PgTime      `json:"triton_setup"`
	Updated               pgtime.PgTime      `json:"updated"`
	UptimeSince           pgtime.PgTime      `json:"uptime_since"`
	Validated             pgtime.PgTime      `json:"validated"`
	Validations           []ValidationReport `json:"validations"`
	LatestReport          interface{}        `json:"latest_report"`
	LatestReportIsInvalid bool               `json:"latest_report_is_invalid"`
	InvalidReport         string             `json:"invalid_report"`
	Disks                 []Disk             `json:"disks"`
}

// DeviceDisk ...
type Disk struct {
	ID           uuid.UUID     `json:"id"`
	Created      pgtime.PgTime `json:"created"`
	Updated      pgtime.PgTime `json:"updated"`
	DriveType    string        `json:"drive_type"`
	Enclosure    string        `json:"enclosure"`
	Firmware     string        `json:"firmware"`
	HBA          string        `json:"hba"`
	Health       string        `json:"health"`
	Model        string        `json:"model"`
	SerialNumber string        `json:"serial_number"`
	Size         int           `json:"size"`
	Slot         int           `json:"slot"`
	Temp         int           `json:"temp"`
	Transport    string        `json:"transport"`
	Vendor       string        `json:"vendor"`
}

// DeviceLocation represents the location of a device, including its datacenter
// and rack
type DeviceLocation struct {
	Datacenter            Datacenter            `json:"datacenter"`
	Rack                  Rack                  `json:"rack"`
	TargetHardwareProduct HardwareProductTarget `json:"target_hardware_product"`
}

type ExtendedDevice struct {
	Device
	IPMI          string                    `json:"ipmi"`
	HardwareName  string                    `json:"hardware_name"`
	SKU           string                    `json:"sku"`
	Enclosures    map[string]map[int]Disk   `json:"enclosures"`
	IsGraduated   bool                      `json:"is_graduated"`
	IsTritonSetup bool                      `json:"is_triton_setup"`
	IsValidated   bool                      `json:"is_validated"`
	Validations   []ValidationPlanExecution `json:"validations"`
}

// GlobalDatacenter represents a datacenter in the global domain
type GlobalDatacenter struct {
	ID         uuid.UUID `json:"id"`
	Vendor     string    `json:"vendor"`
	VendorName string    `json:"vendor_name"`
	Region     string    `json:"region"`
	Location   string    `json:"location"`
	Created    time.Time `json:"created"`
	Updated    time.Time `json:"updated"`
}

// GlobalRack represents a datacenter rack in the global domain
type GlobalRack struct {
	ID               uuid.UUID `json:"id"`
	Created          time.Time `json:"created"`
	Updated          time.Time `json:"updated"`
	DatacenterRoomID uuid.UUID `json:"datacenter_room_id"`
	Name             string    `json:"name"`
	RoleID           uuid.UUID `json:"role"`
	SerialNumber     string    `json:"serial_number"`
	AssetTag         string    `json:"asset_tag"`
}

// GlobalRackLayoutSlot represents an individual rack layout entry
type GlobalRackLayoutSlot struct {
	ID        uuid.UUID `json:"id"`
	Created   time.Time `json:"created"`
	Updated   time.Time `json:"updated"`
	RackID    uuid.UUID `json:"rack_id"`
	ProductID uuid.UUID `json:"product_id"`
	RUStart   int       `json:"ru_start"`
}

// GlobalRackRole represents a rack role in the global domain
type GlobalRackRole struct {
	ID       uuid.UUID `json:"id"`
	Created  time.Time `json:"created"`
	Updated  time.Time `json:"updated"`
	Name     string    `json:"name"`
	RackSize int       `json:"rack_size"`
}

// GlobalRoom represents a datacenter room in the global domain
type GlobalRoom struct {
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
	Prefix            string          `json:"prefix"`
	HardwareVendorID  uuid.UUID       `json:"hardware_vendor_id"`
	GenerationName    string          `json:"generation_name"`
	LegacyProductName string          `json:"legacy_product_name"`
	SKU               string          `json:"sku"`
	Specification     string          `json:"specification"`
	Profile           HardwareProfile `json:"hardware_product_profile"`
	Created           pgtime.PgTime   `json:"created"`
	Updated           pgtime.PgTime   `json:"updated"`
}

// HardwareProductTarget represents the HardwareProduct that a device should
// have based on its location
type HardwareProductTarget struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Alias string    `json:"alias"`
}

// HardwareProfile is a detailed accounting of either the actual hardware or
// intended hardware configuration of a Device, depending on the API endpoint
// in question
type HardwareProfile struct {
	ID           uuid.UUID            `json:"id"`
	BiosFirmware string               `json:"bios_firmware"`
	CPUType      string               `json:"cpu_type"`
	HbaFirmware  string               `json:"hba_firmware"`
	NumCPU       int                  `json:"cpu_num"`
	NumDimms     int                  `json:"dimms_num"`
	NumNics      int                  `json:"nics_num"`
	NumSATA      int                  `json:"sata_num"`
	NumSSD       int                  `json:"ssd_num"`
	NumUSB       int                  `json:"usb_num"`
	Purpose      string               `json:"purpose"`
	SASNum       int                  `json:"sas_num"`
	SizeSAS      int                  `json:"sas_size"`
	SizeSATA     int                  `json:"sata_size"`
	SizeSSD      int                  `json:"ssd_size"`
	SlotsSAS     string               `json:"saas_slots"`
	SlotsSATA    string               `json:"sata_slots"`
	SlotsSSD     string               `json:"ssd_slots"`
	TotalPSU     int                  `json:"psu_total"`
	TotalRAM     int                  `json:"ram_total"`
	RackUnit     int                  `json:"rack_unit"`
	Zpool        HardwareProfileZpool `json:"zpool"`
}

// HardwareProfileZpool represents the layout of the target device's ZFS zpools
type HardwareProfileZpool struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Cache    int       `json:"cache"`
	Log      int       `json:"log"`
	DisksPer int       `json:"disks_per"`
	Spare    int       `json:"spare"`
	VdevN    int       `json:"vdev_n"`
	VdevT    string    `json:"vdev_t"`
}

// HardwareVendor ...
type HardwareVendor struct {
	ID      uuid.UUID     `json:"id"`
	Name    string        `json:"name"`
	Created pgtime.PgTime `json:"created"`
	Updated pgtime.PgTime `json:"updated"`
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

// Rack represents a physical rack
type Rack struct {
	ID           uuid.UUID  `json:"id"`
	Name         string     `json:"name"`
	Role         string     `json:"role"`
	Unit         int        `json:"unit"` // BUG(sungo): This exists because device locations provide rack info, but also slot info. This is a sloppy combination of data streams
	Size         int        `json:"size"`
	Datacenter   string     `json:"datacenter"`
	Slots        []RackSlot `json:"slots,omitempty"`
	SerialNumber string     `json:"serial_number"`
	AssetTag     string     `json:"asset_tag"`
}

// RackSlot represents a physical slot in a physical Rack
type RackSlot struct {
	ID            uuid.UUID `json:"id"`
	Size          int       `json:"size"`
	Name          string    `json:"name"`
	Alias         string    `json:"alias"`
	Vendor        string    `json:"vendor"`
	Occupant      Device    `json:"occupant"`
	RackUnitStart int       `json:"rack_unit_start"`
}

// Room represents a physical area in a datacenter/AZ
type Room struct {
	ID         string `json:"id"`
	AZ         string `json:"az"`
	Alias      string `json:"alias"`
	VendorName string `json:"vendor_name"`
}

// User represents a person able to access the Conch API or UI
type User struct {
	ID      string    `json:"id,omitempty"`
	Email   string    `json:"email"`
	Name    string    `json:"name"`
	Role    string    `json:"role"`
	RoleVia uuid.UUID `json:"role_via,omitempty"`
}

// UserDetailed ...
type UserDetailed struct {
	ID                  uuid.UUID          `json:"id"`
	Name                string             `json:"name"`
	Email               string             `json:"email"`
	Created             pgtime.PgTime      `json:"created"`
	LastLogin           pgtime.PgTime      `json:"last_login"`
	RefuseSessionAuth   bool               `json:"refuse_session_auth"`
	ForcePasswordChange bool               `json:"force_password_change"`
	Workspaces          []WorkspaceAndRole `json:"workspaces,omitempty"`
	IsAdmin             bool               `json:"is_admin"`
}

// Validation represents device validations loaded into Conch
type Validation struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Version     int       `json:"version"`
	Description string    `json:"description"`
	Created     time.Time `json:"created"`
	Updated     time.Time `json:"updated"`
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

// WorkspaceRelay represents a Conch Relay unit, a physical piece of hardware that
// mediates Livesys interactions in the field
type WorkspaceRelay struct {
	ID         string                 `json:"id"` // *not* a UUID
	Created    pgtime.PgTime          `json:"created"`
	Updated    pgtime.PgTime          `json:"updated"`
	Alias      string                 `json:"alias"`
	IPAddr     string                 `json:"ipaddr"`
	SSHPort    int                    `json:"ssh_port"`
	Version    string                 `json:"version"`
	LastSeen   pgtime.PgTime          `json:"last_seen"`
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

// WorkspaceUser ...
type WorkspaceUser struct {
	User
	RoleVia uuid.UUID `json:"role_via,omitempty"`
}
