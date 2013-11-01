
package alipay

import (
	"sort"
	"fmt"
	"io"
	"strings"
	"net/url"
	"crypto/md5"
)

var alipayGatewayNew = `https://mapi.alipay.com/gateway.do?`

type Config struct {
	Partner string
	Key string
}

type Request struct {
	Service string
	PaymentType string
	NotifyUrl string
	ReturnUrl string
	OutTradeNo string
	Subject string
	TotalFee float64
	Body string
	ShowUrl string
	SellerEmail string
}

type Response struct {
	BuyerEmail string
	OutTradeNo string
	TradeStatus string
	Subject string
	TotalFee float64
}

type kvpair struct {
	k, v string
}

type kvpairs []kvpair

func (t kvpairs) Less(i, j int) bool {
	return t[i].k < t[j].k
}

func (t kvpairs) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func (t kvpairs) Len() int {
	return len(t)
}

func (t kvpairs) Sort() {
	sort.Sort(t)
}

func (t kvpairs) RemoveEmpty() (t2 kvpairs) {
	for _, kv := range t {
		if kv.v != "" {
			t2 = append(t2, kv)
		}
	}
	return
}

func (t kvpairs) Join() string {
	var strs []string
	for _, kv := range t {
		strs = append(strs, kv.k+"="+kv.v)
	}
	return strings.Join(strs, "&")
}

func md5Sign(str, key string) string {
	h := md5.New()
	io.WriteString(h, str)
	io.WriteString(h, key)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func verifySign(c Config, u url.Values) (err error) {
	p := kvpairs{}
	sign := ""
	for k := range u {
		v := u.Get(k)
		switch k {
		case "sign":
			sign = v
			continue
		case "sign_type":
			continue
		}
		p = append(p, kvpair{k, v})
	}
	if sign == "" {
		err = fmt.Errorf("sign not found")
		return
	}
	p = p.RemoveEmpty()
	p.Sort()
	fmt.Println(u)
	if md5Sign(p.Join(), c.Key) != sign {
		err = fmt.Errorf("sign invalid")
		return
	}
	return
}

func ParseResponse(c Config, p url.Values) (r Response, err error) {
	if err = verifySign(c, p); err != nil {
		return
	}

	r.BuyerEmail = p.Get("buyer_email")
	r.TradeStatus = p.Get("trade_status")
	r.OutTradeNo = p.Get("out_trade_no")
	r.Subject = p.Get("subject")
	fmt.Sscanf(p.Get("total_fee"), "%f", &r.TotalFee)

	if r.TradeStatus != "TRADE_SUCCESS" && r.TradeStatus != "TRADE_FINISHED" {
		err = fmt.Errorf("trade not success or finnished")
		return
	}
	return
}

func NewPage(c Config, r Request, w io.Writer) {
	p := kvpairs {
		kvpair{`_input_charset`, `utf-8`},
		kvpair{`out_trade_no`, r.OutTradeNo},
		kvpair{`partner`, c.Partner},
		kvpair{`payment_type`, r.PaymentType},
		kvpair{`notify_url`, r.NotifyUrl},
		kvpair{`return_url`, r.ReturnUrl},
		kvpair{`subject`, r.Subject},
		kvpair{`total_fee`, fmt.Sprintf("%.2f", r.TotalFee)},
		kvpair{`body`, r.Body},
		kvpair{`service`, r.Service},
		kvpair{`show_url`, r.ShowUrl},
		kvpair{`seller_email`, r.SellerEmail},
	}
	p = p.RemoveEmpty()
	p.Sort()

	sign := md5Sign(p.Join(), c.Key)
	p = append(p, kvpair{`sign`, sign})
	p = append(p, kvpair{`sign_type`, `MD5`})

	fmt.Fprintln(w, `<html><head>
	<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
	</head><body>`)
	fmt.Fprintf(w, `<form name='alipaysubmit' action='%s_input_charset=utf-8' method='post'> `, alipayGatewayNew)
	for _, kv := range p {
		fmt.Fprintf(w, `<input type='hidden' name='%s' value='%s' />`, kv.k, kv.v)
	}
	fmt.Fprintln(w, `<script>document.forms['alipaysubmit'].submit();</script>`)
	fmt.Fprintln(w, `</body></html>`)
}

