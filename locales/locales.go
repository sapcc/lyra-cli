package locales

import (
	"fmt"
)

var attrDesc = map[string]string{
	"json":                            `Attributes are JSON format.`,
	"selector":                        `Filter used to select on which nodes should the automation be executed. See link https://github.com/pages/monsoon/arc/docs/api/api.html#filter_agents for more information. Basic ex: @identity='{node_id}'.`,
	"run-id":                          `Automation run id.`,
	"job-id":                          `Job id.`,
	"watch":                           `Keep track of the running process.`,
	"automation-id":                   `Automation id.`,
	"automation-name":                 `Describes the template. Should be short and alphanumeric without white spaces.`,
	"automation-repository":           `Describes the place where the automation is being described. Git ist the only suported repository type. Ex: https://github.com/user123/automation-test.git.`,
	"automation-repository-revision":  `Describes the repository branch.`,
	"automation-timeout":              `Describes the time elapsed before a timeout is being triggered.`,
	"automation-log-level":            `Describe the level should be used when logging.`,
	"automation-tags":                 `"Are key value pairs. Key-value pairs are separated by ':' or '='. Following this pattern: 'key1:value1,key2=value2...'."`,
	"automation-runlist":              `Describes the sequence of recipes should be executed. Runlist is an array of strings. Array of strings are separated by ','.`,
	"automation-attributes":           `Attributes are JSON based.`,
	"automation-attributes-from-file": `Path to the file containing the chef attributes in JSON format. Giving a dash '-' will be read from standard input.`,
	"automation-path":                 `Path to the script`,
	"automation-arguments":            `Arguments is an array of strings. Array of strings are separated by ','.`,
	"automation-environment":          `Key-value pairs are separated by ':' or '='. Following this pattern: 'key1:value1,key2=value2...'.`,
	"arc-install-identifier":          "Compute instance id or external instance identifier.",
	"arc-install-os":                  "Instance's operating system. Available OS are linux or windows.",
	"update-service-url":              "Url to the arc update server.",
	"pki-service-url":                 "Url to the arc pki server.",
	"arc-broker-url":                  "Url to the arc broker.",
}

var errMsg = map[string]string{
	"automation-id-missing":       "No automation id provided.",
	"automation-selector-missing": "No automation selector given.",
	"run-id-missing":              "No automation run id given.",
	"job-id-missing":              "No job id provided.",
	"job-missing":                 fmt.Sprint(jobMissingDesc),
}

var cmdShortDescription = map[string]string{
	"arc":                               "Remote job execution framework.",
	"arc-install":                       "Retrieves the script used to install arc nodes on instances.",
	"authenticate":                      "Get an authentication token and endpoints for the automation and arc service.",
	"automation-create-chef":            "Create a new chef automation.",
	"automation-create-script":          "Create a new script automation.",
	"automation-create":                 "Create a new automation.",
	"automation-execute":                "Runs an exsiting automation",
	"automation-list":                   "List all available automations",
	"automation-show":                   "Show a specific automation",
	"automation-update-chef-attributes": "Updates chef attributes",
	"automation-update-chef":            "Updates a chef automation",
	"automation-update":                 "Updates an exsiting automation",
	"automation":                        "Automation service.",
	"bash-completion":                   "Generate completions for bash",
	"job-list":                          "List all jobs",
	"job-log":                           "Shows job log",
	"job-show":                          "Shows an especific job",
	"job":                               "Automation job service.",
	"root":                              "Automation service CLI",
	"run-list":                          "List all automation runs",
	"run-show":                          "Show a specific automation run",
	"run":                               "Automation run service.",
	"version":                           "Show program's version number and exit.",
}

var cmdLongDescription = map[string]string{
	"bash-completion": `Add $(lyra bash-completion) to your .bashrc to enable tab completion for lyra`,
	"root":            `Execute ad-hoc jobs using scripts, Chef and Ansible to configure machines and install the open source IaC service into any other OpenStack.`,
}

func AttributeDescription(id string) string {
	return attrDesc[id]
}

func ErrorMessages(id string) string {
	return errMsg[id]
}

func CmdShortDescription(id string) string {
	return cmdShortDescription[id]
}

func CmdLongDescription(id string) string {
	return cmdLongDescription[id]
}

var jobMissingDesc = `Job not found.

Note:
- Check if the job id matches.
- Jobs older than 30 days will be removed from the system. Check when the automation run is created running following command:
lyra-cli run show --run-id={run_id}
`
