(self.webpackChunkidas_ui=self.webpackChunkidas_ui||[]).push([[416],{85064:function(He,Ae,n){"use strict";n.d(Ae,{Z:function(){return ze}});var pe=n(28991),Je=n(28481),or=n(96156),We=n(81253),re=n(67294),Ir=n(94184),se=n.n(Ir),wr=(0,re.createContext)({}),ge=wr,ur=n(90484),Ze=n(86500),Pe=n(1350),Re=2,w=.16,E=.05,f=.05,Tr=.15,R=5,N=4,cr=[{index:7,opacity:.15},{index:6,opacity:.25},{index:5,opacity:.3},{index:5,opacity:.45},{index:5,opacity:.65},{index:5,opacity:.85},{index:4,opacity:.9},{index:3,opacity:.95},{index:2,opacity:.97},{index:1,opacity:.98}];function be(a){var o=a.r,g=a.g,c=a.b,T=(0,Ze.py)(o,g,c);return{h:T.h*360,s:T.s,v:T.v}}function h(a){var o=a.r,g=a.g,c=a.b;return"#".concat((0,Ze.vq)(o,g,c,!1))}function dr(a,o,g){var c=g/100,T={r:(o.r-a.r)*c+a.r,g:(o.g-a.g)*c+a.g,b:(o.b-a.b)*c+a.b};return T}function ne(a,o,g){var c;return Math.round(a.h)>=60&&Math.round(a.h)<=240?c=g?Math.round(a.h)-Re*o:Math.round(a.h)+Re*o:c=g?Math.round(a.h)+Re*o:Math.round(a.h)-Re*o,c<0?c+=360:c>=360&&(c-=360),c}function Ve(a,o,g){if(a.h===0&&a.s===0)return a.s;var c;return g?c=a.s-w*o:o===N?c=a.s+w:c=a.s+E*o,c>1&&(c=1),g&&o===R&&c>.1&&(c=.1),c<.06&&(c=.06),Number(c.toFixed(2))}function Ye(a,o,g){var c;return g?c=a.v+f*o:c=a.v-Tr*o,c>1&&(c=1),Number(c.toFixed(2))}function Oe(a){for(var o=arguments.length>1&&arguments[1]!==void 0?arguments[1]:{},g=[],c=(0,Pe.uA)(a),T=R;T>0;T-=1){var ue=be(c),ce=h((0,Pe.uA)({h:ne(ue,T,!0),s:Ve(ue,T,!0),v:Ye(ue,T,!0)}));g.push(ce)}g.push(h(c));for(var de=1;de<=N;de+=1){var ae=be(c),Ce=h((0,Pe.uA)({h:ne(ae,de),s:Ve(ae,de),v:Ye(ae,de)}));g.push(Ce)}return o.theme==="dark"?cr.map(function(q){var Ue=q.index,Le=q.opacity,ye=h(dr((0,Pe.uA)(o.backgroundColor||"#141414"),(0,Pe.uA)(g[Ue]),Le*100));return ye}):g}var xe={red:"#F5222D",volcano:"#FA541C",orange:"#FA8C16",gold:"#FAAD14",yellow:"#FADB14",lime:"#A0D911",green:"#52C41A",cyan:"#13C2C2",blue:"#1677FF",geekblue:"#2F54EB",purple:"#722ED1",magenta:"#EB2F96",grey:"#666666"},Y={},je={};Object.keys(xe).forEach(function(a){Y[a]=Oe(xe[a]),Y[a].primary=Y[a][5],je[a]=Oe(xe[a],{theme:"dark",backgroundColor:"#141414"}),je[a].primary=je[a][5]});var oe=Y.red,Ar=Y.volcano,Fe=Y.gold,Pr=Y.orange,Rr=Y.yellow,Fr=Y.lime,Nr=Y.green,Be=Y.cyan,kr=Y.blue,Ne=Y.geekblue,J=Y.purple,Qe=Y.magenta,$r=Y.grey,Lr=Y.grey,Xe=n(80334),fr=n(44958);function vr(a,o){(0,Xe.ZP)(a,"[@ant-design/icons] ".concat(o))}function qe(a){return(0,ur.Z)(a)==="object"&&typeof a.name=="string"&&typeof a.theme=="string"&&((0,ur.Z)(a.icon)==="object"||typeof a.icon=="function")}function _e(){var a=arguments.length>0&&arguments[0]!==void 0?arguments[0]:{};return Object.keys(a).reduce(function(o,g){var c=a[g];switch(g){case"class":o.className=c,delete o.class;break;default:o[g]=c}return o},{})}function Ke(a,o,g){return g?re.createElement(a.tag,(0,pe.Z)((0,pe.Z)({key:o},_e(a.attrs)),g),(a.children||[]).map(function(c,T){return Ke(c,"".concat(o,"-").concat(a.tag,"-").concat(T))})):re.createElement(a.tag,(0,pe.Z)({key:o},_e(a.attrs)),(a.children||[]).map(function(c,T){return Ke(c,"".concat(o,"-").concat(a.tag,"-").concat(T))}))}function er(a){return Oe(a)[0]}function rr(a){return a?Array.isArray(a)?a:[a]:[]}var Er={width:"1em",height:"1em",fill:"currentColor","aria-hidden":"true",focusable:"false"},pr=`
.anticon {
  display: inline-block;
  color: inherit;
  font-style: normal;
  line-height: 0;
  text-align: center;
  text-transform: none;
  vertical-align: -0.125em;
  text-rendering: optimizeLegibility;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}

.anticon > * {
  line-height: 1;
}

.anticon svg {
  display: inline-block;
}

.anticon::before {
  display: none;
}

.anticon .anticon-icon {
  display: block;
}

.anticon[tabindex] {
  cursor: pointer;
}

.anticon-spin::before,
.anticon-spin {
  display: inline-block;
  -webkit-animation: loadingCircle 1s infinite linear;
  animation: loadingCircle 1s infinite linear;
}

@-webkit-keyframes loadingCircle {
  100% {
    -webkit-transform: rotate(360deg);
    transform: rotate(360deg);
  }
}

@keyframes loadingCircle {
  100% {
    -webkit-transform: rotate(360deg);
    transform: rotate(360deg);
  }
}
`,mr=function(){var o=arguments.length>0&&arguments[0]!==void 0?arguments[0]:pr,g=(0,re.useContext)(ge),c=g.csp,T=g.prefixCls,ue=o;T&&(ue=ue.replace(/anticon/g,T)),(0,re.useEffect)(function(){(0,fr.hq)(ue,"@ant-design-icons",{prepend:!0,csp:c})},[])},hr=["icon","className","onClick","style","primaryColor","secondaryColor"],Ie={primaryColor:"#333",secondaryColor:"#E6E6E6",calculated:!1};function e(a){var o=a.primaryColor,g=a.secondaryColor;Ie.primaryColor=o,Ie.secondaryColor=g||er(o),Ie.calculated=!!g}function gr(){return(0,pe.Z)({},Ie)}var ke=function(o){var g=o.icon,c=o.className,T=o.onClick,ue=o.style,ce=o.primaryColor,de=o.secondaryColor,ae=(0,We.Z)(o,hr),Ce=Ie;if(ce&&(Ce={primaryColor:ce,secondaryColor:de||er(ce)}),mr(),vr(qe(g),"icon should be icon definiton, but got ".concat(g)),!qe(g))return null;var q=g;return q&&typeof q.icon=="function"&&(q=(0,pe.Z)((0,pe.Z)({},q),{},{icon:q.icon(Ce.primaryColor,Ce.secondaryColor)})),Ke(q.icon,"svg-".concat(q.name),(0,pe.Z)({className:c,onClick:T,style:ue,"data-icon":q.name,width:"1em",height:"1em",fill:"currentColor","aria-hidden":"true"},ae))};ke.displayName="IconReact",ke.getTwoToneColors=gr,ke.setTwoToneColors=e;var Me=ke;function nr(a){var o=rr(a),g=(0,Je.Z)(o,2),c=g[0],T=g[1];return Me.setTwoToneColors({primaryColor:c,secondaryColor:T})}function yr(){var a=Me.getTwoToneColors();return a.calculated?[a.primaryColor,a.secondaryColor]:a.primaryColor}var Zr=["className","icon","spin","rotate","tabIndex","onClick","twoToneColor"];nr("#1890ff");var $e=re.forwardRef(function(a,o){var g,c=a.className,T=a.icon,ue=a.spin,ce=a.rotate,de=a.tabIndex,ae=a.onClick,Ce=a.twoToneColor,q=(0,We.Z)(a,Zr),Ue=re.useContext(ge),Le=Ue.prefixCls,ye=Le===void 0?"anticon":Le,Dr=Ue.rootClassName,Cr=se()(Dr,ye,(g={},(0,or.Z)(g,"".concat(ye,"-").concat(T.name),!!T.name),(0,or.Z)(g,"".concat(ye,"-spin"),!!ue||T.name==="loading"),g),c),ar=de;ar===void 0&&ae&&(ar=-1);var tr=ce?{msTransform:"rotate(".concat(ce,"deg)"),transform:"rotate(".concat(ce,"deg)")}:void 0,Sr=rr(Ce),lr=(0,Je.Z)(Sr,2),br=lr[0],xr=lr[1];return re.createElement("span",(0,pe.Z)((0,pe.Z)({role:"img","aria-label":T.name},q),{},{ref:o,tabIndex:ar,onClick:ae,className:Cr}),re.createElement(Me,{icon:T,primaryColor:br,secondaryColor:xr,style:tr}))});$e.displayName="AntdIcon",$e.getTwoToneColor=yr,$e.setTwoToneColor=nr;var ze=$e},16894:function(He,Ae,n){"use strict";var pe=n(4582);Ae.ZP=pe.Z},23528:function(He){He.exports={SearchInput:"SearchInput___1BqFd"}},2873:function(He,Ae,n){"use strict";n.r(Ae),n.d(Ae,{default:function(){return sn}});var pe=n(57338),Je=n(273),or=n(18106),We=n(58500),re=n(93224),Ir=n(54421),se=n(38272),wr=n(71194),ge=n(34883),ur=n(57663),Ze=n(71577),Pe=n(62350),Re=n(24565),w=n(11849),E=n(2824),f=n(39428),Tr=n(34792),R=n(48086),N=n(3182),cr=n(30381),be=n.n(cr),h=n(67294),dr=n(21010),ne=n(14300),Ve=n(72709),Ye=["id"];function Oe(y,t){return xe.apply(this,arguments)}function xe(){return xe=(0,N.Z)((0,f.Z)().mark(function y(t,i){return(0,f.Z)().wrap(function(u){for(;;)switch(u.prev=u.next){case 0:return u.abrupt("return",(0,Ve.WY)("/api/v1/sessions",(0,w.Z)({method:"GET",params:(0,w.Z)({},t)},i||{})));case 1:case"end":return u.stop()}},y)})),xe.apply(this,arguments)}function Y(y,t){return je.apply(this,arguments)}function je(){return je=(0,N.Z)((0,f.Z)().mark(function y(t,i){var s,u;return(0,f.Z)().wrap(function(p){for(;;)switch(p.prev=p.next){case 0:return s=t.id,u=(0,re.Z)(t,Ye),p.abrupt("return",(0,Ve.WY)("/api/v1/sessions/".concat(s),(0,w.Z)({method:"DELETE",params:(0,w.Z)({},u)},i||{})));case 2:case"end":return p.stop()}},y)})),je.apply(this,arguments)}var oe=n(10861),Ar=n(88293),Fe=n(93400),Pr=n(55035),Rr=n(28508),Fr=n(79508),Nr=n(51042),Be=n(57119),kr=n(81253),Ne=n(96156),J=n(28991),Qe=n(28481),$r=n(95985),Lr=n(85064),Xe=function(t,i){return h.createElement(Lr.Z,(0,J.Z)((0,J.Z)({},t),{},{ref:i,icon:$r.Z}))};Xe.displayName="FilterOutlined";var fr=h.forwardRef(Xe),vr=n(17503),qe=n(93035),_e=n(92124),Ke=n(88182),er=n(94184),rr=n.n(er),Er=n(97435),pr=n(42489),mr=n(27350),hr=function(t){return(0,Ne.Z)({},t.componentCls,{lineHeight:"30px","&::before":{display:"block",height:0,visibility:"hidden",content:"'.'"},"&-small":{lineHeight:t.lineHeight},"&-container":{display:"flex",flexWrap:"wrap",gap:8},"&-item":(0,Ne.Z)({whiteSpace:"nowrap"},"".concat(t.antCls,"-form-item"),{marginBlock:0}),"&-line":{minWidth:"198px"},"&-line:not(:first-child)":{marginBlockStart:"16px",marginBlockEnd:8},"&-collapse-icon":{width:t.controlHeight,height:t.controlHeight,borderRadius:"50%",display:"flex",alignItems:"center",justifyContent:"center"},"&-effective":(0,Ne.Z)({},"".concat(t.componentCls,"-collapse-icon"),{backgroundColor:t.colorBgTextHover})})};function Ie(y){return(0,mr.Xj)("LightFilter",function(t){var i=(0,J.Z)((0,J.Z)({},t),{},{componentCls:".".concat(y)});return[hr(i)]})}var e=n(85893),gr=["size","collapse","collapseLabel","initialValues","onValuesChange","form","placement","formRef","bordered","ignoreRules","footerRender"],ke=function(t){var i=t.items,s=t.prefixCls,u=t.size,r=u===void 0?"middle":u,p=t.collapse,B=t.collapseLabel,K=t.onValuesChange,d=t.bordered,z=t.values,U=t.footerRender,D=t.placement,V=(0,vr.YB)(),F="".concat(s,"-light-filter"),k=Ie(F),G=k.wrapSSR,C=k.hashId,W=(0,h.useState)(!1),x=(0,Qe.Z)(W,2),M=x[0],me=x[1],v=(0,h.useState)(function(){return(0,J.Z)({},z)}),O=(0,Qe.Z)(v,2),_=O[0],fe=O[1];(0,h.useEffect)(function(){fe((0,J.Z)({},z))},[z]);var H=(0,h.useMemo)(function(){var L=[],ee=[];return i.forEach(function(te){var P=te.props||{},j=P.secondary;j||p?L.push(te):ee.push(te)}),{collapseItems:L,outsideItems:ee}},[t.items]),ie=H.collapseItems,Ee=H.outsideItems,Se=function(){return B||(p?(0,e.jsx)(fr,{className:"".concat(F,"-collapse-icon ").concat(C)}):(0,e.jsx)(qe.Q,{size:r,label:V.getMessage("form.lightFilter.more","\u66F4\u591A\u7B5B\u9009"),expanded:M}))};return G((0,e.jsx)("div",{className:rr()(F,C,"".concat(F,"-").concat(r),(0,Ne.Z)({},"".concat(F,"-effective"),Object.keys(z).some(function(L){return z[L]}))),children:(0,e.jsxs)("div",{className:"".concat(F,"-container ").concat(C),children:[Ee.map(function(L,ee){var te=L.key,P=L.props.fieldProps,j=P!=null&&P.placement?P==null?void 0:P.placement:D;return(0,e.jsx)("div",{className:"".concat(F,"-item ").concat(C),children:h.cloneElement(L,{fieldProps:(0,J.Z)((0,J.Z)({},L.props.fieldProps),{},{placement:j}),proFieldProps:{light:!0,label:L.props.label,bordered:d},bordered:d})},te||ee)}),ie.length?(0,e.jsx)("div",{className:"".concat(F,"-item ").concat(C),children:(0,e.jsx)(_e.M,{padding:24,open:M,onOpenChange:function(ee){me(ee)},placement:D,label:Se(),footerRender:U,footer:{onConfirm:function(){K((0,J.Z)({},_)),me(!1)},onClear:function(){var ee={};ie.forEach(function(te){var P=te.props.name;ee[P]=void 0}),K(ee)}},children:ie.map(function(L){var ee=L.key,te=L.props,P=te.name,j=te.fieldProps,ve=(0,J.Z)((0,J.Z)({},j),{},{onChange:function(he){return fe((0,J.Z)((0,J.Z)({},_),{},(0,Ne.Z)({},P,he!=null&&he.target?he.target.value:he))),!1}});_.hasOwnProperty(P)&&(ve[L.props.valuePropName||"value"]=_[P]);var le=j!=null&&j.placement?j==null?void 0:j.placement:D;return(0,e.jsx)("div",{className:"".concat(F,"-line ").concat(C),children:h.cloneElement(L,{fieldProps:(0,J.Z)((0,J.Z)({},ve),{},{placement:le})})},ee)})})},"more"):null]})}))};function Me(y){var t=y.size,i=y.collapse,s=y.collapseLabel,u=y.initialValues,r=y.onValuesChange,p=y.form,B=y.placement,K=y.formRef,d=y.bordered,z=y.ignoreRules,U=y.footerRender,D=(0,kr.Z)(y,gr),V=(0,h.useContext)(Ke.ZP.ConfigContext),F=V.getPrefixCls,k=F("pro-form"),G=(0,h.useState)(function(){return(0,J.Z)({},u)}),C=(0,Qe.Z)(G,2),W=C[0],x=C[1],M=(0,h.useRef)();return(0,h.useImperativeHandle)(K,function(){return M.current}),(0,e.jsx)(pr.I,(0,J.Z)((0,J.Z)({size:t,initialValues:u,form:p,contentRender:function(v){return(0,e.jsx)(ke,{prefixCls:k,items:v.flatMap(function(O){return(O==null?void 0:O.type.displayName)==="ProForm-Group"?O.props.children:O}),size:t,bordered:d,collapse:i,collapseLabel:s,placement:B,values:W||{},footerRender:U,onValuesChange:function(_){var fe,H,ie=(0,J.Z)((0,J.Z)({},W),_);x(ie),(fe=M.current)===null||fe===void 0||fe.setFieldsValue(ie),(H=M.current)===null||H===void 0||H.submit(),r&&r(_,ie)}})},formRef:M,formItemProps:{colon:!1,labelAlign:"left"},fieldProps:{style:{width:void 0}}},(0,Er.Z)(D,["labelWidth"])),{},{onValuesChange:function(v,O){var _;x(O),r==null||r(v,O),(_=M.current)===null||_===void 0||_.submit()}}))}var nr=n(64317),yr=n(55997),Zr=n(43037),$e=n(85224),ze=n(16894),a=n(71153),o=n(60331),g=n(47673),c=n(28166),T=function(t){var i=t.visible,s=t.onCancel,u=t.user,r=t.parentIntl,p=new Fe.f("keypair.form",r),B=(0,h.useState)(!1),K=(0,E.Z)(B,2),d=K[0],z=K[1],U=(0,h.useState)(),D=(0,E.Z)(U,2),V=D[0],F=D[1],k=(0,h.useState)(""),G=(0,E.Z)(k,2),C=G[0],W=G[1];return(0,h.useEffect)(function(){i&&u&&(z(!1),F(void 0),W(""))},[u,i]),(0,e.jsx)(e.Fragment,{children:(0,e.jsxs)(ge.Z,{title:V?p.t("title.success","Success"):p.t("title","Add User Key"),open:i,onCancel:s,cancelButtonProps:{hidden:Boolean(V)},onOk:(0,N.Z)((0,f.Z)().mark(function x(){var M;return(0,f.Z)().wrap(function(v){for(;;)switch(v.prev=v.next){case 0:if(!(u!=null&&u.id&&!d)){v.next=12;break}if(C){v.next=4;break}return R.default.error(p.t("name.empty","please input name!")),v.abrupt("return");case 4:return v.next=6,(0,oe.bJ)({userId:u.id},{name:C,userId:u.id});case 6:if(M=v.sent,!M.success){v.next=11;break}return F(M.data),z(!0),v.abrupt("return");case 11:return v.abrupt("return");case 12:s();case 13:case"end":return v.stop()}},x)})),children:[(0,e.jsx)(c.Z,{hidden:d,value:C,onChange:function(M){W(M.target.value)},required:!d,placeholder:p.t("name.placeholder","Please Input name or description of user key")}),V&&(0,e.jsxs)("div",{children:[p.t("tips","User key pair, please keep it properly"),":",(0,e.jsxs)("ul",{style:{padding:"0px"},children:[(0,e.jsxs)("li",{children:["AccessKey: ",(0,e.jsx)(o.Z,{children:V.key})]}),(0,e.jsxs)("li",{children:["AccessSecret: ",(0,e.jsx)(o.Z,{children:V.secret})]}),(0,e.jsxs)("li",{children:["Private: ",(0,e.jsx)(o.Z,{children:V.private})]})]})]})]})})},ue=T,ce=n(2844),de=n(62448),ae=n(5966),Ce=n(48736),q=n(27049),Ue=n(71748),Le=n(90860),ye=n(86582),Dr=n(43358),Cr=n(66074),ar=n(402),tr=n(97272),Sr=n(58533),lr=n(44728),br=n(23528),xr=n.n(br),Mr=["type","loading"],zr=function(t){var i=t.onChange,s=t.loading,u=t.granting,r=t.parentIntl,p=t.apps,B=new Fe.f("apps",r);return(0,e.jsx)(se.ZP,{dataSource:p,size:"small",loading:s,renderItem:function(d){var z,U,D,V,F;return(0,e.jsxs)(se.ZP.Item,{actions:i&&u?[(0,e.jsx)("a",{onClick:function(){i(p.filter(function(G){return G.id!=d.id}))},children:B.t("appList.delete","Delete")},"grant")]:[],children:[(0,e.jsx)(se.ZP.Item.Meta,{avatar:(0,e.jsx)(ce.ZP,{src:"".concat(d.avatar)}),title:(z=d.displayName)!==null&&z!==void 0?z:d.name,description:(0,e.jsx)(tr.Z.Paragraph,{type:"secondary",ellipsis:{tooltip:d.description},children:(U=d.description)!==null&&U!==void 0?U:""})}),i&&u&&d.roles&&d.roles.length>0?(0,e.jsx)(Cr.Z,{onSelect:function(G){i(p.map(function(C){return C.id==d.id?(0,w.Z)((0,w.Z)({},C),{},{roleId:G}):C}))},defaultValue:d&&d.roleId?d.roleId:(D=d.roles.find(function(k){return k.isDefault}))===null||D===void 0?void 0:D.id,options:d.roles.map(function(k){return{key:k.id,value:k.id,label:k.name}})}):(0,e.jsx)("div",{children:d&&d.role?d.role:(V=d.roles)===null||V===void 0||(F=V.find(function(k){return k.isDefault}))===null||F===void 0?void 0:F.name})]},d.id)}})},Ur=function(t){var i=t.type,s=i===void 0?"vertical":i,u=t.loading,r=(0,re.Z)(t,Mr),p=r.onChange,B=r.apps,K=r.granting,d=r.parentIntl,z=new Fe.f("grant",d),U=(0,h.useState)(!0),D=(0,E.Z)(U,2),V=D[0],F=D[1],k=(0,h.useState)(0),G=(0,E.Z)(k,2),C=G[0],W=G[1],x=(0,h.useState)(),M=(0,E.Z)(x,2),me=M[0],v=M[1],O=(0,h.useState)([]),_=(0,E.Z)(O,2),fe=_[0],H=_[1],ie=(0,h.useState)(!1),Ee=(0,E.Z)(ie,2),Se=Ee[0],L=Ee[1],ee=function(){var P=(0,N.Z)((0,f.Z)().mark(function j(ve){var le,we,he,Ge,sr;return(0,f.Z)().wrap(function(Q){for(;;)switch(Q.prev=Q.next){case 0:return Q.prev=0,L(!0),Q.next=4,(0,lr.C6)((0,w.Z)({current:C+1,pageSize:20,keywords:me},ve));case 4:le=Q.sent,le&&le.data?(we=le.data,he=le.current,Ge=le.pageSize,sr=le.total,H(function(ir){return[].concat((0,ye.Z)(ir),(0,ye.Z)(we))}),W(he),sr<he*Ge&&F(!1)):F(!1),Q.next=11;break;case 8:Q.prev=8,Q.t0=Q.catch(0),Q.t0.handled||console.error("failed to get app list: ".concat(Q.t0));case 11:return Q.prev=11,L(!1),Q.finish(11);case 14:case"end":return Q.stop()}},j,null,[[0,8,11,14]])}));return function(ve){return P.apply(this,arguments)}}(),te=function(j){H([]),v(j),W(1),ee({keywords:j,current:1})};return(0,h.useEffect)(function(){te()},[]),(0,e.jsxs)("div",{style:{height:s==="vertical"?"calc((100vh - 300px))":"100%",width:"100%"},children:[K&&(0,e.jsxs)(e.Fragment,{children:[(0,e.jsxs)("div",{style:{height:s==="vertical"?"calc(100% / 2)":"100%",display:K?"block":"none",width:s==="vertical"?"100%":"calc((100% - 32px) / 2)",float:s==="horizontal"?"left":"unset"},children:[(0,e.jsx)(c.Z.Search,{className:xr().SearchInput,onSearch:function(j){Se||te(j)},loading:Se}),(0,e.jsx)("div",{id:"scrollableDiv",style:{height:"calc( 100% - 32px )",overflow:"auto",padding:"0 16px",border:"1px solid rgba(140, 140, 140, 0.35)"},children:(0,e.jsx)(Sr.Z,{dataLength:fe.length,next:ee,hasMore:V,loader:(0,e.jsx)(Le.Z,{avatar:!0,paragraph:{rows:1},active:!0}),endMessage:(0,e.jsx)(q.Z,{plain:!0,children:"End"}),scrollableTarget:"scrollableDiv",children:(0,e.jsx)(se.ZP,{dataSource:fe.filter(function(P){return!B.map(function(j){return j.id}).includes(P.id)}),size:"small",loading:Se,renderItem:function(j){var ve,le;return(0,e.jsx)(se.ZP.Item,{actions:[(0,e.jsx)("a",{onClick:function(){p([].concat((0,ye.Z)(B),[(0,w.Z)((0,w.Z)({},j),{},{roleId:""})]))},children:z.t("appList.grant","Grant")},"grant")],children:(0,e.jsx)(se.ZP.Item.Meta,{avatar:(0,e.jsx)(ce.ZP,{src:"".concat(j.avatar)}),title:(ve=j.displayName)!==null&&ve!==void 0?ve:j.name,description:(0,e.jsx)(tr.Z.Paragraph,{type:"secondary",ellipsis:{tooltip:j.description},children:(le=j.description)!==null&&le!==void 0?le:""})})},j.id)}})})})]}),(0,e.jsx)(q.Z,{type:s==="vertical"?"horizontal":"vertical",style:{height:s==="vertical"?"unset":"100%",margin:s==="vertical"?"10px 0":"0 10px",float:s==="horizontal"?"left":"unset",display:K?"":"none"}})]}),(0,e.jsx)("div",{id:"scrollableDiv",style:{height:s!=="vertical"||!K?"100%":"calc((100% - 42px) / 2)",width:s==="horizontal"&&K?"calc((100% - 32px) / 2)":"100%",overflow:"auto",padding:"0 16px",float:s==="horizontal"?"left":"unset",border:"1px solid rgba(140, 140, 140, 0.35)"},children:(0,e.jsx)(zr,{parentIntl:d,apps:B.map(function(P){var j;return(0,w.Z)((0,w.Z)({},(j=fe.find(function(ve){return ve.id==P.id}))!==null&&j!==void 0?j:P),{},{roleId:P.roleId,role:P.role})}),granting:K,loading:u,onChange:p})})]})},Vr=Ur,Gr=["parentIntl"],Hr=function(t){var i=t.parentIntl,s=(0,re.Z)(t,Gr),u=new Fe.f("form",i),r=(0,h.useState)(0),p=(0,E.Z)(r,2),B=p[0],K=p[1],d=s.values,z=s.title,U=s.modalVisible,D=s.onSubmit,V=s.onCancel,F=(0,h.useState)([]),k=(0,E.Z)(F,2),G=k[0],C=k[1];return(0,h.useEffect)(function(){K(0),d!=null&&d.id&&U&&(0,oe.bG)({id:d.id}).then(function(W){var x;W.success&&(x=W.data)!==null&&x!==void 0&&x.apps&&C(W.data.apps)})},[d,U]),(0,e.jsxs)(de.L,{stepsProps:{size:"small"},formProps:{preserve:!1},current:B,onCurrentChange:K,stepsFormRender:function(x,M){return(0,e.jsx)(ge.Z,{width:640,bodyStyle:{padding:"32px 40px 48px"},destroyOnClose:!0,title:z,open:U,footer:M,onCancel:function(){ge.Z.confirm({title:u.t("cancel?","Cancel editing?"),icon:(0,e.jsx)(Be.Z,{}),onOk:function(){V()},maskClosable:!0})},children:x})},onFinish:function(){var W=(0,N.Z)((0,f.Z)().mark(function x(M){return(0,f.Z)().wrap(function(v){for(;;)switch(v.prev=v.next){case 0:return v.abrupt("return",D((0,w.Z)((0,w.Z)({},M),{},{apps:G.map(function(O){return{id:O.id,roleId:O.roleId}}),status:d!=null&&d.status?d.status:ne.J0.normal,isDelete:d!=null&&d.isDelete?d==null?void 0:d.isDelete:!1})));case 1:case"end":return v.stop()}},x)}));return function(x){return W.apply(this,arguments)}}(),children:[(0,e.jsxs)(de.L.StepForm,{initialValues:(0,w.Z)({},d),labelCol:{span:8},wrapperCol:{span:14},layout:"horizontal",title:u.t("title.basicConfig","Basic"),className:"basic-view",children:[(0,e.jsx)(ce.HZ,{label:u.t("avatar.label","Avatar"),name:"avatar"}),(0,e.jsx)(ae.Z,{hidden:!0,name:"id"}),(0,e.jsx)(ae.Z,{hidden:!0,name:"storage"}),(0,e.jsx)(ae.Z,{name:"username",label:u.t("userName.label","Username"),width:"md",rules:[{required:!0,message:u.t("userName.required","Please input username!")},{pattern:/^[-_A-Za-z0-9]+$/,message:u.t("name.invalid","Username format error!")}]}),(0,e.jsx)(ae.Z,{name:"fullName",label:u.t("fullName.label","FullName"),width:"md"}),(0,e.jsx)(ae.Z,{name:"email",label:u.t("email.label","Email"),width:"md"}),(0,e.jsx)(ae.Z,{name:"phoneNumber",label:u.t("phoneNumber.label","Telephone number"),width:"md"})]}),(0,e.jsx)(de.L.StepForm,{className:"grant-view",initialValues:{},title:u.t("app.title","App"),children:(0,e.jsx)(Vr,{granting:!0,parentIntl:u,apps:G,onChange:C})})]})},Jr=Hr,Wr=n(96486),Yr=["success","failed","onClick"],Qr=function(t){var i=t.success,s=t.failed,u=t.onClick,r=(0,re.Z)(t,Yr),p=(0,h.useState)(!1),B=(0,E.Z)(p,2),K=B[0],d=B[1],z=(0,h.useState)(0),U=(0,E.Z)(z,2),D=U[0],V=U[1];return(0,e.jsx)(Ze.Z,(0,w.Z)((0,w.Z)({},r),{},{loading:K,onClick:function(){var F=(0,N.Z)((0,f.Z)().mark(function k(G){var C;return(0,f.Z)().wrap(function(x){for(;;)switch(x.prev=x.next){case 0:if(!u){x.next=16;break}return x.prev=1,d(!0),x.next=5,u(G);case 5:C=x.sent,(0,Wr.isBoolean)(C)&&!C?V(2):V(1),x.next=13;break;case 9:x.prev=9,x.t0=x.catch(1),R.default.error("".concat(x.t0),3),V(2);case 13:return x.prev=13,d(!1),x.finish(13);case 16:case"end":return x.stop()}},k,null,[[1,9,13,16]])}));return function(k){return F.apply(this,arguments)}}(),children:D===1&&i?i:D===2&&s?s:r.children}))},Xr=Qr,qr=["createTime","loginTime","extendedData","updateTime"],_r=function(){var y=(0,N.Z)((0,f.Z)().mark(function t(i){var s;return(0,f.Z)().wrap(function(r){for(;;)switch(r.prev=r.next){case 0:return s=R.default.loading("Adding ..."),r.prev=1,delete i.id,r.next=5,(0,oe.r4)(i);case 5:return s(),R.default.success("Added successfully"),r.abrupt("return",!0);case 10:return r.prev=10,r.t0=r.catch(1),s(),R.default.error("Adding failed, please try again!"),r.abrupt("return",!1);case 15:case"end":return r.stop()}},t,null,[[1,10]])}));return function(i){return y.apply(this,arguments)}}(),Or=function(){var y=(0,N.Z)((0,f.Z)().mark(function t(i){var s;return(0,f.Z)().wrap(function(r){for(;;)switch(r.prev=r.next){case 0:if(s=R.default.loading("Configuring"),r.prev=1,!i.id){r.next=10;break}return r.next=5,(0,oe.Nq)({id:i.id},i);case 5:return s(),R.default.success("update is successful"),r.abrupt("return",!0);case 10:return R.default.success("update failed, system error"),r.abrupt("return",!1);case 12:r.next=19;break;case 14:return r.prev=14,r.t0=r.catch(1),s(),R.default.error("Configuration failed, please try again!"),r.abrupt("return",!1);case 19:case"end":return r.stop()}},t,null,[[1,14]])}));return function(i){return y.apply(this,arguments)}}(),en=function(){var y=(0,N.Z)((0,f.Z)().mark(function t(i){var s;return(0,f.Z)().wrap(function(r){for(;;)switch(r.prev=r.next){case 0:if(s=R.default.loading("enabling ..."),i){r.next=3;break}return r.abrupt("return",!0);case 3:return r.prev=3,r.next=6,(0,oe.Q8)(i.map(function(p){return{id:p.id,status:ne.J0.normal}}));case 6:return s(),R.default.success("Enabled successfully and will refresh soon"),r.abrupt("return",!0);case 11:return r.prev=11,r.t0=r.catch(3),s(),R.default.error("Enabled failed, please try again"),r.abrupt("return",!1);case 16:case"end":return r.stop()}},t,null,[[3,11]])}));return function(i){return y.apply(this,arguments)}}(),rn=function(){var y=(0,N.Z)((0,f.Z)().mark(function t(i){var s;return(0,f.Z)().wrap(function(r){for(;;)switch(r.prev=r.next){case 0:if(s=R.default.loading("Disabling ..."),i){r.next=3;break}return r.abrupt("return",!0);case 3:return r.prev=3,r.next=6,(0,oe.Q8)(i.map(function(p){return{id:p.id,status:ne.J0.disabled}}));case 6:return s(),R.default.success("Disabled successfully and will refresh soon"),r.abrupt("return",!0);case 11:return r.prev=11,r.t0=r.catch(3),s(),R.default.error("Disable failed, please try again"),r.abrupt("return",!1);case 16:case"end":return r.stop()}},t,null,[[3,11]])}));return function(i){return y.apply(this,arguments)}}(),nn=function(){var y=(0,N.Z)((0,f.Z)().mark(function t(i){var s;return(0,f.Z)().wrap(function(r){for(;;)switch(r.prev=r.next){case 0:if(s=R.default.loading("Deleting ..."),i){r.next=3;break}return r.abrupt("return",!0);case 3:return r.prev=3,r.next=6,(0,oe.Vt)(i.map(function(p){return{id:p.id}}));case 6:return s(),R.default.success("Deleted successfully and will refresh soon"),r.abrupt("return",!0);case 11:return r.prev=11,r.t0=r.catch(3),s(),R.default.error("Delete failed, please try again: ".concat(r.t0)),r.abrupt("return",!1);case 16:case"end":return r.stop()}},t,null,[[3,11]])}));return function(i){return y.apply(this,arguments)}}(),an=function(){var y=(0,N.Z)((0,f.Z)().mark(function t(i){var s;return(0,f.Z)().wrap(function(r){for(;;)switch(r.prev=r.next){case 0:if(i){r.next=2;break}return r.abrupt("return",!0);case 2:return s=R.default.loading("Deleting ..."),r.prev=3,r.next=6,Y({id:i});case 6:return s(),R.default.success("Deleted successfully and will refresh soon"),r.abrupt("return",!0);case 11:return r.prev=11,r.t0=r.catch(3),s(),R.default.error("Delete failed, please try again"),r.abrupt("return",!1);case 16:case"end":return r.stop()}},t,null,[[3,11]])}));return function(i){return y.apply(this,arguments)}}(),tn=function(){var y=(0,N.Z)((0,f.Z)().mark(function t(i,s){var u;return(0,f.Z)().wrap(function(p){for(;;)switch(p.prev=p.next){case 0:if(!(!i||!s)){p.next=2;break}return p.abrupt("return",!0);case 2:return u=R.default.loading("Deleting ..."),p.prev=3,p.next=6,(0,oe.OK)({id:i,userId:s});case 6:return u(),R.default.success("Deleted successfully and will refresh soon"),p.abrupt("return",!0);case 11:return p.prev=11,p.t0=p.catch(3),u(),R.default.error("Delete failed, please try again"),p.abrupt("return",!1);case 16:case"end":return p.stop()}},t,null,[[3,11]])}));return function(i,s){return y.apply(this,arguments)}}(),ln=function(){var t=(0,h.useState)(!1),i=(0,E.Z)(t,2),s=i[0],u=i[1],r=(0,h.useState)(!1),p=(0,E.Z)(r,2),B=p[0],K=p[1],d=(0,h.useState)(!1),z=(0,E.Z)(d,2),U=z[0],D=z[1],V=(0,h.useState)("sessions"),F=(0,E.Z)(V,2),k=F[0],G=F[1],C=(0,h.useRef)(),W=(0,h.useRef)(),x=(0,h.useRef)(),M=(0,h.useState)(),me=(0,E.Z)(M,2),v=me[0],O=me[1],_=(0,h.useState)([]),fe=(0,E.Z)(_,2),H=fe[0],ie=fe[1],Ee=(0,h.useState)("all"),Se=(0,E.Z)(Ee,2),L=Se[0],ee=Se[1],te=(0,h.useState)(),P=(0,E.Z)(te,2),j=P[0],ve=P[1],le=(0,h.useState)(!1),we=(0,E.Z)(le,2),he=we[0],Ge=we[1],sr=(0,h.useState)([]),jr=(0,E.Z)(sr,2),Q=jr[0],ir=jr[1],m=new Fe.f("pages.users",(0,dr.YB)()),on=function(){var b=(0,N.Z)((0,f.Z)().mark(function Z(l){var S;return(0,f.Z)().wrap(function($){for(;;)switch($.prev=$.next){case 0:return S=(0,w.Z)((0,w.Z)({},l),{},{status:L!=="all"?L:void 0}),j&&(S=(0,w.Z)({keywords:j},S)),$.abrupt("return",(0,oe.Rf)(S));case 3:case"end":return $.stop()}},Z)}));return function(l){return b.apply(this,arguments)}}();(0,h.useEffect)(function(){B&&(G("sessions"),D(!1),v!=null&&v.id&&(0,oe.bG)({id:v.id}).then(function(b){var Z;b.success&&(Z=b.data)!==null&&Z!==void 0&&Z.apps&&ir(b.data.apps)}))},[v,B]);var un=[{title:m.t("session.title.lastSeen","Last Seen"),dataIndex:"lastSeen",render:function(Z,l){return be()(l.lastSeen).locale(m.locale).fromNow()}},{title:m.t("session.title.expiry","Expiry"),dataIndex:"expiry",render:function(Z,l){return be()(l.expiry).locale(m.locale).fromNow()}},{title:m.t("session.title.loggedOn","Logged on"),dataIndex:"createTime",render:function(Z,l){return be()(l.createTime).locale(m.locale).format("LLL")}},{render:function(Z,l){return[(0,e.jsx)("a",{onClick:(0,N.Z)((0,f.Z)().mark(function S(){var A;return(0,f.Z)().wrap(function(I){for(;;)switch(I.prev=I.next){case 0:return I.next=2,an(l.id);case 2:if(!I.sent){I.next=4;break}(A=W.current)===null||A===void 0||A.reload();case 4:case"end":return I.stop()}},S)})),children:(0,e.jsx)(Pr.Z,{})},"delete")]}}],cn=[{title:m.t("keypair.title.name","Name"),dataIndex:"name",ellipsis:!0},{title:m.t("keypair.title.key","Key"),dataIndex:"key",ellipsis:!0},{title:m.t("keypair.title.createTime","Create Time"),dataIndex:"createTime",width:160,render:function(Z,l){return be()(l.createTime).local().format("LLL")}},{width:40,render:function(Z,l){return[(0,e.jsx)(Re.Z,{title:m.t("keypair.delete.popconfirm","Are you sure you want to delete the key named {name}? After deletion, the service cannot be accessed using this key pair.",void 0,{name:l.name}),onConfirm:(0,N.Z)((0,f.Z)().mark(function S(){var A;return(0,f.Z)().wrap(function(I){for(;;)switch(I.prev=I.next){case 0:return I.next=2,tn(l.id,l.userId);case 2:if(!I.sent){I.next=4;break}(A=x.current)===null||A===void 0||A.reload();case 4:case"end":return I.stop()}},S)})),children:(0,e.jsx)("a",{children:(0,e.jsx)(Rr.Z,{})})},"delete")]}}],Br=(0,Ar.GG)(ne.J0,m,"status.value",{normal:"Success",disable:"Error",inactive:"Warning"}),Kr=[{title:m.t("updateForm.userName.nameLabel","User name"),hideInSearch:!0,dataIndex:"username",render:function(Z,l){return(0,e.jsx)("a",{onClick:function(){O(l),K(!0)},children:Z})}},{title:m.t("title.fullName","FullName"),dataIndex:"fullName",hideInSearch:!0},{title:m.t("title.phoneNumber","Telephone number"),dataIndex:"phoneNumber",hideInSearch:!0},{title:m.t("title.email","Email"),dataIndex:"email",hideInSearch:!0},{title:m.t("title.status","Status"),dataIndex:"status",hideInForm:!0,valueEnum:Br},{title:m.t("title.updatedTime","Last update time"),dataIndex:"updateTime",valueType:"dateTime",hideInTable:!0,hideInSearch:!0,hideInForm:!0},{title:m.t("title.loginTime","Last login time"),dataIndex:"loginTime",valueType:"dateTime",hideInSearch:!0,hideInForm:!0},{title:m.t("title.createTime","Create time"),dataIndex:"createTime",valueType:"dateTime",hideInSearch:!0,hideInTable:!0,hideInForm:!0},{title:m.t("title.option","Operating"),dataIndex:"option",valueType:"option",render:function(Z,l){var S=[(0,e.jsx)("a",{style:{flex:"unset"},onClick:function(){u(!0),O(l)},children:m.t("button.edit","Edit")},"config"),(0,e.jsx)("a",{onClick:function(){Ge(!0),O(l)},style:{flex:"unset"},hidden:l.status!==ne.J0.normal,children:m.t("button.addkey","Add key pair")},"addKey"),(0,e.jsx)(Xr,{type:"link",style:{flex:"unset"},success:(0,e.jsx)(Fr.Z,{color:"green"}),onClick:(0,N.Z)((0,f.Z)().mark(function A(){var $;return(0,f.Z)().wrap(function(X){for(;;)switch(X.prev=X.next){case 0:if(l.email){X.next=2;break}throw new Error(m.t("activate.no-email"," The user has no email."));case 2:return X.next=4,(0,oe.Pw)({userId:l.id});case 4:if($=X.sent,!$.success){X.next=9;break}R.default.success(m.t("activate.succcess","Email sent successfully.")),X.next=10;break;case 9:throw new Error("Email sent failed.");case 10:case"end":return X.stop()}},A)})),hidden:l.status!==ne.J0.user_inactive,children:m.t("button.activate","Activate")},"activate")];return S}}];return(0,e.jsxs)(Zr.Oc,{children:[(0,e.jsx)(ze.ZP,{actionRef:C,rowKey:"id",search:!1,toolbar:{search:{onSearch:function(Z){if(ve(Z),C.current){var l,S;(l=(S=C.current).setPageInfo)===null||l===void 0||l.call(S,(0,w.Z)((0,w.Z)({},C.current.pageInfo),{},{current:1})),C.current.reload()}}},filter:(0,e.jsx)(Me,{onFinish:function(){var b=(0,N.Z)((0,f.Z)().mark(function Z(l){var S,A,$,I,X;return(0,f.Z)().wrap(function(Te){for(;;)switch(Te.prev=Te.next){case 0:return A=l.status,ee((S=A.value)!==null&&S!==void 0?S:"all"),C.current&&(($=(I=C.current).setPageInfo)===null||$===void 0||$.call(I,(0,w.Z)((0,w.Z)({},C.current.pageInfo),{},{current:1})),(X=C.current)===null||X===void 0||X.reload()),Te.abrupt("return",!0);case 4:case"end":return Te.stop()}},Z)}));return function(Z){return b.apply(this,arguments)}}(),children:(0,e.jsx)(nr.Z,{name:"status",label:m.t("title.status","Status"),initialValue:{value:"all",label:m.t("status.all","All")},fieldProps:{labelInValue:!0},valueEnum:(0,w.Z)({all:m.t("status.all","All")},Br)})}),actions:[(0,e.jsxs)(Ze.Z,{type:"primary",onClick:function(){u(!0)},children:[(0,e.jsx)(Nr.Z,{}),m.t("button.create","Create")]},"create")]},request:on,columns:Kr,rowSelection:{onChange:function(Z,l){ie(l)}}}),(H==null?void 0:H.length)>0&&(0,e.jsxs)($e.Z,{extra:(0,e.jsxs)("div",{children:[m.t("chosen","Chosen")," ",(0,e.jsx)("a",{style:{fontWeight:600},children:H.length})," ",m.t("item","Item(s)")]}),children:[(0,e.jsx)(Ze.Z,{danger:!0,onClick:function(){ge.Z.confirm({title:m.t("deleteConfirm","Are you sure you want to delete the following users?            "),icon:(0,e.jsx)(Be.Z,{}),onOk:function(){return(0,N.Z)((0,f.Z)().mark(function l(){var S,A;return(0,f.Z)().wrap(function(I){for(;;)switch(I.prev=I.next){case 0:return I.next=2,nn(H);case 2:ie([]),(S=C.current)===null||S===void 0||(A=S.reloadAndRest)===null||A===void 0||A.call(S);case 4:case"end":return I.stop()}},l)}))()},content:(0,e.jsx)(se.ZP,{dataSource:H,rowKey:"id",renderItem:function(l){return(0,e.jsxs)(se.ZP.Item,{children:[l.username,l.fullName?"(".concat(l.fullName,")"):l.email?"(".concat(l.email,")"):""]})}})})},children:m.t("batchDeletion","Batch deletion")}),H.filter(function(b){return b.status!==ne.J0.disabled}).length>0&&(0,e.jsx)(Ze.Z,{onClick:function(){ge.Z.confirm({title:m.t("disableConfirm","Are you sure you want to disable the following users?"),icon:(0,e.jsx)(Be.Z,{}),onOk:function(){return(0,N.Z)((0,f.Z)().mark(function l(){var S,A;return(0,f.Z)().wrap(function(I){for(;;)switch(I.prev=I.next){case 0:return I.next=2,rn(H.filter(function(X){return X.status!==ne.J0.disabled}));case 2:ie([]),(S=C.current)===null||S===void 0||(A=S.reloadAndRest)===null||A===void 0||A.call(S);case 4:case"end":return I.stop()}},l)}))()},content:(0,e.jsx)(se.ZP,{dataSource:H.filter(function(Z){return Z.status!==ne.J0.disabled}),rowKey:"id",renderItem:function(l){return(0,e.jsxs)(se.ZP.Item,{children:[l.username,l.fullName?"(".concat(l.fullName,")"):l.email?"(".concat(l.email,")"):""]})}})})},children:m.t("batchDisable","Batch disable")}),H.filter(function(b){return b.status!==ne.J0.normal}).length>0&&(0,e.jsx)(Ze.Z,{onClick:function(){ge.Z.confirm({title:m.t("enableConfirm","Are you sure you want to enable the following users?"),icon:(0,e.jsx)(Be.Z,{}),onOk:function(){return(0,N.Z)((0,f.Z)().mark(function l(){var S,A;return(0,f.Z)().wrap(function(I){for(;;)switch(I.prev=I.next){case 0:return I.next=2,en(H.filter(function(X){return X.status!==ne.J0.normal}));case 2:ie([]),(S=C.current)===null||S===void 0||(A=S.reloadAndRest)===null||A===void 0||A.call(S);case 4:case"end":return I.stop()}},l)}))()},content:(0,e.jsx)(se.ZP,{dataSource:H.filter(function(Z){return Z.status!==ne.J0.normal}),rowKey:"id",renderItem:function(l){return(0,e.jsxs)(se.ZP.Item,{children:[l.username,l.fullName?"(".concat(l.fullName,")"):l.email?"(".concat(l.email,")"):""]})}})})},children:m.t("batchEnable","Batch enable")})]}),(0,e.jsx)(ue,{visible:he,user:v,onCancel:function(){Ge(!1)},parentIntl:m}),(0,e.jsx)(Jr,{title:m.t(v?"form.title.userUpdate":"form.title.userCreate",v?"Modify user":"Add user"),onSubmit:function(){var b=(0,N.Z)((0,f.Z)().mark(function Z(l){var S;return(0,f.Z)().wrap(function($){for(;;)switch($.prev=$.next){case 0:return $.next=2,(v?Or:_r)(l);case 2:return S=$.sent,S&&(u(!1),O(void 0),C.current&&C.current.reload()),$.abrupt("return",S);case 5:case"end":return $.stop()}},Z)}));return function(Z){return b.apply(this,arguments)}}(),onCancel:function(){u(!1),B||O(void 0)},modalVisible:s,values:v,parentIntl:m}),(0,e.jsx)(Je.Z,{width:800,open:B,onClose:function(){O(void 0),K(!1)},closable:!1,children:B&&(v==null?void 0:v.username)&&(0,e.jsxs)(e.Fragment,{children:[(0,e.jsx)(yr.ZP,{column:2,title:m.t("detail.title","User Details"),request:(0,N.Z)((0,f.Z)().mark(function b(){return(0,f.Z)().wrap(function(l){for(;;)switch(l.prev=l.next){case 0:return l.abrupt("return",{data:v||{}});case 1:case"end":return l.stop()}},b)})),params:{id:v==null?void 0:v.id},extra:(0,e.jsxs)(e.Fragment,{children:[(0,e.jsx)("a",{onClick:function(){D(!0),G("apps")},style:{flex:"unset"},hidden:U,children:m.t("button.grant","Grant")},"grant"),(0,e.jsx)("a",{onClick:(0,N.Z)((0,f.Z)().mark(function b(){var Z,l,S,A,$,I;return(0,f.Z)().wrap(function(De){for(;;)switch(De.prev=De.next){case 0:return Z=v.createTime,l=v.loginTime,S=v.extendedData,A=v.updateTime,$=(0,re.Z)(v,qr),De.next=3,Or((0,w.Z)((0,w.Z)({},$),{},{apps:Q.map(function(Te){return{id:Te.id,roleId:Te.roleId}})}));case 3:I=De.sent,I&&D(!1);case 5:case"end":return De.stop()}},b)})),style:{flex:"unset"},hidden:!U,children:m.t("button.save","Save")},"save")]}),columns:Kr}),(0,e.jsx)(We.Z,{activeKey:k,onChange:function(Z){G(Z)},items:[{label:m.t("detail.sessions.title","Sessions"),key:"sessions",children:(0,e.jsx)(ze.ZP,{actionRef:W,toolBarRender:!1,request:function(Z){return Oe((0,w.Z)((0,w.Z)({},Z),{},{userId:v.id}))},columns:un,rowKey:"id",search:!1})},{label:m.t("detail.keypairs.title","Key Pairs"),key:"keypairs",children:(0,e.jsx)(ze.ZP,{actionRef:x,toolBarRender:!1,request:function(Z){return(0,oe.V6)((0,w.Z)((0,w.Z)({},Z),{},{userId:v.id}))},columns:cn,rowKey:"id",search:!1})},{label:m.t("detail.apps.title","App"),key:"apps",children:(0,e.jsx)(Vr,{apps:Q,onChange:ir,granting:U,parentIntl:m})}]})]})})]})},sn=ln}}]);