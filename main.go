package main

import (
	"log"

	"github.com/bayugyug/rest-api-booking/config"
	"github.com/bayugyug/rest-api-booking/controllers"
)

func main() {
	log.Println("Start")
	config.NewGlobalConfig().InitConfig()
	//init
	controllers.ApiService, _ = controllers.NewService(
		controllers.WithSvcOptAddress(":"+config.ApiConfig.Port),
		controllers.WithSvcOptDbConf(config.ApiConfig.Driver),
	)
	//run service
	controllers.ApiService.Run()
	log.Println("Done")
}
