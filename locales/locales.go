package locales

import ()

var attrDesc = map[string]string{
	"automation-id":       `Automation id.`,
	"automation-selector": `Filter used to select on which nodes should the automation be runed. See link https://github.com/pages/monsoon/arc/docs/api/api.html#filter_agents for more information. Basic ex: @identity='{node_id}'.`,
}

var errMsg = map[string]string{
	"automation-id-missing":       "No automation id given.",
	"automation-selector-missing": "No automation selector given.",
}

func AttributeDescription(id string) string {
	return attrDesc[id]
}

func ErrorMessages(id string) string {
	return errMsg[id]
}
