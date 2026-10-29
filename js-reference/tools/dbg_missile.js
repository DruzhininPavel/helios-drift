'use strict';
const vm=require('vm'),fs=require('fs'),path=require('path');
const sandbox={console,Math,Date,JSON,isNaN,isFinite,parseInt,parseFloat,Uint8Array,Uint32Array,Uint8ClampedArray};
sandbox.window=sandbox;sandbox.globalThis=sandbox;vm.createContext(sandbox);
for(const f of ['render.js','models.js','galaxy.js','game.js','ui.js']) vm.runInContext(fs.readFileSync(path.join(__dirname,'..','js',f),'utf8'),sandbox,{filename:f});
const E=sandbox.E,GAME=E.game,V=E.V;
GAME.newGame('T'); GAME.launch();
const p=GAME.player;
GAME.world.entities.length=0;
GAME.spawnPirates(1);
const e=GAME.world.entities.find(x=>x.kind==='ship');
e.type='viper'; e.canMissile=true;
e.pos=V.add(p.pos,V.mul([-p.mat[0][2],-p.mat[1][2],-p.mat[2][2]],2000));
e.mat=GAME.lookMatrix(V.norm(V.sub(p.pos,e.pos)));
let missiles=0;
for(let i=0;i<600;i++){
  GAME.update(1/60);
  const m=GAME.world.entities.filter(x=>x.kind==='missile');
  if(m.length && !missiles){ missiles=1; console.log('missile launched at frame',i,'hp',p.energy.toFixed(0)); }
  if(GAME.state!=='flight') break;
}
console.log('after: energy',p.energy.toFixed(0),'shieldF',p.shieldF.toFixed(0),'state',GAME.state,'missiles alive',GAME.world.entities.filter(x=>x.kind==='missile').length);
// ECM test
GAME.launch();
GAME.cmd.equipment.ecm=true;
GAME.enemyMissile({pos:GAME.player.pos, mat:GAME.player.mat, radius:20});
console.log('missile before ecm',GAME.world.entities.filter(x=>x.kind==='missile').length);
GAME.fireECM();
console.log('missile after ecm',GAME.world.entities.filter(x=>x.kind==='missile').length);
