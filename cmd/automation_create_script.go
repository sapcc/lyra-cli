package cmd

import (
	"encoding/json"
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
			Name:               viper.GetString("automation-create-script-name"),
			Repository:         viper.GetString("automation-create-script-repository"),
			RepositoryRevision: viper.GetString("automation-create-script-repository-revision"),
			Timeout:            viper.GetInt("automation-create-script-timeout"),
			Path:               viper.GetString("automation-create-script-path"),
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
		bodyPrint := ""
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
		cmd.Println(bodyPrint)

		return nil
	},
}

func init() {
	AutomationCreateCmd.AddCommand(AutomationCreateScriptCmd)
	initAutomationCreateScriptCmdFlags()
}

func initAutomationCreateScriptCmdFlags() {
	// flags
	AutomationCreateScriptCmd.Flags().StringP("name", "", "", locales.AttributeDescription("automation-name"))
	AutomationCreateScriptCmd.Flags().StringP("repository", "", "", locales.AttributeDescription("automation-repository"))
	AutomationCreateScriptCmd.Flags().StringP("repository-revision", "", "master", locales.AttributeDescription("automation-repository-revision"))
	AutomationCreateScriptCmd.Flags().IntP("timeout", "", 3600, locales.AttributeDescription("automation-timeout"))
	AutomationCreateScriptCmd.Flags().StringP("tags", "", "", locales.AttributeDescription("automation-tags"))
	AutomationCreateScriptCmd.Flags().StringP("path", "", "", locales.AttributeDescription("automation-path"))
	AutomationCreateScriptCmd.Flags().StringP("arguments", "", "", locales.AttributeDescription("automation-arguments"))
	AutomationCreateScriptCmd.Flags().StringP("environment", "", "", locales.AttributeDescription("automation-environment"))
	viper.BindPFlag("automation-create-script-name", AutomationCreateScriptCmd.Flags().Lookup("name"))
	viper.BindPFlag("automation-create-script-repository", AutomationCreateScriptCmd.Flags().Lookup("repository"))
	viper.BindPFlag("automation-create-script-repository-revision", AutomationCreateScriptCmd.Flags().Lookup("repository-revision"))
	viper.BindPFlag("automation-create-script-timeout", AutomationCreateScriptCmd.Flags().Lookup("timeout"))
	viper.BindPFlag("automation-create-script-tags", AutomationCreateScriptCmd.Flags().Lookup("tags"))
	viper.BindPFlag("automation-create-script-path", AutomationCreateScriptCmd.Flags().Lookup("path"))
	viper.BindPFlag("automation-create-script-arguments", AutomationCreateScriptCmd.Flags().Lookup("arguments"))
	viper.BindPFlag("automation-create-script-environment", AutomationCreateScriptCmd.Flags().Lookup("environment"))
}

// private

func setupAutomationScriptAttr(scriptObj *Script) error {
	scriptObj.Tags = helpers.StringTokeyValueMap(viper.GetString("automation-create-script-tags"))
	scriptObj.Arguments = helpers.StringToArray(viper.GetString("automation-create-script-arguments"))
	scriptObj.Environment = helpers.StringTokeyValueMap(viper.GetString("automation-create-script-environment"))
	return nil
}

func automationCreateScript(scriptObj *Script) (string, error) {
	// add the type
	scriptObj.AutomationType = "Script"
	// convert to Json
	body, err := json.Marshal(scriptObj)
	if err != nil {
		return "", err
	}

	response, _, err := RestClient.Services.Automation.Post("automations", url.Values{}, string(body))
	if err != nil {
		return "", err
	}

	return response, nil
}
