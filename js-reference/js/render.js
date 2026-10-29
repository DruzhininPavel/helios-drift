/* ELITE 16 - renderer
 * Software framebuffer + 16-bit console style rasterizer:
 *   320x224 internal resolution (Mega Drive), limited 12-ramp palette,
 *   4x4 ordered dithering for smooth flat shading, 5x7 bitmap font.
 */
(function (global) {
  'use strict';

  var E = global.E || (global.E = {});

  // ---------------------------------------------------------------- palette
  var RAMPS = {
    steel:  ['#070a10', '#1a2230', '#33415a', '#556a88', '#8fa3c0', '#dbe6f7'],
    ash:    ['#0a0a0c', '#22232a', '#43454f', '#6b6e7c', '#9ea2b2', '#e8ecf7'],
    red:    ['#170406', '#400c14', '#7a1a26', '#b83434', '#ee6a4e', '#ffc0a0'],
    orange: ['#1a0a02', '#4d2106', '#8c4310', '#c9761c', '#f5ae44', '#ffe1a8'],
    yellow: ['#1c1502', '#4d3a08', '#8a6d18', '#c9a63a', '#f2d96a', '#fff8c0'],
    green:  ['#03120a', '#0d331c', '#1a6032', '#2c944e', '#5ed27c', '#c0ffcf'],
    teal:   ['#02141a', '#073442', '#0f6274', '#1e93a8', '#3ecfdd', '#b6f6ff'],
    blue:   ['#04081a', '#0e1d52', '#1a37a0', '#2f60d8', '#5f9bff', '#c4dcff'],
    purple: ['#0d0420', '#260e4a', '#481d86', '#7136c2', '#a566f0', '#e2c8ff'],
    brown:  ['#150d05', '#372211', '#5e3c1e', '#8b6031', '#b98c52', '#e8caa0'],
    pink:   ['#1a0512', '#4a0f30', '#8a1f5c', '#c8388c', '#f56ab8', '#ffc6e6'],
    alien:  ['#06140c', '#12402a', '#227a48', '#3aa860', '#7ce07a', '#e0ffb0']
  };

  var RAMP_ID = {}, RAMP_NAMES = Object.keys(RAMPS);
  var PAL = [];
  function hexToU32(h) {
    var n = parseInt(h.slice(1), 16);
    var r = (n >> 16) & 255, g = (n >> 8) & 255, b = n & 255;
    return ((255 << 24) | (b << 16) | (g << 8) | r) >>> 0;
  }
  RAMP_NAMES.forEach(function (name, i) {
    RAMP_ID[name] = i;
    RAMPS[name].forEach(function (h) { PAL.push(hexToU32(h)); });
  });
  PAL = new Uint32Array(PAL);

  var SPACE = hexToU32('#05070f');
  var SPACE_HI = hexToU32('#0d1424');

  // ------------------------------------------------------------- utilities
  var BAYER = [0, 8, 2, 10, 12, 4, 14, 6, 3, 11, 1, 9, 15, 7, 13, 5];

  function clamp(v, a, b) { return v < a ? a : (v > b ? b : v); }
  function lerp(a, b, t) { return a + (b - a) * t; }

  E.clamp = clamp;
  E.lerp = lerp;
  E.rgb = hexToU32;
  E.RAMPS = RAMPS;
  E.RAMP_ID = RAMP_ID;
  E.PAL = PAL;
  E.SPACE = SPACE;

  E.color = function (ramp, shade) {
    var rid = typeof ramp === 'string' ? RAMP_ID[ramp] : ramp;
    var s = clamp(Math.round(shade), 0, 5);
    return PAL[rid * 6 + s];
  };

  // stable rng
  E.rng = function (seed) {
    var a = seed >>> 0;
    return function () {
      a |= 0; a = (a + 0x6D2B79F5) | 0;
      var t = Math.imul(a ^ (a >>> 15), 1 | a);
      t = (t + Math.imul(t ^ (t >>> 7), 61 | t)) ^ t;
      return ((t ^ (t >>> 14)) >>> 0) / 4294967296;
    };
  };

  // ------------------------------------------------------------------ font
  // 5x7 glyphs
  var GLYPHS = {
    '0': '01110 10001 10011 10101 11001 10001 01110',
    '1': '00100 01100 00100 00100 00100 00100 01110',
    '2': '01110 10001 00001 00010 00100 01000 11111',
    '3': '11110 00001 00001 01110 00001 00001 11110',
    '4': '00010 00110 01010 10010 11111 00010 00010',
    '5': '11111 10000 11110 00001 00001 10001 01110',
    '6': '00110 01000 10000 11110 10001 10001 01110',
    '7': '11111 00001 00010 00100 01000 01000 01000',
    '8': '01110 10001 10001 01110 10001 10001 01110',
    '9': '01110 10001 10001 01111 00001 00010 01100',
    'A': '01110 10001 10001 11111 10001 10001 10001',
    'B': '11110 10001 10001 11110 10001 10001 11110',
    'C': '01110 10001 10000 10000 10000 10001 01110',
    'D': '11100 10010 10001 10001 10001 10010 11100',
    'E': '11111 10000 10000 11110 10000 10000 11111',
    'F': '11111 10000 10000 11110 10000 10000 10000',
    'G': '01110 10001 10000 10111 10001 10001 01111',
    'H': '10001 10001 10001 11111 10001 10001 10001',
    'I': '01110 00100 00100 00100 00100 00100 01110',
    'J': '00111 00010 00010 00010 00010 10010 01100',
    'K': '10001 10010 10100 11000 10100 10010 10001',
    'L': '10000 10000 10000 10000 10000 10000 11111',
    'M': '10001 11011 10101 10101 10001 10001 10001',
    'N': '10001 10001 11001 10101 10011 10001 10001',
    'O': '01110 10001 10001 10001 10001 10001 01110',
    'P': '11110 10001 10001 11110 10000 10000 10000',
    'Q': '01110 10001 10001 10001 10101 10010 01101',
    'R': '11110 10001 10001 11110 10100 10010 10001',
    'S': '01111 10000 10000 01110 00001 00001 11110',
    'T': '11111 00100 00100 00100 00100 00100 00100',
    'U': '10001 10001 10001 10001 10001 10001 01110',
    'V': '10001 10001 10001 10001 10001 01010 00100',
    'W': '10001 10001 10001 10101 10101 11011 10001',
    'X': '10001 10001 01010 00100 01010 10001 10001',
    'Y': '10001 10001 01010 00100 00100 00100 00100',
    'Z': '11111 00001 00010 00100 01000 10000 11111',
    ' ': '00000 00000 00000 00000 00000 00000 00000',
    '.': '00000 00000 00000 00000 00000 01100 01100',
    ',': '00000 00000 00000 00000 00110 00100 01000',
    ':': '00000 00100 00100 00000 00100 00100 00000',
    ';': '00000 00100 00100 00000 00100 00100 01000',
    '-': '00000 00000 00000 11111 00000 00000 00000',
    '_': '00000 00000 00000 00000 00000 00000 11111',
    '+': '00000 00100 00100 11111 00100 00100 00000',
    '/': '00001 00010 00010 00100 01000 01000 10000',
    '\\': '10000 01000 01000 00100 00010 00010 00001',
    '%': '11001 11010 00010 00100 01000 01011 10011',
    '(': '00010 00100 01000 01000 01000 00100 00010',
    ')': '01000 00100 00010 00010 00010 00100 01000',
    '[': '01110 01000 01000 01000 01000 01000 01110',
    ']': '01110 00010 00010 00010 00010 00010 01110',
    '*': '00000 00100 10101 01110 10101 00100 00000',
    '\'': '00100 00100 00000 00000 00000 00000 00000',
    '!': '00100 00100 00100 00100 00100 00000 00100',
    '?': '01110 10001 00001 00110 00100 00000 00100',
    '=': '00000 00000 11111 00000 11111 00000 00000',
    '>': '10000 01000 00100 00010 00100 01000 10000',
    '<': '00001 00010 00100 01000 00100 00010 00001',
    '^': '00100 01010 10001 00000 00000 00000 00000',
    '@': '01110 10001 10111 10101 10111 10000 01110'
  };
  var GLYPH_BITS = {};
  for (var gk in GLYPHS) {
    GLYPH_BITS[gk] = GLYPHS[gk].split(' ').map(function (row) {
      var bits = 0;
      for (var i = 0; i < 5; i++) if (row[i] === '1') bits |= (1 << (4 - i));
      return bits;
    });
  }
  E.GLYPH_BITS = GLYPH_BITS;

  // -------------------------------------------------------------- framebuf
  function FB(w, h) {
    this.w = w; this.h = h;
    this.buf = new Uint32Array(w * h);
    this.bytes = new Uint8ClampedArray(this.buf.buffer);
    this.zbuf = new Float32Array(w * h);   // 1/depth, 0 = empty
    this.clipY0 = 0; this.clipY1 = h - 1; this.clipX0 = 0; this.clipX1 = w - 1;
  }
  FB.prototype.setClip = function (x0, y0, x1, y1) {
    this.clipX0 = clamp(x0, 0, this.w - 1); this.clipY0 = clamp(y0, 0, this.h - 1);
    this.clipX1 = clamp(x1, 0, this.w - 1); this.clipY1 = clamp(y1, 0, this.h - 1);
  };
  FB.prototype.resetClip = function () {
    this.clipX0 = 0; this.clipY0 = 0; this.clipX1 = this.w - 1; this.clipY1 = this.h - 1;
  };
  FB.prototype.clear = function (c) {
    this.buf.fill(c === undefined ? SPACE : c);
    this.zbuf.fill(0);
  };
  FB.prototype.clearGradient = function (top, bottom) {
    this.zbuf.fill(0);
    var buf = this.buf, w = this.w;
    for (var y = 0; y < this.h; y++) {
      var t = y / (this.h - 1);
      var i0 = Math.floor(t * 5), f = t * 5 - i0;
      var ca = top[clamp(i0, 0, 5)], cb = top[clamp(i0 + 1, 0, 5)];
      var base = y * w;
      for (var x = 0; x < w; x++) {
        buf[base + x] = (BAYER[(y & 3) * 4 + (x & 3)] < f * 16) ? cb : ca;
      }
    }
  };
  FB.prototype.px = function (x, y, c) {
    x = x | 0; y = y | 0;
    if (x < this.clipX0 || x > this.clipX1 || y < this.clipY0 || y > this.clipY1) return;
    this.buf[y * this.w + x] = c;
  };
  FB.prototype.get = function (x, y) {
    x = x | 0; y = y | 0;
    if (x < 0 || x >= this.w || y < 0 || y >= this.h) return 0;
    return this.buf[y * this.w + x];
  };
  FB.prototype.hline = function (x0, x1, y, c) {
    if (y < this.clipY0 || y > this.clipY1) return;
    if (x0 > x1) { var t = x0; x0 = x1; x1 = t; }
    x0 = Math.max(x0, this.clipX0); x1 = Math.min(x1, this.clipX1);
    if (x1 < x0) return;
    var row = y * this.w;
    this.buf.fill(c, row + x0, row + x1 + 1);
  };
  FB.prototype.vline = function (x, y0, y1, c) {
    if (x < this.clipX0 || x > this.clipX1) return;
    y0 = Math.max(y0, this.clipY0); y1 = Math.min(y1, this.clipY1);
    for (var y = y0; y <= y1; y++) this.buf[y * this.w + x] = c;
  };
  FB.prototype.rect = function (x0, y0, x1, y1, c) {
    this.hline(x0, x1, y0, c); this.hline(x0, x1, y1, c);
    this.vline(x0, y0, y1, c); this.vline(x1, y0, y1, c);
  };
  FB.prototype.fillRect = function (x0, y0, x1, y1, c) {
    x0 = clamp(x0 | 0, this.clipX0, this.clipX1); x1 = clamp(x1 | 0, this.clipX0, this.clipX1);
    y0 = clamp(y0 | 0, this.clipY0, this.clipY1); y1 = clamp(y1 | 0, this.clipY0, this.clipY1);
    if (x1 < x0 || y1 < y0) return;
    for (var y = y0; y <= y1; y++) this.buf.fill(c, y * this.w + x0, y * this.w + x1 + 1);
  };
  FB.prototype.line = function (x0, y0, x1, y1, c) {
    x0 |= 0; y0 |= 0; x1 |= 0; y1 |= 0;
    var dx = Math.abs(x1 - x0), dy = Math.abs(y1 - y0);
    var sx = x0 < x1 ? 1 : -1, sy = y0 < y1 ? 1 : -1, err = dx - dy;
    for (var guard = 0; guard < 4000; guard++) {
      this.px(x0, y0, c);
      if (x0 === x1 && y0 === y1) break;
      var e2 = 2 * err;
      if (e2 > -dy) { err -= dy; x0 += sx; }
      if (e2 < dx) { err += dx; y0 += sy; }
    }
  };

  // scanline span filler; cb(y, x0, x1, iz0, iz1) with raw (unclipped) x and 1/z
  FB.prototype.spans = function (pts, cb) {
    var n = pts.length;
    if (n < 3) return;
    var minY = Infinity, maxY = -Infinity, i;
    var hasZ = pts[0].z !== undefined;
    for (i = 0; i < n; i++) {
      if (pts[i].y < minY) minY = pts[i].y;
      if (pts[i].y > maxY) maxY = pts[i].y;
    }
    if (!isFinite(minY) || !isFinite(maxY)) return;
    var y0 = Math.max(Math.ceil(minY), this.clipY0);
    var y1 = Math.min(Math.floor(maxY), this.clipY1);
    var xs = [];
    for (var y = y0; y <= y1; y++) {
      var sy = y + 0.5;
      xs.length = 0;
      for (i = 0; i < n; i++) {
        var a = pts[i], b = pts[(i + 1) % n];
        if ((a.y <= sy && b.y > sy) || (b.y <= sy && a.y > sy)) {
          var t = (sy - a.y) / (b.y - a.y);
          var iz = 0;
          if (hasZ) {
            var za = a.z > 0.001 ? a.z : 0.001, zb = b.z > 0.001 ? b.z : 0.001;
            iz = (1 - t) / za + t / zb;
          }
          xs.push({ x: a.x + t * (b.x - a.x), iz: iz });
        }
      }
      if (xs.length < 2) continue;
      xs.sort(function (p, q) { return p.x - q.x; });
      for (var k = 0; k + 1 < xs.length; k += 2) {
        cb(y, xs[k].x, xs[k + 1].x, xs[k].iz, xs[k + 1].iz);
      }
    }
  };

  FB.prototype.poly = function (pts, color) {
    var buf = this.buf, w = this.w, self = this;
    this.spans(pts, function (y, x0, x1) {
      var xa = Math.max(Math.ceil(x0 - 0.5), self.clipX0);
      var xb = Math.min(Math.floor(x1 - 0.5), self.clipX1);
      if (xb >= xa) buf.fill(color, y * w + xa, y * w + xb + 1);
    });
  };

  // flat-shaded polygon with ordered dithering between two palette shades.
  // When the input points carry .z (depth) the fill is depth-tested per pixel.
  FB.prototype.polyShaded = function (pts, rampId, shade) {
    var s = clamp(shade, 0, 5);
    var i0 = Math.floor(s); if (i0 > 5) i0 = 5;
    var f = s - i0, i1 = Math.min(5, i0 + 1);
    var c0 = PAL[rampId * 6 + i0], c1 = PAL[rampId * 6 + i1];
    var buf = this.buf, zb = this.zbuf, w = this.w;

    var th = f * 16;
    var flat = (th < 0.5) || (th > 15.5);
    var flatCol = th > 15.5 ? c1 : c0;
    var depth = pts[0].z !== undefined;
    var self = this;

    this.spans(pts, function (y, x0, x1, iz0, iz1) {
      var xa = Math.max(Math.ceil(x0 - 0.5), self.clipX0);
      var xb = Math.min(Math.floor(x1 - 0.5), self.clipX1);
      if (xb < xa) return;
      var row = y * w, by = (y & 3) * 4;
      var dx = x1 - x0;
      var step = dx > 0.0001 ? (iz1 - iz0) / dx : 0;
      var iz = iz0 + (xa - x0) * step;
      if (flat) {
        for (var x = xa; x <= xb; x++, iz += step) {
          if (depth) {
            var idx = row + x;
            if (iz <= zb[idx]) continue;
            zb[idx] = iz;
            buf[idx] = flatCol;
          } else {
            buf[row + x] = flatCol;
          }
        }
        return;
      }
      for (var x2 = xa; x2 <= xb; x2++, iz += step) {
        var c = (BAYER[by + (x2 & 3)] < th) ? c1 : c0;
        if (depth) {
          var id2 = row + x2;
          if (iz <= zb[id2]) continue;
          zb[id2] = iz;
          buf[id2] = c;
        } else {
          buf[row + x2] = c;
        }
      }
    });
  };

  // dithered gradient polygon (for UI plates)
  FB.prototype.polyGrad = function (pts, cTop, cBot) {
    var buf = this.buf, w = this.w;
    var minY = Infinity, maxY = -Infinity;
    for (var i = 0; i < pts.length; i++) {
      if (pts[i].y < minY) minY = pts[i].y;
      if (pts[i].y > maxY) maxY = pts[i].y;
    }
    if (maxY - minY < 1) { this.poly(pts, cTop); return; }
    var self = this;
    this.spans(pts, function (y, x0, x1) {
      var xa = Math.max(Math.ceil(x0 - 0.5), self.clipX0);
      var xb = Math.min(Math.floor(x1 - 0.5), self.clipX1);
      if (xb < xa) return;
      var t = (y - minY) / (maxY - minY);
      var dith = (t * 8) % 1;
      var row = y * w, by = (y & 3) * 4;
      for (var x = xa; x <= xb; x++) {
        buf[row + x] = (BAYER[by + (x & 3)] < dith * 16) ? cBot : cTop;
      }
    });
  };

  FB.prototype.ellipseFill = function (cx, cy, rx, ry, color) {
    for (var y = -ry; y <= ry; y++) {
      var t = 1 - (y * y) / (ry * ry);
      if (t < 0) continue;
      var half = Math.sqrt(t) * rx;
      this.hline(Math.round(cx - half), Math.round(cx + half), Math.round(cy + y), color);
    }
  };

  FB.prototype.ellipse = function (cx, cy, rx, ry, color) {
    var steps = Math.max(24, Math.round((rx + ry) * 2));
    for (var i = 0; i < steps; i++) {
      var a = i / steps * Math.PI * 2;
      this.px(Math.round(cx + Math.cos(a) * rx), Math.round(cy + Math.sin(a) * ry), color);
    }
  };

  // dashed ellipse (scanner rings)
  FB.prototype.ellipseDash = function (cx, cy, rx, ry, color, dash) {
    var steps = Math.max(32, Math.round((rx + ry) * 2.4));
    for (var i = 0; i < steps; i++) {
      if ((i % (dash * 2)) >= dash) continue;
      var a = i / steps * Math.PI * 2;
      this.px(Math.round(cx + Math.cos(a) * rx), Math.round(cy + Math.sin(a) * ry), color);
    }
  };

  FB.prototype.text = function (str, x, y, color, scale) {
    scale = scale || 1;
    str = String(str).toUpperCase();
    var cx = x;
    for (var i = 0; i < str.length; i++) {
      var ch = str[i];
      if (ch === '\n') { y += 8 * scale; cx = x; continue; }
      var gl = GLYPH_BITS[ch] || GLYPH_BITS['?'];
      for (var row = 0; row < 7; row++) {
        var bits = gl[row];
        for (var col = 0; col < 5; col++) {
          if (!(bits & (1 << (4 - col)))) continue;
          if (scale === 1) this.px(cx + col, y + row, color);
          else this.fillRect(cx + col * scale, y + row * scale, cx + col * scale + scale - 1, y + row * scale + scale - 1, color);
        }
      }
      cx += 6 * scale;
    }
  };

  FB.prototype.textWidth = function (str, scale) {
    return String(str).length * 6 * (scale || 1) - (scale || 1);
  };

  FB.prototype.textC = function (str, cx, y, color, scale) {
    this.text(str, Math.round(cx - this.textWidth(str, scale) / 2), y, color, scale);
  };

  FB.prototype.textR = function (str, rx, y, color, scale) {
    this.text(str, Math.round(rx - this.textWidth(str, scale)), y, color, scale);
  };

  // nearest-neighbour blit with per-pixel lighting multiply (planet sprites)
  // nearest-neighbour scaled blit (planet sprites)
  FB.prototype.blitScaled = function (src, dx, dy, dw, dh, depth) {
    if (dw <= 0 || dh <= 0) return;
    var x0 = Math.round(dx), y0 = Math.round(dy);
    var w1 = Math.round(dw), h1 = Math.round(dh);
    var sw = src.w, sh = src.h;
    var sxStep = sw / w1, syStep = sh / h1;
    for (var y = 0; y < h1; y++) {
      var ty = y0 + y;
      if (ty < this.clipY0 || ty > this.clipY1) continue;
      var syi = (y * syStep) | 0;
      if (syi >= sh) syi = sh - 1;
      var row = syi * sw;
      for (var x = 0; x < w1; x++) {
        var tx = x0 + x;
        if (tx < this.clipX0 || tx > this.clipX1) continue;
        var sxi = (x * sxStep) | 0;
        if (sxi >= sw) sxi = sw - 1;
        var c = src.buf[row + sxi];
        if (((c >>> 24) & 255) === 0) continue;
        var di = ty * this.w + tx;
        this.buf[di] = c;
        if (depth) this.zbuf[di] = 1 / depth;
      }
    }
  };

  // copy a source rectangle (nearest neighbour, integer scale, x-wrap)
  FB.prototype.blitPart = function (src, sx, sy, dw, dh, dx, dy, scale) {
    scale = scale || 1;
    var W = src.w, H = src.h;
    for (var y = 0; y < dh; y++) {
      var ty = dy + y;
      if (ty < this.clipY0 || ty > this.clipY1) continue;
      var srow = (((sy + ((y / scale) | 0)) % H) + H) % H * W;
      for (var x = 0; x < dw; x++) {
        var tx = dx + x;
        if (tx < this.clipX0 || tx > this.clipX1) continue;
        var sxi = (((sx + ((x / scale) | 0)) % W) + W) % W;
        var c = src.buf[srow + sxi];
        if (((c >>> 24) & 255) === 0) continue;
        this.buf[ty * this.w + tx] = c;
      }
    }
  };

  FB.prototype.blit = function (src, dx, dy, scale, tint, alphaMask) {
    // tint: optional function(r,g,b,x,y,dist) -> [r,g,b] applied per pixel
    var w = this.w, h = this.h;
    var sw = src.w, sh = src.h;
    var x0 = Math.round(dx), y0 = Math.round(dy);
    for (var y = 0; y < sh * scale; y++) {
      var ty = y0 + y;
      if (ty < this.clipY0 || ty > this.clipY1) continue;
      var syi = (y / scale) | 0;
      for (var x = 0; x < sw * scale; x++) {
        var tx = x0 + x;
        if (tx < this.clipX0 || tx > this.clipX1) continue;
        var sxi = (x / scale) | 0;
        var c = src.buf[syi * sw + sxi];
        var a = (c >>> 24) & 255;
        if (a === 0) continue;
        if (tint) c = tint(c);
        this.buf[ty * w + tx] = c;
      }
    }
  };

  E.FB = FB;
  E.BAYER = BAYER;

  // --------------------------------------------------------------- 3D math
  var V = {
    add: function (a, b) { return [a[0] + b[0], a[1] + b[1], a[2] + b[2]]; },
    sub: function (a, b) { return [a[0] - b[0], a[1] - b[1], a[2] - b[2]]; },
    mul: function (a, s) { return [a[0] * s, a[1] * s, a[2] * s]; },
    dot: function (a, b) { return a[0] * b[0] + a[1] * b[1] + a[2] * b[2]; },
    cross: function (a, b) {
      return [a[1] * b[2] - a[2] * b[1], a[2] * b[0] - a[0] * b[2], a[0] * b[1] - a[1] * b[0]];
    },
    len: function (a) { return Math.sqrt(a[0] * a[0] + a[1] * a[1] + a[2] * a[2]); },
    norm: function (a) {
      var l = V.len(a); if (l < 1e-9) return [0, 0, -1];
      return [a[0] / l, a[1] / l, a[2] / l];
    }
  };
  E.V = V;

  // 3x3 matrix, row-major arrays of arrays
  var M3 = {
    id: function () { return [[1, 0, 0], [0, 1, 0], [0, 0, 1]]; },
    apply: function (m, v) {
      return [
        m[0][0] * v[0] + m[0][1] * v[1] + m[0][2] * v[2],
        m[1][0] * v[0] + m[1][1] * v[1] + m[1][2] * v[2],
        m[2][0] * v[0] + m[2][1] * v[1] + m[2][2] * v[2]
      ];
    },
    applyT: function (m, v) {
      return [
        m[0][0] * v[0] + m[1][0] * v[1] + m[2][0] * v[2],
        m[0][1] * v[0] + m[1][1] * v[1] + m[2][1] * v[2],
        m[0][2] * v[0] + m[1][2] * v[1] + m[2][2] * v[2]
      ];
    },
    mul: function (a, b) {
      var r = [[0, 0, 0], [0, 0, 0], [0, 0, 0]];
      for (var i = 0; i < 3; i++) for (var j = 0; j < 3; j++) {
        var s = 0;
        for (var k = 0; k < 3; k++) s += a[i][k] * b[k][j];
        r[i][j] = s;
      }
      return r;
    },
    rotX: function (a) {
      var c = Math.cos(a), s = Math.sin(a);
      return [[1, 0, 0], [0, c, -s], [0, s, c]];
    },
    rotY: function (a) {
      var c = Math.cos(a), s = Math.sin(a);
      return [[c, 0, s], [0, 1, 0], [-s, 0, c]];
    },
    rotZ: function (a) {
      var c = Math.cos(a), s = Math.sin(a);
      return [[c, -s, 0], [s, c, 0], [0, 0, 1]];
    },
    orthonormalize: function (m) {
      // columns are the axes
      var x = V.norm([m[0][0], m[1][0], m[2][0]]);
      var y = [m[0][1], m[1][1], m[2][1]];
      var z = V.norm(V.cross(x, y));
      y = V.cross(z, x);
      return [
        [x[0], y[0], z[0]],
        [x[1], y[1], z[1]],
        [x[2], y[2], z[2]]
      ];
    }
  };
  E.M3 = M3;

  // ----------------------------------------------------------- 3D renderer
  function Camera(w, h, fovDeg) {
    this.w = w; this.h = h;
    this.cx = w / 2; this.cy = h / 2;
    this.focal = (h / 2) / Math.tan((fovDeg || 70) * Math.PI / 360);
    this.near = 3;
    this.pos = [0, 0, 0];
    this.mat = M3.id();   // world -> view rotation (transpose applied to points)
  }
  Camera.prototype.toView = function (world) {
    var p = V.sub(world, this.pos);
    return M3.applyT(this.mat, p);
  };
  Camera.prototype.projectView = function (v) {
    var depth = -v[2];                  // camera looks down -Z
    if (depth < 0.001) depth = 0.001;
    return { x: this.cx + v[0] * this.focal / depth, y: this.cy - v[1] * this.focal / depth, z: depth };
  };
  Camera.prototype.project = function (world) {
    return this.projectView(this.toView(world));
  };
  E.Camera = Camera;

  // clip polygon (array of view-space [x,y,z]) keeping depth >= near (z <= -near)
  function clipNear(poly, near) {
    var lim = -near;
    var out = [], n = poly.length;
    for (var i = 0; i < n; i++) {
      var a = poly[i], b = poly[(i + 1) % n];
      var ain = a[2] <= lim, bin = b[2] <= lim;
      if (ain) out.push(a);
      if (ain !== bin) {
        var t = (lim - a[2]) / (b[2] - a[2]);
        out.push([a[0] + (b[0] - a[0]) * t, a[1] + (b[1] - a[1]) * t, lim]);
      }
    }
    return out;
  }

  // direction from surface towards the light (light sits above-left, behind viewer)
  var LIGHT = V.norm([-0.40, 0.50, 0.77]);

  // draw a mesh instance. ent: {mesh, pos, mat, scale?}
  // returns screen-space projected radius (for targeting etc.)
  function drawMesh(fb, cam, mesh, pos, mat, opts) {
    opts = opts || {};
    var scale = opts.scale || 1;
    var faces = mesh.faces;
    var list = [];
    var viewPos = cam.toView(pos);
    if (viewPos[2] > mesh.radius * scale + 100) return 0; // fully behind camera
    var bounds = { x0: 1e9, y0: 1e9, x1: -1e9, y1: -1e9 };

    for (var f = 0; f < faces.length; f++) {
      var face = faces[f];
      var pv = [], centre = [0, 0, 0], i;
      var culled = false;
      for (i = 0; i < face.p.length; i++) {
        var local = face.p[i];
        var wv = M3.apply(mat, [local[0] * scale, local[1] * scale, local[2] * scale]);
        wv = V.add(wv, pos);
        var vv = cam.toView(wv);
        pv.push(vv);
        centre[0] += vv[0]; centre[1] += vv[1]; centre[2] += vv[2];
      }
      centre[0] /= pv.length; centre[1] /= pv.length; centre[2] /= pv.length;

      // face normal in view space
      var nrm = null;
      if (face.p.length >= 3) {
        var worldN = M3.apply(mat, face.n);
        var viewN = M3.applyT(cam.mat, worldN);
        var toFace = V.norm(centre);
        // outward-facing test against the camera at the origin
        if (V.dot(viewN, toFace) > 0 && !face.double) { culled = true; }
        else nrm = V.norm(viewN);
      }
      if (culled) continue;

      if (pv.every(function (p) { return p[2] > -cam.near; })) continue;

      // shade
      var lam = nrm ? V.dot(nrm, LIGHT) : 0.5;
      var shade = 1.35 + Math.max(0, lam) * 3.85;
      if (face.shade !== undefined) shade = face.shade;
      if (face.tint) shade *= face.tint;
      // distance fade: slight darkening with range only
      var dist = V.len(centre);
      var fog = clamp(1 - dist / (opts.fogDist || 60000), 0, 1);
      shade = shade * (0.80 + 0.20 * fog);

      list.push({ pv: pv, z: centre[2], ramp: face.ramp, shade: shade, centre: centre });
    }

    list.sort(function (a, b) { return a.z - b.z; }); // far (more negative? see below)

    // view space: camera looks down -Z, so far = more negative z. Draw most negative first.
    list.sort(function (a, b) { return a.z - b.z; });

    var rampId = RAMP_ID;
    for (var j = 0; j < list.length; j++) {
      var it = list[j];
      var clipped = clipNear(it.pv, cam.near);
      if (clipped.length < 3) continue;
      var pts = clipped.map(function (p) { return cam.projectView(p); });
      var ok = true;
      for (var q = 0; q < pts.length; q++) {
        if (!isFinite(pts[q].x) || !isFinite(pts[q].y)) { ok = false; break; }
      }
      if (!ok) continue;
      var rid = typeof it.ramp === 'string' ? rampId[it.ramp] : it.ramp;
      fb.polyShaded(pts, rid, it.shade);
    }

    return bounds;
  }

  E.drawMesh = drawMesh;
  E.clipNear = clipNear;

  // sun / glow helper: concentric dithered discs
  E.drawGlow = function (fb, cx, cy, radius, ramp, maxShade) {
    var steps = 6;
    for (var i = steps; i >= 1; i--) {
      var r = radius * (i / steps);
      var shade = 0.6 + (maxShade - 0.6) * (1 - i / steps) * 0.5;
      var c = E.color(ramp, Math.min(5, shade + (steps - i) * 0.55));
      // dithered ring
      for (var y = -r; y <= r; y++) {
        var t = 1 - (y * y) / (r * r);
        if (t < 0) continue;
        var half = Math.sqrt(t) * r;
        var yy = Math.round(cy + y);
        for (var x = Math.round(cx - half); x <= Math.round(cx + half); x++) {
          var b = BAYER[((yy & 3) * 4 + ((x & 3)))];
          if (i > 1 && b > 5) continue;
          fb.px(x, yy, c);
        }
      }
    }
  };
})(typeof window !== 'undefined' ? window : globalThis);
