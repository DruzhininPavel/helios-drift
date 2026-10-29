'use strict';
const vm = require('vm');
const fs = require('fs');
const path = require('path');
const { writePNG } = require('./png.js');
const sandbox = { console, Math, Date, JSON, isNaN, isFinite, parseInt, parseFloat, Uint8Array, Uint32Array, Uint8ClampedArray };
sandbox.window = sandbox; sandbox.globalThis = sandbox;
vm.createContext(sandbox);
for (const f of ['render.js', 'models.js', 'galaxy.js', 'game.js', 'ui.js']) {
  vm.runInContext(fs.readFileSync(path.join(__dirname, '..', 'js', f), 'utf8'), sandbox, { filename: f });
}
const E = sandbox.E, GAME = E.game, V = E.V, M3 = E.M3, UI = E.ui;
const fb = new E.FB(320, 224);
function shot(n) {
  UI.draw(fb, 1 / 60);
  writePNG(path.join(__dirname, '..', 'build', n + '.png'), 320, 224, Buffer.from(fb.bytes.buffer, 0, fb.bytes.length));
  console.log('wrote', n);
}
function lookAt(from, target) {
  GAME.player.pos = from.slice();
  GAME.player.mat = GAME.lookMatrix(V.norm(V.sub(target, from)));
}

GAME.newGame('JAMESON');
GAME.launch();
const p = GAME.player, st = GAME.world.station;
const normal = [st.mat[0][2], st.mat[1][2], st.mat[2][2]];
GAME.messages.length = 0;
GAME.world.entities.length = 0;
GAME.time = 8;
p.speed = 180;

// 1. station from 1.6 km with traffic
lookAt(V.add(st.pos, V.mul(normal, 1600)), st.pos);
GAME.spawnTraffic(2);
GAME.world.entities.forEach((e, i) => {
  if (e.kind !== 'ship') return;
  e.pos = V.add(p.pos, V.mul(V.norm([0.45 - i * 0.9, 0.15, -1]), 900 + i * 500));
  e.mat = GAME.lookMatrix(V.norm(V.sub(p.pos, e.pos)));
  e.speed = 120;
});
GAME.message('DOCKING CLEARANCE GRANTED', 'good');
shot('30-approach');

// 2. station close, slot lined up
const slotPoint = V.add(st.pos, V.mul(normal, 520));
lookAt(slotPoint, st.pos);
GAME.world.entities.length = 0;
GAME.message('ALIGN WITH THE SLOT', 'info');
shot('31-slot');

// 3. combat
lookAt(V.add(st.pos, V.mul(normal, 2200)), st.pos);
GAME.world.entities.length = 0;
GAME.spawnPirates(3);
GAME.world.entities.forEach((e, i) => {
  if (e.kind !== 'ship') return;
  const off = [[-260, 40, -900], [180, -60, -1300], [60, 120, -700]][i % 3];
  e.pos = V.add(p.pos, off);
  e.mat = GAME.lookMatrix(V.norm(V.sub(p.pos, e.pos)));
  e.speed = 260;
  e.hostile = true;
});
GAME.messages.length = 0;
GAME.message('HOSTILE CONTACT', 'bad');
shot('32-combat');

// 4. planet
const pl = GAME.world.planet;
lookAt(V.add(pl.pos, [0, 0, pl.radius * 3.4]), pl.pos);
GAME.world.entities.length = 0;
shot('33-planet');

// 5. sun with scoop
const sun = GAME.world.sun;
lookAt(V.add(sun.pos, [sun.radius * 2.2, 0, sun.radius * 2.2]), sun.pos);
shot('34-sun');

// 6. docked screens
GAME.state = 'docked';
GAME.dockedAt();
GAME.screen = 'market';
shot('35-market');
GAME.screen = 'chart';
UI.sel.chartCursor = GAME.cmd.system;
shot('36-chart');
GAME.screen = 'launch';
shot('37-station-menu');

// 7. title
GAME.state = 'title';
GAME.time = 3.5;
shot('38-title');
