package controllers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/bayugyug/rest-api-booking/config"
	"github.com/bayugyug/rest-api-booking/models"
	"github.com/bayugyug/rest-api-booking/utils"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

func (api *ApiHandler) UpdateLocation(w http.ResponseWriter, r *http.Request) {

	token := api.GetAuthToken(r)
	var loc models.Location
	err := json.NewDecoder(r.Body).Decode(&loc)
	if err != nil || loc.Mobile == "" || loc.Type == "" || loc.Latitude == 0 || loc.Longitude == 0 {
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

	switch loc.Type {
	case models.UserTypeCustomer:
		api.UpdateCustomerCoords(w, r, loc)
	case models.UserTypeDriver:
		api.UpdateDriverCoords(w, r, loc)
	default:
		utils.Dumper("MISSING_REQUIRED_PARAMETERS", loc)
		//206
		api.ReplyErrContent(w, r, http.StatusPartialContent, "Missing required parameters")
		return
	}

}

func (api *ApiHandler) GetLocation(w http.ResponseWriter, r *http.Request) {

	token := api.GetAuthToken(r)
	mobile := strings.TrimSpace(chi.URLParam(r, "mobile"))
	who := strings.TrimSpace(chi.URLParam(r, "who"))

	utils.Dumper("get-id", mobile, who, ",tok:", token)

	//token mismatched
	if mobile != token || token == "" || mobile == "" {
		utils.Dumper("INVALID_TOKEN:", token, mobile)
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
		utils.Dumper("MISSING_REQUIRED_PARAMETERS")
		//206
		api.ReplyErrContent(w, r, http.StatusPartialContent, "Missing required parameters")
		return
	}

}

func (api *ApiHandler) UpdateCustomerCoords(w http.ResponseWriter, r *http.Request, loc models.Location) {

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
	row.Latitude = loc.Latitude
	row.Longitude = loc.Longitude
	//update
	oks, err := user.UpdateCustomerCoords(ApiService.Context, ApiService.DB, row)
	if !oks || err != nil {
		utils.Dumper("UPDATE_FAILED:", loc.Latitude, loc.Longitude, loc.Mobile)
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
	row.Latitude = loc.Latitude
	row.Longitude = loc.Longitude
	//update
	oks, err := user.UpdateDriverCoords(ApiService.Context, ApiService.DB, row)
	if !oks || err != nil {
		utils.Dumper("UPDATE_FAILED:", loc.Latitude, loc.Longitude, loc.Mobile)
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

func (api *ApiHandler) GetCustomerCoords(w http.ResponseWriter, r *http.Request, mobile string) {

	user := &models.Customer{}
	row, err := user.GetCustomer(ApiService.Context, ApiService.DB, mobile)
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

	//reply
	render.JSON(w, r, APIResponse{
		Code:   http.StatusOK,
		Status: "Location is found",
		Result: map[string]interface{}{"latitude": row.Latitude, "longitude": row.Longitude},
	})
}

func (api *ApiHandler) GetAddress(w http.ResponseWriter, r *http.Request) {

	data := &models.Location{}
	if err := render.Bind(r, data); err != nil {
		utils.Dumper("BIND_FAILED::INVALID_ADDRESS:", err)
		//204
		api.ReplyErrContent(w, r, http.StatusNoContent, "Invalid required parameters")
		return
	}
	utils.Dumper(data)
	//get token
	bmobile := api.GetAuthToken(r)
	utils.Dumper("TOKEN: ", bmobile)
	if len(bmobile) <= 0 {
		utils.Dumper("INVALID_TOKEN:", bmobile)
		//403
		api.ReplyErrContent(w, r, http.StatusForbidden, "Invalid token")
		return
	}

	if len(data.Address) <= 0 {
		utils.Dumper("INVALID_ADDRESS:", data.Address)
		//204
		api.ReplyErrContent(w, r, http.StatusNoContent, "Invalid required parameters")
		return
	}
	data.Mobile = bmobile
	//try get coords
	ok, location, err := utils.NewGoogleMapGeoCode(config.ApiConfig.GoogleApiKey).GetCoordinates(data.Address)
	if !ok || err != nil {
		utils.Dumper("INVALID_ADDRESS:", data.Address)
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
