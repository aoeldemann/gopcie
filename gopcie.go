//
// Project:        gopcie
// File:           gopcie.go
// Date Create:    June 9th 2017
// Date Modified:  July 4th 2017
// Author:         Andreas Oeldemann, TUM <andreas.oeldemann@tum.de>
//

/*
Package gopcie implements data transfer to/from a PCIExpress hardware device on
a Linux-based system. Data can be transferred via:

  1) PCIExpress Direct Memory Access (DMA) transfers (requires kernel-space
	 device driver) or
  2) PCIExpress Base Address Register (BAR) accesses.

The BAR resource file identification is based on Andre Richter's easy-pci-mmap:
https://github.com/andre-richter/easy-pci-mmap

*/
package gopcie

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"unsafe"
)

const (
	PCIE_ACCESS_READ  = 1
	PCIE_ACCESS_WRITE = 2
)

// PCIeDMA implements PCIExpress DMA reads and writes.
type PCIeDMA struct {
	fd         *os.File
	accessMode int
}

// PCIeDMAOpen opens a PCIExpress DMA device. The function expects the
// device name (i.e. /dev/...) and an access mode flag (read/write).
func PCIeDMAOpen(devName string, accessMode int) (*PCIeDMA, error) {
	// check if access mode is set
	if (accessMode & (PCIE_ACCESS_READ | PCIE_ACCESS_WRITE)) == 0 {
		return nil, errors.New("invalid access mode")
	}

	// open device
	fd, err := os.OpenFile(devName, os.O_RDWR, 0600)
	if err != nil {
		return nil,
			errors.New("could not open device")
	}

	// create device
	dev := PCIeDMA{
		fd:         fd,
		accessMode: accessMode,
	}

	return &dev, nil
}

// Close closes the PCIExpress DMA device.
func (dev *PCIeDMA) Close() {
	dev.fd.Close()
}

// Write performs a DMA write transfer.
func (dev *PCIeDMA) Write(addr uint64, data []byte) error {
	// check if access mode allows writing
	if (dev.accessMode & PCIE_ACCESS_WRITE) == 0 {
		return errors.New("access mode does not allow writing")
	}

	// perform write transfer
	nBytesWritten, err := dev.fd.WriteAt(data, int64(addr))
	if err != nil || nBytesWritten != len(data) {
		return errors.New("could not write to device")
	}
	return nil
}

// Read performs a DMA read transfer.
func (dev *PCIeDMA) Read(addr uint64, data []byte) error {
	// check if access mode allows reading
	if (dev.accessMode & PCIE_ACCESS_READ) == 0 {
		return errors.New("access mode does not allow reading")
	}

	// perform read transfer
	nBytesRead, err := dev.fd.ReadAt(data, int64(addr))
	if err != nil || nBytesRead != len(data) {
		return errors.New("could not read from device")
	}
	return nil
}

// PCIeBAR implements reads and writes from/to a PCIExpress base address
// registers.
type PCIeBAR struct {
	fd  *os.File
	bar []byte
}

// PCIeBAROpen opens the PCIExpress base address register. The function expects
// the function, vendor, device and bar ID of the bar to be opened.
func PCIeBAROpen(functionId, vendorId, deviceId, barId uint) (*PCIeBAR, error) {
	barFilename := ""

	// list system devices directory
	devDirs, err := ioutil.ReadDir("/sys/bus/pci/devices")
	if err != nil {
		return nil, errors.New("could not read /sys/bus/pci/devices directory")
	}

	// iterate over all devices
	for _, devDir := range devDirs {

		// not the device we are looking for if directory name does not start
		// with "0000:"
		if devDir.Name()[0:5] != "0000:" {
			continue
		}

		// get function id and see if it matches the one we are looking for
		functionIdFound, err := strconv.ParseUint(
			devDir.Name()[len(devDir.Name())-1:], 10, 32)
		if err != nil || uint(functionIdFound) != functionId {
			continue
		}

		// read device vendor file
		vendorFilename := filepath.Join("/sys/bus/pci/devices", devDir.Name(),
			"vendor")
		vendorFile, err := ioutil.ReadFile(vendorFilename)
		if err != nil {
			return nil, errors.New("could not open pci vendor file")
		}
		vendorFileStr := string(vendorFile)

		// device vendor file should have only one line and start with "0x"
		if vendorFileStr[0:2] != "0x" ||
			strings.Index(vendorFileStr, "\n") != len(vendorFileStr)-1 {
			continue
		}
		vendorFileStr = vendorFileStr[0 : len(vendorFileStr)-1]

		// get vendor id
		vendorIdFound, err := strconv.ParseInt(vendorFileStr[2:], 16, 32)
		if err != nil || uint(vendorIdFound) != vendorId {
			continue
		}

		// read pci device file
		deviceFilename := filepath.Join("/sys/bus/pci/devices", devDir.Name(),
			"device")
		deviceFile, err := ioutil.ReadFile(deviceFilename)
		if err != nil {
			return nil, errors.New("could not open pci device file")
		}
		deviceFileStr := string(deviceFile)

		// device file should have only one line and start with "0x"
		if deviceFileStr[0:2] != "0x" ||
			strings.Index(deviceFileStr, "\n") != len(deviceFileStr)-1 {
			continue
		}
		deviceFileStr = deviceFileStr[0 : len(deviceFileStr)-1]

		// get device id
		deviceIdFound, err := strconv.ParseInt(deviceFileStr[2:], 16, 32)
		if err != nil || uint(deviceIdFound) != deviceId {
			continue
		}

		// all ids matched. found it!
		barFilename = filepath.Join("/sys/bus/pci/devices", devDir.Name(),
			fmt.Sprintf("resource%d", barId))
		break
	}

	// check if the BAR resource file was found
	if len(barFilename) == 0 {
		return nil, errors.New("could not find BAR")
	}

	// stat the BAR resource file to get its size
	barFileInfo, err := os.Stat(barFilename)
	if err != nil {
		return nil, errors.New("could not stat BAR resource file")
	}

	// open BAR resource file
	fd, err := os.OpenFile(barFilename, os.O_RDWR|os.O_SYNC, 0666)
	if err != nil {
		return nil, errors.New("could not open BAR resource file")
	}

	// memory-map the BAR
	bar, err := syscall.Mmap(int(fd.Fd()), 0, int(barFileInfo.Size()),
		syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	if err != nil {
		return nil, errors.New("could not memory-map BAR")
	}

	return &PCIeBAR{fd, bar}, nil
}

// Close closes the PCIExpress base address register.
func (bar *PCIeBAR) Close() error {
	// un-memory map the BAR
	err := syscall.Munmap(bar.bar)
	if err != nil {
		return errors.New("could not un-memory-map BAR")
	}
	// close BAR resource file
	bar.fd.Close()
	return nil
}

// Write writes data to a PCIExpress base address register.
func (bar *PCIeBAR) Write(addr, data uint32) {
	*(*uint32)(unsafe.Pointer(uintptr(unsafe.Pointer(&bar.bar[0])) +
		uintptr(addr))) = data
}

// WriteMask writes data to a PCIExpress base address register. The specified
// mask determins which bits of the register shall be written.
func (bar *PCIeBAR) WriteMask(addr, data, mask uint32) {
	rd_data := bar.Read(addr)
	wr_data := (rd_data & ^mask) | (data & mask)
	bar.Write(addr, wr_data)
}

// Read reads data from a PCIExpress base address register.
func (bar *PCIeBAR) Read(addr uint32) uint32 {
	return *(*uint32)(unsafe.Pointer(uintptr(unsafe.Pointer(&bar.bar[0])) +
		uintptr(addr)))
}
