'use strict';
const vm=require('vm'),fs=require('fs'),path=require('path');
const sandbox={console,Math,Date,JSON,isNaN,isFinite,parseFloat,parseInt};
sandbox.window=sandbox;sandbox.globalThis=sandbox;vm.createContext(sandbox);
for(const f of ['render.js','models.js']) vm.runInContext(fs.readFileSync(path.join(__dirname,'..','js',f),'utf8'),sandbox,{filename:f});
const E=sandbox.E;
const b=new E.models.Builder().at([0,0,0]);
E.models.box(b,0,0,0,40,40,40,'steel');
const mesh=b.mesh();
const W=200,H=100,fb=new E.FB(W,H);fb.clear(0xff000000);
const cam=new E.Camera(W,H,60);
const mat=E.M3.mul(E.M3.rotX(-0.5),E.M3.rotY(0.6));
E.drawMesh(fb,cam,mesh,[0,0,-140],mat,{});
let x0=1e9,y0=1e9,x1=-1,y1=-1,n=0;
for(let y=0;y<H;y++)for(let x=0;x<W;x++){const c=fb.buf[y*W+x];if(c!==0xff000000){n++;if(x<x0)x0=x;if(x>x1)x1=x;if(y<y0)y0=y;if(y>y1)y1=y;}}
console.log('painted',n,'bbox',x0,y0,x1,y1,'center',((x0+x1)/2).toFixed(1),((y0+y1)/2).toFixed(1));
// row histogram
for(let y=y0;y<=y1;y+=4){let c=0;for(let x=0;x<W;x++)if(fb.buf[y*W+x]!==0xff000000)c++;console.log('row',y,'count',c);}
