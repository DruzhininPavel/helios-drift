// Mesh audit: closed-shell check, signed volume (outward normals), open edges.
'use strict';
const vm = require('vm');
const fs = require('fs');
const path = require('path');
const sandbox = { console, Math, Date, JSON, isNaN, isFinite, parseInt, parseFloat };
sandbox.window = sandbox; sandbox.globalThis = sandbox;
vm.createContext(sandbox);
for (const f of ['render.js', 'models.js']) {
  vm.runInContext(fs.readFileSync(path.join(__dirname, '..', 'js', f), 'utf8'), sandbox, { filename: f });
}
const E = sandbox.E;

function keyOf(p) { return p.map(v => Math.round(v * 100) / 100).join(','); }

function audit(name, mesh) {
  const parent = {};
  function find(a) { let r = a; while (parent[r] !== r) r = parent[r]; while (parent[a] !== r) { const n = parent[a]; parent[a] = r; a = n; } return r; }
  const edgeCount = {};
  for (const f of mesh.faces) {
    const k0 = keyOf(f.p[0]);
    if (!(k0 in parent)) parent[k0] = k0;
    for (let j = 1; j < f.p.length; j++) {
      const kj = keyOf(f.p[j]);
      if (!(kj in parent)) parent[kj] = kj;
      const ra = find(k0), rb = find(kj);
      if (ra !== rb) parent[rb] = ra;
    }
    for (let j = 0; j < f.p.length; j++) {
      const a = keyOf(f.p[j]), b = keyOf(f.p[(j + 1) % f.p.length]);
      const ek = a < b ? a + '|' + b : b + '|' + a;
      edgeCount[ek] = (edgeCount[ek] || 0) + 1;
    }
  }
  const shells = {};
  for (const f of mesh.faces) {
    const r = find(keyOf(f.p[0]));
    const s = shells[r] || (shells[r] = { faces: 0, vol: 0, verts: new Set() });
    s.faces++;
    let vol = 0;
    for (let j = 1; j + 1 < f.p.length; j++) {
      const a = f.p[0], b = f.p[j], c = f.p[j + 1];
      vol += (a[0] * (b[1] * c[2] - b[2] * c[1]) - a[1] * (b[0] * c[2] - b[2] * c[0]) + a[2] * (b[0] * c[1] - b[1] * c[0])) / 6;
    }
    s.vol += vol;
    for (const p of f.p) s.verts.add(keyOf(p));
  }
  let open = 0, nonManifold = 0;
  for (const ek in edgeCount) {
    if (edgeCount[ek] === 1) open++;
    else if (edgeCount[ek] > 2) nonManifold++;
  }
  const bad = Object.values(shells).filter(s => s.vol < -1e-3);
  const flat = Object.values(shells).filter(s => Math.abs(s.vol) <= 1e-3);
  console.log(
    name.padEnd(11),
    'faces=' + String(mesh.faces.length).padStart(3),
    'shells=' + String(Object.keys(shells).length).padStart(2),
    'openEdges=' + String(open).padStart(3),
    'nonManifold=' + String(nonManifold).padStart(3),
    'invertedShells=' + bad.length,
    'flatShells=' + flat.length,
    'badVol=' + bad.map(s => s.vol.toFixed(1)).join(',')
  );
  return { open, bad: bad.length };
}

let problems = 0;
for (const [name, mesh] of Object.entries(E.models.meshes)) {
  const r = audit(name, mesh);
  problems += r.open + r.bad;
}
console.log('\ntotal problems:', problems);
