/* ELITE 16 - boot, input, main loop, sound */
(function (global) {
  'use strict';
  var E = global.E, GAME = E.game, UI = E.ui;

  var W = 320, H = 224;
  var canvas, ctx, fb, img, running = false;
  var keys = {};
  var last = 0;
  var mouse = { x: 0, y: 0, inside: false, active: 0 };
  var scale = 1;

  // ------------------------------------------------------------------ audio
  var AudioMod = (function () {
    var ctxA = null, master = null, muted = false;
    function ensure() {
      if (ctxA) return ctxA;
      try {
        var AC = global.AudioContext || global.webkitAudioContext;
        ctxA = new AC();
        master = ctxA.createGain();
        master.gain.value = 0.22;
        master.connect(ctxA.destination);
      } catch (e) { ctxA = null; }
      return ctxA;
    }
    function blip(freq, dur, type, vol, slide) {
      if (muted) return;
      var a = ensure(); if (!a) return;
      var o = a.createOscillator(), g = a.createGain();
      o.type = type || 'square';
      o.frequency.setValueAtTime(freq, a.currentTime);
      if (slide) o.frequency.exponentialRampToValueAtTime(Math.max(30, slide), a.currentTime + dur);
      g.gain.setValueAtTime(vol === undefined ? 0.5 : vol, a.currentTime);
      g.gain.exponentialRampToValueAtTime(0.001, a.currentTime + dur);
      o.connect(g); g.connect(master);
      o.start(); o.stop(a.currentTime + dur + 0.02);
    }
    function noise(dur, vol, filterFreq) {
      if (muted) return;
      var a = ensure(); if (!a) return;
      var len = Math.floor(a.sampleRate * dur);
      var buf = a.createBuffer(1, len, a.sampleRate);
      var d = buf.getChannelData(0);
      for (var i = 0; i < len; i++) d[i] = (Math.random() * 2 - 1) * (1 - i / len);
      var src = a.createBufferSource();
      src.buffer = buf;
      var f = a.createBiquadFilter();
      f.type = 'lowpass';
      f.frequency.value = filterFreq || 1800;
      var g = a.createGain();
      g.gain.value = vol === undefined ? 0.6 : vol;
      src.connect(f); f.connect(g); g.connect(master);
      src.start();
    }
    return {
      laser: function () { blip(1400, 0.07, 'square', 0.35, 260); },
      hit: function () { noise(0.12, 0.5, 2600); },
      explode: function () { noise(0.55, 0.9, 900); blip(90, 0.4, 'sawtooth', 0.4, 40); },
      missile: function () { blip(300, 0.5, 'sawtooth', 0.3, 900); },
      dock: function () { blip(520, 0.15, 'square', 0.3); setTimeout(function () { blip(780, 0.25, 'square', 0.3); }, 140); },
      hyper: function () { blip(120, 1.6, 'sawtooth', 0.35, 1800); },
      error: function () { blip(180, 0.2, 'square', 0.35, 90); },
      toggleMute: function () { muted = !muted; return muted; },
      unlock: function () { var a = ensure(); if (a && a.state === 'suspended') a.resume(); }
    };
  })();
  global.SFX = AudioMod;

  // ------------------------------------------------------------------- boot
  function boot() {
    if (running) return;
    canvas = document.getElementById('screen');
    if (!canvas) return;
    canvas.width = W; canvas.height = H;
    ctx = canvas.getContext('2d', { alpha: false });
    ctx.imageSmoothingEnabled = false;
    img = ctx.createImageData(W, H);
    fb = new E.FB(W, H);
    E.fb = fb;
    resize();
    global.addEventListener('resize', resize);
    bindInput();
    GAME.state = 'title';
    GAME.time = 0;
    GAME.onEvent = function (name) {
      if (name === 'laser') { AudioMod.laser(); GAME.beamFlash = 0.05; }
      else if (name === 'explode') { AudioMod.explode(); GAME.shake = 0.7; }
      else if (name === 'hit') AudioMod.hit();
      else if (name === 'dock') AudioMod.dock();
    };
    last = performance.now();
    running = true;
    requestAnimationFrame(frame);
  }

  function resize() {
    var vw = global.innerWidth - 24, vh = global.innerHeight - 96;
    scale = Math.max(1, Math.min(Math.floor(vw / W), Math.floor(vh / H)));
    canvas.style.width = (W * scale) + 'px';
    canvas.style.height = (H * scale) + 'px';
  }

  // ------------------------------------------------------------------ input
  function bindInput() {
    global.addEventListener('keydown', function (e) {
      AudioMod.unlock();
      var k = e.key;
      if (k === ' ' || k === 'Tab' || k.indexOf('Arrow') === 0) e.preventDefault();
      if (e.repeat) return;
      keys[k.toLowerCase ? k.toLowerCase() : k] = true;
      keys[k] = true;

      var handled = false;
      if (GAME.state !== 'flight') {
        handled = UI.key(k);
      } else if (GAME.overlay) {
        handled = overlayKey(k);
      } else {
        if (k === 'c' || k === 'C') { GAME.overlay = 'chart'; UI.sel.chartCursor = GAME.cmd.system; handled = true; }
        else if (k === 'v' || k === 'V') { GAME.overlay = 'status'; handled = true; }
        else if (k === 'Escape') { GAME.overlay = 'pause'; handled = true; }
        else if (k === 'm' || k === 'M') { GAME.fireMissile(); handled = true; }
        else if (k === 'e' || k === 'E') { GAME.fireECM(); handled = true; }
        else if (k === 'Tab') { GAME.fireBomb(); handled = true; }
        else if (k === 'l' || k === 'L') { GAME.toggleAutopilot(); handled = true; }
        else if (k === 'j' || k === 'J') { GAME.jumpInSystem(); handled = true; }
      }
      if (k.toLowerCase() === 'm' && e.shiftKey) { /* shift+M mutes */ }
      if (k === 'MediaMute') handled = true;
      if (handled) e.preventDefault();
    });

    global.addEventListener('keyup', function (e) {
      keys[e.key] = false;
      if (e.key.toLowerCase) keys[e.key.toLowerCase()] = false;
    });

    canvas.addEventListener('mousemove', function (e) {
      var r = canvas.getBoundingClientRect();
      mouse.x = (e.clientX - r.left) / r.width;
      mouse.y = (e.clientY - r.top) / r.height;
      mouse.inside = true;
      mouse.active = 0.35;
    });
    canvas.addEventListener('mouseleave', function () { mouse.inside = false; });
    canvas.addEventListener('mousedown', function (e) {
      AudioMod.unlock();
      if (e.button === 0) { keys[' '] = true; e.preventDefault(); }
    });
    global.addEventListener('mouseup', function () { keys[' '] = false; });
    canvas.addEventListener('contextmenu', function (e) { e.preventDefault(); });
    canvas.addEventListener('touchstart', function (e) { AudioMod.unlock(); }, { passive: true });
  }

  function overlayKey(k) {
    if (k === 'Escape') { GAME.overlay = null; return true; }
    if (GAME.overlay === 'chart') {
      var gal = E.galaxy.get(GAME.cmd.galaxy);
      var curId = UI.sel.chartCursor === null ? GAME.cmd.system : UI.sel.chartCursor;
      if (k.indexOf('Arrow') === 0) {
        UI.sel.chartCursor = UI.nearestInDirection(gal, curId, k);
        return true;
      }
      if (k === 'h' || k === 'H') {
        if (GAME.hyperspace(UI.sel.chartCursor)) { GAME.overlay = null; AudioMod.hyper(); }
        return true;
      }
      if (k === 'c' || k === 'C') { GAME.overlay = null; return true; }
      return true;
    }
    if (k === 'v' || k === 'V') { GAME.overlay = null; return true; }
    return true;
  }

  function readInput(dt) {
    if (GAME.state !== 'flight' || GAME.overlay) {
      GAME.in.pitch = 0; GAME.in.roll = 0; GAME.in.accel = 0; GAME.in.fire = false;
      return;
    }
    var pitch = 0, roll = 0, accel = 0;
    if (keys['ArrowUp'] || keys['w']) pitch += 1;
    if (keys['ArrowDown'] || keys['s']) pitch -= 1;
    if (keys['ArrowLeft'] || keys['a']) roll -= 1;
    if (keys['ArrowRight'] || keys['d']) roll += 1;
    if (keys['x'] || keys['PageUp'] || keys['=']) accel += 1;
    if (keys['z'] || keys['PageDown'] || keys['-']) accel -= 1;

    if (pitch === 0 && roll === 0) {
      mouse.active -= dt;
      if (mouse.inside && mouse.active > 0) {
        var nx = (mouse.x - 0.5) * 2, ny = (mouse.y - 0.5) * 2;
        if (Math.abs(nx) > 0.06) roll = Math.max(-1, Math.min(1, (nx - Math.sign(nx) * 0.06) * 1.6));
        if (Math.abs(ny) > 0.06) pitch = -Math.max(-1, Math.min(1, (ny - Math.sign(ny) * 0.06) * 1.6));
        // invert / deadzone handled above
      }
    }
    GAME.in.pitch = pitch;
    GAME.in.roll = roll;
    GAME.in.accel = accel;
    GAME.in.fire = !!keys[' '];
  }

  // ------------------------------------------------------------------- loop
  function frame(now) {
    if (!running) return;
    var dt = (now - last) / 1000;
    last = now;
    if (dt > 0.06) dt = 0.06;
    if (dt < 0) dt = 0;

    readInput(dt);
    if (!GAME.paused && !GAME.overlay) {
      GAME.update(dt);
    } else {
      GAME.time += dt; // keep animations alive on menus
    }
    if (GAME.state === 'title') GAME.time += dt;

    fb.clear(E.SPACE);
    fb.resetClip();
    UI.draw(fb, dt);
    img.data.set(fb.bytes);
    ctx.putImageData(img, 0, 0);
    requestAnimationFrame(frame);
  }

  // expose for tests
  E.boot = boot;
  E.frameStep = function (dt) {
    readInput(dt);
    GAME.update(dt);
    UI.draw(fb, dt);
    return fb;
  };
  E.getFB = function () { return fb; };

  if (typeof document !== 'undefined') {
    if (document.readyState === 'loading') document.addEventListener('DOMContentLoaded', boot);
    else boot();
  }
})(typeof window !== 'undefined' ? window : globalThis);
