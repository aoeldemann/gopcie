//
// Project:        gopcie
// File:           main.go
// Date Create:    June 9th 2017
// Date Modified:  July 4th 2017
// Author:         Andreas Oeldemann, TUM <andreas.oeldemann@tum.de>
//
// Description:
//
// Utility to write a PCIExpress Base Address Register (BAR) from the
// command-line.
//

package main

import (
	"flag"
	"fmt"
	"github.com/aoeldemann/gopcie"
)

func main() {
	// read command line arguments
	var functionIdStr, vendorIdStr, deviceIdStr, barIdStr, addrStr,
		dataStr string
	flag.StringVar(&addrStr, "addr", "", "addr")
	flag.StringVar(&dataStr, "data", "", "data")
	flag.StringVar(&functionIdStr, "functionId", "", "device function ID")
	flag.StringVar(&vendorIdStr, "vendorId", "", "device vendor ID")
	flag.StringVar(&deviceIdStr, "deviceId", "", "device ID")
	flag.StringVar(&barIdStr, "barId", "", "device BAR ID")
	flag.Parse()

	// make sure parameters are set
	if len(addrStr) == 0 || len(dataStr) == 0 || len(functionIdStr) == 0 ||
		len(vendorIdStr) == 0 || len(deviceIdStr) == 0 || len(barIdStr) == 0 {
		flag.Usage()
		return
	}

	// convert hex string values to int
	addr, err := gopcie.HexStringToInt(addrStr)
	if err != nil {
		panic("invalid address")
	}
	data, err := gopcie.HexStringToInt(dataStr)
	if err != nil {
		panic("invalid data")
	}
	functionId, err := gopcie.HexStringToInt(functionIdStr)
	if err != nil {
		panic("invalid device function ID")
	}
	vendorId, err := gopcie.HexStringToInt(vendorIdStr)
	if err != nil {
		panic("invalid device vendor ID")
	}
	deviceId, err := gopcie.HexStringToInt(deviceIdStr)
	if err != nil {
		panic("invalid device ID")
	}
	barId, err := gopcie.HexStringToInt(barIdStr)
	if err != nil {
		panic("invalid BAR ID")
	}

	// create and open pcie bar
	pcieBAR, err := gopcie.PCIeBAROpen(uint(functionId), uint(vendorId),
		uint(deviceId), uint(barId))
	if err != nil {
		panic(err.Error())
	}
	defer pcieBAR.Close()

	// write data
	pcieBAR.Write(uint32(addr), uint32(data))

	// print write address and data
	fmt.Printf("Addr: 0x%08x\n", addr)
	fmt.Printf("Data: 0x%08x\n", data)
}
