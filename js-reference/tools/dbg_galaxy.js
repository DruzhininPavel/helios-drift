'use strict';
const vm=require('vm'),fs=require('fs'),path=require('path');
const sandbox={console,Math,Date,JSON,isNaN,isFinite,parseInt,parseFloat};
sandbox.window=sandbox;sandbox.globalThis=sandbox;vm.createContext(sandbox);
for(const f of ['render.js','models.js','galaxy.js']) vm.runInContext(fs.readFileSync(path.join(__dirname,'..','js',f),'utf8'),sandbox,{filename:f});
const E=sandbox.E,G=E.galaxy;
for(let gi=0;gi<8;gi++){
  const gal=G.get(gi);
  // BFS from system 0 with 7 LY jumps
  const seen=new Set([0]); const q=[0];
  while(q.length){ const cur=gal.systems[q.shift()];
    for(const s of gal.systems){ if(seen.has(s.id)) continue; if(G.distance(cur,s)<=7){ seen.add(s.id); q.push(s.id);} } }
  // average nearest-neighbour distance
  let sum=0, minD=1e9, maxD=0;
  for(const s of gal.systems){ let best=1e9; for(const t of gal.systems){ if(t===s) continue; const d=G.distance(s,t); if(d<best) best=d; } sum+=best; if(best<minD)minD=best; if(best>maxD)maxD=best; }
  console.log('galaxy',gi,'reachable',seen.size+'/256','avg nearest',(sum/256).toFixed(2),'min',minD.toFixed(2),'max',maxD.toFixed(2));
}
// market profitability check
const gal=G.get(0);
const markets=gal.systems.map(s=>G.makeMarket(s,1));
let best=null;
for(const a of gal.systems) for(const b of gal.systems){
  if(a.id===b.id) continue;
  if(G.distance(a,b)>7) continue;
  const ma=markets[a.id], mb=markets[b.id];
  for(let i=0;i<17;i++){
    const profit=mb[i].sell-ma[i].price;
    if(!best || profit>best.profit) best={profit,from:a.name,to:b.name,item:ma[i].name,buy:ma[i].price,sell:mb[i].sell,d:G.distance(a,b).toFixed(1)};
  }
}
console.log('best single-unit arbitrage:',JSON.stringify(best));
console.log('avg commodity price at sys0:', (markets[0].reduce((s,m)=>s+m.price,0)/17).toFixed(0));
