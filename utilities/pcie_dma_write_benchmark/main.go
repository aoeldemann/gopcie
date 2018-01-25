//
// Project:        gopcie
// File:           main.go
// Date Create:    January 25th 2018
// Date Modified:  January 25th 2018
// Author:         Andreas Oeldemann, TUM <andreas.oeldemann@tum.de>
//
// Description:
//
// tbd
//

package main

import (
	"flag"
	"fmt"
	"github.com/aoeldemann/gopcie"
	"time"
)

func main() {
	// read command line arguments
	var addrStr, sizeStr, device string
	flag.StringVar(&addrStr, "addr", "", "addr")
	flag.StringVar(&sizeStr, "size", "", "size")
	flag.StringVar(&device, "device", "", "device")
	flag.Parse()

	// make sure parameters are set
	if len(addrStr) == 0 || len(sizeStr) == 0 || len(device) == 0 {
		flag.Usage()
		return
	}

	// convert hex addr string to int
	addr, err := gopcie.HexStringToInt(addrStr)
	if err != nil {
		panic("invalid address")
	}

	// convert hex size string to int
	size, err := gopcie.HexStringToInt(sizeStr)
	if err != nil {
		panic("invalid address")
	}

	// create and open pcie device
	dev, err := gopcie.PCIeDMAOpen(device, gopcie.PCIE_ACCESS_WRITE)
	if err != nil {
		panic(err.Error())
	}
	defer dev.Close()

	// create a byte slice containing write data (all zeros)
	data := make([]byte, size)

	for {
		// record time before transfer
		transferStartTime := time.Now()

		// write to pcie dev
		dev.Write(addr, data)

		// get duration since transfer start
		transferDuration := time.Since(transferStartTime)

		// calculate average throughput in Gbps
		transferThroughput :=
			8.0 * float64(size) / transferDuration.Seconds() / 1e9

		// print out infos
		fmt.Printf("%d bytes; %s; %f Gbps\n", size, transferDuration,
			transferThroughput)
	}
}
