'use strict';
const vm=require('vm'),fs=require('fs'),path=require('path');
const sandbox={console,Math,Date,JSON,isNaN,isFinite,parseInt,parseFloat,Uint8Array,Uint32Array,Uint8ClampedArray};
sandbox.window=sandbox;sandbox.globalThis=sandbox;vm.createContext(sandbox);
for(const f of ['render.js','models.js','galaxy.js','game.js','ui.js']) vm.runInContext(fs.readFileSync(path.join(__dirname,'..','js',f),'utf8'),sandbox,{filename:f});
const E=sandbox.E,GAME=E.game,UI=E.ui,V=E.V,M3=E.M3;
const fb=new E.FB(320,224);
fb.clear();
const cam=new E.Camera(320,150,72);
cam.cx=160;cam.cy=75;
const station=E.models.meshes.station;
console.log('station faces',station.faces.length,'radius',station.radius.toFixed(1));
E.drawMesh(fb,cam,station,[0,0,-900],M3.rotY(0.3),{fogDist:60000});
let painted=0;
for(let y=0;y<150;y++)for(let x=0;x<320;x++){const c=fb.buf[y*320+x];if((c&255)+((c>>8)&255)+((c>>16)&255)>90)painted++;}
console.log('painted pixels', painted);
// check first face projection
const f=station.faces[0];
const pv=f.p.map(p=>cam.toView(E.V.add(M3.apply(M3.rotY(0.3),p),[0,0,-900])));
console.log('z range',Math.min(...pv.map(v=>v[2])).toFixed(1),Math.max(...pv.map(v=>v[2])).toFixed(1));
const pr=pv.map(v=>cam.projectView(v));
console.log('screen',pr.map(p=>p.x.toFixed(0)+','+p.y.toFixed(0)).join(' '));
console.log('zbuf sample',fb.zbuf[75*320+160]);
