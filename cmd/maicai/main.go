package main

import (
	"flag"
	"log"
	"time"

	"github.com/importcjj/ddxq/pkg/api"
)

var cookie = flag.String("cookie", "", "叮咚cookie")

func main() {
	flag.Parse()

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
			time.Sleep(1 * time.Minute)
			continue
		}

		if len(cart.NewOrderProductList) == 0 {
			log.Println("购物车无可购买商品")
			time.Sleep(1 * time.Minute)
			continue
		}

		break
	}

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

	log.Println("正在自动下单中...")

	checkOrder, err := ddapi.CheckOrder(stationId, addressId, cart.NewOrderProductList[0])
	if err != nil {
		log.Println("检查订单失败", err)
		goto GetCart
	}

	order, err := ddapi.AddNewOrder(stationId, addressId, 6, cart, reserveTime, checkOrder)
	if err != nil {
		log.Println("下单失败", err)
		goto GetCart
	}

	log.Println("下单成功", order)
	goto GetCart
}
