package utils

import (
	"slices"
	"sync"
	"time"

	"github.com/gofiber/fiber/v3/client"
)

var vpnCache sync.Map
var _ = startCacheCleaner()

var loopbacks = []string{
	"127.0.0.1",
	"::1",
}

func startCacheCleaner() bool {
	go func() {
		ticker := time.NewTicker(3 * time.Hour)
		for range ticker.C {
			vpnCache.Clear()
		}
	}()
	return true
}

func VPNCheck(ip string) bool {
	if slices.Contains(loopbacks, ip) {
		return false
	}

	if value, ok := vpnCache.Load(ip); ok {
		return value.(bool)
	}

	cl := client.New()
	resp, err := cl.Get("http://ip-api.com/json/" + ip + "?fields=proxy,hosting,mobile")

	if err != nil {
		return false
	}

	if resp.StatusCode() != 200 {
		return false
	}

	var result struct {
		Proxy   bool `json:"proxy"`
		Hosting bool `json:"hosting"`
		Mobile  bool `json:"mobile"`
	}

	if err := resp.JSON(&result); err != nil {
		return false
	}

	if result.Proxy || result.Hosting || result.Mobile {
		vpnCache.Store(ip, true)
		return true
	}

	vpnCache.Store(ip, false)

	return false
}
