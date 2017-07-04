//
// Project:        gopcie
// File:           main.go
// Date Create:    June 9th 2017
// Date Modified:  July 4th 2017
// Author:         Andreas Oeldemann, TUM <andreas.oeldemann@tum.de>
//
// Description:
//
// Command-line utility to write data to a PCIExpress device via a DMA
// transfer.

package main

import (
	"flag"
	"github.com/aoeldemann/gopcie"
	"io/ioutil"
	"os"
)

func main() {
	// read command line arguments
	var addrStr string
	var filename, device string
	flag.StringVar(&addrStr, "addr", "", "addr")
	flag.StringVar(&filename, "file", "", "source filename")
	flag.StringVar(&device, "device", "", "device")
	flag.Parse()

	// make sure parameters are set
	if len(addrStr) == 0 || len(filename) == 0 || len(device) == 0 {
		flag.Usage()
		return
	}

	// convert hex addr string to int
	addr, err := gopcie.HexStringToInt(addrStr)
	if err != nil {
		panic("invalid address")
	}

	// create and open pcie device
	dev, err := gopcie.PCIeDMAOpen(device, gopcie.PCIE_ACCESS_WRITE)
	if err != nil {
		panic(err.Error())
	}
	defer dev.Close()

	// open input file
	file, err := os.Open(filename)
	if err != nil {
		panic("could not open input file")
	}

	// get input file size
	fileInfo, err := file.Stat()
	if err != nil {
		panic("could not stat input file")
	}
	fileSize := fileInfo.Size()

	// close file
	file.Close()

	// read input file
	data, err := ioutil.ReadFile(filename)
	if err != nil || len(data) != int(fileSize) {
		panic("could not read input file")
	}

	// write to pcie dev
	dev.Write(addr, data)
}
