package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/bayugyug/rest-api-booking/models"
	"github.com/bayugyug/rest-api-booking/utils"
	"github.com/go-chi/render"
)

func (api *ApiHandler) UpdateOtp(w http.ResponseWriter, r *http.Request) {

	var otp models.Otp
	err := json.NewDecoder(r.Body).Decode(&otp)
	if err != nil || otp.Otp == "" || otp.Mobile == "" || otp.Type == "" {
		utils.Dumper("MISSING_REQUIRED_PARAMETERS", err, otp)
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
		utils.Dumper("MISSING_REQUIRED_PARAMETERS", otp)
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
		utils.Dumper("RECORD_NOT_FOUND", err)
		//404
		api.ReplyErrContent(w, r, http.StatusNotFound, "Record not found")
		return
	}

	if row.Status != models.UserStatusPending {
		utils.Dumper("INVALID_STATUS", err)
		//206
		api.ReplyErrContent(w, r, http.StatusPartialContent, "Status is not pending")
		return

	}
	//mismatch
	if row.Otp != otp.Otp {
		utils.Dumper("INVALID_OTP", row.Otp, " != ", otp.Otp)
		//401
		api.ReplyErrContent(w, r, http.StatusUnauthorized, "Mismatch Otp")
		return

	}
	//expired
	if row.OtpExpired > 0 {
		utils.Dumper("TIME_EXPIRED", row.OtpExpiry)
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
		utils.Dumper("RECORD_NOT_FOUND", err)
		//404
		api.ReplyErrContent(w, r, http.StatusNotFound, "Record not found")
		return
	}

	if row.Status != models.UserStatusPending {
		utils.Dumper("INVALID_STATUS", err)
		//206
		api.ReplyErrContent(w, r, http.StatusPartialContent, "Status is not pending")
		return

	}
	//mismatch
	if row.Otp != otp.Otp {
		utils.Dumper("INVALID_OTP", row.Otp, " != ", otp.Otp)
		//401
		api.ReplyErrContent(w, r, http.StatusUnauthorized, "Mismatch Otp")
		return

	}
	//expired
	if row.OtpExpired > 0 {
		utils.Dumper("TIME_EXPIRED", row.OtpExpiry)
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
