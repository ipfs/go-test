package random_test

import (
	"strings"
	"testing"

	"github.com/ipfs/go-test/random"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multihash"
	"github.com/stretchr/testify/require"
)

func TestBytes(t *testing.T) {
	b1 := random.Bytes(32)
	b2 := random.Bytes(32)
	require.Len(t, b1, 32)
	require.NotEqual(t, b1, b2)
	t.Logf("bytes: %x", b1)
	t.Logf("bytes: %x", b2)
}

func TestBlocksOfSize(t *testing.T) {
	const blockSize = 64
	blocks := random.BlocksOfSize(3, blockSize)
	require.Len(t, blocks, 3)
	for _, b := range blocks {
		require.Equal(t, blockSize, len(b.RawData()))
	}
}

func TestHttpMultiaddrs(t *testing.T) {
	hms := random.HttpMultiaddrs(3)
	require.Len(t, hms, 3)
	for _, ma := range hms {
		maStr := ma.String()
		require.True(t, strings.HasSuffix(maStr, "http"))
		t.Log("http multiaddr:", maStr)
	}
}

func TestCids(t *testing.T) {
	cids := random.Cids(3)
	require.Len(t, cids, 3)
	for _, c := range cids {
		require.True(t, strings.HasPrefix(c.String(), "baguqeera"))
		t.Log("cid:", c.String())
	}
}

func TestIdentity(t *testing.T) {
	id, pvtKey, pubKey := random.Identity()

	testID, err := peer.IDFromPrivateKey(pvtKey)
	require.NoError(t, err)
	require.Equal(t, id, testID)
	testPub := pvtKey.GetPublic()
	require.Equal(t, pubKey, testPub)
}

func TestMultihashes(t *testing.T) {
	mhs := random.Multihashes(3)
	require.Len(t, mhs, 3)
	for _, mh := range mhs {
		decMh, err := multihash.Decode(mh)
		require.NoError(t, err)
		require.Equal(t, int(decMh.Code), multihash.SHA2_256)
	}
}

func TestPeers(t *testing.T) {
	peerIDs := random.Peers(3)
	require.Len(t, peerIDs, 3)
	for _, peerID := range peerIDs {
		peerStr := peerID.String()
		require.True(t, strings.HasPrefix(peerStr, "12D3Koo"))
		t.Log("peerID:", peerStr)
	}
}

// TestSeed tests that setting seed back to its initial value generates the
// same results.
func TestSeed(t *testing.T) {
	initSeed := random.Seed()

	random.SetSeed(initSeed)
	b1 := random.Bytes(32)

	random.SetSeed(initSeed)
	b2 := random.Bytes(32)

	require.Equal(t, b1, b2)
}
