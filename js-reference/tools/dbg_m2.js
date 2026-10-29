'use strict';
const vm=require('vm'),fs=require('fs'),path=require('path');
const sandbox={console,Math,Date,JSON,isNaN,isFinite,parseInt,parseFloat,Uint8Array,Uint32Array,Uint8ClampedArray};
sandbox.window=sandbox;sandbox.globalThis=sandbox;vm.createContext(sandbox);
for(const f of ['render.js','models.js','galaxy.js','game.js','ui.js']) vm.runInContext(fs.readFileSync(path.join(__dirname,'..','js',f),'utf8'),sandbox,{filename:f});
const E=sandbox.E,GAME=E.game,V=E.V,M3=E.M3;
GAME.newGame('T'); GAME.launch();
const p=GAME.player;
GAME.world.entities.length=0;
GAME.spawnPirates(1);
const e=GAME.world.entities.find(x=>x.kind==='ship');
e.type='viper'; e.canMissile=true; e.cooldown=0;
e.pos=V.add(p.pos,V.mul([-p.mat[0][2],-p.mat[1][2],-p.mat[2][2]],2000));
e.mat=GAME.lookMatrix(V.norm(V.sub(p.pos,e.pos)));
console.log('behaviour',e.behaviour,'canMissile',e.canMissile,'cooldown',e.cooldown);
for(let i=0;i<30;i++){
  GAME.update(1/60);
  const rel=V.sub(p.pos,e.pos);
  const loc=M3.applyT(e.mat,rel);
  const dist=V.len(rel);
  if(i<3||i%10===0) console.log(i,'dist',dist.toFixed(0),'loc',loc.map(x=>x.toFixed(1)).join(','),'missiles',GAME.world.entities.filter(x=>x.kind==='missile').length,'cool',e.cooldown.toFixed(2));
}
