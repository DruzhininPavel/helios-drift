'use strict';
const vm=require('vm'),fs=require('fs'),path=require('path');
const sandbox={console,Math,Date,JSON,isNaN,isFinite,parseFloat,parseInt};
sandbox.window=sandbox;sandbox.globalThis=sandbox;vm.createContext(sandbox);
for(const f of ['render.js','models.js']) vm.runInContext(fs.readFileSync(path.join(__dirname,'..','js',f),'utf8'),sandbox,{filename:f});
const E=sandbox.E;
const b=new E.models.Builder().at([0,0,0]);
E.models.box(b,0,0,0,40,40,40,'steel');
const mesh=b.mesh();
const cam=new E.Camera(120,60,60);
const mat=E.M3.mul(E.M3.rotX(-0.5),E.M3.rotY(0.6));
console.log('mat',JSON.stringify(mat));
for(const f of mesh.faces.slice(0,2)){
  console.log('face ramp',f.ramp,'n',f.n.map(x=>+x.toFixed(2)),'pts',f.p.length);
  for(const p of f.p){
    const w=E.V.add(E.M3.apply(mat,p),[0,0,-140]);
    const v=cam.toView(w);
    const s=cam.projectView(v);
    console.log('  local',p,'world',w.map(x=>+x.toFixed(1)),'screen',+s.x.toFixed(1),+s.y.toFixed(1));
  }
}
