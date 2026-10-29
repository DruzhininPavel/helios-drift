/* ELITE 16 - galaxy and market model
 * 8 galaxies x 256 systems, Elite-style names, economy/government/tech level,
 * 17 commodities with economy-driven price gradients and per-dock volatility.
 */
(function (global) {
  'use strict';
  var E = global.E;
  var G = E.galaxy = {};

  var PAIRS = [
    'AB', 'OU', 'SE', 'IT', 'IL', 'ET', 'ST', 'ON', 'LO', 'NU', 'TH', 'NO',
    'AL', 'LE', 'XE', 'GE', 'ZA', 'CE', 'BI', 'SO', 'US', 'ES', 'AR', 'MA',
    'IN', 'DI', 'RE', 'A', 'ER', 'AT', 'EN', 'BE', 'RA', 'LA', 'VE', 'TI',
    'ED', 'OR', 'QU', 'AN', 'TE', 'IS', 'RI', 'ON'
  ];

  var ECONOMY_NAMES = ['RICH INDUSTRIAL', 'AVERAGE INDUSTRIAL', 'POOR INDUSTRIAL',
    'MAINLY INDUSTRIAL', 'MAINLY AGRICULTURAL', 'RICH AGRICULTURAL',
    'AVERAGE AGRICULTURAL', 'POOR AGRICULTURAL'];
  var GOV_NAMES = ['ANARCHY', 'FEUDAL', 'MULTI-GOVERNMENT', 'DICTATORSHIP',
    'COMMUNIST', 'CONFEDERACY', 'DEMOCRACY', 'CORPORATE STATE'];

  var COMMODITIES = [
    { name: 'FOOD', base: 19, grad: 0.55, qty: 14, unit: 'T' },
    { name: 'TEXTILES', base: 20, grad: 0.45, qty: 14, unit: 'T' },
    { name: 'RADIOACTIVES', base: 65, grad: 0.30, qty: 8, unit: 'T' },
    { name: 'SLAVES', base: 40, grad: -0.25, qty: 6, unit: 'T' },
    { name: 'LIQUOR/WINES', base: 83, grad: 0.40, qty: 8, unit: 'T' },
    { name: 'LUXURIES', base: 196, grad: -0.35, qty: 6, unit: 'T' },
    { name: 'NARCOTICS', base: 235, grad: -0.30, qty: 5, unit: 'T' },
    { name: 'COMPUTERS', base: 154, grad: -0.60, qty: 8, unit: 'T' },
    { name: 'MACHINERY', base: 117, grad: -0.55, qty: 8, unit: 'T' },
    { name: 'ALLOYS', base: 78, grad: -0.40, qty: 10, unit: 'T' },
    { name: 'FIREARMS', base: 124, grad: 0.20, qty: 6, unit: 'T' },
    { name: 'FURS', base: 176, grad: 0.35, qty: 6, unit: 'T' },
    { name: 'MINERALS', base: 32, grad: -0.20, qty: 12, unit: 'T' },
    { name: 'GOLD', base: 97, grad: 0.15, qty: 6, unit: 'KG' },
    { name: 'PLATINUM', base: 171, grad: 0.10, qty: 5, unit: 'KG' },
    { name: 'GEM-STONES', base: 45, grad: -0.10, qty: 4, unit: 'KG' },
    { name: 'ALIEN ITEMS', base: 53, grad: -0.90, qty: 2, unit: 'T' }
  ];
  G.COMMODITIES = COMMODITIES;
  G.ECONOMY_NAMES = ECONOMY_NAMES;
  G.GOV_NAMES = GOV_NAMES;

  function makeName(rnd) {
    var n = 2 + Math.floor(rnd() * 3);
    var s = '';
    for (var i = 0; i < n; i++) s += PAIRS[Math.floor(rnd() * PAIRS.length)];
    if (s.length > 9) s = s.slice(0, 9);
    if (s.length < 4) s += PAIRS[Math.floor(rnd() * PAIRS.length)];
    return s;
  }

  function generate(index) {
    var rnd = E.rng(0x5a4a27 + index * 0x9e3779b9);
    var systems = [];
    var used = {};
    for (var i = 0; i < 256; i++) {
      var name = makeName(rnd);
      var guard = 0;
      while (used[name] && guard++ < 40) name = makeName(rnd);
      if (used[name]) name = name + (i % 10);
      used[name] = true;

      var economy = Math.floor(rnd() * 8);
      var government = Math.floor(rnd() * 8);
      // anarchy/feudal tends to be poorer and more dangerous, corporates rich
      if (government >= 6 && rnd() < 0.5) economy = Math.floor(rnd() * 4);
      if (government <= 1 && rnd() < 0.5) economy = 4 + Math.floor(rnd() * 4);

      var tech = Math.max(1, Math.min(15,
        Math.round(1 + (7 - economy) * 0.6 + government * 0.35 + rnd() * 4)));
      var pop = +(0.4 + rnd() * 11).toFixed(1);
      var prod = Math.round((economy <= 3 ? 1 : 0.6) * tech * pop * 4 + rnd() * 400);
      var radius = Math.round(2000 + rnd() * 7000);
      systems.push({
        id: i,
        name: name,
        x: Math.floor(rnd() * 256),
        y: Math.floor(rnd() * 256),
        economy: economy,
        government: government,
        tech: tech,
        pop: pop,
        prod: prod,
        radius: radius,
        kind: E.models.planetKindFor(i * 977 + index * 31 + 5),
        seed: (i * 2654435761 + index * 40503) >>> 0
      });
    }
    // spread positions out a little so the chart is not a uniform fog
    for (var k = 0; k < 256; k++) {
      var s0 = systems[k];
      for (var j = 0; j < 8; j++) {
        var other = systems[Math.floor(rnd() * 256)];
        if (other === s0) continue;
        var dx = s0.x - other.x, dy = s0.y - other.y;
        var d = Math.sqrt(dx * dx + dy * dy);
        if (d < 2.4) { s0.x = (other.x + 3 + Math.floor(rnd() * 9)) % 256; s0.y = (other.y + 3 + Math.floor(rnd() * 9)) % 256; }
      }
    }
    // final de-dupe pass: no two systems on top of each other
    for (var pass = 0; pass < 6; pass++) {
      var moved = 0;
      for (var a = 0; a < systems.length; a++) {
        for (var b = a + 1; b < systems.length; b++) {
          var A = systems[a], B = systems[b];
          var ddx = A.x - B.x, ddy = A.y - B.y;
          if (ddx * ddx + ddy * ddy < 4) {
            B.x = (B.x + 5 + Math.floor(rnd() * 7)) % 256;
            B.y = (B.y + 5 + Math.floor(rnd() * 7)) % 256;
            moved++;
          }
        }
      }
      if (!moved) break;
    }
    ensureConnected(systems, rnd);
    return systems;
  }

  // Elite's galaxy is fully navigable: nudge any system that is out of jump
  // range of the rest so every system can be reached in 7 LY hops.
  function ensureConnected(systems, rnd) {
    for (var pass = 0; pass < 4; pass++) {
      var seen = { 0: true };
      var queue = [0];
      while (queue.length) {
        var cur = systems[queue.shift()];
        for (var i = 0; i < systems.length; i++) {
          var s = systems[i];
          if (seen[s.id]) continue;
          if (lyBetween(cur, s) <= 6.8) { seen[s.id] = true; queue.push(s.id); }
        }
      }
      var reach = Object.keys(seen).length;
      if (reach === systems.length) return;
      var lonely = [];
      for (var k = 0; k < systems.length; k++) if (!seen[systems[k].id]) lonely.push(systems[k]);
      for (var m = 0; m < lonely.length; m++) {
        var host = systems[Math.floor(rnd() * systems.length)];
        if (!seen[host.id]) host = systems[0];
        var ang = rnd() * Math.PI * 2;
        var dist = 3 + rnd() * 3;
        lonely[m].x = (host.x + Math.cos(ang) * dist * 4 + 256) % 256 | 0;
        lonely[m].y = (host.y + Math.sin(ang) * dist * 4 + 256) % 256 | 0;
      }
    }
  }

  function lyBetween(a, b) {
    var dx = a.x - b.x, dy = a.y - b.y;
    return Math.sqrt(dx * dx + dy * dy) / 4;
  }

  var cache = {};
  G.get = function (index) {
    index = ((index % 8) + 8) % 8;
    if (!cache[index]) {
      cache[index] = { index: index, systems: generate(index) };
    }
    return cache[index];
  };

  // chart distance in light years: chart is 256 units across, 1 LY = 4 units
  G.distance = function (a, b) {
    var dx = a.x - b.x, dy = a.y - b.y;
    return Math.sqrt(dx * dx + dy * dy) / 4;
  };

  G.findByName = function (gal, name) {
    for (var i = 0; i < gal.systems.length; i++) {
      if (gal.systems[i].name === name) return gal.systems[i];
    }
    return null;
  };

  // market: base price level per system, with a volatility salt per dock
  G.makeMarket = function (sys, salt) {
    var rnd = E.rng((sys.seed ^ ((salt || 0) * 2654435761)) >>> 0);
    var ind = (7 - sys.economy) / 7;    // 0 = agricultural, 1 = industrial
    var out = [];
    for (var i = 0; i < COMMODITIES.length; i++) {
      var c = COMMODITIES[i];
      var level = 1 + c.grad * (ind - 0.45) * 1.15;
      level = Math.max(0.35, level);
      var fluct = 0.88 + rnd() * 0.26;
      var price = Math.max(1, Math.round(c.base * level * fluct));
      var q = Math.round(c.qty * (0.45 + rnd() * 1.5) * (1.15 - 0.45 * c.grad * (ind - 0.45)));
      q = Math.max(0, q);
      out.push({
        name: c.name, unit: c.unit,
        price: price,
        sell: Math.max(1, Math.round(price * 0.80)),
        qty: q
      });
    }
    return out;
  };

  // equipment availability / price by tech level
  G.EQUIPMENT = [
    { id: 'missile', name: 'MISSILE', price: 30, tech: 1, repeat: true },
    { id: 'largebay', name: 'LARGE CARGO BAY', price: 400, tech: 2 },
    { id: 'ecm', name: 'E.C.M. SYSTEM', price: 600, tech: 3 },
    { id: 'fuelscoop', name: 'FUEL SCOOP', price: 525, tech: 4 },
    { id: 'laserpulse', name: 'PULSE LASER', price: 400, tech: 2 },
    { id: 'laserbeam', name: 'BEAM LASER', price: 1000, tech: 5 },
    { id: 'lasermilitary', name: 'MILITARY LASER', price: 6000, tech: 10 },
    { id: 'mininglaser', name: 'MINING LASER', price: 800, tech: 4 },
    { id: 'bomb', name: 'ENERGY BOMB', price: 900, tech: 7 },
    { id: 'escape', name: 'ESCAPE CAPSULE', price: 1000, tech: 6 },
    { id: 'extraenergy', name: 'EXTRA ENERGY UNIT', price: 1500, tech: 8 },
    { id: 'dockcomp', name: 'DOCKING COMPUTER', price: 1500, tech: 9 },
    { id: 'galdrive', name: 'GALACTIC HYPERDRIVE', price: 5000, tech: 11 }
  ];

  G.availableEquipment = function (sys) {
    return G.EQUIPMENT.filter(function (e) { return sys.tech >= e.tech; });
  };

  G.fuelPricePerLY = 2; // credits per light year

  G.rankFor = function (kills) {
    var ranks = ['HARMLESS', 'MOSTLY HARMLESS', 'POOR', 'AVERAGE', 'ABOVE AVERAGE',
      'COMPETENT', 'DANGEROUS', 'DEADLY', 'LEGEND'];
    var th = [0, 8, 16, 32, 64, 128, 512, 2560, 6400];
    var r = 0;
    for (var i = 0; i < th.length; i++) if (kills >= th[i]) r = i;
    return ranks[r];
  };
})(typeof window !== 'undefined' ? window : globalThis);
