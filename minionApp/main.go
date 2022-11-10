package main

import (
  "context"
  "time"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
  "log"
  "fmt"
)

func main() {
	instance_env := os.Getenv("CONNSTR")
  IdleEnv := os.Getenv("IDLE")

  Idle, err := time.ParseDuration(IdleEnv)
  
  if err != nil{
    log.Println("set IDLE to a valid time format as 60s, 1m, 2h... Using the default 1m")
    Idle, _ = time.ParseDuration("1m")
  }

  log.Printf("The server will automatically shutdown in %s if no request is received\n", Idle.String())

	hostname, _ := os.Hostname()
	router := gin.Default()

  // start a new timer
  timer1 := time.NewTimer(Idle)
  
  // if a request is received reset the timer
  instance :=	router.Group(fmt.Sprintf("session/%s",hostname))
	{
		instance.GET("/ping", func(c *gin.Context) {
    
    // as soon as we receive a request proceed to reset the timer
    timer1.Reset(Idle)


		c.JSON(http.StatusOK, gin.H{
				"message":           "pong",
				"connection string": instance_env,
			})
		})
	}

  srv := &http.Server {
    Addr: ":9000",
    Handler: router,
  }

  // listen in a separate go routine
  go func(){
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
