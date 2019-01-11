package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/bayugyug/rest-api-booking/config"
	"github.com/bayugyug/rest-api-booking/models"
	"github.com/bayugyug/rest-api-booking/utils"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
)

const (
	MsgLoginOkay = "Login successful"
	MsgStatusOK  = "Success"
	MsgStatusNOK = "Error"
)

type APIResponse struct {
	Code   int
	Status string
	Result interface{}
}

type ApiHandler struct {
}

func (api *ApiHandler) IndexPage(w http.ResponseWriter, r *http.Request) {
	//reply
	render.JSON(w, r,
		map[string]string{
			"Greeting": "Welcome!",
		})
}

func (api *ApiHandler) CreateCustomer(w http.ResponseWriter, r *http.Request) {

	var err error
	data := &models.Customer{}
	if err = render.Bind(r, data); err != nil {
		log.Println("BIND_FAILED", err)
		//206
		api.ReplyErrContent(w, r, http.StatusPartialContent, http.StatusText(http.StatusPartialContent))
		return
	}
	defer r.Body.Close()
	if data.Type != models.UserTypeCustomer ||
		data.Pass == "" || data.Firstname == "" ||
		data.Lastname == "" || data.Mobile == "" {
		log.Println("MISSING_REQUIRED_PARAMS", data.Type)
		//206
		api.ReplyErrContent(w, r, http.StatusPartialContent, http.StatusText(http.StatusPartialContent))
	}
	if data.Pass == "" {
	}
	log.Println(fmt.Sprintf("%+#v", data))

	//reply
	render.JSON(w, r,
		map[string]string{
			"status": "Ok",
		})
}

func (api *ApiHandler) UpdateCustomer(w http.ResponseWriter, r *http.Request) {
	log.Println("TOKEN: ", api.GetAuthToken(r))
	//reply
	render.JSON(w, r,
		map[string]string{
			"status": "Ok",
		})
}

func (api *ApiHandler) GetCustomer(w http.ResponseWriter, r *http.Request) {
	log.Println("TOKEN: ", api.GetAuthToken(r))
	mobile := strings.TrimSpace(chi.URLParam(r, "id"))
	log.Println("get", mobile)

	//get 1
	data := &models.Customer{}
	usr, err := data.GetCustomer(ApiService.Context, ApiService.DB, mobile)

	//sanity
	if err != nil {
		log.Println(err)
		//403
		api.ReplyErrContent(w, r, http.StatusForbidden, "Invalid token")
		return
	}

	//reply
	render.JSON(w, r, APIResponse{
		Code:   http.StatusOK,
		Status: MsgStatusOK,
		Result: usr,
	})
}

func (api *ApiHandler) DeleteCustomer(w http.ResponseWriter, r *http.Request) {
	log.Println("TOKEN: ", api.GetAuthToken(r))
	//reply
	render.JSON(w, r,
		map[string]string{
			"status": "Ok",
		})
}

func (api *ApiHandler) CreateDriver(w http.ResponseWriter, r *http.Request) {
	log.Println("TOKEN: ", api.GetAuthToken(r))
	//reply
	render.JSON(w, r,
		map[string]string{
			"status": "Ok",
		})
}

func (api *ApiHandler) UpdateDriver(w http.ResponseWriter, r *http.Request) {
	log.Println("TOKEN: ", api.GetAuthToken(r))
	//reply
	render.JSON(w, r,
		map[string]string{
			"status": "Ok",
		})
}

func (api *ApiHandler) GetDriver(w http.ResponseWriter, r *http.Request) {
	log.Println("TOKEN: ", api.GetAuthToken(r))
	id := strings.TrimSpace(chi.URLParam(r, "id"))
	log.Println("get", id)
	//reply
	render.JSON(w, r,
		map[string]string{
			"status": "Ok",
		})
}

func (api *ApiHandler) DeleteDriver(w http.ResponseWriter, r *http.Request) {
	log.Println("TOKEN: ", api.GetAuthToken(r))
	//reply
	render.JSON(w, r,
		map[string]string{
			"status": "Ok",
		})
}

func (api *ApiHandler) GetDriversList(w http.ResponseWriter, r *http.Request) {
	log.Println("TOKEN: ", api.GetAuthToken(r))
	id := strings.TrimSpace(chi.URLParam(r, "address"))
	log.Println("getlist-address", id)
	//reply
	render.JSON(w, r,
		map[string]string{
			"status": "Ok",
		})
}

func (api *ApiHandler) UpdateLocation(w http.ResponseWriter, r *http.Request) {
	log.Println("TOKEN: ", api.GetAuthToken(r))
	//reply
	render.JSON(w, r,
		map[string]string{
			"status": "Ok",
		})
}

func (api *ApiHandler) GetLocation(w http.ResponseWriter, r *http.Request) {
	log.Println("TOKEN: ", api.GetAuthToken(r))
	id := strings.TrimSpace(chi.URLParam(r, "id"))
	who := strings.TrimSpace(chi.URLParam(r, "who"))

	log.Println("get-id", id, who)
	//reply
	render.JSON(w, r,
		map[string]string{
			"status": "Ok",
		})
}

func (api *ApiHandler) CreateBooking(w http.ResponseWriter, r *http.Request) {
	log.Println("TOKEN: ", api.GetAuthToken(r))
	//reply
	render.JSON(w, r,
		map[string]string{
			"status": "Ok",
		})
}

func (api *ApiHandler) UpdateBooking(w http.ResponseWriter, r *http.Request) {
	log.Println("TOKEN: ", api.GetAuthToken(r))
	//reply
	render.JSON(w, r,
		map[string]string{
			"status": "Ok",
		})
}

func (api *ApiHandler) GetBooking(w http.ResponseWriter, r *http.Request) {
	log.Println("TOKEN: ", api.GetAuthToken(r))
	id := strings.TrimSpace(chi.URLParam(r, "id"))
	log.Println("get", id)
	//reply
	render.JSON(w, r,
		map[string]string{
			"status": "Ok",
		})
}

func (api *ApiHandler) GetAddress(w http.ResponseWriter, r *http.Request) {

	data := &models.Location{}
	if err := render.Bind(r, data); err != nil {
		log.Println("BIND_FAILED::INVALID_ADDRESS:", err)
		//204
		api.ReplyErrContent(w, r, http.StatusNoContent, "Invalid required parameters")
		return
	}
	log.Println(data)
	//get token
	bmobile := api.GetAuthToken(r)
	log.Println("TOKEN: ", bmobile)
	if len(bmobile) <= 0 {
		log.Println("INVALID_TOKEN:", bmobile)
		//403
		api.ReplyErrContent(w, r, http.StatusForbidden, "Invalid token")
		return
	}

	if len(data.Address) <= 0 {
		log.Println("INVALID_ADDRESS:", data.Address)
		//204
		api.ReplyErrContent(w, r, http.StatusNoContent, "Invalid required parameters")
		return
	}
	data.Mobile = bmobile
	//try get coords
	ok, location, err := utils.NewGoogleMapGeoCode(config.ApiConfig.GoogleApiKey).GetCoordinates(data.Address)
	if !ok || err != nil {
		log.Println("INVALID_ADDRESS:", data.Address)
		//400
		api.ReplyErrContent(w, r, http.StatusBadRequest, "No coordinates found")
		return
	}
	//reply
	render.JSON(w, r, APIResponse{
		Code:   http.StatusOK,
		Status: "Coordinates found",
		Result: location,
	})
}

func (api *ApiHandler) Login(w http.ResponseWriter, r *http.Request) {

	var err error
	data := &models.UserLogin{}
	if err = render.Bind(r, data); err != nil {
		log.Println(err)
		//403
		api.ReplyErrContent(w, r, http.StatusForbidden, "Invalid token")
		return
	}
	defer r.Body.Close()

	log.Println("Login", data.Mobile, data.Hash, data.Type)
	var usr *models.User
	usr, err = data.GetUserInfo(ApiService.Context, ApiService.DB, data.Type, data.Mobile)

	//sanity
	if err != nil {
		log.Println(err)
		//403
		api.ReplyErrContent(w, r, http.StatusForbidden, "Invalid token")
		return
	}
	//good then check password match
	if data.Hash != usr.Pass {
		log.Println("LOGIN_PASSWORD_MISMATCH")
		//403
		api.ReplyErrContent(w, r, http.StatusForbidden, http.StatusText(http.StatusForbidden))
		return
	}
	//good then check password match
	if usr.Status != "active" {
		log.Println("LOGIN_ACCOUNT_NOT_ACTIVE", usr.Status)
		//403
		api.ReplyErrContent(w, r, http.StatusForbidden, http.StatusText(http.StatusForbidden))
		return
	}
	//generate new token
	token, err := utils.AppJwtToken.GenToken(
		jwt.MapClaims{
			"mobile": usr.Mobile,
			"exp":    jwtauth.ExpireIn(24 * time.Hour),
		},
	)
	if err != nil {
		log.Println("ERROR_TOKEN", err)
		//403
		api.ReplyErrContent(w, r, http.StatusForbidden, "Invalid token")
		return
	}
	//set flag
	_ = data.SetUserLogStatus(ApiService.Context, ApiService.DB, data.Type, data.Mobile, 1)
	//token send
	render.JSON(w, r, APIResponse{
		Code:   http.StatusOK,
		Status: MsgLoginOkay,
		Result: token,
	})
}

func (api ApiHandler) GetAuthToken(r *http.Request) string {
	_, claims, _ := jwtauth.FromContext(r.Context())
	s := claims["mobile"]
	if s != nil {
		return s.(string)
	}
	return ""
}

//ReplyErrContent send 204 msg
//
//  http.StatusNoContent
//  http.StatusText(http.StatusNoContent)
func (api ApiHandler) ReplyErrContent(w http.ResponseWriter, r *http.Request, code int, msg string) {
	render.JSON(w, r, APIResponse{
		Code:   code,
		Status: msg,
	})
}
