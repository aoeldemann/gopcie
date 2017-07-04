# gopcie PCIExpress library

Golang package gopcie implements data transfer to/from a PCIExpress hardware
device on a Linux-based system. Data can be transferred via:

  1) PCIExpress Direct Memory Access (DMA) transfers (requires kernel-space
     device driver) or
  2) PCIExpress Base Address Register (BAR) accesses.

The BAR resource file identification is based on Andre Richter's
[easy-pci-mmap](https://github.com/andre-richter/easy-pci-mmap).

