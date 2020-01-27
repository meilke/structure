package main

import "testing"

func assertParentL(result AnalysisResult, expected OEItem, t *testing.T) {
	if result.ParentL.Id != expected.Id {
		t.Errorf("Wanted parent OE (L) %s, got: %s", expected.Id, result.ParentL.Id)
	}
}

func assertNoParentL(result AnalysisResult, t *testing.T) {
	if result.ParentL != nil {
		t.Errorf("Wanted no parent OE (L), got: %s", result.ParentL.Id)
	}
}

func assertHasParentLRef(result AnalysisResult, expected bool, t *testing.T) {
	if result.HasParentLRef != expected {
		t.Errorf("Wanted parent OE (L) ref: %t, got: %t", expected, result.HasParentLRef)
	}
}

func assertParentF(result AnalysisResult, expected OEItem, t *testing.T) {
	if result.ParentF.Id != expected.Id {
		t.Errorf("Wanted parent OE (F) %s, got: %s", expected.Id, result.ParentF.Id)
	}
}

func assertNoParentF(result AnalysisResult, t *testing.T) {
	if result.ParentF != nil {
		t.Errorf("Wanted no parent OE (F), got: %s", result.ParentF.Id)
	}
}

func assertHasParentFRef(result AnalysisResult, expected bool, t *testing.T) {
	if result.HasParentFRef != expected {
		t.Errorf("Wanted parent OE (F) ref: %t, got: %t", expected, result.HasParentFRef)
	}
}

func assertRelatedFS(result AnalysisResult, expected FSItem, t *testing.T) {
	if result.RelatedFS.Id != expected.Id {
		t.Errorf("Wanted related FS %s, got: %s", expected.Id, result.RelatedFS.Id)
	}
}

func assertNoRelatedFS(result AnalysisResult, t *testing.T) {
	if result.RelatedFS != nil {
		t.Errorf("Wanted no related FS, got: %s", result.RelatedFS.Id)
	}
}

func assertHasFSRef(result AnalysisResult, expected bool, t *testing.T) {
	if result.HasFSRef != expected {
		t.Errorf("Wanted FS ref: %t, got: %t", expected, result.HasFSRef)
	}
}

func assertRelatedKU(result AnalysisResult, expected KUItem, t *testing.T) {
	if result.RelatedKU.Id != expected.Id {
		t.Errorf("Wanted related KU %s, got: %s", expected.Id, result.RelatedKU.Id)
	}
}

func assertNoRelatedKU(result AnalysisResult, t *testing.T) {
	if result.RelatedKU != nil {
		t.Errorf("Wanted no related KU, got: %s", result.RelatedKU.Id)
	}
}

func assertHasKURef(result AnalysisResult, expected bool, t *testing.T) {
	if result.HasKURef != expected {
		t.Errorf("Wanted KU ref: %t, got: %t", expected, result.HasKURef)
	}
}

func TestAnalyzeOEEmptyMaps(t *testing.T) {
	maps := Maps{
		KU: make(map[string]KUItem),
		FS: make(map[string]FSItem),
		OE: make(map[string]OEItem),
	}

	oeItem := OEItem{
		Id:        "oe1",
		FSId:      "fs1",
		KUId:      "ku1",
		ParentLId: "oe10",
		ParentFId: "oe20",
	}

	result := analyzeOE(oeItem, maps)
	assertNoParentL(result, t)
	assertHasParentLRef(result, true, t)
	assertNoParentF(result, t)
	assertHasParentFRef(result, true, t)
	assertNoRelatedFS(result, t)
	assertHasFSRef(result, true, t)
	assertNoRelatedKU(result, t)
	assertHasKURef(result, true, t)
}

func TestAnalyzeOENoKU(t *testing.T) {
	maps := Maps{
		KU: make(map[string]KUItem),
		FS: make(map[string]FSItem),
		OE: make(map[string]OEItem),
	}
	maps.FS["fs1"] = FSItem{Id: "fs1"}
	maps.OE["oe10"] = OEItem{Id: "oe10"}
	maps.OE["oe20"] = OEItem{Id: "oe20"}

	oeItem := OEItem{
		Id:        "oe1",
		FSId:      "fs1",
		KUId:      "ku1",
		ParentLId: "oe10",
		ParentFId: "oe20",
	}

	result := analyzeOE(oeItem, maps)
	assertParentL(result, maps.OE["oe10"], t)
	assertHasParentLRef(result, true, t)
	assertParentF(result, maps.OE["oe20"], t)
	assertHasParentFRef(result, true, t)
	assertRelatedFS(result, maps.FS["fs1"], t)
	assertHasFSRef(result, true, t)
	assertNoRelatedKU(result, t)
	assertHasKURef(result, true, t)
}

func TestAnalyzeOENoFS(t *testing.T) {
	maps := Maps{
		KU: make(map[string]KUItem),
		FS: make(map[string]FSItem),
		OE: make(map[string]OEItem),
	}
	maps.KU["ku1"] = KUItem{Id: "ku1"}
	maps.OE["oe10"] = OEItem{Id: "oe10"}
	maps.OE["oe20"] = OEItem{Id: "oe20"}

	oeItem := OEItem{
		Id:        "oe1",
		FSId:      "fs1",
		KUId:      "ku1",
		ParentLId: "oe10",
		ParentFId: "oe20",
	}

	result := analyzeOE(oeItem, maps)
	assertParentL(result, maps.OE["oe10"], t)
	assertHasParentLRef(result, true, t)
	assertParentF(result, maps.OE["oe20"], t)
	assertHasParentFRef(result, true, t)
	assertNoRelatedFS(result, t)
	assertHasFSRef(result, true, t)
	assertRelatedKU(result, maps.KU["ku1"], t)
	assertHasKURef(result, true, t)
}

func TestAnalyzeOENoParentL(t *testing.T) {
	maps := Maps{
		KU: make(map[string]KUItem),
		FS: make(map[string]FSItem),
		OE: make(map[string]OEItem),
	}
	maps.FS["fs1"] = FSItem{Id: "fs1"}
	maps.KU["ku1"] = KUItem{Id: "ku1"}
	maps.OE["oe20"] = OEItem{Id: "oe20"}

	oeItem := OEItem{
		Id:        "oe1",
		FSId:      "fs1",
		KUId:      "ku1",
		ParentLId: "oe10",
		ParentFId: "oe20",
	}

	result := analyzeOE(oeItem, maps)
	assertNoParentL(result, t)
	assertHasParentLRef(result, true, t)
	assertParentF(result, maps.OE["oe20"], t)
	assertHasParentFRef(result, true, t)
	assertRelatedFS(result, maps.FS["fs1"], t)
	assertHasFSRef(result, true, t)
	assertRelatedKU(result, maps.KU["ku1"], t)
	assertHasKURef(result, true, t)
}

func TestAnalyzeOENoParentF(t *testing.T) {
	maps := Maps{
		KU: make(map[string]KUItem),
		FS: make(map[string]FSItem),
		OE: make(map[string]OEItem),
	}
	maps.FS["fs1"] = FSItem{Id: "fs1"}
	maps.KU["ku1"] = KUItem{Id: "ku1"}
	maps.OE["oe10"] = OEItem{Id: "oe10"}

	oeItem := OEItem{
		Id:        "oe1",
		FSId:      "fs1",
		KUId:      "ku1",
		ParentLId: "oe10",
		ParentFId: "oe20",
	}

	result := analyzeOE(oeItem, maps)
	assertParentL(result, maps.OE["oe10"], t)
	assertHasParentLRef(result, true, t)
	assertNoParentF(result, t)
	assertHasParentFRef(result, true, t)
	assertRelatedFS(result, maps.FS["fs1"], t)
	assertHasFSRef(result, true, t)
	assertRelatedKU(result, maps.KU["ku1"], t)
	assertHasKURef(result, true, t)
}

func TestAnalyzeOENoParentLRef(t *testing.T) {
	maps := Maps{
		KU: make(map[string]KUItem),
		FS: make(map[string]FSItem),
		OE: make(map[string]OEItem),
	}

	oeItem := OEItem{
		Id:        "oe1",
		FSId:      "fs1",
		KUId:      "ku1",
		ParentFId: "oe20",
	}

	result := analyzeOE(oeItem, maps)
	assertHasParentLRef(result, false, t)
	assertHasParentFRef(result, true, t)
	assertHasFSRef(result, true, t)
	assertHasKURef(result, true, t)
}

func TestAnalyzeOENoParentFRef(t *testing.T) {
	maps := Maps{
		KU: make(map[string]KUItem),
		FS: make(map[string]FSItem),
		OE: make(map[string]OEItem),
	}

	oeItem := OEItem{
		Id:        "oe1",
		FSId:      "fs1",
		KUId:      "ku1",
		ParentLId: "oe10",
	}

	result := analyzeOE(oeItem, maps)
	assertHasParentLRef(result, true, t)
	assertHasParentFRef(result, false, t)
	assertHasFSRef(result, true, t)
	assertHasKURef(result, true, t)
}

func TestAnalyzeOENoFSRef(t *testing.T) {
	maps := Maps{
		KU: make(map[string]KUItem),
		FS: make(map[string]FSItem),
		OE: make(map[string]OEItem),
	}

	oeItem := OEItem{
		Id:        "oe1",
		KUId:      "ku1",
		ParentLId: "oe10",
		ParentFId: "oe20",
	}

	result := analyzeOE(oeItem, maps)
	assertHasParentLRef(result, true, t)
	assertHasParentFRef(result, true, t)
	assertHasFSRef(result, false, t)
	assertHasKURef(result, true, t)
}

func TestAnalyzeOENoKURef(t *testing.T) {
	maps := Maps{
		KU: make(map[string]KUItem),
		FS: make(map[string]FSItem),
		OE: make(map[string]OEItem),
	}

	oeItem := OEItem{
		Id:        "oe1",
		FSId:      "fs1",
		ParentLId: "oe10",
		ParentFId: "oe20",
	}

	result := analyzeOE(oeItem, maps)
	assertHasParentLRef(result, true, t)
	assertHasParentFRef(result, true, t)
	assertHasFSRef(result, true, t)
	assertHasKURef(result, false, t)
}
