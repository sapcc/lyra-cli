package locales

import (
	"fmt"
)

var attrDesc = map[string]string{
	"json":                            `Attributes are JSON format.`,
	"selector":                        `Filter used to select on which nodes should the automation be executed. See link https://github.com/pages/monsoon/arc/docs/api/api.html#filter_agents for more information. Basic ex: @identity='{node_id}'.`,
	"run-id":                          `Automation run identity.`,
	"job-id":                          `Job identity.`,
	"watch":                           `Keep track of the running process.`,
	"automation-id":                   `Automation identity.`,
	"automation-name":                 `Describes the template. Should be short and alphanumeric without white spaces.`,
	"automation-repository":           `Describes the place where the automation is being described. Git is the only supported repository type. Ex: https://github.com/user123/automation-test.git.`,
	"automation-repository-revision":  `Describes the repository branch.`,
	"automation-timeout":              `Describes the time elapsed before a timeout is being triggered.`,
	"automation-log-level":            `Describes the level should be used when logging.`,
	"automation-debug":                `Debug mode will not delete the temporary working directory on the instance when the automation job exists. This allows you to inspect the bundled automation artifacts, modify them and run the automation manually. Enabling debug mode for an extended period of time can exhaust  your instances disk space as each automation run will leave a directory behind. Also be aware that the payload may contain secrets which are persisted to disk indefinitely when debug mode is enabled. (default false)`,
	"automation-tags":                 `"Are key value pairs. Key-value pairs are separated by ':' or '='. Following this pattern: 'key1:value1,key2=value2...'."`,
	"automation-runlist":              `Describes the sequence of recipes should be executed. Runlist is an array of strings. Array of strings are separated by ','.`,
	"automation-chef-version":         `Specifies the Chef version should be installed in case no Chef is already been installed. (default latest)`,
	"automation-attributes":           `Attributes are JSON based.`,
	"automation-attributes-from-file": `Path to the file containing the chef attributes in JSON format. Giving a dash '-' will be read from standard input.`,
	"automation-path":                 `Path to the script`,
	"automation-arguments":            `Arguments is an array of strings. Array of strings are separated by ','.`,
	"automation-environment":          `Key-value pairs are separated by ':' or '='. Following this pattern: 'key1:value1,key2=value2...'.`,
	"install-format":                  `Installation script format. Supported: linux,windows,cloud-config,json.`,
	"node-id":                         `Node identity.`,
	"node-selector":                   `Filter nodes. See link https://github.com/pages/monsoon/arc/docs/api/api.html#filter_agents for more information. Basic ex: @identity='{node_id}'.`,
}

var errMsg = map[string]string{
	"automation-id-missing":       "No automation identity provided.",
	"automation-selector-missing": "No automation selector given.",
	"run-id-missing":              "No automation run identity given.",
	"job-id-missing":              "No job identity provided.",
	"node-id-missing":             "No node identity provided.",
	"job-missing":                 fmt.Sprint(jobMissingDesc),
	"node-missing":                fmt.Sprint(nodeMissingDesc),
	"flag-missing":                "Please make sure to provide following flags: ",
}

var cmdShortDescription = map[string]string{
	"arc":                               "Remote job execution framework.",
	"arc-node-install":                  "Retrieves the script used to install arc nodes on instances. User authentication flags are mandatory.",
	"arc-node-list":                     "List all nodes.",
	"arc-node-show":                     "Shows an especific node.",
	"arc-node-delete":                   "Deletes an especific node.",
	"arc-node-tag":                      "Node tags.",
	"arc-node-tag-list":                 "List all tags from an especific node.",
	"arc-node-tag-add":                  "Add tags to a given node.",
	"arc-node-tag-delete":               "Deletes tags from a given node.",
	"arc-node-fact":                     "Node facts.",
	"arc-node-fact-list":                "List all facts from an especific node.",
	"authenticate":                      "Get an authentication token and endpoints for the automation and arc service.",
	"automation-create-chef":            "Create a new chef automation.",
	"automation-create-script":          "Create a new script automation.",
	"automation-create":                 "Create a new automation.",
	"automation-execute":                "Runs an exsiting automation",
	"automation-list":                   "List all available automations",
	"automation-show":                   "Show a specific automation",
	"automation-delete":                 "Deletes a specific automation.",
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
	"bash-completion":     `Add $(lyra bash-completion) to your .bashrc to enable tab completion for lyra`,
	"root":                `Execute ad-hoc jobs using scripts, Chef and Ansible to configure machines and install the open source IaC service into any other OpenStack.`,
	"arc-node-delete":     "Deletes an especific node. \nThis will just delete the entry in the data base. For a permanent deletion you have to remove the node itself from the instance.",
	"arc-node-tag-add":    fmt.Sprint(nodeTagAddCmdLongDescription),
	"arc-node-tag-delete": fmt.Sprint(nodeTagDeleteCmdLongDescription),
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

var nodeTagDeleteCmdLongDescription = `Deletes tags from a given node.
Add the keys from the desired tags as command arguments.

Example:
lyra node tag delete --node-id 123456789 pool name plan"
`

var nodeTagAddCmdLongDescription = `Add tags to a given node.
Tags are key value pairs separated by the first "=" or ":" and added as command arguments. When using spacial characters use quotations.

Example:
lyra node tag add --node-id 123456789 pool:green name=db "plan=test new"`

var jobMissingDesc = `Job not found.

Note:
- Check if the job identity matches.
- Jobs older than 30 days will be removed from the system. Check when the automation run is created by running following command:
lyra-cli run show --run-id={run_id}
`

var nodeMissingDesc = `Node not found.

Note:
- Check if the node identity matches.
`
