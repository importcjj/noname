package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/importcjj/ddxq/pkg/api"
	"github.com/importcjj/ddxq/pkg/dingding"
)

var (
	cookie       = flag.String("cookie", "", "叮咚cookie")
	dingdinghook = flag.String("dingding", "", "钉钉机器人")
)

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
		if address.IsDefault {
			stationId = address.StationID
			addressId = address.ID
			log.Printf("[%s] %s", address.StationInfo.CityName, address.Location.Address)
			break
		}

	}
	var cart *api.CartInfo
	var reserveTime api.ReserveTime

GetCart:
	log.Println("正在获取购物车详情中...")
	for {
		cart, err = ddapi.Cart(stationId)
		if err != nil {
			log.Println("购物车获取失败", err)
			time.Sleep(10 * time.Second)
			continue
		}

		// 勾选有货的商品
		if effective := cart.Product.Effective; len(effective) > 0 {
			list := effective[0]
			for _, item := range list.Products {
				cart, err = ddapi.UpdateCheck(stationId, item.ID, item.CartID)
				if err != nil {
					log.Println(err)
				}

				fmt.Println(cart)
			}
		}

		if len(cart.NewOrderProductList) == 0 {
			log.Println("购物车无可购买商品")
			time.Sleep(10 * time.Second)
			continue
		}

		break
	}

	fmt.Println(cart.NewOrderProductList)

	log.Println("正在获取可用运力中...")
GetTime:
	for {
		times, err := ddapi.GetMultiReverseTime(stationId, cart.NewOrderProductList[0].Products)
		if err != nil {
			log.Println("获取运力失败", err)
		} else {
			for _, item := range *times {
				for _, day := range item.Time {
					for _, time := range day.Times {
						if !time.FullFlag {
							reserveTime = time
							break GetTime
						}
					}
				}
			}
		}
		time.Sleep(1500 * time.Millisecond)
	}

	dingdingbot.Send(context.Background(), reserveTime.SelectMsg)
	log.Println("正在自动下单中...")

	checkOrder, err := ddapi.CheckOrder(stationId, addressId, cart.NewOrderProductList[0])
	if err != nil {
		log.Println("检查订单失败", err)
		goto GetCart
	}

	order, err := ddapi.AddNewOrder(stationId, addressId, api.PayTypeAlipay, cart, reserveTime, checkOrder)
	if err != nil {
		log.Println("下单失败", err)
		goto GetCart
	}

	log.Println("下单成功", order)
	dingdingbot.Send(context.Background(), "下单成功, 请付款")

	goto GetCart
}
