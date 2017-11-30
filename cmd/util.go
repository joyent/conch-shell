package cmd

import (
	"errors"
	config "github.com/joyent/conch-shell/config"
	conch "github.com/joyent/go-conch"
	"github.com/mkideal/cli"
)

var (
	ConchConfigurationError = errors.New("No configuration data found or parse error")
	ConchNoApiSessionData   = errors.New("No API session data found")
)

type CliArgs struct {
	Local  interface{}
	Global *GlobalArgs
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
