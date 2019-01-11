package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/bayugyug/rest-api-booking/config"
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
	pStillRunning    = true
	pBuildTime       = "0"
	pVersion         = "0.1.0" + "-" + pBuildTime
	pShowConsole     = true
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

//initRecov is for dumpIng segv in
func initRecov() {
	//might help u
	defer func() {
		recvr := recover()
		if recvr != nil {
			fmt.Println("MAIN-RECOV-INIT: ", recvr)
		}
	}()
}

//os.Stdout, os.Stdout, os.Stderr
func initLogger(i, w, e io.Writer) {
	//just in case
	if !pShowConsole {
		infoLog = makeLogger(i, pLogDir, "imgur", "INFO: ")
		warnLog = makeLogger(w, pLogDir, "imgur", "WARN: ")
		errorLog = makeLogger(e, pLogDir, "imgur", "ERROR: ")
	} else {
		infoLog = log.New(i,
			"INFO: ",
			log.Ldate|log.Ltime|log.Lmicroseconds)
		warnLog = log.New(w,
			"WARN: ",
			log.Ldate|log.Ltime|log.Lshortfile)
		errorLog = log.New(e,
			"ERROR: ",
			log.Ldate|log.Ltime|log.Lshortfile)
	}
}

//initEnvParams enable all OS envt vars to reload internally
func initEnvParams() {
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
func initConfig() {
	initEnvParams()
	utils.HttpInit()
	//try to reconfigure if there is passed params, otherwise use the default
	if pParamsAppConfig != "" {
		config.ApiConfig = config.FormatAppConfig(pParamsAppConfig)
	}

	//jwt
	utils.AppJwtToken = utils.NewAppJwtConfig()

	//load defaults
	if config.ApiConfig == nil {
		config.ApiConfig = &config.AppConfig{
			Driver: driver.DbConnectorConfig{
				User: "restapi",
				Pass: "r3stapi",
				Host: "127.0.0.1",
				Port: 3306,
				Name: "restapi",
			},
			GoogleApiKey: utils.GoogleApiKey,
			Port:         "8989",
		}
	}
}

//formatLogger try to init all filehandles for logs
func formatLogger(fdir, fname, pfx string) string {
	t := time.Now()
	r := regexp.MustCompile("[^a-zA-Z0-9]")
	p := t.Format("2006-01-02") + "-" + r.ReplaceAllString(strings.ToLower(pfx), "")
	s := path.Join(pLogDir, fdir)
	if _, err := os.Stat(s); os.IsNotExist(err) {
		//mkdir -p
		os.MkdirAll(s, os.ModePerm)
	}
	return path.Join(s, p+"-"+fname+".log")
}

//makeLogger initialize the logger either via file or console
func makeLogger(w io.Writer, ldir, fname, pfx string) *log.Logger {
	logFile := w
	if !pShowConsole {
		var err error
		logFile, err = os.OpenFile(formatLogger(ldir, fname, pfx), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0664)
		if err != nil {
			log.Println(err)
		}
	}
	//give it
	return log.New(logFile,
		pfx,
		log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)

}

//showUsage
func showUsage() {
	fmt.Println("Version:", pVersion)
	flag.PrintDefaults()
	os.Exit(0)
}
