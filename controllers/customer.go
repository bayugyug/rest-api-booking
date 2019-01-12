package controllers

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/bayugyug/rest-api-booking/models"
	"github.com/bayugyug/rest-api-booking/utils"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

func (api *ApiHandler) CreateCustomer(w http.ResponseWriter, r *http.Request) {

	data := &models.Customer{}
	if err := render.Bind(r, data); err != nil {
		utils.Dumper("BIND_FAILED", err)
		//206
		api.ReplyErrContent(w, r, http.StatusPartialContent, http.StatusText(http.StatusPartialContent))
		return
	}
	defer r.Body.Close()
	if data.Type != models.UserTypeCustomer ||
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
	oks, err := data.CreateCustomer(ApiService.Context, ApiService.DB, data)
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

func (api *ApiHandler) UpdateCustomer(w http.ResponseWriter, r *http.Request) {

	token := api.GetAuthToken(r)
	data := &models.Customer{}
	if err := render.Bind(r, data); err != nil {
		utils.Dumper("BIND_FAILED", err)
		//206
		api.ReplyErrContent(w, r, http.StatusPartialContent, http.StatusText(http.StatusPartialContent))
		return
	}
	defer r.Body.Close()
	if data.Type != models.UserTypeCustomer ||
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

	row, err := data.GetCustomer(ApiService.Context, ApiService.DB, data.Mobile)
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
	oks, err := data.UpdateCustomer(ApiService.Context, ApiService.DB, data)
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

func (api *ApiHandler) GetCustomer(w http.ResponseWriter, r *http.Request) {
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
	data := &models.Customer{}
	usr, err := data.GetCustomer(ApiService.Context, ApiService.DB, mobile)

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

func (api *ApiHandler) DeleteCustomer(w http.ResponseWriter, r *http.Request) {

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
	
	data := &models.Customer{}
	row, err := data.GetCustomer(ApiService.Context, ApiService.DB, mobile)
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
	oks, err := data.DeleteCustomer(ApiService.Context, ApiService.DB, mobile)
	if !oks || err != nil {
		utils.Dumper("RECORD_DELETE_FAILED", mobile, err)
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

func (api *ApiHandler) UpdateCustomerPassword(w http.ResponseWriter, r *http.Request) {

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

	user := &models.Customer{}
	row, err := user.GetCustomer(ApiService.Context, ApiService.DB, loc.Mobile)
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
	oks, err := user.UpdateCustomerPass(ApiService.Context, ApiService.DB, row)
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

func (api *ApiHandler) UpdateCustomerStatus(w http.ResponseWriter, r *http.Request) {

	token := api.GetAuthToken(r)
	mobile := strings.TrimSpace(chi.URLParam(r, "id"))
	status := strings.TrimSpace(chi.URLParam(r, "status"))
	utils.Dumper("mobile", mobile, "token", token, "status",status)
	
	//token mismatched
	if mobile!= token || token == "" || mobile == "" {
		utils.Dumper("INVALID_TOKEN:", token, mobile)
		//403
		api.ReplyErrContent(w, r, http.StatusForbidden, "Invalid token")
		return
	}

	user := &models.Customer{}
	row, err := user.GetCustomer(ApiService.Context, ApiService.DB, mobile)
	//sanity
	if err != nil {
		utils.Dumper("RECORD_NOT_FOUND", err)
		//404
		api.ReplyErrContent(w, r, http.StatusNotFound, "Record not found")
		return
	}

	//sanity
	switch status {
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
	if row.Status == status {
		utils.Dumper("STATUS_ALREADY_SET")
		//404
		api.ReplyErrContent(w, r, http.StatusNotFound, "Status already set")
		return
	}

	//update
	oks, err := user.UpdateCustomerStatus(ApiService.Context, ApiService.DB, status, mobile)
	if !oks || err != nil {
		utils.Dumper("UPDATE_FAILED:", mobile)
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
