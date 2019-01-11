package controllers

import (
	"net/http"
	"time"

	"github.com/bayugyug/rest-api-booking/models"
	"github.com/bayugyug/rest-api-booking/utils"
	jwt "github.com/dgrijalva/jwt-go"
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
	Result interface{} `json:",omitempty"`
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

func (api *ApiHandler) Login(w http.ResponseWriter, r *http.Request) {

	var err error
	data := &models.UserLogin{}
	if err = render.Bind(r, data); err != nil {
		utils.Dumper("MISSING_REQUIRED_PARAMETERS", err)
		//204
		api.ReplyErrContent(w, r, http.StatusNoContent, "Invalid required parameters")
		return
	}
	defer r.Body.Close()

	utils.Dumper("Login", data.Mobile, data.Hash, data.Type)
	var usr *models.User
	usr, err = data.GetUserInfo(ApiService.Context, ApiService.DB, data.Type, data.Mobile)

	//sanity
	if err != nil {
		utils.Dumper("LOGIN_FAILED", err)
		//404
		api.ReplyErrContent(w, r, http.StatusNotFound, "Record not found")
		return
	}
	//good then check password match
	if data.Hash != usr.Pass {
		utils.Dumper("LOGIN_PASSWORD_MISMATCH")
		//403
		api.ReplyErrContent(w, r, http.StatusForbidden, "Password mismatch or invalid")
		return
	}
	//good then check password match
	if usr.Status != "active" {
		utils.Dumper("LOGIN_ACCOUNT_NOT_ACTIVE", usr.Status)
		//403
		api.ReplyErrContent(w, r, http.StatusForbidden, "Account is not active")
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
		utils.Dumper("ERROR_TOKEN", err)
		//500
		api.ReplyErrContent(w, r, http.StatusInternalServerError, "Failed to generate token")
		return
	}
	//set flag
	_ = data.SetUserLogStatus(ApiService.Context, ApiService.DB, data.Type, data.Mobile, 1)

	//token send
	render.JSON(w, r, APIResponse{
		Code:   http.StatusOK,
		Status: MsgLoginOkay,
		Result: map[string]interface{}{"Token": token},
	})

}

func (api *ApiHandler) Logout(w http.ResponseWriter, r *http.Request) {
	//NOTE: not yet implemented for token invalidation ;-)
	//reply
	render.JSON(w, r,
		map[string]string{
			"Greeting": "Bye!",
		})
}

func (api ApiHandler) GetAuthToken(r *http.Request) string {
	_, claims, _ := jwtauth.FromContext(r.Context())
	
	//try checking it
	if token, ok := claims["mobile"].(string); ok  {
		return token
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

