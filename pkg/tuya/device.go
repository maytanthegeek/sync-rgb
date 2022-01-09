// Copyright 2019 py60800.
// Use of this source code is governed by Apache-2 licence
// license that can be found in the LICENSE file.

package tuya

import (
	"sync"
	"sync/atomic"
)

type SyncMsg struct {
	Code int
	Dev  Device
}
type SyncChannel chan SyncMsg

type Device interface {
	Type() string
	Name() string
	Subscribe(SyncChannel) int64
	Unsubscribe(int64)
	Configure(*Appliance, *ConfigurationData)
	ProcessResponse(int, []byte)
	Init(string, *Appliance, *ConfigurationData)
}

// Code for tuya messages
const (
	CodeMsgSet        = 7
	CodeMsgStatus     = 10
	CodeMsgPing       = 9
	CodeMsgAutoStatus = 8
)

// to be embedded in Device
type BaseDevice struct {
	sync.Mutex
	typ  string
	name string
	App  *Appliance
	// Publish subscribe
	subscribers map[int64]SyncChannel
}

// BaseDevice initialization to be invoked during configation
func (b *BaseDevice) Init(typ string, a *Appliance, c *ConfigurationData) {
	b.typ = typ
	b.App = a
	b.name = c.Name
	b.subscribers = make(map[int64]SyncChannel)
}

// Implementation of Device interface provided by BaseDevice
func (b *BaseDevice) Type() string {
	return b.typ
}
func (b *BaseDevice) Name() string {
	return b.name
}

// Publish subscribe
var keyIndexCpt int64

func MakeSyncChannel() SyncChannel {
	return SyncChannel(make(chan SyncMsg, 1))
}

func (b *BaseDevice) Subscribe(c SyncChannel) int64 {
	b.Lock()
	defer b.Unlock()
	key := atomic.AddInt64(&keyIndexCpt, 1) // ignore overflow
	b.subscribers[key] = c
	return key
}
func (b *BaseDevice) Unsubscribe(key int64) {
	b.Lock()
	defer b.Unlock()
	delete(b.subscribers, key)
}
func (b *BaseDevice) Notify(code int, d Device) {
	b.Lock()
	defer b.Unlock()
	syncMsg := SyncMsg{code, d}
	for _, c := range b.subscribers {
		select {
		case c <- syncMsg:
		default:
		}
	}
}
