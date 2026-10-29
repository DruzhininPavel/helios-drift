// Package sfx synthesises the game's sound effects as 16-bit PCM and plays
// them through Ebiten's audio context. Everything is generated at start-up:
// square/saw/sine tones plus filtered noise, 16-bit console style.
package sfx

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"

	"github.com/hajimehoshi/ebiten/v2/audio"
)

const SampleRate = 44100

var ctx *audio.Context

// Bank holds one voice pool per effect.
type Bank struct {
	Laser   *Voice
	Hit     *Voice
	Explode *Voice
	Missile *Voice
	ECM     *Voice
	Bomb    *Voice
	Dock    *Voice
	Hyper   *Voice
	Warning *Voice
	Engine  *Loop
	Muted   bool
	Master  float64
}

var B *Bank

// Voice is a round-robin pool so a sound can retrigger while it is playing.
type Voice struct {
	players []*audio.Player
	next    int
	vol     float64
	base    float64
}

func (v *Voice) Play() {
	if v == nil || B == nil || B.Muted || len(v.players) == 0 {
		return
	}
	p := v.players[v.next]
	v.next = (v.next + 1) % len(v.players)
	p.SetVolume(v.base * B.Master)
	if p.IsPlaying() {
		_ = p.Rewind()
	}
	_ = p.Seek(0)
	p.Play()
}

// Loop is a continuously looping bed (the engine hum).
type Loop struct {
	player *audio.Player
	target float64
	cur    float64
}

func (l *Loop) Start() {
	if l == nil || l.player == nil {
		return
	}
	l.player.SetVolume(0)
	l.player.Play()
}

func (l *Loop) SetVolume(v float64) {
	if l == nil {
		return
	}
	if B != nil && B.Muted {
		v = 0
	}
	l.target = v
}

// Update smooths the loop volume so it does not click.
func (l *Loop) Update() {
	if l == nil || l.player == nil {
		return
	}
	l.cur += (l.target - l.cur) * 0.12
	l.player.SetVolume(l.cur)
}

// ---------------------------------------------------------------- synthesis

func env(i, n int) float64 {
	// fast attack, exponential decay
	t := float64(i) / float64(n)
	a := math.Min(1, t*40)
	return a * math.Exp(-3.2*t)
}

func oscillator(wave int, phase float64) float64 {
	switch wave {
	case 0: // square
		if math.Mod(phase, 1) < 0.5 {
			return 1
		}
		return -1
	case 1: // saw
		return 2*math.Mod(phase, 1) - 1
	case 2: // sine
		return math.Sin(phase * 2 * math.Pi)
	case 3: // triangle
		p := math.Mod(phase, 1)
		if p < 0.5 {
			return 4*p - 1
		}
		return 3 - 4*p
	}
	return 0
}

func pcm(samples []float64) []byte {
	out := make([]byte, len(samples)*4) // stereo 16-bit
	for i, s := range samples {
		if s > 1 {
			s = 1
		}
		if s < -1 {
			s = -1
		}
		v := int16(s * 32000)
		out[i*4+0] = byte(v)
		out[i*4+1] = byte(v >> 8)
		out[i*4+2] = byte(v)
		out[i*4+3] = byte(v >> 8)
	}
	return out
}

// sweep makes a tone gliding from f0 to f1 with a decay envelope.
func sweep(f0, f1, dur float64, wave int, vol float64) []float64 {
	n := int(dur * SampleRate)
	out := make([]float64, n)
	phase := 0.0
	for i := 0; i < n; i++ {
		t := float64(i) / float64(n)
		f := f0 * math.Pow(f1/f0, t)
		phase += f / SampleRate
		out[i] = oscillator(wave, phase) * vol * env(i, n)
	}
	return out
}

// noiseBurst makes filtered noise with a decay envelope.
func noiseBurst(dur float64, vol, cutoff float64, seed int64) []float64 {
	n := int(dur * SampleRate)
	out := make([]float64, n)
	r := rand.New(rand.NewSource(seed))
	prev := 0.0
	k := cutoff / float64(SampleRate) * 6.0
	if k > 1 {
		k = 1
	}
	for i := 0; i < n; i++ {
		w := r.Float64()*2 - 1
		prev += (w - prev) * k
		out[i] = prev * vol * env(i, n)
	}
	return out
}

func mix(a, b []float64) []float64 {
	n := len(a)
	if len(b) > n {
		n = len(b)
	}
	out := make([]float64, n)
	for i := 0; i < n; i++ {
		var s float64
		if i < len(a) {
			s += a[i]
		}
		if i < len(b) {
			s += b[i]
		}
		out[i] = s
	}
	return out
}

func silence(dur float64) []float64 {
	return make([]float64, int(dur*SampleRate))
}

func concat(parts ...[]float64) []float64 {
	var out []float64
	for _, p := range parts {
		out = append(out, p...)
	}
	return out
}

func newVoice(samples []float64, voices int) *Voice { return newVoiceVol(samples, voices, 0.55) }

func newVoiceVol(samples []float64, voices int, vol float64) *Voice {
	data := pcm(samples)
	v := &Voice{vol: 1, base: vol}
	for i := 0; i < voices; i++ {
		p := ctx.NewPlayerFromBytes(data)
		p.SetVolume(vol)
		v.players = append(v.players, p)
	}
	return v
}

// SetMaster scales every voice (0..1).
func (b *Bank) SetMaster(m float64) {
	if b == nil {
		return
	}
	if m < 0 {
		m = 0
	}
	if m > 1 {
		m = 1
	}
	b.Master = m
	for _, v := range []*Voice{b.Laser, b.Hit, b.Explode, b.Missile, b.ECM, b.Bomb, b.Dock, b.Hyper, b.Warning} {
		if v == nil {
			continue
		}
		for _, p := range v.players {
			p.SetVolume(v.base * m)
		}
	}
}

// Init creates the audio context and synthesises every effect. On a machine
// with no sound server it returns an empty bank: every voice is nil and Play
// is a no-op, so the game keeps running silently.
func Init() *Bank {
	if !AudioAvailable() || os.Getenv("HELIOS_NO_AUDIO") != "" {
		fmt.Println("audio: no sound server found, running silent")
		B = &Bank{}
		return B
	}
	ctx = audio.NewContext(SampleRate)
	b := &Bank{Master: 0.8}
	b.Laser = newVoiceVol(sweep(1500, 260, 0.075, 0, 0.55), 4, 0.35)
	b.Hit = newVoice(noiseBurst(0.11, 0.9, 2600, 11), 3)
	b.Explode = newVoice(mix(noiseBurst(0.65, 1.0, 700, 23), sweep(120, 38, 0.6, 2, 0.7)), 3)
	b.Missile = newVoice(mix(sweep(240, 900, 0.5, 1, 0.35), noiseBurst(0.5, 0.35, 1500, 31)), 2)
	b.ECM = newVoice(sweep(900, 260, 0.35, 3, 0.4), 2)
	b.Bomb = newVoice(mix(noiseBurst(1.1, 0.9, 500, 41), sweep(200, 30, 1.1, 1, 0.5)), 1)
	b.Dock = newVoice(concat(sweep(520, 520, 0.12, 0, 0.4), silence(0.04), sweep(780, 780, 0.2, 0, 0.4)), 1)
	b.Hyper = newVoice(mix(sweep(110, 1700, 1.7, 1, 0.4), noiseBurst(1.7, 0.25, 2200, 53)), 1)
	b.Warning = newVoice(concat(sweep(520, 520, 0.09, 0, 0.45), silence(0.09),
		sweep(520, 520, 0.09, 0, 0.45), silence(0.09), sweep(520, 520, 0.09, 0, 0.45)), 2)

	// engine bed: low rumble, looped forever
	eng := mix(sweep(58, 58, 1.0, 1, 0.5), noiseBurst(1.0, 0.35, 240, 67))
	// taper the loop seam
	n := len(eng)
	for i := 0; i < n; i++ {
		e := math.Min(1, math.Min(float64(i), float64(n-1-i))/float64(SampleRate/60))
		eng[i] *= e
	}
	loopData := pcm(eng)
	src := bytes.NewReader(loopData)
	infinite := audio.NewInfiniteLoop(src, int64(len(loopData)))
	p, err := ctx.NewPlayer(infinite)
	if err == nil {
		b.Engine = &Loop{player: p}
		b.Engine.Start()
	}
	B = b
	return b
}

// Update pumps the engine volume; call once per frame.
func (b *Bank) Update(engineLevel float64) {
	if b == nil || b.Engine == nil {
		return
	}
	b.Engine.SetVolume(engineLevel * b.Master)
	b.Engine.Update()
}

// ToggleMute flips global muting and returns the new state.
func (b *Bank) ToggleMute() bool {
	if b == nil {
		return true
	}
	b.Muted = !b.Muted
	if b.Muted && b.Engine != nil {
		b.Engine.SetVolume(0)
		b.Engine.Update()
	}
	return b.Muted
}
