//
// Project:        gopcie
// File:           helpers.go
// Date Create:    June 9th 2017
// Date Modified:  July 4th 2017
// Author:         Andreas Oeldemann, TUM <andreas.oeldemann@tum.de>
//
// Description:
//
// Helper functions.
//

package gopcie

import (
	"strconv"
)

// HexStringToInt converts a hex string (which may start with a '0x' prefix) to
// an integer.
func HexStringToInt(hexStr string) (value uint64, err error) {
	// if hex addr string starts with 0x remove the prefix
	if (len(hexStr) > 1) && (hexStr[0:2] == "0x") {
		hexStr = hexStr[2:]
	}

	// convert hex string to integer
	value, err = strconv.ParseUint(hexStr, 16, 64)
	return value, err
}
