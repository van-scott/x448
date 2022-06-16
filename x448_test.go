// The MIT License (MIT)
//
// Copyright (c) 2014-2015 Cryptography Research, Inc.
// Copyright (c) 2015-2019 Yawning Angel
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

package x448

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"testing"
)

// Cowardly refuse to run the full slow test vector case unless this is set
// at compile time, because the timeout for the test harness needs to be
// adjusted at runtime.
var reallyRunSlowTest = false

func TestX448(t *testing.T) {
	type KATVectors struct {
		scalar [x448Bytes]byte
		base   [x448Bytes]byte
		answer [x448Bytes]byte
	}

	vectors := []KATVectors{
		{
			[x448Bytes]byte{
				0x3d, 0x26, 0x2f, 0xdd, 0xf9, 0xec, 0x8e, 0x88,
				0x49, 0x52, 0x66, 0xfe, 0xa1, 0x9a, 0x34, 0xd2,
				0x88, 0x82, 0xac, 0xef, 0x04, 0x51, 0x04, 0xd0,
				0xd1, 0xaa, 0xe1, 0x21, 0x70, 0x0a, 0x77, 0x9c,
				0x98, 0x4c, 0x24, 0xf8, 0xcd, 0xd7, 0x8f, 0xbf,
				0xf4, 0x49, 0x43, 0xeb, 0xa3, 0x68, 0xf5, 0x4b,
				0x29, 0x25, 0x9a, 0x4f, 0x1c, 0x60, 0x0a, 0xd3,
			},
			[x448Bytes]byte{
				0x06, 0xfc, 0xe6, 0x40, 0xfa, 0x34, 0x87, 0xbf,
				0xda, 0x5f, 0x6c, 0xf2, 0xd5, 0x26, 0x3f, 0x8a,
				0xad, 0x88, 0x33, 0x4c, 0xbd, 0x07, 0x43, 0x7f,
				0x02, 0x0f, 0x08, 0xf9, 0x81, 0x4d, 0xc0, 0x31,
				0xdd, 0xbd, 0xc3, 0x8c, 0x19, 0xc6, 0xda, 0x25,
				0x83, 0xfa, 0x54, 0x29, 0xdb, 0x94, 0xad, 0xa1,
				0x8a, 0xa7, 0xa7, 0xfb, 0x4e, 0xf8, 0xa0, 0x86,
			},
			[x448Bytes]byte{
				0xce, 0x3e, 0x4f, 0xf9, 0x5a, 0x60, 0xdc, 0x66,
				0x97, 0xda, 0x1d, 0xb1, 0xd8, 0x5e, 0x6a, 0xfb,
				0xdf, 0x79, 0xb5, 0x0a, 0x24, 0x12, 0xd7, 0x54,
				0x6d, 0x5f, 0x23, 0x9f, 0xe1, 0x4f, 0xba, 0xad,
				0xeb, 0x44, 0x5f, 0xc6, 0x6a, 0x01, 0xb0, 0x77,
				0x9d, 0x98, 0x22, 0x39, 0x61, 0x11, 0x1e, 0x21,
				0x76, 0x62, 0x82, 0xf7, 0x3d, 0xd9, 0x6b, 0x6f,
			},
		},
		{
			[x448Bytes]byte{
				0x20, 0x3d, 0x49, 0x44, 0x28, 0xb8, 0x39, 0x93,
				0x52, 0x66, 0x5d, 0xdc, 0xa4, 0x2f, 0x9d, 0xe8,
				0xfe, 0xf6, 0x00, 0x90, 0x8e, 0x0d, 0x46, 0x1c,
				0xb0, 0x21, 0xf8, 0xc5, 0x38, 0x34, 0x5d, 0xd7,
				0x7c, 0x3e, 0x48, 0x06, 0xe2, 0x5f, 0x46, 0xd3,
				0x31, 0x5c, 0x44, 0xe0, 0xa5, 0xb4, 0x37, 0x12,
				0x82, 0xdd, 0x2c, 0x8d, 0x5b, 0xe3, 0x09, 0x5f,
			},
			[x448Bytes]byte{
				0x0f, 0xbc, 0xc2, 0xf9, 0x93, 0xcd, 0x56, 0xd3,
				0x30, 0x5b, 0x0b, 0x7d, 0x9e, 0x55, 0xd4, 0xc1,
				0xa8, 0xfb, 0x5d, 0xbb, 0x52, 0xf8, 0xe9, 0xa1,
				0xe9, 0xb6, 0x20, 0x1b, 0x16, 0x5d, 0x01, 0x58,
				0x94, 0xe5, 0x6c, 0x4d, 0x35, 0x70, 0xbe, 0xe5,
				0x2f, 0xe2, 0x05, 0xe2, 0x8a, 0x78, 0xb9, 0x1c,
				0xdf, 0xbd, 0xe7, 0x1c, 0xe8, 0xd1, 0x57, 0xdb,
			},
			[x448Bytes]byte{
				0x88, 0x4a, 0x02, 0x57, 0x62, 0x39, 0xff, 0x7a,
				0x2f, 0x2f, 0x63, 0xb2, 0xdb, 0x6a, 0x9f, 0xf3,
				0x70, 0x47, 0xac, 0x13, 0x56, 0x8e, 0x1e, 0x30,
				0xfe, 0x63, 0xc4, 0xa7, 0xad, 0x1b, 0x3e, 0xe3,
				0xa5, 0x70, 0x0d, 0xf3, 0x43, 0x21, 0xd6, 0x20,
				0x77, 0xe6, 0x36, 0x33, 0xc5, 0x75, 0xc1, 0xc9,
				0x54, 0x51, 0x4e, 0x99, 0xda, 0x7c, 0x17, 0x9d,
			},
		},
	}

	var out [x448Bytes]byte
	for i, vec := range vectors {
		ScalarMult(&out, &vec.scalar, &vec.base)
		if !bytes.Equal(out[:], vec.answer[:]) {
			t.Errorf("KAT[%d]: Mismatch", i)
		}
	}
}

func TestX448IETFDraft(t *testing.T) {
	// Run the other test vectors from 5.2 of the IETF draft.

	// WARNING: The full version of the test will easily take longer than the
	// default 10 min test timeout, even on a moderately powerful box.
	//
	// Unless reallyRunSlowTest is set in the source code, it will cowardly
	// refuse to run the full 1 million iterations, and the `go test`
	// timeout will need to be increased (`go test -timeout 30m`).

	var k, u, out [x448Bytes]byte
	copy(k[:], basePoint[:])
	copy(u[:], basePoint[:])

	for i := 0; i < 1000000; i++ {
		ScalarMult(&out, &k, &u)
		switch i + 1 {
		case 1:
			known, _ := hex.DecodeString("3f482c8a9f19b01e6c46ee9711d9dc14fd4bf67af30765c2ae2b846a4d23a8cd0db897086239492caf350b51f833868b9bc2b3bca9cf4113")
			if !bytes.Equal(out[:], known) {
				t.Fatalf("Iterated[%d]: Mismatch", i)
			}
		case 1000:
			known, _ := hex.DecodeString("aa3b4749d55b9daf1e5b00288826c467274ce3ebbdd5c17b975e09d4af6c67cf10d087202db88286e2b79fceea3ec353ef54faa26e219f38")
			if !bytes.Equal(out[:], known) {
				t.Fatalf("Iterated[%d]: Mismatch", i)
			}
			if testing.Short() || !reallyRunSlowTest {
				t.Skipf("Short test requested, skipping remaining, was correct at 1k")
			}
		}
		copy(u[:], k[:])
		copy(k[:], out[:])
	}
	known, _ := hex.DecodeString("077f453681caca3693198420bbe515cae0002472519b3e67661a7e89cab94695c8f4bcd66e61b9b9c946da8d524de3d69bd9d9d66b997e37")
	if !bytes.Equal(k[:], known) {
		t.Fatal("Final value mismatch")
	}
}

func TestCurve448(t *testing.T) {
	alicePriv := [x448Bytes]byte{
		0x9a, 0x8f, 0x49, 0x25, 0xd1, 0x51, 0x9f, 0x57,
		0x75, 0xcf, 0x46, 0xb0, 0x4b, 0x58, 0x00, 0xd4,
		0xee, 0x9e, 0xe8, 0xba, 0xe8, 0xbc, 0x55, 0x65,
		0xd4, 0x98, 0xc2, 0x8d, 0xd9, 0xc9, 0xba, 0xf5,
		0x74, 0xa9, 0x41, 0x97, 0x44, 0x89, 0x73, 0x91,
		0x00, 0x63, 0x82, 0xa6, 0xf1, 0x27, 0xab, 0x1d,
		0x9a, 0xc2, 0xd8, 0xc0, 0xa5, 0x98, 0x72, 0x6b,
	}

	alicePub := [x448Bytes]byte{
		0x9b, 0x08, 0xf7, 0xcc, 0x31, 0xb7, 0xe3, 0xe6,
		0x7d, 0x22, 0xd5, 0xae, 0xa1, 0x21, 0x07, 0x4a,
		0x27, 0x3b, 0xd2, 0xb8, 0x3d, 0xe0, 0x9c, 0x63,
		0xfa, 0xa7, 0x3d, 0x2c, 0x22, 0xc5, 0xd9, 0xbb,
		0xc8, 0x36, 0x64, 0x72, 0x41, 0xd9, 0x53, 0xd4,
		0x0c, 0x5b, 0x12, 0xda, 0x88, 0x12, 0x0d, 0x53,
		0x17, 0x7f, 0x80, 0xe5, 0x32, 0xc4, 0x1f, 0xa0,
	}

	bobPriv := [x448Bytes]byte{
		0x1c, 0x30, 0x6a, 0x7a, 0xc2, 0xa0, 0xe2, 0xe0,
		0x99, 0x0b, 0x29, 0x44, 0x70, 0xcb, 0xa3, 0x39,
		0xe6, 0x45, 0x37, 0x72, 0xb0, 0x75, 0x81, 0x1d,
		0x8f, 0xad, 0x0d, 0x1d, 0x69, 0x27, 0xc1, 0x20,
		0xbb, 0x5e, 0xe8, 0x97, 0x2b, 0x0d, 0x3e, 0x21,
		0x37, 0x4c, 0x9c, 0x92, 0x1b, 0x09, 0xd1, 0xb0,
		0x36, 0x6f, 0x10, 0xb6, 0x51, 0x73, 0x99, 0x2d,
	}

	bobPub := [x448Bytes]byte{
		0x3e, 0xb7, 0xa8, 0x29, 0xb0, 0xcd, 0x20, 0xf5,
		0xbc, 0xfc, 0x0b, 0x59, 0x9b, 0x6f, 0xec, 0xcf,
		0x6d, 0xa4, 0x62, 0x71, 0x07, 0xbd, 0xb0, 0xd4,
		0xf3, 0x45, 0xb4, 0x30, 0x27, 0xd8, 0xb9, 0x72,
		0xfc, 0x3e, 0x34, 0xfb, 0x42, 0x32, 0xa1, 0x3c,
		0xa7, 0x06, 0xdc, 0xb5, 0x7a, 0xec, 0x3d, 0xae,
		0x07, 0xbd, 0xc1, 0xc6, 0x7b, 0xf3, 0x36, 0x09,
	}

	aliceBob := [x448Bytes]byte{
		0x07, 0xff, 0xf4, 0x18, 0x1a, 0xc6, 0xcc, 0x95,
		0xec, 0x1c, 0x16, 0xa9, 0x4a, 0x0f, 0x74, 0xd1,
		0x2d, 0xa2, 0x32, 0xce, 0x40, 0xa7, 0x75, 0x52,
		0x28, 0x1d, 0x28, 0x2b, 0xb6, 0x0c, 0x0b, 0x56,
		0xfd, 0x24, 0x64, 0xc3, 0x35, 0x54, 0x39, 0x36,
		0x52, 0x1c, 0x24, 0x40, 0x30, 0x85, 0xd5, 0x9a,
		0x44, 0x9a, 0x50, 0x37, 0x51, 0x4a, 0x87, 0x9d,
	}

	var out [x448Bytes]byte
	ScalarBaseMult(&out, &alicePriv)
	if !bytes.Equal(out[:], alicePub[:]) {
		t.Error("Alice: ScalarBaseMult Mismatch")
	}
	ScalarBaseMult(&out, &bobPriv)
	if !bytes.Equal(out[:], bobPub[:]) {
		t.Error("Bob: ScalarBaseMult Mismatch")
	}
	ScalarMult(&out, &bobPriv, &alicePub)
	if !bytes.Equal(out[:], aliceBob[:]) {
		t.Error("Bob: ScalarMult Mismatch")
	}
	ScalarMult(&out, &alicePriv, &bobPub)
	if !bytes.Equal(out[:], aliceBob[:]) {
		t.Error("Alice: ScalarMult Mismatch")
	}
}

func BenchmarkECDH(b *testing.B) {
	var sa, sb, pa, pb, ab, ba [x448Bytes]byte

	_, _ = rand.Read(sa[:])
	_, _ = rand.Read(sb[:])
	b.ResetTimer()
	b.StopTimer()
	for i := 0; i < b.N; i++ {
		ScalarBaseMult(&pa, &sa)
		ScalarBaseMult(&pb, &sb)
		b.StartTimer()
		ScalarMult(&ab, &sa, &pb)
		b.StopTimer()
		ScalarMult(&ba, &sb, &pa)
		if !bytes.Equal(ab[:], ba[:]) {
			b.Fatal("Alice/Bob: Mismatch")
		}
		copy(sa[:], pa[:])
		copy(sb[:], pb[:])
	}
}
