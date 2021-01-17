// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tls

const (
	MaxHandshake = 65536 // maximum handshake we support (protocol max is 16 MB)
)

// TLS extension numbers
const (
	extensionServerName uint16 = 0
)
