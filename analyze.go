package main

type Maps struct {
	OE map[string]OEItem
	FS map[string]FSItem
	KU map[string]KUItem
}

type AnalysisResult struct {
	ParentL       *OEItem
	HasParentLRef bool
	ParentF       *OEItem
	HasParentFRef bool
	RelatedFS     *FSItem
	HasFSRef      bool
	RelatedKU     *KUItem
	HasKURef      bool
}

func analyzeOE(oeItem OEItem, maps Maps) AnalysisResult {
	result := AnalysisResult{
		HasParentLRef: true,
		HasParentFRef: true,
		HasFSRef:      true,
		HasKURef:      true,
	}

	if oeItem.ParentLId == "" {
		result.HasParentLRef = false
	} else if parentL, ok := maps.OE[oeItem.ParentLId]; ok {
		result.ParentL = &parentL
	}

	if oeItem.ParentFId == "" {
		result.HasParentFRef = false
	} else if parentF, ok := maps.OE[oeItem.ParentFId]; ok {
		result.ParentF = &parentF
	}

	if oeItem.FSId == "" {
		result.HasFSRef = false
	} else if fs, ok := maps.FS[oeItem.FSId]; ok {
		result.RelatedFS = &fs
	}

	if oeItem.KUId == "" {
		result.HasKURef = false
	} else if ku, ok := maps.KU[oeItem.KUId]; ok {
		result.RelatedKU = &ku
	}

	return result
}
