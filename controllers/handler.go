package controllers

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
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

func (api *ApiHandler) CreateCustomer(w http.ResponseWriter, r *http.Request) {

	data := &models.Customer{}
	if err := render.Bind(r, data); err != nil {
		log.Println("BIND_FAILED", err)
		//206
		api.ReplyErrContent(w, r, http.StatusPartialContent, http.StatusText(http.StatusPartialContent))
		return
	}
	defer r.Body.Close()
	if data.Type != models.UserTypeCustomer ||
		data.Pass == "" || data.Firstname == "" ||
		data.Lastname == "" || data.Mobile == "" {
		log.Println("MISSING_REQUIRED_PARAMS", data)
		//206
		api.ReplyErrContent(w, r, http.StatusPartialContent, http.StatusText(http.StatusPartialContent))
		return
	}

	//md5
	data.Pass = fmt.Sprintf("%x", md5.Sum([]byte(data.Pass)))
	log.Println(fmt.Sprintf("%+#v", data))

	//exists
	old := data.Exists(ApiService.Context, ApiService.DB, data.Mobile)
	if old > 0 {
		log.Println("RECORD_EXISTS", data.Mobile)
		//409
		api.ReplyErrContent(w, r, http.StatusConflict, "Record already exists")
		return
	}

	//add pin
	now := time.Now().Local()
	data.Otp = fmt.Sprintf("%05d", rand.Intn(9999))
	data.OtpExpiry = now.Add(time.Minute * time.Duration(5)).Format("2006-01-02 15:04:05")
	data.Created = now.Format("2006-01-02 15:04:05")

	//create
	oks, err := data.CreateCustomer(ApiService.Context, ApiService.DB, data)
	if !oks || err != nil {
		log.Println("RECORD_CREATE_FAILED", data.Mobile, err)
		//400
		api.ReplyErrContent(w, r, http.StatusBadRequest, "Record create failed")
		return
	}

	//reply
	render.JSON(w, r, APIResponse{
		Code:   http.StatusOK,
		Status: "Create successful",
	})
}

func (api *ApiHandler) UpdateOtp(w http.ResponseWriter, r *http.Request) {

	var otp models.Otp
	err := json.NewDecoder(r.Body).Decode(&otp)
	if err != nil || otp.Otp == "" || otp.Mobile == "" || otp.Type == "" {
		log.Println("MISSING_REQUIRED_PARAMTERS", err, otp)
		//206
		api.ReplyErrContent(w, r, http.StatusPartialContent, "Missing required parameters")
		return
	}
	switch otp.Type {
	case models.UserTypeCustomer:
		api.UpdateCustomerOtp(w, r, otp)
	case models.UserTypeDriver:
		api.UpdateDriverOtp(w, r, otp)
	default:
		log.Println("MISSING_REQUIRED_PARAMTERS", otp)
		//206
		api.ReplyErrContent(w, r, http.StatusPartialContent, "Missing required parameters")
		return

	}
}

func (api *ApiHandler) UpdateCustomerOtp(w http.ResponseWriter, r *http.Request, otp models.Otp) {

	user := &models.Customer{}
	row, err := user.GetCustomer(ApiService.Context, ApiService.DB, otp.Mobile)
	//sanity
	if err != nil {
		log.Println("RECORD_NOT_FOUND", err)
		//404
		api.ReplyErrContent(w, r, http.StatusNotFound, "Record not found")
		return
	}

	if row.Status != models.UserStatusPending {
		log.Println("INVALID_STATUS", err)
		//206
		api.ReplyErrContent(w, r, http.StatusPartialContent, "Status is not pending")
		return

	}
	//mismatch
	if row.Otp != otp.Otp {
		log.Println("INVALID_OTP", row.Otp, " != ", otp.Otp)
		//401
		api.ReplyErrContent(w, r, http.StatusUnauthorized, "Mismatch Otp")
		return

	}
	//expired
	if row.OtpExpired > 0 {
		log.Println("TIME_EXPIRED", row.OtpExpiry)
		//406
		api.ReplyErrContent(w, r, http.StatusNotAcceptable, "Otp expired")
		return

	}

	//update status:active
	_, _ = user.UpdateCustomerStatus(ApiService.Context, ApiService.DB, models.UserStatusActive, row.Mobile)

	//update otpexpiry
	_, _ = user.UpdateCustomerOtpExpiry(ApiService.Context, ApiService.DB, row)

	//reply
	render.JSON(w, r, APIResponse{
		Code:   http.StatusOK,
		Status: "Otp successful",
	})
}

func (api *ApiHandler) UpdateDriverOtp(w http.ResponseWriter, r *http.Request, otp models.Otp) {

	user := &models.Driver{}
	row, err := user.GetDriver(ApiService.Context, ApiService.DB, otp.Mobile)
	//sanity
	if err != nil {
		log.Println("RECORD_NOT_FOUND", err)
		//404
		api.ReplyErrContent(w, r, http.StatusNotFound, "Record not found")
		return
	}

	if row.Status != models.UserStatusPending {
		log.Println("INVALID_STATUS", err)
		//206
		api.ReplyErrContent(w, r, http.StatusPartialContent, "Status is not pending")
		return

	}
	//mismatch
	if row.Otp != otp.Otp {
		log.Println("INVALID_OTP", row.Otp, " != ", otp.Otp)
		//401
		api.ReplyErrContent(w, r, http.StatusUnauthorized, "Mismatch Otp")
		return

	}
	//expired
	if row.OtpExpired > 0 {
		log.Println("TIME_EXPIRED", row.OtpExpiry)
		//406
		api.ReplyErrContent(w, r, http.StatusNotAcceptable, "Otp expired")
		return

	}

	//update status:active
	_, _ = user.UpdateDriverStatus(ApiService.Context, ApiService.DB, models.UserStatusActive, row.Mobile)

	//update otpexpiry
	_, _ = user.UpdateDriverOtpExpiry(ApiService.Context, ApiService.DB, row)

	//reply
	render.JSON(w, r, APIResponse{
		Code:   http.StatusOK,
		Status: "Otp successful",
	})
}

func (api *ApiHandler) UpdateCustomer(w http.ResponseWriter, r *http.Request) {

	token := api.GetAuthToken(r)
	data := &models.Customer{}
	if err := render.Bind(r, data); err != nil {
		log.Println("BIND_FAILED", err)
		//206
		api.ReplyErrContent(w, r, http.StatusPartialContent, http.StatusText(http.StatusPartialContent))
		return
	}
	defer r.Body.Close()
	if data.Type != models.UserTypeCustomer ||
		data.Firstname == "" ||
		data.Lastname == "" ||
		data.Mobile == "" {
		log.Println("MISSING_REQUIRED_PARAMS", data)
		//206
		api.ReplyErrContent(w, r, http.StatusPartialContent, http.StatusText(http.StatusPartialContent))
		return
	}

	//token mismatched
	if data.Mobile != token || token == "" || data.Mobile == "" {
		log.Println("INVALID_TOKEN:", token, data.Mobile)
		//403
		api.ReplyErrContent(w, r, http.StatusForbidden, "Invalid token")
		return
	}

	row, err := data.GetCustomer(ApiService.Context, ApiService.DB, data.Mobile)
	//sanity
	if err != nil {
		log.Println("RECORD_NOT_FOUND", err)
		//404
		api.ReplyErrContent(w, r, http.StatusNotFound, "Record not found")
		return
	}
	if row.Status != models.UserStatusActive {
		log.Println("INVALID_STATUS", err)
		//206
		api.ReplyErrContent(w, r, http.StatusPartialContent, "Status is not active")
		return

	}
	//optional
	data.Modified = time.Now().Local().Format("2006-01-02 15:04:05")
	if data.Latitude <= 0 {
		data.Latitude = row.Latitude
	}
	if data.Latitude <= 0 {
		data.Longitude = row.Longitude
	}
	oks, err := data.UpdateCustomer(ApiService.Context, ApiService.DB, data)
	if !oks || err != nil {
		log.Println("RECORD_CREATE_FAILED", data.Mobile, err)
		//400
		api.ReplyErrContent(w, r, http.StatusBadRequest, "Record create failed")
		return
	}

	//reply
	render.JSON(w, r, APIResponse{
		Code:   http.StatusOK,
		Status: "Update successful",
	})

}

func (api *ApiHandler) GetCustomer(w http.ResponseWriter, r *http.Request) {
	token := api.GetAuthToken(r)
	mobile := strings.TrimSpace(chi.URLParam(r, "id"))
	log.Println("mobile", mobile, ",token", token)

	//token mismatched
	if mobile != token || token == "" || mobile == "" {
		log.Println("INVALID_TOKEN:", token, mobile)
		//403
		api.ReplyErrContent(w, r, http.StatusForbidden, "Invalid token")
		return
	}
	//get 1
	data := &models.Customer{}
	usr, err := data.GetCustomer(ApiService.Context, ApiService.DB, mobile)

	//sanity
	if err != nil {
		log.Println(err)
		//404
		api.ReplyErrContent(w, r, http.StatusNotFound, "Record not found")
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

	token := api.GetAuthToken(r)
	data := &models.Customer{}
	if err := render.Bind(r, data); err != nil {
		log.Println("BIND_FAILED", err)
		//206
		api.ReplyErrContent(w, r, http.StatusPartialContent, http.StatusText(http.StatusPartialContent))
		return
	}
	defer r.Body.Close()
	if data.Type != models.UserTypeCustomer ||
		data.Mobile == "" {
		log.Println("MISSING_REQUIRED_PARAMS", data)
		//206
		api.ReplyErrContent(w, r, http.StatusPartialContent, http.StatusText(http.StatusPartialContent))
		return
	}

	//token mismatched
	if data.Mobile != token || token == "" || data.Mobile == "" {
		log.Println("INVALID_TOKEN:", token, data.Mobile)
		//403
		api.ReplyErrContent(w, r, http.StatusForbidden, "Invalid token")
		return
	}

	row, err := data.GetCustomer(ApiService.Context, ApiService.DB, data.Mobile)
	//sanity
	if err != nil {
		log.Println("RECORD_NOT_FOUND", err)
		//404
		api.ReplyErrContent(w, r, http.StatusNotFound, "Record not found")
		return
	}
	if row.Status != models.UserStatusActive {
		log.Println("INVALID_STATUS", err)
		//206
		api.ReplyErrContent(w, r, http.StatusPartialContent, "Status is not active")
		return

	}
	//delete
	oks, err := data.DeleteCustomer(ApiService.Context, ApiService.DB, data.Mobile)
	if !oks || err != nil {
		log.Println("RECORD_DELETE_FAILED", data.Mobile, err)
		//400
		api.ReplyErrContent(w, r, http.StatusBadRequest, "Record delete failed")
		return
	}

	//reply
	render.JSON(w, r, APIResponse{
		Code:   http.StatusOK,
		Status: "Delete successful",
	})
}

func (api *ApiHandler) CreateDriver(w http.ResponseWriter, r *http.Request) {
	data := &models.Driver{}
	if err := render.Bind(r, data); err != nil {
		log.Println("BIND_FAILED", err)
		//206
		api.ReplyErrContent(w, r, http.StatusPartialContent, http.StatusText(http.StatusPartialContent))
		return
	}
	defer r.Body.Close()
	if data.Type != models.UserTypeDriver ||
		data.Pass == "" || data.Firstname == "" ||
		data.Lastname == "" || data.Mobile == "" {
		log.Println("MISSING_REQUIRED_PARAMS", data)
		//206
		api.ReplyErrContent(w, r, http.StatusPartialContent, http.StatusText(http.StatusPartialContent))
		return
	}

	//md5
	data.Pass = fmt.Sprintf("%x", md5.Sum([]byte(data.Pass)))
	log.Println(fmt.Sprintf("%+#v", data))

	//exists
	old := data.Exists(ApiService.Context, ApiService.DB, data.Mobile)
	if old > 0 {
		log.Println("RECORD_EXISTS", data.Mobile)
		//409
		api.ReplyErrContent(w, r, http.StatusConflict, "Record already exists")
		return
	}

	//add pin
	now := time.Now().Local()
	data.Otp = fmt.Sprintf("%05d", rand.Intn(9999))
	data.OtpExpiry = now.Add(time.Minute * time.Duration(5)).Format("2006-01-02 15:04:05")
	data.Created = now.Format("2006-01-02 15:04:05")

	//create
	oks, err := data.CreateDriver(ApiService.Context, ApiService.DB, data)
	if !oks || err != nil {
		log.Println("RECORD_CREATE_FAILED", data.Mobile, err)
		//400
		api.ReplyErrContent(w, r, http.StatusBadRequest, "Record create failed")
		return
	}

	//reply
	render.JSON(w, r, APIResponse{
		Code:   http.StatusOK,
		Status: "Create successful",
	})

}

func (api *ApiHandler) UpdateDriver(w http.ResponseWriter, r *http.Request) {

	token := api.GetAuthToken(r)
	data := &models.Driver{}
	if err := render.Bind(r, data); err != nil {
		log.Println("BIND_FAILED", err)
		//206
		api.ReplyErrContent(w, r, http.StatusPartialContent, http.StatusText(http.StatusPartialContent))
		return
	}
	defer r.Body.Close()
	if data.Type != models.UserTypeDriver ||
		data.Firstname == "" ||
		data.Lastname == "" ||
		data.Mobile == "" {
		log.Println("MISSING_REQUIRED_PARAMS", data)
		//206
		api.ReplyErrContent(w, r, http.StatusPartialContent, http.StatusText(http.StatusPartialContent))
		return
	}

	//token mismatched
	if data.Mobile != token || token == "" || data.Mobile == "" {
		log.Println("INVALID_TOKEN:", token, data.Mobile)
		//403
		api.ReplyErrContent(w, r, http.StatusForbidden, "Invalid token")
		return
	}

	row, err := data.GetDriver(ApiService.Context, ApiService.DB, data.Mobile)
	//sanity
	if err != nil {
		log.Println("RECORD_NOT_FOUND", err)
		//404
		api.ReplyErrContent(w, r, http.StatusNotFound, "Record not found")
		return
	}
	if row.Status != models.UserStatusActive {
		log.Println("INVALID_STATUS", err)
		//206
		api.ReplyErrContent(w, r, http.StatusPartialContent, "Status is not active")
		return

	}
	//optional
	data.Modified = time.Now().Local().Format("2006-01-02 15:04:05")
	if data.Latitude <= 0 {
		data.Latitude = row.Latitude
	}
	if data.Latitude <= 0 {
		data.Longitude = row.Longitude
	}
	oks, err := data.UpdateDriver(ApiService.Context, ApiService.DB, data)
	if !oks || err != nil {
		log.Println("RECORD_CREATE_FAILED", data.Mobile, err)
		//400
		api.ReplyErrContent(w, r, http.StatusBadRequest, "Record create failed")
		return
	}

	//reply
	render.JSON(w, r, APIResponse{
		Code:   http.StatusOK,
		Status: "Update successful",
	})

}

func (api *ApiHandler) GetDriver(w http.ResponseWriter, r *http.Request) {
	token := api.GetAuthToken(r)
	mobile := strings.TrimSpace(chi.URLParam(r, "id"))
	log.Println("mobile", mobile, ",token", token)

	//token mismatched
	if mobile != token || token == "" || mobile == "" {
		log.Println("INVALID_TOKEN:", token, mobile)
		//403
		api.ReplyErrContent(w, r, http.StatusForbidden, "Invalid token")
		return
	}
	//get 1
	data := &models.Driver{}
	usr, err := data.GetDriver(ApiService.Context, ApiService.DB, mobile)

	//sanity
	if err != nil {
		log.Println(err)
		//404
		api.ReplyErrContent(w, r, http.StatusNotFound, "Record not found")
		return
	}

	//reply
	render.JSON(w, r, APIResponse{
		Code:   http.StatusOK,
		Status: MsgStatusOK,
		Result: usr,
	})
}

func (api *ApiHandler) DeleteDriver(w http.ResponseWriter, r *http.Request) {
	token := api.GetAuthToken(r)
	data := &models.Driver{}
	if err := render.Bind(r, data); err != nil {
		log.Println("BIND_FAILED", err)
		//206
		api.ReplyErrContent(w, r, http.StatusPartialContent, http.StatusText(http.StatusPartialContent))
		return
	}
	defer r.Body.Close()
	if data.Type != models.UserTypeDriver ||
		data.Mobile == "" {
		log.Println("MISSING_REQUIRED_PARAMS", data)
		//206
		api.ReplyErrContent(w, r, http.StatusPartialContent, http.StatusText(http.StatusPartialContent))
		return
	}

	//token mismatched
	if data.Mobile != token || token == "" || data.Mobile == "" {
		log.Println("INVALID_TOKEN:", token, data.Mobile)
		//403
		api.ReplyErrContent(w, r, http.StatusForbidden, "Invalid token")
		return
	}

	row, err := data.GetDriver(ApiService.Context, ApiService.DB, data.Mobile)
	//sanity
	if err != nil {
		log.Println("RECORD_NOT_FOUND", err)
		//404
		api.ReplyErrContent(w, r, http.StatusNotFound, "Record not found")
		return
	}
	if row.Status != models.UserStatusActive {
		log.Println("INVALID_STATUS", err)
		//206
		api.ReplyErrContent(w, r, http.StatusPartialContent, "Status is not active")
		return

	}
	//delete
	oks, err := data.DeleteDriver(ApiService.Context, ApiService.DB, data.Mobile)
	if !oks || err != nil {
		log.Println("RECORD_DELETE_FAILED", data.Mobile, err)
		//400
		api.ReplyErrContent(w, r, http.StatusBadRequest, "Record delete failed")
		return
	}

	//reply
	render.JSON(w, r, APIResponse{
		Code:   http.StatusOK,
		Status: "Delete successful",
	})

}

func (api *ApiHandler) GetDriversList(w http.ResponseWriter, r *http.Request) {

	lat := strings.TrimSpace(chi.URLParam(r, "lat"))
	lon := strings.TrimSpace(chi.URLParam(r, "lon"))
	token := api.GetAuthToken(r)

	log.Println("getlist", token, lat, lon)

	//token mismatched
	if token == "" {
		log.Println("INVALID_TOKEN:", token)
		//403
		api.ReplyErrContent(w, r, http.StatusForbidden, "Invalid token")
		return
	}
	if lat == "" || lon == "" {
		log.Println("MISSING_REQUIRED_PARAMS")
		//206
		api.ReplyErrContent(w, r, http.StatusPartialContent, http.StatusText(http.StatusPartialContent))
		return
	}
	latitude, _ := strconv.ParseFloat(lat, 64)
	longitude, _ := strconv.ParseFloat(lon, 64)
	if latitude == 0 || longitude == 0 {
		log.Println("MISSING_REQUIRED_PARAMS::ZERO", latitude, longitude)
		//206
		api.ReplyErrContent(w, r, http.StatusPartialContent, http.StatusText(http.StatusPartialContent))
		return
	}

	row := &models.Driver{}
	drivers, err := row.GetDriversNearestLocation(ApiService.Context, ApiService.DB, latitude, longitude, models.NearestDistance)
	if err != nil {
		log.Println("QUERY_NEAREST_DRIVERS_FAILED:", latitude, longitude)
		//500
		api.ReplyErrContent(w, r, http.StatusInternalServerError, "Nearest drivers look-up failed")
		return
	}

	if len(drivers) <= 0 {
		log.Println("QUERY_NEAREST_DRIVERS_FAILED:", latitude, longitude)
		//404
		api.ReplyErrContent(w, r, http.StatusNotFound, "Nearest drivers not found")
		return
	}

	//reply
	render.JSON(w, r, APIResponse{
		Code:   http.StatusOK,
		Status: "Drivers List found",
		Result: map[string]interface{}{"total": len(drivers), "list": drivers},
	})

}

func (api *ApiHandler) UpdateLocation(w http.ResponseWriter, r *http.Request) {

	token := api.GetAuthToken(r)
	var loc models.Location
	err := json.NewDecoder(r.Body).Decode(&loc)
	if err != nil || loc.Mobile == "" || loc.Type == "" || loc.Latitude == 0 || loc.Longitude == 0 {
		utils.Dumper("MISSING_REQUIRED_PARAMTERS", err, loc)
		//206
		api.ReplyErrContent(w, r, http.StatusPartialContent, "Missing required parameters")
		return
	}
	defer r.Body.Close()
	//token mismatched
	if loc.Mobile != token || token == "" || loc.Mobile == "" {
		log.Println("INVALID_TOKEN:", token, loc.Mobile)
		//403
		api.ReplyErrContent(w, r, http.StatusForbidden, "Invalid token")
		return
	}

	switch loc.Type {
	case models.UserTypeCustomer:
		api.UpdateCustomerCoords(w, r, loc)
	case models.UserTypeDriver:
		api.UpdateDriverCoords(w, r, loc)
	default:
		log.Println("MISSING_REQUIRED_PARAMTERS", loc)
		//206
		api.ReplyErrContent(w, r, http.StatusPartialContent, "Missing required parameters")
		return
	}

}

func (api *ApiHandler) GetLocation(w http.ResponseWriter, r *http.Request) {

	token := api.GetAuthToken(r)
	mobile := strings.TrimSpace(chi.URLParam(r, "mobile"))
	who := strings.TrimSpace(chi.URLParam(r, "who"))

	log.Println("get-id", mobile, who, ",tok:", token)

	//token mismatched
	if mobile != token || token == "" || mobile == "" {
		log.Println("INVALID_TOKEN:", token, mobile)
		//403
		api.ReplyErrContent(w, r, http.StatusForbidden, "Invalid token")
		return
	}
	switch who {
	case models.UserTypeCustomer:
		api.GetCustomerCoords(w, r, mobile)
	case models.UserTypeDriver:
		api.GetDriverCoords(w, r, mobile)
	default:
		log.Println("MISSING_REQUIRED_PARAMTERS")
		//206
		api.ReplyErrContent(w, r, http.StatusPartialContent, "Missing required parameters")
		return
	}

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
		//204
		api.ReplyErrContent(w, r, http.StatusNoContent, "Invalid required parameters")
		return
	}
	defer r.Body.Close()

	log.Println("Login", data.Mobile, data.Hash, data.Type)
	var usr *models.User
	usr, err = data.GetUserInfo(ApiService.Context, ApiService.DB, data.Type, data.Mobile)

	//sanity
	if err != nil {
		log.Println(err)
		//204
		api.ReplyErrContent(w, r, http.StatusNoContent, "Invalid required parameters")
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

func (api *ApiHandler) UpdateCustomerCoords(w http.ResponseWriter, r *http.Request, loc models.Location) {

	user := &models.Customer{}
	row, err := user.GetCustomer(ApiService.Context, ApiService.DB, loc.Mobile)
	//sanity
	if err != nil {
		log.Println("RECORD_NOT_FOUND", err)
		//404
		api.ReplyErrContent(w, r, http.StatusNotFound, "Record not found")
		return
	}

	if row.Status != models.UserStatusActive {
		log.Println("INVALID_STATUS", err)
		//206
		api.ReplyErrContent(w, r, http.StatusPartialContent, "Status is not active")
		return

	}
	row.Latitude = loc.Latitude
	row.Longitude = loc.Longitude
	//update
	oks, err := user.UpdateCustomerCoords(ApiService.Context, ApiService.DB, row)
	if !oks || err != nil {
		log.Println("UPDATE_FAILED:", loc.Latitude, loc.Longitude, loc.Mobile)
		//500
		api.ReplyErrContent(w, r, http.StatusInternalServerError, "Location update failed")
		return
	}

	//reply
	render.JSON(w, r, APIResponse{
		Code:   http.StatusOK,
		Status: "Location update successful",
	})
}

func (api *ApiHandler) UpdateDriverCoords(w http.ResponseWriter, r *http.Request, loc models.Location) {

	user := &models.Driver{}
	row, err := user.GetDriver(ApiService.Context, ApiService.DB, loc.Mobile)
	//sanity
	if err != nil {
		log.Println("RECORD_NOT_FOUND", err)
		//404
		api.ReplyErrContent(w, r, http.StatusNotFound, "Record not found")
		return
	}

	if row.Status != models.UserStatusActive {
		log.Println("INVALID_STATUS", err)
		//206
		api.ReplyErrContent(w, r, http.StatusPartialContent, "Status is not active")
		return

	}
	row.Latitude = loc.Latitude
	row.Longitude = loc.Longitude
	//update
	oks, err := user.UpdateDriverCoords(ApiService.Context, ApiService.DB, row)
	if !oks || err != nil {
		log.Println("UPDATE_FAILED:", loc.Latitude, loc.Longitude, loc.Mobile)
		//500
		api.ReplyErrContent(w, r, http.StatusInternalServerError, "Location update failed")
		return
	}

	//reply
	render.JSON(w, r, APIResponse{
		Code:   http.StatusOK,
		Status: "Location update successful",
	})
}

func (api *ApiHandler) UpdatePassword(w http.ResponseWriter, r *http.Request) {

	//206
	api.ReplyErrContent(w, r, http.StatusPartialContent, "Missing required parameters")
}
func (api *ApiHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {

	//206
	api.ReplyErrContent(w, r, http.StatusPartialContent, "Missing required parameters")
}

func (api *ApiHandler) GetCustomerCoords(w http.ResponseWriter, r *http.Request, mobile string) {

	user := &models.Customer{}
	row, err := user.GetCustomer(ApiService.Context, ApiService.DB, mobile)
	//sanity
	if err != nil {
		log.Println("RECORD_NOT_FOUND", err)
		//404
		api.ReplyErrContent(w, r, http.StatusNotFound, "Record not found")
		return
	}

	if row.Status != models.UserStatusActive {
		log.Println("INVALID_STATUS", err)
		//206
		api.ReplyErrContent(w, r, http.StatusPartialContent, "Status is not active")
		return

	}

	//reply
	render.JSON(w, r, APIResponse{
		Code:   http.StatusOK,
		Status: "Location is found",
		Result: map[string]interface{}{"latitude": row.Latitude, "longitude": row.Longitude},
	})
}

func (api *ApiHandler) GetDriverCoords(w http.ResponseWriter, r *http.Request, mobile string) {

	user := &models.Driver{}
	row, err := user.GetDriver(ApiService.Context, ApiService.DB, mobile)
	//sanity
	if err != nil {
		log.Println("RECORD_NOT_FOUND", err)
		//404
		api.ReplyErrContent(w, r, http.StatusNotFound, "Record not found")
		return
	}

	if row.Status != models.UserStatusActive {
		log.Println("INVALID_STATUS", err)
		//206
		api.ReplyErrContent(w, r, http.StatusPartialContent, "Status is not active")
		return

	}

	//reply
	render.JSON(w, r, APIResponse{
		Code:   http.StatusOK,
		Status: "Location is found",
		Result: map[string]interface{}{"latitude": row.Latitude, "longitude": row.Longitude},
	})
}
