package main

import (
	"testing"
)

func assertError(expected Error, errors []*Error, t *testing.T) {
	var found bool
	for _, e := range errors {
		found = e.Message == expected.Message && e.Type == expected.Type
		if found {
			return
		}
	}

	if !found {
		t.Errorf("Wanted error: %s (%d), did not get it", expected.Message, expected.Type)
	}
}

func assertErrorCount(expected int, errors []*Error, t *testing.T) {
	if len(errors) != expected {
		t.Errorf("Wanted %d errors, got %d", expected, len(errors))
	}
}

func assertNoRelatedFS(item *OEItem, t *testing.T) {
	if item.FS != nil {
		t.Errorf("Wanted no related FS, got: %s", item.FS.Id)
	}
}

func assertHasFSRef(item *OEItem, fsItem *FSItem, t *testing.T) {
	if item.FS == nil {
		t.Errorf("Wanted FS ref: %s, got no ref", fsItem.Id)
	}
}

func assertNoRelatedKU(item *OEItem, t *testing.T) {
	if item.KU != nil {
		t.Errorf("Wanted no related KU, got: %s", item.KU.Id)
	}
}

func assertHasKURef(item *OEItem, kuItem *KUItem, t *testing.T) {
	if item.KU == nil {
		t.Errorf("Wanted KU ref: %s, got no ref", kuItem.Id)
	}
}

func assertNoParentLRef(item *OEItem, t *testing.T) {
	if item.ParentL != nil {
		t.Errorf("Wanted no parent L, got: %s", item.ParentL.Id)
	}
}

func assertHasParentLRef(item *OEItem, other *OEItem, t *testing.T) {
	if item.ParentL == nil {
		t.Errorf("Wanted parent L ref: %s, got no ref", other.Id)
	}
}

func assertNoParentFRef(item *OEItem, t *testing.T) {
	if item.ParentF != nil {
		t.Errorf("Wanted no parent F, got: %s", item.ParentF.Id)
	}
}

func assertHasParentFRef(item *OEItem, other *OEItem, t *testing.T) {
	if item.ParentF == nil {
		t.Errorf("Wanted parent F ref: %s, got no ref", other.Id)
	}
}

func assertNoParentKU(item *KUItem, t *testing.T) {
	if item.Parent != nil {
		t.Errorf("Wanted no parent, got: %s", item.Parent.Id)
	}
}

func assertHasParentKURef(item *KUItem, other *KUItem, t *testing.T) {
	if item.Parent == nil {
		t.Errorf("Wanted parent ref: %s, got no ref", other.Id)
	}
}

func assertNoParentFS(item *FSItem, t *testing.T) {
	if item.Parent != nil {
		t.Errorf("Wanted no parent, got: %s", item.Parent.Id)
	}
}

func assertHasParentFSRef(item *FSItem, other *FSItem, t *testing.T) {
	if item.Parent == nil {
		t.Errorf("Wanted parent ref: %s, got no ref", other.Id)
	}
}

func TestKUMissingParent(t *testing.T) {
	kuMap := make(map[string]*KUItem)
	fsMap := make(map[string]*FSItem)
	oeMap := make(map[string]*OEItem)

	kuItem := &KUItem{
		Id:       "ku1",
		ParentId: "ku10",
	}
	kuMap[kuItem.Id] = kuItem

	buildTrees(oeMap, kuMap, fsMap)
	analyzeTrees(oeMap, kuMap, fsMap)

	assertErrorCount(1, kuItem.Errors, t)
	assertNoParentKU(kuItem, t)
	assertError(Error{Message: NON_EXISTING_PARENT, Type: NonExistingReference}, kuItem.Errors, t)
}

func TestKUGood(t *testing.T) {
	kuMap := make(map[string]*KUItem)
	fsMap := make(map[string]*FSItem)
	oeMap := make(map[string]*OEItem)

	kuItem := &KUItem{
		Id:       "ku1",
		ParentId: "ku10",
	}
	kuMap[kuItem.Id] = kuItem

	kuMap["ku10"] = &KUItem{Id: "ku10"}

	buildTrees(oeMap, kuMap, fsMap)
	analyzeTrees(oeMap, kuMap, fsMap)

	assertErrorCount(0, kuItem.Errors, t)
	assertHasParentKURef(kuItem, kuMap["ku10"], t)
}

func TestKUCycle(t *testing.T) {
	kuMap := make(map[string]*KUItem)
	fsMap := make(map[string]*FSItem)
	oeMap := make(map[string]*OEItem)

	kuItem1 := &KUItem{
		Id:       "ku1",
		ParentId: "ku2",
	}
	kuMap[kuItem1.Id] = kuItem1

	fsItem2 := &KUItem{
		Id:       "ku2",
		ParentId: "ku1",
	}
	kuMap[fsItem2.Id] = fsItem2

	buildTrees(oeMap, kuMap, fsMap)
	analyzeTrees(oeMap, kuMap, fsMap)

	assertErrorCount(1, kuItem1.Errors, t)
	assertError(Error{Message: CYCLE_REFERENCE, Type: CycleError}, kuItem1.Errors, t)
	assertErrorCount(1, fsItem2.Errors, t)
	assertError(Error{Message: CYCLE_REFERENCE, Type: CycleError}, fsItem2.Errors, t)
}

func TestFSMissingParent(t *testing.T) {
	kuMap := make(map[string]*KUItem)
	fsMap := make(map[string]*FSItem)
	oeMap := make(map[string]*OEItem)

	fsItem := &FSItem{
		Id:       "fs1",
		ParentId: "fs10",
	}
	fsMap[fsItem.Id] = fsItem

	buildTrees(oeMap, kuMap, fsMap)
	analyzeTrees(oeMap, kuMap, fsMap)

	assertErrorCount(1, fsItem.Errors, t)
	assertNoParentFS(fsItem, t)
	assertError(Error{Message: NON_EXISTING_PARENT, Type: NonExistingReference}, fsItem.Errors, t)
}

func TestFSGood(t *testing.T) {
	kuMap := make(map[string]*KUItem)
	fsMap := make(map[string]*FSItem)
	oeMap := make(map[string]*OEItem)

	fsItem := &FSItem{
		Id:       "fs1",
		ParentId: "fs10",
	}
	fsMap[fsItem.Id] = fsItem

	fsMap["fs10"] = &FSItem{Id: "fs10"}

	buildTrees(oeMap, kuMap, fsMap)
	analyzeTrees(oeMap, kuMap, fsMap)

	assertErrorCount(0, fsItem.Errors, t)
	assertHasParentFSRef(fsItem, fsMap["fs10"], t)
}

func TestFSCycle(t *testing.T) {
	kuMap := make(map[string]*KUItem)
	fsMap := make(map[string]*FSItem)
	oeMap := make(map[string]*OEItem)

	fsItem1 := &FSItem{
		Id:       "fs1",
		ParentId: "fs2",
	}
	fsMap[fsItem1.Id] = fsItem1

	fsItem2 := &FSItem{
		Id:       "fs2",
		ParentId: "fs1",
	}
	fsMap[fsItem2.Id] = fsItem2

	buildTrees(oeMap, kuMap, fsMap)
	analyzeTrees(oeMap, kuMap, fsMap)

	assertErrorCount(1, fsItem1.Errors, t)
	assertError(Error{Message: CYCLE_REFERENCE, Type: CycleError}, fsItem1.Errors, t)
	assertErrorCount(1, fsItem2.Errors, t)
	assertError(Error{Message: CYCLE_REFERENCE, Type: CycleError}, fsItem2.Errors, t)
}

func TestOENoRefsAtAll(t *testing.T) {
	kuMap := make(map[string]*KUItem)
	fsMap := make(map[string]*FSItem)
	oeMap := make(map[string]*OEItem)

	oeItem := &OEItem{
		Id: "oe1",
	}
	oeMap[oeItem.Id] = oeItem

	buildTrees(oeMap, kuMap, fsMap)
	analyzeTrees(oeMap, kuMap, fsMap)

	assertErrorCount(2, oeItem.Errors, t)
	assertNoRelatedKU(oeItem, t)
	assertError(Error{Message: NO_RELATED_KU_ID, Type: MissingReference}, oeItem.Errors, t)
	assertNoRelatedFS(oeItem, t)
	assertError(Error{Message: NO_RELATED_FS_ID, Type: MissingReference}, oeItem.Errors, t)
	assertNoParentLRef(oeItem, t)
	assertNoParentFRef(oeItem, t)
}

func TestOENonExistingRefs(t *testing.T) {
	kuMap := make(map[string]*KUItem)
	fsMap := make(map[string]*FSItem)
	oeMap := make(map[string]*OEItem)

	oeItem := &OEItem{
		Id:        "oe1",
		ParentLId: "oe10",
		ParentFId: "oe20",
		KUId:      "ku1",
		FSId:      "fs1",
	}
	oeMap[oeItem.Id] = oeItem

	buildTrees(oeMap, kuMap, fsMap)
	analyzeTrees(oeMap, kuMap, fsMap)

	assertErrorCount(4, oeItem.Errors, t)
	assertNoRelatedKU(oeItem, t)
	assertError(Error{Message: NON_EXISTING_RELATED_KU, Type: NonExistingReference}, oeItem.Errors, t)
	assertNoRelatedFS(oeItem, t)
	assertError(Error{Message: NON_EXISTING_RELATED_FS, Type: NonExistingReference}, oeItem.Errors, t)
	assertNoParentLRef(oeItem, t)
	assertError(Error{Message: NON_EXISTING_PARENT_L, Type: NonExistingReference}, oeItem.Errors, t)
	assertNoParentFRef(oeItem, t)
	assertError(Error{Message: NON_EXISTING_PARENT_F, Type: NonExistingReference}, oeItem.Errors, t)
}

func TestOEFullPicture(t *testing.T) {
	kuMap := make(map[string]*KUItem)
	fsMap := make(map[string]*FSItem)
	oeMap := make(map[string]*OEItem)

	oeItem := &OEItem{
		Id:        "oe1",
		ParentLId: "oe10",
		ParentFId: "oe20",
		KUId:      "ku1",
		FSId:      "fs1",
	}
	oeMap[oeItem.Id] = oeItem

	oeMap["oe10"] = &OEItem{Id: "oe10"}
	oeMap["oe20"] = &OEItem{Id: "oe20"}

	fsMap["fs1"] = &FSItem{Id: "fs1"}
	kuMap["ku1"] = &KUItem{Id: "ku1"}

	buildTrees(oeMap, kuMap, fsMap)
	analyzeTrees(oeMap, kuMap, fsMap)

	assertErrorCount(0, oeItem.Errors, t)
	assertHasKURef(oeItem, kuMap["ku1"], t)
	assertHasFSRef(oeItem, fsMap["fs1"], t)
	assertHasParentLRef(oeItem, oeMap["oe10"], t)
	assertHasParentFRef(oeItem, oeMap["oe20"], t)
}

func TestOECycle(t *testing.T) {
	kuMap := make(map[string]*KUItem)
	fsMap := make(map[string]*FSItem)
	oeMap := make(map[string]*OEItem)

	fsMap["fs1"] = &FSItem{Id: "fs1"}
	kuMap["ku1"] = &KUItem{Id: "ku1"}

	oeItem1 := &OEItem{
		Id:        "oe1",
		ParentLId: "oe2",
		KUId:      "ku1",
		FSId:      "fs1",
	}
	oeMap[oeItem1.Id] = oeItem1

	oeItem2 := &OEItem{
		Id:        "oe2",
		ParentFId: "oe3",
		KUId:      "ku1",
		FSId:      "fs1",
	}
	oeMap[oeItem2.Id] = oeItem2

	oeItem3 := &OEItem{
		Id:        "oe3",
		ParentLId: "oe1",
		KUId:      "ku1",
		FSId:      "fs1",
	}
	oeMap[oeItem3.Id] = oeItem3

	buildTrees(oeMap, kuMap, fsMap)
	analyzeTrees(oeMap, kuMap, fsMap)

	assertErrorCount(1, oeItem1.Errors, t)
	assertError(Error{Message: CYCLE_REFERENCE, Type: CycleError}, oeItem1.Errors, t)
	assertErrorCount(1, oeItem2.Errors, t)
	assertError(Error{Message: CYCLE_REFERENCE, Type: CycleError}, oeItem2.Errors, t)
	assertErrorCount(1, oeItem3.Errors, t)
	assertError(Error{Message: CYCLE_REFERENCE, Type: CycleError}, oeItem3.Errors, t)
}
