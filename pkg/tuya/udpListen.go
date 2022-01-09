// Copyright 2019 py60800.
// Use of this source code is governed by Apache-2 licence
// license that can be found in the LICENSE file.

package tuya

import (
	"crypto/md5"
	"log"
	"net"
)

func udpListener(dm *DeviceManager, encrypted bool) {
	port := ":6666"
	if encrypted {
		port = ":6667"
	}
	cnx, err := net.ListenPacket("udp", port)
	if err != nil {
		log.Fatal("UDP Listener failed:", err)
	}
	for {
		buffer := make([]byte, 1024)
		n, _, err := cnx.ReadFrom(buffer)
		buffer = buffer[:n]
		if err == nil && len(buffer) > 16 {
			if uiRd(buffer) == uint(0x55aa) {
				sz := uiRd(buffer[12:])
				if sz <= uint(len(buffer)-16) {
					//discard potential leading 0
					sz = sz - 8 // discard CRC and end marker
					is := 16
					for ; buffer[is] == byte(0) && is < (int(sz)+16); is++ {
					}
					if encrypted {
						// https://github.com/codetheweb/tuyapi/blob/5a08481689c5062e17ff9a280d0e51815e2cafb7/lib/config.js
						key := md5.Sum([]byte("yGAdlopoPVldABfn"))
						buffer, err = aesDecrypt(buffer[is:16+sz], key[:], false)
						if err != nil {
							log.Print("UPD packet decryption failed: ", err)
						}
					} else {
						buffer = buffer[is : 16+sz]
					}
					dm.applianceUpdate(buffer)
				}
			}
		}
	}
}
