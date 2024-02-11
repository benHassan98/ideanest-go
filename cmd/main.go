package main

import (
	"Ideanest/pkg"
	"Ideanest/pkg/utils"
	"log"
)

func main() {

	config, err := utils.ReadConfig()

	if err != nil {
		log.Fatal("config error:", err)
		return
	}

	log.Println("config:", config)

	if err := pkg.Init(config); err != nil {
		log.Fatal("app error:", err)
		return
	}

}
