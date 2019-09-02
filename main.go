package main

import (
	"encoding/hex"
	"fmt"
	"os"

	ttnsdk "github.com/NeuralSpaz/go-app-sdk"
	ttnlog "github.com/TheThingsNetwork/go-utils/log"
	"github.com/TheThingsNetwork/go-utils/log/apex"
)

const (
	sdkClientName = "myApp"
	appID         = "dragino-node"
	appAccessKey  = "ttn-account-v2.lES1w4HL9arKraDzzb9Fq-1VfIh_io3ZOjxkiyu9Zlk"
)

func main() {
	log := apex.Stdout()
	log.MustParseLevel("debug")
	ttnlog.Set(log)

	if appID == "" || appAccessKey == "" {
		os.Exit(0)
	}

	config := ttnsdk.NewCommunityConfig(sdkClientName)
	config.ClientVersion = "2.0.5"

	client := config.NewClient(appID, appAccessKey)
	defer client.Close()

	devices, err := client.ManageDevices()
	if err != nil {
		log.WithError(err).Fatalf("%s: could not read CA certificate file", sdkClientName)
	}

	deviceList, err := devices.List(100, 0)
	if err != nil {
		log.WithError(err).Fatal("could not get devices")
	}

	log.Info("ethtrack: found devices")
	for _, device := range deviceList {
		fmt.Printf("- %s", device.DevID)
	}

	dev := new(ttnsdk.Device)
	dev.SparseDevice.AppID = appID
	dev.SparseDevice.DevID = "dragino-device"
	dev.SparseDevice.AppEUI = ttnsdk.AppEUI{0x70, 0xB3, 0xD5, 0x7E, 0xD0, 0x02, 0x17, 0xA7} // Use the real AppEUI here
	dev.SparseDevice.DevEUI = ttnsdk.DevEUI{0x00, 0xB0, 0x5D, 0x90, 0xD5, 0xB2, 0x88, 0xAE} // Use the real DevEUI here
	//random.FillBytes(dev.AppKey[:])                                                         // Generate a random AppKey

	err = devices.Set(dev)
	if err != nil {
		log.WithError(err).Fatalf("%s: could not create device", sdkClientName)
	}

	dev, err = devices.Get("dragino-device")
	if err != nil {
		log.WithError(err).Fatalf("%s: could not get device", sdkClientName)
	}

	log.Infof("Found dev %v", dev)

	/*PubSub ret ApplicationPubSub, error
	* interface ApplicationPubSub
	*	- Publish
	*	- Device
	*	- AllDevices
	*	- Close
	 */
	pubsub, err := client.PubSub()
	if err != nil {
		log.WithError(err).Fatal("ethtrack: could not get application pub/sub")
	}

	// allDevices := pubsub.AllDevices()
	// defer allDevices.Close()

	// activations, err := allDevices.SubscribeActivations()
	// if err != nil {
	// 	log.WithError(err).Fatalf("%s: could not subscribe to activations", sdkClientName)
	// }

	devicePubSub := pubsub.Device("dragino-device")

	uplinks, err := devicePubSub.SubscribeUplink()
	if err != nil {
		fmt.Println(err)
	}

	events, err := devicePubSub.SubscribeEvents()
	if err != nil {
		fmt.Println(err)
	}

	go func() {
		for {
			select {
			// case a := <-activations:
			// 	fmt.Printf("Activation: %+v\n", a)
			case u := <-uplinks:
				fmt.Printf("Uplink: %+v\n", u)
			case e := <-events:
				s := fmt.Sprintf("%v", e)
				log.Infof("s from e %v", s)
				str, err := hex.DecodeString(s)
				log.Infof("string from s %v", str)
				if err != nil {
					log.Fatal("Error conversion")
				}
				hexPayload := hex.EncodeToString([]byte(str))
				fmt.Printf("Event: %+v\n", hexPayload)
			}
		}
	}()

	err = pubsub.Publish("dragino-device", &ttnsdk.DownlinkMessage{
		AppID:      "dragino-node",
		DevID:      "dragino-device",
		PayloadRaw: []byte{0x01, 0x02, 0x03, 0x04},
	})
	if err != nil {
		log.WithError(err).Fatalf("%s: could not schedule downlink message", sdkClientName)
	}

	select {}
}
