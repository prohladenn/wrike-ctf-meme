import{z as v,B as i,O as x,A as T,C as N,X as A,E as s,o,y as l,w as e,b as t,G as B,V as M,H as U,d as n,t as m,M as I,aw as R,J as j,a as O,aN as P,x as S,aO as h,aP as z}from"./DIu47odP.js";import{u as D}from"./4_Vd_w68.js";import{u as E}from"./CpZuf1_j.js";import{V as p}from"./sVARX7W8.js";import{V}from"./DGsDAbiX.js";import{V as _}from"./Dhe5bGlE.js";const F=v({__name:"[id]",setup(G){const a=i(null),d=i("Unknown"),f=i(!0),y=x(),g=T(),c=Number(y.params.id),{fetchTemplate:C,getTemplateImageUrl:w}=D(),{fetchUser:b}=E(),k=()=>{g.push(`/meme/new/${c}`)};return N(async()=>{if(a.value=await C(c),a.value){const u=await b(a.value.owner_id);u&&(d.value=u.username)}}),A(a,()=>{a.value&&(f.value=!1)}),(u,r)=>s(f)?(o(),l(p,{key:1},{default:e(()=>[t(_,{justify:"center"},{default:e(()=>[t(V,{cols:"12",md:"8",lg:"6"},{default:e(()=>[t(z,{indeterminate:"",color:"primary"})]),_:1})]),_:1})]),_:1})):(o(),l(p,{key:0},{default:e(()=>[t(_,{justify:"center"},{default:e(()=>[t(V,{cols:"12",md:"8",lg:"6"},{default:e(()=>[s(a)?(o(),l(B,{key:0},{default:e(()=>[t(M,{src:s(w)(s(a).id),"aspect-ratio":"1",class:"white--text align-end"},{default:e(()=>[t(U,{class:"bg-black bg-opacity-50"},{default:e(()=>[n(m(s(a).name),1)]),_:1}),s(a).private_info?(o(),l(I,{key:0,class:"bg-grey bg-opacity-50 mt-2 pt-2"},{default:e(()=>[n(m(s(a).private_info),1)]),_:1})):R("",!0)]),_:1},8,["src"]),t(j,{class:"mt-2"},{default:e(()=>[r[0]||(r[0]=n(" Created by ")),O("strong",null,m(s(d)),1)]),_:1}),t(P,null,{default:e(()=>[t(S,{color:"primary",onClick:k},{default:e(()=>r[1]||(r[1]=[n(" Create Meme ")])),_:1})]),_:1})]),_:1})):(o(),l(h,{key:1,type:"error",dismissible:""},{default:e(()=>r[2]||(r[2]=[n(" Template not found. ")])),_:1}))]),_:1})]),_:1})]),_:1}))}});export{F as default};