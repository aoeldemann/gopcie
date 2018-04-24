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
