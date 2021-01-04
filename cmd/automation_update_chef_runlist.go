package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"path"

	"github.com/sapcc/lyra-cli/helpers"
	"github.com/sapcc/lyra-cli/locales"
	"github.com/sapcc/lyra-cli/print"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var AutomationUpdateChefRunlistCmd = &cobra.Command{
	Use:   "runlist",
	Short: locales.CmdShortDescription("automation-update-chef-runlist"),
	Long:  locales.CmdLongDescription("automation-update-chef-runlist"),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// check required automation id
		if len(viper.GetString("automation-update-chef-automation-id")) == 0 {
			return errors.New(locales.ErrorMessages("automation-id-missing"))
		}
		return nil
	},

	RunE: func(cmd *cobra.Command, args []string) error {
		chef = Chef{}

		// set chef runlist
		chef.Runlist = helpers.StringToArray(viper.GetString("automation-update-chef-runlist"))

		// update automation
		response, err := automationUpdateChefRunlist(&chef)
		if err != nil {
			return err
		}

		// convert data to struct
		var dataStruct map[string]interface{}
		err = helpers.JSONStringToStructure(response, &dataStruct)
		if err != nil {
			return err
		}

		// print the data out
		printer := print.Print{Data: dataStruct}
		var bodyPrint string
		if viper.GetBool("json") {
			bodyPrint, err = printer.JSON()
			if err != nil {
				return err
			}
		} else {
			bodyPrint, err = printer.Table()
			if err != nil {
				return err
			}
		}

		// print response
		fmt.Println(bodyPrint)

		return nil
	},
}

func init() {
	AutomationUpdateChefCmd.AddCommand(AutomationUpdateChefRunlistCmd)
	initAutomationUpdateChefRunlistCmdFlags()
}

func initAutomationUpdateChefRunlistCmdFlags() {
	AutomationUpdateChefRunlistCmd.Flags().StringP("runlist", "", "", locales.AttributeDescription("automation-runlist"))
	AutomationUpdateChefRunlistCmd.Flags().StringP(FLAG_AUTOMATION_ID, "", "", locales.AttributeDescription("automation-id"))
	helpers.CheckErrAndPrintToStdErr(viper.BindPFlag("automation-update-chef-runlist", AutomationUpdateChefRunlistCmd.Flags().Lookup("runlist")), "BindPFlag:")
	helpers.CheckErrAndPrintToStdErr(viper.BindPFlag("automation-update-chef-automation-id", AutomationUpdateChefRunlistCmd.Flags().Lookup("automation-id")), "BindPFlag:")
}

func automationUpdateChefRunlist(chefObj *Chef) (string, error) {
	automationService := RestClient.Services["automation"]

	response, code, err := automationService.Get(path.Join("automations", viper.GetString("automation-update-chef-automation-id")), url.Values{}, false)
	if err != nil {
		return "", err
	}

	if int(code) >= 400 {
		return "", errors.New(response)
	}

	// get the existing data
	oldChef := Chef{}
	respByt := []byte(response)
	if err := json.Unmarshal(respByt, &oldChef); err != nil {
		return "", err
	}

	// change runlist
	oldChef.Runlist = chefObj.Runlist

	// convert to Json
	body, err := json.Marshal(oldChef)
	if err != nil {
		return "", err
	}

	// send data back
	newResp, _, err := automationService.Put(path.Join("automations", viper.GetString("automation-update-chef-automation-id")), url.Values{}, string(body))
	if err != nil {
		return "", err
	}

	return newResp, nil
}
