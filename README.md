# gopcie PCIExpress library

Golang package gopcie implements data transfer to/from a PCIExpress hardware
device on a Linux-based system. Data can be transferred via:

  1) PCIExpress Direct Memory Access (DMA) transfers (requires kernel-space
     device driver) or
  2) PCIExpress Base Address Register (BAR) accesses.

The BAR resource file identification is based on Andre Richter's
[easy-pci-mmap](https://github.com/andre-richter/easy-pci-mmap).

## Utilities

* `pcie_bar_read`: Command-line utility to read data from PCIExpress Base
Address Register
* `pcie_bar_write`: Command-line utility to write data to PCIExpress Base
Address Register
* `pcie_dma_read`: Command-line utility to read data from PCIExpress device via
Direct Memory Access transfer
* `pcie_dma_write`: Command-line utility to write data to PCIExpress device via
Direct Memory Access transfer
