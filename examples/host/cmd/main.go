package main

import (
	"log"

	"github.com/filanov/stateswitch/examples/host/hardware"
	"github.com/filanov/stateswitch/examples/host/host"
	"github.com/filanov/stateswitch/examples/host/models"
	"github.com/go-openapi/swag"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/sirupsen/logrus"
)

func main() {
	db, err := gorm.Open("sqlite3", ":memory:")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	//db = db.Debug()
	if err := db.AutoMigrate(&models.Host{}).Error; err != nil {
		log.Fatal(err)
	}

	hwValidator := hardware.New()
	hapi := host.New(db, hwValidator)

	h1 := &models.Host{ID: uuid.New()}
	h2 := &models.Host{ID: uuid.New()}
	h3 := &models.Host{ID: uuid.New()}

	if err := hapi.Register(h1); err != nil {
		log.Fatal(err)
	}
	if err := hapi.Register(h2); err != nil {
		log.Fatal(err)
	}
	if err := hapi.Register(h3); err != nil {
		log.Fatal(err)
	}

	logHosts := func(s string) {
		logrus.Info(s)
		hosts, err := hapi.List()
		if err != nil {
			logrus.Fatal(err)
		}
		for _, h := range hosts {
			logrus.Infof("id: %s status: %s hw: %t", h.ID, h.Status, swag.BoolValue(h.HwInfo))
		}
	}

	logHosts("Before Changes")

	if err := hapi.SetHwInfo(h1, true); err != nil {
		log.Fatal(err)
	}
	if err := hapi.SetHwInfo(h2, true); err != nil {
		log.Fatal(err)
	}
	if err := hapi.SetHwInfo(h3, false); err != nil {
		log.Fatal(err)
	}
	logHosts("After setting Hw Info")

	if err := hapi.Register(h1); err != nil {
		log.Fatal(err)
	}
	if err := hapi.Register(h2); err != nil {
		log.Fatal(err)
	}
	logHosts("After register again")
}
