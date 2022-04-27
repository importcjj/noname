package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/importcjj/ddxq/internal/config"
	"github.com/importcjj/ddxq/pkg/api"
	"github.com/importcjj/ddxq/pkg/dingding"
	"github.com/importcjj/ddxq/pkg/notify"
	"github.com/importcjj/ddxq/pkg/serverchan"
)

var (
	configFile = flag.String("config", "config.yml", "配置文件")
	boostMode  = flag.Bool("boost", false, "boost模式")
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

func Sleep(d time.Duration) {
	log.Printf("sleep %s", d)
	time.Sleep(d)
}

func intervalCheckHomePage(ddapi *api.API, mode *config.Mode, notify notify.Notify) {
	var m map[string]struct{}
	for {
		homeflow, err := ddapi.HomeFlowDetail()
		if err != nil {
			log.Println("获取首页推荐商品失败", err)
		} else {

			var firstRun bool
			if m == nil {
				firstRun = true
				m = make(map[string]struct{})
			}

			var findNew bool
			for _, product := range homeflow.List {
				_, ok := m[product.ID]
				if !ok {
					m[product.ID] = struct{}{}
					findNew = true
				}
			}

			if !firstRun && findNew {
				notify.Send(context.Background(), "首页检测到新商品")
			}
		}

		Sleep(mode.HomeInterval())
	}
}

func intervalUpdateCart(ddapi *api.API, config config.Config, mode *config.Mode) {

	for {
		// 下单中不刷新购物车，避免抢单时冲突
		if makingOrderProcess {
			continue
		}
		log.Println("正在更新购物车详情...")
		cart, err := ddapi.Cart()
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
		}

		Sleep(mode.CartInterval())
	}
}

func main() {
	flag.Parse()

	config, err := config.Load(*configFile)
	if err != nil {
		log.Fatal(err)
	}
	config.BoostMode.Enable = *boostMode || config.BoostMode.Enable
	log.Printf("%#v", config)

	mode, err := config.NewMode()
	if err != nil {
		log.Fatalf("无法创建boost: %v", err)
	}

	notify := notify.Combine(
		dingding.NewRobot(config.Dingding),
		serverchan.NewRobot(config.ServerChan),
	)

	ddapi, err := api.NewAPI(config.API)
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

	var inAddress api.Address
	for _, address := range userAddress.ValidAddress {

		if address.IsDefault {
			inAddress = address
			log.Printf("[%s] %s %s", address.StationInfo.CityName, address.Location.Address, address.AddrDetail)
			break
		}

	}

	ddapi.SetAddress(inAddress)

	//定期更新购物车
	go intervalUpdateCart(ddapi, config, mode)
	// 定期检查首页商品
	go intervalCheckHomePage(ddapi, mode, notify)

	if mode.BoostMode.Enable() {
		log.Println("注意！boost模式已启动，到时我将彻底疯狂！！！")
	}
	log.Println("开始运行...")

	var reserveTime api.ReserveTime

CheckTime:
	for {
		// boost模式非疯狂时间不请求接口
		if mode.BoostMode.Enable() &&
			!mode.BoostMode.WarmUpBoostTime() &&
			!mode.BoostMode.BoostTime() {
			continue
		}

		cart := globalCart.Get()
		if len(cart.NewOrderProductList) == 0 {
			continue
		}
		log.Println("正在获取运力...")
		times, err := ddapi.GetMultiReserveTime(cart.NewOrderProductList[0].Products)
		if err != nil {
			log.Println("获取运力失败", err)
		} else {
			time, ok := times.FirstUsableTime()
			if ok {
				reserveTime = time
				s := fmt.Sprintln("预约时间 -> ", reserveTime.SelectMsg)
				notify.Send(context.Background(), s)

				goto MakeOrder
			}

			log.Println("当前暂无可用运力...")
		}

		Sleep(mode.ReserveInterval())
	}

MakeOrder:
	log.Println("开始自动下单...")
	cart := globalCart.Get()
	if len(cart.NewOrderProductList) == 0 {
		log.Println("购物车内无可购买商品, 终止下单...")
		goto CheckTime
	}
	makingOrderProcess = true

CheckOrder:
	checkOrder, err := ddapi.CheckOrder(cart.NewOrderProductList[0], false)
	if err != nil {
		log.Println("检查订单失败", err)
		if mode.BoostMode.Enable() &&
			(mode.BoostMode.WarmUpBoostTime() || mode.BoostMode.BoostTime()) {
			Sleep(mode.RecheckInterval())
			log.Println("重新检查订单", err)
			goto CheckOrder
		}
		goto CheckTime
	}
	log.Println("检查订单成功，开始下单")

NewOrder:
	if mode.BoostMode.Enable() && mode.BoostMode.WarmUpBoostTime() {
		goto NewOrder
	}

	_, err = ddapi.AddNewOrder(api.PayTypeAlipay, cart, reserveTime, checkOrder)
	if err != nil {
		log.Println("下单失败", err)

		if mode.BoostMode.Enable() && mode.BoostMode.BoostTime() {
			Sleep(mode.ReorderInterval())
			log.Println("重新下单", err)
			goto NewOrder
		}
		goto CheckTime

	}

	makingOrderProcess = false
	notify.Send(context.Background(), "下单成功, 快抓紧付钱！")

	var continueY string
	fmt.Println("是否退出[y/n]?")
	fmt.Scanln(&continueY)

	if continueY == "n" {
		goto CheckTime
	}

	log.Println("停止运行并退出")
}
