import{u as O}from"./DljAJ1KK.js";import{A as he,C as be,D as E,E as a,G as I,H as T,I as g,J as X,g as Y,K as j,L as d,M as R,N as fe,O as ve,P as ge,Q as G,R as me,s as F,S as pe,v as W,T as L,U as B,V as U,W as m,X as we,j as xe,l as ye,c as ke,a as A,b as Se,m as _e,o as Ce}from"./UPofMcO5.js";import{u as Re}from"./BFvCYdQA.js";import{N as Be}from"./B1Vfolkc.js";import{u as $e}from"./BMjm02Nn.js";import{N as Ve}from"./DmLDoBJx.js";import"./B9vZVzT3.js";import"./B-p6aW7q.js";import"./CYgBSJtP.js";import"./b9Y9N2M4.js";import"./DE4Qtn4z.js";const ze={buttonHeightSmall:"14px",buttonHeightMedium:"18px",buttonHeightLarge:"22px",buttonWidthSmall:"14px",buttonWidthMedium:"18px",buttonWidthLarge:"22px",buttonWidthPressedSmall:"20px",buttonWidthPressedMedium:"24px",buttonWidthPressedLarge:"28px",railHeightSmall:"18px",railHeightMedium:"22px",railHeightLarge:"26px",railWidthSmall:"32px",railWidthMedium:"40px",railWidthLarge:"48px"};function Fe(e){const{primaryColor:c,opacityDisabled:u,borderRadius:r,textColor3:h}=e;return Object.assign(Object.assign({},ze),{iconColor:h,textColor:"white",loadingColor:c,opacityDisabled:u,railColor:"rgba(0, 0, 0, .14)",railColorActive:c,buttonBoxShadow:"0 1px 4px 0 rgba(0, 0, 0, 0.3), inset 0 0 1px 0 rgba(0, 0, 0, 0.05)",buttonColor:"#FFF",railBorderRadiusSmall:r,railBorderRadiusMedium:r,railBorderRadiusLarge:r,buttonBorderRadiusSmall:r,buttonBorderRadiusMedium:r,buttonBorderRadiusLarge:r,boxShadowFocus:`0 0 0 2px ${be(c,{alpha:.2})}`})}const We={common:he,self:Fe},Me=E("switch",`
 height: var(--n-height);
 min-width: var(--n-width);
 vertical-align: middle;
 user-select: none;
 -webkit-user-select: none;
 display: inline-flex;
 outline: none;
 justify-content: center;
 align-items: center;
`,[a("children-placeholder",`
 height: var(--n-rail-height);
 display: flex;
 flex-direction: column;
 overflow: hidden;
 pointer-events: none;
 visibility: hidden;
 `),a("rail-placeholder",`
 display: flex;
 flex-wrap: none;
 `),a("button-placeholder",`
 width: calc(1.75 * var(--n-rail-height));
 height: var(--n-rail-height);
 `),E("base-loading",`
 position: absolute;
 top: 50%;
 left: 50%;
 transform: translateX(-50%) translateY(-50%);
 font-size: calc(var(--n-button-width) - 4px);
 color: var(--n-loading-color);
 transition: color .3s var(--n-bezier);
 `,[I({left:"50%",top:"50%",originalTransform:"translateX(-50%) translateY(-50%)"})]),a("checked, unchecked",`
 transition: color .3s var(--n-bezier);
 color: var(--n-text-color);
 box-sizing: border-box;
 position: absolute;
 white-space: nowrap;
 top: 0;
 bottom: 0;
 display: flex;
 align-items: center;
 line-height: 1;
 `),a("checked",`
 right: 0;
 padding-right: calc(1.25 * var(--n-rail-height) - var(--n-offset));
 `),a("unchecked",`
 left: 0;
 justify-content: flex-end;
 padding-left: calc(1.25 * var(--n-rail-height) - var(--n-offset));
 `),T("&:focus",[a("rail",`
 box-shadow: var(--n-box-shadow-focus);
 `)]),g("round",[a("rail","border-radius: calc(var(--n-rail-height) / 2);",[a("button","border-radius: calc(var(--n-button-height) / 2);")])]),X("disabled",[X("icon",[g("rubber-band",[g("pressed",[a("rail",[a("button","max-width: var(--n-button-width-pressed);")])]),a("rail",[T("&:active",[a("button","max-width: var(--n-button-width-pressed);")])]),g("active",[g("pressed",[a("rail",[a("button","left: calc(100% - var(--n-offset) - var(--n-button-width-pressed));")])]),a("rail",[T("&:active",[a("button","left: calc(100% - var(--n-offset) - var(--n-button-width-pressed));")])])])])])]),g("active",[a("rail",[a("button","left: calc(100% - var(--n-button-width) - var(--n-offset))")])]),a("rail",`
 overflow: hidden;
 height: var(--n-rail-height);
 min-width: var(--n-rail-width);
 border-radius: var(--n-rail-border-radius);
 cursor: pointer;
 position: relative;
 transition:
 opacity .3s var(--n-bezier),
 background .3s var(--n-bezier),
 box-shadow .3s var(--n-bezier);
 background-color: var(--n-rail-color);
 `,[a("button-icon",`
 color: var(--n-icon-color);
 transition: color .3s var(--n-bezier);
 font-size: calc(var(--n-button-height) - 4px);
 position: absolute;
 left: 0;
 right: 0;
 top: 0;
 bottom: 0;
 display: flex;
 justify-content: center;
 align-items: center;
 line-height: 1;
 `,[I()]),a("button",`
 align-items: center; 
 top: var(--n-offset);
 left: var(--n-offset);
 height: var(--n-button-height);
 width: var(--n-button-width-pressed);
 max-width: var(--n-button-width);
 border-radius: var(--n-button-border-radius);
 background-color: var(--n-button-color);
 box-shadow: var(--n-button-box-shadow);
 box-sizing: border-box;
 cursor: inherit;
 content: "";
 position: absolute;
 transition:
 background-color .3s var(--n-bezier),
 left .3s var(--n-bezier),
 opacity .3s var(--n-bezier),
 max-width .3s var(--n-bezier),
 box-shadow .3s var(--n-bezier);
 `)]),g("active",[a("rail","background-color: var(--n-rail-color-active);")]),g("loading",[a("rail",`
 cursor: wait;
 `)]),g("disabled",[a("rail",`
 cursor: not-allowed;
 opacity: .5;
 `)])]),Ne=Object.assign(Object.assign({},G.props),{size:String,value:{type:[String,Number,Boolean],default:void 0},loading:Boolean,defaultValue:{type:[String,Number,Boolean],default:!1},disabled:{type:Boolean,default:void 0},round:{type:Boolean,default:!0},"onUpdate:value":[Function,Array],onUpdateValue:[Function,Array],checkedValue:{type:[String,Number,Boolean],default:!0},uncheckedValue:{type:[String,Number,Boolean],default:!1},railStyle:Function,rubberBand:{type:Boolean,default:!0},spinProps:Object,onChange:[Function,Array]});let z;const Pe=Y({name:"Switch",props:Ne,slots:Object,setup(e){z===void 0&&(typeof CSS<"u"?typeof CSS.supports<"u"?z=CSS.supports("width","max(1px)"):z=!1:z=!0);const{mergedClsPrefixRef:c,inlineThemeDisabled:u,mergedComponentPropsRef:r}=ge(e),h=G("Switch","-switch",Me,We,e,c),f=me(e,{mergedSize(t){var w,x;if(e.size!==void 0)return e.size;if(t)return t.mergedSize.value;const C=(x=(w=r==null?void 0:r.value)===null||w===void 0?void 0:w.Switch)===null||x===void 0?void 0:x.size;return C||"medium"}}),{mergedSizeRef:S,mergedDisabledRef:v}=f,_=F(e.defaultValue),i=we(e,"value"),l=$e(i,_),b=W(()=>l.value===e.checkedValue),o=F(!1),s=F(!1),n=W(()=>{const{railStyle:t}=e;if(t)return t({focused:s.value,checked:b.value})});function p(t){const{"onUpdate:value":w,onChange:x,onUpdateValue:C}=e,{nTriggerFormInput:M,nTriggerFormChange:N}=f;w&&L(w,t),C&&L(C,t),x&&L(x,t),_.value=t,M(),N()}function J(){const{nTriggerFormFocus:t}=f;t()}function Q(){const{nTriggerFormBlur:t}=f;t()}function q(){e.loading||v.value||(l.value!==e.checkedValue?p(e.checkedValue):p(e.uncheckedValue))}function Z(){s.value=!0,J()}function ee(){s.value=!1,Q(),o.value=!1}function te(t){e.loading||v.value||t.key===" "&&(l.value!==e.checkedValue?p(e.checkedValue):p(e.uncheckedValue),o.value=!1)}function ae(t){e.loading||v.value||t.key===" "&&(t.preventDefault(),o.value=!0)}const K=W(()=>{const{value:t}=S,{self:{opacityDisabled:w,railColor:x,railColorActive:C,buttonBoxShadow:M,buttonColor:N,boxShadowFocus:ie,loadingColor:oe,textColor:ne,iconColor:re,[B("buttonHeight",t)]:y,[B("buttonWidth",t)]:le,[B("buttonWidthPressed",t)]:se,[B("railHeight",t)]:k,[B("railWidth",t)]:V,[B("railBorderRadius",t)]:de,[B("buttonBorderRadius",t)]:ce},common:{cubicBezierEaseInOut:ue}}=h.value;let P,D,H;return z?(P=`calc((${k} - ${y}) / 2)`,D=`max(${k}, ${y})`,H=`max(${V}, calc(${V} + ${y} - ${k}))`):(P=U((m(k)-m(y))/2),D=U(Math.max(m(k),m(y))),H=m(k)>m(y)?V:U(m(V)+m(y)-m(k))),{"--n-bezier":ue,"--n-button-border-radius":ce,"--n-button-box-shadow":M,"--n-button-color":N,"--n-button-width":le,"--n-button-width-pressed":se,"--n-button-height":y,"--n-height":D,"--n-offset":P,"--n-opacity-disabled":w,"--n-rail-border-radius":de,"--n-rail-color":x,"--n-rail-color-active":C,"--n-rail-height":k,"--n-rail-width":V,"--n-width":H,"--n-box-shadow-focus":ie,"--n-loading-color":oe,"--n-text-color":ne,"--n-icon-color":re}}),$=u?pe("switch",W(()=>S.value[0]),K,e):void 0;return{handleClick:q,handleBlur:ee,handleFocus:Z,handleKeyup:te,handleKeydown:ae,mergedRailStyle:n,pressed:o,mergedClsPrefix:c,mergedValue:l,checked:b,mergedDisabled:v,cssVars:u?void 0:K,themeClass:$==null?void 0:$.themeClass,onRender:$==null?void 0:$.onRender}},render(){const{mergedClsPrefix:e,mergedDisabled:c,checked:u,mergedRailStyle:r,onRender:h,$slots:f}=this;h==null||h();const{checked:S,unchecked:v,icon:_,"checked-icon":i,"unchecked-icon":l}=f,b=!(j(_)&&j(i)&&j(l));return d("div",{role:"switch","aria-checked":u,class:[`${e}-switch`,this.themeClass,b&&`${e}-switch--icon`,u&&`${e}-switch--active`,c&&`${e}-switch--disabled`,this.round&&`${e}-switch--round`,this.loading&&`${e}-switch--loading`,this.pressed&&`${e}-switch--pressed`,this.rubberBand&&`${e}-switch--rubber-band`],tabindex:this.mergedDisabled?void 0:0,style:this.cssVars,onClick:this.handleClick,onFocus:this.handleFocus,onBlur:this.handleBlur,onKeyup:this.handleKeyup,onKeydown:this.handleKeydown},d("div",{class:`${e}-switch__rail`,"aria-hidden":"true",style:r},R(S,o=>R(v,s=>o||s?d("div",{"aria-hidden":!0,class:`${e}-switch__children-placeholder`},d("div",{class:`${e}-switch__rail-placeholder`},d("div",{class:`${e}-switch__button-placeholder`}),o),d("div",{class:`${e}-switch__rail-placeholder`},d("div",{class:`${e}-switch__button-placeholder`}),s)):null)),d("div",{class:`${e}-switch__button`},R(_,o=>R(i,s=>R(l,n=>d(fe,null,{default:()=>this.loading?d(ve,Object.assign({key:"loading",clsPrefix:e,strokeWidth:20},this.spinProps)):this.checked&&(s||o)?d("div",{class:`${e}-switch__button-icon`,key:s?"checked-icon":"icon"},s||o):!this.checked&&(n||o)?d("div",{class:`${e}-switch__button-icon`,key:n?"unchecked-icon":"icon"},n||o):null})))),R(S,o=>o&&d("div",{key:"checked",class:`${e}-switch__checked`},o)),R(v,o=>o&&d("div",{key:"unchecked",class:`${e}-switch__unchecked`},o)))))}}),De={class:"card"},Ye=Y({__name:"users",setup(e){const c=F([]),u=F(!1),r=Re(),h=xe();ye(async()=>{u.value=!0;try{const{data:i}=await O.list();c.value=i}finally{u.value=!1}});const f=[{label:"Студент",value:"student"},{label:"Преподаватель",value:"teacher"},{label:"Администратор",value:"admin"}];async function S(i,l){var b,o,s;if(i===((b=h.user)==null?void 0:b.id)){r.warning("Нельзя изменить свою роль");return}try{await O.updateRole(i,l),r.success("Роль обновлена");const n=c.value.find(p=>p.id===i);n&&(n.role=l)}catch(n){r.error(((s=(o=n==null?void 0:n.response)==null?void 0:o.data)==null?void 0:s.error)||"Ошибка")}}async function v(i,l){var b,o,s;if(i===((b=h.user)==null?void 0:b.id)){r.warning("Нельзя заблокировать себя");return}try{await O.block(i,!l);const n=c.value.find(p=>p.id===i);n&&(n.is_active=!l)}catch(n){r.error(((s=(o=n==null?void 0:n.response)==null?void 0:o.data)==null?void 0:s.error)||"Ошибка")}}const _=[{title:"Имя",key:"full_name"},{title:"Email",key:"email"},{title:"Роль",key:"role",width:200,render:i=>d(Be,{value:i.role,size:"small",options:f,onUpdateValue:l=>S(i.id,l)})},{title:"Активен",key:"is_active",width:100,render:i=>d(Pe,{value:i.is_active,onUpdateValue:()=>v(i.id,i.is_active)})},{title:"Дата регистрации",key:"created_at",width:150,render:i=>new Date(i.created_at).toLocaleDateString("ru-RU")}];return(i,l)=>(Ce(),ke("div",null,[l[0]||(l[0]=A("div",{class:"page-header"},[A("h1",null,"Пользователи")],-1)),A("div",De,[Se(_e(Ve),{columns:_,data:c.value,loading:u.value},null,8,["data","loading"])])]))}});export{Ye as default};
