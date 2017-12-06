package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/mkideal/cli"
)

type getSettingsArgs struct {
	cli.Helper
}

var GetSettingsCmd = &cli.Command{
	Name: "get_settings",
	Desc: "Get the settings for the current user",
	Argv: func() interface{} { return new(getSettingsArgs) },
	Fn: func(ctx *cli.Context) error {
		args, _, api, err := GetStarted(&getSettingsArgs{}, ctx)

		if err != nil {
			return err
		}

		settings, err := api.GetUserSettings()
		if err != nil {
			return err
		}

		if args.Global.JSON == true {
			j, err := json.Marshal(settings)

			if err != nil {
				return err
			}

			fmt.Println(string(j))
		} else {
			if len(settings) == 0 {
				fmt.Println("No settings found")
			} else {
				for k, v := range settings {
					fmt.Printf("%s: %v\n", k, v)
				}
			}
		}
		return nil
	},
}
