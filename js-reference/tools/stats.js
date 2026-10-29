'use strict';
const vm=require('vm'),fs=require('fs'),path=require('path');
const sandbox={console,Math,Date,JSON,isNaN,isFinite,parseInt,parseFloat,Uint8Array,Uint32Array,Uint8ClampedArray};
sandbox.window=sandbox;sandbox.globalThis=sandbox;vm.createContext(sandbox);
for(const f of ['render.js','models.js','galaxy.js','game.js','ui.js']) vm.runInContext(fs.readFileSync(path.join(__dirname,'..','js',f),'utf8'),sandbox,{filename:f});
const E=sandbox.E,GAME=E.game,UI=E.ui,V=E.V;
const fb=new E.FB(320,224);

function stats(label, region){
  let painted=0,total=0;
  const shades={};
  const x0=region?region[0]:0,x1=region?region[2]:319,y0=region?region[1]:0,y1=region?region[3]:149;
  for(let y=y0;y<=y1;y++)for(let x=x0;x<=x1;x++){
    total++;
    const c=fb.buf[y*320+x];
    const r=c&255,g=(c>>8)&255,b=(c>>16)&255;
    if(r+g+b > 90) painted++;
    shades[c]=(shades[c]||0)+1;
  }
  const top=Object.entries(shades).sort((a,b)=>b[1]-a[1]).slice(0,6);
  console.log(label.padEnd(22), 'painted', (painted/total*100).toFixed(1)+'%', 'colours', Object.keys(shades).length,
    'top', top.map(([c,n])=>'#'+((c&255).toString(16).padStart(2,'0'))+(((c>>8)&255).toString(16).padStart(2,'0'))+(((c>>16)&255).toString(16).padStart(2,'0'))+':'+n).join(' '));
}

GAME.newGame('S'); GAME.launch();
const p=GAME.player,st=GAME.world.station;
const n=[st.mat[0][2],st.mat[1][2],st.mat[2][2]];
GAME.world.entities.length=0;
GAME.messages.length=0;
// station at 900 units, filling a good part of the frame
p.pos=V.add(st.pos,V.mul(n,900)); p.pos[1]+=60;
p.mat=GAME.lookMatrix(V.norm(V.sub(st.pos,p.pos)));
GAME.time=5;
UI.draw(fb,1/60); stats('station 900u',[0,0,319,149]);
// empty space
GAME.world=JSON.parse(JSON.stringify({sys:GAME.world.sys}));
