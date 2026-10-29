'use strict';
const vm=require('vm'),fs=require('fs'),path=require('path');
const sandbox={console,Math,Date,JSON,isNaN,isFinite,parseInt,parseFloat,Uint8Array,Uint32Array,Uint8ClampedArray};
sandbox.window=sandbox;sandbox.globalThis=sandbox;vm.createContext(sandbox);
for(const f of ['render.js','models.js','galaxy.js','game.js','ui.js']) vm.runInContext(fs.readFileSync(path.join(__dirname,'..','js',f),'utf8'),sandbox,{filename:f});
const E=sandbox.E,GAME=E.game,UI=E.ui,V=E.V;
const fb=new E.FB(320,224);
GAME.newGame('P'); GAME.launch();
// busy scene: station + 8 ships + asteroids in view
const st=GAME.world.station, p=GAME.player;
const n=[st.mat[0][2],st.mat[1][2],st.mat[2][2]];
p.pos=V.add(st.pos,V.mul(n,1600));
p.mat=GAME.lookMatrix(V.norm(V.sub(st.pos,p.pos)));
for(let i=0;i<6;i++) GAME.spawnPirates(1);
GAME.world.entities.forEach((e,i)=>{ if(e.kind==='ship'){ e.pos=V.add(p.pos,V.mul(V.norm([0.4-i*0.15,0.1-i*0.05,-1]),600+i*200)); e.mat=GAME.lookMatrix(V.norm(V.sub(p.pos,e.pos))); }});
GAME.world.planetFrames=null;
const t0=process.hrtime.bigint();
const N=200;
for(let i=0;i<N;i++){ GAME.update(1/60); UI.draw(fb,1/60); }
const t1=process.hrtime.bigint();
console.log('avg frame ms', (Number(t1-t0)/1e6/N).toFixed(2), 'objects', GAME.world.entities.length);
