'use strict';
const vm=require('vm'),fs=require('fs'),path=require('path');
const sandbox={console,Math,Date,JSON,isNaN,isFinite,parseInt,parseFloat,Uint8Array,Uint32Array,Uint8ClampedArray};
sandbox.window=sandbox;sandbox.globalThis=sandbox;vm.createContext(sandbox);
for(const f of ['render.js','models.js','galaxy.js','game.js','ui.js']) vm.runInContext(fs.readFileSync(path.join(__dirname,'..','js',f),'utf8'),sandbox,{filename:f});
const E=sandbox.E,GAME=E.game,V=E.V,M3=E.M3;
GAME.newGame('T');
GAME.launch();
GAME.cmd.equipment.dockcomp=true;
GAME.player.pos=[3000,900,-3000];
GAME.toggleAutopilot();
console.log('autopilot',GAME.autopilot,'pos',GAME.player.pos.map(x=>x.toFixed(0)));
const st=GAME.world.station;
console.log('station',st.pos.map(x=>x.toFixed(0)),'normal',[st.mat[0][2],st.mat[1][2],st.mat[2][2]].map(x=>x.toFixed(2)));
for(let i=0;i<1200;i++){
  GAME.update(1/60);
  if(i%120===0){
    const p=GAME.player;
    const rel=V.sub(p.pos,st.pos);
    const local=M3.applyT(st.mat,rel);
    const fwd=[-p.mat[0][2],-p.mat[1][2],-p.mat[2][2]];
    const normal=[st.mat[0][2],st.mat[1][2],st.mat[2][2]];
    console.log(i, 'dist',V.len(rel).toFixed(0),'local',local.map(x=>x.toFixed(0)).join(','),'spd',p.speed.toFixed(0),'fwd.-n',V.dot(fwd,V.mul(normal,-1)).toFixed(2),'state',GAME.state, 'auto',GAME.autopilot);
  }
  if(GAME.state!=='flight') {console.log('ENDED',GAME.state, i); break;}
}
