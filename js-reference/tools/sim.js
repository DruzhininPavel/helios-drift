/* Headless playtest: boots the real game code with a stub DOM, drives the real
 * input handlers, runs the real update/draw loop and writes screenshots. */
'use strict';
const vm = require('vm');
const fs = require('fs');
const path = require('path');
const { writePNG } = require('./png.js');

const ROOT = path.join(__dirname, '..');
const OUT = path.join(ROOT, 'build');
fs.mkdirSync(OUT, { recursive: true });

// ---------------------------------------------------------------- stub DOM
const listeners = { window: {}, canvas: {} };
function mkTarget(store) {
  return {
    addEventListener(type, fn) { (store[type] = store[type] || []).push(fn); },
    removeEventListener() { }
  };
}
const canvasStub = Object.assign(mkTarget(listeners.canvas), {
  width: 320, height: 224,
  style: {},
  getContext: () => ({
    imageSmoothingEnabled: false,
    createImageData: (w, h) => ({ width: w, height: h, data: new Uint8ClampedArray(w * h * 4) }),
    putImageData() { },
    getImageData: (x, y, w, h) => ({ width: w, height: h, data: new Uint8ClampedArray(w * h * 4) })
  }),
  getBoundingClientRect: () => ({ left: 0, top: 0, width: 320, height: 224 })
});
const documentStub = {
  readyState: 'complete',
  getElementById: (id) => (id === 'screen' ? canvasStub : null),
  addEventListener() { }
};
const sandbox = {
  console, Math, Date, JSON, isNaN, isFinite, parseInt, parseFloat, Uint8Array,
  Uint32Array, Uint8ClampedArray, Int32Array, Float32Array, Array, Object, String, Number,
  document: documentStub,
  performance: { now: () => Date.now() },
  requestAnimationFrame: () => 0,
  innerWidth: 1400, innerHeight: 900,
  addEventListener: mkTarget(listeners.window).addEventListener,
  localStorage: (() => {
    const m = {};
    return {
      getItem: (k) => (k in m ? m[k] : null),
      setItem: (k, v) => { m[k] = String(v); },
      removeItem: (k) => { delete m[k]; }
    };
  })()
};
sandbox.window = sandbox;
sandbox.globalThis = sandbox;
sandbox.setTimeout = (fn) => { };
vm.createContext(sandbox);

for (const f of ['render.js', 'models.js', 'galaxy.js', 'game.js', 'ui.js', 'main.js']) {
  vm.runInContext(fs.readFileSync(path.join(ROOT, 'js', f), 'utf8'), sandbox, { filename: f });
}
const E = sandbox.E, GAME = E.game;

const fb = E.getFB();

function keyDown(key) {
  (listeners.window.keydown || []).forEach(fn => fn({ key, shiftKey: false, repeat: false, preventDefault() { } }));
}
function keyUp(key) {
  (listeners.window.keyup || []).forEach(fn => fn({ key, preventDefault() { } }));
}
function shot(name) {
  writePNG(path.join(OUT, name + '.png'), fb.w, fb.h, Buffer.from(fb.bytes.buffer, 0, fb.bytes.length));
}
function run(frames, opts) {
  opts = opts || {};
  for (let i = 0; i < frames; i++) {
    E.frameStep(opts.dt || 1 / 60);
  }
}
function stats(label) {
  const p = GAME.player;
  const w = GAME.world;
  console.log(
    label,
    '| state=' + GAME.state,
    'screen=' + GAME.screen,
    p ? 'pos=' + p.pos.map(v => v.toFixed(0)).join(',') : '',
    p ? 'spd=' + p.speed.toFixed(0) : '',
    'ents=' + (w ? w.entities.length : '-'),
    'cash=' + GAME.cmd.cash.toFixed(1)
  );
}

const errors = [];
function stage(name, fn) {
  try { fn(); } catch (e) { console.log('FAIL at ' + name + ': ' + e.message); console.log(e.stack); errors.push(e); }
}

// ------------------------------------------------------------------ tests
try {

console.log('--- boot');
run(2);
shot('01-title');
stats('title');

console.log('--- new game');
keyDown('1'); keyUp('1');
run(2);
stats('docked');
shot('02-station');

console.log('--- market');
keyDown('1'); keyUp('1');      // menu -> market
run(2);
shot('03-market');
const foodBefore = GAME.cmd.cash;
keyDown('b'); keyUp('b');      // buy 1t food
keyDown('ArrowDown'); keyUp('ArrowDown');
keyDown('B'); keyUp('B');      // buy 10t textiles
run(2);
shot('04-market-bought');
console.log('bought, cash', GAME.cmd.cash.toFixed(1), 'from', foodBefore.toFixed(1),
  'cargo', JSON.stringify(GAME.cmd.cargo));

console.log('--- sell');
keyDown('ArrowUp'); keyUp('ArrowUp');
keyDown('s'); keyUp('s');
run(2);

console.log('--- equip');
keyDown('Escape'); keyUp('Escape');
keyDown('2'); keyUp('2');
run(2);
shot('05-equip');
keyDown('Enter'); keyUp('Enter');
run(2);
console.log('equip bought:', JSON.stringify(GAME.cmd.equipment), 'missiles', GAME.cmd.missiles);

console.log('--- chart');
keyDown('Escape'); keyUp('Escape');
keyDown('3'); keyUp('3');
run(2);
shot('06-chart');
keyDown('ArrowRight'); keyUp('ArrowRight');
keyDown('ArrowRight'); keyUp('ArrowRight');
run(2);
shot('07-chart-select');
console.log('chart cursor', E.ui.sel.chartCursor);

console.log('--- status');
keyDown('Escape'); keyUp('Escape');
keyDown('4'); keyUp('4');
run(2);
shot('08-status');

console.log('--- launch');
keyDown('Escape'); keyUp('Escape');
keyDown('5'); keyUp('5');
run(3);
stats('after launch');
shot('09-flight');

console.log('--- fly');
keyDown('ArrowUp'); run(30); keyUp('ArrowUp');
keyDown('x'); run(120); keyUp('x');
run(60);
stats('flying');
shot('10-flight-moving');

// spin to find the station
keyDown('ArrowLeft'); run(40); keyUp('ArrowLeft');
run(30);
shot('11-flight-turned');

// shoot
keyDown(' '); run(20); keyUp(' ');
run(5);
shot('12-flight-firing');

// docking computer
GAME.cmd.equipment.dockcomp = true;
GAME.player.pos = E.V.add(GAME.world.station.pos, [1200, 300, 6000]);
keyDown('l'); keyUp('l');
let docked = false;
for (let i = 0; i < 7000; i++) {
  E.frameStep(1 / 60);
  if (GAME.state === 'docked') { docked = true; console.log('docked after', (i / 60).toFixed(1), 's'); break; }
  if (GAME.state === 'dead') { console.log('died while docking:', GAME.dead.reason); break; }
}
stats('after autopilot');
shot('13-docking-result');

// hyperspace
if (GAME.state === 'docked') {
  console.log('--- hyperspace');
  keyDown('3'); keyUp('3');
  run(2);
  const gal = E.galaxy.get(GAME.cmd.galaxy);
  const cur = gal.systems[GAME.cmd.system];
  let target = null;
  for (const s of gal.systems) {
    const d = E.galaxy.distance(cur, s);
    if (s.id !== cur.id && d < 6.5) { target = s; break; }
  }
  E.ui.sel.chartCursor = target.id;
  GAME.cmd.fuel = 7;
  keyDown('h'); keyUp('h');
  run(4);
  console.log('hyper state', GAME.state, GAME.hyper && GAME.hyper.dist.toFixed(2));
  shot('14-hyperspace');
  for (let i = 0; i < 400; i++) {
    E.frameStep(1 / 60);
    if (GAME.state === 'flight') break;
  }
  stats('arrived');
  run(60);
  shot('15-arrival');
}

// combat sandbox: spawn pirates and let them shoot
console.log('--- combat');
if (GAME.state === 'flight') {
  GAME.spawnPirates(3);
  for (let i = 0; i < 900; i++) {
    E.frameStep(1 / 60);
    if (i % 120 === 0) GAME.in.fire = true;
    if (i === 600) shot('16-combat');
    if (GAME.state !== 'flight') break;
  }
  stats('combat end');
  shot('17-combat-end');
}

// HUD check at zero state
console.log('--- final');
run(2);
shot('18-final');
console.log('errors:', errors.length);
if (errors.length) console.log(errors[0].stack);

} catch (e) { console.log("EXCEPTION:", e.stack); }
