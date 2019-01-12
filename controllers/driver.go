package controllers

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/bayugyug/rest-api-booking/models"
	"github.com/bayugyug/rest-api-booking/utils"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

func (api *ApiHandler) CreateDriver(w http.ResponseWriter, r *http.Request) {
	data := &models.Driver{}
	if err := render.Bind(r, data); err != nil {
		utils.Dumper("BIND_FAILED", err)
		//206
		api.ReplyErrContent(w, r, http.StatusPartialContent, http.StatusText(http.StatusPartialContent))
		return
	}
	defer r.Body.Close()
	if data.Type != models.UserTypeDriver ||
		data.Pass == "" || data.Firstname == "" ||
		data.Lastname == "" || data.Mobile == "" {
		utils.Dumper("MISSING_REQUIRED_PARAMS", data)
		//206
		api.ReplyErrContent(w, r, http.StatusPartialContent, http.StatusText(http.StatusPartialContent))
		return
	}

	//md5
	data.Pass = fmt.Sprintf("%x", md5.Sum([]byte(data.Pass)))

	//exists
	old := data.Exists(ApiService.Context, ApiService.DB, data.Mobile)
	if old > 0 {
		utils.Dumper("RECORD_EXISTS", data.Mobile)
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
		utils.Dumper("RECORD_CREATE_FAILED", data.Mobile, err)
		//400
		api.ReplyErrContent(w, r, http.StatusBadRequest, "Record create failed")
		return
	}

	//reply
	render.JSON(w, r, APIResponse{
		Code:   http.StatusOK,
		Status: "Create successful",
		Result: map[string]interface{}{
			"otp":        data.Otp,
			"otp-expiry": data.OtpExpiry,
			"mobile":     data.Mobile,
			"uid":        data.ID,
		},
	})

}

func (api *ApiHandler) UpdateDriver(w http.ResponseWriter, r *http.Request) {

	token := api.GetAuthToken(r)
	data := &models.Driver{}
	if err := render.Bind(r, data); err != nil {
		utils.Dumper("BIND_FAILED", err)
		//206
		api.ReplyErrContent(w, r, http.StatusPartialContent, http.StatusText(http.StatusPartialContent))
		return
	}
	defer r.Body.Close()
	if data.Type != models.UserTypeDriver ||
		data.Firstname == "" ||
		data.Lastname == "" ||
		data.Mobile == "" {
		utils.Dumper("MISSING_REQUIRED_PARAMS", data)
		//206
		api.ReplyErrContent(w, r, http.StatusPartialContent, http.StatusText(http.StatusPartialContent))
		return
	}

	//token mismatched
	if data.Mobile != token || token == "" || data.Mobile == "" {
		utils.Dumper("INVALID_TOKEN:", token, data.Mobile)
		//403
		api.ReplyErrContent(w, r, http.StatusForbidden, "Invalid token")
		return
	}

	row, err := data.GetDriver(ApiService.Context, ApiService.DB, data.Mobile)
	//sanity
	if err != nil {
		utils.Dumper("RECORD_NOT_FOUND", err)
		//404
		api.ReplyErrContent(w, r, http.StatusNotFound, "Record not found")
		return
	}
	if row.Status != models.UserStatusActive {
		utils.Dumper("INVALID_STATUS", err)
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
		utils.Dumper("RECORD_CREATE_FAILED", data.Mobile, err)
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
	utils.Dumper("mobile", mobile, ",token", token)

	//token mismatched
	if mobile != token || token == "" || mobile == "" {
		utils.Dumper("INVALID_TOKEN:", token, mobile)
		//403
		api.ReplyErrContent(w, r, http.StatusForbidden, "Invalid token")
		return
	}
	//get 1
	data := &models.Driver{}
	usr, err := data.GetDriver(ApiService.Context, ApiService.DB, mobile)

	//sanity
	if err != nil {
		utils.Dumper(err)
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
		utils.Dumper("BIND_FAILED", err)
		//206
		api.ReplyErrContent(w, r, http.StatusPartialContent, http.StatusText(http.StatusPartialContent))
		return
	}
	defer r.Body.Close()
	
	//token mismatched
	if data.Mobile != token || token == "" || data.Mobile == "" {
		utils.Dumper("INVALID_TOKEN:", token, data.Mobile)
		//403
		api.ReplyErrContent(w, r, http.StatusForbidden, "Invalid token")
		return
	}

	row, err := data.GetDriver(ApiService.Context, ApiService.DB, data.Mobile)
	//sanity
	if err != nil {
		utils.Dumper("RECORD_NOT_FOUND", err)
		//404
		api.ReplyErrContent(w, r, http.StatusNotFound, "Record not found")
		return
	}
	if row.Status != models.UserStatusActive {
		utils.Dumper("INVALID_STATUS", err)
		//206
		api.ReplyErrContent(w, r, http.StatusPartialContent, "Status is not active")
		return

	}
	//delete
	oks, err := data.DeleteDriver(ApiService.Context, ApiService.DB, data.Mobile)
	if !oks || err != nil {
		utils.Dumper("RECORD_DELETE_FAILED", data.Mobile, err)
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

	utils.Dumper("getlist", token, lat, lon)

	//token mismatched
	if token == "" {
		utils.Dumper("INVALID_TOKEN:", token)
		//403
		api.ReplyErrContent(w, r, http.StatusForbidden, "Invalid token")
		return
	}
	if lat == "" || lon == "" {
		utils.Dumper("MISSING_REQUIRED_PARAMS")
		//206
		api.ReplyErrContent(w, r, http.StatusPartialContent, http.StatusText(http.StatusPartialContent))
		return
	}
	latitude, _ := strconv.ParseFloat(lat, 64)
	longitude, _ := strconv.ParseFloat(lon, 64)
	if latitude == 0 || longitude == 0 {
		utils.Dumper("MISSING_REQUIRED_PARAMS::ZERO", latitude, longitude)
		//206
		api.ReplyErrContent(w, r, http.StatusPartialContent, http.StatusText(http.StatusPartialContent))
		return
	}

	row := &models.Driver{}
	drivers, err := row.GetDriversNearestLocation(ApiService.Context, ApiService.DB, latitude, longitude, models.NearestDistance)
	if err != nil {
		utils.Dumper("QUERY_NEAREST_DRIVERS_FAILED:", latitude, longitude)
		//500
		api.ReplyErrContent(w, r, http.StatusInternalServerError, "Nearest drivers look-up failed")
		return
	}

	if len(drivers) <= 0 {
		utils.Dumper("QUERY_NEAREST_DRIVERS_FAILED:", latitude, longitude)
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

func (api *ApiHandler) UpdateDriverPassword(w http.ResponseWriter, r *http.Request) {

	token := api.GetAuthToken(r)
	var loc models.UserLogin
	err := json.NewDecoder(r.Body).Decode(&loc)
	if err != nil || loc.Mobile == "" || loc.Pass == "" {
		utils.Dumper("MISSING_REQUIRED_PARAMETERS", err, loc)
		//206
		api.ReplyErrContent(w, r, http.StatusPartialContent, "Missing required parameters")
		return
	}

	defer r.Body.Close()
	//token mismatched
	if loc.Mobile != token || token == "" || loc.Mobile == "" {
		utils.Dumper("INVALID_TOKEN:", token, loc.Mobile)
		//403
		api.ReplyErrContent(w, r, http.StatusForbidden, "Invalid token")
		return
	}

	user := &models.Driver{}
	row, err := user.GetDriver(ApiService.Context, ApiService.DB, loc.Mobile)
	//sanity
	if err != nil {
		utils.Dumper("RECORD_NOT_FOUND", err)
		//404
		api.ReplyErrContent(w, r, http.StatusNotFound, "Record not found")
		return
	}

	if row.Status != models.UserStatusActive {
		utils.Dumper("INVALID_STATUS", err)
		//206
		api.ReplyErrContent(w, r, http.StatusPartialContent, "Status is not active")
		return

	}
	row.Pass = fmt.Sprintf("%x", md5.Sum([]byte(loc.Pass)))
	//update
	oks, err := user.UpdateDriverPass(ApiService.Context, ApiService.DB, row)
	if !oks || err != nil {
		utils.Dumper("UPDATE_FAILED:", loc.Mobile)
		//500
		api.ReplyErrContent(w, r, http.StatusInternalServerError, "Password update failed")
		return
	}

	//reply
	render.JSON(w, r, APIResponse{
		Code:   http.StatusOK,
		Status: "Password update successful",
	})
}

func (api *ApiHandler) UpdateVehicleStatus(w http.ResponseWriter, r *http.Request) {

	token := api.GetAuthToken(r)
	var loc models.VehicleStatusInfo
	err := json.NewDecoder(r.Body).Decode(&loc)
	if err != nil || loc.Mobile == "" || loc.Status == "" || loc.Latitude == 0 || loc.Longitude == 0 {
		utils.Dumper("MISSING_REQUIRED_PARAMETERS", err, loc)
		//206
		api.ReplyErrContent(w, r, http.StatusPartialContent, "Missing required parameters")
		return
	}
	defer r.Body.Close()
	//token mismatched
	if loc.Mobile != token || token == "" || loc.Mobile == "" {
		utils.Dumper("INVALID_TOKEN:", token, loc.Mobile)
		//403
		api.ReplyErrContent(w, r, http.StatusForbidden, "Invalid token")
		return
	}

	//sanity check
	switch loc.Status {
	case models.VehicleStatusOpen:
	case models.VehicleStatusBooked:
	case models.VehicleStatusCanceled:
	case models.VehicleStatusTripStart:
	case models.VehicleStatusTripEnd:
	case models.VehicleStatusCompleted:
	case models.VehicleStatusGasUp:
	case models.VehicleStatusPanic:
	default:
		utils.Dumper("MISSING_REQUIRED_PARAMETERS", loc)
		//206
		api.ReplyErrContent(w, r, http.StatusPartialContent, "Missing required parameters")
		return
	}

	user := &models.Driver{}
	row, err := user.GetDriver(ApiService.Context, ApiService.DB, loc.Mobile)
	//sanity check
	if err != nil {
		utils.Dumper("RECORD_NOT_FOUND", err)
		//404
		api.ReplyErrContent(w, r, http.StatusNotFound, "Record not found")
		return
	}
	//sanity check
	if row.Status != models.UserStatusActive {
		utils.Dumper("INVALID_STATUS", err)
		//206
		api.ReplyErrContent(w, r, http.StatusPartialContent, "Status is not active")
		return

	}

	//update
	oks, err := user.UpdateDriverVehicleStatus(ApiService.Context, ApiService.DB, loc.Status, loc.Mobile, loc.Latitude, loc.Longitude)
	if !oks || err != nil {
		utils.Dumper("UPDATE_FAILED:", loc.Mobile)
		//500
		api.ReplyErrContent(w, r, http.StatusInternalServerError, "VehicleStatus update failed")
		return
	}

	//reply
	render.JSON(w, r, APIResponse{
		Code:   http.StatusOK,
		Status: "VehicleStatus update successful",
	})
}

func (api *ApiHandler) UpdateDriverStatus(w http.ResponseWriter, r *http.Request) {

	token := api.GetAuthToken(r)
	var loc models.User
	err := json.NewDecoder(r.Body).Decode(&loc)
	if err != nil || loc.Mobile == "" || loc.Status == "" {
		utils.Dumper("MISSING_REQUIRED_PARAMETERS", err, loc)
		//206
		api.ReplyErrContent(w, r, http.StatusPartialContent, "Missing required parameters")
		return
	}
	defer r.Body.Close()
	//token mismatched
	if loc.Mobile != token || token == "" || loc.Mobile == "" {
		utils.Dumper("INVALID_TOKEN:", token, loc.Mobile)
		//403
		api.ReplyErrContent(w, r, http.StatusForbidden, "Invalid token")
		return
	}

	user := &models.Driver{}
	row, err := user.GetDriver(ApiService.Context, ApiService.DB, loc.Mobile)
	//sanity
	if err != nil {
		utils.Dumper("RECORD_NOT_FOUND", err)
		//404
		api.ReplyErrContent(w, r, http.StatusNotFound, "Record not found")
		return
	}

	//sanity
	switch loc.Status {
	case models.UserStatusPending:
	case models.UserStatusActive:
	case models.UserStatusDeleted:
	default:
		utils.Dumper("STATUS_INVALID")
		//404
		api.ReplyErrContent(w, r, http.StatusNotFound, "Invalid status")
		return
	}

	//same
	if row.Status == loc.Status {
		utils.Dumper("STATUS_ALREADY_SET")
		//404
		api.ReplyErrContent(w, r, http.StatusNotFound, "Status already set")
		return
	}

	//update
	oks, err := user.UpdateDriverStatus(ApiService.Context, ApiService.DB, loc.Status, loc.Mobile)
	if !oks || err != nil {
		utils.Dumper("UPDATE_FAILED:", loc.Mobile)
		//500
		api.ReplyErrContent(w, r, http.StatusInternalServerError, "Status update failed")
		return
	}

	//reply
	render.JSON(w, r, APIResponse{
		Code:   http.StatusOK,
		Status: "Status update successful",
	})
}
