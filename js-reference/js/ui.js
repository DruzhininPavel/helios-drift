/* ELITE 16 - interface: dashboard, 3D view, station screens */
(function (global) {
  'use strict';
  var E = global.E;
  var V = E.V, M3 = E.M3, G = E.galaxy, GAME = E.game;
  var UI = E.ui = {};

  var COL = {
    text: E.rgb('#9fe8ff'),
    dim: E.rgb('#3f6f8c'),
    bg: E.rgb('#0a1220'),
    bg2: E.rgb('#111d2e'),
    frame: E.rgb('#2d5f7a'),
    frameHi: E.rgb('#79c9e6'),
    yellow: E.rgb('#ffd45e'),
    green: E.rgb('#6ef08a'),
    red: E.rgb('#ff6a4e'),
    white: E.rgb('#e8f6ff'),
    orange: E.rgb('#ff9c3a'),
    purple: E.rgb('#c48cff')
  };

  UI.sel = { market: 0, equip: 0, menu: 0, chart: null, chartCursor: null, name: 'JAMESON' };

  function pad(s, n) { s = String(s); while (s.length < n) s += ' '; return s; }
  function padl(s, n) { s = String(s); while (s.length < n) s = ' ' + s; return s; }
  function money(v) { return (Math.round(v * 10) / 10).toFixed(1); }

  // ------------------------------------------------------------ chrome
  function plate(fb, x, y, w, h, title) {
    fb.polyGrad([[x, y], [x + w, y], [x + w, y + h], [x, y + h]], COL.bg2, COL.bg);
    fb.rect(x, y, x + w - 1, y + h - 1, COL.frame);
    fb.hline(x + 1, x + w - 2, y + 1, E.rgb('#1b3348'));
    // corner ticks
    fb.hline(x, x + 5, y, COL.frameHi);
    fb.hline(x + w - 6, x + w - 1, y, COL.frameHi);
    if (title) {
      fb.fillRect(x + 4, y, x + 4 + fb.textWidth(title) + 6, y + 8, COL.bg);
      fb.text(title, x + 7, y - 3, COL.yellow);
    }
  }

  function segBar(fb, x, y, w, h, frac, ramp, warn) {
    frac = Math.max(0, Math.min(1, frac));
    var segs = 8, gap = 1;
    var sw = Math.floor((w - gap * (segs - 1)) / segs);
    var filled = frac * segs;
    for (var i = 0; i < segs; i++) {
      var sx = x + i * (sw + gap);
      var f = Math.max(0, Math.min(1, filled - i));
      var rid = E.RAMP_ID[ramp];
      if (f <= 0) {
        fb.fillRect(sx, y, sx + sw - 1, y + h - 1, E.rgb('#0d1a28'));
        fb.hline(sx, sx + sw - 1, y, E.rgb('#16283c'));
      } else {
        var shade = 1.6 + f * 2.4;
        var col = E.color(rid, shade);
        var col2 = E.color(rid, shade + 1);
        for (var yy = 0; yy < h; yy++) {
          for (var xx = 0; xx < sw; xx++) {
            var b = E.BAYER[((yy & 3) * 4) + ((xx & 3))];
            fb.px(sx + xx, y + yy, (f < 0.99 && b < 4) ? col2 : col);
          }
        }
      }
    }
    if (warn) {
      fb.rect(x - 1, y - 1, x + w, y + h, COL.red);
    }
  }

  // --------------------------------------------------------------- scanner
  function scanner(fb, cx, cy, rx, ry) {
    var p = GAME.player, w = GAME.world;
    var range = GAME.C.radarRange;
    fb.ellipseFill(cx, cy, rx, ry, E.rgb('#04121a'));
    // grid
    fb.ellipseDash(cx, cy, rx, ry, E.rgb('#1b5a6a'), 3);
    fb.ellipseDash(cx, cy, rx * 0.66, ry * 0.66, E.rgb('#123f4d'), 2);
    fb.ellipseDash(cx, cy, rx * 0.33, ry * 0.33, E.rgb('#123f4d'), 2);
    fb.hline(Math.round(cx - rx), Math.round(cx + rx), Math.round(cy), E.rgb('#123f4d'));
    fb.vline(Math.round(cx), Math.round(cy - ry), Math.round(cy + ry), E.rgb('#123f4d'));
    fb.ellipse(cx, cy, rx, ry, E.rgb('#3fbcd8'));

    function plot(pos, col, size, stick) {
      var rel = V.sub(pos, p.pos);
      var v = M3.applyT(p.mat, rel);
      var depth = -v[2];
      if (depth < -600) return;
      var side = v[0] / range;
      var fwdN = depth / range;
      if (Math.abs(side) > 1 || Math.abs(fwdN) > 1) return;
      var sx = Math.round(cx + side * rx);
      var sy = Math.round(cy - fwdN * ry);
      // altitude stick: v[1] up positive
      var alt = Math.max(-1, Math.min(1, v[1] / (range * 0.45)));
      var topY = Math.round(sy - alt * (ry * 0.55));
      if (stick) {
        fb.vline(sx, Math.min(sy, topY), Math.max(sy, topY), E.rgb('#2b6f80'));
      }
      if (size === 1) fb.px(sx, topY, col);
      else fb.fillRect(sx - 1, topY - 1, sx, topY, col);
    }

    // station
    plot(w.station.pos, E.color('teal', 5), 2, true);
    // planet / sun edges
    plot(w.planet.pos, E.color('blue', 3), 2, false);
    plot(w.sun.pos, E.color('yellow', 4), 2, false);

    var ents = w.entities;
    for (var i = 0; i < ents.length; i++) {
      var e = ents[i];
      if (e.kind === 'explosion') continue;
      var col;
      if (e.kind === 'missile') col = COL.white;
      else if (e.kind === 'canister') col = E.color('yellow', 4);
      else if (e.kind === 'asteroid') col = E.color('brown', 4);
      else if (e.hostile) col = COL.red;
      else if (e.behaviour === 'police') col = E.color('blue', 5);
      else col = E.color('green', 4);
      plot(e.pos, col, 1, true);
    }
  }

  function compass(fb, cx, cy, r) {
    var p = GAME.player, st = GAME.world.station;
    fb.ellipseFill(cx, cy, r, r, E.rgb('#04121a'));
    fb.ellipse(cx, cy, r, r, E.rgb('#3fbcd8'));
    var rel = V.sub(st.pos, p.pos);
    var v = M3.applyT(p.mat, rel);
    var depth = -v[2];
    var len = Math.sqrt(v[0] * v[0] + v[1] * v[1] + depth * depth) || 1;
    var nx = v[0] / len, ny = -v[1] / len;
    var sx = Math.round(cx + nx * (r - 2));
    var sy = Math.round(cy + ny * (r - 2));
    var col = depth > 0 ? COL.green : E.color('teal', 2);
    if (depth > 0) fb.fillRect(sx - 1, sy - 1, sx + 1, sy + 1, col);
    else fb.px(sx, sy, col);
    fb.px(cx, cy, E.rgb('#1b5a6a'));
  }

  // ---------------------------------------------------------------- flight
  UI.flight = function (fb, dt) {
    var p = GAME.player, w = GAME.world;
    if (!w) return;
    var shakeX = 0, shakeY = 0;
    if (GAME.shake > 0) {
      shakeX = Math.round((Math.random() - 0.5) * GAME.shake * 6);
      shakeY = Math.round((Math.random() - 0.5) * GAME.shake * 6);
      GAME.shake -= dt * 3;
    }

    fb.setClip(0, 0, 319, 149);
    // deep space gradient
    fb.clearGradient([E.rgb('#05070f'), E.rgb('#070c16'), E.rgb('#08101c'),
      E.rgb('#0a1422'), E.rgb('#0b1726'), E.rgb('#0c1a2b')], E.rgb('#0c1a2b'));

    // nebula backdrop (pans with yaw)
    if (!w.nebula) w.nebula = E.models.makeNebula(w.sys.seed, GAME.cmd.galaxy);
    var fwdW = [-p.mat[0][2], -p.mat[1][2], -p.mat[2][2]];
    var yawW = Math.atan2(fwdW[0], -fwdW[2]);
    var pitchW = Math.asin(Math.max(-1, Math.min(1, fwdW[1])));
    var nebOff = Math.round((yawW / (Math.PI * 2)) * w.nebula.w * 4);
    var nebY = Math.round(pitchW * 9);
    if (nebY < -4) nebY = -4; if (nebY > 4) nebY = 4;
    fb.blitPart(w.nebula, nebOff, 0, 320, 150, 0, -nebY, 2);

    var cam = new E.Camera(GAME.C.viewW, GAME.C.viewH, 72);
    cam.cx = 160 + shakeX; cam.cy = 75 + shakeY;
    cam.pos = p.pos;
    cam.mat = p.mat;
    var focal = cam.focal;

    // ---- stars
    for (var i = 0; i < w.stars.length; i++) {
      var s = w.stars[i];
      var world = V.add(p.pos, V.mul(s.dir, 400000));
      var vp = cam.toView(world);
      if (vp[2] > -1000) continue;
      var sp = cam.projectView(vp);
      if (sp.x < -2 || sp.x > 322 || sp.y < -2 || sp.y > 152) continue;
      var bright = s.bright;
      var col = bright > 0.86 ? E.color('steel', 5) : (bright > 0.6 ? E.color('steel', 4) : E.color('steel', 3));
      fb.px(sp.x, sp.y, col);
      if (s.size === 2 && bright > 0.8) {
        fb.px(sp.x + 1, sp.y, E.color('steel', 3));
        fb.px(sp.x, sp.y + 1, E.color('steel', 3));
      }
    }

    // ---- distant planet
    var planet = w.planet;
    var pv = cam.toView(planet.pos);
    var pdepth = -pv[2];
    if (pdepth > planet.radius * 0.5) {
      var psp = cam.projectView(pv);
      var sr = planet.radius * focal / pdepth;
      if (!w.planetFrames) {
        w.planetFrames = E.models.makePlanetFrames(128, planet.seed, w.sys.kind, 8);
        w.planetFrame = 0;
      }
      var frameIdx = Math.floor((GAME.time * 0.6) % 8);
      if (sr > 3) {
        var frames = w.planetFrames;
        fb.blitScaled(frames[frameIdx], psp.x - sr, psp.y - sr, sr * 2, sr * 2,
          Math.max(1, pdepth - planet.radius * 0.6));
        if (sr > 7) {
          // atmosphere halo, tinted by world type
          var rimRamp = ['teal', 'orange', 'steel', 'yellow', 'red'][w.sys.kind] || 'teal';
          for (var ri = 0; ri < 3; ri++) {
            var rad = sr + 1.2 + ri * 1.6;
            var colr = E.color(rimRamp, 3 - ri);
            var steps = Math.max(28, Math.round(rad * 3.4));
            var thr = 6 + ri * 4;
            for (var ai = 0; ai < steps; ai++) {
              var ang = ai / steps * Math.PI * 2;
              var qx = Math.round(psp.x + Math.cos(ang) * rad);
              var qy = Math.round(psp.y + Math.sin(ang) * rad);
              if (qx < 0 || qx > 319 || qy < 0 || qy > 149) continue;
              if (E.BAYER[((qy & 3) * 4 + (qx & 3))] > thr) fb.px(qx, qy, colr);
            }
          }
        }
      } else {
        fb.fillRect(psp.x - 1, psp.y - 1, psp.x, psp.y, E.color('blue', 4));
      }
    }

    // ---- sun
    var sv = cam.toView(w.sun.pos);
    var sdepth = -sv[2];
    if (sdepth > w.sun.radius) {
      var ssp = cam.projectView(sv);
      var sunR = Math.min(160, w.sun.radius * focal / sdepth);
      if (ssp.x > -sunR && ssp.x < 320 + sunR && ssp.y > -sunR && ssp.y < 150 + sunR) {
        E.drawGlow(fb, ssp.x, ssp.y, sunR, 'yellow', 5);
        fb.ellipseFill(ssp.x, ssp.y, Math.max(1, sunR * 0.45), Math.max(1, sunR * 0.45), E.color('yellow', 5));
      }
    }

    // ---- dust (motion cue)
    for (var d = 0; d < w.dust.length; d++) {
      var dp = w.dust[d];
      var rel = V.sub(dp.pos, p.pos);
      if (V.len(rel) > 2600) {
        var dir = V.norm([Math.random() - 0.5, Math.random() - 0.5, -1]);
        dp.pos = V.add(p.pos, V.mul(dir, 400 + Math.random() * 2200));
        continue;
      }
      var dv = cam.toView(dp.pos);
      if (dv[2] > -40) continue;
      var dsp = cam.projectView(dv);
      if (dsp.x < 0 || dsp.x > 319 || dsp.y < 0 || dsp.y > 149) continue;
      var frac = 1 - Math.min(1, (-dv[2]) / 2600);
      var col2 = E.color('steel', 2 + Math.floor(frac * 3));
      fb.px(dsp.x, dsp.y, col2);
      if (p.speed > 250 && frac > 0.45) fb.px(dsp.x, dsp.y + 1, E.color('steel', 2));
    }

    // ---- entities sorted far to near
    var list = w.entities.slice();
    list.push(w.station);
    list.sort(function (a, b) {
      return V.len(V.sub(b.pos, p.pos)) - V.len(V.sub(a.pos, p.pos));
    });

    for (var k = 0; k < list.length; k++) {
      var e = list[k];
      if (e.kind === 'explosion') {
        var ev = cam.toView(e.pos);
        if (ev[2] > 0) continue;
        var esp = cam.projectView(ev);
        var t = e.age / e.life;
        var rad = (e.radius * (0.4 + t * 2.2)) * focal / Math.max(20, -ev[2]);
        if (rad < 1) rad = 1;
        var rings = 3;
        for (var rr = 0; rr < rings; rr++) {
          var fr = (t + rr * 0.18) % 1;
          var r2 = rad * (0.4 + fr);
          var ramp = fr < 0.4 ? 'yellow' : (fr < 0.7 ? 'orange' : 'red');
          var col = E.color(ramp, 5 - Math.floor(fr * 3));
          var steps = Math.max(10, Math.round(r2 * 3));
          for (var a = 0; a < steps; a++) {
            var ang = a / steps * Math.PI * 2;
            var px = Math.round(esp.x + Math.cos(ang) * r2);
            var py = Math.round(esp.y + Math.sin(ang) * r2);
            if (E.BAYER[((py & 3) * 4 + (px & 3))] < 6) fb.px(px, py, col);
          }
        }
        continue;
      }
      if (e.kind === 'station') {
        E.drawMesh(fb, cam, E.models.meshes.station, e.pos, e.mat, { fogDist: 60000 });
        // slot beacon
        var nrm = [e.mat[0][2], e.mat[1][2], e.mat[2][2]];
        var beacon = V.add(e.pos, V.mul(nrm, GAME.C.stationS + 2));
        var bv = cam.toView(beacon);
        if (bv[2] < 0) {
          var bsp = cam.projectView(bv);
          var blink = (Math.sin(GAME.time * 4) > 0);
          var bc = blink ? COL.red : E.color('red', 2);
          fb.fillRect(bsp.x - 1, bsp.y - 1, bsp.x + 1, bsp.y + 1, bc);
        }
        continue;
      }
      var dist = V.len(V.sub(e.pos, p.pos));
      if (dist > 40000) continue;
      if (e.mesh) {
        E.drawMesh(fb, cam, e.mesh, e.pos, e.mat, { fogDist: 30000 });
        // engine flare at the tail
        if (e.speed > 25 && dist < 18000) {
          var back = [e.mat[0][2], e.mat[1][2], e.mat[2][2]];
          var tail = V.add(e.pos, V.mul(back, (e.radius || 20) * 0.85));
          var tv = cam.toView(tail);
          if (tv[2] < -20) {
            var tp = cam.projectView(tv);
            var glowR = Math.max(1, Math.min(7, (e.radius || 20) * 0.28 * focal / (-tv[2])));
            UI.glowDot(fb, tp.x, tp.y, glowR, e.hostile ? 'orange' : 'teal');
          }
        }
      }
    }

    // ---- target reticle
    UI.reticle(fb, cam);

    // ---- laser hit flash / beams
    if (GAME.beamFlash > 0) {
      GAME.beamFlash -= dt;
      var cx = 160, cy = 75;
      var flick = Math.random() * 0.6 + 0.4;
      fb.line(90, 149, cx - 3, cy + 4, E.color('red', Math.round(4 * flick) + 1));
      fb.line(230, 149, cx + 3, cy + 4, E.color('red', Math.round(4 * flick) + 1));
      fb.px(cx - 4, cy, COL.white); fb.px(cx + 4, cy, COL.white);
    }
    if (GAME.ecmFlash > 0) {
      for (var y = 0; y < 150; y += 2) {
        var f = (Math.sin(y * 0.35 + GAME.time * 40) * 0.5 + 0.5);
        if (f > 0.6) fb.hline(0, 319, y, E.color('teal', 2 + Math.round(f * 2)));
      }
    }
    if (GAME.bombFlash > 0) {
      var ba = Math.min(1, GAME.bombFlash);
      for (var yy = 0; yy < 150; yy++) {
        var f2 = 1 - Math.abs(yy - 75) / 75;
        if (Math.random() < f2 * ba) fb.hline(0, 319, yy, E.color('yellow', 4 + Math.round(f2 * 1)));
      }
    }

    fb.resetClip();
    UI.cockpit(fb);

    // ---- crosshair
    UI.crosshair(fb, 160, 75);

    // ---- messages
    var my = 6;
    for (var m = GAME.messages.length - 1; m >= 0; m--) {
      var msg = GAME.messages[m];
      var colm = msg.kind === 'bad' ? COL.red : (msg.kind === 'good' ? COL.green : COL.yellow);
      var alpha = Math.min(1, (msg.life - msg.t) * 2);
      if (alpha < 0.4 && Math.sin(msg.t * 30) < 0) continue;
      var tw = fb.textWidth(msg.text);
      fb.fillRect(160 - tw / 2 - 3, my - 1, 160 + tw / 2 + 2, my + 7, E.rgb('#050b12'));
      fb.textC(msg.text, 160, my, colm);
      my += 9;
    }

    UI.hud(fb, dt);
  };

  UI.glowDot = function (fb, x, y, r, ramp) {
    var rBase = E.RAMP_ID[ramp] || E.RAMP_ID.orange;
    for (var yy = -r; yy <= r; yy++) {
      for (var xx = -r; xx <= r; xx++) {
        var d = Math.sqrt(xx * xx + yy * yy);
        if (d > r + 0.4) continue;
        var b = E.BAYER[(((y + yy) & 3) * 4) + ((x + xx) & 3)];
        if (d > r * 0.55 && b > 9) continue;
        var shade = d < r * 0.45 ? 5 : (d < r * 0.8 ? 4 : 3);
        fb.px(x + xx, y + yy, E.color(rBase, shade));
      }
    }
  };

  UI.reticle = function (fb, cam) {
    var p = GAME.player, w = GAME.world;
    if (!w) return;
    var best = null, bestScore = 1e9, bestScreen = null;
    for (var i = 0; i < w.entities.length; i++) {
      var e = w.entities[i];
      if (e.kind !== 'ship' || !e.hostile) continue;
      var vv = cam.toView(e.pos);
      if (vv[2] > -60) continue;
      var s = cam.projectView(vv);
      if (s.x < 0 || s.x > 319 || s.y < 0 || s.y > 149) continue;
      var score = Math.sqrt(vv[0] * vv[0] + vv[1] * vv[1] + vv[2] * vv[2]);
      if (score < bestScore) { bestScore = score; best = e; bestScreen = s; }
    }
    if (!best) return;
    var e = best, s = bestScreen;
    var depth = -cam.toView(e.pos)[2];
    var rad = Math.max(5, Math.min(40, (e.radius || 24) * cam.focal / depth));
    var x0 = Math.round(s.x - rad), x1 = Math.round(s.x + rad);
    var y0 = Math.round(s.y - rad * 0.8), y1 = Math.round(s.y + rad * 0.8);
    var col = E.color('red', 4);
    var len = Math.max(3, Math.round(rad * 0.35));
    // corner brackets
    fb.hline(x0, x0 + len, y0, col); fb.vline(x0, y0, y0 + len, col);
    fb.hline(x1 - len, x1, y0, col); fb.vline(x1, y0, y0 + len, col);
    fb.hline(x0, x0 + len, y1, col); fb.vline(x0, y1 - len, y1, col);
    fb.hline(x1 - len, x1, y1, col); fb.vline(x1, y1 - len, y1, col);
    var label = (e.type || 'SHIP').toUpperCase() + ' ' + Math.round(depth);
    fb.text(label, Math.round(s.x) - fb.textWidth(label) / 2, Math.round(y1) + 3, E.color('red', 3));
  };

  UI.cockpit = function (fb) {
    var mid = E.rgb('#0a1622'), dark = E.rgb('#040a11'), hi = E.color('steel', 3), low = E.color('steel', 1);
    var legs = [[0, 0, 1, 1], [319, 0, -1, 1], [0, 149, 1, -1], [319, 149, -1, -1]];
    for (var i = 0; i < legs.length; i++) {
      var x = legs[i][0], y = legs[i][1], sx = legs[i][2], sy = legs[i][3];
      var w = 54, h = 26;
      fb.poly([[x, y], [x + sx * w, y], [x + sx * w, y + sy * 8], [x + sx * 16, y + sy * h]], mid);
      fb.line(x + sx * w, y, x + sx * 16, y + sy * h, hi);
      fb.line(x + sx * 16, y + sy * h, x, y + sy * h, hi);
      fb.line(x + sx * w, y + sy * 8, x + sx * 16, y + sy * h, low);
      for (var r = 1; r <= 3; r++) {
        fb.px(x + sx * (10 + r * 13), y + sy * (3 + r * 4), E.color('steel', 4));
      }
    }
    // edge darkening
    for (var yy = 0; yy < 150; yy++) {
      for (var xx = 0; xx < 6; xx++) {
        var b1 = E.BAYER[((yy & 3) * 4 + (xx & 3))];
        if (b1 < (6 - xx) * 3) {
          fb.px(xx, yy, dark);
          fb.px(319 - xx, yy, dark);
        }
      }
    }
    for (var y2 = 0; y2 < 5; y2++) {
      for (var x2 = 0; x2 < 320; x2++) {
        var b2 = E.BAYER[((y2 & 3) * 4 + (x2 & 3))];
        if (b2 < (5 - y2) * 3) {
          fb.px(x2, y2, dark);
          fb.px(x2, 149 - y2, dark);
        }
      }
    }
  };

  UI.crosshair = function (fb, cx, cy) {
    var c = E.color('green', 4);
    var c2 = E.color('green', 3);
    fb.px(cx, cy, COL.white);
    // brackets
    fb.hline(cx - 12, cx - 6, cy, c);
    fb.hline(cx + 6, cx + 12, cy, c);
    fb.vline(cx, cy - 9, cy - 4, c2);
    fb.vline(cx, cy + 4, cy + 9, c2);
    fb.px(cx - 14, cy, c2); fb.px(cx + 14, cy, c2);
    fb.px(cx, cy - 11, c2); fb.px(cx, cy + 11, c2);
  };

  // ------------------------------------------------------------------ HUD
  UI.hud = function (fb, dt) {
    var p = GAME.player, cmd = GAME.cmd, C = GAME.C;
    var y0 = C.viewH;
    fb.fillRect(0, y0, 319, 223, COL.bg);
    fb.polyGrad([[0, y0 + 1], [320, y0 + 1], [320, 224], [0, 224]], E.rgb('#0d1a2a'), E.rgb('#070d16'));
    fb.hline(0, 319, y0, COL.frameHi);
    fb.hline(0, 319, y0 + 1, E.rgb('#1b3348'));

    // ---- left gauges
    var gx = 6, bx = 22, bw = 58;
    fb.text('FS', gx, y0 + 8, COL.text);
    segBar(fb, bx, y0 + 8, bw, 5, p.shieldF / C.maxShield, 'blue', p.shieldF < 22);
    fb.text('AS', gx, y0 + 17, COL.text);
    segBar(fb, bx, y0 + 17, bw, 5, p.shieldA / C.maxShield, 'blue', p.shieldA < 22);
    fb.text('FU', gx, y0 + 26, COL.text);
    segBar(fb, bx, y0 + 26, bw, 5, cmd.fuel / C.fuelMax, 'green', cmd.fuel < 0.8);
    fb.text('EN', gx, y0 + 35, COL.text);
    segBar(fb, bx, y0 + 35, bw, 5, p.energy / C.maxEnergy, 'yellow', p.energy < 22);
    fb.text('HT', gx, y0 + 44, COL.text);
    segBar(fb, bx, y0 + 44, bw, 5, p.laserTemp, 'red', p.laserTemp > 0.85);

    fb.text('SPD ' + Math.round(p.speed), gx, y0 + 54, COL.dim);
    fb.text(cmd.fuel.toFixed(1) + 'LY', bx + bw - 26, y0 + 54, COL.green);
    fb.text('SYS ' + GAME.world.sys.name, gx, y0 + 63, E.rgb('#2d5f7a'));

    // ---- scanner
    scanner(fb, 160, y0 + 40, 56, 26);

    // ---- right block
    var rx = 220;      // label column
    var vx = 238;      // value column
    fb.text('MSL', rx, y0 + 8, COL.text);
    for (var i = 0; i < C.missilesMax; i++) {
      var mx = vx + i * 10;
      if (i < cmd.missiles) {
        fb.fillRect(mx, y0 + 8, mx + 5, y0 + 13, E.color('ash', 4));
        fb.vline(mx + 2, y0 + 9, y0 + 12, COL.red);
        fb.hline(mx, mx + 5, y0 + 8, E.color('ash', 5));
      } else {
        fb.rect(mx, y0 + 8, mx + 5, y0 + 13, E.rgb('#22415a'));
      }
    }
    fb.text('LSR', rx, y0 + 17, COL.text);
    segBar(fb, vx, y0 + 17, 34, 5, p.laserTemp, 'red');
    fb.text(['PULSE', 'BEAM', 'MILITARY', 'MINING'][cmd.laser], vx + 38, y0 + 17, E.rgb('#2d5f7a'));

    // status lamps
    var lamps = [
      ['E', cmd.equipment.ecm],
      ['S', cmd.equipment.fuelscoop],
      ['D', cmd.equipment.dockcomp],
      ['B', cmd.bombs > 0]
    ];
    for (var l = 0; l < lamps.length; l++) {
      var lx = rx + l * 13;
      var on = lamps[l][1];
      fb.text(lamps[l][0], lx, y0 + 26, on ? COL.green : E.rgb('#274458'));
      if (on) fb.px(lx + 5, y0 + 26, E.color('green', 5));
    }

    fb.textR(money(cmd.cash) + ' CR', 318, y0 + 35, COL.yellow);
    var legal = cmd.legal;
    var legalTxt = legal >= 500 ? 'FUGITIVE' : (legal > 0 ? 'OFFENDER' : 'CLEAN');
    fb.textR(legalTxt, 318, y0 + 44, legal > 0 ? COL.red : COL.green);
    fb.textR('HOLD ' + GAME.cargoUsed(cmd) + '/' + cmd.cargoCap, 318, y0 + 53, COL.dim);

    compass(fb, 230, y0 + 55, 8);
    var dist = V.len(V.sub(GAME.world.station.pos, p.pos));
    fb.text('STN ' + (dist > 9999 ? (dist / 1000).toFixed(1) + 'K' : Math.round(dist)), 246, y0 + 62, COL.dim);
  };

  // ------------------------------------------------------------------ title
  UI.title = function (fb, dt) {
    fb.clearGradient([E.rgb('#04060d'), E.rgb('#070d18'), E.rgb('#0a1422'),
      E.rgb('#0c1a2b'), E.rgb('#0a1422'), E.rgb('#05070f')], E.rgb('#05070f'));
    // stars
    for (var i = 0; i < 70; i++) {
      var sx = (i * 73 + (GAME.time * 12 * (0.3 + (i % 5) * 0.2))) % 320;
      var sy = (i * 41) % 224;
      fb.px(sx, sy, E.color('steel', 2 + (i % 3)));
    }
    // rotating cobra
    var cam = new E.Camera(320, 224, 55);
    cam.cx = 160; cam.cy = 128;
    var mat = M3.mul(M3.rotX(0.35 + Math.sin(GAME.time * 0.3) * 0.15), M3.rotY(GAME.time * 0.5));
    E.drawMesh(fb, cam, E.models.meshes.cobra3, [0, 6, -230], mat, {});
    // logo
    fb.textC('HELIOS', 160, 14, E.color('steel', 5), 5);
    fb.textC('DRIFT', 160, 53, E.color('steel', 4), 3);
    fb.textC('SIXTEEN BIT EDITION', 160, 68, COL.yellow, 1);
    fb.textC('THE GALAXY IS YOURS, COMMANDER', 160, 80, COL.dim, 1);
    // menu
    fb.textC('1  NEW GAME', 160, 168, COL.text, 1);
    fb.textC('2  LOAD COMMANDER', 160, 180, GAME.hasSave() ? COL.text : COL.dim, 1);
    fb.textC('ARROWS / MOUSE STEER   SPACE FIRE   X/Z SPEED', 160, 198, COL.dim, 1);
    fb.textC('C CHART   J JUMP   ESC CONTROLS', 160, 209, E.rgb('#2d5f7a'), 1);
  };

  // -------------------------------------------------------------- game over
  UI.gameOver = function (fb) {
    fb.clear(E.rgb('#05070f'));
    fb.textC('GAME OVER', 160, 60, COL.red, 3);
    fb.textC(GAME.dead ? GAME.dead.reason : 'SHIP LOST', 160, 100, COL.text, 1);
    fb.textC('CASH ' + money(GAME.cmd.cash) + ' CR   KILLS ' + GAME.cmd.kills, 160, 120, COL.dim, 1);
    fb.textC('PRESS ENTER TO CONTINUE', 160, 160, COL.yellow, 1);
  };

  // --------------------------------------------------------------- hyper
  UI.hyper = function (fb, dt) {
    fb.clear(E.rgb('#02030a'));
    var t = GAME.hyper ? GAME.hyper.t : 1;
    var cx = 160, cy = 112;
    var n = 70;
    for (var i = 0; i < n; i++) {
      var a = (i / n) * Math.PI * 2 + t * 0.4;
      var r0 = 8 + ((i * 37 + t * 900) % 260);
      var r1 = r0 + 10 + t * 40;
      var col = E.color('steel', 2 + (i % 3));
      var ca = Math.cos(a), sa = Math.sin(a);
      fb.line(cx + ca * r0, cy + sa * r0 * 0.75, cx + ca * r1, cy + sa * r1 * 0.75, col);
    }
    fb.textC('HYPERSPACE', 160, 20, COL.white, 2);
    var gal = G.get(GAME.cmd.galaxy);
    var tgt = gal.systems[GAME.hyper ? GAME.hyper.target : GAME.cmd.targetSystem];
    fb.textC('TO ' + tgt.name, 160, 40, COL.yellow, 1);
    fb.textC(GAME.hyper ? (GAME.hyper.dist.toFixed(1) + ' LIGHT YEARS') : '', 160, 200, COL.dim, 1);
  };

  // ------------------------------------------------------- docked screens
  UI.chrome = function (fb, title) {
    fb.clearGradient([E.rgb('#04070e'), E.rgb('#081120'), E.rgb('#0b1a2c'),
      E.rgb('#0a1626'), E.rgb('#071019'), E.rgb('#04070e')], E.rgb('#04070e'));
    fb.fillRect(0, 0, 319, 13, E.rgb('#08131f'));
    fb.hline(0, 319, 13, COL.frame);
    fb.text(title, 6, 4, COL.yellow);
    fb.textR('GALAXY ' + (GAME.cmd.galaxy + 1) + '   ' + money(GAME.cmd.cash) + ' CR', 314, 4, COL.text);
  };

  UI.footer = function (fb) {
    fb.fillRect(0, 211, 319, 223, E.rgb('#08131f'));
    fb.hline(0, 319, 211, COL.frame);
    fb.textC(UI.hint || '', 160, 215, COL.dim);
  };

  UI.docked = function (fb, dt) {
    var sys = GAME.world ? GAME.world.sys : null;
    UI.chrome(fb, sys ? sys.name + ' STATION' : 'STATION');
    switch (GAME.screen) {
      case 'market': UI.market(fb); break;
      case 'equip': UI.equip(fb); break;
      case 'chart': UI.chart(fb); break;
      case 'status': UI.status(fb); break;
      default: UI.stationMenu(fb); break;
    }
    UI.footer(fb);
  };

  UI.pause = function (fb) {
    UI.chrome(fb, 'PAUSED');
    plate(fb, 40, 26, 240, 170, 'FLIGHT CONTROLS');
    var rows = [
      ['ARROWS / MOUSE', 'PITCH AND ROLL'],
      ['X / Z', 'ACCELERATE / BRAKE'],
      ['SPACE or LMB', 'FIRE LASERS'],
      ['M', 'LAUNCH MISSILE'],
      ['E', 'E.C.M.'],
      ['TAB', 'ENERGY BOMB'],
      ['L', 'DOCKING COMPUTER'],
      ['J', 'IN-SYSTEM JUMP'],
      ['C', 'GALACTIC CHART'],
      ['V', 'COMMANDER STATUS'],
      ['ESC', 'RESUME']
    ];
    for (var i = 0; i < rows.length; i++) {
      fb.text(rows[i][0], 52, 40 + i * 12, COL.yellow);
      fb.textR(rows[i][1], 268, 40 + i * 12, COL.text);
    }
  };

  UI.stationMenu = function (fb) {
    UI.hint = '1-6 SELECT   UP/DOWN   F REFUEL   ENTER CONFIRM';
    plate(fb, 10, 20, 150, 108, 'COMMANDER ' + GAME.cmd.name);
    fb.text('CASH', 18, 34, COL.dim); fb.textR(money(GAME.cmd.cash) + ' CR', 152, 34, COL.yellow);
    fb.text('FUEL', 18, 44, COL.dim); fb.textR(GAME.cmd.fuel.toFixed(1) + ' LY', 152, 44, COL.text);
    fb.text('HOLD', 18, 54, COL.dim); fb.textR(GAME.cargoUsed(GAME.cmd) + ' / ' + GAME.cmd.cargoCap + ' T', 152, 54, COL.text);
    fb.text('MISSILES', 18, 64, COL.dim); fb.textR(String(GAME.cmd.missiles), 152, 64, COL.text);
    fb.text('LASER', 18, 74, COL.dim); fb.textR(GAME.C.laserName[GAME.cmd.laser], 152, 74, COL.text);
    fb.text('RANK', 18, 84, COL.dim); fb.textR(G.rankFor(GAME.cmd.kills), 152, 84, COL.green);
    fb.text('LEGAL', 18, 94, COL.dim);
    fb.textR(GAME.cmd.legal >= 500 ? 'FUGITIVE' : (GAME.cmd.legal > 0 ? 'OFFENDER' : 'CLEAN'),
      152, 94, GAME.cmd.legal > 0 ? COL.red : COL.green);
    fb.text('TARGET', 18, 104, COL.dim); fb.textR(G.get(GAME.cmd.galaxy).systems[GAME.cmd.targetSystem].name, 152, 104, COL.text);

    plate(fb, 168, 20, 142, 108, 'STATION SERVICES');
    var items = ['BUY / SELL CARGO', 'EQUIP SHIP', 'GALACTIC CHART', 'COMMANDER STATUS',
      'LAUNCH', 'SAVE COMMANDER'];
    for (var i = 0; i < items.length; i++) {
      var y = 36 + i * 14;
      var sel = UI.sel.menu === i;
      if (sel) fb.fillRect(174, y - 3, 304, y + 7, E.rgb('#12384f'));
      fb.text((i + 1) + '  ' + items[i], 176, y, sel ? COL.yellow : COL.text);
    }
    fb.text('REFUEL  [F]  ' + G.fuelPricePerLY + ' CR/LY', 176, 124, COL.dim);
  };

  UI.market = function (fb) {
    UI.hint = 'UP/DOWN SELECT   B BUY 1T   SHIFT+B 10T   S SELL   ESC BACK';
    var cmd = GAME.cmd;
    plate(fb, 4, 18, 312, 190, 'GALACTIC MARKET');
    fb.text('COMMODITY', 12, 30, COL.dim);
    fb.textR('UNIT', 176, 30, COL.dim);
    fb.textR('PRICE', 218, 30, COL.dim);
    fb.textR('QTY', 252, 30, COL.dim);
    fb.textR('HOLD', 292, 30, COL.dim);
    fb.hline(10, 310, 37, COL.frame);

    var list = GAME.market;
    for (var i = 0; i < list.length; i++) {
      var m = list[i];
      var y = 42 + i * 9;
      var sel = UI.sel.market === i;
      if (sel) fb.fillRect(8, y - 1, 311, y + 7, E.rgb('#12384f'));
      var have = cmd.cargo[i] || 0;
      var col = sel ? COL.yellow : COL.text;
      fb.text(m.name, 12, y, col);
      fb.textR(m.unit, 176, y, COL.dim);
      fb.textR(String(m.price), 218, y, sel ? COL.white : COL.green);
      fb.textR(String(m.qty), 252, y, COL.dim);
      fb.textR(String(have), 292, y, have > 0 ? COL.yellow : COL.dim);
      // profit hint
      if (have > 0) {
        var diff = m.sell - m.price;
      }
    }
    fb.text('HOLD ' + GAME.cargoUsed(cmd) + '/' + cmd.cargoCap + ' T', 240, 197, COL.yellow);
  };

  UI.equip = function (fb) {
    UI.hint = 'UP/DOWN SELECT   ENTER BUY   ESC BACK';
    var cmd = GAME.cmd;
    plate(fb, 4, 18, 312, 190, 'EQUIP SHIP  -  TECH LEVEL ' + (GAME.world ? GAME.world.sys.tech : 0));
    var items = G.availableEquipment(GAME.world.sys);
    // always show missiles
    for (var i = 0; i < items.length; i++) {
      var it = items[i];
      var y = 32 + i * 12;
      var sel = UI.sel.equip === i;
      if (sel) fb.fillRect(8, y - 2, 311, y + 8, E.rgb('#12384f'));
      var owned = it.id === 'missile' ? (cmd.missiles >= GAME.C.missilesMax)
        : (it.id.indexOf('laser') === 0 ? false : !!cmd.equipment[it.id]);
      fb.text(it.name, 14, y, owned ? COL.dim : (sel ? COL.yellow : COL.text));
      fb.textR(String(it.price) + ' CR', 240, y, cmd.cash >= it.price ? COL.green : COL.red);
      if (it.id === 'missile') fb.textR('x' + cmd.missiles, 290, y, COL.text);
      else if (owned) fb.textR('FITTED', 290, y, COL.green);
    }

  };

  UI.chart = function (fb) {
    UI.hint = 'ARROWS SELECT SYSTEM   H HYPERSPACE   ESC BACK';
    var cmd = GAME.cmd;
    var gal = G.get(cmd.galaxy);
    var cur = gal.systems[cmd.system];
    var sel = UI.sel.chartCursor === null ? cur : gal.systems[UI.sel.chartCursor];
    if (!sel) sel = cur;

    plate(fb, 4, 18, 196, 190, 'GALACTIC CHART ' + (cmd.galaxy + 1));
    var ox = 14, oy = 34, cw = 176, ch = 164;
    fb.rect(ox, oy, ox + cw, oy + ch, E.rgb('#16394a'));
    // range circle (7 LY = 28 chart units)
    var cxp = ox + cur.x / 256 * cw, cyp = oy + ch - cur.y / 256 * ch;
    var rr = 28 / 256 * cw;
    fb.ellipse(cxp, cyp, rr, rr, E.rgb('#1f5e6e'));

    for (var i = 0; i < gal.systems.length; i++) {
      var s = gal.systems[i];
      var x = Math.round(ox + s.x / 256 * cw);
      var y = Math.round(oy + ch - s.y / 256 * ch);
      var isCur = s.id === cur.id;
      var isSel = s.id === sel.id;
      var col = isCur ? COL.white : (isSel ? COL.yellow : E.color('teal', 3));
      if (isCur) {
        fb.fillRect(x - 1, y - 1, x + 1, y + 1, COL.white);
      } else if (isSel) {
        fb.fillRect(x - 1, y - 1, x + 1, y + 1, COL.yellow);
        fb.rect(x - 3, y - 3, x + 3, y + 3, COL.orange);
      } else {
        fb.px(x, y, col);
      }
    }

    plate(fb, 206, 18, 110, 190, 'SYSTEM DATA');
    var ty = 32;
    fb.text(sel.name, 212, ty, COL.yellow); ty += 11;
    fb.text('ECONOMY', 212, ty, COL.dim); ty += 8;
    fb.text(G.ECONOMY_NAMES[sel.economy], 212, ty, COL.text); ty += 11;
    fb.text('GOVERNMENT', 212, ty, COL.dim); ty += 8;
    fb.text(G.GOV_NAMES[sel.government], 212, ty, COL.text); ty += 11;
    fb.text('TECH LEVEL ' + sel.tech, 212, ty, COL.text); ty += 9;
    fb.text('POP ' + sel.pop + ' B', 212, ty, COL.text); ty += 9;
    fb.text('PROD ' + sel.prod + ' MC', 212, ty, COL.text); ty += 9;
    fb.text('RADIUS ' + sel.radius + ' KM', 212, ty, COL.text); ty += 11;
    var d = G.distance(cur, sel);
    fb.text('DIST ' + d.toFixed(1) + ' LY', 212, ty, d <= 7 ? COL.green : COL.red); ty += 9;
    fb.text('FUEL ' + cmd.fuel.toFixed(1) + ' LY', 212, ty, COL.text); ty += 11;
    if (sel.id !== cur.id) {
      fb.text(d <= 7 && cmd.fuel >= d ? 'H - HYPERSPACE' : 'OUT OF RANGE', 212, ty, d <= 7 ? COL.yellow : COL.dim);
    }

  };

  UI.status = function (fb) {
    UI.hint = 'ESC BACK';
    var cmd = GAME.cmd;
    plate(fb, 4, 18, 152, 190, 'COMMANDER');
    var y = 32;
    var rows = [
      ['NAME', cmd.name], ['CASH', money(cmd.cash) + ' CR'],
      ['FUEL', cmd.fuel.toFixed(1) + ' LY'], ['CARGO', GAME.cargoUsed(cmd) + '/' + cmd.cargoCap + ' T'],
      ['MISSILES', String(cmd.missiles)], ['LASER', GAME.C.laserName[cmd.laser]],
      ['RANK', G.rankFor(cmd.kills)], ['KILLS', String(cmd.kills)],
      ['LEGAL', cmd.legal >= 500 ? 'FUGITIVE' : (cmd.legal > 0 ? 'OFFENDER' : 'CLEAN')],
      ['SYSTEM', GAME.world.sys.name], ['GALAXY', String(cmd.galaxy + 1)]
    ];
    for (var i = 0; i < rows.length; i++) {
      fb.text(rows[i][0], 12, y, COL.dim);
      fb.textR(rows[i][1], 150, y, COL.text);
      y += 12;
    }
    plate(fb, 160, 18, 156, 190, 'EQUIPMENT');
    var eq = [['E.C.M.', cmd.equipment.ecm], ['FUEL SCOOP', cmd.equipment.fuelscoop],
    ['LARGE BAY', cmd.equipment.largebay], ['ESCAPE CAPSULE', cmd.equipment.escape],
    ['ENERGY UNIT', cmd.equipment.extraenergy], ['DOCKING COMP', cmd.equipment.dockcomp],
    ['GAL. HYPERDRIVE', cmd.equipment.galdrive], ['ENERGY BOMB', cmd.bombs > 0]];
    var yy = 32;
    for (var k = 0; k < eq.length; k++) {
      fb.text(eq[k][0], 168, yy, eq[k][1] ? COL.green : E.rgb('#2d4a5c'));
      fb.textR(eq[k][1] ? 'YES' : '-', 308, yy, eq[k][1] ? COL.green : COL.dim);
      yy += 11;
    }
    yy += 6;
    fb.text('CARGO:', 168, yy, COL.dim); yy += 10;
    var any = false;
    for (var c in cmd.cargo) {
      if (cmd.cargo[c] > 0) {
        any = true;
        fb.text(GAME.commoditiesName(+c) + ' x' + cmd.cargo[c], 168, yy, COL.text);
        yy += 9;
      }
    }
    if (!any) fb.text('EMPTY', 168, yy, COL.dim);
  };

  // --------------------------------------------------------------- dispatch
  UI.draw = function (fb, dt) {
    if (GAME.state === 'flight' && GAME.overlay) {
      if (GAME.overlay === 'pause') { UI.pause(fb); return; }
      UI.chrome(fb, GAME.overlay === 'chart' ? 'GALACTIC CHART - IN FLIGHT' : 'COMMANDER STATUS');
      if (GAME.overlay === 'chart') UI.chart(fb); else UI.status(fb);
      UI.footer(fb);
      return;
    }
    switch (GAME.state) {
      case 'title': UI.title(fb, dt); break;
      case 'dead': UI.gameOver(fb); break;
      case 'hyper': UI.hyper(fb, dt); break;
      case 'docked': UI.docked(fb, dt); break;
      default: UI.flight(fb, dt); break;
    }
  };

  // ------------------------------------------------------------------ keys
  UI.key = function (k) {
    var cmd = GAME.cmd;
    if (GAME.state === 'title') {
      if (k === '1' || k === 'Enter') { GAME.newGame(UI.sel.name); return true; }
      if (k === '2') { if (GAME.load()) return true; }
      return true;
    }
    if (GAME.state === 'dead') {
      if (k === 'Enter') { GAME.state = 'title'; }
      return true;
    }
    if (GAME.state === 'docked') {
      if (GAME.screen === 'launch') {
        if (k === 'ArrowUp') { UI.sel.menu = (UI.sel.menu + 5) % 6; return true; }
        if (k === 'ArrowDown') { UI.sel.menu = (UI.sel.menu + 1) % 6; return true; }
        if (k === 'f' || k === 'F') { GAME.refuel(); return true; }
        var go = null;
        if (k === '1') go = 'market';
        if (k === '2') go = 'equip';
        if (k === '3') go = 'chart';
        if (k === '4') go = 'status';
        if (k === '5') { GAME.launch(); return true; }
        if (k === '6') { GAME.message(GAME.save() ? 'SAVED' : 'SAVE FAILED', 'info'); return true; }
        if (k === 'Enter') {
          go = ['market', 'equip', 'chart', 'status', 'launch', 'save'][UI.sel.menu];
          if (go === 'launch') { GAME.launch(); return true; }
          if (go === 'save') { GAME.message(GAME.save() ? 'SAVED' : 'SAVE FAILED', 'info'); return true; }
        }
        if (go) { GAME.screen = go; if (go === 'chart') UI.sel.chartCursor = cmd.system; return true; }
        return true;
      }
      if (k === 'Escape') { GAME.screen = 'launch'; return true; }
      if (GAME.screen === 'market') {
        if (k === 'ArrowUp') { UI.sel.market = (UI.sel.market + 16) % 17; return true; }
        if (k === 'ArrowDown') { UI.sel.market = (UI.sel.market + 1) % 17; return true; }
        if (k === 'b') { GAME.buy(UI.sel.market, 1); return true; }
        if (k === 'B') { GAME.buy(UI.sel.market, 10); return true; }
        if (k === 's') { GAME.sell(UI.sel.market, 1); return true; }
        if (k === 'S') { GAME.sell(UI.sel.market, 10); return true; }
        return true;
      }
      if (GAME.screen === 'equip') {
        var items = G.availableEquipment(GAME.world.sys);
        if (k === 'ArrowUp') { UI.sel.equip = (UI.sel.equip + items.length - 1) % items.length; return true; }
        if (k === 'ArrowDown') { UI.sel.equip = (UI.sel.equip + 1) % items.length; return true; }
        if (k === 'Enter') { GAME.buyEquipment(items[UI.sel.equip].id); return true; }
        return true;
      }
      if (GAME.screen === 'chart') {
        var gal = G.get(cmd.galaxy);
        var curId = UI.sel.chartCursor === null ? cmd.system : UI.sel.chartCursor;
        if (k === 'ArrowUp' || k === 'ArrowDown' || k === 'ArrowLeft' || k === 'ArrowRight') {
          UI.sel.chartCursor = nearestInDirection(gal, curId, k);
          return true;
        }
        if (k === 'h' || k === 'H') {
          var target = UI.sel.chartCursor === null || UI.sel.chartCursor === undefined ? cmd.system : UI.sel.chartCursor;
          if (GAME.hyperspace(target)) GAME.overlay = null;
          return true;
        }
        return true;
      }
      return true;
    }
    // in flight
    if (k === 'c' || k === 'C') { return false; }
    return false;
  };

  function nearestInDirection(gal, fromId, key) {
    var from = gal.systems[fromId];
    var dirs = { ArrowUp: [0, 1], ArrowDown: [0, -1], ArrowLeft: [-1, 0], ArrowRight: [1, 0] };
    var d = dirs[key];
    if (!d) return fromId;
    var best = null, bestScore = 1e9;
    for (var i = 0; i < gal.systems.length; i++) {
      var s = gal.systems[i];
      if (s.id === fromId) continue;
      // wrap-around chart
      var dx = s.x - from.x, dy = s.y - from.y;
      if (dx > 128) dx -= 256; if (dx < -128) dx += 256;
      if (dy > 128) dy -= 256; if (dy < -128) dy += 256;
      var along = dx * d[0] + dy * d[1];
      if (along <= 2) continue;
      var perp = Math.abs(dx * d[1] - dy * d[0]);
      var score = along + perp * 2.5;
      if (score < bestScore) { bestScore = score; best = s.id; }
    }
    return best === null ? fromId : best;
  }

  UI.nearestInDirection = nearestInDirection;
})(typeof window !== 'undefined' ? window : globalThis);
