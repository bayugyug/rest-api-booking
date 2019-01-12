package config

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/bayugyug/rest-api-booking/driver"
	"github.com/bayugyug/rest-api-booking/utils"
)

const (
	//status
	StatusInProgress = "in-progress"
	StatusSuccess    = "success"
	StatusFailed     = "failed"
	StatusPending    = "pending"
	StatusComplete   = "complete"
	StatusOnGoing    = "on-going"
	usageConfig      = "use to set the config file parameter with db-userinfos/googlemap-api-key"
)

var (

	//loggers
	infoLog  *log.Logger
	warnLog  *log.Logger
	errorLog *log.Logger

	//signal flag
	pLogDir          = "."
	pBuildTime       = "0"
	ApiVersion       = "0.1.0" + "-" + pBuildTime
	pParamsAppConfig = ""
	pEnvVars         = map[string]*string{
		"REST_API_BOOKING_CONFIG": &pParamsAppConfig,
	}
)

//internal system initialize
func init() {
	//uniqueness
	rand.Seed(time.Now().UnixNano())

}

type GlobalConfig struct {
}

func NewGlobalConfig() *GlobalConfig {

	//NOTE: @populate misc key/values from os ENV
	//
	// os.Getenv("API.NEAREST.DISTANCE")
	// os.Getenv("API.SQL.USER")
	// os.Getenv("API.SQL.PASS")
	// os.Getenv("API.SQL.NAME")
	// os.Getenv("API.SQL.HOST")
	// os.Getenv("API.SQL.PORT")
	// os.Getenv("API.HTTP.PORT")
	//
	return &GlobalConfig{}
}

//initRecov is for dumpIng segv in
func (g *GlobalConfig) InitRecov() {
	//might help u
	defer func() {
		recvr := recover()
		if recvr != nil {
			fmt.Println("MAIN-RECOV-INIT: ", recvr)
		}
	}()
}

//initEnvParams enable all OS envt vars to reload internally
func (g *GlobalConfig) InitEnvParams() {
	//just in-case, over-write from ENV
	for k, v := range pEnvVars {
		if os.Getenv(k) != "" {
			*v = os.Getenv(k)
		}
	}
	//get options
	flag.StringVar(&pParamsAppConfig, "config", pParamsAppConfig, usageConfig)
	flag.Parse()
}

//initConfig set defaults for initial reqmts
func (g *GlobalConfig) InitConfig() {
	g.InitEnvParams()
	utils.HttpInit()
	//try to reconfigure if there is passed params, otherwise use the default
	if pParamsAppConfig != "" {
		ApiConfig = FormatAppConfig(pParamsAppConfig)
	}

	//jwt
	utils.AppJwtToken = utils.NewAppJwtConfig()

	//load defaults
	if ApiConfig == nil {
		ApiConfig = &AppConfig{
			Driver: driver.DbConnectorConfig{
				User: "restapi",
				Pass: "r3stapi",
				Host: "127.0.0.1",
				Port: 3306,
				Name: "restapi",
			},
			GoogleApiKey: utils.GoogleApiKey,
			Port:         "8989",
			Showlog:      true,
		}
	}
	//set it
	utils.ShowMeLog = ApiConfig.Showlog
}
