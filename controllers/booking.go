package controllers

import (
	"net/http"
	"strings"

	"github.com/bayugyug/rest-api-booking/models"
	"github.com/bayugyug/rest-api-booking/utils"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

func (api *ApiHandler) CreateBooking(w http.ResponseWriter, r *http.Request) {

	token := api.GetAuthToken(r)
	utils.Dumper("TOKEN: ", token)
	data := &models.Booking{}
	if err := render.Bind(r, data); err != nil {
		utils.Dumper("FAILED_BIND:", err)
		//204
		api.ReplyErrContent(w, r, http.StatusNoContent, "Invalid required parameters")
		return
	}

	if data.MobileCustomer == "" || data.MobileDriver == "" ||
		data.Src == "" || data.SrcLatitude == 0 || data.SrcLongitude == 0 ||
		data.Dst == "" || data.DstLatitude == 0 || data.DstLongitude == 0 {
		utils.Dumper("MISSING_REQUIRED_PARAMS", data)
		//206
		api.ReplyErrContent(w, r, http.StatusPartialContent, http.StatusText(http.StatusPartialContent))
		return
	}

	if data.MobileCustomer != token || token == "" || data.MobileCustomer == "" {
		utils.Dumper("INVALID_TOKEN:", token, data.MobileCustomer)
		//403
		api.ReplyErrContent(w, r, http.StatusForbidden, "Invalid token")
		return
	}

	//check if have open trip for driver
	driver := &models.Driver{}
	row, err := driver.GetDriver(ApiService.Context, ApiService.DB, data.MobileDriver)

	//sanity
	if err != nil {
		utils.Dumper("RECORD_NOT_FOUND", err)
		//404
		api.ReplyErrContent(w, r, http.StatusNotFound, "Driver Record not found")
		return
	}

	//check if open
	if row.VehicleStatus != models.VehicleStatusOpen {
		//400
		api.ReplyErrContent(w, r, http.StatusBadRequest, "Driver vehicle status is not valid")
		return
	}

	//active
	if row.Status != models.UserStatusActive {
		//400
		api.ReplyErrContent(w, r, http.StatusBadRequest, "Driver status is not valid")
		return
	}

	//customer
	customer := &models.Customer{}
	prow, err := customer.GetCustomer(ApiService.Context, ApiService.DB, data.MobileCustomer)
	//sanity
	if err != nil {
		utils.Dumper("RECORD_NOT_FOUND", err)
		//404
		api.ReplyErrContent(w, r, http.StatusNotFound, "Customer Record not found")
		return
	}

	//active
	if prow.Status != models.UserStatusActive {
		//400
		api.ReplyErrContent(w, r, http.StatusBadRequest, "Customer status is not valid")
		return
	}
	//create
	data.Status = models.VehicleStatusBooked
	if oks, err := data.CreateBooking(ApiService.Context, ApiService.DB, data); !oks || err != nil {
		utils.Dumper("RECORD_CREATE_FAILED", err)
		//400
		api.ReplyErrContent(w, r, http.StatusBadRequest, "Record create failed")
		return
	}

	//update vehicle-status
	if oks, err := driver.UpdateDriverVehicleStatus(ApiService.Context, ApiService.DB,
		models.VehicleStatusBooked,
		data.MobileDriver,
		data.DstLatitude,
		data.DstLongitude); !oks || err != nil {
		utils.Dumper("RECORD_UPDATE_FAILED", err)
		//400
		api.ReplyErrContent(w, r, http.StatusBadRequest, "Record update vehicle-status failed")
		return
	}

	//reply
	render.JSON(w, r, APIResponse{
		Code:   http.StatusOK,
		Status: "Create successful",
		Result: map[string]interface{}{"booking": data.ID},
	})
}

func (api *ApiHandler) UpdateBooking(w http.ResponseWriter, r *http.Request) {
	utils.Dumper("TOKEN: ", api.GetAuthToken(r))
	//reply
	render.JSON(w, r,
		map[string]string{
			"status": "Ok",
		})
}

func (api *ApiHandler) GetBooking(w http.ResponseWriter, r *http.Request) {

	token := api.GetAuthToken(r)
	id := strings.TrimSpace(chi.URLParam(r, "id"))

	utils.Dumper("get", id, token)

	if token == "" {
		utils.Dumper("INVALID_TOKEN:", token)
		//403
		api.ReplyErrContent(w, r, http.StatusForbidden, "Invalid token")
		return
	}

	data := &models.Booking{}

	//get 1
	book, err := data.GetBooking(ApiService.Context, ApiService.DB, id)

	//sanity
	if err != nil {
		utils.Dumper("RECORD_NOT_FOUND", err)
		//404
		api.ReplyErrContent(w, r, http.StatusNotFound, "Booking Record not found")
		return
	}

	//reply
	render.JSON(w, r, APIResponse{
		Code:   http.StatusOK,
		Status: "Booking found",
		Result: map[string]interface{}{"booking": book},
	})
}
