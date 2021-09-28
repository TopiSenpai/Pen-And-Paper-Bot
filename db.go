package main

import (
	"encoding/json"
	"github.com/DisgoOrg/log"
	"os"
)

const dbFile = "sounds.json"

func loadSounds() {
	file, err := os.Open(dbFile)
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		log.Errorf("failed to open %s. err: %s", dbFile, err)
		return
	}
	defer file.Close()
	err = json.NewDecoder(file).Decode(&sounds)
	if err != nil {
		log.Errorf("failed to decode sounds from %s. err :%s", dbFile, err)
	}
}

func saveSounds() {
	file, err := os.OpenFile(dbFile, os.O_CREATE, os.ModePerm)
	if err != nil {
		if os.IsNotExist(err) {
			file, err = os.Create(dbFile)
			if err != nil {
				log.Errorf("failed to create %s. err: ", dbFile, err)
			}
		} else {
			log.Errorf("failed to create %s. err: ", dbFile, err)
		}
	}
	defer file.Close()
	err = json.NewEncoder(file).Encode(&sounds)
	if err != nil {
		log.Errorf("failed to encode sounds into %s. err :%s", dbFile, err)
	}
}
