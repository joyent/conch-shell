package cmd

import (
	"errors"
	config "github.com/joyent/conch-shell/config"
	conch "github.com/joyent/go-conch"
	pgtime "github.com/joyent/go-conch/pg_time"
	"github.com/mkideal/cli"
	"github.com/olekukonko/tablewriter"
	"time"
)

var (
	ConchConfigurationError = errors.New("No configuration data found or parse error")
	ConchNoApiSessionData   = errors.New("No API session data found")
)

type CliArgs struct {
	Local  interface{}
	Global *GlobalArgs
}

type MinimalDevice struct {
	Id       string             `json:"id"`
	AssetTag string             `json:"asset_tag"`
	Created  pgtime.ConchPgTime `json:"created,int"`
	LastSeen pgtime.ConchPgTime `json:"last_seen,int"`
	Health   string             `json:"health"`
	Flags    string             `json:"flags"`
}

// GetStarted handles the initial logic of parsing arguments, loading the JSON
// config file and verifying that the login credentials are still valid.
// Pretty much every command should start by using this function.
//
// Pro-tip: To cast args.Local to your command's arguments struct, use:
//   argv := args.Local.(*myLocalArgs)
func GetStarted(argv interface{}, ctx *cli.Context) (args *CliArgs, cfg *config.ConchConfig, api *conch.Conch, err error) {
	globals := &GlobalArgs{}
	if err := ctx.GetArgvList(argv, globals); err != nil {
		return nil, nil, nil, err
	}

	args = &CliArgs{
		Local:  argv,
		Global: globals,
	}

	cfg, err = config.NewFromJsonFile(globals.ConfigPath)
	if err != nil {
		return nil, nil, nil, ConchConfigurationError
	}

	if cfg.Session == "" {
		return nil, nil, nil, ConchNoApiSessionData
	}

	api = &conch.Conch{
		BaseUrl: cfg.Api,
		User:    cfg.User,
		Session: cfg.Session,
	}

	if err = api.VerifyLogin(); err != nil {
		return nil, nil, nil, err
	}

	return args, cfg, api, nil
}

func GenerateDeviceFlags(d conch.ConchDevice) (flags string) {
	flags = ""

	if !d.Deactivated.IsZero() {
		flags += "X"
	}

	if !d.Validated.IsZero() {
		flags += "v"
	}

	if !d.Graduated.IsZero() {
		flags += "g"
	}
	return flags
}

func TableizeMinimalDevices(devices []MinimalDevice, table *tablewriter.Table) *tablewriter.Table {
	table.SetHeader([]string{
		"ID",
		"Asset Tag",
		"Created",
		"Last Seen",
		"Health",
		"Flags",
	})

	for _, d := range devices {
		last_seen := ""
		if !d.LastSeen.IsZero() {
			last_seen = d.LastSeen.Format(time.UnixDate)
		}

		table.Append([]string{
			d.Id,
			d.AssetTag,
			d.Created.Format(time.UnixDate),
			last_seen,
			d.Health,
			d.Flags,
		})
	}

	return table
}
