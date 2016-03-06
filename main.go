package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
)

const fileTTL = 60 * time.Second

func main() {
	addr := flag.String("listen", "localhost:3336", "address for server to listen")
	flag.Parse()

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(mw.Logger())
	e.Use(mw.Recover())

	// Routes
	e.Get("/", ping)
	e.Post("/", gather)
	e.Get("/download/:name", download)

	// Start server
	e.Run(*addr)
}

func ping(c *echo.Context) error {
	return c.String(http.StatusOK, "pong\n")
}

type downloadRequest struct {
	Name string   `json:"name"`
	URLs []string `json:"urls"`
}

type downloadResponse struct {
	Name string `json:"name"`
}

func gather(c *echo.Context) error {
	var req downloadRequest
	if err := c.Bind(&req); err != nil {
		return err
	}
	t := time.Now()
	name := req.Name + "-" + t.Format("20060102150405") + ".zip"
	Download(name, req.URLs)
	// setup delayed cleanup
	go cleanup(name)
	resp := &downloadResponse{name}
	return c.JSON(http.StatusCreated, resp)
}

func download(c *echo.Context) error {
	name := c.Param("name")
	return c.File(name, name, true)
}

func cleanup(name string) {
	time.AfterFunc(fileTTL, func() {
		log.Printf("Deleting file %s", name)
		os.Remove(name)
	})
}