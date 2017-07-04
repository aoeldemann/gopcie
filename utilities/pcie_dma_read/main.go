//
// Project:        gopcie
// File:           main.go
// Date Create:    June 9th 2017
// Date Modified:  July 4th 2017
// Author:         Andreas Oeldemann, TUM <andreas.oeldemann@tum.de>
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
