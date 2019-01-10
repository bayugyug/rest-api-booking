package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/bayugyug/rest-api-booking/models"
	"github.com/bayugyug/rest-api-booking/utils"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
)

const (
	MsgLoginOkay = "Login successful"
)

type APIResponse struct {
	Code    int
	Message string
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
	log.Println("TOKEN: ", api.GetAuthToken(r))
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
	id := strings.TrimSpace(chi.URLParam(r, "id"))
	log.Println("get", id)
	//reply
	render.JSON(w, r,
		map[string]string{
			"status": "Ok",
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
	log.Println("TOKEN: ", api.GetAuthToken(r))
	loc := strings.TrimSpace(chi.URLParam(r, "loc"))
	log.Println("location", loc)
	//reply
	render.JSON(w, r,
		map[string]string{
			"status": "Ok",
		})
}

func (api *ApiHandler) Login(w http.ResponseWriter, r *http.Request) {

	var err error
	data := &models.UserLogin{}
	//if err := render.DecodeJSON(r.Body, &data); err != nil {
	if err = render.Bind(r, data); err != nil {
		log.Println(err)
		//203
		api.ReplyErrContent(w, r, http.StatusNonAuthoritativeInfo, http.StatusText(http.StatusNonAuthoritativeInfo))
		return
	}
	defer r.Body.Close()

	log.Println("Login", data.Mobile, data.Hash, data.Type)
	var usr *models.User
	usr, err = data.GetUserInfo(ApiService.Context, ApiService.DB, data.Type, data.Mobile)

	//sanity
	if err != nil {
		log.Println(err)
		//203
		api.ReplyErrContent(w, r, http.StatusNonAuthoritativeInfo, http.StatusText(http.StatusNonAuthoritativeInfo))
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
		//203
		api.ReplyErrContent(w, r, http.StatusNonAuthoritativeInfo, http.StatusText(http.StatusNonAuthoritativeInfo))
		return
	}
	//token send
	render.JSON(w, r,
		map[string]string{
			"Code":    fmt.Sprintf("%d", http.StatusOK),
			"Message": MsgLoginOkay,
			"Token":   token,
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
		Code:    code,
		Message: msg,
	})
}
