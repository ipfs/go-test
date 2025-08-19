package random_test

import (
	"strings"
	"sync"
	"testing"

	blocks "github.com/ipfs/go-block-format"
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
	blks := random.BlocksOfSize(3, blockSize)
	require.Len(t, blks, 3)
	for _, b := range blks {
		require.Equal(t, blockSize, len(b.RawData()))
		require.Equal(t, blocks.NewBlock(b.RawData()).Cid(), b.Cid(), "block CIDs mismatch")
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
	firstNum := random.SequenceNext()

	random.SetSeed(initSeed)
	b2 := random.Bytes(32)
	secondNum := random.SequenceNext()

	require.Equal(t, b1, b2)
	require.Equal(t, firstNum, secondNum)
}

func TestSequence(t *testing.T) {
	const (
		seqCount = 5
		seqSize  = 10
	)
	firstNum := random.SequenceNext()
	secondNum := random.Sequence(1)[0]
	require.Equal(t, firstNum+1, secondNum)

	seen := make(map[uint64]struct{}, seqSize*seqCount+1)
	seen[firstNum] = struct{}{}
	seen[secondNum] = struct{}{}

	seqs := make(chan []uint64)
	startGate := make(chan struct{})
	var ready sync.WaitGroup
	ready.Add(seqCount)

	for range seqCount {
		go func() {
			ready.Done()
			<-startGate
			seq := random.Sequence(seqSize)
			seqs <- seq
		}()
	}

	ready.Wait()
	close(startGate)

	for range seqCount {
		seq := <-seqs
		for _, num := range seq {
			_, found := seen[num]
			require.Falsef(t, found, "sequence number %d is not unique", num)
			seen[num] = struct{}{}
		}
	}
	lastNum := random.SequenceNext()
	t.Log("first seq num:", firstNum)
	t.Log("last seq num: ", lastNum)
	require.Equal(t, secondNum+(seqCount*seqSize)+1, lastNum, "expected lastnum to be %d + %d + 1, was %d", secondNum, seqCount*seqSize, lastNum)
}

func TestNewRand(t *testing.T) {
	rng1 := random.NewSeededRand(137)
	rng2 := random.NewSeededRand(137)
	allEqual := true
	for range 100 {
		n1 := rng1.Int()
		n2 := rng2.Int()
		if n1 != n2 {
			allEqual = false
			break
		}
	}
	require.True(t, allEqual)

	rng1 = random.NewRand()
	rng2 = random.NewRand()
	for range 100 {
		n1 := rng1.Int()
		n2 := rng2.Int()
		if n1 != n2 {
			allEqual = false
			break
		}
	}
	require.False(t, allEqual)
}

func TestRead(t *testing.T) {
	buf1 := make([]byte, 16)
	buf2 := make([]byte, 16)

	random.NewRand().Read(buf1)
	random.NewRand().Read(buf2)
	require.NotEqual(t, buf1, buf2)
}
