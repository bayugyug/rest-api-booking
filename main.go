package main

import (
	"log"

	"github.com/bayugyug/rest-api-booking/config"
	"github.com/bayugyug/rest-api-booking/controllers"
)

func main() {
	log.Println("Start")
	initConfig()
	//init
	controllers.ApiService, _ = controllers.NewService(
		controllers.WithSvcOptAddress(":8989"),
		controllers.WithSvcOptDbConf(config.ApiConfig.Driver),
	)
	//run service
	controllers.ApiService.Run()
	log.Println("Done")
}
