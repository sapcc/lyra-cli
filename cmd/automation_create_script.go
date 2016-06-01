package cmd

import (
	"encoding/json"
	"net/url"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/sapcc/lyra-cli/helpers"
	"github.com/sapcc/lyra-cli/print"
)

// createCmd represents the create command
var AutomationCreateScriptCmd = &cobra.Command{
	Use:   "script",
	Short: "Create a new script automation.",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command.`,
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
	AutomationCreateScriptCmd.Flags().StringP("name", "", "", "Describes the template. Should be short and alphanumeric without white spaces.")
	AutomationCreateScriptCmd.Flags().StringP("repository", "", "", "Describes the place where the automation is being described. Git ist the only suported repository type. Ex: https://github.com/user123/automation-test.git.")
	AutomationCreateScriptCmd.Flags().StringP("repository-revision", "", "master", "Describes the repository branch.")
	AutomationCreateScriptCmd.Flags().IntP("timeout", "", 3600, "Describes the time elapsed before a timeout is being triggered.")
	AutomationCreateScriptCmd.Flags().StringP("tags", "", "", "Are key value pairs. Key-value pairs are separated by ':' or '='. Following this pattern: 'key1:value1,key2=value2...'.")
	AutomationCreateScriptCmd.Flags().StringP("path", "", "", "Path to the script")
	AutomationCreateScriptCmd.Flags().StringP("arguments", "", "", "Arguments is an array of strings. Array of strings are separated by ','.")
	AutomationCreateScriptCmd.Flags().StringP("environment", "", "", "Key-value pairs are separated by ':' or '='. Following this pattern: 'key1:value1,key2=value2...'.")
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
