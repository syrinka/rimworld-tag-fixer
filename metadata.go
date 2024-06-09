package main

import (
	"github.com/beevik/etree"
)

type ModMetaData struct {
	path string
	doc  *etree.Document
	root *etree.Element
}

func (meta *ModMetaData) Init(path string) error {
	meta.path = path
	meta.doc = etree.NewDocument()
	err := meta.doc.ReadFromFile(path)
	if err != nil {
		return err
	}
	meta.root = meta.doc.SelectElement("ModMetaData")
	if meta.root == nil {
		// some mod use non-standard tag
		// for example 2135654735
		meta.root = meta.doc.SelectElement("modMetaData")
	}
	return nil
}

func (meta ModMetaData) Name() string {
	return meta.root.SelectElement("name").Text()
}

func (meta ModMetaData) Author() string {
	return meta.root.SelectElement("author").Text()
}

func (meta ModMetaData) Id() string {
	return meta.root.SelectElement("packageId").Text()
}

func (meta ModMetaData) GetVersionTags() []string {
	var elements = meta.root.SelectElement("supportedVersions").ChildElements()
	var tags = make([]string, len(elements)+1)

	for i, e := range elements {
		tags[i] = e.Text()
	}
	return tags
}

func (meta ModMetaData) ContainVersionTag(tag string) bool {
	for _, i := range meta.GetVersionTags() {
		if i == tag {
			return true
		}
	}
	return false
}

func (meta *ModMetaData) AddVersionTag(tag string) {
	if meta.ContainVersionTag(tag) {
		return
	}
	var e = meta.root.SelectElement("supportedVersions").CreateElement("li")
	e.SetText(tag)
}

func (meta ModMetaData) Update() {
	meta.doc.Indent(2)
	meta.doc.WriteToFile(meta.path)
}
