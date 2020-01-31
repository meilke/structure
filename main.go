package main

import (
	"flag"
	log "github.com/sirupsen/logrus"
	"os"
)

func exitOnError(err error) {
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("exiting")
		os.Exit(1)
	}
}

var logLevels = map[string]log.Level{
	"panic":   log.PanicLevel,
	"fatal":   log.FatalLevel,
	"error":   log.ErrorLevel,
	"warning": log.WarnLevel,
	"info":    log.InfoLevel,
	"debug":   log.DebugLevel,
	"trace":   log.TraceLevel,
}

func main() {
	fsPath := flag.String("fs", "XML_FS.xml", "path to XML_FS.xml file")
	kuPath := flag.String("ku", "XML_KU.xml", "path to XML_KU.xml file")
	oePath := flag.String("oe", "XML_OE.xml", "path to XML_OE.xml file")
	logLevel := flag.String("log", "info", "log level")
	flag.Parse()

	log.SetOutput(os.Stdout)
	log.SetLevel(logLevels[*logLevel])

	fsMap, err := parseFS(*fsPath)
	exitOnError(err)

	kuMap, err := parseKU(*kuPath)
	exitOnError(err)

	oeItems, oeMap, err := parseOE(*oePath)
	exitOnError(err)

	buildTrees(oeMap, kuMap, fsMap)
	analyzeTrees(oeMap, kuMap, fsMap)

	for _, item := range oeItems {
		if len(item.Errors) > 0 {
			for _, e := range item.Errors {
				log.WithFields(log.Fields{
					"id":      item.Id,
					"name":    item.OrgKZ,
					"message": e.Message,
				}).Info("OE with errors")
			}
		}
	}

	for _, item := range kuMap {
		if len(item.Errors) > 0 {
			for _, e := range item.Errors {
				log.WithFields(log.Fields{
					"id":      item.Id,
					"name":    item.NameLong,
					"message": e.Message,
				}).Info("KU with errors")
			}
		}
	}

	for _, item := range fsMap {
		if len(item.Errors) > 0 {
			for _, e := range item.Errors {
				log.WithFields(log.Fields{
					"id":      item.Id,
					"name":    item.NameLong,
					"message": e.Message,
				}).Info("FS with errors")
			}
		}
	}

}
