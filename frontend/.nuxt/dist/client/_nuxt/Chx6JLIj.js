import{bO as r,bP as t}from"./UPofMcO5.js";const n=r(()=>{const e=localStorage.getItem("user");if(!e)return t("/auth/login");if(JSON.parse(e).role!=="admin")return t("/")});export{n as default};
