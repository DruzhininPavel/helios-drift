'use strict';
const vm=require('vm'),fs=require('fs'),path=require('path');
const sandbox={console,Math,Date,JSON,isNaN,isFinite,parseFloat,parseInt};
sandbox.window=sandbox;sandbox.globalThis=sandbox;vm.createContext(sandbox);
for(const f of ['render.js','models.js']) vm.runInContext(fs.readFileSync(path.join(__dirname,'..','js',f),'utf8'),sandbox,{filename:f});
const E=sandbox.E;
const b=new E.models.Builder().at([0,0,0]);
E.models.box(b,0,0,0,40,40,40,'steel');
const mesh=b.mesh();
const W=120,H=60,fb=new E.FB(W,H);fb.clear(0xff000000);
const cam=new E.Camera(W,H,60);
const mat=E.M3.mul(E.M3.rotX(-0.5),E.M3.rotY(0.6));
E.drawMesh(fb,cam,mesh,[0,0,-140],mat,{});
const chars=' .:-=+*#%@';
for(let y=0;y<H;y+=2){let l='';for(let x=0;x<W;x++){const c=fb.buf[y*fb.w+x];const v=((c&255)*0.3+((c>>8)&255)*0.59+((c>>16)&255)*0.11)/255;l+=chars[Math.min(9,Math.round(v*10))];}console.log(l.replace(/\s+$/,''));}
