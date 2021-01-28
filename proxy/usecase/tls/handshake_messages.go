// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tls

import (
	"strings"

	"golang.org/x/crypto/cryptobyte"
)

// readUint8LengthPrefixed acts like s.ReadUint8LengthPrefixed, but targets a
// []byte instead of a cryptobyte.String.
func readUint8LengthPrefixed(s *cryptobyte.String, out *[]byte) bool {
	return s.ReadUint8LengthPrefixed((*cryptobyte.String)(out))
}

type ClientHelloMsg struct {
	raw        []byte
	vers       uint16
	random     []byte
	sessionId  []byte
	ServerName string
}

func (m *ClientHelloMsg) Unmarshal(data []byte) bool {
	*m = ClientHelloMsg{raw: data}
	s := cryptobyte.String(data)

	if !s.Skip(4) || // message type and uint24 length field
		!s.ReadUint16(&m.vers) || !s.ReadBytes(&m.random, 32) ||
		!readUint8LengthPrefixed(&s, &m.sessionId) {
		return false
	}

	var cipherSuites cryptobyte.String
	if !s.ReadUint16LengthPrefixed(&cipherSuites) {
		return false
	}
	for !cipherSuites.Empty() {
		var suite uint16
		if !cipherSuites.ReadUint16(&suite) {
			return false
		}
	}

	if s.Empty() {
		// ClientHello is optionally followed by extension data
		return true
	}

	s.Skip(2) // Compression Methods
	s.Skip(2) // Extensions Length

	var extensions = s
	for !extensions.Empty() {
		var extension uint16
		var extData cryptobyte.String
		if !extensions.ReadUint16(&extension) ||
			!extensions.ReadUint16LengthPrefixed(&extData) {
			return false
		}

		switch extension {
		case extensionServerName:
			// RFC 6066, Section 3
			var nameList cryptobyte.String
			if !extData.ReadUint16LengthPrefixed(&nameList) || nameList.Empty() {
				return false
			}
			for !nameList.Empty() {
				var nameType uint8
				var serverName cryptobyte.String
				if !nameList.ReadUint8(&nameType) ||
					!nameList.ReadUint16LengthPrefixed(&serverName) ||
					serverName.Empty() {
					return false
				}
				if nameType != 0 {
					continue
				}
				if len(m.ServerName) != 0 {
					// Multiple names of the same name_type are prohibited.
					return false
				}
				m.ServerName = string(serverName)
				// An SNI value may not include a trailing dot.
				if strings.HasSuffix(m.ServerName, ".") {
					return false
				}
				return true
			}
		default:
			// Ignore unknown extensions.
			continue
		}

		if !extData.Empty() {
			return false
		}
	}

	return true
}
