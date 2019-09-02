package main

import (
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

	pubsub, err := client.PubSub()
	if err != nil {
		log.WithError(err).Fatal("ethtrack: could not get application pub/sub")
	}

	all := pubsub.AllDevices()
	defer all.Close()
	// activations, err := all.SubscribeActivations()
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// uplinks, err := all.SubscribeUplink()
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// events, err := all.SubscribeEvents()
	// if err != nil {
	// 	fmt.Println(err)
	// }

	/*
		for {
			select {
			case a := <-activations:
				fmt.Printf("Activation: %+v\n", a)
			case u := <-uplinks:
				fmt.Printf("Uplink: %+v\n", u)
			case e := <-events:
				fmt.Printf("Event: %+v\n", e)
			}
		}
	*/

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

	log.Info("Set device")
	log.Infof("dev %v", dev)

	// pubsub, err := client.PubSub()
	// if err != nil {
	// 	log.WithError(err).Fatalf("%s: could not get application pub/sub", sdkClientName)
	// }

	// allDevicesPubSub := pubsub.AllDevices()

	// activations, err := allDevicesPubSub.SubscribeActivations()
	// if err != nil {
	// 	log.WithError(err).Fatalf("%s: could not subscribe to activations", sdkClientName)
	// }
	// go func() {
	// 	for activation := range activations {
	// 		log.WithFields(ttnlog.Fields{
	// 			"appEUI":  activation.AppEUI.String(),
	// 			"devEUI":  activation.DevEUI.String(),
	// 			"devAddr": activation.DevAddr.String(),
	// 		}).Info("my-amazing-app: received activation")
	// 	}
	// }()

	// err = allDevicesPubSub.UnsubscribeActivations()
	// if err != nil {
	// 	log.WithError(err).Fatalf("%s: could not unsubscribe from activations", sdkClientName)
	// }

	// myNewDevicePubSub := pubsub.Device("my-new-device")

	// uplink, err := myNewDevicePubSub.SubscribeUplink()
	// if err != nil {
	// 	log.WithError(err).Fatalf("%s: could not subscribe to uplink messages", sdkClientName)
	// }
	// go func() {
	// 	for message := range uplink {
	// 		hexPayload := hex.EncodeToString(message.PayloadRaw)
	// 		log.WithField("data", hexPayload).Infof("%s: received uplink", sdkClientName)
	// 	}
	// }()

	// err = myNewDevicePubSub.UnsubscribeUplink()
	// if err != nil {
	// 	log.WithError(err).Fatalf("%s: could not unsubscribe from uplink", sdkClientName)
	// }

	// err = myNewDevicePubSub.Publish(&types.DownlinkMessage{
	// 	PayloadRaw: []byte{0xaa, 0xbc},
	// 	FPort:      10,
	// })
	// if err != nil {
	// 	log.WithError(err).Fatalf("%s: could not schedule downlink message", sdkClientName)
	// }
}
