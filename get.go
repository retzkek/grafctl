package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
	"github.com/retzkek/grafanactl/gapi"
)

var getCmd = &Command{
	Name:    "get",
	Usage:   "[OPTIONS] [DASHBOARD...]",
	Summary: "Retrieve dashboards and save to file.",
	Help: `The get command retrieves dashboards and saves them to file.
If no dashboards are specified, retrieve all available dashboards.`,
}

func getFunc(client *gapi.Client, cmd *Command, args []string) error {
	var dashboards []string
	if len(args) == 0 {
		dashboards = make([]string, 0)
		dl, err := client.ListDashboards()
		if err != nil {
			log.Error(err)
			return fmt.Errorf("error getting dashboard list")
		}
		for _, d := range *dl {
			dashboards = append(dashboards, d.URI)
		}
	} else {
		dashboards = args
	}
	for _, d := range dashboards {
		dash, err := client.Dashboard(d)
		if err != nil {
			log.WithField("dashboard", d).Error(err)
			return fmt.Errorf("error getting dashboard")
		}
		filename := filepath.Join(*path, d) + ".json"
		log.WithFields(log.Fields{
			"dashboard": d,
			"file":      filename,
		}).Info("saving dashboard")
		if err := writeDashboard(dash, filename); err != nil {
			log.WithField("dashboard", d).Error(err)
			return fmt.Errorf("error saving dashboard to file")
		}
	}
	return nil
}

func writeDashboard(dash *gapi.Dashboard, filename string) error {
	ll := log.WithFields(log.Fields{
		"dashboard": dash.Meta.Slug,
		"file":      filename,
		"where":     "writeDashboard",
	})

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	d, err := json.MarshalIndent(dash.Model, "", "\t")
	if err != nil {
		ll.Error(err)
		return fmt.Errorf("error marshalling dashboard to JSON")
	}

	if _, err = f.Write(d); err != nil {
		ll.Error(err)
		return fmt.Errorf("error writing JSON to file")
	}

	ll.Debug("successfully wrote file")
	return nil
}

func init() {
	getCmd.Function = getFunc
}
