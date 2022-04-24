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
	signer *Signer

	ddmcUid string
}

func NewAPI(cookie string) (*API, error) {
	if len(cookie) == 0 {
		return nil, errors.New("无效的cookie")
	}

	signer, err := NewSigner("./sign.js")
	if err != nil {
		return nil, err
	}

	return &API{
		Cookie: cookie,
		client: http.DefaultClient,
		signer: signer,
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
	err = api.do(request, nil, address)
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
	err = api.do(request, nil, cart)
	if err != nil {
		return nil, err
	}

	return cart, nil
}

func (api *API) GetMultiReverseTime(stationId, addressId string, products []ProductListItem) (*MultiReserveTime, error) {
	_, err := json.Marshal([]interface{}{products})
	if err != nil {
		return nil, err
	}

	params := api.newURLEncodedForm()
	params.Add("address_id", addressId)
	params.Add("group_config_id", ``)
	params.Add("products", `[[{"type":1,"id":"612cc0982c34fab505117d4e","price":"828.00","count":1,"description":"","sizes":[],"cart_id":"612cc0982c34fab505117d4e","parent_id":"","parent_batch_type":-1,"category_path":"","manage_category_path":"411,412,413","activity_id":"","sku_activity_id":"","conditions_num":"","product_name":"洋河蓝色经典梦之蓝M6+52度白酒 550ml/瓶","product_type":0,"small_image":"https://ddfs-public.ddimg.mobi/img/blind/product-management/202108/1242efbb2a37470aa081683513fb3677.jpg?width=800&height=800","total_price":"828.00","origin_price":"828.00","total_origin_price":"828.00","no_supplementary_price":"828.00","no_supplementary_total_price":"828.00","size_price":"0.00","buy_limit":0,"price_type":0,"promotion_num":0,"instant_rebate_money":"0.00","is_invoice":1,"sub_list":[],"is_booking":0,"is_bulk":0,"view_total_weight":"瓶","net_weight":"550","net_weight_unit":"ml","storage_value_id":0,"temperature_layer":"","sale_batches":{"batch_type":-1},"is_shared_station_product":0,"is_gift":0,"supplementary_list":[],"order_sort":1,"is_presale":0}]]`)
	params.Add("isBridge", `false`)
	// params.Add("nars", `f109a4692e2ce8404d3a2a96f0a3b199`)
	// params.Add("sesi", `RYx3k3E6cacc851110cf3c3063dc18b383f75aa`)

	url, err := url.ParseRequestURI("https://maicai.api.ddxq.mobi/order/getMultiReserveTime")
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(http.MethodPost, url.String(), nil)
	if err != nil {
		return nil, err
	}

	var header = api.newHeader()
	header.Set("ddmc-station-id", stationId)

	request.Header = header
	var times = new(MultiReserveTime)
	err = api.do(request, params, times)
	if err != nil {
		return nil, err
	}

	return times, nil
}

func (api *API) UpdateCheck(stationId string, productId string, cartId string) (*CartInfo, error) {

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
	urlForm.Set("station_id", stationId)
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

	var header = api.newHeader()
	header.Set("ddmc-station-id", stationId)
	request.Header = header
	var cart = new(CartInfo)
	err = api.do(request, urlForm, cart)
	if err != nil {
		return nil, err
	}

	return cart, nil
}

func (api *API) CheckOrder(stationId, addressId string, productList ProductList) (*CheckOrder, error) {
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

	request, err := http.NewRequest(http.MethodPost, url.String(), nil)
	if err != nil {
		return nil, err
	}

	var header = api.newHeader()
	header.Set("ddmc-station-id", stationId)
	request.Header = header
	var checkOrder = new(CheckOrder)
	err = api.do(request, urlForm, checkOrder)
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

	var header = api.newHeader()
	header.Set("ddmc-station-id", stationId)
	request.Header = header
	var addNewOrder = new(AddNewOrder)
	err = api.do(request, params, addNewOrder)
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
	header.Set("ddmc-latitude", "23.109281")
	header.Set("ddmc-ip", "")

	header.Set("ddmc-device-id", "osP8I0RgncVIhrJLWwUCb0gi9uDQ")
	header.Set("ddmc-longitude", "113.415302")
	header.Set("ddmc-os-version", "[object Undefined]")
	header.Set("ddmc-time", strconv.FormatInt(time.Now().Unix(), 10))
	if len(api.ddmcUid) > 0 {
		header.Set("ddmc-uid", api.ddmcUid)
	}

	return header
}

func (api *API) newURLEncodedForm() url.Values {
	var params = url.Values{}
	params.Add("uid", `5db2faa481eef77f04ab13e1`)
	params.Add("longitude", `121.409128`)
	params.Add("latitude", `31.306508`)
	params.Add("station_id", `5bc5a799716de1a94f8b6fb4`)
	params.Add("city_number", `0101`)
	params.Add("api_version", `9.50.0`)
	params.Add("app_version", `2.83.0`)
	params.Add("applet_source", ``)
	params.Add("channel", `applet`)
	params.Add("app_client_id", `4`)
	params.Add("device_token", `WHJMrwNw1k/FKPjcOOgRd+Ed/O2S3GOkz07Wa1UPcfbDL2PfhzepFdBa/QF9u539PLLYm6SKU+84w6mApK0aXmA9Vne9MFdf+dCW1tldyDzmauSxIJm5Txg==1487582755342`)

	// me
	params.Add("sharer_uid", ``)
	params.Add("s_id", `4606726bbe6337d4094e1dec808431d9`)
	params.Add("openid", `osP8I0RgncVIhrJLWwUCb0gi9uDQ`)
	params.Add("h5_source", ``)
	t := strconv.FormatInt(time.Now().Unix(), 10)
	params.Add("time", t)

	return params
}

func (api *API) do(req *http.Request, form url.Values, data interface{}) error {
	if form != nil {
		var m = make(map[string]string)
		for k, v := range form {
			m[k] = v[0]
		}

		signResult, err := api.signer.Sign(m)
		if err != nil {
			return err
		}

		form.Set("nars", signResult.Nars)
		form.Set("sesi", signResult.Sesi)

		req.Body = io.NopCloser(strings.NewReader(form.Encode()))
	}

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
