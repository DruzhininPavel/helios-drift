// Renders the model sheet and prints an ASCII luminance preview so the
// geometry can be sanity-checked without opening an image.
'use strict';
const { execFileSync } = require('child_process');
const fs = require('fs');
const path = require('path');

execFileSync(process.execPath, [path.join(__dirname, 'modelsheet.js')], { stdio: 'inherit' });

// re-render a single ship straight to ASCII
const vm = require('vm');
const sandbox = { console, Math, Date, JSON, isNaN, isFinite, parseInt, parseFloat };
sandbox.window = sandbox; sandbox.globalThis = sandbox;
vm.createContext(sandbox);
for (const f of ['render.js', 'models.js']) {
  vm.runInContext(fs.readFileSync(path.join(__dirname, '..', 'js', f), 'utf8'), sandbox, { filename: f });
}
const E = sandbox.E;
const chars = ' .:-=+*#%@';

function ascii(mesh, w = 72, h = 30, matR = Math.PI * 0.78, rotX = 0.34, distMul = 2.9) {
  const fb = new E.FB(w * 2, h * 2);
  fb.clear(0xff000000);
  const cam = new E.Camera(w * 2, h * 2, 60);
  const mat = E.M3.mul(E.M3.rotX(rotX), E.M3.rotY(matR));
  const d = Math.max(mesh.radius * distMul, 60);
  E.drawMesh(fb, cam, mesh, [0, 0, -d], mat, {});
  let out = '';
  for (let y = 0; y < h; y++) {
    let line = '';
    for (let x = 0; x < w; x++) {
      let sum = 0;
      for (let yy = 0; yy < 2; yy++) for (let xx = 0; xx < 2; xx++) {
        const c = fb.buf[(y * 2 + yy) * fb.w + (x * 2 + xx)];
        const r = c & 255, g = (c >> 8) & 255, b = (c >> 16) & 255;
        sum += (r * 0.3 + g * 0.59 + b * 0.11) / 255;
      }
      const lum = sum / 4;
      line += chars[Math.min(9, Math.max(0, Math.round(lum * 12)))];
    }
    out += line.replace(/\s+$/, '') + '\n';
  }
  return out;
}

const which = process.argv[2];
const list = which ? [which] : ['cobra3', 'sidewinder', 'krait', 'station', 'thargoid', 'python'];
for (const name of list) {
  console.log('=== ' + name + ' ===');
  console.log(ascii(E.models.meshes[name]));
}
