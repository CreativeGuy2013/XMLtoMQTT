package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/clbanning/mxj"
	MQTTClient "github.com/yosssi/gmq/mqtt/client"
)

type MQTTConfig struct {
	Host    string `json:"host"`
	Port    string `json:"port"`
	BaseKey string `json:"base_key"`
	Name    string `json:"client_name"`
}

type Map map[string]interface{}

var client *MQTTClient.Client
var baseKey string

func main() {
	fmt.Println("Starting client.")

	var configuration MQTTConfig

	//Read config file and unmarshal
	loadedConfig, _ := ioutil.ReadFile("config.json")
	json.Unmarshal(loadedConfig, &configuration)

	baseKey = configuration.BaseKey
	initMQTT(configuration)

	parseXML()
}

func initMQTT(config MQTTConfig) {
	fmt.Println("Starting MQTT")

	//Initialize a new MQTT client
	client = MQTTClient.New(&MQTTClient.Options{
		// Define the error handler.
		ErrorHandler: func(err error) {
			fmt.Println(err)
		},
	})

	//If the programm stops close the MQTT client
	defer client.Terminate()

	//Connect to the broker with the config provided
	err := client.Connect(&MQTTClient.ConnectOptions{
		Network:  "tcp",
		Address:  config.Host + ":" + config.Port,
		ClientID: []byte(config.Name),
	})
	if err != nil {
		panic(err)
	}
}

func parseXML() {
	xmlFile, _ := os.Open("config.xml")

	mv, err := mxj.NewMapXmlReader(xmlFile) // repeated calls, as with an os.File Reader, will process stream
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	jsonVal, err := mv.Json()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(string(jsonVal))

}
