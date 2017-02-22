package renter

import (
	"bytes"
	"crypto/rand"
	"testing"
)

// TestRSEncode tests the rsCode type.
func TestRSEncode(t *testing.T) {
	badParams := []struct {
		data, parity int
	}{
		{-1, -1},
		{-1, 0},
		{0, -1},
		{0, 0},
		{0, 1},
		{1, 0},
	}
	for _, ps := range badParams {
		if _, err := NewRSCode(ps.data, ps.parity); err == nil {
			t.Error("expected bad parameter error, got nil")
		}
	}

	rsc, err := NewRSCode(10, 3)
	if err != nil {
		t.Fatal(err)
	}
	var needed []uint64
	for i := 0; i < rsc.NumPieces(); i++ {
		needed = append(needed, uint64(i))
	}

	data := make([]byte, 777)
	rand.Read(data)

	pieces, err := rsc.Encode(data, needed)
	if err != nil {
		t.Fatal(err)
	}
	_, err = rsc.Encode(nil, needed)
	if err == nil {
		t.Fatal("expected nil data error, got nil")
	}

	buf := new(bytes.Buffer)
	err = rsc.Recover(pieces, 777, buf)
	if err != nil {
		t.Fatal(err)
	}
	err = rsc.Recover(nil, 777, buf)
	if err == nil {
		t.Fatal("expected nil pieces error, got nil")
	}

	if !bytes.Equal(data, buf.Bytes()) {
		t.Fatal("recovered data does not match original")
	}
}

func BenchmarkRSEncode(b *testing.B) {
	rsc, err := NewRSCode(10, 20)
	if err != nil {
		b.Fatal(err)
	}
	var needed []uint64
	for i := 0; i < rsc.NumPieces(); i++ {
		needed = append(needed, uint64(i))
	}
	data := make([]byte, 10*1<<22)
	rand.Read(data)

	b.SetBytes(10 * 1 << 22)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := rsc.Encode(data, needed)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkRSRecover(b *testing.B) {
	rsc, err := NewRSCode(10, 20)
	if err != nil {
		b.Fatal(err)
	}
	var needed []uint64
	for i := 20; i < rsc.NumPieces(); i++ {
		needed = append(needed, uint64(i))
	}
	data := make([]byte, 10*1<<22)
	rand.Read(data)
	pieces, err := rsc.Encode(data, needed)
	if err != nil {
		b.Fatal(err)
	}

	b.SetBytes(10 * 1 << 22)
	buf := new(bytes.Buffer)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := rsc.Recover(pieces, 10*1<<22, buf)
		if err != nil {
			b.Fatal(err)
		}
	}
}
