package main

type OEError struct {
	*Error
	OE *OEItem
}

type KUError struct {
	*Error
	KU *KUItem
}

type FSError struct {
	*Error
	FS *FSItem
}

type Errors struct {
	OEErrors []*OEError
	KUErrors []*KUError
	FSErrors []*FSError
}

func buildTrees(oeMap map[string]*OEItem, kuMap map[string]*KUItem, fsMap map[string]*FSItem) {
	for _, v := range kuMap {
		if v.ParentId != "" {
			v.Parent = kuMap[v.ParentId]
			if v.Parent != nil {
				v.Parent.Children = append(v.Parent.Children, v)
			}
		}
	}

	for _, v := range fsMap {
		if v.ParentId != "" {
			v.Parent = fsMap[v.ParentId]
			if v.Parent != nil {
				v.Parent.Children = append(v.Parent.Children, v)
			}
		}
	}

	for _, v := range oeMap {
		if v.ParentLId != "" {
			v.ParentL = oeMap[v.ParentLId]
			if v.ParentL != nil {
				v.ParentL.LChildren = append(v.ParentL.LChildren, v)
			}
		}

		if v.ParentFId != "" {
			v.ParentF = oeMap[v.ParentFId]
			if v.ParentF != nil {
				v.ParentF.FChildren = append(v.ParentF.FChildren, v)
			}
		}

		if v.KUId != "" {
			v.KU = kuMap[v.KUId]
			if v.KU != nil {
				v.KU.OE = append(v.KU.OE, v)
			}
		}

		if v.FSId != "" {
			v.FS = fsMap[v.FSId]
			if v.FS != nil {
				v.FS.OE = append(v.FS.OE, v)
			}
		}
	}
}

func addMissingReferenceError(item *ItemWithError, message string) *Error {
	e := &Error{Message: message, Type: MissingReference}
	item.Errors = append(item.Errors, e)
	return e
}

func addNonExistingReferenceError(item *ItemWithError, message string) *Error {
	e := &Error{Message: message, Type: NonExistingReference}
	item.Errors = append(item.Errors, e)
	return e
}

func addKUError(item *KUItem, e *Error, errors *Errors) {
	_ = append(errors.KUErrors, &KUError{KU: item, Error: e})
}

func addFSError(item *FSItem, e *Error, errors *Errors) {
	_ = append(errors.FSErrors, &FSError{FS: item, Error: e})
}

func addOEError(item *OEItem, e *Error, errors *Errors) {
	_ = append(errors.OEErrors, &OEError{OE: item, Error: e})
}

const (
	NO_PARENT_ID            = "no parent id"
	NON_EXISTING_PARENT     = "non-existing parent"
	NO_PARENT_L_ID          = "no parent L id"
	NON_EXISTING_PARENT_L   = "non-existing parent L"
	NO_PARENT_F_ID          = "no parent F id"
	NON_EXISTING_PARENT_F   = "non-existing parent F"
	NO_RELATED_KU_ID        = "no related KU id"
	NON_EXISTING_RELATED_KU = "non-existing related FS"
	NO_RELATED_FS_ID        = "no related FS id"
	NON_EXISTING_RELATED_FS = "non-existing related FS"
)

func analyzeTrees(oeMap map[string]*OEItem, kuMap map[string]*KUItem, fsMap map[string]*FSItem) *Errors {
	result := &Errors{}

	for _, v := range kuMap {
		if v.ParentId != "" {
			if _, ok := kuMap[v.ParentId]; !ok {
				addKUError(v, addNonExistingReferenceError(&v.ItemWithError, NON_EXISTING_PARENT), result)
			}
		}
	}

	for _, v := range fsMap {
		if v.ParentId != "" {
			if _, ok := fsMap[v.ParentId]; !ok {
				addFSError(v, addNonExistingReferenceError(&v.ItemWithError, NON_EXISTING_PARENT), result)
			}
		}
	}

	for _, v := range oeMap {
		if v.ParentLId != "" {
			if _, ok := oeMap[v.ParentLId]; !ok {
				addOEError(v, addNonExistingReferenceError(&v.ItemWithError, NON_EXISTING_PARENT_L), result)
			}
		}

		if v.ParentFId != "" {
			if _, ok := oeMap[v.ParentFId]; !ok {
				addOEError(v, addNonExistingReferenceError(&v.ItemWithError, NON_EXISTING_PARENT_F), result)
			}
		}

		if v.KUId == "" {
			addOEError(v, addMissingReferenceError(&v.ItemWithError, NO_RELATED_KU_ID), result)
		} else if v.KU == nil {
			addOEError(v, addNonExistingReferenceError(&v.ItemWithError, NON_EXISTING_RELATED_KU), result)
		}

		if v.FSId == "" {
			addOEError(v, addMissingReferenceError(&v.ItemWithError, NO_RELATED_FS_ID), result)
		} else if v.FS == nil {
			addOEError(v, addNonExistingReferenceError(&v.ItemWithError, NON_EXISTING_RELATED_FS), result)
		}
	}

	return result
}
