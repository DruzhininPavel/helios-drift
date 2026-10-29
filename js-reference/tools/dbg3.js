'use strict';
const vm=require('vm'),fs=require('fs'),path=require('path');
const sandbox={console,Math,Date,JSON,isNaN,isFinite,parseFloat,parseInt};
sandbox.window=sandbox;sandbox.globalThis=sandbox;vm.createContext(sandbox);
for(const f of ['render.js','models.js']) vm.runInContext(fs.readFileSync(path.join(__dirname,'..','js',f),'utf8'),sandbox,{filename:f});
const E=sandbox.E;
const mesh=E.models.meshes.cobra3;
const cam=new E.Camera(200,100,60);
console.log('focal',cam.focal,'cx',cam.cx,'cy',cam.cy);
const mat=E.M3.rotY(0); // identity
const d=86;
const pts=[]; for(const f of mesh.faces) for(const p of f.p) pts.push(p);
const xs=pts.map(p=>p[0]),ys=pts.map(p=>p[1]),zs=pts.map(p=>p[2]);
console.log('bounds x',Math.min(...xs).toFixed(1),Math.max(...xs).toFixed(1),'y',Math.min(...ys).toFixed(1),Math.max(...ys).toFixed(1),'z',Math.min(...zs).toFixed(1),Math.max(...zs).toFixed(1));
for(const p of [[0,0,-34],[0,0,34],[11,0,8],[-11,0,8],[0,3,0]]){
  const v=cam.toView(E.V.add(E.M3.apply(mat,p),[0,0,-d]));
  console.log(p,'->view',v.map(n=>+n.toFixed(1)),'screen',(()=>{const s=cam.projectView(v);return [ +s.x.toFixed(1), +s.y.toFixed(1)];})());
}
