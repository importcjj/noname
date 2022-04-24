package api

import "encoding/json"

type Response struct {
	Success bool            `json:"success"`
	Code    interface{}     `json:"code"`
	Msg     string          `json:"msg"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

// 用户信息相关

type UserDetail struct {
	UserInfo UserInfo `json:"user_info"`
}

type UserInfo struct {
	ID string `json:"id"`
}

// 用户收货地址相关

type UserAddress struct {
	ValidAddress    []Address `json:"valid_address"`
	InvalidAddress  []Address `json:"invalid_address"`
	MaxAddressCount int       `json:"max_address_count"`
	CanAddAddress   bool      `json:"can_add_address"`
}

type Address struct {
	ID          string      `json:"id"`
	Gender      int         `json:"gender"`
	Mobile      string      `json:"mobile"`
	Location    Location    `json:"location"`
	Label       string      `json:"label"`
	UserName    string      `json:"user_name"`
	AddrDetail  string      `json:"addr_detail"`
	StationID   string      `json:"station_id"`
	StationName string      `json:"station_name"`
	IsDefault   bool        `json:"is_default"`
	StationInfo StationInfo `json:"station_info"`
}

type Location struct {
	Typecode string    `json:"typecode"`
	Address  string    `json:"address"`
	Name     string    `json:"name"`
	Location []float64 `json:"location"`
	ID       string    `json:"id"`
}

type StationInfo struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Phone      string `json:"phone"`
	CityName   string `json:"city_name"`
	CityNumber string `json:"city_number"`
}

// 购物车相关

type CartInfo struct {
	Product             CartProduct     `json:"product"`
	NewOrderProductList []ProductList   `json:"new_order_product_list"`
	ParentOrderInfo     ParentOrderInfo `json:"parent_order_info"`
}

type ParentOrderInfo struct {
	ParentOrderSign string `json:"parent_order_sign"`
}

type CartProduct struct {
	Effective []CartProductList `json:"effective"`
}

type CartProductList struct {
	Products []CartProductItem `json:"products"`
}

type CartProductItem struct {
	ProductName string     `json:"product_name"`
	ID          string     `json:"id"`
	CartID      string     `json:"cart_id"`
	Sizes       []struct{} `json:"sizes"`
	IsCheck     int        `json:"is_check"`
}

type ProductList struct {
	Products               []ProductListItem `json:"products"`
	TotalMoney             string            `json:"total_money"`
	TotalOriginMoney       string            `json:"total_origin_money"`
	GoodsRealMoney         string            `json:"goods_real_money"`
	TotalCount             int               `json:"total_count"`
	CartCount              int               `json:"cart_count"`
	IsPresale              int               `json:"is_presale"`
	InstantRebateMoney     string            `json:"instant_rebate_money"`
	CouponRebateMoney      string            `json:"coupon_rebate_money"`
	TotalRebateMoney       string            `json:"total_rebate_money"`
	UsedBalanceMoney       string            `json:"used_balance_money"`
	CanUsedBalanceMoney    string            `json:"can_used_balance_money"`
	UsedPointNum           int               `json:"used_point_num"`
	UsedPointMoney         string            `json:"used_point_money"`
	CanUsedPointNum        int               `json:"can_used_point_num"`
	CanUsedPointMoney      string            `json:"can_used_point_money"`
	IsShareStation         int               `json:"is_share_station"`
	OnlyTodayProducts      []struct{}        `json:"only_today_products"`
	OnlyTomorrowProducts   []struct{}        `json:"only_tomorrow_products"`
	PackageType            int               `json:"package_type"`
	PackageID              int               `json:"package_id"`
	FrontPackageText       string            `json:"front_package_text"`
	FrontPackageType       int               `json:"front_package_type"`
	FrontPackageStockColor string            `json:"front_package_stock_color"`
	FrontPackageBgColor    string            `json:"front_package_bg_color"`
}

type ProductListItem struct {
	Type               int        `json:"type"`
	ProductType        int        `json:"product_type"`
	ID                 string     `json:"id"`
	CategoryPath       string     `json:"category_path"`
	TotalMoney         string     `json:"total_money"`
	TotalOriginMoney   string     `json:"total_origin_money"`
	TotalPrice         string     `json:"total_price"`
	TotalOriginPrice   string     `json:"total_origin_price"`
	InstantRebateMoney string     `json:"instant_rebate_money"`
	ActivityID         string     `json:"activity_id"`
	ConditionsNum      string     `json:"conditions_num"`
	Price              string     `json:"price"`
	OriginPrice        string     `json:"origin_price"`
	PriceType          int        `json:"price_type"`
	BatchType          int        `json:"batch_type"`
	SubList            []struct{} `json:"sub_list"`
	Count              int        `json:"count"`
	Description        string     `json:"description"`
	ParentID           string     `json:"parent_id"`
	Sizes              []struct{} `json:"sizes"`
	CartID             string     `json:"cart_id"`
	IsBooking          int        `json:"is_booking"`
	ProductName        string     `json:"product_name"`
	SmallImage         string     `json:"small_image"`
	SaleBatches        struct {
		BatchType int `json:"batch_type"`
	} `json:"sale_batches"`
	OrderSort int `json:"order_sort"`
}

// 下单相关

type MultiReserveTime []ReserveTimeItem

type ReserveTimeItem struct {
	Time []struct {
		DateStr          string        `json:"date_str"`
		DateStrTimestamp int64         `json:"date_str_timestamp"`
		Day              string        `json:"day"`
		Times            []ReserveTime `json:"times"`
	} `json:"time"`
}

type ReserveTime struct {
	Type           int    `json:"type"`
	FullFlag       bool   `json:"fullFlag"`
	StartTime      string `json:"start_time"`
	StartTimestamp int64  `json:"start_timestamp"`
	EndTimestamp   int64  `json:"end_timestamp"`
	EndTime        string `json:"end_time"`
	SelectMsg      string `json:"select_msg"`
}

type CheckOrder struct {
	Order struct {
		Freights             []FreightItem `json:"freights"`
		TotalMoney           string        `json:"total_money"`
		FreightDiscountMoney string        `json:"freight_discount_money"`
		FreightMoney         string        `json:"freight_money"`
		FreightRealMoney     string        `json:"freight_real_money"`
		DefaultCoupon        struct {
			Name  string `json:"name"`
			Money string `json:"money"`
			ID    string `json:"_id"`
		} `json:"default_coupon"`
	} `json:"order"`
}

type FreightItem struct {
	Freight struct {
		Type                 int    `json:"typ"`
		Remark               string `json:"remark"`
		FreightMoney         string `json:"freight_money"`
		DiscountFreightMoney string `json:"discount_freight_money"`
		FreightRealMoney     string `json:"freight_real_money"`
	} `json:"freight"`
	PackageId int `json:"package_id"`
}

type AddNewOrder struct {
}
