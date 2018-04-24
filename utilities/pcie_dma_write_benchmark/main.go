//
// The MIT License
//
// Copyright (c) 2017-2018 by the author(s)
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.
//
// Author(s):
//   - Andreas Oeldemann <andreas.oeldemann@tum.de>
//
// Date Created:        January 25th 2018
// Date Last Modified:  January 25th 2018
//
// Description:
//
// Tool benchmarks PCIe DMA write throughput.
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
