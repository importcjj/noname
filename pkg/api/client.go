package api

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

type API struct {
	Cookie string
	client *http.Client

	ddmcUid string
}

func NewAPI(cookie string) (*API, error) {
	if len(cookie) == 0 {
		return nil, errors.New("无效的cookie")
	}
	return &API{
		Cookie: cookie,
		client: http.DefaultClient,
	}, nil
}

func (api *API) UserDetail() (*UserDetail, error) {

	url, err := url.ParseRequestURI("https://sunquan.api.ddxq.mobi/api/v1/user/detail/")
	if err != nil {
		return nil, err
	}

	var query = url.Query()
	query.Set("api_version", "9.50.0")
	query.Set("app_version", "2.83.0")
	query.Set("applet_source", "")
	query.Set("channel", "applet")
	query.Set("app_client_id", "4")
	url.RawQuery = query.Encode()

	request, err := http.NewRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		return nil, err
	}

	var header = api.newHeader()
	header.Set("host", "sunquan.api.ddxq.mobi")
	request.Header = header

	var detail = new(UserDetail)
	err = api.do(request, detail)
	if err != nil {
		return nil, err
	}

	api.ddmcUid = detail.UserInfo.ID

	return detail, nil
}

func (api *API) UserAddress() (*UserAddress, error) {
	url, err := url.ParseRequestURI("https://sunquan.api.ddxq.mobi/api/v1/user/address/")
	if err != nil {
		return nil, err
	}

	var query = url.Query()
	query.Set("api_version", "9.50.0")
	query.Set("app_version", "2.83.0")
	query.Set("applet_source", "")
	query.Set("channel", "applet")
	query.Set("app_client_id", "4")
	url.RawQuery = query.Encode()

	request, err := http.NewRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		return nil, err
	}
	var header = api.newHeader()
	header.Set("host", "sunquan.api.ddxq.mobi")
	request.Header = header
	var address = new(UserAddress)
	err = api.do(request, address)
	if err != nil {
		return nil, err
	}

	return address, nil
}

func (api *API) Cart(stationId string) (*CartInfo, error) {
	url, err := url.ParseRequestURI("https://maicai.api.ddxq.mobi/cart/index")
	if err != nil {
		return nil, err
	}

	var query = url.Query()
	query.Set("station_id", stationId)
	query.Set("is_load", "1")
	query.Set("api_version", "9.50.0")
	query.Set("app_version", "2.83.0")
	query.Set("applet_source", "")
	query.Set("channel", "applet")
	query.Set("app_client_id", "4")
	url.RawQuery = query.Encode()

	request, err := http.NewRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		return nil, err
	}

	var header = api.newHeader()
	request.Header = header
	var cart = new(CartInfo)
	err = api.do(request, cart)
	if err != nil {
		return nil, err
	}

	return cart, nil
}

func (api *API) GetMultiReverseTime(stationId string, products []ProductListItem) (*MultiReserveTime, error) {
	productsData, err := json.Marshal([]interface{}{products})
	if err != nil {
		return nil, err
	}

	var urlForms = url.Values{}
	// urlForms.Set("uid", "")
	// urlForms.Set("longitude", "")
	// urlForms.Set("latitude", "")
	urlForms.Set("station_id", stationId)
	// urlForms.Set("city_number", "")
	urlForms.Set("api_version", "9.50.0")
	urlForms.Set("app_version", "2.83.0")
	urlForms.Set("applet_source", "")
	urlForms.Set("channel", "applet")
	urlForms.Set("app_client_id", "4")
	// urlForms.Set("sharer_uid", "")
	// urlForms.Set("s_id", "")
	// urlForms.Set("openid", "")
	// urlForms.Set("h5_source", "")
	// urlForms.Set("time", "")
	// urlForms.Set("device_token", "")
	// urlForms.Set("address_id", "")
	// urlForms.Set("group_config_id", "")
	urlForms.Set("products", string(productsData))
	// urlForms.Set("isBridge", "false")
	// urlForms.Set("nars", "")
	// urlForms.Set("sesi", "")

	url, err := url.ParseRequestURI("https://maicai.api.ddxq.mobi/order/getMultiReserveTime")
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(http.MethodPost, url.String(), io.NopCloser(strings.NewReader(urlForms.Encode())))
	if err != nil {
		return nil, err
	}

	var header = api.newHeader()
	header.Set("ddmc-station-id", stationId)
	request.Header = header
	var times = new(MultiReserveTime)
	err = api.do(request, times)
	if err != nil {
		return nil, err
	}

	return times, nil
}

func (api *API) CheckOrder(stationId, addressId string, productList ProductList) (*CheckOrder, error) {
	type ReservedTime struct {
		ReservedTimeStart *int64 `json:"reserved_time_start"`
		ReservedTimeEnd   *int64 `json:"reserved_time_end"`
	}

	for i := range productList.Products {
		product := &productList.Products[i]
		product.TotalOriginMoney = product.TotalOriginPrice
		product.TotalMoney = product.TotalPrice
	}

	var data = struct {
		ProductList
		ReservedTime ReservedTime `json:"reserved_time"`
	}{
		ProductList:  productList,
		ReservedTime: ReservedTime{},
	}

	packagesData, err := json.Marshal([]interface{}{data})
	if err != nil {
		return nil, err
	}

	var urlForm = api.newURLEncodedForm()
	urlForm.Set("station_id", stationId)
	urlForm.Set("address_id", addressId)
	urlForm.Set("user_ticket_id", "default")
	urlForm.Set("freight_ticket_id", "default")
	urlForm.Set("is_use_point", "0")
	urlForm.Set("is_use_balance", "0")
	urlForm.Set("is_buy_vip", "0")
	urlForm.Set("coupons_id", "")
	urlForm.Set("is_buy_coupons", "0")
	urlForm.Set("packages", string(packagesData))

	url, err := url.ParseRequestURI("https://maicai.api.ddxq.mobi/order/checkOrder")
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(http.MethodPost, url.String(), io.NopCloser(strings.NewReader(urlForm.Encode())))
	if err != nil {
		return nil, err
	}

	var header = api.newHeader()
	header.Set("ddmc-station-id", stationId)
	request.Header = header
	var checkOrder = new(CheckOrder)
	err = api.do(request, checkOrder)
	if err != nil {
		return nil, err
	}

	return checkOrder, nil
}

func (api *API) AddNewOrder(stationId, addressId string, payType int, cartInfo *CartInfo, reserveTime ReserveTime, checkOrder *CheckOrder) (*AddNewOrder, error) {
	var payment = struct {
		ReservedTimeStart    int64       `json:"reserved_time_start"`
		ReservedTimeEnd      int64       `json:"reserved_time_end"`
		Price                string      `json:"price"`
		FreightDiscountMoney string      `json:"freight_discount_money"`
		FreightMoney         string      `json:"freight_money"`
		OrderFreight         string      `json:"order_freight"`
		ParentOrderSign      string      `json:"parent_order_sign"`
		ProductType          int         `json:"product_type"`
		AddressID            string      `json:"address_id"`
		FormID               string      `json:"form_id"`
		ReceiptWithoutSku    interface{} `json:"receipt_without_sku"`
		PayType              int         `json:"pay_type"`
		UserTicketID         string      `json:"user_ticket_id"`
		VipMoney             string      `json:"vip_money"`
		VipBuyUserTicketID   string      `json:"vip_buy_user_ticket_id"`
		CouponsMoney         string      `json:"coupons_money"`
		CouponsID            string      `json:"coupons_id"`
	}{
		ReservedTimeStart:    reserveTime.StartTimestamp,
		ReservedTimeEnd:      reserveTime.EndTimestamp,
		ParentOrderSign:      cartInfo.ParentOrderInfo.ParentOrderSign,
		AddressID:            addressId,
		PayType:              payType,
		ProductType:          1,
		FormID:               strings.ReplaceAll(uuid.New().String(), "-", ""),
		ReceiptWithoutSku:    nil,
		VipMoney:             "",
		VipBuyUserTicketID:   "",
		CouponsMoney:         "",
		CouponsID:            "",
		Price:                checkOrder.Order.TotalMoney,
		FreightDiscountMoney: checkOrder.Order.FreightDiscountMoney,
		FreightMoney:         checkOrder.Order.FreightMoney,
		OrderFreight:         checkOrder.Order.Freights[0].Freight.FreightRealMoney,
		UserTicketID:         checkOrder.Order.DefaultCoupon.ID,
	}

	var pkg = struct {
		ProductList
		ReservedTimeStart    int64  `json:"reserved_time_start"`
		ReservedTimeEnd      int64  `json:"reserved_time_end"`
		EtaTraceID           string `json:"eta_trace_id"`
		SoonArrival          string `json:"soon_arrival"`
		FirstSelectedBigTime int64  `json:"first_selected_big_time"`
		ReceiptWithoutSku    int    `json:"receipt_without_sku"`
	}{
		ProductList:          cartInfo.NewOrderProductList[0],
		ReservedTimeStart:    reserveTime.StartTimestamp,
		ReservedTimeEnd:      payment.ReservedTimeEnd,
		EtaTraceID:           "",
		SoonArrival:          "",
		FirstSelectedBigTime: 0,
		ReceiptWithoutSku:    0,
	}

	data, err := json.Marshal(map[string]interface{}{
		"payment_order": payment,
		"packages":      []interface{}{pkg},
	})

	if err != nil {
		return nil, err
	}

	var urlForm = api.newURLEncodedForm()
	urlForm.Set("station_id", stationId)
	urlForm.Set("package_order", string(data))

	url, err := url.ParseRequestURI("https://maicai.api.ddxq.mobi/order/addNewOrder")
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(http.MethodPost, url.String(), io.NopCloser(strings.NewReader(urlForm.Encode())))
	if err != nil {
		return nil, err
	}

	var header = api.newHeader()
	header.Set("ddmc-station-id", stationId)
	request.Header = header
	var addNewOrder = new(AddNewOrder)
	err = api.do(request, addNewOrder)
	if err != nil {
		return nil, err
	}

	return addNewOrder, nil
}

func (api *API) newHeader() http.Header {
	header := http.Header{}
	header.Set("host", "maicai.api.ddxq.mobi")
	header.Set("Cookie", api.Cookie)
	header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/53.0.2785.143 Safari/537.36 MicroMessenger/7.0.9.501 NetType/WIFI MiniProgramEnv/Windows WindowsWechat")
	header.Set("content-type", "application/x-www-form-urlencoded")
	header.Set("Referer", "https://servicewechat.com/wx1e113254eda17715/425/page-frame.html")

	header.Set("ddmc-api-version", "9.50.0")
	header.Set("ddmc-app-client-id", "4")
	header.Set("ddmc-app-version", "2.83.0")
	header.Set("ddmc-channel", "applet")
	header.Set("ddmc-time", strconv.FormatInt(time.Now().Unix(), 10))
	if len(api.ddmcUid) > 0 {
		header.Set("ddmc-uid", api.ddmcUid)
	}

	return header
}

func (api *API) newURLEncodedForm() url.Values {
	var urlForm = url.Values{}
	urlForm.Set("api_version", "9.50.0")
	urlForm.Set("app_version", "2.83.0")
	urlForm.Set("applet_source", "")
	urlForm.Set("channel", "applet")
	urlForm.Set("app_client_id", "4")
	urlForm.Set("h5_source", "")
	urlForm.Set("time", strconv.FormatInt(time.Now().Unix(), 10))

	return urlForm
}

func (api *API) do(req *http.Request, data interface{}) error {
	resp, err := api.client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	var response Response
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return err
	}

	if !response.Success {
		return NewResponseError(response.Code, response.Message)
	}

	return json.Unmarshal(response.Data, data)
}
