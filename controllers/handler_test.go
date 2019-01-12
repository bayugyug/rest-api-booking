package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bayugyug/rest-api-booking/config"
	"github.com/bayugyug/rest-api-booking/utils"
	jwt "github.com/dgrijalva/jwt-go"
)

type bookingInfo APIResponse

func TestSomeHandler(t *testing.T) {
	t.Log("Sanity checking ....")

	config.NewGlobalConfig().InitConfig()

	//init
	ApiService, _ = NewService(
		WithSvcOptAddress(":"+config.ApiConfig.Port),
		WithSvcOptDbConf(config.ApiConfig.Driver),
	)

	utils.ShowMeLog = false
	utils.AppJwtToken = utils.NewAppJwtConfig()
	token1, _ := utils.AppJwtToken.GenToken(jwt.MapClaims{"mobile": "6500000000"})
	token2, _ := utils.AppJwtToken.GenToken(jwt.MapClaims{"mobile": "6500000001"})

	t.Log(token1, token2)
	ts := httptest.NewServer(ApiService.Router)
	defer ts.Close()

	//try
	clearTable()

	mockLists := []struct {
		Method    string
		URL       string
		Bearer    string
		Mobile    string
		Ctx       context.Context
		Body      io.Reader
		URLParams string
	}{
		{
			Method: "POST",
			URL:    "/v1/api/customer",
			Mobile: "6500000000",
			Body:   bytes.NewBufferString(`{"mobile":"6500000000","pass":"8888","latitude":1.308761,"longitude":103.921434,"firstname":"customer","lastname": "lover aguy dabis"}`),
		},
		{
			Method: "POST",
			URL:    "/v1/api/login",
			Body:   bytes.NewBufferString(`{"mobile":"6500000000","pass":"8888","type":"customer"}`),
		},
		{
			Method: "GET",
			URL:    "/v1/api/customer/6500000000",
			Bearer: token1,
		},
		{
			Method: "PUT",
			URL:    "/v1/api/password/customer",
			Bearer: token1,
			Body:   bytes.NewBufferString(`{"mobile":"6500000000","pass":"1234"}`),
		},
		{
			Method: "PUT",
			URL:    "/v1/api/location",
			Bearer: token1,
			Body:   bytes.NewBufferString(`{"mobile":"6500000000","type":"customer","latitude":1.35821,"longitude":103.85615}`),
		},
		{
			Method: "PUT",
			URL:    "/v1/api/customer",
			Bearer: token1,
			Body:   bytes.NewBufferString(`{"mobile":"6500000000","latitude":1.304832,"longitude":103.852855,"firstname":"upd8-customer","lastname": "upd8-dabis"}`),
		},
		{
			Method: "DELETE",
			URL:    "/v1/api/customer",
			Bearer: token1,
			Body:   bytes.NewBufferString(`{"mobile":"6500000000"}`),
		},
		{
			Method: "PUT",
			URL:    "/v1/api/status/customer",
			Bearer: token1,
			Body:   bytes.NewBufferString(`{"mobile":"6500000000","status":"deleted"}`),
		},
		{
			Method: "PUT",
			URL:    "/v1/api/status/customer",
			Bearer: token1,
			Body:   bytes.NewBufferString(`{"mobile":"6500000000","status":"active"}`),
		},
		//DRIVER
		{
			Method: "POST",
			URL:    "/v1/api/driver",
			Mobile: "6500000001",
			Body:   bytes.NewBufferString(`{"mobile":"6500000001","pass":"8888","latitude":1.304832,"longitude":103.852844,"firstname":"driverme","lastname": "aguy dabis"}`),
		},
		{
			Method: "POST",
			URL:    "/v1/api/login",
			Body:   bytes.NewBufferString(`{"mobile":"6500000001","pass":"8888","type":"driver"}`),
			Bearer: token2,
		},
		{
			Method: "GET",
			URL:    "/v1/api/driver/6500000001",
			Body:   bytes.NewBufferString(`{"mobile":"6500000001","pass":"8888","type":"driver"}`),
			Bearer: token2,
		},
		{
			Method: "PUT",
			URL:    "/v1/api/location",
			Body:   bytes.NewBufferString(`{"mobile":"6500000001","type":"driver","latitude":1.35971,"longitude":102.88615}`),
			Bearer: token2,
		},
		{
			Method: "PUT",
			URL:    "/v1/api/driver",
			Body:   bytes.NewBufferString(`{"mobile":"6500000001","latitude":1.303832,"longitude":103.842844,"firstname":"upd8-driverme","lastname": "upd8-aguy dabis"}`),
			Bearer: token2,
		},
		{
			Method: "DELETE",
			URL:    "/v1/api/driver",
			Body:   bytes.NewBufferString(`{"mobile":"6500000001"}`),
			Bearer: token2,
		},
		{
			Method: "PUT",
			URL:    "/v1/api/status/driver",
			Body:   bytes.NewBufferString(`{"mobile":"6500000001","status":"active"}`),
			Bearer: token2,
		},
		{
			Method: "PUT",
			URL:    "/v1/api/password/driver",
			Body:   bytes.NewBufferString(`{"mobile":"6500000001","pass":"1234"}`),
			Bearer: token2,
		},
		//Driver List Within Nearest 50 KM Radius /drivers/{LATITUDE}/{LONGITUDE}
		{
			Method: "GET",
			URL:    "/v1/api/drivers/1.336209/103.737326",
			Bearer: token2,
		},
		{
			Method: "PUT",
			URL:    "/v1/api/vehiclestatus",
			Body:   bytes.NewBufferString(`{"mobile":"6500000001","status":"canceled","latitude":1.35991,"longitude":102.85615}`),
			Bearer: token2,
		},
		{
			Method: "PUT",
			URL:    "/v1/api/vehiclestatus",
			Body:   bytes.NewBufferString(`{"mobile":"6500000001","status":"open","latitude":1.35991,"longitude":102.85615}`),
			Bearer: token2,
		},
		{
			Method: "POST",
			URL:    "/v1/api/booking",
			Body:   bytes.NewBufferString(`{"mobile_customer":"6500000000","src":"kembangan","src_latitude":1.371572,"src_longitude":103.956551,"mobile_driver":"6500000001","dst":"bugis","dst_latitude":1.371572,"dst_longitude":103.956551}`),
			Bearer: token1,
		},
		{
			Method: "GET",
			URL:    "/v1/api/booking/",
			Body:   bytes.NewBufferString(""),
			Bearer: token1,
		},
		/**
		{
			Method: "PUT",
			URL:    "/v1/api/booking/status/customer/", //canceled by customer
			Body:   bytes.NewBufferString(""),
			Bearer: token1,
		},
		{
			Method: "PUT",
			URL:    "/v1/api/booking/status/driver/",
			Body:   bytes.NewBufferString(""),
			Bearer: token2,
			URLParams: "canceled",
		},**/
		{
			Method: "PUT",
			URL:    "/v1/api/booking/pickup-time/",
			Body:   bytes.NewBufferString(""),
			Bearer: token2,
		},
		{
			Method: "PUT",
			URL:    "/v1/api/booking/dropoff-time/",
			Body:   bytes.NewBufferString(""),
			Bearer: token2,
		},
	}

	var bookingID string

	for _, rec := range mockLists {
		var urlStr = rec.URL
		if bookingID != "" {
			urlStr = rec.URL + bookingID
			if rec.URLParams != "" {
				urlStr = rec.URL + bookingID + "/" + rec.URLParams
			}
		}
		ret, body := testRequest(t, ts, rec.Method, urlStr, rec.Body, rec.Bearer)

		//no-otp :-)
		switch rec.URL {
		case "/v1/api/customer":
			byPassOtp(rec.Mobile, "customers")
		case "/v1/api/driver":
			byPassOtp(rec.Mobile, "drivers")
		case "/v1/api/booking":
			bookingID = parseBookingResultId(t, body)
		}
		t.Log("bookid", bookingID, rec.URL+bookingID, body)
		if ret.StatusCode != http.StatusOK {
			t.Fatalf("Request status:%s", ret.StatusCode)
		}

	}

	//try
	clearTable()
}

func clearTable() {
	ApiService.DB.Exec("DELETE FROM customers WHERE mobile='6500000000'")
	ApiService.DB.Exec("DELETE FROM drivers   WHERE mobile='6500000001'")
	ApiService.DB.Exec("DELETE FROM bookings  WHERE mobile_driver IN('6500000001') OR mobile_customer IN('6500000000')")
}

func byPassOtp(mobile, table string) {
	ApiService.DB.Exec("UPDATE " + table + " SET status='active' WHERE mobile='" + mobile + "'")
}

func testRequest(t *testing.T, ts *httptest.Server, method, path string, body io.Reader, auth string) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, body)
	if err != nil {
		t.Fatal(err)
		return nil, ""
	}
	if auth != "" {
		req.Header.Add("Authorization", "Bearer "+auth)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
		return nil, ""
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
		return nil, ""
	}
	defer resp.Body.Close()

	return resp, string(respBody)
}

func parseBookingResultId(t *testing.T, body string) string {
	var booking bookingInfo
	if err := json.NewDecoder(bytes.NewBufferString(body)).Decode(&booking); err != nil {
		t.Fatal("Oops", err, booking)
		return ""
	}

	if booking.Result == nil {
		t.Fatal("Oops nil json", booking.Result)
		return ""
	}

	book, oks := booking.Result.(map[string]interface{})
	if !oks {
		t.Fatal("Oops nil convert", oks)
		return ""
	}
	bval, oks := book["booking"]
	if !oks {
		t.Fatal("Oops nil convert", oks, book)
		return ""
	}
	bid, _ := bval.(float64)
	return fmt.Sprintf("%.0f", bid)
}
