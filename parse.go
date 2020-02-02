package main

import (
	"encoding/xml"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"time"
)

type ErrorType int

const (
	MissingReference ErrorType = iota
	NonExistingReference
	CycleError
)

type Error struct {
	Message string
	Type    ErrorType
}

type customTime struct {
	time.Time
}

type FS struct {
	XMLName xml.Name  `xml:"vw_FS"`
	Items   []*FSItem `xml:"FS"`
}

type ItemWithError struct {
	Errors []*Error `xml:"-"`
}

type FSItem struct {
	ItemWithError
	XMLName   xml.Name   `xml:"FS"`
	NameShort string     `xml:"FS_KURZ,attr"`
	NameLong  string     `xml:"FSLANG,attr"`
	Id        string     `xml:"s_NODE_FS_ID,attr"`
	Depth     int        `xml:"DEPTH,attr"`
	ParentId  string     `xml:"s_NODE_PARENT_ID,attr"`
	Parent    *FSItem    `xml:"-"`
	Children  []*FSItem  `xml:"-"`
	OE        []*OEItem  `xml:"-"`
	From      customTime `xml:"GAB,attr"`
	Until     customTime `xml:"GBIS,attr"`
}

type KU struct {
	XMLName xml.Name  `xml:"vw_KU"`
	Items   []*KUItem `xml:"KU"`
}

type KUItem struct {
	ItemWithError
	XMLName  xml.Name   `xml:"KU"`
	NameLong string     `xml:"KULANG,attr"`
	Id       string     `xml:"s_NODE_KU_ID,attr"`
	Depth    int        `xml:"DEPTH,attr"`
	ParentId string     `xml:"s_NODE_PARENT_ID,attr"`
	Parent   *KUItem    `xml:"-"`
	Children []*KUItem  `xml:"-"`
	OE       []*OEItem  `xml:"-"`
	From     customTime `xml:"GAB,attr"`
	Until    customTime `xml:"GBIS,attr"`
}

type OE struct {
	XMLName xml.Name  `xml:"OETBL"`
	Items   []*OEItem `xml:"OE"`
}

type OEItem struct {
	ItemWithError
	XMLName      xml.Name   `xml:"OE"`
	Id           string     `xml:"s_OE_ID,attr"`
	KUId         string     `xml:"s_KU_ID,attr"`
	KU           *KUItem    `xml:"-"`
	FSId         string     `xml:"s_FS_ID,attr"`
	FS           *FSItem    `xml:"-"`
	ParentLId    string     `xml:"s_PARENTOE_L_ID,attr"`
	ParentL      *OEItem    `xml:"-"`
	LChildren    []*OEItem  `xml:"-"`
	ParentFId    string     `xml:"s_PARENTOE_F_ID,attr"`
	ParentF      *OEItem    `xml:"-"`
	FChildren    []*OEItem  `xml:"-"`
	PSId         int        `xml:"PS_OEID,attr"`
	FSStart      int        `xml:"FS_START,attr"`
	From         customTime `xml:"Gültig_x0020_ab,attr"`
	Until        customTime `xml:"Gültig_x0020_bis,attr"`
	Type         string     `xml:"Typ,attr"`
	KUName       string     `xml:"Konzernunternehmen,attr"`
	FSName       string     `xml:"Führungsstruktur,attr"`
	OrgKZ        string     `xml:"Org-Kz,attr"`
	OrgName1     string     `xml:"Org-Bez1,attr"`
	OrgName2     string     `xml:"Org-Bez2,attr"`
	OrgName3     string     `xml:"Org-Bez3,attr"`
	Location     string     `xml:"Standort,attr"`
	CompanyName1 string     `xml:"Firmierung1,attr"`
	CompanyName2 string     `xml:"Firmierung2,attr"`
}

func (c *customTime) UnmarshalXMLAttr(attr xml.Attr) error {
	const shortForm = "2006-01-02T15:04:05"
	parse, err := time.Parse(shortForm, attr.Value)
	if err != nil {
		log.WithFields(log.Fields{
			"input": attr.Value,
			"err":   err,
		}).Error("date parsing error")
		return err
	}
	*c = customTime{parse}
	return nil
}

func readFile(path string) ([]byte, error) {
	xmlFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer xmlFile.Close()
	byteValue, _ := ioutil.ReadAll(xmlFile)
	return byteValue, err
}

func parseBytes(data []byte, v interface{}) error {
	xml.Unmarshal(data, &v)
	return nil
}

func parseFS(path string) (map[string]*FSItem, error) {
	data, err := readFile(path)
	if err != nil {
		return nil, err
	}

	return parseFSBytes(data)
}

func parseFSBytes(data []byte) (map[string]*FSItem, error) {
	var fs FS
	err := parseBytes(data, &fs)
	if err != nil {
		return nil, err
	}

	log.WithFields(log.Fields{
		"count": len(fs.Items),
	}).Debug("parsed FS items")

	fsMap := make(map[string]*FSItem)
	for _, item := range fs.Items {
		fsMap[item.Id] = item
	}

	return fsMap, nil
}

func parseKU(path string) (map[string]*KUItem, error) {
	data, err := readFile(path)
	if err != nil {
		return nil, err
	}

	return parseKUBytes(data)
}

func parseKUBytes(data []byte) (map[string]*KUItem, error) {
	var ku KU
	err := parseBytes(data, &ku)
	if err != nil {
		return nil, err
	}
	log.WithFields(log.Fields{
		"count": len(ku.Items),
	}).Debug("parsed KU items")

	kuMap := make(map[string]*KUItem)
	for _, item := range ku.Items {
		kuMap[item.Id] = item
	}

	return kuMap, nil
}

func parseOE(path string) ([]*OEItem, map[string]*OEItem, error) {
	data, err := readFile(path)
	if err != nil {
		return nil, nil, err
	}

	return parseOEBytes(data)
}

func parseOEBytes(data []byte) ([]*OEItem, map[string]*OEItem, error) {
	var oe OE
	err := parseBytes(data, &oe)
	if err != nil {
		return nil, nil, err
	}

	log.WithFields(log.Fields{
		"count": len(oe.Items),
	}).Debug("parsed OE items")

	oeMap := make(map[string]*OEItem)
	for _, item := range oe.Items {
		oeMap[item.Id] = item
	}

	return oe.Items, oeMap, nil
}
