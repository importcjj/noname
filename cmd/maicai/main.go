package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/importcjj/ddxq/pkg/api"
	"github.com/importcjj/ddxq/pkg/dingding"
)

var (
	cookie       = flag.String("cookie", "DDXQSESSID=4c165b29f0466c6acbc48ee07b0b992c", "叮咚cookie， 抓包小程序可得")
	dingdinghook = flag.String("dingding", "", "钉钉机器人")
	sid          = flag.String("sid", "4c165b29f0466c6acbc48ee07b0b992c", "抓包小程序可得")
	openid       = flag.String("openid", "osP8I0f05BPiuikzy0HQeSMubrg4", "抓包小程序可得")
	deviceID     = flag.String("device_id", "osP8I0f05BPiuikzy0HQeSMubrg4", "抓包小程序可得")
	deviceToken  = flag.String("device_token", "WHJMrwNw1k/F0qdLNvE01AUwlQtiEc7qol6Nyikv9NlcvvPrivpurzjH084wLhpKazlcN+OgJ5CPOGBpFtSvdQi8BwFXwUa90dCW1tldyDzmauSxIJm5Txg==1487582755342", "抓包小程序可得")
	ua           = flag.String("ua", "Mozilla/5.0 (iPhone; CPU iPhone OS 11_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E217 MicroMessenger/6.8.0(0x16080000) NetType/WIFI Language/en Branch/Br_trunk MiniProgramEnv/Mac", "User-Agent, 抓包小程序可得")
	boostMode    = flag.Bool("boost", true, "彻底疯狂！！！！！")
)

var globalCart = NewCart()
var makingOrderProcess = false

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

func intervalUpdateCart(ddapi *api.API) {

	for {
		// 下单中不刷新购物车，避免抢单时冲突
		if makingOrderProcess {
			continue
		}
		log.Println("正在更新购物车详情...")
		cart, err := ddapi.Cart()
		if err != nil {
			log.Println("购物车获取失败", err)
			if *boostMode && preBoostTime() {
				time.Sleep(10 * time.Second)
			} else {
				time.Sleep(2 * time.Minute)
			}
		} else {

			// 勾选有货的商品
			if effective := cart.Product.Effective; len(effective) > 0 {
				list := effective[0]
				for _, item := range list.Products {
					// 不重复勾选
					if item.IsCheck == 1 {
						continue
					}
					cart, err = ddapi.UpdateCheck(item.ID, item.CartID)
					if err != nil {
						log.Println(err)
					} else {
						log.Printf("已添加 %s", item.ProductName)
					}

				}
			} else {
				log.Println("购物车暂无可下单商品")
			}
			globalCart.Set(cart)
			time.Sleep(2 * time.Minute)
		}

		//time.Sleep(550 * time.Millisecond)
	}
}

func main() {
	flag.Parse()

	dingdingbot := dingding.NewRobot("信号", *dingdinghook)

	ddapi, err := api.NewAPI(*cookie)
	if err != nil {
		log.Fatal(err)
	}

	ddapi.
		SetSID(*sid).
		SetOpenID(*openid).
		SetDeviceID(*deviceID).
		SetDeviceToken(*deviceToken).
		SetUserAgent(*ua)

	userDetail, err := ddapi.UserDetail()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("用户ID: %s", userDetail.UserInfo.ID)

	userAddress, err := ddapi.UserAddress()
	if err != nil {
		log.Fatal(err)
	}

	var inAddress api.Address
	for _, address := range userAddress.ValidAddress {

		if address.IsDefault {
			inAddress = address
			log.Printf("[%s] %s %s", address.StationInfo.CityName, address.Location.Address, address.AddrDetail)
			break
		}

	}

	ddapi.SetAddress(inAddress)

	// 定期更新购物车
	go intervalUpdateCart(ddapi)

	log.Println("开始运行...")
	if *boostMode {
		log.Println("注意！boost模式已启动，到时我将彻底疯狂！！！")
	}

	var reserveTime api.ReserveTime

CheckTime:

	for {
		// boost模式非疯狂时间不请求接口
		if *boostMode && !boostTime() {
			continue
		}

		cart := globalCart.Get()
		if len(cart.NewOrderProductList) == 0 {
			continue
		}
		log.Println("正在获取运力...")
		times, err := ddapi.GetMultiReverseTime(cart.NewOrderProductList[0].Products)
		if err != nil {
			log.Println("获取运力失败", err)
		} else {
			for _, item := range *times {
				for _, day := range item.Time {
					for _, time := range day.Times {
						if !time.FullFlag {
							reserveTime = time
							log.Println("预约时间 -> ", time)
							dingdingbot.Send(context.Background(), reserveTime.SelectMsg)

							goto MakeOrder
						}
					}
				}
			}

			log.Println("当前暂无可用运力...")
		}

		if boostTime() {
			time.Sleep(550 * time.Millisecond)
		} else {
			time.Sleep(2000 * time.Millisecond)
		}

	}

MakeOrder:
	log.Println("开始自动下单...")
	cart := globalCart.Get()
	if len(cart.NewOrderProductList) == 0 {
		log.Println("购物车内无可购买商品, 终止下单...")
		goto CheckTime
	}
	makingOrderProcess = true
	checkOrder, err := ddapi.CheckOrder(cart.NewOrderProductList[0])
	if err != nil {
		log.Println("检查订单失败", err)
		if *boostMode && boostTime() {
			checkOrderSuccess := false
			for !checkOrderSuccess {
				log.Println("重新检查订单", err)
				checkOrder, err = ddapi.CheckOrder(cart.NewOrderProductList[0])
				if err != nil {
					log.Println("检查订单失败", err)
					time.Sleep(500 * time.Millisecond)
				} else {
					checkOrderSuccess = true
				}
			}
		} else {
			goto CheckTime
		}
	}
	log.Println("检查订单成功，开始下单")
	order, err := ddapi.AddNewOrder(api.PayTypeAlipay, cart, reserveTime, checkOrder)
	if err != nil {
		newOrderSuccess := false
		log.Println("下单失败", err)
		if *boostMode && boostTime() {
			for !newOrderSuccess {
				log.Println("重新下单", err)
				order, err = ddapi.AddNewOrder(api.PayTypeAlipay, cart, reserveTime, checkOrder)
				if err != nil {
					log.Println("下单失败", err)
					time.Sleep(500 * time.Millisecond)
				} else {
					newOrderSuccess = true
				}
			}
		} else {
			goto CheckTime
		}
	}

	log.Println("下单成功", order)
	makingOrderProcess = false
	dingdingbot.Send(context.Background(), "下单成功, 请付款")

	var continueY string
	fmt.Println("是否退出[y/n]?")
	fmt.Scanln(&continueY)

	if continueY == "n" {
		goto CheckTime
	}

	log.Println("停止运行并退出")
}

func boostTime() bool {

	now := time.Now()
	if now.Hour() == 6 && now.Minute() <= 5 {
		return true
	}

	if now.Hour() == 8 && now.Minute() >= 30 && now.Minute() <= 35 {
		return true
	}

	return false
}

func preBoostTime() bool {
	now := time.Now()
	if now.Hour() == 5 && now.Minute() >= 58 {
		return true
	}

	if now.Hour() == 8 && now.Minute() >= 28 && now.Minute() <= 30 {
		return true
	}

	return false
}
