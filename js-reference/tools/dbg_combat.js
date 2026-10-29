'use strict';
const vm=require('vm'),fs=require('fs'),path=require('path');
const sandbox={console,Math,Date,JSON,isNaN,isFinite,parseInt,parseFloat,Uint8Array,Uint32Array,Uint8ClampedArray};
sandbox.window=sandbox;sandbox.globalThis=sandbox;vm.createContext(sandbox);
for(const f of ['render.js','models.js','galaxy.js','game.js','ui.js']) vm.runInContext(fs.readFileSync(path.join(__dirname,'..','js',f),'utf8'),sandbox,{filename:f});
const E=sandbox.E,GAME=E.game,V=E.V,M3=E.M3;
GAME.newGame('T'); GAME.launch();
const p=GAME.player;
// put a sidewinder 800 units dead ahead
const fwd=[-p.mat[0][2],-p.mat[1][2],-p.mat[2][2]];
const s=GAME.spawnPirates(1);
const e=GAME.world.entities.find(x=>x.kind==='ship');
console.log('enemy',e.type,'hp',e.energy,'pos',e.pos.map(v=>v.toFixed(0)));
e.pos=V.add(p.pos,V.mul(fwd,800));
// fire 20 times through the normal update path
GAME.in.fire=true;
for(let i=0;i<60;i++){ GAME.update(1/30); }
console.log('after 2s fire: enemy hp',e.energy.toFixed(1),'alive',GAME.world.entities.indexOf(e)>=0,'laserTemp',p.laserTemp.toFixed(2));
// direct damage
GAME.cmd.laser=2;
for(let i=0;i<200;i++){ GAME.update(1/30); if(GAME.world.entities.indexOf(e)<0) break; }
console.log('military laser: alive',GAME.world.entities.indexOf(e)>=0,'cash',GAME.cmd.cash,'kills',GAME.cmd.kills);
// missile test
GAME.launch();
const e2=GAME.spawnPirates(1);
const t=GAME.world.entities.find(x=>x.kind==='ship');
t.pos=V.add(GAME.player.pos,V.mul([-GAME.player.mat[0][2],-GAME.player.mat[1][2],-GAME.player.mat[2][2]],900));
console.log('missiles',GAME.cmd.missiles,'target',!!GAME.missileTargetPick());
GAME.fireMissile();
for(let i=0;i<200;i++){ GAME.update(1/30); if(GAME.world.entities.indexOf(t)<0) break; }
console.log('missile result: target alive',GAME.world.entities.indexOf(t)>=0,'missiles',GAME.cmd.missiles);
// ECM
GAME.fireECM();
console.log('ecm ok');
// damage to player
const before=GAME.player.energy+GAME.player.shieldF;
GAME.player.invulnerable=0;
GAME.damagePlayer(40,true);
console.log('player took damage',before-(GAME.player.energy+GAME.player.shieldF));
