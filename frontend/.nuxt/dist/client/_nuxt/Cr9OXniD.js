import{u as X}from"./B9vZVzT3.js";import{g as H}from"./CiHFR58E.js";import{Y as pe,A as ge,D as u,I,H as P,E as v,g as G,P as le,Q as ce,L as z,Z as fe,_ as ve,M as ee,$ as te,i as he,a0 as _e,S as ye,v as N,U as R,a1 as ze,j as xe,l as be,m as n,c as b,b as M,a as t,t as f,y as k,w as $,q as x,a2 as K,F as ie,r as ne,n as oe,p as Ce,d as T,s as B,o as p,B as q}from"./UPofMcO5.js";import{u as we}from"./CUAbPExj.js";import{N as Se,a as ke}from"./B1Vfolkc.js";import{N as $e}from"./CYgBSJtP.js";import{f as Te}from"./B-p6aW7q.js";import{u as Be}from"./BFvCYdQA.js";import{N as Ne}from"./Ba3mRcTH.js";import{N as Ie}from"./BMjm02Nn.js";import{_ as Pe}from"./DlAUqK2U.js";let se=!1;function Re(){if(pe&&window.CSS&&!se&&(se=!0,"registerProperty"in(window==null?void 0:window.CSS)))try{CSS.registerProperty({name:"--n-color-start",syntax:"<color>",inherits:!1,initialValue:"#0000"}),CSS.registerProperty({name:"--n-color-end",syntax:"<color>",inherits:!1,initialValue:"#0000"})}catch{}}const Le={titleMarginMedium:"0 0 6px 0",titleMarginLarge:"-2px 0 6px 0",titleFontSizeMedium:"14px",titleFontSizeLarge:"16px",iconSizeMedium:"14px",iconSizeLarge:"14px"};function Me(o){const{textColor3:c,infoColor:i,errorColor:d,successColor:r,warningColor:l,textColor1:h,textColor2:g,railColor:_,fontWeightStrong:C,fontSize:y}=o;return Object.assign(Object.assign({},Le),{contentFontSize:y,titleFontWeight:C,circleBorder:`2px solid ${c}`,circleBorderInfo:`2px solid ${i}`,circleBorderError:`2px solid ${d}`,circleBorderSuccess:`2px solid ${r}`,circleBorderWarning:`2px solid ${l}`,iconColor:c,iconColorInfo:i,iconColorError:d,iconColorSuccess:r,iconColorWarning:l,titleTextColor:h,contentTextColor:g,metaTextColor:c,lineColor:_})}const je={common:ge,self:Me},re=1.25,Fe=u("timeline",`
 position: relative;
 width: 100%;
 display: flex;
 flex-direction: column;
 line-height: ${re};
`,[I("horizontal",`
 flex-direction: row;
 `,[P(">",[u("timeline-item",`
 flex-shrink: 0;
 padding-right: 40px;
 `,[I("dashed-line-type",[P(">",[u("timeline-item-timeline",[v("line",`
 background-image: linear-gradient(90deg, var(--n-color-start), var(--n-color-start) 50%, transparent 50%, transparent 100%);
 background-size: 10px 1px;
 `)])])]),P(">",[u("timeline-item-content",`
 margin-top: calc(var(--n-icon-size) + 12px);
 `,[P(">",[v("meta",`
 margin-top: 6px;
 margin-bottom: unset;
 `)])]),u("timeline-item-timeline",`
 width: 100%;
 height: calc(var(--n-icon-size) + 12px);
 `,[v("line",`
 left: var(--n-icon-size);
 top: calc(var(--n-icon-size) / 2 - 1px);
 right: 0px;
 width: unset;
 height: 2px;
 `)])])])])]),I("right-placement",[u("timeline-item",[u("timeline-item-content",`
 text-align: right;
 margin-right: calc(var(--n-icon-size) + 12px);
 `),u("timeline-item-timeline",`
 width: var(--n-icon-size);
 right: 0;
 `)])]),I("left-placement",[u("timeline-item",[u("timeline-item-content",`
 margin-left: calc(var(--n-icon-size) + 12px);
 `),u("timeline-item-timeline",`
 left: 0;
 `)])]),u("timeline-item",`
 position: relative;
 `,[P("&:last-child",[u("timeline-item-timeline",[v("line",`
 display: none;
 `)]),u("timeline-item-content",[v("meta",`
 margin-bottom: 0;
 `)])]),u("timeline-item-content",[v("title",`
 margin: var(--n-title-margin);
 font-size: var(--n-title-font-size);
 transition: color .3s var(--n-bezier);
 font-weight: var(--n-title-font-weight);
 color: var(--n-title-text-color);
 `),v("content",`
 transition: color .3s var(--n-bezier);
 font-size: var(--n-content-font-size);
 color: var(--n-content-text-color);
 `),v("meta",`
 transition: color .3s var(--n-bezier);
 font-size: 12px;
 margin-top: 6px;
 margin-bottom: 20px;
 color: var(--n-meta-text-color);
 `)]),I("dashed-line-type",[u("timeline-item-timeline",[v("line",`
 --n-color-start: var(--n-line-color);
 transition: --n-color-start .3s var(--n-bezier);
 background-color: transparent;
 background-image: linear-gradient(180deg, var(--n-color-start), var(--n-color-start) 50%, transparent 50%, transparent 100%);
 background-size: 1px 10px;
 `)])]),u("timeline-item-timeline",`
 width: calc(var(--n-icon-size) + 12px);
 position: absolute;
 top: calc(var(--n-title-font-size) * ${re} / 2 - var(--n-icon-size) / 2);
 height: 100%;
 `,[v("circle",`
 border: var(--n-circle-border);
 transition:
 background-color .3s var(--n-bezier),
 border-color .3s var(--n-bezier);
 width: var(--n-icon-size);
 height: var(--n-icon-size);
 border-radius: var(--n-icon-size);
 box-sizing: border-box;
 `),v("icon",`
 color: var(--n-icon-color);
 font-size: var(--n-icon-size);
 height: var(--n-icon-size);
 width: var(--n-icon-size);
 display: flex;
 align-items: center;
 justify-content: center;
 `),v("line",`
 transition: background-color .3s var(--n-bezier);
 position: absolute;
 top: var(--n-icon-size);
 left: calc(var(--n-icon-size) / 2 - 1px);
 bottom: 0px;
 width: 2px;
 background-color: var(--n-line-color);
 `)])])]),Ve=Object.assign(Object.assign({},ce.props),{horizontal:Boolean,itemPlacement:{type:String,default:"left"},size:{type:String,default:"medium"},iconSize:Number}),de=fe("n-timeline"),Ae=G({name:"Timeline",props:Ve,setup(o,{slots:c}){const{mergedClsPrefixRef:i}=le(o),d=ce("Timeline","-timeline",Fe,je,o,i);return ve(de,{props:o,mergedThemeRef:d,mergedClsPrefixRef:i}),()=>{const{value:r}=i;return z("div",{class:[`${r}-timeline`,o.horizontal&&`${r}-timeline--horizontal`,`${r}-timeline--${o.size}-size`,!o.horizontal&&`${r}-timeline--${o.itemPlacement}-placement`]},c)}}}),Ee={time:[String,Number],title:String,content:String,color:String,lineType:{type:String,default:"default"},type:{type:String,default:"default"}},Oe=G({name:"TimelineItem",props:Ee,slots:Object,setup(o){const c=he(de);c||_e("timeline-item","`n-timeline-item` must be placed inside `n-timeline`."),Re();const{inlineThemeDisabled:i}=le(),d=N(()=>{const{props:{size:l,iconSize:h},mergedThemeRef:g}=c,{type:_}=o,{self:{titleTextColor:C,contentTextColor:y,metaTextColor:w,lineColor:j,titleFontWeight:F,contentFontSize:V,[R("iconSize",l)]:A,[R("titleMargin",l)]:L,[R("titleFontSize",l)]:E,[R("circleBorder",_)]:O,[R("iconColor",_)]:W},common:{cubicBezierEaseInOut:D}}=g.value;return{"--n-bezier":D,"--n-circle-border":O,"--n-icon-color":W,"--n-content-font-size":V,"--n-content-text-color":y,"--n-line-color":j,"--n-meta-text-color":w,"--n-title-font-size":E,"--n-title-font-weight":F,"--n-title-margin":L,"--n-title-text-color":C,"--n-icon-size":Te(h)||A}}),r=i?ye("timeline-item",N(()=>{const{props:{size:l,iconSize:h}}=c,{type:g}=o;return`${l[0]}${h||"a"}${g[0]}`}),d,c.props):void 0;return{mergedClsPrefix:c.mergedClsPrefixRef,cssVars:i?void 0:d,themeClass:r==null?void 0:r.themeClass,onRender:r==null?void 0:r.onRender}},render(){const{mergedClsPrefix:o,color:c,onRender:i,$slots:d}=this;return i==null||i(),z("div",{class:[`${o}-timeline-item`,this.themeClass,`${o}-timeline-item--${this.type}-type`,`${o}-timeline-item--${this.lineType}-line-type`],style:this.cssVars},z("div",{class:`${o}-timeline-item-timeline`},z("div",{class:`${o}-timeline-item-timeline__line`}),ee(d.icon,r=>r?z("div",{class:`${o}-timeline-item-timeline__icon`,style:{color:c}},r):z("div",{class:`${o}-timeline-item-timeline__circle`,style:{borderColor:c}}))),z("div",{class:`${o}-timeline-item-content`},ee(d.header,r=>r||this.title?z("div",{class:`${o}-timeline-item-content__title`},r||this.title):null),z("div",{class:`${o}-timeline-item-content__content`},te(d.default,()=>[this.content])),z("div",{class:`${o}-timeline-item-content__meta`},te(d.footer,()=>[this.time]))))}}),ae={list:o=>X().get(`/grades/${o}/comments`),create:(o,c)=>X().post(`/grades/${o}/comments`,c)},We={key:0,class:"center"},De={key:1,class:"grade-detail"},Ue={class:"page-header"},He={class:"grade-num"},Ke={class:"header-actions"},qe={class:"grade-layout"},Ge={class:"grade-main"},Qe={class:"card",style:{"margin-bottom":"16px"}},Ye={class:"grade-notes"},Ze={class:"card"},Je={key:0,class:"empty-comments"},Xe={class:"comment-list"},et={class:"comment__header"},tt={class:"comment__date"},it={class:"comment__body"},nt={class:"comment-form"},ot={class:"comment-form__footer"},st={class:"grade-aside"},rt={class:"card"},at={class:"meta-row"},lt={class:"meta-row"},ct={class:"priority-tag"},dt={class:"meta-row"},mt={class:"meta-row"},ut={class:"meta-row"},pt={class:"meta-row"},gt={key:0,class:"meta-row"},ft={class:"score-value"},vt={key:0,class:"card history-card"},ht=G({__name:"[id]",setup(o){const c=ze(),i=we(),d=xe(),r=Be(),l=c.params.id,h=B([]),g=B(""),_=B(!1),C=B(!1),y=B(!1),w=B(null);be(async()=>{await i.fetchById(l);const{data:s}=await ae.list(l);h.value=s});const j=N(()=>{var s,e;return d.isTeacher&&!((s=i.current)!=null&&s.teacher_id)&&((e=i.current)==null?void 0:e.status)==="pending"}),F=N(()=>{var s;return d.isTeacher&&((s=i.current)==null?void 0:s.status)!=="closed"}),V=N(()=>{var s,e,a;return d.isStudent&&((s=i.current)==null?void 0:s.student_id)===((e=d.user)==null?void 0:e.id)&&((a=i.current)==null?void 0:a.status)==="graded"}),A={pending:["graded"],graded:["certified","not_certified","pending"],not_certified:["graded"],certified:["closed"],closed:[]},L={pending:"Ожидает",graded:"Оценено",not_certified:"Не аттестован",certified:"Аттестован",closed:"Завершено"},E=N(()=>{var e;const s=((e=i.current)==null?void 0:e.status)||"";return(A[s]||[]).map(a=>({label:L[a],value:a}))});async function O(){var s,e;y.value=!0;try{await H.assign(l),await i.fetchById(l),r.success("Оценка назначена на вас")}catch(a){r.error(((e=(s=a==null?void 0:a.response)==null?void 0:s.data)==null?void 0:e.error)||"Ошибка")}finally{y.value=!1}}async function W(s){var e,a;try{await H.updateStatus(l,s),await i.fetchById(l),w.value=null,r.success("Статус обновлён")}catch(S){r.error(((a=(e=S==null?void 0:S.response)==null?void 0:e.data)==null?void 0:a.error)||"Ошибка")}}async function D(){var s,e;y.value=!0;try{await H.acknowledge(l),await i.fetchById(l),r.success("Ознакомление подтверждено")}catch(a){r.error(((e=(s=a==null?void 0:a.response)==null?void 0:s.data)==null?void 0:e.error)||"Ошибка")}finally{y.value=!1}}async function me(){var s,e;if(g.value.trim()){C.value=!0;try{const{data:a}=await ae.create(l,{body:g.value,is_internal:_.value});h.value.push(a),g.value="",_.value=!1}catch(a){r.error(((e=(s=a==null?void 0:a.response)==null?void 0:s.data)==null?void 0:e.error)||"Ошибка")}finally{C.value=!1}}}const ue=s=>L[s]||s,U=s=>new Date(s).toLocaleString("ru-RU",{day:"2-digit",month:"2-digit",year:"numeric",hour:"2-digit",minute:"2-digit"});return(s,e)=>{var a,S,Q,Y,Z;return n(i).loading?(p(),b("div",We,[M(n(Ne))])):n(i).current?(p(),b("div",De,[t("div",Ue,[t("div",null,[t("h1",null,f(n(i).current.task_name),1),t("span",He,"#"+f(n(i).current.number),1)]),t("div",Ke,[n(j)?(p(),k(n(q),{key:0,size:"small",loading:n(y),onClick:O},{default:$(()=>[...e[3]||(e[3]=[T(" Взять в работу ",-1)])]),_:1},8,["loading"])):x("",!0),n(F)?(p(),k(n(Se),{key:1,value:n(w),"onUpdate:value":[e[0]||(e[0]=m=>K(w)?w.value=m:null),W],options:n(E),placeholder:"Сменить статус",style:{width:"200px"},size:"small"},null,8,["value","options"])):x("",!0),n(V)?(p(),k(n(q),{key:2,size:"small",type:"primary",loading:n(y),onClick:D},{default:$(()=>[...e[4]||(e[4]=[T(" Подтвердить ознакомление ",-1)])]),_:1},8,["loading"])):x("",!0)])]),t("div",qe,[t("div",Ge,[t("div",Qe,[t("div",Ye,f(n(i).current.notes),1)]),t("div",Ze,[e[8]||(e[8]=t("h3",{class:"section-title"},"Комментарии",-1)),n(h).length?x("",!0):(p(),b("div",Je,"Нет комментариев")),t("div",Xe,[(p(!0),b(ie,null,ne(n(h),m=>{var J;return p(),b("div",{key:m.id,class:oe(["comment",{"comment--internal":m.is_internal}])},[t("div",et,[t("strong",null,f((J=m.author)==null?void 0:J.full_name),1),t("span",tt,f(U(m.created_at)),1),m.is_internal?(p(),k(n(ke),{key:0,size:"small",type:"warning"},{default:$(()=>[...e[5]||(e[5]=[T("Внутренний",-1)])]),_:1})):x("",!0)]),t("div",it,f(m.body),1)],2)}),128))]),t("div",nt,[M(n(Ie),{value:n(g),"onUpdate:value":e[1]||(e[1]=m=>K(g)?g.value=m:null),type:"textarea",rows:3,placeholder:"Написать комментарий..."},null,8,["value"]),t("div",ot,[n(d).isTeacher?(p(),k(n($e),{key:0,checked:n(_),"onUpdate:checked":e[2]||(e[2]=m=>K(_)?_.value=m:null)},{default:$(()=>[...e[6]||(e[6]=[T("Внутренняя заметка",-1)])]),_:1},8,["checked"])):x("",!0),M(n(q),{type:"primary",loading:n(C),disabled:!n(g).trim(),onClick:me},{default:$(()=>[...e[7]||(e[7]=[T(" Отправить ",-1)])]),_:1},8,["loading","disabled"])])])])]),t("aside",st,[t("div",rt,[t("div",at,[e[9]||(e[9]=t("span",{class:"meta-label"},"Статус оценки",-1)),t("span",{class:oe(["status-badge",`status-badge--${n(i).current.status}`])},f(ue(n(i).current.status)),3)]),t("div",lt,[e[10]||(e[10]=t("span",{class:"meta-label"},"Тип оценки",-1)),t("span",ct,[t("span",{class:"dot",style:Ce({background:(a=n(i).current.grade_type)==null?void 0:a.color})},null,4),T(" "+f((S=n(i).current.grade_type)==null?void 0:S.name),1)])]),t("div",dt,[e[11]||(e[11]=t("span",{class:"meta-label"},"Дисциплина",-1)),t("span",null,f((Q=n(i).current.discipline)==null?void 0:Q.name),1)]),t("div",mt,[e[12]||(e[12]=t("span",{class:"meta-label"},"Студент",-1)),t("span",null,f((Y=n(i).current.student)==null?void 0:Y.full_name),1)]),t("div",ut,[e[13]||(e[13]=t("span",{class:"meta-label"},"Преподаватель",-1)),t("span",null,f(((Z=n(i).current.teacher)==null?void 0:Z.full_name)||"—"),1)]),t("div",pt,[e[14]||(e[14]=t("span",{class:"meta-label"},"Добавлено",-1)),t("span",null,f(U(n(i).current.created_at)),1)]),n(i).current.score!=null?(p(),b("div",gt,[e[15]||(e[15]=t("span",{class:"meta-label"},"Балл",-1)),t("span",ft,f(n(i).current.score),1)])):x("",!0)]),n(i).history.length?(p(),b("div",vt,[e[16]||(e[16]=t("h4",{class:"section-title"},"История",-1)),M(n(Ae),null,{default:$(()=>[(p(!0),b(ie,null,ne(n(i).history,m=>(p(),k(n(Oe),{key:m.id,content:`${m.old_value||"—"} → ${m.new_value}`,time:U(m.changed_at)},null,8,["content","time"]))),128))]),_:1})])):x("",!0)])])])):x("",!0)}}}),Bt=Pe(ht,[["__scopeId","data-v-7dc251eb"]]);export{Bt as default};
