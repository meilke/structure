package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

type customTime struct {
	time.Time
}

type FS struct {
	XMLName xml.Name `xml:"vw_FS"`
	Items   []FSItem `xml:"FS"`
}

type FSItem struct {
	XMLName   xml.Name   `xml:"FS"`
	NameShort string     `xml:"FS_KURZ,attr"`
	NameLong  string     `xml:"FSLANG,attr"`
	Id        string     `xml:"s_NODE_FS_ID,attr"`
	Depth     int        `xml:"DEPTH,attr"`
	ParentId  string     `xml:"s_NODE_PARENT_ID,attr"`
	From      customTime `xml:"GAB,attr"`
	Until     customTime `xml:"GBIS,attr"`
}

type KU struct {
	XMLName xml.Name `xml:"vw_KU"`
	Items   []KUItem `xml:"KU"`
}

type KUItem struct {
	XMLName  xml.Name   `xml:"KU"`
	NameLong string     `xml:"KULANG,attr"`
	Id       string     `xml:"s_NODE_KU_ID,attr"`
	Depth    int        `xml:"DEPTH,attr"`
	ParentId string     `xml:"s_NODE_PARENT_ID,attr"`
	From     customTime `xml:"GAB,attr"`
	Until    customTime `xml:"GBIS,attr"`
}

type OE struct {
	XMLName xml.Name `xml:"OETBL"`
	Items   []OEItem `xml:"OE"`
}

type OEItem struct {
	XMLName      xml.Name   `xml:"OE"`
	Id           string     `xml:"s_OE_ID,attr"`
	KUId         string     `xml:"s_KU_ID,attr"`
	FSId         string     `xml:"s_FS_ID,attr"`
	ParentLId    string     `xml:"s_PARENTOE_L_ID,attr"`
	ParentFId    string     `xml:"s_PARENTOE_F_ID,attr"`
	PSId         int        `xml:"PS_OEID,attr"`
	FSStart      int        `xml:"FS_START,attr"`
	From         customTime `xml:"Gültig_x0020_ab,attr"`
	Until        customTime `xml:"Gültig_x0020_bis,attr"`
	Type         string     `xml:"Typ,attr"`
	KU           string     `xml:"Konzernunternehmen,attr"`
	FS           string     `xml:"Führungsstruktur,attr"`
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
		fmt.Println(err)
		return err
	}
	*c = customTime{parse}
	return nil
}

func parseFile(path string, v interface{}) error {
	xmlFile, err := os.Open(path)
	if err != nil {
		return err
	}

	defer xmlFile.Close()
	byteValue, _ := ioutil.ReadAll(xmlFile)

	xml.Unmarshal(byteValue, &v)

	return nil
}

func parseFS(path string) (map[string]FSItem, error) {
	var fs FS
	err := parseFile(path, &fs)
	if err != nil {
		return nil, err
	}

	//TODO logging!
	fmt.Printf("parsed FS: %d items\n", len(fs.Items))

	fsMap := make(map[string]FSItem)
	for _, item := range fs.Items {
		fsMap[item.Id] = item
	}

	return fsMap, nil
}

func parseKU(path string) (map[string]KUItem, error) {
	var ku KU
	err := parseFile(path, &ku)
	if err != nil {
		return nil, err
	}
	fmt.Printf("parsed KU: %d items\n", len(ku.Items))

	kuMap := make(map[string]KUItem)
	for _, item := range ku.Items {
		kuMap[item.Id] = item
	}

	return kuMap, nil
}

func parseOE(path string) ([]OEItem, map[string]OEItem, error) {
	var oe OE
	err := parseFile(path, &oe)
	if err != nil {
		return nil, nil, err
	}
	fmt.Printf("parsed OE: %d items\n", len(oe.Items))

	oeMap := make(map[string]OEItem)
	for _, item := range oe.Items {
		oeMap[item.Id] = item
	}

	return oe.Items, oeMap, nil
}
