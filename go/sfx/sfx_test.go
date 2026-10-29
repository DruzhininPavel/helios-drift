package sfx

import (
	"math"
	"testing"
)

// Every effect must actually contain signal: non-empty, audible peak, and no
// gross clipping (the mixer clamps, but the raw voice should stay near unity).
func TestSynthVoices(t *testing.T) {
	cases := map[string][]float64{
		"laser":   sweep(1500, 260, 0.075, 0, 0.55),
		"hit":     noiseBurst(0.11, 0.9, 2600, 11),
		"explode": mix(noiseBurst(0.65, 1.0, 700, 23), sweep(120, 38, 0.6, 2, 0.7)),
		"missile": mix(sweep(240, 900, 0.5, 1, 0.35), noiseBurst(0.5, 0.35, 1500, 31)),
		"ecm":     sweep(900, 260, 0.35, 3, 0.4),
		"bomb":    mix(noiseBurst(1.1, 0.9, 500, 41), sweep(200, 30, 1.1, 1, 0.5)),
		"hyper":   mix(sweep(110, 1700, 1.7, 1, 0.4), noiseBurst(1.7, 0.25, 2200, 53)),
	}
	for name, s := range cases {
		if len(s) < SampleRate/20 {
			t.Errorf("%s: too short (%d samples)", name, len(s))
			continue
		}
		peak := 0.0
		energy := 0.0
		for _, v := range s {
			a := math.Abs(v)
			if a > peak {
				peak = a
			}
			energy += v * v
		}
		rms := math.Sqrt(energy / float64(len(s)))
		if peak < 0.05 {
			t.Errorf("%s: silent (peak %.3f)", name, peak)
		}
		if peak > 1.6 {
			t.Errorf("%s: clipping hard (peak %.2f)", name, peak)
		}
		if rms < 0.01 {
			t.Errorf("%s: almost no energy (rms %.4f)", name, rms)
		}
		t.Logf("%-8s %6.0f ms  peak %.2f  rms %.3f", name, float64(len(s))/SampleRate*1000, peak, rms)
	}
}

func TestPCMEncoding(t *testing.T) {
	s := []float64{0, 0.5, -0.5, 1.5}
	b := pcm(s)
	if len(b) != len(s)*4 {
		t.Fatalf("expected %d bytes, got %d", len(s)*4, len(b))
	}
	// 0.5 -> 16000, channels interleaved
	if b[0] != 0 || b[2] != 0 {
		t.Errorf("dc sample not zero: %v %v", b[0], b[2])
	}
	left := int16(uint16(b[4]) | uint16(b[5])<<8)
	if left != 16000 {
		t.Errorf("expected 16000 for 0.5, got %d", left)
	}
	// clipping guard
	clipped := int16(uint16(b[12]) | uint16(b[13])<<8)
	if clipped != 32000 {
		t.Errorf("expected clamp to 32000, got %d", clipped)
	}
}

func TestEngineLoopSeam(t *testing.T) {
	eng := mix(sweep(58, 58, 1.0, 1, 0.5), noiseBurst(1.0, 0.35, 240, 67))
	n := len(eng)
	for i := 0; i < n; i++ {
		e := math.Min(1, math.Min(float64(i), float64(n-1-i))/float64(SampleRate/60))
		eng[i] *= e
	}
	if math.Abs(eng[0]) > 0.001 || math.Abs(eng[n-1]) > 0.001 {
		t.Errorf("loop seam not tapered: start %.3f end %.3f", eng[0], eng[n-1])
	}
	mid := eng[n/2]
	if math.Abs(mid) < 0.01 {
		t.Errorf("engine bed silent in the middle (%.3f)", mid)
	}
}
