package vivo

import (
	"net/url"
	"time"

	"resty.dev/v3"
)

func (vivo *Vivo) newHttpClient() *resty.Request {
	client := resty.New()
	client.AddRequestMiddleware(func(c *resty.Client, r *resty.Request) error {
		u, e := url.Parse(r.URL)
		if e != nil {
			return e
		}
		signing_string := r.Method + "\n" + u.Path + "\n" + r.QueryParams.Encode() + "\n" + vivo.appID + "\n" + r.Header.Get("X-AI-GATEWAY-TIMESTAMP") + "\n" + "x-ai-gateway-app-id" + ":" + vivo.appID + "\n" + "x-ai-gateway-timestamp" + ":" + r.Header.Get("X-AI-GATEWAY-TIMESTAMP") + "\n" + "x-ai-gateway-nonce" + ":" + r.Header.Get("X-AI-GATEWAY-NONCE")
		r.Header.Set("X-AI-GATEWAY-SIGNATURE", base64encode(hMACSHA256HEX(signing_string, vivo.appKey)))
		return nil
	})
	req := client.R()
	req.Header.Add("X-AI-GATEWAY-APP-ID", vivo.appID)
	req.Header.Add("X-AI-GATEWAY-TIMESTAMP", int64toString(time.Now().Unix()))
	req.Header.Add("X-AI-GATEWAY-NONCE", generateRandomString(8))
	req.Header.Add("X-AI-GATEWAY-SIGNED-HEADERS", "x-ai-gateway-app-id;x-ai-gateway-timestamp;x-ai-gateway-nonce")
	return req
}
