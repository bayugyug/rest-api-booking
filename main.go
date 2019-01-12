package main

import (
	"log"
	"time"

	"github.com/bayugyug/rest-api-booking/config"
	"github.com/bayugyug/rest-api-booking/controllers"
)

func main() {
	start := time.Now()
	log.Println("Ver: ", config.ApiVersion)

	var err error

	//init
	config.NewGlobalConfig().InitConfig()
	controllers.ApiService, err = controllers.NewService(
		controllers.WithSvcOptAddress(":"+config.ApiConfig.Port),
		controllers.WithSvcOptDbConf(config.ApiConfig.Driver),
	)

	//check
	if err != nil {
		log.Fatal("Oops! config might be missing", err)
	}

	//run service
	controllers.ApiService.Run()
	log.Println("Since", time.Since(start))
	log.Println("Done")
}
