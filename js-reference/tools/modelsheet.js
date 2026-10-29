// Minimal PNG writer + model contact sheet renderer (headless verification)
'use strict';
const zlib = require('zlib');
const fs = require('fs');
const path = require('path');
const vm = require('vm');

function crc32(buf) {
  let c, crc = 0xffffffff;
  for (let n = 0; n < buf.length; n++) {
    c = (crc ^ buf[n]) & 0xff;
    for (let k = 0; k < 8; k++) c = c & 1 ? 0xedb88320 ^ (c >>> 1) : c >>> 1;
    crc = (crc >>> 8) ^ c;
  }
  return (crc ^ 0xffffffff) >>> 0;
}

function chunk(type, data) {
  const len = Buffer.alloc(4);
  len.writeUInt32BE(data.length, 0);
  const t = Buffer.from(type, 'ascii');
  const body = Buffer.concat([t, data]);
  const crc = Buffer.alloc(4);
  crc.writeUInt32BE(crc32(body), 0);
  return Buffer.concat([len, body, crc]);
}

function writePNG(file, w, h, rgba) {
  const raw = Buffer.alloc((w * 4 + 1) * h);
  for (let y = 0; y < h; y++) {
    raw[y * (w * 4 + 1)] = 0;
    rgba.copy(raw, y * (w * 4 + 1) + 1, y * w * 4, (y + 1) * w * 4);
  }
  const ihdr = Buffer.alloc(13);
  ihdr.writeUInt32BE(w, 0);
  ihdr.writeUInt32BE(h, 4);
  ihdr[8] = 8; ihdr[9] = 6; ihdr[10] = 0; ihdr[11] = 0; ihdr[12] = 0;
  const png = Buffer.concat([
    Buffer.from([0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a]),
    chunk('IHDR', ihdr),
    chunk('IDAT', zlib.deflateSync(raw, { level: 9 })),
    chunk('IEND', Buffer.alloc(0))
  ]);
  fs.writeFileSync(file, png);
}

// load game modules into a shared sandbox
const sandbox = { console, Math, Date, JSON, isNaN, isFinite, parseInt, parseFloat };
sandbox.window = sandbox;
sandbox.globalThis = sandbox;
vm.createContext(sandbox);
for (const f of ['render.js', 'models.js']) {
  const code = fs.readFileSync(path.join(__dirname, '..', 'js', f), 'utf8');
  vm.runInContext(code, sandbox, { filename: f });
}
const E = sandbox.E;

const CELL_W = 200, CELL_H = 160, COLS = 4;
const names = Object.keys(E.models.meshes);
const rows = Math.ceil(names.length / COLS);
const sheet = new E.FB(CELL_W * COLS, CELL_H * rows);
sheet.clear(E.rgb('#05070f'));

names.forEach((name, i) => {
  const mesh = E.models.meshes[name];
  const cell = new E.FB(CELL_W, CELL_H);
  cell.clear(E.rgb('#080c16'));
  const cam = new E.Camera(CELL_W, CELL_H, 60);
  cam.cx = CELL_W / 2; cam.cy = CELL_H / 2 - 6;
  // 3/4 front view
  const mat = E.M3.mul(E.M3.rotX(0.34), E.M3.rotY(Math.PI * 0.78));
  const dist = Math.max(mesh.radius * 3.4, 70);
  E.drawMesh(cell, cam, mesh, [0, 0, -dist], mat, { scale: 1 });
  // label
  cell.text(name, 4, CELL_H - 12, E.rgb('#9ef'), 1);
  sheet.blit(cell, (i % COLS) * CELL_W, Math.floor(i / COLS) * CELL_H, 1);
});

const outPath = path.join(__dirname, '..', 'build', 'models.png');
fs.mkdirSync(path.dirname(outPath), { recursive: true });
writePNG(outPath, sheet.w, sheet.h, Buffer.from(sheet.bytes.buffer, 0, sheet.bytes.length));
console.log('wrote', outPath, sheet.w + 'x' + sheet.h);

// also a big hero shot: cobra + station + planet
const hero = new E.FB(480, 300);
hero.clear(E.rgb('#05070f'));
const cam2 = new E.Camera(480, 300, 70);
cam2.cx = 240; cam2.cy = 140;
E.drawMesh(hero, cam2, E.models.meshes.station, [40, 20, -1300], E.M3.mul(E.M3.rotX(0.2), E.M3.rotY(0.5)), {});
E.drawMesh(hero, cam2, E.models.meshes.cobra3, [-120, -40, -420], E.M3.mul(E.M3.rotX(0.2), E.M3.rotY(Math.PI * 0.85)), {});
E.drawMesh(hero, cam2, E.models.meshes.viper, [180, 30, -700], E.M3.mul(E.M3.rotX(0.1), E.M3.rotY(Math.PI * 0.6)), {});
const pframes = E.models.makePlanetFrames(96, 1234, 0, 1);
hero.blit(pframes[0], 470, 20, 1);
writePNG(path.join(__dirname, '..', 'build', 'hero.png'), hero.w, hero.h, Buffer.from(hero.bytes.buffer, 0, hero.bytes.length));
console.log('wrote hero.png');

// hero stats
let hn = 0;
for (let i = 0; i < hero.buf.length; i++) if (hero.buf[i] !== E.rgb('#05070f')) hn++;
console.log('hero painted pixels', hn, 'of', hero.buf.length);
