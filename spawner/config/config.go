package config

import (
	"flag"
	"time"
)

var PathRouting bool
var CookieRouting bool
var ContainerPrefix string
var RedirectTimeout time.Duration
var CookieKey string
var CookieFallBackUrl string
var HttpProbeTimeOut time.Duration
var HttpProbeOkStatus int
var HealthCheckInterval time.Duration
var HealthCheckRetries int
var TraefikCheckEnabled bool
var TraefikBaseUrl string
var TraefikPlatform string

var EnableMgMtAuth bool
var UsersPassFile string
var TokenExpireTime time.Duration
var TokenCleanUpLoop time.Duration
func LoadFlags() {

	flag.BoolVar(&PathRouting, "use-path-routing", false, "Route to your spawned container using request path (see container prefix)")
	flag.BoolVar(&CookieRouting, "use-cookie-routing", true, "Route to your spawned container using a cookie")
	flag.StringVar(&ContainerPrefix, "container-prefix", "", "The URL prefix that route the HTTP request to the specific container")
	flag.DurationVar(&RedirectTimeout, "redirect-timeout", 2*time.Second, "Time to wait for a container to start before redirecting the request")
	flag.StringVar(&CookieKey, "cookie-key", "instance", "Value used from routing to match the container")
	flag.StringVar(&CookieFallBackUrl, "cookie-fallback-url", "http://localhost/home", "Set the url the client should be redirected to when the cookie is invalid or non-existent")
	flag.DurationVar(&HealthCheckInterval, "health-check-interval", 30*time.Second, "Number of seconds to wait between each healthcheck")
	flag.IntVar(&HealthCheckRetries, "health-check-retries", 3, "Number of healthcheck retries")
	flag.IntVar(&HttpProbeOkStatus, "http-probe-ok-status", 200, "Valid response code for an HTTP healthcheck")
	flag.DurationVar(&HttpProbeTimeOut, "http-probe-timeout", 60*time.Second, "Http request timeout")
  flag.StringVar(&TraefikBaseUrl, "traefik-base-url", "http://traefik:8080", "The traefick api endpoint")
  flag.StringVar(&TraefikPlatform, "traefik-platform", "docker", "The traefick endpoint tipe eg spawner@docker")
  flag.BoolVar(&TraefikCheckEnabled, "traefik-check-enabled", true, "Enable checking traefick route before redirect")
  flag.BoolVar(&EnableMgMtAuth, "enable-mgmt-auth", true, "Enable management APIs authentication")
  flag.StringVar(&UsersPassFile, "userspass-file", "/passwd.txt", "The file that contains user a passwords sha256 hashes")
  flag.DurationVar(&TokenExpireTime, "token-expire-time", 60 * time.Minute, "The validity time of a token")
  flag.DurationVar(&TokenCleanUpLoop, "token-cleanup-interval", 60 * time.Minute, "Time between expired token eviction")
  flag.Parse()

}
