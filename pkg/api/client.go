package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

type API struct {
	config *Config

	client *http.Client
	signer *Signer

	address *Address
	ddmcUid string
}

func NewAPI(config Config) (*API, error) {
	err := config.check()
	if err != nil {
		return nil, err
	}

	signer, err := NewSigner("./sign.js")
	if err != nil {
		return nil, err
	}

	return &API{
		client: http.DefaultClient,
		signer: signer,
		config: &config,
	}, nil
}

func (api *API) SetUserAgent(ua string) *API {
	if len(ua) > 0 {
		api.config.UserAgent = ua
	}

	return api
}

func (api *API) SetSID(sid string) *API {
	if len(sid) > 0 {
		api.config.SID = sid
	}

	return api
}

func (api *API) SetOpenID(openid string) *API {
	if len(openid) > 0 {
		api.config.OpenID = openid
	}
	return api
}

func (api *API) SetDeviceID(id string) *API {
	if len(id) > 0 {
		api.config.DeviceID = id
	}

	return api
}

func (api *API) SetDeviceToken(token string) *API {
	if len(token) > 0 {
		api.config.DeviceToken = token
	}

	return api
}

func (api *API) SetAddress(address Address) *API {
	api.address = &address
	return api
}

func (api *API) SetDebugTime(time string) *API {
	if len(time) > 0 {
		api.config.DebugTime = time
	}
	return api
}

func (api *API) getTime() string {
	if len(api.config.DebugTime) > 0 {
		return api.config.DebugTime
	}

	return strconv.FormatInt(time.Now().Unix(), 10)
}

func (api *API) getLocation() ([]string, error) {
	if api.address == nil {
		return nil, errors.New("请先使用SetAddress设置地址")
	}

	return []string{
		fmt.Sprint(api.address.Location.Location[0]),
		fmt.Sprint(api.address.Location.Location[1]),
	}, nil
}

func (api *API) UserDetail() (*UserDetail, error) {
	url, err := url.ParseRequestURI("https://sunquan.api.ddxq.mobi/api/v1/user/detail/")
	if err != nil {
		return nil, err
	}

	var query = url.Query()
	query.Set("api_version", api.config.APIVersion)
	query.Set("app_version", api.config.APPVersion)
	query.Set("applet_source", "")
	query.Set("channel", api.config.Channel)
	query.Set("app_client_id", api.config.ClientID)
	url.RawQuery = query.Encode()

	request, err := http.NewRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		return nil, err
	}

	var header = api.newBaseHeader()
	if err != nil {
		return nil, err
	}

	header.Set("host", "sunquan.api.ddxq.mobi")
	request.Header = header

	var detail = new(UserDetail)
	err = api.do(request, nil, detail)
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
	query.Set("api_version", api.config.APIVersion)
	query.Set("app_version", api.config.APPVersion)
	query.Set("applet_source", "")
	query.Set("channel", api.config.Channel)
	query.Set("app_client_id", api.config.ClientID)
	url.RawQuery = query.Encode()

	request, err := http.NewRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		return nil, err
	}
	header, err := api.newHeader()
	if err != nil {
		return nil, err
	}

	header.Set("host", "sunquan.api.ddxq.mobi")
	request.Header = header
	var address = new(UserAddress)
	err = api.do(request, nil, address)
	if err != nil {
		return nil, err
	}

	return address, nil
}

func (api *API) Cart() (*CartInfo, error) {
	if api.address == nil {
		return nil, errors.New("需先使用SetAddress绑定地址信息")
	}

	url, err := url.ParseRequestURI("https://maicai.api.ddxq.mobi/cart/index")
	if err != nil {
		return nil, err
	}

	var query = url.Query()
	query.Set("station_id", api.address.StationInfo.ID)
	query.Set("is_load", "1")
	query.Set("api_version", api.config.APIVersion)
	query.Set("app_version", api.config.APPVersion)
	query.Set("applet_source", "")
	query.Set("channel", api.config.Channel)
	query.Set("app_client_id", api.config.ClientID)
	url.RawQuery = query.Encode()

	request, err := http.NewRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		return nil, err
	}

	header, err := api.newHeader()
	if err != nil {
		return nil, err
	}
	request.Header = header
	var cart = new(CartInfo)
	err = api.do(request, nil, cart, true)
	if err != nil {
		return nil, err
	}

	return cart, nil
}

func (api *API) GetMultiReserveTime(products []ProductListItem) (*MultiReserveTime, error) {
	if api.address == nil {
		return nil, errors.New("需先使用SetAddress绑定地址信息")
	}

	data, err := json.Marshal([]interface{}{products})
	if err != nil {
		return nil, err
	}

	params := api.newURLEncodedForm()
	params.Set("station_id", api.address.StationInfo.ID)
	params.Set("address_id", api.address.ID)
	params.Set("group_config_id", ``)
	params.Set("products", string(data))
	params.Set("isBridge", `false`)

	url, err := url.ParseRequestURI("https://maicai.api.ddxq.mobi/order/getMultiReserveTime")
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(http.MethodPost, url.String(), nil)
	if err != nil {
		return nil, err
	}

	header, err := api.newHeader()
	if err != nil {
		return nil, err
	}

	request.Header = header
	var times = new(MultiReserveTime)
	err = api.do(request, params, times)
	if err != nil {
		return nil, err
	}

	return times, nil
}

func (api *API) UpdateCheck(productId string, cartId string) (*CartInfo, error) {

	var data = struct {
		ID      string     `json:"id"`
		CartID  string     `json:"cart_id"`
		IsCheck bool       `json:"is_check"`
		Sizes   []struct{} `json:"sizes"`
	}{
		ID:      productId,
		CartID:  cartId,
		IsCheck: true,
		Sizes:   []struct{}{},
	}

	packagesData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	var urlForm = api.newURLEncodedForm()
	urlForm.Set("product", string(packagesData))
	urlForm.Set("is_load", "1")
	urlForm.Set("ab_config", `{"key_onion":"D","key_cart_discount_price":"C"}`)

	url, err := url.ParseRequestURI("https://maicai.api.ddxq.mobi/cart/updateCheck")
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(http.MethodPost, url.String(), nil)
	if err != nil {
		return nil, err
	}

	header, err := api.newHeader()
	if err != nil {
		return nil, err
	}

	request.Header = header
	var cart = new(CartInfo)
	err = api.do(request, urlForm, cart)
	if err != nil {
		return nil, err
	}

	return cart, nil
}

func (api *API) CheckOrder(productList ProductList, useBalance bool) (*CheckOrder, error) {
	if len(productList.Products) == 0 {
		return nil, errors.New("没有可购买商品")
	}

	for i := range productList.Products {
		product := &productList.Products[i]
		product.TotalOriginMoney = product.TotalOriginPrice
		product.TotalMoney = product.TotalPrice
	}

	type ReservedTime struct {
		ReservedTimeStart *int64 `json:"reserved_time_start"`
		ReservedTimeEnd   *int64 `json:"reserved_time_end"`
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

	var params = api.newURLEncodedForm()
	params.Set("user_ticket_id", "default")
	params.Set("freight_ticket_id", "default")
	params.Set("is_use_point", "0")
	params.Set("is_use_balance", "0")
	if useBalance {
		params.Set("is_use_balance", "1")
	}
	params.Set("is_buy_vip", "0")
	params.Set("coupons_id", "")
	params.Set("is_buy_coupons", "0")
	params.Set("packages", string(packagesData))

	url, err := url.ParseRequestURI("https://maicai.api.ddxq.mobi/order/checkOrder")
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(http.MethodPost, url.String(), nil)
	if err != nil {
		return nil, err
	}

	header, err := api.newHeader()
	if err != nil {
		return nil, err
	}

	request.Header = header
	var checkOrder = new(CheckOrder)
	err = api.do(request, params, checkOrder, true)
	if err != nil {
		return nil, err
	}

	return checkOrder, nil
}

func (api *API) AddNewOrder(payType int, cartInfo *CartInfo, reserveTime ReserveTime, checkOrder *CheckOrder) (*AddNewOrder, error) {
	if api.address == nil {
		return nil, errors.New("需先使用SetAddress绑定地址信息")
	}

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
		AddressID:            api.address.ID,
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

	goodsRealMoney, _ := strconv.ParseFloat(checkOrder.Order.GoodsRealMoney, 64)
	orderFreight, _ := strconv.ParseFloat(payment.OrderFreight, 64)

	price := strconv.FormatFloat(goodsRealMoney+orderFreight, 'f', 2, 64)
	log.Println("订单总价", checkOrder.Order.TotalMoney)

	payment.Price = price

	if len(payment.FreightDiscountMoney) == 0 {
		payment.FreightDiscountMoney = "0.00"
	}

	var pl = cartInfo.NewOrderProductList[0]
	pl.TotalMoney = checkOrder.Order.TotalMoney
	pl.GoodsRealMoney = checkOrder.Order.GoodsRealMoney
	pl.TotalOriginMoney = checkOrder.Order.GoodsOriginMoney
	pl.InstantRebateMoney = checkOrder.Order.InstantRebateMoney
	pl.UsedBalanceMoney = checkOrder.Order.UsedBalanceMoney
	pl.CanUsedBalanceMoney = checkOrder.Order.CanUsedBalanceMoney

	var pkg = struct {
		ProductList
		ReservedTimeStart    int64  `json:"reserved_time_start"`
		ReservedTimeEnd      int64  `json:"reserved_time_end"`
		EtaTraceID           string `json:"eta_trace_id"`
		SoonArrival          string `json:"soon_arrival"`
		FirstSelectedBigTime int64  `json:"first_selected_big_time"`
		ReceiptWithoutSku    int    `json:"receipt_without_sku"`
	}{
		ProductList:          pl,
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

	var params = api.newURLEncodedForm()

	params.Set("showMsg", "false")
	params.Set("showData", "true")
	params.Set("ab_config", `{"key_onion": "C"}`)
	params.Set("package_order", string(data))

	url, err := url.ParseRequestURI("https://maicai.api.ddxq.mobi/order/addNewOrder")
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(http.MethodPost, url.String(), nil)
	if err != nil {
		return nil, err
	}

	header, err := api.newHeader()
	if err != nil {
		return nil, err
	}

	request.Header = header
	var addNewOrder = new(AddNewOrder)
	err = api.do(request, params, addNewOrder, true)
	if err != nil {
		return nil, err
	}

	return addNewOrder, nil
}

func (api *API) newBaseHeader() http.Header {

	header := http.Header{}
	header.Set("host", "maicai.api.ddxq.mobi")
	header.Set("User-Agent", api.config.UserAgent)
	header.Set("content-type", "application/x-www-form-urlencoded")
	header.Set("Referer", "https://servicewechat.com/wx1e113254eda17715/425/page-frame.html")

	header.Set("ddmc-api-version", api.config.APIVersion)
	header.Set("ddmc-app-client-id", api.config.ClientID)
	header.Set("ddmc-build-version", api.config.APPVersion)
	header.Set("ddmc-channel", api.config.Channel)
	header.Set("ddmc-os-version", "[object Undefined]")

	header.Set("ddmc-ip", "")
	header.Set("ddmc-time", api.getTime())

	header.Set("ddmc-device-id", api.config.DeviceID)
	header.Set("Cookie", api.config.Cookie)

	return header
}

func (api *API) newHeader() (http.Header, error) {
	header := api.newBaseHeader()

	if len(api.ddmcUid) == 0 {
		return nil, errors.New("用户id未设置")
	}

	if len(api.ddmcUid) > 0 {
		header.Set("ddmc-uid", api.ddmcUid)
	}

	if api.address != nil {
		header.Set("ddmc-station-id", api.address.StationInfo.ID)
		header.Set("ddmc-city-number", api.address.StationInfo.CityNumber)

		location, _ := api.getLocation()
		header.Set("ddmc-longitude", location[0])
		header.Set("ddmc-latitude", location[1])
	}

	return header, nil
}

func (api *API) newURLEncodedForm() url.Values {
	var params = url.Values{}

	params.Set("api_version", api.config.APIVersion)
	params.Set("app_version", api.config.APPVersion)
	params.Set("applet_source", ``)
	params.Set("channel", api.config.Channel)
	params.Set("app_client_id", api.config.ClientID)
	params.Set("device_token", api.config.DeviceToken)

	// me
	params.Set("sharer_uid", ``)
	params.Set("s_id", api.config.SID)
	params.Set("openid", api.config.OpenID)
	params.Set("h5_source", ``)
	params.Set("time", api.getTime())

	if len(api.ddmcUid) > 0 {
		params.Set("uid", api.ddmcUid)
	}

	if api.address != nil {
		params.Set("station_id", api.address.StationInfo.ID)
		params.Set("city_number", api.address.StationInfo.CityNumber)
		location, _ := api.getLocation()
		params.Set("longitude", location[0])
		params.Set("latitude", location[1])
	}

	return params
}

func debugMode(debug ...bool) bool {
	return len(debug) > 0 && debug[0]
}

func (api *API) do(req *http.Request, form url.Values, data interface{}, debug ...bool) error {
	if form != nil {
		var m = make(map[string]interface{})
		for k, v := range form {
			m[k] = v[0]
		}

		signResult, err := api.signer.Sign(m)
		if err != nil {
			return err
		}

		form.Set("nars", signResult.Nars)
		form.Set("sesi", signResult.Sesi)

		if debugMode(debug...) {
			log.Println(signResult)
		}

		req.Body = io.NopCloser(strings.NewReader(form.Encode()))
	}

	if debugMode(debug...) {
		for k, v := range req.Header {
			fmt.Printf("%s\t%s\n", k, v[0])
		}
	}

	resp, err := api.client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var response Response
	if err := json.NewDecoder(bytes.NewReader(body)).Decode(&response); err != nil {
		return err
	}

	if !response.Success {
		log.Println(string(body))
		return NewResponseError(response.Code, response.Message)
	} else if debugMode(debug...) {
		log.Println(string(body))
	}

	return json.Unmarshal(response.Data, data)
}
