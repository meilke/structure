package main

import (
	"fmt"
	"os"
)

func exitOnError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func main() {
	fsMap, err := parseFS("/path/to/XML_FS.xml")
	exitOnError(err)

	kuMap, err := parseKU("/path/to/XML_KU.xml")
	exitOnError(err)

	oeItems, oeMap, err := parseOE("/path/to/XML_OE.xml")
	exitOnError(err)

	for _, item := range oeItems {

		kuItem, kuok := kuMap[item.KUId]
		fsItem, fsok := fsMap[item.FSId]
		parentOELItem, lok := oeMap[item.ParentLId]
		parentOEFItem, fok := oeMap[item.ParentFId]
		fmt.Printf("OE %s (%s)\n", item.OrgKZ, item.Id)
		if kuok {
			fmt.Printf("KU %s (%s)\n", kuItem.NameLong, kuItem.Id)
		} else {
			fmt.Println("no KU!")
		}

		if fsok {
			fmt.Printf("FS %s (%s)\n", fsItem.NameLong, fsItem.Id)
		} else {
			fmt.Println("no FS!")
		}

		if lok {
			fmt.Printf("OE (L) %s (%s)\n", parentOELItem.OrgKZ, parentOELItem.Id)
		} else {
			fmt.Println("no OE L!")
		}

		if fok {
			fmt.Printf("OE (F) %s (%s)\n", parentOEFItem.OrgKZ, parentOEFItem.Id)
		} else {
			fmt.Println("no OE F!")
		}

		fmt.Println()
	}
}
