package cmd

import (
	"fmt"
	Helper "github.com/XNXKTech/athena/cmd/helper"
	"github.com/XNXKTech/athena/cmd/resource"
	"github.com/spf13/cobra"
)

// controllerCmd represents the controller command
var controllerCmd = &cobra.Command{
	Use:   "controller",
	Short: "Create a new controller class",
	Long:  `Create a new controller class`,
	Run: func(cmd *cobra.Command, args []string) {
		//根据模板生成控制器
		Helper.GenFile(Helper.UnGzip(resource.CONTROLLER_TPL),
			fmt.Sprintf("/src/classes/%sClass.go", Helper.Ucfirst(args[0])),
			map[string]interface{}{
				"ControllerName": args[0],
			},
		)
		fmt.Println(fmt.Sprintf("Controller [src/classes/%sClass.go] created successfully.", Helper.Ucfirst(args[0])))
	},
	Args: cobra.MinimumNArgs(1),
}

func init() {
	newCmd.AddCommand(controllerCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// controllerCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// controllerCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
