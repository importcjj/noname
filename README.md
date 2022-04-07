比较简单，不是自动下单

## **准备工作**

1. Charles抓包工具（安装ssl证书并信任)

2. 微信windows客户端

## **具体操作**

 1.  开启charles https监听本机请求

1.  使用微信windows客户端打开叮咚小程序，切换到购物车页面（如果提示请求异常，请检查charles的https监听是否配置正常）
2. 找到如下请求后(https://maicai.api.ddxq.mobi/order/getMultiReserveTime)，右键copy cURL request复制

![Untitled](https://s3-us-west-2.amazonaws.com/secure.notion-static.com/3e1278c0-9793-4cec-9264-87f0453cc169/Untitled.png)

1. 新建raw_request.txt放入复制的内容后保存
2. 运行程序，有运力的时候会通知，立刻去购物车下单

## 补充

可以使用钉钉群机器人webhook通知
