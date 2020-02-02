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

	log.WithFields(log.Fields{
		"path": *fsPath,
	}).Info("parsing FS data...")
	fsMap, err := parseFS(*fsPath)
	exitOnError(err)
	log.Info("successfully parsed FS data!")

	log.WithFields(log.Fields{
		"path": *kuPath,
	}).Info("parsing KU data...")
	kuMap, err := parseKU(*kuPath)
	exitOnError(err)
	log.Info("successfully parsed KU data!")

	log.WithFields(log.Fields{
		"path": *oePath,
	}).Info("parsing OE data...")
	oeItems, oeMap, err := parseOE(*oePath)
	exitOnError(err)
	log.Info("successfully parsed OE data!")

	log.Info("building trees...")
	buildTrees(oeMap, kuMap, fsMap)
	log.Info("successfully built trees!")
	log.Info("analyzing trees...")
	analyzeTrees(oeMap, kuMap, fsMap)
	log.Info("successfully analyzed trees!")

	var foundError bool

	foundError = false
	for _, item := range oeItems {
		if len(item.Errors) > 0 {
			foundError = true
			for _, e := range item.Errors {
				log.WithFields(log.Fields{
					"id":      item.Id,
					"name":    item.OrgKZ,
					"message": e.Message,
				}).Info("OE with errors")
			}
		}
	}

	if !foundError {
		log.Info("did not find any errors in OE data!")
	}

	foundError = false
	for _, item := range kuMap {
		if len(item.Errors) > 0 {
			foundError = true
			for _, e := range item.Errors {
				log.WithFields(log.Fields{
					"id":      item.Id,
					"name":    item.NameLong,
					"message": e.Message,
				}).Info("KU with errors")
			}
		}
	}

	if !foundError {
		log.Info("did not find any errors in KU data!")
	}

	foundError = false
	for _, item := range fsMap {
		if len(item.Errors) > 0 {
			foundError = true
			for _, e := range item.Errors {
				log.WithFields(log.Fields{
					"id":      item.Id,
					"name":    item.NameLong,
					"message": e.Message,
				}).Info("FS with errors")
			}
		}
	}

	if !foundError {
		log.Info("did not find any errors in FS data!")
	}

}
