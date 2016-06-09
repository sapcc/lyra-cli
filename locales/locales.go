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
	"automation-repository":           `Describes the place where the automation is being described. Git is the only supported repository type. Ex: https://github.com/user123/automation-test.git.`,
	"automation-repository-revision":  `Describes the repository branch.`,
	"automation-timeout":              `Describes the time elapsed before a timeout is being triggered.`,
	"automation-log-level":            `Describes the level should be used when logging.`,
	"automation-tags":                 `"Are key value pairs. Key-value pairs are separated by ':' or '='. Following this pattern: 'key1:value1,key2=value2...'."`,
	"automation-runlist":              `Describes the sequence of recipes should be executed. Runlist is an array of strings. Array of strings are separated by ','.`,
	"automation-attributes":           `Attributes are JSON based.`,
	"automation-attributes-from-file": `Path to the file containing the chef attributes in JSON format. Giving a dash '-' will be read from standard input.`,
	"automation-path":                 `Path to the script`,
	"automation-arguments":            `Arguments is an array of strings. Array of strings are separated by ','.`,
	"automation-environment":          `Key-value pairs are separated by ':' or '='. Following this pattern: 'key1:value1,key2=value2...'.`,
}

var errMsg = map[string]string{
	"automation-id-missing":       "No automation id provided.",
	"automation-selector-missing": "No automation selector given.",
	"run-id-missing":              "No automation run id given.",
	"job-id-missing":              "No job id provided.",
	"job-missing":                 fmt.Sprint(jobMissingDesc),
}

func AttributeDescription(id string) string {
	return attrDesc[id]
}

func ErrorMessages(id string) string {
	return errMsg[id]
}

var jobMissingDesc = `Job not found.

Note:
- Check if the job id matches.
- Jobs older than 30 days will be removed from the system. Check when the automation run is created by running following command:
lyra-cli run show --run-id={run_id}
`
