package xhttp

import (
	"github.com/BetaGoRobot/BetaGo-Redefine/internal/infrastructure/config"
	"github.com/go-resty/resty/v2"
)

var (
	HttpClient          = resty.New()
	HttpClientWithProxy = resty.New().SetProxy(config.Get().ProxyConfig.PrivateProxy)
)
