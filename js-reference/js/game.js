/* ELITE 16 - game core
 * Flight model, combat, docking, hyperspace, ship AI, commander state.
 * Units: 1 unit ~ 1 metre. Speed in units/second. Time in seconds.
 */
(function (global) {
  'use strict';
  var E = global.E;
  var V = E.V, M3 = E.M3, G = E.galaxy;

  var GAME = E.game = {};

  // ------------------------------------------------------------- constants
  var C = GAME.C = {
    viewW: 320,
    viewH: 150,
    maxSpeed: 420,          // cobra mk III
    pitchRate: 1.55,
    rollRate: 2.35,
    stationS: (E.models && E.models.CORIOLIS) ? E.models.CORIOLIS.S : 220,
    slotW: (E.models && E.models.CORIOLIS) ? E.models.CORIOLIS.slotW : 0.31,
    slotH: (E.models && E.models.CORIOLIS) ? E.models.CORIOLIS.slotH : 0.072,
    laserRange: [2600, 3400, 4600, 2000],
    laserDamage: [9, 16, 26, 7],
    laserHeat: [0.055, 0.10, 0.13, 0.05],
    laserCooldown: [0.14, 0.075, 0.06, 0.2],
    laserName: ['PULSE LASER', 'BEAM LASER', 'MILITARY LASER', 'MINING LASER'],
    maxEnergy: 100,
    maxShield: 100,
    energyRecharge: 1.7,
    shieldRecharge: 1.9,
    radarRange: 26000,
    cargoBase: 20,
    fuelMax: 7.0,
    missilesMax: 4
  };

  var SHIPS = GAME.SHIPS = {
    sidewinder: { mesh: 'sidewinder', radius: 26, speed: 380, energy: 34, shield: 0, bounty: 5, cargo: 2, laser: 1, fireRate: 0.55, turn: 1.7, role: 'pirate' },
    krait: { mesh: 'krait', radius: 27, speed: 340, energy: 60, shield: 0, bounty: 10, cargo: 3, laser: 1, fireRate: 0.5, turn: 1.6, role: 'pirate' },
    mamba: { mesh: 'mamba', radius: 28, speed: 420, energy: 48, shield: 10, bounty: 15, cargo: 3, laser: 1, fireRate: 0.45, turn: 1.8, role: 'pirate' },
    adder: { mesh: 'adder', radius: 24, speed: 260, energy: 40, shield: 0, bounty: 8, cargo: 4, laser: 0, fireRate: 0.8, turn: 1.3, role: 'trader' },
    cobra3: { mesh: 'cobra3', radius: 34, speed: 350, energy: 80, shield: 20, bounty: 20, cargo: 6, laser: 1, fireRate: 0.5, turn: 1.6, role: 'pirate' },
    python: { mesh: 'python', radius: 46, speed: 200, energy: 140, shield: 30, bounty: 30, cargo: 12, laser: 1, fireRate: 0.7, turn: 1.0, role: 'trader' },
    viper: { mesh: 'viper', radius: 28, speed: 480, energy: 70, shield: 25, bounty: 15, cargo: 2, laser: 1, fireRate: 0.4, turn: 2.0, role: 'police' },
    asp: { mesh: 'asp', radius: 32, speed: 400, energy: 100, shield: 25, bounty: 40, cargo: 4, laser: 2, fireRate: 0.4, turn: 1.8, role: 'pirate' },
    moray: { mesh: 'moray', radius: 24, speed: 320, energy: 44, shield: 0, bounty: 12, cargo: 3, laser: 1, fireRate: 0.6, turn: 1.5, role: 'pirate' },
    ferdelance: { mesh: 'ferdelance', radius: 30, speed: 480, energy: 100, shield: 30, bounty: 40, cargo: 4, laser: 2, fireRate: 0.35, turn: 2.1, role: 'pirate' },
    thargoid: { mesh: 'thargoid', radius: 36, speed: 400, energy: 260, shield: 40, bounty: 100, cargo: 0, laser: 2, fireRate: 0.35, turn: 1.6, role: 'thargoid' },
    shuttle: { mesh: 'shuttle', radius: 22, speed: 200, energy: 26, shield: 0, bounty: 0, cargo: 2, laser: 0, fireRate: 1.0, turn: 1.0, role: 'trader' }
  };

  GAME.NEW_COMMANDER = function () {
    return {
      name: 'JAMESON',
      cash: 100,
      fuel: 7.0,
      missiles: 3,
      cargo: {},              // commodity index -> tonnes
      cargoCap: 20,
      laser: 0,               // index into laser arrays
      rearLaser: false,
      kills: 0,
      legal: 0,               // bounty in credits; >0 = offender/fugitive
      galaxy: 0,
      system: 0,
      targetSystem: 0,
      equipment: {
        ecm: false, fuelscoop: false, largebay: false, bomb: false,
        escape: false, extraenergy: false, dockcomp: false, galdrive: false,
        mining: false
      },
      bombs: 0,
      deaths: 0
    };
  };

  // --------------------------------------------------------------- helpers
  function rnd() { return Math.random(); }
  function randRange(a, b) { return a + rnd() * (b - a); }
  function randDir() {
    var z = randRange(-1, 1), t = randRange(0, Math.PI * 2), r = Math.sqrt(1 - z * z);
    return [Math.cos(t) * r, Math.sin(t) * r, z];
  }

  GAME.rnd = rnd;

  // --------------------------------------------------------------- state
  GAME.state = 'title';
  GAME.cmd = GAME.NEW_COMMANDER();
  GAME.world = null;
  GAME.messages = [];
  GAME.time = 0;
  GAME.in = { pitch: 0, roll: 0, accel: 0, fire: false, mouse: null };
  GAME.screen = 'flight';   // flight | market | chart | equip | status | launch
  GAME.market = [];
  GAME.dockSalt = 1;
  GAME.autopilot = false;
  GAME.autopilotPath = null;
  GAME.paused = false;
  GAME.overlay = null;
  GAME.shake = 0;
  GAME.beamFlash = 0;
  GAME.ecmFlash = 0;
  GAME.bombFlash = 0;
  GAME.hyper = null;
  GAME.dead = null;
  GAME.arrivalPhase = 0;
  GAME.encounterTimer = 30;
  GAME.targetLocked = null;

  GAME.message = function (text, kind) {
    GAME.messages.push({ text: text, t: 0, life: 4.5, kind: kind || 'info' });
    if (GAME.messages.length > 4) GAME.messages.shift();
  };

  GAME.emit = function (name) {
    if (GAME.onEvent) GAME.onEvent(name);
  };

  // ------------------------------------------------------------ world setup
  GAME.enterSystem = function (systemIndex, opts) {
    opts = opts || {};
    var gal = G.get(GAME.cmd.galaxy);
    var sys = gal.systems[systemIndex % 256];
    GAME.cmd.system = sys.id;
    var Rp = 4200 + (sys.seed % 2600);
    var planetPos = [0, 0, -Rp * 4];
    var stationPos = [Rp * 0.6, Rp * 0.18, -Rp * 1.2];

    var spawnDist = opts.spawnDist || 9000;
    var spawnPos = [0, 0, 0];
    var spawnMat = M3.id();

    if (opts.mode === 'launch') {
      // just outside the slot, moving away from the station
      var sm = M3.rotY(randRange(0, Math.PI * 2));
      var normal = [sm[0][2], sm[1][2], sm[2][2]];
      spawnPos = V.add(stationPos, V.mul(normal, C.stationS + 420));
      spawnMat = lookMatrix(normal);
    } else {
      // witchpoint: station dead ahead, planet behind it
      spawnPos = V.add(stationPos, [0, 0, spawnDist]);
      spawnMat = lookMatrix(V.norm(V.sub(stationPos, spawnPos)));
    }

    var world = {
      sys: sys,
      planet: { kind: 'planet', pos: planetPos, radius: Rp, seed: sys.seed, name: sys.name },
      sun: { kind: 'sun', pos: [Rp * 26, Rp * 6, -Rp * 30], radius: Rp * 7.5 },
      entities: [],
      stars: makeStars(),
      station: {
        kind: 'station',
        pos: stationPos,
        mat: M3.rotY(randRange(0, 6.28)),
        radius: C.stationS,
        hp: 100000
      },
      dust: makeDust()
    };
    GAME.world = world;

    var p = GAME.player = {
      pos: spawnPos,
      mat: spawnMat,
      speed: opts.mode === 'launch' ? 90 : 260,
      energy: C.maxEnergy,
      shieldF: C.maxShield,
      shieldA: C.maxShield,
      laserTemp: 0,
      fireCooldown: 0,
      missileTarget: null,
      invulnerable: opts.mode === 'launch' ? 3 : 1.5
    };

    GAME.state = 'flight';
    GAME.screen = 'flight';
    GAME.autopilot = false;
    GAME.encounterTimer = opts.mode === 'launch' ? 25 : 12;
    GAME.messages.length = 0;
    GAME.targetLocked = null;

    if (opts.mode !== 'launch') {
      GAME.message('ARRIVED: ' + sys.name, 'good');
      if (opts.misjump) {
        GAME.message('HYPERSPACE MISJUMP!', 'bad');
        spawnThargoids();
      } else if (rnd() < 0.55) {
        GAME.spawnTraffic(1 + Math.floor(rnd() * 2));
      }
      if (rnd() < dangerLevel(sys) * 0.6) {
        GAME.spawnPirates(1 + Math.floor(rnd() * 2));
      }
    } else {
      GAME.message('LAUNCHED FROM ' + sys.name + ' STATION', 'info');
    }
  };

  function makeStars() {
    var stars = [];
    for (var i = 0; i < 90; i++) {
      var d = randDir();
      stars.push({ dir: d, bright: rnd(), size: rnd() < 0.12 ? 2 : 1 });
    }
    return stars;
  }
  function makeDust() {
    var dust = [];
    for (var i = 0; i < 40; i++) {
      dust.push({ pos: [randRange(-1500, 1500), randRange(-1000, 1000), randRange(-3000, -200)], });
    }
    return dust;
  }

  function lookMatrix(forward) {
    // columns are the ship's local axes in world space (matches M3.apply =
    // local -> world, M3.applyT = world -> local)
    var zAxis = V.mul(V.norm(forward), -1);
    var up = [0, 1, 0];
    if (Math.abs(V.dot(zAxis, up)) > 0.98) up = [1, 0, 0];
    var x = V.norm(V.cross(up, zAxis));
    var y = V.cross(zAxis, x);
    return [
      [x[0], y[0], zAxis[0]],
      [x[1], y[1], zAxis[1]],
      [x[2], y[2], zAxis[2]]
    ];
  }
  GAME.lookMatrix = lookMatrix;

  function dangerLevel(sys) {
    return 0.15 + (7 - sys.government) * 0.09 + (sys.tech < 5 ? 0.12 : 0);
  }

  // -------------------------------------------------------------- spawning
  function baseShip(typeKey, pos, behaviour) {
    var def = SHIPS[typeKey];
    var sys = GAME.world.sys;
    var e = {
      kind: 'ship',
      type: typeKey,
      mesh: E.models.meshes[def.mesh],
      pos: pos.slice(),
      mat: lookMatrix(randDir()),
      speed: def.speed * 0.4,
      maxSpeed: def.speed,
      energy: def.energy,
      maxEnergy: def.energy,
      shield: def.shield || 0,
      maxShield: def.shield || 0,
      radius: def.radius,
      bounty: def.bounty,
      cargoType: def.cargo,
      behaviour: behaviour || def.role,
      laser: def.laser,
      fireRate: def.fireRate,
      turn: def.turn,
      cooldown: randRange(0.2, 1.5),
      hostile: behaviour === 'pirate' || behaviour === 'thargoid' || behaviour === 'police',
      age: 0,
      canMissile: ['viper', 'asp', 'ferdelance', 'thargoid', 'cobra3', 'python'].indexOf(typeKey) >= 0,
      role: def.role
    };
    return e;
  }

  function spawnAtDistance(type, dist, behaviour) {
    var d = randDir();
    // bias spawns toward the front hemisphere
    d[2] = -Math.abs(d[2]) * 0.8 - 0.2;
    d = V.norm(d);
    var pos = V.add(GAME.player.pos, V.mul(d, dist));
    return baseShip(type, pos, behaviour);
  }

  GAME.spawnPirates = function (n) {
    var types = ['sidewinder', 'krait', 'mamba', 'adder', 'cobra3', 'moray', 'ferdelance', 'asp'];
    var sys = GAME.world.sys;
    for (var i = 0; i < n; i++) {
      var t = types[Math.floor(rnd() * types.length)];
      if (sys.government <= 1 && rnd() < 0.3) t = 'ferdelance';
      var s = spawnAtDistance(t, randRange(5500, 9000), 'pirate');
      GAME.world.entities.push(s);
    }
    GAME.message('HOSTILE CONTACT', 'bad');
  };

  GAME.spawnTraffic = function (n) {
    var types = ['adder', 'python', 'cobra3', 'shuttle', 'moray'];
    for (var i = 0; i < n; i++) {
      var t = types[Math.floor(rnd() * types.length)];
      var s = spawnAtDistance(t, randRange(4000, 11000), 'trader');
      s.hostile = false;
      GAME.world.entities.push(s);
    }
  };

  function spawnPolice(n) {
    for (var i = 0; i < n; i++) {
      var s = spawnAtDistance('viper', randRange(4000, 7000), 'police');
      GAME.world.entities.push(s);
    }
    GAME.message('POLICE VIPER ON INTERCEPT', 'bad');
  }

  function spawnThargoids() {
    for (var i = 0; i < 2; i++) {
      var s = spawnAtDistance('thargoid', randRange(3000, 6000), 'thargoid');
      GAME.world.entities.push(s);
    }
    GAME.message('THARGOID ATTACK', 'bad');
  }

  GAME.spawnCargo = function (pos, cargoType) {
    GAME.world.entities.push({
      kind: 'canister',
      pos: pos.slice(),
      mat: lookMatrix(randDir()),
      speed: randRange(10, 40),
      radius: 8,
      cargoType: cargoType,
      age: 0
    });
  };

  GAME.spawnAsteroid = function (dist) {
    var d = randDir(); d[2] = -Math.abs(d[2]) - 0.1;
    var pos = V.add(GAME.player.pos, V.mul(V.norm(d), dist));
    GAME.world.entities.push({
      kind: 'asteroid',
      pos: pos,
      mat: lookMatrix(randDir()),
      speed: randRange(20, 60),
      radius: 20,
      energy: 60,
      maxEnergy: 60,
      spin: [randRange(-0.6, 0.6), randRange(-0.6, 0.6), randRange(-0.6, 0.6)],
      mesh: rnd() < 0.5 ? E.models.meshes.asteroid : E.models.meshes.asteroid2,
      age: 0
    });
  };

  // ---------------------------------------------------------------- combat
  GAME.fireLaser = function () {
    if (GAME.state !== 'flight') return;
    var p = GAME.player;
    if (p.laserTemp > 0.98) return;
    var cmd = GAME.cmd;
    var dmg = C.laserDamage[cmd.laser];
    var range = C.laserRange[cmd.laser];
    var heat = C.laserHeat[cmd.laser];
    p.laserTemp = Math.min(1, p.laserTemp + heat);

    // find the closest target inside the firing cone
    var best = null, bestDepth = 1e9;
    var fwd = [-p.mat[0][2], -p.mat[1][2], -p.mat[2][2]];
    var ents = GAME.world.entities;
    for (var i = 0; i < ents.length; i++) {
      var e = ents[i];
      if (e.kind === 'planet' || e.kind === 'sun') continue;
      var rel = V.sub(e.pos, p.pos);
      var depth = V.dot(rel, fwd);
      if (depth < 30 || depth > range) continue;
      var lateral = V.len(V.sub(rel, V.mul(fwd, depth)));
      var angularRadius = (e.radius || 20) / depth;
      var spread = Math.max(angularRadius * 1.25, 0.012);
      if (lateral < depth * spread && depth < bestDepth) {
        best = e; bestDepth = depth;
      }
    }
    if (!best) return;
    GAME.hitEntity(best, dmg, 'laser');
  };

  GAME.hitEntity = function (e, dmg, source) {
    if (e.kind === 'missile') { removeEntity(e); return; }
    if (e.kind === 'station') return; // station is indestructible
    if (e.shield > 0) {
      e.shield -= dmg;
      if (e.shield < 0) { e.energy += e.shield; e.shield = 0; }
    } else {
      e.energy -= dmg;
    }
    if (e.kind === 'canister') { e.energy = 0; }
    if (e.energy <= 0) {
      GAME.destroyEntity(e, source);
    } else if (e.hostile && e.kind === 'ship') {
      // getting shot at makes traders run / police hostile
      if (e.behaviour === 'trader') { e.behaviour = 'flee'; e.hostile = false; }
    }
  };

  GAME.destroyEntity = function (e, source) {
    var cmd = GAME.cmd;
    if (e.kind === 'ship') {
      var def = SHIPS[e.type];
      // loot
      if (e.cargoType > 0 && rnd() < 0.6) {
        for (var k = 0; k < 1 + Math.floor(rnd() * Math.min(3, e.cargoType)); k++) {
          GAME.spawnCargo(V.add(e.pos, V.mul(randDir(), 20)), Math.floor(rnd() * 12));
        }
      }
      if (e.behaviour === 'pirate' || e.behaviour === 'thargoid') {
        cmd.cash += def.bounty;
        cmd.kills++;
        if (def.bounty > 0) GAME.message('BOUNTY ' + def.bounty + ' CR', 'good');
      } else if (e.behaviour === 'police') {
        cmd.legal += 50;
        GAME.message('POLICE SHIP DESTROYED  +50 CR BOUNTY ON YOU', 'bad');
      } else {
        cmd.legal += 20;
        GAME.message('TRADER DESTROYED  LEGAL RECORD WORSENED', 'bad');
      }
      GAME.addExplosion(e.pos, e.radius);
      GAME.emit('explode');
    } else if (e.kind === 'asteroid') {
      GAME.addExplosion(e.pos, e.radius);
      if (cmd.equipment.mining || rnd() < 0.4) {
        GAME.spawnCargo(V.add(e.pos, V.mul(randDir(), 30)), 12); // minerals
      }
    } else if (e.kind === 'canister') {
      GAME.addExplosion(e.pos, 10);
    }
    removeEntity(e);
  };

  function removeEntity(e) {
    var arr = GAME.world.entities;
    var i = arr.indexOf(e);
    if (i >= 0) arr.splice(i, 1);
  }

  GAME.addExplosion = function (pos, radius) {
    GAME.world.entities.push({
      kind: 'explosion',
      pos: pos.slice(),
      radius: radius,
      age: 0,
      life: 0.6 + radius / 90
    });
  };

  GAME.fireMissile = function () {
    var cmd = GAME.cmd, p = GAME.player;
    if (cmd.missiles <= 0) { GAME.message('NO MISSILES', 'bad'); return; }
    var target = GAME.missileTargetPick();
    if (!target) { GAME.message('NO TARGET LOCKED', 'bad'); return; }
    cmd.missiles--;
    var m = {
      kind: 'missile',
      pos: V.add(p.pos, [0, -8, 0]),
      mat: p.mat.slice(),
      speed: 300,
      maxSpeed: 1100,
      target: target,
      radius: 6,
      age: 0,
      life: 9,
      energy: 10,
      mesh: E.models.meshes.missile
    };
    GAME.world.entities.push(m);
    GAME.message('MISSILE AWAY', 'good');
  };

  GAME.enemyMissile = function (from) {
    var m = {
      kind: 'missile',
      pos: V.add(from.pos, V.mul([from.mat[0][2], from.mat[1][2], from.mat[2][2]], from.radius)),
      mat: from.mat.map(function (r) { return r.slice(); }),
      speed: 260,
      maxSpeed: 900,
      target: 'player',
      radius: 6,
      age: 0,
      life: 10,
      energy: 6,
      mesh: E.models.meshes.missile
    };
    GAME.world.entities.push(m);
    GAME.message('INCOMING MISSILE', 'bad');
    GAME.emit('warning');
  };

  GAME.missileTargetPick = function () {
    var p = GAME.player;
    var fwd = [-p.mat[0][2], -p.mat[1][2], -p.mat[2][2]];
    var best = null, bestScore = -1;
    var ents = GAME.world.entities;
    for (var i = 0; i < ents.length; i++) {
      var e = ents[i];
      if (e.kind !== 'ship') continue;
      var rel = V.sub(e.pos, p.pos);
      var depth = V.dot(rel, fwd);
      if (depth < 0 || depth > 30000) continue;
      var lateral = V.len(V.sub(rel, V.mul(fwd, depth)));
      var ang = lateral / Math.max(depth, 1);
      var score = Math.cos(Math.min(Math.PI / 2, ang)) * (e.hostile ? 1.5 : 1);
      if (score > bestScore) { bestScore = score; best = e; }
    }
    return best;
  };

  GAME.fireECM = function () {
    if (!GAME.cmd.equipment.ecm) { GAME.message('NO E.C.M. FITTED', 'bad'); return; }
    var killed = 0;
    GAME.world.entities.forEach(function (e) {
      if (e.kind === 'missile') { e.dead = true; killed++; }
    });
    GAME.world.entities = GAME.world.entities.filter(function (e) { return !e.dead; });
    GAME.player.energy = Math.max(0, GAME.player.energy - 4);
    GAME.message(killed ? 'E.C.M. ACTIVE - MISSILES DESTROYED' : 'E.C.M. ACTIVE', 'info');
    GAME.ecmFlash = 0.4;
  };

  GAME.fireBomb = function () {
    var cmd = GAME.cmd;
    if (cmd.bombs <= 0) { GAME.message('NO ENERGY BOMB', 'bad'); return; }
    cmd.bombs--;
    var killed = 0;
    GAME.world.entities.slice().forEach(function (e) {
      if (e.kind === 'ship' && e.type !== 'thargoid') {
        GAME.addExplosion(e.pos, e.radius);
        removeEntity(e);
        killed++;
      }
    });
    GAME.bombFlash = 1.2;
    GAME.message('ENERGY BOMB DETONATED  (' + killed + ' DESTROYED)', 'good');
  };

  // ----------------------------------------------------------- player damage
  GAME.damagePlayer = function (dmg, fromFront) {
    var p = GAME.player;
    if (p.invulnerable > 0) return;
    var sh = fromFront ? p.shieldF : p.shieldA;
    if (sh > 0) {
      sh -= dmg;
      if (sh < 0) { p.energy += sh; sh = 0; }
      if (fromFront) p.shieldF = sh; else p.shieldA = sh;
    } else {
      p.energy -= dmg;
    }
    if (p.energy <= 0) GAME.playerDeath();
  };

  GAME.playerDeath = function () {
    var cmd = GAME.cmd;
    if (cmd.equipment.escape && !GAME.escapeUsed) {
      GAME.escapeUsed = true;
      cmd.equipment.escape = false;
      cmd.cash = Math.max(0, Math.round(cmd.cash * 0.5));
      cmd.cargo = {};
      cmd.missiles = 0;
      GAME.message('ESCAPE CAPSULE LAUNCHED - RESCUED', 'bad');
      GAME.state = 'docked';
      GAME.screen = 'launch';
      GAME.dead = null;
      GAME.dockedAt();
      return;
    }
    GAME.state = 'dead';
    GAME.dead = { reason: GAME.lastDamageReason || 'SHIP DESTROYED', t: 0 };
    cmd.deaths++;
  };

  GAME.dockedAt = function () {
    GAME.screen = 'launch';
    GAME.market = G.makeMarket(GAME.world.sys, ++GAME.dockSalt);
    GAME.autopilot = false;
    GAME.autopilotPath = null;
    GAME.player.speed = 0;
  };

  // --------------------------------------------------------------- docking
  GAME.tryDock = function () {
    var w = GAME.world, p = GAME.player;
    if (!w) return false;
    var st = w.station;
    var rel = V.sub(p.pos, st.pos);
    var local = M3.applyT(st.mat, rel);
    var S = C.stationS;
    var slotHalfW = C.slotW * S, slotHalfH = C.slotH * S;

    // inside the slot mouth?
    if (local[2] > 0 && local[2] < S + 90) {
      if (Math.abs(local[0]) < slotHalfW + 10 && Math.abs(local[1]) < slotHalfH + 10) {
        var normal = [st.mat[0][2], st.mat[1][2], st.mat[2][2]];
        var fwd = [-p.mat[0][2], -p.mat[1][2], -p.mat[2][2]];
        var align = V.dot(fwd, V.mul(normal, -1));
        if (align > 0.80 && p.speed < C.maxSpeed * 0.8) {
          GAME.dock();
          return true;
        }
      }
    }
    return false;
  };

  GAME.dock = function () {
    GAME.state = 'docked';
    GAME.screen = 'launch';
    GAME.player.speed = 0;
    GAME.dockedAt();
    GAME.message('DOCKED', 'good');
    GAME.emit('dock');
    if (GAME.cmd.legal > 0) {
      // pay off bounty on docking? Elite: fines handled via legal status
    }
  };

  GAME.launch = function () {
    GAME.escapeUsed = false;
    GAME.enterSystem(GAME.cmd.system, { mode: 'launch' });
  };

  // ------------------------------------------------------------- hyperspace
  GAME.hyperspace = function (targetId) {
    var cmd = GAME.cmd;
    var gal = G.get(cmd.galaxy);
    var from = gal.systems[cmd.system];
    var to = gal.systems[targetId];
    var d = G.distance(from, to);
    if (to.id === from.id) { GAME.message('ALREADY IN SYSTEM', 'bad'); return false; }
    if (d > 7) { GAME.message('TARGET OUT OF RANGE (MAX 7.0 LY)', 'bad'); return false; }
    if (cmd.fuel < d) { GAME.message('INSUFFICIENT FUEL', 'bad'); return false; }
    cmd.fuel = Math.max(0, +(cmd.fuel - d).toFixed(1));
    cmd.targetSystem = to.id;
    GAME.state = 'hyper';
    GAME.hyper = { t: 0, phase: 0, target: to.id, misjump: rnd() < 0.03, dist: d };
    GAME.message('HYPERSPACE TO ' + to.name, 'info');
    return true;
  };

  GAME.galacticJump = function () {
    if (!GAME.cmd.equipment.galdrive) { GAME.message('NO GALACTIC HYPERDRIVE', 'bad'); return; }
    GAME.cmd.galdriveUsed = true;
    GAME.cmd.galaxy = (GAME.cmd.galaxy + 1) % 8;
    GAME.cmd.system = Math.floor(rnd() * 256);
    GAME.message('GALACTIC HYPERSPACE - GALAXY ' + (GAME.cmd.galaxy + 1), 'good');
    GAME.enterSystem(GAME.cmd.system, { mode: 'arrive', spawnDist: 9000 });
  };

  // --------------------------------------------------------------- in-system
  GAME.jumpInSystem = function () {
    var cmd = GAME.cmd;
    if (GAME.state !== 'flight') return;
    if (cmd.fuel < 0.1) { GAME.message('NOT ENOUGH FUEL FOR JUMP', 'bad'); return; }
    var w = GAME.world;
    var st = w.station;
    var dist = V.len(V.sub(st.pos, GAME.player.pos));
    if (dist < 2500) { GAME.message('ALREADY AT THE STATION', 'bad'); return; }
    cmd.fuel = Math.max(0, +(cmd.fuel - 0.1).toFixed(1));
    // arrive a short way out, facing the station
    var dir = V.norm(V.sub(st.pos, GAME.player.pos));
    var drop = V.sub(st.pos, V.mul(dir, 3000));
    GAME.player.pos = drop;
    GAME.player.mat = lookMatrix(dir);
    GAME.player.speed = 200;
    GAME.message('IN-SYSTEM JUMP', 'info');
    if (rnd() < dangerLevel(w.sys) * 0.5) GAME.spawnPirates(1 + Math.floor(rnd() * 2));
  };

  // ------------------------------------------------------------------ update
  GAME.update = function (dt) {
    if (GAME.state === 'dead') { GAME.dead.t += dt; return; }
    if (GAME.state === 'docked') { updateDocked(dt); return; }
    if (GAME.state === 'hyper') { updateHyper(dt); return; }
    if (GAME.state !== 'flight') return;

    GAME.time += dt;
    var p = GAME.player;
    var cmd = GAME.cmd;
    var w = GAME.world;

    // ---- flight controls / docking computer
    if (GAME.autopilot && GAME.autopilotPath) {
      runAutopilot(dt);
    } else {
      var pitch = GAME.in.pitch, roll = GAME.in.roll;
      if (pitch) p.mat = M3.mul(p.mat, M3.rotX(pitch * C.pitchRate * dt));
      if (roll) p.mat = M3.mul(p.mat, M3.rotZ(-roll * C.rollRate * dt));
      if (pitch || roll) p.mat = M3.orthonormalize(p.mat);

      var accel = GAME.in.accel;
      if (accel > 0) p.speed = Math.min(C.maxSpeed, p.speed + 260 * dt);
      else if (accel < 0) p.speed = Math.max(0, p.speed - 300 * dt);
      p.pos = V.add(p.pos, V.mul([-p.mat[0][2], -p.mat[1][2], -p.mat[2][2]], p.speed * dt));
    }

    // ---- energy / shields / heat
    var recharge = C.energyRecharge * (cmd.equipment.extraenergy ? 1.6 : 1);
    p.energy = Math.min(C.maxEnergy, p.energy + recharge * dt);
    p.shieldF = Math.min(C.maxShield, p.shieldF + C.shieldRecharge * dt);
    p.shieldA = Math.min(C.maxShield, p.shieldA + C.shieldRecharge * dt);
    p.laserTemp = Math.max(0, p.laserTemp - 0.35 * dt);
    if (p.invulnerable > 0) p.invulnerable -= dt;

    if (GAME.in.fire && !GAME.autopilot && p.fireCooldown <= 0) {
      p.fireCooldown = C.laserCooldown[cmd.laser];
      GAME.fireLaser();
      GAME.emit('laser');
    }
    if (p.fireCooldown > 0) p.fireCooldown -= dt;

    // ---- planet / sun
    var dPlanet = V.len(V.sub(p.pos, w.planet.pos));
    if (dPlanet < w.planet.radius * 1.02) {
      GAME.lastDamageReason = 'BURNED UP IN THE ATMOSPHERE';
      GAME.playerDeath();
    } else if (dPlanet < w.planet.radius * 1.18) {
      GAME.message('WARNING: ATMOSPHERE', 'bad');
      GAME.damagePlayer(14 * dt, true);
    }
    var dSun = V.len(V.sub(p.pos, w.sun.pos));
    if (dSun < w.sun.radius * 1.05) {
      GAME.lastDamageReason = 'CONSUMED BY THE STAR';
      GAME.playerDeath();
    }
    // fuel scooping near the sun
    if (cmd.equipment.fuelscoop && dSun < w.sun.radius * 2.2 && cmd.fuel < C.fuelMax) {
      cmd.fuel = Math.min(C.fuelMax, +(cmd.fuel + 0.35 * dt).toFixed(2));
      if (rnd() < dt * 0.8) GAME.message('FUEL SCOOPING', 'good');
    }

    // ---- station
    var st = w.station;
    st.mat = M3.mul(st.mat, M3.rotY(0.05 * dt));   // slow rotation
    var dst = V.len(V.sub(p.pos, st.pos));
    if (dst < C.stationS * 1.06) {
      if (!GAME.tryDock()) {
        // did we clip the hull?
        var local = M3.applyT(st.mat, V.sub(p.pos, st.pos));
        var outsideSlot = Math.abs(local[0]) > C.slotW * C.stationS + 10 || Math.abs(local[1]) > C.slotH * C.stationS + 10;
        if (dst < C.stationS * 1.02 && outsideSlot) {
          GAME.lastDamageReason = 'CRASHED INTO THE STATION';
          GAME.playerDeath();
        }
      }
    } else {
      GAME.tryDock();
    }

    // ---- entities
    updateEntities(dt);

    // ---- encounters
    GAME.encounterTimer -= dt;
    if (GAME.encounterTimer <= 0) {
      GAME.encounterTimer = randRange(25, 55);
      var danger = dangerLevel(w.sys);
      if (rnd() < danger * 0.7 && countType('pirate') < 4) {
        GAME.spawnPirates(1 + Math.floor(rnd() * 2));
      } else if (cmd.legal >= 100 && countType('police') < 3 && rnd() < 0.5) {
        spawnPolice(1);
      } else if (rnd() < 0.25) {
        GAME.spawnTraffic(1);
      } else if (rnd() < 0.12 && countType('asteroid') < 6) {
        for (var a = 0; a < 2 + Math.floor(rnd() * 3); a++) GAME.spawnAsteroid(randRange(2000, 6000));
      }
    }

    // messages
    for (var i = GAME.messages.length - 1; i >= 0; i--) {
      GAME.messages[i].t += dt;
      if (GAME.messages[i].t > GAME.messages[i].life) GAME.messages.splice(i, 1);
    }
    if (GAME.ecmFlash > 0) GAME.ecmFlash -= dt;
    if (GAME.bombFlash > 0) GAME.bombFlash -= dt;
  };

  function countType(behaviour) {
    var n = 0;
    for (var i = 0; i < GAME.world.entities.length; i++) {
      var e = GAME.world.entities[i];
      if (e.kind === 'ship' && e.behaviour === behaviour) n++;
    }
    return n;
  }

  function updateDocked(dt) {
    GAME.time += dt;
  }

  function updateHyper(dt) {
    var h = GAME.hyper;
    h.t += dt;
    if (h.phase === 0 && h.t > 1.6) {
      h.phase = 1;
    }
    if (h.phase === 1 && h.t > 3.2) {
      GAME.hyper = null;
      GAME.enterSystem(h.target, { mode: 'arrive', misjump: h.misjump, spawnDist: 9000 });
    }
  }

  // ------------------------------------------------------------------- AI
  function shipTurnTowards(s, targetVec, dt) {
    var local = M3.applyT(s.mat, targetVec);
    var ax = Math.atan2(local[0], -local[2]);
    var ay = Math.atan2(local[1], Math.sqrt(local[0] * local[0] + local[2] * local[2]));
    var roll = Math.max(-1, Math.min(1, ax * 1.6));
    var pitch = Math.max(-1, Math.min(1, ay * 1.6));
    s.mat = M3.mul(s.mat, M3.rotZ(roll * 1.8 * dt));
    s.mat = M3.mul(s.mat, M3.rotX(pitch * 1.6 * dt));
    s.mat = M3.orthonormalize(s.mat);
    return { ax: ax, ay: ay };
  }

  function updateEntities(dt) {
    var p = GAME.player;
    var w = GAME.world;
    var ents = w.entities;
    for (var i = ents.length - 1; i >= 0; i--) {
      var e = ents[i];
      e.age = (e.age || 0) + dt;

      if (e.kind === 'explosion') {
        if (e.age > e.life) ents.splice(i, 1);
        continue;
      }

      var fwd = [-e.mat[0][2], -e.mat[1][2], -e.mat[2][2]];

      if (e.kind === 'missile') {
        updateMissile(e, dt);
        if (e.dead) ents.splice(i, 1);
        continue;
      }

      if (e.kind === 'canister') {
        e.pos = V.add(e.pos, V.mul(fwd, e.speed * dt));
        e.mat = M3.mul(e.mat, M3.rotY(0.4 * dt));
        var dc = V.len(V.sub(p.pos, e.pos));
        if (dc < 40) {
          if (GAME.cmd.equipment.fuelscoop) {
            GAME.scoopCargo(e);
            ents.splice(i, 1);
          } else if (dc < 22) {
            GAME.damagePlayer(2, true);
            ents.splice(i, 1);
            GAME.message('CARGO CANISTER HIT HULL', 'bad');
          }
        }
        if (e.age > 180) ents.splice(i, 1);
        continue;
      }

      if (e.kind === 'asteroid') {
        e.pos = V.add(e.pos, V.mul(fwd, e.speed * dt));
        e.mat = M3.mul(e.mat, M3.rotX(e.spin[0] * dt));
        e.mat = M3.mul(e.mat, M3.rotY(e.spin[1] * dt));
        e.mat = M3.mul(e.mat, M3.rotZ(e.spin[2] * dt));
        var da = V.len(V.sub(p.pos, e.pos));
        if (da < e.radius * 0.8) {
          GAME.lastDamageReason = 'COLLIDED WITH AN ASTEROID';
          GAME.playerDeath();
        }
        continue;
      }

      if (e.kind !== 'ship') continue;

      // ---------------- ship AI
      var toPlayer = V.sub(p.pos, e.pos);
      var distP = V.len(toPlayer);

      if (e.behaviour === 'trader') {
        // head for the station
        var toStation = V.sub(w.station.pos, e.pos);
        var ds = V.len(toStation);
        var aim = V.norm(toStation);
        var angle = shipTurnTowards(e, aim, dt);
        e.speed = e.maxSpeed * (ds > 1200 ? 0.85 : 0.4);
        if (ds < C.stationS * 2.2) {
          // docked / gone
          ents.splice(i, 1);
          continue;
        }
      } else if (e.behaviour === 'flee') {
        var away = V.norm(V.mul(V.sub(e.pos, p.pos), 1));
        shipTurnTowards(e, away, dt);
        e.speed = e.maxSpeed;
      } else if (e.behaviour === 'police') {
        var aimP = V.norm(toPlayer);
        var angP = shipTurnTowards(e, aimP, dt);
        e.speed = distP > 900 ? e.maxSpeed : e.maxSpeed * 0.7;
        e.cooldown -= dt;
        if (e.cooldown <= 0 && distP < 3500 && Math.abs(angP.ax) < 0.12 && Math.abs(angP.ay) < 0.12) {
          e.cooldown = e.fireRate * randRange(0.8, 1.5);
          if (rnd() < Math.max(0.25, 1.05 - distP / 5200)) GAME.damagePlayer(randRange(3, 6), true);
        }
      } else {
        // pirates and thargoids: attack
        var aimH = V.norm(toPlayer);
        var angH = shipTurnTowards(e, aimH, dt);
        var desired = e.maxSpeed;
        if (distP < 700) desired = e.maxSpeed * 0.45;
        // weave a bit
        e.speed += (desired - e.speed) * Math.min(1, dt * 1.5);
        e.cooldown -= dt;
        if (e.canMissile && e.cooldown <= 0 && distP > 900 && distP < 5000 &&
          Math.abs(angH.ax) < 0.3 && Math.abs(angH.ay) < 0.3 && rnd() < 0.25) {
          e.canMissile = false;
          GAME.enemyMissile(e);
          e.cooldown = 2.5;
        }
        if (e.cooldown <= 0 && distP < 4000 && Math.abs(angH.ax) < 0.14 && Math.abs(angH.ay) < 0.14) {
          e.cooldown = e.fireRate * randRange(0.7, 1.6);
          if (rnd() < Math.max(0.22, 1.0 - distP / 5600)) {
            var dmg = (e.type === 'thargoid' ? randRange(6, 11) : randRange(2, 6));
            GAME.damagePlayer(dmg, true);
            GAME.lastLastName = '';
            GAME.lastDamageReason = 'DESTROYED BY ' + (e.type === 'thargoid' ? 'THARGOID' : 'PIRATE');
            GAME.emit('hit');
            GAME.shake = 0.35;
          }
        }
        // thargoids launch drones / no missiles for simplicity
      }

      e.pos = V.add(e.pos, V.mul([-e.mat[0][2], -e.mat[1][2], -e.mat[2][2]], e.speed * dt));

      // avoid the station
      if (V.len(V.sub(e.pos, w.station.pos)) < C.stationS * 1.3 && e.speed > 100) {
        e.mat = M3.mul(e.mat, M3.rotX(0.5 * dt));
      }
    }
  }

  function updateMissile(m, dt) {
    if (m.target === 'player') {
      var relP = V.sub(GAME.player.pos, m.pos);
      var dP = V.len(relP);
      if (dP < 70) {
        GAME.damagePlayer(70, false);
        GAME.addExplosion(m.pos, 18);
        GAME.lastDamageReason = 'DESTROYED BY A MISSILE';
        GAME.emit('explode');
        m.dead = true;
        return;
      }
      var aimP = V.norm(relP);
      var locP = M3.applyT(m.mat, aimP);
      var axP = Math.atan2(locP[0], -locP[2]);
      var ayP = Math.atan2(locP[1], Math.sqrt(locP[0] * locP[0] + locP[2] * locP[2]));
      m.mat = M3.mul(m.mat, M3.rotZ(Math.max(-1, Math.min(1, axP * 2.4)) * 2.2 * dt));
      m.mat = M3.mul(m.mat, M3.rotX(Math.max(-1, Math.min(1, ayP * 2.4)) * 2.2 * dt));
      m.mat = M3.orthonormalize(m.mat);
      m.speed = Math.min(m.maxSpeed, m.speed + 700 * dt);
      m.pos = V.add(m.pos, V.mul([-m.mat[0][2], -m.mat[1][2], -m.mat[2][2]], m.speed * dt));
      m.age += dt;
      if (m.age > m.life) m.dead = true;
      return;
    }
    var t = m.target;
    if (!t || GAME.world.entities.indexOf(t) < 0) {
      // retarget
      t = m.target = GAME.missileTargetPick();
    }
    if (t) {
      var rel = V.sub(t.pos, m.pos);
      var d = V.len(rel);
      if (d < 80) {
        GAME.hitEntity(t, 120, 'missile');
        GAME.addExplosion(m.pos, 20);
        m.dead = true;
        return;
      }
      var aim = V.norm(rel);
      var local = M3.applyT(m.mat, aim);
      var ax = Math.atan2(local[0], -local[2]);
      var ay = Math.atan2(local[1], Math.sqrt(local[0] * local[0] + local[2] * local[2]));
      m.mat = M3.mul(m.mat, M3.rotZ(Math.max(-1, Math.min(1, -ax * 2.4)) * 2.2 * dt));
      m.mat = M3.mul(m.mat, M3.rotX(Math.max(-1, Math.min(1, ay * 2.4)) * 2.2 * dt));
      m.mat = M3.orthonormalize(m.mat);
    }
    m.speed = Math.min(m.maxSpeed, m.speed + 700 * dt);
    m.pos = V.add(m.pos, V.mul([-m.mat[0][2], -m.mat[1][2], -m.mat[2][2]], m.speed * dt));
    m.age += dt;
    if (m.age > m.life) m.dead = true;
    // the missile can hit the player if it was fired at a target near us - skip
  }

  GAME.scoopCargo = function (e) {
    var cmd = GAME.cmd;
    var used = cargoUsed(cmd);
    if (used >= cmd.cargoCap) { GAME.message('CARGO HOLD FULL', 'bad'); return; }
    var idx = e.cargoType;
    cmd.cargo[idx] = (cmd.cargo[idx] || 0) + 1;
    GAME.message('SCOOPED: ' + GAME.commoditiesName(idx), 'good');
  };

  function cargoUsed(cmd) {
    var n = 0;
    for (var k in cmd.cargo) n += cmd.cargo[k];
    return n;
  }
  GAME.cargoUsed = cargoUsed;

  // ---------------------------------------------------------------- autopilot
  // Two-stage approach: line up on the slot axis from 2.5km out, then fly
  // straight through the mouth while tracking the station's rotation.
  function autopilotState() {
    var p = GAME.player, w = GAME.world, st = w.station, S = C.stationS;
    var normal = [st.mat[0][2], st.mat[1][2], st.mat[2][2]];
    var rel = V.sub(p.pos, st.pos);
    var local = M3.applyT(st.mat, rel);
    var dist = V.len(rel);
    var lateral = Math.sqrt(local[0] * local[0] + local[1] * local[1]);
    var inLine = local[2] > S + 60 && lateral < 55;
    var gate = V.add(st.pos, V.mul(normal, S + 1400));
    var aim = inLine ? V.mul(normal, -1) : V.norm(V.sub(gate, p.pos));
    return { aim: aim, inLine: inLine, dist: dist, lateral: lateral, local: local, normal: normal, gate: gate };
  }

  function autopilotInput() {
    var s = autopilotState();
    var local = M3.applyT(GAME.player.mat, s.aim);
    var ax = Math.atan2(local[0], -local[2]);
    var ay = Math.atan2(local[1], Math.sqrt(local[0] * local[0] + local[2] * local[2]));
    return {
      pitch: Math.max(-1, Math.min(1, ay * 2.2)),
      roll: Math.max(-1, Math.min(1, ax * 2.2)),
      ax: ax, ay: ay, inLine: s.inLine, dist: s.dist
    };
  }

  function autopilotThrottle() {
    var p = GAME.player;
    var s = autopilotInput();
    if (s.inLine) return p.speed < 105 ? 1 : (p.speed > 145 ? -1 : 0);
    if (Math.abs(s.ax) > 0.5 || Math.abs(s.ay) > 0.5) return s.dist > 3000 ? 1 : -1;
    if (s.dist > 3000) return 1;
    return p.speed < 240 ? 1 : -1;
  }

  GAME.toggleAutopilot = function () {
    if (!GAME.cmd.equipment.dockcomp) {
      GAME.message('DOCKING COMPUTER NOT FITTED', 'bad');
      return;
    }
    if (GAME.autopilot) {
      GAME.autopilot = false;
      GAME.autopilotPath = null;
      GAME.message('MANUAL CONTROL', 'info');
      return;
    }
    var p = GAME.player;
    var s = autopilotState();
    if (s.dist > 3000 || s.local[2] < 0 || s.lateral > 300) {
      p.pos = s.gate;
      p.speed = 70;
      p.mat = lookMatrix(V.mul(s.normal, -1));
    }
    GAME.autopilot = true;
    GAME.autopilotPath = {
      t: 0,
      dur: 6.0,
      from: p.pos.slice(),
      fromMat: p.mat.map(function (r) { return r.slice(); })
    };
    GAME.message('DOCKING COMPUTER ON', 'info');
  };

  function smoothstep(t) { return t * t * (3 - 2 * t); }

  function runAutopilot(dt) {
    var p = GAME.player, st = GAME.world.station, S = C.stationS;
    var path = GAME.autopilotPath;
    path.t += dt / path.dur;
    var k = path.t >= 1 ? 1 : smoothstep(path.t);
    var normal = [st.mat[0][2], st.mat[1][2], st.mat[2][2]];
    var end = V.add(st.pos, V.mul(normal, S + 12));
    var old = p.pos;
    p.pos = [
      path.from[0] + (end[0] - path.from[0]) * k,
      path.from[1] + (end[1] - path.from[1]) * k,
      path.from[2] + (end[2] - path.from[2]) * k
    ];
    p.speed = V.len(V.sub(p.pos, old)) / Math.max(dt, 0.0001);
    var target = lookMatrix(V.mul(normal, -1));
    var m = [[0, 0, 0], [0, 0, 0], [0, 0, 0]];
    for (var i = 0; i < 3; i++) for (var j = 0; j < 3; j++) {
      m[i][j] = path.fromMat[i][j] * (1 - k) + target[i][j] * k;
    }
    p.mat = M3.orthonormalize(m);
    if (path.t >= 1) {
      GAME.autopilot = false;
      GAME.autopilotPath = null;
      GAME.dock();
    }
  }

  // --------------------------------------------------------------- buy/sell
  GAME.buy = function (idx, qty) {
    var cmd = GAME.cmd, m = GAME.market[idx];
    if (!m) return;
    var free = cmd.cargoCap - cargoUsed(cmd);
    qty = Math.min(qty, m.qty, free, Math.floor(cmd.cash / m.price));
    if (qty <= 0) { GAME.message('CANNOT BUY', 'bad'); return; }
    cmd.cash -= qty * m.price;
    cmd.cargo[idx] = (cmd.cargo[idx] || 0) + qty;
    m.qty -= qty;
    GAME.message('BOUGHT ' + qty + 'T ' + m.name, 'good');
  };

  GAME.sell = function (idx, qty) {
    var cmd = GAME.cmd, m = GAME.market[idx];
    if (!m) return;
    var have = cmd.cargo[idx] || 0;
    qty = Math.min(qty, have);
    if (qty <= 0) { GAME.message('NOTHING TO SELL', 'bad'); return; }
    cmd.cash += qty * m.sell;
    cmd.cargo[idx] = have - qty;
    m.qty += qty;
    GAME.message('SOLD ' + qty + 'T ' + m.name + ' FOR ' + (qty * m.sell) + ' CR', 'good');
  };

  GAME.refuel = function () {
    var cmd = GAME.cmd;
    var need = C.fuelMax - cmd.fuel;
    if (need <= 0.01) { GAME.message('TANKS ALREADY FULL', 'bad'); return; }
    var cost = Math.ceil(need * G.fuelPricePerLY);
    if (cmd.cash < cost) {
      // partial fill
      var canBuy = cmd.cash / G.fuelPricePerLY;
      if (canBuy < 0.1) { GAME.message('NOT ENOUGH CREDITS', 'bad'); return; }
      cmd.fuel = +(cmd.fuel + canBuy).toFixed(1);
      cmd.cash = 0;
      GAME.message('PARTIAL REFUEL', 'info');
      return;
    }
    cmd.cash -= cost;
    cmd.fuel = C.fuelMax;
    GAME.message('REFUELLED FOR ' + cost + ' CR', 'good');
  };

  GAME.buyEquipment = function (id) {
    var cmd = GAME.cmd;
    var item = null;
    for (var i = 0; i < G.EQUIPMENT.length; i++) if (G.EQUIPMENT[i].id === id) item = G.EQUIPMENT[i];
    if (!item) return;
    if (GAME.world && GAME.world.sys.tech < item.tech) { GAME.message('NOT AVAILABLE HERE', 'bad'); return; }
    if (cmd.cash < item.price) { GAME.message('NOT ENOUGH CREDITS', 'bad'); return; }

    if (id === 'missile') {
      if (cmd.missiles >= C.missilesMax) { GAME.message('MISSILE RACKS FULL', 'bad'); return; }
      cmd.missiles++;
    } else if (id === 'laserpulse' || id === 'laserbeam' || id === 'lasermilitary') {
      var lvl = id === 'laserpulse' ? 0 : (id === 'laserbeam' ? 1 : 2);
      if (cmd.laser >= lvl) { GAME.message('ALREADY FITTED', 'bad'); return; }
      cmd.laser = lvl;
    } else if (id === 'mininglaser') {
      cmd.equipment.mining = true;
      cmd.laser = 3;
    } else if (id === 'largebay') {
      if (cmd.equipment.largebay) { GAME.message('ALREADY FITTED', 'bad'); return; }
      cmd.equipment.largebay = true;
      cmd.cargoCap = 35;
    } else if (id === 'bomb') {
      cmd.bombs++;
      cmd.equipment.bomb = true;
    } else {
      if (cmd.equipment[id]) { GAME.message('ALREADY FITTED', 'bad'); return; }
      cmd.equipment[id] = true;
    }
    cmd.cash -= item.price;
    GAME.message('PURCHASED: ' + item.name, 'good');
  };

  // ------------------------------------------------------------------ save
  GAME.save = function () {
    try {
      var data = {
        version: 1,
        cmd: GAME.cmd,
        system: GAME.cmd.system,
        savedAt: Date.now()
      };
      global.localStorage.setItem('elite16.save', JSON.stringify(data));
      return true;
    } catch (err) {
      return false;
    }
  };

  GAME.load = function () {
    try {
      var raw = global.localStorage.getItem('elite16.save');
      if (!raw) return false;
      var data = JSON.parse(raw);
      GAME.cmd = data.cmd;
      GAME.enterSystem(GAME.cmd.system, { mode: 'arrive', spawnDist: 9000 });
      GAME.state = 'docked';
      GAME.screen = 'launch';
      GAME.dockedAt();
      return true;
    } catch (err) {
      return false;
    }
  };

  GAME.hasSave = function () {
    try { return !!global.localStorage.getItem('elite16.save'); } catch (e) { return false; }
  };

  GAME.commoditiesName = function (idx) {
    return G.COMMODITIES[idx] ? G.COMMODITIES[idx].name : 'CARGO';
  };

  // start a fresh game, docked at the starting system
  GAME.newGame = function (name) {
    GAME.cmd = GAME.NEW_COMMANDER();
    if (name) GAME.cmd.name = String(name).toUpperCase().slice(0, 12);
    GAME.cmd.system = 0;
    GAME.cmd.galaxy = 0;
    GAME.escapeUsed = false;
    GAME.enterSystem(0, { mode: 'arrive', spawnDist: 9000 });
    GAME.state = 'docked';
    GAME.screen = 'launch';
    GAME.dockedAt();
  };
})(typeof window !== 'undefined' ? window : globalThis);
