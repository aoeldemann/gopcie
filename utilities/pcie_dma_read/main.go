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
// Date Created:        June 9th 2017
// Date Last Modified:  July 4th 2017
//
// Description:
//
// Command-line utility to read data from a PCIExpress device via a DMA
// transfer.
//

package main

import (
	"flag"
	"github.com/aoeldemann/gopcie"
	"os"
)

func main() {
	// read command line arguments
	var addrStr, sizeStr string
	var filename, device string
	flag.StringVar(&addrStr, "addr", "", "address")
	flag.StringVar(&sizeStr, "size", "", "size")
	flag.StringVar(&filename, "file", "", "target filename")
	flag.StringVar(&device, "device", "", "device")
	flag.Parse()

	// make sure parameters are set
	if len(addrStr) == 0 || len(sizeStr) == 0 || len(filename) == 0 ||
		len(device) == 0 {
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
		panic("invalid size")
	}

	// create and open pcie device
	dev, err := gopcie.PCIeDMAOpen(device, gopcie.PCIE_ACCESS_READ)
	if err != nil {
		panic(err.Error())
	}
	defer dev.Close()

	// create output file
	file, err := os.Create(filename)
	if err != nil {
		panic("could not create output file")
	}
	defer file.Close()

	// read data from pcie device
	data := make([]byte, size)
	dev.Read(addr, data)

	// write data to output file
	nBytesWritten, err := file.Write(data)
	if err != nil || uint64(nBytesWritten) != size {
		panic("could not write output file")
	}
}
