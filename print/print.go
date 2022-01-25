package print

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/mitchellh/go-wordwrap"
	"github.com/olekukonko/tablewriter"
	"github.com/sapcc/lyra-cli/helpers"
)

type Print struct {
	Data interface{}
}

var ErrTypeAssertion = fmt.Errorf("not able to convert the data	")

func (p *Print) TableList(showColumns []string) (string, error) {
	// create table
	var buf bytes.Buffer
	table := tablewriter.NewWriter(&buf)
	table.SetColWidth(40)
	table.SetAlignment(3)
	table.SetHeader(showColumns)

	arrayStruct, ok := p.Data.([]interface{})
	if !ok {
		return "", ErrTypeAssertion
	}

	for _, valueMap := range arrayStruct {
		mapStruct, ok := valueMap.(map[string]interface{})
		if !ok {
			return "", ErrTypeAssertion
		}

		tableRow := []string{}
		for _, v := range showColumns {
			tableRow = append(tableRow, fmt.Sprintf("%v", mapStruct[v]))
		}
		table.Append(tableRow)
	}

	// print out
	table.Render()

	return buf.String(), nil
}

func (p *Print) Table() (string, error) {
	dataStruct, ok := p.Data.(map[string]interface{})
	if !ok {
		return "", ErrTypeAssertion
	}

	// sort map
	var keys []string
	for k := range dataStruct {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// create table
	var buf bytes.Buffer
	table := tablewriter.NewWriter(&buf)
	table.SetColWidth(20)
	table.SetAlignment(3)

	// set header
	table.SetHeader([]string{"Key", "Value"})

	// set body
	for _, k := range keys {
		value := wordwrap.WrapString(fmt.Sprintf("%v", dataStruct[k]), 150)
		table.Append([]string{fmt.Sprintf("%v", k), value})
	}

	// print out
	table.Render()

	return buf.String(), nil
}

func (p *Print) JSON() (string, error) {
	// convert data
	jsonData, err := helpers.StructureToJSON(p.Data)
	if err != nil {
		return "", err
	}

	var out bytes.Buffer
	err = json.Indent(&out, []byte(jsonData), "", "  ")
	if err != nil {
		return "", err
	}

	return out.String(), nil
}
