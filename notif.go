// Copyright for portions of this fork are held by [Protocol Labs, Inc., 2016] as
// part of the original go-libp2p-kad-dht project. All other copyright for
// this fork are held by [The BDWare Authors, 2020]. All rights reserved.
// Use of this source code is governed by MIT license that can be
// found in the LICENSE file.

package dht

import (
	"context"

	"github.com/libp2p/go-libp2p-core/helpers"
	"github.com/libp2p/go-libp2p-core/network"

	ma "github.com/multiformats/go-multiaddr"
	mstream "github.com/multiformats/go-multistream"
)

// netNotifiee defines methods to be used with the IpfsDHT
type netNotifiee IpfsDHT

func (nn *netNotifiee) DHT() *IpfsDHT {
	return (*IpfsDHT)(nn)
}

func (nn *netNotifiee) Connected(n network.Network, v network.Conn) {
	dht := nn.DHT()
	select {
	case <-dht.Process().Closing():
		return
	default:
	}

	p := v.RemotePeer()
	protos, err := dht.peerstore.SupportsProtocols(p, dht.protocolStrs()...)
	if err == nil && len(protos) != 0 {
		// We lock here for consistency with the lock in testConnection.
		// This probably isn't necessary because (dis)connect
		// notifications are serialized but it's nice to be consistent.
		dht.plk.Lock()
		defer dht.plk.Unlock()
		if dht.host.Network().Connectedness(p) == network.Connected {
			refresh := dht.routingTable.Size() <= minRTRefreshThreshold
			dht.Update(dht.Context(), p)
			if refresh && dht.autoRefresh {
				select {
				case dht.triggerRtRefresh <- nil:
				default:
				}
			}
		}
		return
	}

	// Note: Unfortunately, the peerstore may not yet know that this peer is
	// a DHT server. So, if it didn't return a positive response above, test
	// manually.
	go nn.testConnection(v)
}

func (nn *netNotifiee) testConnection(v network.Conn) {
	dht := nn.DHT()
	p := v.RemotePeer()

	// Forcibly use *this* connection. Otherwise, if we have two connections, we could:
	// 1. Test it twice.
	// 2. Have it closed from under us leaving the second (open) connection untested.
	s, err := v.NewStream()
	if err != nil {
		// Connection error
		return
	}
	defer helpers.FullClose(s)

	selected, err := mstream.SelectOneOf(dht.protocolStrs(), s)
	if err != nil {
		// Doesn't support the protocol
		return
	}
	// Remember this choice (makes subsequent negotiations faster)
	dht.peerstore.AddProtocols(p, selected)

	// We lock here as we race with disconnect. If we didn't lock, we could
	// finish processing a connect after handling the associated disconnect
	// event and add the peer to the routing table after removing it.
	dht.plk.Lock()
	defer dht.plk.Unlock()
	if dht.host.Network().Connectedness(p) == network.Connected {
		refresh := dht.routingTable.Size() <= minRTRefreshThreshold
		dht.Update(dht.Context(), p)
		if refresh && dht.autoRefresh {
			select {
			case dht.triggerRtRefresh <- nil:
			default:
			}
		}
	}
}

func (nn *netNotifiee) Disconnected(n network.Network, v network.Conn) {
	dht := nn.DHT()
	select {
	case <-dht.Process().Closing():
		return
	default:
	}

	p := v.RemotePeer()

	// Lock and check to see if we're still connected. We lock to make sure
	// we don't concurrently process a connect event.
	dht.plk.Lock()
	defer dht.plk.Unlock()
	if dht.host.Network().Connectedness(p) == network.Connected {
		// We're still connected.
		return
	}

	dht.routingTable.Remove(p)
	dht.host.ConnManager().Unprotect(p, "routingTable")
	if dht.routingTable.Size() < minRTRefreshThreshold {
		// TODO: Actively bootstrap. For now, just try to add the currently connected peers.
		for _, p := range dht.host.Network().Peers() {
			// Don't bother probing, we do that on connect.
			protos, err := dht.peerstore.SupportsProtocols(p, dht.protocolStrs()...)
			if err == nil && len(protos) != 0 {
				dht.Update(dht.Context(), p)
			}
		}
	}

	dht.smlk.Lock()
	defer dht.smlk.Unlock()
	ms, ok := dht.strmap[p]
	if !ok {
		return
	}
	delete(dht.strmap, p)

	// Do this asynchronously as ms.lk can block for a while.
	go func() {
		ms.lk.Lock(context.Background())
		defer ms.lk.Unlock()
		ms.invalidate()
	}()
}

func (nn *netNotifiee) OpenedStream(n network.Network, v network.Stream) {}
func (nn *netNotifiee) ClosedStream(n network.Network, v network.Stream) {}
func (nn *netNotifiee) Listen(n network.Network, a ma.Multiaddr)         {}
func (nn *netNotifiee) ListenClose(n network.Network, a ma.Multiaddr)    {}
