'use strict';
const vm=require('vm'),fs=require('fs'),path=require('path');
const sandbox={console,Math,Date,JSON,isNaN,isFinite,parseFloat,parseInt};
sandbox.window=sandbox;sandbox.globalThis=sandbox;vm.createContext(sandbox);
for(const f of ['render.js','models.js']) vm.runInContext(fs.readFileSync(path.join(__dirname,'..','js',f),'utf8'),sandbox,{filename:f});
const E=sandbox.E;
const hero=new E.FB(480,300);
const space=E.rgb('#05070f');
hero.clear(space);
let diff=0; for(const c of hero.buf) if(c!==space) diff++;
console.log('after clear diff',diff, 'space',space.toString(16));
let hist={};
for(let y=0;y<300;y+=7)for(let x=0;x<480;x+=7){const c=hero.buf[y*480+x];hist[c]=(hist[c]||0)+1;}
console.log(hist);
