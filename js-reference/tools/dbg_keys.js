'use strict';
const vm=require('vm'),fs=require('fs'),path=require('path');
const listeners={w:{},c:{}};
const canvas={width:320,height:224,style:{},addEventListener:(t,f)=>{(listeners.c[t]=listeners.c[t]||[]).push(f);},getContext:()=>({imageSmoothingEnabled:false,createImageData:(w,h)=>({width:w,height:h,data:new Uint8ClampedArray(w*h*4)}),putImageData(){}}),getBoundingClientRect:()=>({left:0,top:0,width:320,height:224})};
const sandbox={console,Math,Date,JSON,isNaN,isFinite,parseInt,parseFloat,Uint8Array,Uint32Array,Uint8ClampedArray,
 document:{readyState:'complete',getElementById:()=>canvas,addEventListener(){}},
 performance:{now:()=>Date.now()},requestAnimationFrame:()=>0,innerWidth:1400,innerHeight:900,
 addEventListener:(t,f)=>{(listeners.w[t]=listeners.w[t]||[]).push(f);},
 localStorage:{getItem:()=>null,setItem(){},removeItem(){}}};
sandbox.window=sandbox;sandbox.globalThis=sandbox;vm.createContext(sandbox);
for(const f of ['render.js','models.js','galaxy.js','game.js','ui.js','main.js']) vm.runInContext(fs.readFileSync(path.join(__dirname,'..','js',f),'utf8'),sandbox,{filename:f});
const E=sandbox.E,GAME=E.game;
function kd(k){(listeners.w.keydown||[]).forEach(fn=>fn({key:k,shiftKey:false,repeat:false,preventDefault(){}}));}
function ku(k){(listeners.w.keyup||[]).forEach(fn=>fn({key:k,preventDefault(){}}));}
kd('1');ku('1');E.frameStep(0.016);
console.log('state',GAME.state,'screen',GAME.screen);
kd('1');ku('1');E.frameStep(0.016);
console.log('screen',GAME.screen,'sel',E.ui.sel.market,'prices',GAME.market.slice(0,2).map(m=>m.name+' '+m.price).join(','));
kd('b');ku('b');E.frameStep(0.016);
console.log('after b: cash',GAME.cmd.cash,'cargo',JSON.stringify(GAME.cmd.cargo),'sel',E.ui.sel.market);
kd('ArrowDown');ku('ArrowDown');E.frameStep(0.016);
console.log('after down: sel',E.ui.sel.market);
kd('B');ku('B');E.frameStep(0.016);
console.log('after B: cash',GAME.cmd.cash,'cargo',JSON.stringify(GAME.cmd.cargo));
