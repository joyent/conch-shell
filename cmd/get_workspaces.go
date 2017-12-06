package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/mkideal/cli"
	"github.com/olekukonko/tablewriter"
	"os"
)

type getWorkspacesArgs struct {
	cli.Helper
}

var GetWorkspacesCmd = &cli.Command{
	Name: "get_workspaces",
	Desc: "Get a list of workspaces and their IDs",
	Argv: func() interface{} { return new(getWorkspacesArgs) },
	Fn: func(ctx *cli.Context) error {
		args, _, api, err := GetStarted(&getWorkspacesArgs{}, ctx)

		if err != nil {
			return err
		}

		//argv := args.Local.(*getWorkspacesArgs)

		workspaces, err := api.GetWorkspaces()
		if err != nil {
			return err
		}

		if args.Global.JSON == true {

			j, err := json.Marshal(workspaces)

			if err != nil {
				return err
			}

			fmt.Println(string(j))

		} else {
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Role", "Id", "Name", "Description"})

			table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
			table.SetCenterSeparator("|")

			for _, w := range workspaces {
				table.Append([]string{w.Role, w.Id, w.Name, w.Description})
			}

			table.Render()
		}
		return nil
	},
}
