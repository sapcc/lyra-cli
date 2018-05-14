package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/sapcc/lyra-cli/helpers"
	"github.com/sapcc/lyra-cli/locales"
	"github.com/sapcc/lyra-cli/print"
)

// createCmd represents the create command
var AutomationCreateScriptCmd = &cobra.Command{
	Use:   "script",
	Short: locales.CmdShortDescription("automation-create-script"),
	RunE: func(cmd *cobra.Command, args []string) error {
		script = Script{
			Automation: Automation{
				Name:               viper.GetString("automation-create-script-name"),
				Repository:         viper.GetString("automation-create-script-repository"),
				RepositoryRevision: viper.GetString("automation-create-script-repository-revision"),
				Timeout:            viper.GetInt("automation-create-script-timeout"),
			},
			Path: viper.GetString("automation-create-script-path"),
		}

		// setup automation create script attributes
		err := setupAutomationScriptAttr(&script)
		if err != nil {
			return err
		}

		// create automation
		response, err := automationCreateScript(&script)
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

		// Print response
		fmt.Println(bodyPrint)

		return nil
	},
}

func init() {
	AutomationCreateCmd.AddCommand(AutomationCreateScriptCmd)
	initAutomationCreateScriptCmdFlags()
}

func initAutomationCreateScriptCmdFlags() {
	// flags
	AutomationCreateScriptCmd.Flags().String("name", "", locales.AttributeDescription("automation-name"))
	AutomationCreateScriptCmd.Flags().String("repository", "", locales.AttributeDescription("automation-repository"))
	AutomationCreateScriptCmd.Flags().String("repository-revision", "master", locales.AttributeDescription("automation-repository-revision"))
	AutomationCreateScriptCmd.Flags().Int("timeout", 3600, locales.AttributeDescription("automation-timeout"))
	AutomationCreateScriptCmd.Flags().String("path", "", locales.AttributeDescription("automation-path"))
	AutomationCreateScriptCmd.Flags().StringArray("arg", nil, locales.AttributeDescription("automation-argument"))
	AutomationCreateScriptCmd.Flags().StringArray("env", nil, locales.AttributeDescription("automation-environment"))
	viper.BindPFlag("automation-create-script-name", AutomationCreateScriptCmd.Flags().Lookup("name"))
	viper.BindPFlag("automation-create-script-repository", AutomationCreateScriptCmd.Flags().Lookup("repository"))
	viper.BindPFlag("automation-create-script-repository-revision", AutomationCreateScriptCmd.Flags().Lookup("repository-revision"))
	viper.BindPFlag("automation-create-script-timeout", AutomationCreateScriptCmd.Flags().Lookup("timeout"))
	viper.BindPFlag("automation-create-script-path", AutomationCreateScriptCmd.Flags().Lookup("path"))
	viper.BindPFlag("automation-create-script-argument", AutomationCreateScriptCmd.Flags().Lookup("arg"))
	viper.BindPFlag("automation-create-script-environment", AutomationCreateScriptCmd.Flags().Lookup("env"))
}

// private

func setupAutomationScriptAttr(scriptObj *Script) (err error) {
	scriptObj.Arguments = viper.GetStringSlice("automation-create-script-argument")

	scriptObj.Environment, err = helpers.StringSliceKeyValueMap(viper.GetStringSlice("automation-create-script-environment"))
	return
}

func automationCreateScript(scriptObj *Script) (string, error) {
	// add the type
	scriptObj.AutomationType = "Script"
	// convert to Json
	body, err := json.Marshal(scriptObj)
	if err != nil {
		return "", err
	}

	automationService := RestClient.Services["automation"]
	response, _, err := automationService.Post("automations", url.Values{}, http.Header{}, string(body))
	if err != nil {
		return "", err
	}

	return response, nil
}
