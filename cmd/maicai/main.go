package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/importcjj/ddxq/pkg/api"
)

var cookie = flag.String("cookie", "DDXQSESSID=b6b755f09c045e9732dcc31ac9b67203", "叮咚cookie")

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
		}
		log.Printf("[%s]", address.StationInfo.CityName)
	}

	cart, err := ddapi.Cart(stationId)
	if err != nil {
		log.Fatal(err)
	}

	times, err := ddapi.GetMultiReverseTime(stationId, cart.NewOrderProductList[0].Products)
	if err != nil {
		log.Fatal(err)
	}

	var reserveTime api.ReserveTime
	for _, t := range *times {
		reserveTime = t.Time[0].Times[0]
		break
	}

	checkOrder, err := ddapi.CheckOrder(stationId, addressId, cart.NewOrderProductList[0])
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("=====")
	order, err := ddapi.AddNewOrder(stationId, addressId, 6, cart, reserveTime, checkOrder)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(order)
}
