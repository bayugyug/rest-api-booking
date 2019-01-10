package controllers

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/bayugyug/rest-api-booking/config"
	"github.com/bayugyug/rest-api-booking/driver"
	"github.com/bayugyug/rest-api-booking/utils"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
)

const (
	svcOptionWithHandler  = "svc-opts-handler"
	svcOptionWithAddress  = "svc-opts-address"
	svcOptionWithDbConfig = "svc-opts-db-config"
)

type Service struct {
	Router  *chi.Mux
	Address string
	Api     *ApiHandler
	Config  driver.DbConnectorConfig
	DB      *sql.DB
	Context context.Context
}

//api global handler
var ApiService *Service

//WithSvcOptHandler opts for handler
func WithSvcOptHandler(r *ApiHandler) *config.Option {
	return config.NewOption(svcOptionWithHandler, r)
}

//WithSvcOptAddress opts for port#
func WithSvcOptAddress(r string) *config.Option {
	return config.NewOption(svcOptionWithAddress, r)
}

//WithSvcOptDbConf opts for db connector
func WithSvcOptDbConf(r driver.DbConnectorConfig) *config.Option {
	return config.NewOption(svcOptionWithDbConfig, r)
}

//NewService service new instance
func NewService(opts ...*config.Option) (*Service, error) {

	//default
	svc := &Service{
		Address: ":8989",
		Api:     &ApiHandler{},
		Context: context.Background(),
	}

	//add options if any
	for _, o := range opts {
		switch o.Name() {
		case svcOptionWithHandler:
			if s := o.Value().(*ApiHandler); s != nil {
				svc.Api = s
			}
		case svcOptionWithAddress:
			if s := o.Value().(string); s != "" {
				svc.Address = s
			}
		case svcOptionWithDbConfig:
			s := o.Value().(driver.DbConnectorConfig)
			svc.Config = s
		}
	}

	//set the actual router
	svc.Router = svc.MapRoute()

	//get db
	dbh, err := driver.NewDbConnector(svc.Config)
	if err != nil {
		return svc, err
	}

	//save
	svc.DB = dbh
	return svc, nil
}

//Run run the http server based on settings
func (svc *Service) Run() {

	//gracious timing
	srv := &http.Server{
		Addr:         svc.Address,
		Handler:      svc.Router,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	//async run
	go func() {
		log.Println("Listening on port ", svc.Address)
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("listen: %s\n", err)
			os.Exit(0)
		}

	}()

	//watcher
	stopChan := make(chan os.Signal)
	signal.Notify(stopChan, os.Interrupt)

	<-stopChan
	log.Println("Shutting down service...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	srv.Shutdown(ctx)
	defer cancel()
	log.Println("Server gracefully stopped!")
}

//MapRoute route map all endpoints
func (svc *Service) MapRoute() *chi.Mux {

	// Multiplexer
	router := chi.NewRouter()

	// Basic settings
	router.Use(
		render.SetContentType(render.ContentTypeJSON),
		middleware.Logger,
		middleware.DefaultCompress,
		middleware.StripSlashes,
		middleware.Recoverer,
		middleware.RequestID,
		middleware.RealIP,
	)

	// Basic gracious timing
	router.Use(middleware.Timeout(60 * time.Second))

	// Basic CORS
	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
	})

	router.Use(cors.Handler)

	router.Get("/", svc.Api.IndexPage)

	/*
			@Driver
			POST     /v1/driver/
			PUT      /v1/driver/
			GET      /v1/driver/{mobile}
			DELETE   /v1/driver/{mobile}
			GET      /v1/drivers/


			@Customer
			POST     /v1/customer/
			PUT      /v1/customer/
			GET      /v1/customer/{mobile}
			DELETE   /v1/customer/{mobile}


			@Location
			POST     /v1/location/
			PUT      /v1/location/
			GET      /v1/location/{who}/{id}
		    GET      /v1/address/{who}/{id}

			@Booking
			POST     /v1/booking/
			PUT      /v1/booking/
			GET      /v1/booking/{booking_id}

	*/

	// Protected routes
	router.Route("/v1", func(r chi.Router) {
		r.Use(svc.SetContextKeyVal("api.version", "v1"))
		r.Mount("/api/driver",
			func(api *ApiHandler) *chi.Mux {
				sr := chi.NewRouter()
				sr.Use(jwtauth.Verifier(utils.AppJwtToken.TokenAuth), svc.BearerChecker)
				sr.Put("/", api.UpdateDriver)
				//sr.Post("/", api.CreateDriver)
				sr.Get("/{id}", api.GetDriver)
				sr.Delete("/{id}", api.DeleteDriver)
				return sr
			}(svc.Api))
		r.Mount("/api/drivers",
			func(api *ApiHandler) *chi.Mux {
				sr := chi.NewRouter()
				sr.Use(jwtauth.Verifier(utils.AppJwtToken.TokenAuth), svc.BearerChecker)
				sr.Get("/{address}", api.GetDriversList)
				return sr
			}(svc.Api))
		r.Mount("/api/customer",
			func(api *ApiHandler) *chi.Mux {
				sr := chi.NewRouter()
				sr.Use(jwtauth.Verifier(utils.AppJwtToken.TokenAuth), svc.BearerChecker)
				sr.Put("/", api.UpdateCustomer)
				//sr.Post("/", api.CreateCustomer)
				sr.Get("/{id}", api.GetCustomer)
				sr.Delete("/{id}", api.DeleteCustomer)
				return sr
			}(svc.Api))
		r.Mount("/api/location",
			func(api *ApiHandler) *chi.Mux {
				sr := chi.NewRouter()
				sr.Use(jwtauth.Verifier(utils.AppJwtToken.TokenAuth), svc.BearerChecker)
				sr.Post("/", api.UpdateLocation)
				sr.Get("/{who}/{id}", api.GetLocation)
				return sr
			}(svc.Api))
		r.Mount("/api/booking",
			func(api *ApiHandler) *chi.Mux {
				sr := chi.NewRouter()
				sr.Use(jwtauth.Verifier(utils.AppJwtToken.TokenAuth), svc.BearerChecker)
				sr.Put("/", api.UpdateBooking)
				sr.Post("/", api.CreateBooking)
				sr.Get("/{id}", api.GetBooking)
				return sr
			}(svc.Api))
		r.Mount("/api/address",
			func(api *ApiHandler) *chi.Mux {
				sr := chi.NewRouter()
				sr.Use(jwtauth.Verifier(utils.AppJwtToken.TokenAuth), svc.BearerChecker)
				sr.Get("/{loc}", api.GetAddress)
				return sr
			}(svc.Api))
		/*
			r.Mount("/api/login",
				func(api *ApiHandler) *chi.Mux {
					sr := chi.NewRouter()
					sr.Post("/", api.Login)
					return sr
				}(svc.Api))
		*/

	})

	router.Group(func(r chi.Router) {
		r.Post("/v1/api/login", svc.Api.Login)
		r.Post("/v1/api/customer", svc.Api.CreateCustomer)
		r.Post("/v1/api/driver", svc.Api.CreateDriver)
	})
	return router
}

//SetContextKeyVal version context
func (svc *Service) SetContextKeyVal(k, v string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(context.WithValue(r.Context(), k, v))
			next.ServeHTTP(w, r)
		})
	}
}

//BearerChecker check token
func (svc *Service) BearerChecker(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, dmp, err := jwtauth.FromContext(r.Context())
		log.Println("CLAIMS:", dmp, r.Method, r.URL.Path, r.URL.RawPath)
		if err != nil {
			switch err {
			default:
				log.Println("ERROR:", err)
				svc.Api.ReplyErrContent(w, r, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
				return
			case jwtauth.ErrExpired:
				log.Println("ERROR: Expired")
				http.Error(w, "Expired", http.StatusUnauthorized)
				svc.Api.ReplyErrContent(w, r, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
				return
			case jwtauth.ErrUnauthorized:
				log.Println("ERROR: ErrUnauthorized")
				svc.Api.ReplyErrContent(w, r, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
				return
			}
		}

		if token == nil || !token.Valid {
			svc.Api.ReplyErrContent(w, r, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
			return
		}
		log.Println("TOKEN:", fmt.Sprintf("%#v", token.Raw))
		// Token is authenticated, pass it through
		next.ServeHTTP(w, r)
	})

}