package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"net"
	"net/http"
	"os"
	"time"
  "mock-app/build"
)

func main() {
  build.LoadInfo()

  instance_env := os.Getenv("CONNSTR")
	IdleEnv := os.Getenv("IDLE")


	startTime := time.Now()

	Idle, err := time.ParseDuration(IdleEnv)

	if err != nil {
		log.Println("set IDLE to a valid time format as 60s, 1m, 2h... Using the default 1m")
		Idle, _ = time.ParseDuration("1m")
	}

	log.Printf("The server will automatically shutdown in %s if no request is received\n", Idle.String())

	hostname, _ := os.Hostname()

	ifaces, err := net.InterfaceAddrs()
	addr := []string{}
	for _, i := range ifaces {
		addr = append(addr, i.String())
	}

	router := gin.Default()

	// start a new timer

	timer1 := time.NewTimer(Idle)

	// if a request is received reset the timer
	instance := router.Group("/")
	{
		instance.GET("/", func(c *gin.Context) {

			// as soon as we receive a request proceed to reset the timer
			timer1.Reset(Idle)
			startTime = time.Now()

			c.JSON(http.StatusOK, gin.H{
				"hostname":          hostname,
				"IP address":        addr,
				"connection string": instance_env,
        "requested": c.Request.Host+c.Request.URL.Path,
			})
		})

		instance.GET("/status", func(c *gin.Context) {
			// insert a delay to make sure healthcheck is also delayed
			time.Sleep(5 * time.Second)
			delta := time.Now().Sub(startTime)
			c.JSON(http.StatusOK, gin.H{
				"status": "ok",
				"left":   (Idle - delta).String(),
        "requested": c.Request.Host+c.Request.URL.Path,
			})
		})
	}
	srv := &http.Server{
		Addr:    ":9000",
		Handler: router,
	}

	// listen in a separate go routine
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Handle the shutdown with a context allowing connection draining"
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Wait for the timer to tick
	<-timer1.C
	log.Println("Timer 1 expired the application is going down")

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown some requests are pending: ", err)
	} else {
		log.Print("Server shutdown success")
	}
}
