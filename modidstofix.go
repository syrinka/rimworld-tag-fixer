package main

import (
	"encoding/xml"
	"os"
)

type ModIdsToFix struct {
	XMLName xml.Name `xml:"ModIdsToFix"`
	Ids     []string `xml:"li"`
}

func collectFixable(path string) []string {
	// if file not exist, return empty list
	data, err := os.ReadFile(path)
	if err != nil {
		return []string{}
	}
	var st ModIdsToFix
	xml.Unmarshal(data, &st)
	return st.Ids
}
