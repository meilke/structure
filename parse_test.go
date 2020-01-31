package main

import (
	"testing"
	"time"
)

func parseTime(t string) customTime {
	tt, _ := time.Parse("2006-01-02T15:04:05", t)
	return customTime{tt}
}

func TestParseFS(t *testing.T) {
	input := `<vw_FS>
				<FS
					s_NODE_FS_ID="1E34442D-D6CD-47BE-8C41-FED7F4DD60C9"
					s_NODE_PARENT_ID="1E34442D-D6CD-47BE-8C41-FED7F4DD60C8"
					DEPTH="1"
					GAB="1900-01-01T00:00:00"
					GBIS="2025-12-31T00:00:00"
					FS_KURZ="extern"
					FSLANG="externe Firma" />
			  </vw_FS>`

	itemMap, err := parseFSBytes([]byte(input))
	if err != nil {
		t.Errorf("wanted no parsing error, got: %s", err)
	}

	if len(itemMap) != 1 {
		t.Errorf("wanted 1 item, got: %d", len(itemMap))
	}

	item := itemMap["1E34442D-D6CD-47BE-8C41-FED7F4DD60C9"]

	if item.Id != "1E34442D-D6CD-47BE-8C41-FED7F4DD60C9" {
		t.Errorf("wanted id 1E34442D-D6CD-47BE-8C41-FED7F4DD60C9, got: %s", item.Id)
	}

	if item.ParentId != "1E34442D-D6CD-47BE-8C41-FED7F4DD60C8" {
		t.Errorf("wanted parent id 1E34442D-D6CD-47BE-8C41-FED7F4DD60C8, got: %s", item.ParentId)
	}

	if item.Depth != 1 {
		t.Errorf("wanted depth 1 got: %d", item.Depth)
	}

	if item.From != parseTime("1900-01-01T00:00:00") {
		t.Errorf("wanted id %s got: %s", parseTime("1900-01-01T00:00:00"), item.From)
	}

	if item.Until != parseTime("2025-12-31T00:00:00") {
		t.Errorf("wanted id %s got: %s", parseTime("2025-12-31T00:00:00"), item.From)
	}

	if item.NameShort != "extern" {
		t.Errorf("wanted short name extern got: %s", item.NameShort)
	}

	if item.NameLong != "externe Firma" {
		t.Errorf("wanted long name externe Firma, got: %s", item.NameLong)
	}
}

func TestParseKU(t *testing.T) {
	input := `<vw_KU>
				<KU
					s_NODE_KU_ID="66470697-873F-4BF8-B762-72B028E5951C"
					s_NODE_PARENT_ID="66470697-873F-4BF8-B762-72B028E5951B"
					DEPTH="1"
					GAB="1900-01-01T00:00:00"
					GBIS="9999-12-31T00:00:00"
					KULANG="Bundeseisenbahnvermögen" />
			  </vw_KU>`

	itemMap, err := parseKUBytes([]byte(input))
	if err != nil {
		t.Errorf("wanted no parsing error, got: %s", err)
	}

	if len(itemMap) != 1 {
		t.Errorf("wanted 1 item, got: %d", len(itemMap))
	}

	item := itemMap["66470697-873F-4BF8-B762-72B028E5951C"]

	if item.Id != "66470697-873F-4BF8-B762-72B028E5951C" {
		t.Errorf("wanted id 66470697-873F-4BF8-B762-72B028E5951C, got: %s", item.Id)
	}

	if item.ParentId != "66470697-873F-4BF8-B762-72B028E5951B" {
		t.Errorf("wanted parent id 66470697-873F-4BF8-B762-72B028E5951B, got: %s", item.ParentId)
	}

	if item.Depth != 1 {
		t.Errorf("wanted depth 1 got: %d", item.Depth)
	}

	if item.From != parseTime("1900-01-01T00:00:00") {
		t.Errorf("wanted id %s got: %s", parseTime("1900-01-01T00:00:00"), item.From)
	}

	if item.Until != parseTime("9999-12-31T00:00:00") {
		t.Errorf("wanted id %s got: %s", parseTime("9999-12-31T00:00:00"), item.From)
	}

	if item.NameLong != "Bundeseisenbahnvermögen" {
		t.Errorf("wanted long name Bundeseisenbahnvermögen got: %s", item.NameLong)
	}
}

func TestParseOE(t *testing.T) {
	input := `<OETBL>
				<OE
					s_OE_ID="oe1"
					s_KU_ID="ku1"
					s_FS_ID="fs1"
					s_PARENTOE_L_ID="oe10"
					s_PARENTOE_F_ID="oe20"
					PS_OEID="7433"
					FS_START="0"
					Gültig_x0020_ab="2002-07-18T00:00:00"
					Gültig_x0020_bis="9999-12-31T00:00:00"
					Typ="Regionalbereich"
					Konzernunternehmen="DB Station&amp;Service AG"
					Führungsstruktur="Personenbahnhöfe"
					Org-Kz="I.SV-O"
					Org-Bez1="Leitung Regionalbereich Ost"
					Standort="Bln" />
			  </OETBL>`

	itemList, itemMap, err := parseOEBytes([]byte(input))
	if err != nil {
		t.Errorf("wanted no parsing error, got: %s", err)
	}

	if len(itemMap) != 1 {
		t.Errorf("wanted 1 item, got: %d", len(itemMap))
	}

	if len(itemList) != 1 {
		t.Errorf("wanted 1 item, got: %d", len(itemList))
	}

	item := itemMap["oe1"]

	if item.Id != "oe1" {
		t.Errorf("wanted id oe1, got: %s", item.Id)
	}

	if item.ParentLId != "oe10" {
		t.Errorf("wanted parent L id oe10, got: %s", item.ParentLId)
	}

	if item.ParentFId != "oe20" {
		t.Errorf("wanted parent F id oe20, got: %s", item.ParentFId)
	}

	if item.KUId != "ku1" {
		t.Errorf("wanted KU id ku1, got: %s", item.KUId)
	}

	if item.FSId != "fs1" {
		t.Errorf("wanted FS id fs1, got: %s", item.FSId)
	}

	if item.From != parseTime("2002-07-18T00:00:00") {
		t.Errorf("wanted id %s got: %s", parseTime("2002-07-18T00:00:00"), item.From)
	}

	if item.Until != parseTime("9999-12-31T00:00:00") {
		t.Errorf("wanted id %s got: %s", parseTime("9999-12-31T00:00:00"), item.From)
	}

	if item.Type != "Regionalbereich" {
		t.Errorf("wanted type Regionalbereich, got: %s", item.Type)
	}

	if item.KUName != "DB Station&Service AG" {
		t.Errorf("wanted KU name DB Station&Service AG, got: %s", item.KUName)
	}

	if item.FSName != "Personenbahnhöfe" {
		t.Errorf("wanted FS name Personenbahnhöfe, got: %s", item.FSName)
	}

	if item.OrgKZ != "I.SV-O" {
		t.Errorf("wanted org KZ I.SV-O, got: %s", item.OrgKZ)
	}

	// to be continued ;-)
}
