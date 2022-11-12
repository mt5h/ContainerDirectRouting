package utils

import(
  "flag"
)

var ContainerPrefix string 
var RedirectTimeout int
var RedirectHealthcheckTimeout int

func LoadFlags() {

  flag.StringVar(&ContainerPrefix, "container-prefix", "session", "The URL prefix that route the HTTP request to the specific container")
  flag.IntVar(&RedirectTimeout, "redirect-timeout", 2, "Time to wait for a container to start before redirecting the request")
  flag.IntVar(&RedirectHealthcheckTimeout, "redirect-healthcheck-timeout", 2, "Time to wait for a container to start before redirecting the request")
  flag.Parse()

}
