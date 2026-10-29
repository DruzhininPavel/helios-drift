/* ELITE 16 - models
 * Low-poly flat-shaded meshes. All faces carry an outward normal so the
 * renderer can do backface culling + per-face lighting.
 * Coordinate convention: -Z is forward (ship nose), +Y up, +X right.
 */
(function (global) {
  'use strict';
  var E = global.E;
  var V = E.V, M3 = E.M3;
  var M = E.models = {};

  var EPS = 1e-6;

  function Builder() {
    this.faces = [];
    this.center = [0, 0, 0];
    this.radius = 0;
  }
  Builder.prototype.at = function (c) { this.center = c; return this; };

  Builder.prototype.face = function (pts, ramp, opts) {
    opts = opts || {};
    if (pts.length < 3) return this;
    // drop degenerate faces
    var n = V.cross(V.sub(pts[1], pts[0]), V.sub(pts[2], pts[0]));
    if (V.len(n) < EPS) return this;
    n = V.norm(n);
    var c = [0, 0, 0];
    for (var i = 0; i < pts.length; i++) { c[0] += pts[i][0]; c[1] += pts[i][1]; c[2] += pts[i][2]; }
    c[0] /= pts.length; c[1] /= pts.length; c[2] /= pts.length;
    if (V.dot(n, V.sub(c, this.center)) < 0) {
      n = V.mul(n, -1);
      pts = pts.slice().reverse();
    }
    this.faces.push({
      p: pts, ramp: ramp, n: n,
      shade: opts.shade, double: opts.double, glow: opts.glow, tint: opts.tint
    });
    for (var k = 0; k < pts.length; k++) {
      var d = V.len(pts[k]);
      if (d > this.radius) this.radius = d;
    }
    return this;
  };

  Builder.prototype.mesh = function () {
    orientShells(this.faces);
    return { faces: this.faces, radius: this.radius };
  };

  // Normals are authored by winding order, which is easy to get wrong. Group
  // faces into shells (parts share no vertices) and flip each shell's faces so
  // they point outwards from that shell's own centroid. Flat patches (docking
  // face panels, slot walls) fall back to the model origin.
  function orientShells(faces) {
    var parent = {};
    function keyOf(p) {
      return (Math.round(p[0] * 100) / 100) + ',' + (Math.round(p[1] * 100) / 100) + ',' + (Math.round(p[2] * 100) / 100);
    }
    function find(a) {
      var root = a;
      while (parent[root] !== root) root = parent[root];
      while (parent[a] !== root) { var n = parent[a]; parent[a] = root; a = n; }
      return root;
    }
    var i, j;
    for (i = 0; i < faces.length; i++) {
      var f = faces[i];
      var k0 = keyOf(f.p[0]);
      if (!(k0 in parent)) parent[k0] = k0;
      for (j = 1; j < f.p.length; j++) {
        var kj = keyOf(f.p[j]);
        if (!(kj in parent)) parent[kj] = kj;
        var ra = find(k0), rb = find(kj);
        if (ra !== rb) parent[rb] = ra;
      }
    }
    // per-shell vertex sum, area-weighted normal sum
    var acc = {};
    for (i = 0; i < faces.length; i++) {
      var fc = faces[i];
      var root = find(keyOf(fc.p[0]));
      var a = acc[root] || (acc[root] = { sx: 0, sy: 0, sz: 0, n: 0, nx: 0, ny: 0, nz: 0 });
      for (j = 0; j < fc.p.length; j++) { a.sx += fc.p[j][0]; a.sy += fc.p[j][1]; a.sz += fc.p[j][2]; a.n++; }
      var u = V.sub(fc.p[1], fc.p[0]), v = V.sub(fc.p[2], fc.p[0]);
      var area = V.len(V.cross(u, v)) * 0.5;
      a.nx += fc.n[0] * area; a.ny += fc.n[1] * area; a.nz += fc.n[2] * area;
      a.area = (a.area || 0) + area;
    }
    for (i = 0; i < faces.length; i++) {
      var face = faces[i];
      var r = find(keyOf(face.p[0]));
      var s = acc[r];
      var cx = s.sx / s.n, cy = s.sy / s.n, cz = s.sz / s.n;
      var nlen = Math.sqrt(s.nx * s.nx + s.ny * s.ny + s.nz * s.nz);
      var flatness = s.area > 0 ? nlen / s.area : 0;
      var ref;
      if (flatness > 0.6) ref = [0, 0, 0];           // flat patch: away from origin
      else ref = [cx, cy, cz];                        // closed shell: away from its own centre
      var fc2 = [0, 0, 0];
      for (j = 0; j < face.p.length; j++) {
        fc2[0] += face.p[j][0]; fc2[1] += face.p[j][1]; fc2[2] += face.p[j][2];
      }
      fc2[0] /= face.p.length; fc2[1] /= face.p.length; fc2[2] /= face.p.length;
      var dir = V.sub(fc2, ref);
      if (V.len(dir) < 0.001) continue;
      if (V.dot(face.n, dir) < 0) {
        face.n = V.mul(face.n, -1);
        face.p = face.p.slice().reverse();
      }
    }
  }

  // loft an elliptical-section hull through stations [{z, w, h, y}]
  function loft(b, stations, ramp, opts) {
    opts = opts || {};
    var sides = opts.sides || 8;
    var phase = opts.phase === undefined ? Math.PI / sides : opts.phase;
    var rings = [];
    var i, j;
    for (i = 0; i < stations.length; i++) {
      var st = stations[i];
      var ring = [];
      for (j = 0; j < sides; j++) {
        var a = phase + (j / sides) * Math.PI * 2;
        ring.push([
          Math.cos(a) * (st.w || 0),
          (st.y || 0) + Math.sin(a) * (st.h || 0),
          st.z
        ]);
      }
      rings.push(ring);
    }
    for (i = 0; i + 1 < rings.length; i++) {
      for (j = 0; j < sides; j++) {
        var a0 = rings[i][j], a1 = rings[i][(j + 1) % sides];
        var b0 = rings[i + 1][j], b1 = rings[i + 1][(j + 1) % sides];
        var shade = opts.shadeFor ? undefined : undefined;
        b.face([a0, a1, b1, b0], ramp, opts);
      }
    }
    if (opts.capNose !== false) b.face(rings[0].slice().reverse(), ramp, opts);
    if (opts.capTail !== false) b.face(rings[rings.length - 1].slice(), ramp, opts);
    return rings;
  }

  // extrude a 2D polygon defined in the XZ plane along Y
  function prism(b, poly, thick, ramp, opts) {
    opts = opts || {};
    var y = thick / 2;
    var top = poly.map(function (p) { return [p[0], y, p[1]]; });
    var bot = poly.map(function (p) { return [p[0], -y, p[1]]; });
    b.face(top, ramp, opts);
    b.face(bot, ramp, opts);
    for (var i = 0; i < poly.length; i++) {
      var j = (i + 1) % poly.length;
      b.face([top[i], top[j], bot[j], bot[i]], ramp, opts);
    }
  }

  function box(b, cx, cy, cz, sx, sy, sz, ramp, opts) {
    var x = sx / 2, y = sy / 2, z = sz / 2;
    var v = [
      [cx - x, cy + y, cz - z], [cx + x, cy + y, cz - z],
      [cx + x, cy - y, cz - z], [cx - x, cy - y, cz - z],
      [cx - x, cy + y, cz + z], [cx + x, cy + y, cz + z],
      [cx + x, cy - y, cz + z], [cx - x, cy - y, cz + z]
    ];
    b.face([v[0], v[1], v[2], v[3]], ramp, opts);
    b.face([v[5], v[4], v[7], v[6]], ramp, opts);
    b.face([v[4], v[0], v[3], v[7]], ramp, opts);
    b.face([v[1], v[5], v[6], v[2]], ramp, opts);
    b.face([v[4], v[5], v[1], v[0]], ramp, opts);
    b.face([v[3], v[2], v[6], v[7]], ramp, opts);
  }

  // ------------------------------------------------------------- the ships
  // Cobra Mk III - flat wide delta wedge with raised cockpit
  function cobra3() {
    var b = new Builder().at([0, 0, 0]);
    loft(b, [
      { z: -34, w: 1.2, h: 1.0 },
      { z: -24, w: 5.5, h: 2.6 },
      { z: -10, w: 9.0, h: 3.4, y: 0.4 },
      { z: 8, w: 11.5, h: 3.2, y: 0.2 },
      { z: 26, w: 9.5, h: 2.6 },
      { z: 31, w: 4.0, h: 1.8 }
    ], 'steel', { sides: 8 });
    // wings: sweeping delta planforms
    prism(b, [[4, -18], [30, 16], [31, 24], [9, 26], [7, 4]], 1.8, 'steel');
    prism(b, [[-4, -18], [-30, 16], [-31, 24], [-9, 26], [-7, 4]], 1.8, 'steel');
    // wing tip cannons
    box(b, 29, 0, 18, 2.0, 2.0, 12, 'ash');
    box(b, -29, 0, 18, 2.0, 2.0, 12, 'ash');
    // cockpit canopy
    loft(b, [
      { z: -22, w: 1.0, h: 0.8, y: 2.4 },
      { z: -14, w: 4.2, h: 2.2, y: 2.8 },
      { z: 0, w: 4.0, h: 1.6, y: 3.0 },
      { z: 8, w: 2.0, h: 0.8, y: 2.6 }
    ], 'teal', { sides: 6, shade: 4.2 });
    // engine bells
    box(b, 4.4, 0.2, 30, 4.0, 3.0, 4.0, 'ash');
    box(b, -4.4, 0.2, 30, 4.0, 3.0, 4.0, 'ash');
    // spine detail
    box(b, 0, 3.1, 18, 2.0, 1.4, 12, 'steel');
    return b.mesh();
  }

  // Sidewinder - compact wedge fighter
  function sidewinder() {
    var b = new Builder().at([0, 0, 0]);
    loft(b, [
      { z: -28, w: 1.0, h: 0.8 },
      { z: -14, w: 6.0, h: 2.4 },
      { z: 2, w: 8.0, h: 3.4 },
      { z: 18, w: 6.0, h: 3.0 },
      { z: 24, w: 3.0, h: 2.0 }
    ], 'ash', { sides: 6 });
    prism(b, [[5, -6], [24, 20], [23, 25], [7, 20]], 1.4, 'ash');
    prism(b, [[-5, -6], [-24, 20], [-23, 25], [-7, 20]], 1.4, 'ash');
    box(b, 0, 3.2, -6, 2.0, 1.6, 8, 'steel');
    box(b, 0, 0, 23, 5.0, 3.0, 3.0, 'red', { shade: 2.2 });
    return b.mesh();
  }

  // Krait - diamond planform, sharp at both ends
  function krait() {
    var b = new Builder().at([0, 0, 0]);
    loft(b, [
      { z: -30, w: 0.8, h: 0.7 },
      { z: -18, w: 5.0, h: 2.4 },
      { z: -4, w: 13.0, h: 4.0 },
      { z: 6, w: 14.0, h: 4.2, y: 0.4 },
      { z: 18, w: 7.0, h: 2.6 },
      { z: 26, w: 2.5, h: 1.4 }
    ], 'green', { sides: 6 });
    box(b, 0, 3.0, -8, 2.6, 2.4, 12, 'alien', { shade: 4.0 });
    box(b, 4.0, 0, 23, 3.4, 2.6, 5.0, 'ash');
    box(b, -4.0, 0, 23, 3.4, 2.6, 5.0, 'ash');
    return b.mesh();
  }

  // Mamba - sleek dart
  function mamba() {
    var b = new Builder().at([0, 0, 0]);
    loft(b, [
      { z: -32, w: 0.8, h: 0.7 },
      { z: -16, w: 4.0, h: 2.0 },
      { z: 2, w: 7.5, h: 3.2 },
      { z: 20, w: 6.0, h: 2.6 },
      { z: 27, w: 3.0, h: 1.8 }
    ], 'purple', { sides: 6 });
    prism(b, [[6, -2], [22, 22], [20, 26], [8, 20]], 1.2, 'purple');
    prism(b, [[-6, -2], [-22, 22], [-20, 26], [-8, 20]], 1.2, 'purple');
    // long nose spike
    loft(b, [
      { z: -40, w: 0.4, h: 0.4 },
      { z: -28, w: 1.6, h: 1.2 }
    ], 'ash', { sides: 5 });
    box(b, 0, 2.8, 0, 2.0, 1.6, 12, 'teal', { shade: 3.6 });
    return b.mesh();
  }

  // Adder - boxy trader
  function adder() {
    var b = new Builder().at([0, 0, 0]);
    loft(b, [
      { z: -24, w: 2.0, h: 2.0 },
      { z: -18, w: 7.0, h: 3.0 },
      { z: 10, w: 8.5, h: 4.0, y: 0.4 },
      { z: 22, w: 7.0, h: 3.4 }
    ], 'brown', { sides: 8 });
    prism(b, [[5, -8], [22, 14], [22, 20], [5, 18]], 2.6, 'brown');
    prism(b, [[-5, -8], [-22, 14], [-22, 20], [-5, 18]], 2.6, 'brown');
    box(b, 0, 4.2, -8, 3.0, 2.0, 10, 'steel');
    box(b, 0, 0, 23, 8.0, 4.0, 3.0, 'orange', { shade: 2.4 });
    return b.mesh();
  }

  // Python - big slow trader
  function python() {
    var b = new Builder().at([0, 0, 0]);
    loft(b, [
      { z: -40, w: 3.0, h: 3.0 },
      { z: -30, w: 9.0, h: 7.0 },
      { z: -10, w: 13.0, h: 9.0 },
      { z: 20, w: 13.0, h: 8.0 },
      { z: 34, w: 9.0, h: 6.0 },
      { z: 38, w: 5.0, h: 4.0 }
    ], 'orange', { sides: 8 });
    box(b, 15, -1, 12, 8.0, 6.0, 24, 'brown');
    box(b, -15, -1, 12, 8.0, 6.0, 24, 'brown');
    box(b, 0, 8.8, -2, 5.0, 2.6, 16, 'steel');
    box(b, 0, 0, 37, 14.0, 8.0, 4.0, 'ash');
    return b.mesh();
  }

  // Viper - police interceptor
  function viper() {
    var b = new Builder().at([0, 0, 0]);
    loft(b, [
      { z: -30, w: 1.6, h: 1.6 },
      { z: -20, w: 5.0, h: 3.0 },
      { z: -4, w: 8.0, h: 4.0 },
      { z: 14, w: 7.0, h: 3.4 },
      { z: 24, w: 4.0, h: 2.4 }
    ], 'blue', { sides: 8 });
    prism(b, [[6, -4], [26, 16], [26, 22], [7, 18]], 1.6, 'blue');
    prism(b, [[-6, -4], [-26, 16], [-26, 22], [-7, 18]], 1.6, 'blue');
    // police stripes
    box(b, 0, 4.2, -6, 3.2, 1.2, 14, 'ash');
    box(b, 4.2, 0, 25, 3.0, 2.6, 3.0, 'yellow', { shade: 4.4 });
    box(b, -4.2, 0, 25, 3.0, 2.6, 3.0, 'yellow', { shade: 4.4 });
    return b.mesh();
  }

  // Asp Mk II - fast heavy fighter
  function asp() {
    var b = new Builder().at([0, 0, 0]);
    loft(b, [
      { z: -36, w: 1.0, h: 1.0 },
      { z: -20, w: 5.0, h: 3.4 },
      { z: 0, w: 8.5, h: 4.6 },
      { z: 20, w: 9.0, h: 4.0 },
      { z: 30, w: 6.0, h: 3.0 }
    ], 'teal', { sides: 8 });
    // three-engine tail fins
    prism(b, [[7, 10], [20, 30], [20, 33], [8, 26]], 1.4, 'teal');
    prism(b, [[-7, 10], [-20, 30], [-20, 33], [-8, 26]], 1.4, 'teal');
    box(b, 0, 4.6, 4, 3.0, 2.0, 14, 'ash');
    box(b, 0, 0, 29.5, 4.0, 3.0, 4.0, 'pink', { shade: 4.0 });
    return b.mesh();
  }

  // Moray Starboat - amphibious shape
  function moray() {
    var b = new Builder().at([0, 0, 0]);
    loft(b, [
      { z: -28, w: 2.0, h: 1.4 },
      { z: -12, w: 7.0, h: 2.6 },
      { z: 6, w: 8.0, h: 3.0 },
      { z: 22, w: 5.0, h: 2.0 }
    ], 'teal', { sides: 8 });
    prism(b, [[7, -2], [18, 10], [17, 18], [7, 16]], 1.2, 'teal');
    prism(b, [[-7, -2], [-18, 10], [-17, 18], [-7, 16]], 1.2, 'teal');
    box(b, 0, 2.8, -4, 2.6, 1.8, 12, 'green', { shade: 4.2 });
    prism(b, [[2, -18], [2, 16], [-1, 22], [-3, 16], [-3, -18]], 1.2, 'teal');
    box(b, 0, -3.2, 2, 3.0, 4.0, 18, 'steel', { shade: 2.2 });
    return b.mesh();
  }

  // Fer-de-Lance - luxury killer
  function ferdelance() {
    var b = new Builder().at([0, 0, 0]);
    loft(b, [
      { z: -30, w: 1.0, h: 1.0 },
      { z: -14, w: 5.0, h: 3.0 },
      { z: 4, w: 10.0, h: 4.0 },
      { z: 18, w: 8.0, h: 3.0 },
      { z: 26, w: 4.0, h: 2.0 }
    ], 'pink', { sides: 8 });
    prism(b, [[9, -4], [30, 14], [30, 20], [10, 18]], 1.4, 'pink');
    prism(b, [[-9, -4], [-30, 14], [-30, 20], [-10, 18]], 1.4, 'pink');
    loft(b, [
      { z: -16, w: 1.0, h: 0.6, y: 2.6 },
      { z: -6, w: 3.2, h: 1.6, y: 3.0 },
      { z: 6, w: 2.0, h: 1.0, y: 2.8 }
    ], 'yellow', { sides: 6, shade: 4.4 });
    box(b, 5, 0, 25, 3.0, 2.4, 3.0, 'yellow', { shade: 3.4 });
    box(b, -5, 0, 25, 3.0, 2.4, 3.0, 'yellow', { shade: 3.4 });
    return b.mesh();
  }

  // Thargoid - flattened octagonal saucer with a dorsal dome
  function thargoid() {
    var b = new Builder().at([0, 0, 0]);
    loft(b, [
      { z: -22, w: 1.5, h: 1.2 },
      { z: -16, w: 9.0, h: 3.0 },
      { z: -7, w: 17.0, h: 4.2 },
      { z: 2, w: 18.5, h: 4.4 },
      { z: 11, w: 14.0, h: 3.6 },
      { z: 19, w: 6.0, h: 2.4 }
    ], 'alien', { sides: 8 });
    loft(b, [
      { z: -6, w: 5.0, h: 2.2, y: 3.4 },
      { z: 1, w: 6.5, h: 3.4, y: 4.2 },
      { z: 8, w: 4.0, h: 2.0, y: 3.4 }
    ], 'green', { sides: 8, shade: 3.8 });
    for (var i = 0; i < 6; i++) {
      var ang = i * Math.PI / 3;
      box(b, Math.cos(ang) * 13, -0.5, Math.sin(ang) * 9 + 1, 3.6, 3.0, 3.4, 'alien', { shade: 4.4 });
    }
    return b.mesh();
  }

  // Shuttle - small boxy craft
  function shuttle() {
    var b = new Builder().at([0, 0, 0]);
    box(b, 0, 0, 0, 10, 6, 22, 'yellow');
    loft(b, [
      { z: -20, w: 2.0, h: 2.0 },
      { z: -11, w: 5.0, h: 3.0 }
    ], 'yellow', { sides: 6 });
    box(b, 8, -1, 3, 8, 2.2, 12, 'ash');
    box(b, -8, -1, 3, 8, 2.2, 12, 'ash');
    box(b, 0, 3.2, 2, 4.0, 2.0, 10, 'steel');
    return b.mesh();
  }

  // ------------------------------------------------------------ other stuff
  function asteroidMesh(seed) {
    var b = new Builder().at([0, 0, 0]);
    var t = (1 + Math.sqrt(5)) / 2;
    var v = [
      [-1, t, 0], [1, t, 0], [-1, -t, 0], [1, -t, 0],
      [0, -1, t], [0, 1, t], [0, -1, -t], [0, 1, -t],
      [t, 0, -1], [t, 0, 1], [-t, 0, -1], [-t, 0, 1]
    ];
    var f = [
      [0, 11, 5], [0, 5, 1], [0, 1, 7], [0, 7, 10], [0, 10, 11],
      [1, 5, 9], [5, 11, 4], [11, 10, 2], [10, 7, 6], [7, 1, 8],
      [3, 9, 4], [3, 4, 2], [3, 2, 6], [3, 6, 8], [3, 8, 9],
      [4, 9, 5], [2, 4, 11], [6, 2, 10], [8, 6, 7], [9, 8, 1]
    ];
    var rnd = E.rng(seed || 7);
    var scaled = v.map(function (p) {
      var j = 0.66 + rnd() * 0.6;
      return [p[0] * j * 11, p[1] * j * 11, p[2] * j * 11];
    });
    for (var i = 0; i < f.length; i++) {
      b.face([scaled[f[i][0]], scaled[f[i][1]], scaled[f[i][2]]], 'brown');
    }
    return b.mesh();
  }

  function canisterMesh() {
    var b = new Builder().at([0, 0, 0]);
    loft(b, [
      { z: -7, w: 2.2, h: 2.2 },
      { z: -5, w: 4.0, h: 4.0 },
      { z: 5, w: 4.0, h: 4.0 },
      { z: 7, w: 2.2, h: 2.2 }
    ], 'yellow', { sides: 8 });
    box(b, 0, 0, 0, 8.3, 8.3, 3.0, 'orange', { shade: 3.0 });
    return b.mesh();
  }

  function capsuleMesh() {
    var b = new Builder().at([0, 0, 0]);
    loft(b, [
      { z: -9, w: 1.4, h: 1.4 },
      { z: -5, w: 3.4, h: 3.4 },
      { z: 3, w: 4.0, h: 3.6 },
      { z: 8, w: 2.0, h: 2.0 }
    ], 'red', { sides: 8 });
    box(b, 0, 3.0, -1, 2.0, 1.4, 5.0, 'ash');
    return b.mesh();
  }

  function missileMesh() {
    var b = new Builder().at([0, 0, 0]);
    loft(b, [
      { z: -10, w: 0.3, h: 0.3 },
      { z: -4, w: 1.3, h: 1.3 },
      { z: 6, w: 1.3, h: 1.3 },
      { z: 8, w: 0.9, h: 0.9 }
    ], 'ash', { sides: 6 });
    box(b, 0, 0, 7, 0.4, 5.0, 3.0, 'red', { shade: 3.4 });
    box(b, 0, 0, 7, 5.0, 0.4, 3.0, 'red', { shade: 3.4 });
    return b.mesh();
  }

  // Coriolis station: cuboctahedron (6 diamonds + 8 triangles), rotating,
  // docking slot cut through the +Z face.
  M.CORIOLIS = { S: 220, slotW: 0.31, slotH: 0.072 };
  function coriolisMesh() {
    var b = new Builder().at([0, 0, 0]);
    var S = M.CORIOLIS.S;
    var ramp = 'ash';
    var V12 = [];
    var ax;
    for (ax = 0; ax < 3; ax++) {
      for (var s1 = -1; s1 <= 1; s1 += 2) for (var s2 = -1; s2 <= 1; s2 += 2) {
        var p = [0, 0, 0];
        var o = [0, 1, 2].filter(function (z) { return z !== ax; });
        p[o[0]] = s1 * S; p[o[1]] = s2 * S;
        V12.push(p);
      }
    }
    // diamonds on the five non-slot faces
    var axes = [[2, 1], [2, -1], [0, 1], [0, -1], [1, 1], [1, -1]];
    for (var i = 0; i < axes.length; i++) {
      var a = axes[i][0], sgn = axes[i][1];
      if (a === 2 && sgn === 1) continue; // slot face handled separately
      var pts = [];
      var others = [0, 1, 2].filter(function (z) { return z !== a; });
      var ring = [[1, 0], [0, 1], [-1, 0], [0, -1]];
      for (var r = 0; r < 4; r++) {
        var q = [0, 0, 0];
        q[a] = sgn * S;
        q[others[0]] = ring[r][0] * S;
        q[others[1]] = ring[r][1] * S;
        pts.push(q);
      }
      b.face(pts, ramp, { tint: 0.9 + (i % 3) * 0.09 });
    }
    // 8 corner triangles
    for (var cx = -1; cx <= 1; cx += 2) for (var cy = -1; cy <= 1; cy += 2) for (var cz = -1; cz <= 1; cz += 2) {
      var c = [cx, cy, cz];
      var tri = [];
      for (var k = 0; k < 3; k++) {
        var pt = [0, 0, 0];
        pt[k] = 0;
        var o2 = [0, 1, 2].filter(function (z) { return z !== k; });
        pt[o2[0]] = c[o2[0]] * S;
        pt[o2[1]] = c[o2[1]] * S;
        tri.push(pt);
      }
      b.face(tri, ramp, { tint: 1.12 });
    }
    // +Z face with docking slot: diamond |x|+|y| <= S, slot is a rect
    var w = M.CORIOLIS.slotW * S, h = M.CORIOLIS.slotH * S;
    var dh = S - h;
    // left / right of the slot (pentagons down to the diamond corners)
    b.face([[-S, 0, S], [-dh, h, S], [-w, h, S], [-w, -h, S], [-dh, -h, S]], ramp);
    b.face([[S, 0, S], [dh, -h, S], [w, -h, S], [w, h, S], [dh, h, S]], ramp);
    // above / below the slot
    b.face([[-dh, h, S], [dh, h, S], [0, S, S]], ramp);
    b.face([[-dh, -h, S], [0, -S, S], [dh, -h, S]], ramp);
    // slot tunnel: walls running inwards, then a dark interior plate
    var T = S - 54;
    b.face([[-w, h, S], [-w, h, T], [w, h, T], [w, h, S]], 'steel', { shade: 1.2 });
    b.face([[-w, -h, S], [w, -h, S], [w, -h, T], [-w, -h, T]], 'steel', { shade: 1.2 });
    b.face([[-w, h, S], [-w, -h, S], [-w, -h, T], [-w, h, T]], 'steel', { shade: 1.6 });
    b.face([[w, h, S], [w, h, T], [w, -h, T], [w, -h, S]], 'steel', { shade: 1.6 });
    b.face([[-w, h, T], [-w, -h, T], [w, -h, T], [w, h, T]], 'steel', { shade: 0.35 });
    // greebles
    box(b, 0, 140, 140, 60, 22, 60, 'steel', { shade: 2.4 });
    box(b, 0, -140, 140, 60, 22, 60, 'steel', { shade: 2.4 });
    box(b, 140, 0, -140, 46, 46, 46, 'steel', { shade: 2.0 });
    box(b, -140, 0, -140, 46, 46, 46, 'steel', { shade: 2.0 });
    box(b, 0, 0, -S + 12, 60, 60, 18, 'steel', { shade: 2.2 });
    return b.mesh();
  }

  M.meshes = {
    cobra3: cobra3(),
    sidewinder: sidewinder(),
    krait: krait(),
    mamba: mamba(),
    adder: adder(),
    python: python(),
    viper: viper(),
    asp: asp(),
    moray: moray(),
    ferdelance: ferdelance(),
    thargoid: thargoid(),
    shuttle: shuttle(),
    asteroid: asteroidMesh(11),
    asteroid2: asteroidMesh(29),
    canister: canisterMesh(),
    capsule: capsuleMesh(),
    missile: missileMesh(),
    station: coriolisMesh()
  };
  M.Builder = Builder;
  M.loft = loft;
  M.prism = prism;
  M.box = box;

  // ------------------------------------------------------------- planets
  // Value noise on the unit sphere, pre-rendered into a sprite with a
  // rotating texture and a hard terminator. 16-bit style: dither only.
  function hash3(x, y, z) {
    var n = Math.sin(x * 127.1 + y * 311.7 + z * 74.7) * 43758.5453;
    return n - Math.floor(n);
  }
  function vnoise(x, y, z) {
    var xi = Math.floor(x), yi = Math.floor(y), zi = Math.floor(z);
    var xf = x - xi, yf = y - yi, zf = z - zi;
    function s(t) { return t * t * (3 - 2 * t); }
    var u = s(xf), v = s(yf), w = s(zf);
    var c = 0;
    for (var dx = 0; dx <= 1; dx++) for (var dy = 0; dy <= 1; dy++) for (var dz = 0; dz <= 1; dz++) {
      var h = hash3(xi + dx, yi + dy, zi + dz);
      var wx = dx ? u : 1 - u, wy = dy ? v : 1 - v, wz = dz ? w : 1 - w;
      c += h * wx * wy * wz;
    }
    return c;
  }
  function fbm(x, y, z, oct) {
    var s = 0, a = 0.5, f = 1;
    for (var i = 0; i < oct; i++) { s += a * vnoise(x * f, y * f, z * f); a *= 0.5; f *= 2.05; }
    return s;
  }

  // kind: 0 = terran (blue/green), 1 = desert, 2 = ice, 3 = gas giant, 4 = lava
  // Built in two passes: a lon/lat texture table once, then cheap rotated
  // frame renders that sample it, so entering a system costs a few ms.
  var planetCache = {};

  function buildPlanetTexture(seed, kind) {
    var W = 256, H = 128;
    var rnd = E.rng(seed);
    var ox = rnd() * 100, oy = rnd() * 100, oz = rnd() * 100;
    var lum = new Uint8Array(W * H);
    var rampIdx = new Uint8Array(W * H);
    var ramps = { 0: 'blue', 1: 'teal', 2: 'green', 3: 'brown', 4: 'ash', 5: 'steel', 6: 'orange', 7: 'yellow', 8: 'red' };
    var rampOrder = ['blue', 'teal', 'green', 'brown', 'ash', 'steel', 'orange', 'yellow', 'red'];
    for (var j = 0; j < H; j++) {
      var lat = (j / (H - 1) - 0.5) * Math.PI;
      var cl = Math.cos(lat), sl = Math.sin(lat);
      for (var i = 0; i < W; i++) {
        var lon = (i / W) * Math.PI * 2;
        var nx = Math.sin(lon) * cl, ny = sl, nz = Math.cos(lon) * cl;
        var h = fbm(nx * 2.1 + ox, ny * 2.1 + oy, nz * 2.1 + oz, 3);
        var al = Math.abs(ny);
        var ramp, shade;
        if (kind === 3) {
          var band = Math.sin(ny * 11 + h * 3.0) * 0.5 + 0.5;
          ramp = band > 0.66 ? 'orange' : (band > 0.36 ? 'brown' : 'yellow');
          shade = 1.6 + band * 2.4;
        } else if (kind === 2) {
          ramp = h > 0.5 ? 'steel' : 'blue';
          shade = 1.8 + h * 2.2;
        } else if (kind === 1) {
          ramp = h > 0.58 ? 'brown' : (h > 0.44 ? 'orange' : 'yellow');
          shade = 1.6 + h * 2.4;
        } else if (kind === 4) {
          ramp = h > 0.6 ? 'red' : (h > 0.4 ? 'orange' : 'brown');
          shade = 1.3 + h * 2.8;
        } else {
          if (h < 0.46) { ramp = 'blue'; shade = 1.1 + h * 3.2; }
          else if (h < 0.51) { ramp = 'teal'; shade = 2.3; }
          else if (h < 0.63) { ramp = 'green'; shade = 1.7 + (h - 0.51) * 7; }
          else if (h < 0.71) { ramp = 'brown'; shade = 2.1 + (h - 0.63) * 6; }
          else { ramp = 'ash'; shade = 3.6; }
          if (al > 0.87) { ramp = 'steel'; shade = 3.3 + h; }
        }
        lum[j * W + i] = Math.max(0, Math.min(5, Math.round(shade)));
        rampIdx[j * W + i] = rampOrder.indexOf(ramp);
      }
    }
    return { W: W, H: H, lum: lum, ramp: rampIdx, ramps: rampOrder };
  }

  M.makePlanetFrames = function (size, seed, kind, frames) {
    var key = seed + ':' + kind + ':' + size;
    if (planetCache[key]) return planetCache[key];
    frames = frames || 8;
    var tex = buildPlanetTexture(seed, kind);
    var out = [];
    var R = size / 2;
    var sun = V.norm([-0.55, 0.42, 0.72]);
    var TW = tex.W, TH = tex.H;
    for (var fr = 0; fr < frames; fr++) {
      var fb = new E.FB(size, size);
      var rot = fr / frames * Math.PI * 2;
      for (var y = 0; y < size; y++) {
        var ny = -(y - R) / R;
        for (var x = 0; x < size; x++) {
          var nx = (x - R) / R;
          var d2 = nx * nx + ny * ny;
          if (d2 > 1) continue;
          var nz = Math.sqrt(1 - d2);
          var lam = nx * sun[0] + ny * sun[1] + nz * sun[2];
          if (lam <= -0.03) continue;
          var lat = Math.asin(Math.max(-1, Math.min(1, ny)));
          var lon = Math.atan2(nx, nz) + rot;
          var u = ((lon / (Math.PI * 2)) % 1 + 1) % 1;
          var v = (lat / Math.PI + 0.5);
          var ti = (u * TW) | 0; if (ti >= TW) ti = TW - 1;
          var tj = (v * (TH - 1)) | 0; if (tj >= TH - 1) tj = TH - 1;
          var idx = tj * TW + ti;
          var shade = tex.lum[idx] * (0.20 + 0.80 * Math.max(0, lam)) + 0.15;
          var i0 = Math.max(0, Math.min(5, Math.floor(shade)));
          var i1 = Math.min(5, i0 + 1);
          var frac = shade - i0;
          var rid = E.RAMP_ID[tex.ramps[tex.ramp[idx]]];
          var bayer = E.BAYER[((y & 3) * 4 + (x & 3))];
          fb.buf[y * size + x] = E.PAL[rid * 6 + ((frac * 16 > bayer) ? i1 : i0)];
        }
      }
      out.push(fb);
    }
    planetCache[key] = out;
    return out;
  };

  // ------------------------------------------------------------- nebula
  var nebulaCache = {};
  M.makeNebula = function (seed, galaxy) {
    var key = seed + ':' + galaxy;
    if (nebulaCache[key]) return nebulaCache[key];
    var W = 320, H = 75;
    var fb = new E.FB(W, H);
    var rnd = E.rng(seed * 7 + galaxy * 131);
    var ox = rnd() * 40, oy = rnd() * 40;
    var ramp = galaxy % 2 === 0 ? 'purple' : 'blue';
    var ramp2 = galaxy % 3 === 0 ? 'teal' : 'blue';
    for (var y = 0; y < H; y++) {
      for (var x = 0; x < W; x++) {
        var n = fbm(x * 0.015 + ox, y * 0.032 + oy, 1.7, 2);
        var n2 = fbm(x * 0.028 + 9, y * 0.06 + 4, 4.1, 2);
        var d = n * 0.75 + n2 * 0.25;
        if (d < 0.46) continue;
        var t = (d - 0.46) / 0.22;
        var useRamp = n2 > 0.5 ? ramp : ramp2;
        var shade = Math.min(2.2, 0.35 + t * 1.9);
        var ci = shade | 0;
        var frac = shade - ci;
        var b = E.BAYER[((y & 3) * 4 + (x & 3))];
        if (frac * 16 < b) continue;              // sparse edges: starfield shows through
        fb.buf[y * W + x] = E.color(useRamp, ci);
      }
    }
    nebulaCache[key] = fb;
    return fb;
  };

  M.planetKindFor = function (seed) {
    var r = E.rng(seed)();
    if (r < 0.45) return 0;
    if (r < 0.62) return 1;
    if (r < 0.74) return 2;
    if (r < 0.92) return 3;
    return 4;
  };
})(typeof window !== 'undefined' ? window : globalThis);
