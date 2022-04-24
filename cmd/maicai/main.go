package main

import (
	"context"
	"flag"
	"log"
	"sync"
	"time"

	"github.com/importcjj/ddxq/pkg/api"
	"github.com/importcjj/ddxq/pkg/dingding"
)

var (
	cookie       = flag.String("cookie", "", "叮咚cookie")
	dingdinghook = flag.String("dingding", "", "钉钉机器人")
)

var globalCart = NewCart()

type Cart struct {
	cart *api.CartInfo
	mu   sync.Mutex
}

func NewCart() *Cart {
	return &Cart{
		cart: new(api.CartInfo),
	}
}

func (cart *Cart) Set(newCart *api.CartInfo) {
	cart.mu.Lock()
	defer cart.mu.Unlock()

	cart.cart = newCart
}

func (cart *Cart) Get() *api.CartInfo {
	cart.mu.Lock()
	defer cart.mu.Unlock()

	return cart.cart
}

func intervalUpdateCart(stationId string, ddapi *api.API) {

	for {
		log.Println("正在更新购物车详情...")
		cart, err := ddapi.Cart(stationId)
		if err != nil {
			log.Println("购物车获取失败", err)
		} else {

			// 勾选有货的商品
			if effective := cart.Product.Effective; len(effective) > 0 {
				list := effective[0]
				for _, item := range list.Products {
					// 不重复勾选
					if item.IsCheck == 1 {
						continue
					}
					cart, err = ddapi.UpdateCheck(stationId, item.ID, item.CartID)
					if err != nil {
						log.Println(err)
					} else {
						log.Printf("已添加 %s", item.ProductName)
					}

				}
			}

			globalCart.Set(cart)
		}

		time.Sleep(1 * time.Minute)
	}
}

func main() {
	flag.Parse()

	dingdingbot := dingding.NewRobot("信号", *dingdinghook)

	ddapi, err := api.NewAPI(*cookie)
	if err != nil {
		log.Fatal(err)
	}

	userDetail, err := ddapi.UserDetail()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("用户ID: %s", userDetail.UserInfo.ID)

	userAddress, err := ddapi.UserAddress()
	if err != nil {
		log.Fatal(err)
	}

	var stationId string
	var addressId string
	for _, address := range userAddress.ValidAddress {

		// if address.IsDefault {
		stationId = address.StationID
		addressId = address.ID
		log.Printf("[%s] %s", address.StationInfo.CityName, address.Location.Address)
		break
		// }

	}

	// 定期更新购物车
	go intervalUpdateCart(stationId, ddapi)

	log.Println("开始运行...")
	var reserveTime api.ReserveTime

CheckTime:
	log.Println("正在检查运力...")
	for {

		cart := globalCart.Get()
		if len(cart.NewOrderProductList) == 0 {
			continue
		}

		times, err := ddapi.GetMultiReverseTime(stationId, addressId, cart.NewOrderProductList[0].Products)
		if err != nil {
			log.Println("获取运力失败", err)
		} else {
			for _, item := range *times {
				for _, day := range item.Time {
					for _, time := range day.Times {
						if !time.FullFlag {
							reserveTime = time
							dingdingbot.Send(context.Background(), reserveTime.SelectMsg)

							goto MakeOrder
						}
					}
				}
			}
		}
		time.Sleep(1500 * time.Millisecond)
	}

MakeOrder:
	log.Println("开始自动下单...")
	cart := globalCart.Get()
	if len(cart.NewOrderProductList) == 0 {
		log.Println("购物车内无可购买商品, 终止下单")
		goto CheckTime
	}

	checkOrder, err := ddapi.CheckOrder(stationId, addressId, cart.NewOrderProductList[0])
	if err != nil {
		log.Println("检查订单失败", err)
		goto CheckTime
	}

	order, err := ddapi.AddNewOrder(stationId, addressId, api.PayTypeAlipay, cart, reserveTime, checkOrder)
	if err != nil {
		log.Println("下单失败", err)
		goto CheckTime
	}

	log.Println("下单成功", order)
	dingdingbot.Send(context.Background(), "下单成功, 请付款")

	goto CheckTime
}
