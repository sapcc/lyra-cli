package locales

import (
	"fmt"
)

var attrDesc = map[string]string{
	"automation-id":       `Automation id.`,
	"automation-selector": `Filter used to select on which nodes should the automation be runed. See link https://github.com/pages/monsoon/arc/docs/api/api.html#filter_agents for more information. Basic ex: @identity='{node_id}'.`,
	"page":                `Set the pagination page.`,
	"per-page":            `Set the elements per page.`,
	"run-id":              `Automation run id.`,
	"job-id":              `Job id.`,
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
- Jobs older than 30 days will be removed from the system. Check when the automation run is created running following command:
lyra-cli run show --run-id={run_id}
`
