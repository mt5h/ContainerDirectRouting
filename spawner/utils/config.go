package utils

import(
  "flag"
)

var PathRouting bool 
var CookieRouting bool
var ContainerPrefix string 
var RedirectTimeout int
var RedirectHealthcheckTimeout int
var CookieKey string
var CookieMaxAge int
var CookieSecure bool
var CookieHttpOnly bool
func LoadFlags() {

  flag.BoolVar(&PathRouting, "use-path-routing", false, "Route to your spawned container using request path (see container prefix)")
  flag.BoolVar(&CookieRouting, "use-cookie-routing", true, "Route to your spawned container using a cookie (default method)")
  flag.StringVar(&ContainerPrefix, "container-prefix", "", "The URL prefix that route the HTTP request to the specific container")
  flag.IntVar(&RedirectTimeout, "redirect-timeout", 2, "Time to wait for a container to start before redirecting the request")
  flag.IntVar(&RedirectHealthcheckTimeout, "redirect-healthcheck-timeout", 2, "Time to wait for a container to start before redirecting the request")
  flag.StringVar(&CookieKey, "cookie-key", "instance", "Value used from routing to match the container") 
  flag.IntVar(&CookieMaxAge, "cookie-max-age", 24*60*60, "Maxage set for the cookie used to route requests")
  flag.BoolVar(&CookieSecure, "cookie-secure", false, "Set cookie only over a secure connection")
  flag.BoolVar(&CookieHttpOnly, "cookie-http-only", true, "Exclude the cookie if the cookie-string is being generated for a non-HTTP API")

  flag.Parse()

}
