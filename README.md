alipay
======

支付宝集成接口的 golang 实现

[![Build Status](https://drone.io/github.com/go-av/alipay/status.png)](https://drone.io/github.com/go-av/alipay/latest)

示例：即时到账收款接口
======

  	r := Request{
  		NotifyUrl: `http://192.168.1.1/dirpay/notify_url.php`, // 付款后异步通知页面
  		ReturnUrl: `http://192.168.1.1/dirpay/return_url.php`, // 付款后返回页面
  		OutTradeNo: `138329153`, // 订单号
  		SellerEmail: `xxx@qq.com`, // 支付宝用户名
  		Service: `create_direct_pay_by_user`, // 不可改
  		PaymentType: `1`, // 不可改
  		Subject: `测试`, // 商品名称
  		TotalFee: 1.03, // 价格
  	}
  	
  	c := Config{
	  	Partner: `xxxxxx`,   // 支付宝合作者身份 ID
	  	Key: `xxxxxxxxxx`, // 支付宝交易安全校验码
	}
	  
	// 输出的是 html 页面，会自动跳转到支付界面
	NewPage(c, r, os.Stdout)
	  
