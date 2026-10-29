'use strict';
const vm=require('vm'),fs=require('fs'),path=require('path');
const sandbox={console,Math,Date,JSON,isNaN,isFinite,parseFloat,parseInt};
sandbox.window=sandbox;sandbox.globalThis=sandbox;vm.createContext(sandbox);
for(const f of ['render.js','models.js']) vm.runInContext(fs.readFileSync(path.join(__dirname,'..','js',f),'utf8'),sandbox,{filename:f});
const E=sandbox.E;
const mesh=E.models.meshes.cobra3;
const W=200,H=100;
const fb=new E.FB(W,H);fb.clear(0xff000000);
const cam=new E.Camera(W,H,60);
const mat=E.M3.mul(E.M3.rotX(0.35),E.M3.rotY(Math.PI*0.8));
const dist=mesh.radius*2.2;
E.drawMesh(fb,cam,mesh,[0,0,-dist],mat,{});
// count faces that survive culling by replicating
let kept=0,culled=0;
for(const face of mesh.faces){
  const pv=face.p.map(lp=>cam.toView(E.V.add(E.M3.apply(mat,lp),[0,0,-dist])));
  const c=[0,0,0]; for(const v of pv){c[0]+=v[0];c[1]+=v[1];c[2]+=v[2];}
  c[0]/=pv.length;c[1]/=pv.length;c[2]/=pv.length;
  const viewN=E.M3.applyT(cam.mat,E.M3.apply(mat,face.n));
  if(E.V.dot(viewN,E.V.norm(c))>0&&!face.double) culled++; else kept++;
}
console.log('kept',kept,'culled',culled,'dist',dist.toFixed(1));
const chars=' .:-=+*#%@';
for(let y=0;y<H;y+=2){
  let line='';
  for(let x=0;x<W;x+=1){
    let s=0;for(let yy=0;yy<2;yy++){const cc=fb.buf[(y+yy)*fb.w+x];s+=(cc&255)*0.3+((cc>>8)&255)*0.59+((cc>>16)&255)*0.11;}
    line+=chars[Math.min(9,Math.max(0,Math.round(s/2/255*11)))];
  }
  console.log(line.replace(/\s+$/,''));
}
