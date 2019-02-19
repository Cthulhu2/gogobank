// Golang gogobank project
//
//     Schemes: http, https
//     BasePath: /v1
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
// swagger:meta
package main

import (
	"encoding/json"
	"log"
	"runtime"

	"app/jsonconfig"
	"app/model"
	"app/route"
	"app/server"
)

// config the settings variable
var config = &configuration{}

type configuration struct {
	Server server.Server `json:"Server"`
}

func init() {
	// Verbose logging with file name and line number
	log.SetFlags(log.Lshortfile)

	// Use all CPU cores
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	// Load the configuration file
	jsonconfig.Load("config/config.json", config)

	model.Init()

	// Start the listener
	server.Run(route.LoadHTTP(), route.LoadHTTPS(), config.Server)
}

// ParseJSON unmarshals bytes to structs
func (c *configuration) ParseJSON(b []byte) error {
	return json.Unmarshal(b, &c)
}
