(self.webpackChunkidas_ui=self.webpackChunkidas_ui||[]).push([[532],{55356:function(hr,oe,t){"use strict";t.r(oe),t.d(oe,{default:function(){return dr}});var D=t(2824),I=t(39428),u=t(11849),yr=t(34792),B=t(48086),G=t(3182),h=t(67294),Se=t(21010),ae=t(44728),ne=t(93400),xe=t(90949),be=t(75362),se=t(91220),je=t(2844),W=t(14300),de=t(88293),Ce=t(69610),Ie=t(54941),Te=t(56188),ve=function(){function b(){(0,Ce.Z)(this,b),this.generator=void 0,this.generator=new Te.C(0,1,{workerIdBits:5,datacenterIdBits:5,sequenceBits:12})}return(0,Ie.Z)(b,null,[{key:"generate",value:function(){return this.getInstance().generator.nextId().toString()}},{key:"getInstance",value:function(){return this.instance?this.instance:(this.instance=new b,this.instance)}}]),b}();ve.instance=void 0;var re=function(){return ve.generate()},Y=t(62448),Q=t(5966),Pe=t(90672),ce=t(64317),Ve=t(61674),Ee=t(7717),gr=t(57663),Ne=t(71577),me=t(86582),Zr=t(66456),we=t(78587),Sr=t(47673),fe=t(28166),xr=t(43358),Fe=t(66074),br=t(77883),De=t(49747),Re=t(32059),pe=t(93224),jr=t(9715),ie=t(63202),$e=t(80454),le=t(96486),Ae=t.n(le),ue=t(64140),Me=t(42762),Ue=t(82061),Be=t(51042),he=t(63434),Le=t(31649),Ge=t(50317),Ke=t.n(Ge),e=t(85893),He=["title","editable","children","dataIndex","record","handleSave","options","inputType"],Oe=["className","style"],ze=(0,ue.W6)(function(){return(0,e.jsx)(Me.Z,{style:{cursor:"grab",color:"#999"}})}),ye=h.createContext(null),We=(0,ue.W8)(function(b){var a=ie.Z.useForm(),r=(0,D.Z)(a,1),o=r[0];return(0,e.jsx)(ie.Z,{form:o,component:!1,children:(0,e.jsx)(ye.Provider,{value:o,children:(0,e.jsx)("tr",(0,u.Z)({className:"editable-row"},b))})})}),Je=(0,ue.JN)(function(b){return(0,e.jsx)("tbody",(0,u.Z)({},b))}),Qe=["GET","POST","PUT","PATCH","DELETE","OPTIONS","HEAD"].map(function(b){return{value:b,label:b}}),Ye=function b(a){if((0,le.isArray)(a)){var r=(0,se.Z)(a),o;try{for(r.s();!(o=r.n()).done;){var c=o.value;if(!b(c))return!1}}catch(s){r.e(s)}finally{r.f()}return!0}else if((0,le.isString)(a))return!Boolean(a.trim());return!Boolean(a)},Xe=function(a){var r=a.title,o=a.editable,c=a.children,s=a.dataIndex,y=a.record,i=a.handleSave,R=a.options,w=a.inputType,T=(0,pe.Z)(a,He),E=(0,h.useState)(!1),P=(0,D.Z)(E,2),j=P[0],f=P[1],g=(0,h.useRef)(null),$=(0,h.useContext)(ye);(0,h.useEffect)(function(){j&&g.current.focus()},[j]);var U=function(){f(!j),$.setFieldsValue((0,Re.Z)({},s,y[s]))},V=function(){var A=(0,G.Z)((0,I.Z)().mark(function v(){var F;return(0,I.Z)().wrap(function(M){for(;;)switch(M.prev=M.next){case 0:return M.prev=0,M.next=3,$.validateFields();case 3:F=M.sent,U(),i((0,u.Z)((0,u.Z)({},y),F)),M.next=11;break;case 8:M.prev=8,M.t0=M.catch(0),console.error("Save failed:",M.t0);case 11:case"end":return M.stop()}},v,null,[[0,8]])}));return function(){return A.apply(this,arguments)}}(),N=function(){switch(w){case"number":return(0,e.jsx)(De.Z,{ref:g,size:"small",onPressEnter:V,onBlur:V});case"select":return(0,e.jsx)(Fe.Z,{ref:g,size:"small",onBlur:V,options:R});default:return(0,e.jsx)(fe.Z,{ref:g,size:"small",onPressEnter:V,onBlur:V})}}(),L=c;return o&&(L=j?(0,e.jsx)(ie.Z.Item,{className:Ke().ProxyUrlInput,style:{margin:0},name:s,rules:[{required:!0,message:""}],children:N}):(0,e.jsx)("div",{className:"editable-cell-value-wrap",style:{paddingRight:24},onClick:U,children:Ye(c)?(0,e.jsx)(e.Fragment,{children:"\xA0"}):c})),(0,e.jsx)("td",(0,u.Z)((0,u.Z)({},T),{},{children:L}))},ke=function(a){var r=a.parentIntl,o=a.dataSource,c=a.setDataSource,s=(0,h.useState)(240),y=(0,D.Z)(s,2),i=y[0],R=y[1],w=o.domain,T=o.urls,E=o.upstream,P=o.insecureSkipVerify,j=o.transparentServerName,f=(0,h.useState)(0),g=(0,D.Z)(f,2),$=g[0],U=g[1],V=(0,h.useState)(),N=(0,D.Z)(V,2),L=N[0],A=N[1],v=(0,h.useMemo)(function(){return new ne.f("proxy",r)},[r]);(0,h.useEffect)(function(){var m=function(){window.innerHeight>890?R(window.innerHeight-690):R(200)};m(),window.addEventListener("resize",m)},[]),(0,h.useEffect)(function(){var m="",d=(0,se.Z)(T),l;try{for(d.s();!(l=d.n()).done;){var n=l.value;(!n.name||!n.name.trim())&&(m=v.t("name.required","name cannot be empty!")),(!n.method||!n.method.trim())&&(m=v.t("method.required","method cannot be empty!")),(!n.url||!n.url.trim())&&(m=v.t("url.required","URL cannot be empty!"))}}catch(Z){d.e(Z)}finally{d.f()}A(m)},[v,T]);var F=function(d){var l=d.oldIndex,n=d.newIndex;if(l!==n){var Z=(0,$e.q)(T.slice(),l,n).filter(function(S){return!!S});c((0,u.Z)((0,u.Z)({},o),{},{urls:Z}))}},p=function(d){return(0,e.jsx)(Je,(0,u.Z)({useDragHandle:!0,disableAutoscroll:!0,helperClass:"row-dragging",onSortEnd:F},d))},M=function(d){var l=d.className,n=d.style,Z=(0,pe.Z)(d,Oe),S=T.findIndex(function(C){return C.id===Z["data-row-key"]});return(0,e.jsx)(We,(0,u.Z)({index:S},Z))},te=function(d){c(function(l){var n,Z=(n=l==null?void 0:l.urls.map(function(S){return S.id===d.id?d:S}))!==null&&n!==void 0?n:[];return(0,u.Z)((0,u.Z)({},l),{},{urls:Z})})},J=[{dataIndex:"sort",width:30,className:"drag-visible",render:function(){return(0,e.jsx)(ze,{})}},{width:150,title:v.t("name.title","Name"),dataIndex:"name",inputType:"text",editable:!0},{title:v.t("method.title","Method"),dataIndex:"method",inputType:"select",width:100,options:Qe,editable:!0},{title:v.t("url.title","URL"),dataIndex:"url",inputType:"text",editable:!0},{width:30,render:function(d,l){return[(0,e.jsx)("a",{onClick:(0,G.Z)((0,I.Z)().mark(function n(){return(0,I.Z)().wrap(function(S){for(;;)switch(S.prev=S.next){case 0:c(function(C){return(0,u.Z)((0,u.Z)({},C),{},{urls:C.urls.filter(function(O){return O.id!==l.id})})});case 1:case"end":return S.stop()}},n)})),children:(0,e.jsx)(Ue.Z,{})},"delete")]}}],X=J.map(function(m){return m.editable?(0,u.Z)((0,u.Z)({},m),{},{onCell:function(l){return{record:l,editable:m.editable,dataIndex:m.dataIndex,inputType:m.inputType,options:m.options,title:m.title,handleSave:te}}}):m});return(0,e.jsxs)(e.Fragment,{children:[(0,e.jsx)(Q.Z,{name:["proxy","domain"],label:v.t("domain.label","Domain"),colProps:{span:12},initialValue:w,rules:[{required:!0,message:v.t("domain.required","domain cannot be empty!")}]}),(0,e.jsx)(Q.Z,{name:["proxy","upstream"],label:v.t("upstream.label","Upstream"),colProps:{span:12},initialValue:E,tooltip:(0,e.jsxs)(e.Fragment,{children:[v.t("upstream.example","Example"),":",(0,e.jsx)("br",{}),(0,e.jsx)("li",{children:"abc.com"}),(0,e.jsx)("li",{children:"http://abc.com:80"}),(0,e.jsx)("li",{children:"http://1.2.3.4"}),(0,e.jsx)("li",{children:"https://abc.com"})]}),rules:[{required:!0,message:v.t("upstream.required","upstream cannot be empty!")}]}),(0,e.jsx)(he.Z,{name:["proxy","insecureSkipVerify"],label:v.t("insecureSkipVerify.label","Skip TLS Verify"),colProps:{span:12},initialValue:P,tooltip:v.t("insecureSkipVerify.describe","When requesting back-end servers, the certificate verification is ignored (insecure).")}),(0,e.jsx)(he.Z,{name:["proxy","transparentServerName"],label:v.t("transparentServerName.label","Transparent Server Name"),colProps:{span:12},initialValue:j,tooltip:v.t("transparentServerName.describe","When requesting the backend, the domain name requested by the client will be transparently transmitted.")}),(0,e.jsxs)(Le.Z,{colProps:{span:24},label:v.t("url.label","URL"),children:[(0,e.jsx)(we.Z,{pagination:!1,dataSource:T,columns:X,scroll:{y:i},rowKey:"id",size:"small",components:{body:{wrapper:p,row:M,cell:Xe}}}),(0,e.jsx)(Ne.Z,{onClick:function(){c((0,u.Z)((0,u.Z)({},o),{},{urls:[].concat((0,me.Z)(T),[{name:"",id:re(),method:"*",url:""}])})),U($+1)},type:"dashed",block:!0,children:(0,e.jsx)(Be.Z,{})}),(0,e.jsx)("div",{style:{color:"red"},children:L})]})]})},qe=ke,Cr=t(18106),ge=t(58500),Ir=t(49111),_e=t(19650),Tr=t(77576),er=t(12028),Pr=t(13062),rr=t(71230),Vr=t(63185),Ze=t(9676),tr=t(86504),ar=function(a){var r=a.role,o=a.setRole,c=a.activeKey,s=a.setActiveKey,y=(0,h.useRef)(null),i=(0,h.useState)(!1),R=(0,D.Z)(i,2),w=R[0],T=R[1];return(0,h.useEffect)(function(){if(c===r.id&&!r.name&&T(!0),c===r.id&&w){var E;(E=y.current)===null||E===void 0||E.focus({cursor:"all"})}},[w,c,r]),(0,e.jsxs)(e.Fragment,{children:[(0,e.jsx)(fe.Z,{hidden:!(w&&r.id===c),style:{width:80,padding:0},ref:y,onKeyDown:function(P){P.key=="Enter"&&(o({id:r.id,name:P.currentTarget.value}),T(!1))},onBlur:function(P){o({id:r.id,name:P.target.value}),T(!1)},autoFocus:!0,defaultValue:r.name?r.name:"New Role"}),(0,e.jsxs)("div",{hidden:w&&r.id===c,style:{width:80,overflow:"hidden",textShadow:"0 0 0 aliceblue"},title:r.name,onDoubleClick:function(){s(r.id),T(!0)},children:[(0,e.jsx)(tr.Z,{style:{color:"orange"},hidden:!r.isDefault}),r.name]})]})},nr=function(a){var r=a.urls,o=r===void 0?[]:r,c=a.value,s=c===void 0?[]:c,y=a.onChange,i=(0,h.useState)(),R=(0,D.Z)(i,2),w=R[0],T=R[1],E={add:function(){var f="new-role-".concat(re());T(f),y==null||y([].concat((0,me.Z)(s),[{name:"",id:f}]))},remove:function(f){Ae().isString(f)?y==null||y(s.filter(function(g){return g.id!=f})):B.default.warning("system error: ".concat(f," is not string"))}},P=function(){var j=(0,G.Z)((0,I.Z)().mark(function f(g){return(0,I.Z)().wrap(function(U){for(;;)switch(U.prev=U.next){case 0:return U.next=2,y==null?void 0:y(s.map(function(V){return g.id==V.id?(0,u.Z)((0,u.Z)({},V),g):V}));case 2:case"end":return U.stop()}},f)}));return function(g){return j.apply(this,arguments)}}();return(0,e.jsx)(ge.Z,{type:"editable-card",tabPosition:"top",style:{height:350},activeKey:w,onChange:T,onEdit:function(f,g){E[g](f)},children:s.map(function(j){return(0,e.jsx)(ge.Z.TabPane,{tab:(0,e.jsx)(ar,{role:j,setRole:P,activeKey:w,setActiveKey:T}),tabKey:j.id,children:(0,e.jsxs)(_e.Z,{direction:"vertical",children:[(0,e.jsxs)("div",{style:{display:o.length===0?"none":"unset"},children:[(0,e.jsx)("div",{style:{paddingBottom:7},children:"Permission"}),(0,e.jsx)(Ze.Z.Group,{onChange:function(g){P((0,u.Z)((0,u.Z)({},j),{},{urls:g}))},defaultValue:j.urls,children:o.map(function(f){return(0,e.jsx)(rr.Z,{children:(0,e.jsx)(Ze.Z,{value:f.id,children:f.name})},f.id)})})]}),(0,e.jsxs)("div",{style:{padding:10},children:["\u9ED8\u8BA4:",(0,e.jsx)(er.Z,{checked:j.isDefault,onChange:function(g){y==null||y(s.map(function($){return j.id==$.id?(0,u.Z)((0,u.Z)({},$),{},{isDefault:g}):g?(0,u.Z)((0,u.Z)({},$),{},{isDefault:!g}):$}))}})]})]})},j.id)})})},sr=function(a){var r=a.values,o=a.onSubmit,c=a.parentIntl,s=a.loading,y=a.disabled,i=new ne.f("form",c),R=(0,h.useState)([]),w=(0,D.Z)(R,2),T=w[0],E=w[1],P=(0,h.useRef)(),j=(0,h.useState)([]),f=(0,D.Z)(j,2),g=f[0],$=f[1],U=(0,h.useState)([]),V=(0,D.Z)(U,2),N=V[0],L=V[1],A=(0,h.useState)({domain:"",upstream:"",urls:[{id:re(),method:"*",name:"default",url:"/"}],insecureSkipVerify:!1,transparentServerName:!0}),v=(0,D.Z)(A,2),F=v[0],p=v[1],M=(0,h.useState)(0),te=(0,D.Z)(M,2),J=te[0],X=te[1];(0,h.useEffect)(function(){var d,l,n,Z,S;X(0),L((d=r==null||(l=r.roles)===null||l===void 0?void 0:l.map(function(C){return{id:C.id,urls:C.urls,name:C.name,isDefault:C.isDefault}}))!==null&&d!==void 0?d:[]),$((n=r==null?void 0:r.users)!==null&&n!==void 0?n:[]),p((0,u.Z)({domain:"",upstream:"",transparentServerName:!0,insecureSkipVerify:!1,urls:[{id:re(),method:"*",name:"default",url:"/"}]},(Z=r==null?void 0:r.proxy)!==null&&Z!==void 0?Z:{})),E((S=r==null?void 0:r.grantType)!==null&&S!==void 0?S:[])},[r]);var m=function(l){return(l??[]).includes(W.qJ.proxy)};return(0,e.jsxs)(Y.L,{stepsProps:{size:"small"},formProps:{preserve:!1,disabled:y,loading:s},onCurrentChange:function(){var d=(0,G.Z)((0,I.Z)().mark(function l(n){var Z;return(0,I.Z)().wrap(function(C){for(;;)switch(C.prev=C.next){case 0:n==1&&!m((Z=P.current)===null||Z===void 0?void 0:Z.getFieldValue("grantType"))?X(n+(J===0?1:-1)):X(n);case 1:case"end":return C.stop()}},l)}));return function(l){return d.apply(this,arguments)}}(),current:J,onFinish:function(){var d=(0,G.Z)((0,I.Z)().mark(function l(n){var Z,S,C,O,k,q,_,z,ee,K;return(0,I.Z)().wrap(function(x){for(;;)switch(x.prev=x.next){case 0:return console.log(n),Z=W.FQ.manual,S=n.grantMode,C=n.grantType,O=n.status,k=F.domain,q=F.upstream,_=F.urls,z=F.transparentServerName,ee=F.insecureSkipVerify,K=_.map(function(H){var cr=H.id,mr=H.method,fr=H.name,pr=H.url;return{id:cr,method:mr,name:fr,url:pr}}),x.abrupt("return",o((0,u.Z)((0,u.Z)({},n),{},{users:g.map(function(H){return{id:H.id,roleId:H.roleId}}),status:O!==W.w$.unknown?O:W.w$.normal,roles:N.map(function(H){return m(C)?H:(0,u.Z)((0,u.Z)({},H),{},{urls:void 0})}),grantMode:S??Z,grantType:C??[W.qJ.none],proxy:m(C)?{domain:k,upstream:q,urls:K,insecureSkipVerify:ee,transparentServerName:z}:void 0})));case 6:case"end":return x.stop()}},l)}));return function(l){return d.apply(this,arguments)}}(),children:[(0,e.jsxs)(Y.L.StepForm,{initialValues:r,labelCol:{span:8},formRef:P,wrapperCol:{span:14},layout:"vertical",title:i.t("basicConfig.title","Basic"),onFinish:(0,G.Z)((0,I.Z)().mark(function d(){return(0,I.Z)().wrap(function(n){for(;;)switch(n.prev=n.next){case 0:if(!(r&&!r.id)){n.next=3;break}return B.default.error(i.t("app-id.empty","System error, application ID is empty")),n.abrupt("return",!1);case 3:return n.abrupt("return",!0);case 4:case"end":return n.stop()}},d)})),children:[(0,e.jsx)(je.HZ,{label:i.t("avatar.label","Avatar"),name:"avatar"}),(0,e.jsx)(Q.Z,{hidden:!0,name:"id"}),(0,e.jsx)(Q.Z,{name:"name",label:i.t("name.label","Name"),width:"md",rules:[{required:!0,message:i.t("name.required","Please input app name!")},{pattern:/^[-_A-Za-z0-9]+$/,message:i.t("name.invalid","App name format error!")}],disabled:!!(r!=null&&r.name)}),(0,e.jsx)(Q.Z,{name:"displayName",label:i.t("displayName.label","Display Name"),width:"md"}),(0,e.jsx)(Pe.Z,{name:"description",label:i.t("description.label","Description"),width:"md"}),(0,e.jsx)(Q.Z,{name:"url",label:i.t("url.label","URL"),width:"md",rules:[{required:!0,message:i.t("url.required","Please input url!")}]}),(0,e.jsx)(ce.Z,{name:"grantType",label:i.t("grantType.label","Grant Type"),width:"md",mode:"tags",options:(0,de.MM)(W.qJ,c,"grantType.value"),rules:[{required:!0,message:i.t("grantType.required","Please select Grant Type!")}],fieldProps:{onChange:E}}),(0,e.jsx)(ce.Z,{name:"grantMode",label:i.t("grantMode.label","Grant Mode"),width:"md",options:(0,de.MM)(W.FQ,c,"grantMode.value"),rules:[{required:!0,message:i.t("grantMode.required","Please select Grant Mode!")}]})]}),(0,e.jsx)(Y.L.StepForm,{initialValues:{},layout:"vertical",title:i.t("proxy.title","Proxy"),grid:!0,onFinish:function(){var d=(0,G.Z)((0,I.Z)().mark(function l(n){var Z,S,C,O,k,q,_,z,ee,K;return(0,I.Z)().wrap(function(x){for(;;)switch(x.prev=x.next){case 0:S=(Z=n.proxy)!==null&&Z!==void 0?Z:{},C=S.domain,O=S.upstream,k=S.transparentServerName,q=S.insecureSkipVerify,p((0,u.Z)((0,u.Z)({},F),{},{domain:C,upstream:O,transparentServerName:k,insecureSkipVerify:q})),_=F.urls,z=(0,se.Z)(_),x.prev=4,z.s();case 6:if((ee=z.n()).done){x.next=19;break}if(K=ee.value,!(!K.name||!K.name.trim())){x.next=11;break}return B.default.error(i.t("name.required","name cannot be empty!")),x.abrupt("return",!1);case 11:if(!(!K.method||!K.method.trim())){x.next=14;break}return B.default.error(i.t("proxy.method.required","method cannot be empty!")),x.abrupt("return",!1);case 14:if(!(!K.url||!K.url.trim())){x.next=17;break}return B.default.error(i.t("proxy.url.required","URL cannot be empty!")),x.abrupt("return",!1);case 17:x.next=6;break;case 19:x.next=24;break;case 21:x.prev=21,x.t0=x.catch(4),z.e(x.t0);case 24:return x.prev=24,z.f(),x.finish(24);case 27:return x.abrupt("return",!0);case 28:case"end":return x.stop()}},l,null,[[4,21,24,27]])}));return function(l){return d.apply(this,arguments)}}(),children:J===1&&(0,e.jsx)(qe,{dataSource:F,setDataSource:p,parentIntl:i})}),(0,e.jsx)(Y.L.StepForm,{initialValues:{},layout:"vertical",title:i.t("role.title","Role"),children:J===2&&(0,e.jsx)(Ve.Z,{name:"rules",rules:[{validator:function(){if(N.length>0){if(N.filter(function(n){return n.isDefault}).length!=1)return Promise.reject(new Error("No default role specified!"));var l=N.map(function(n){return n.name});if(l.some(function(n,Z){return l.includes(n,Z+1)}))return Promise.reject(new Error("Contains duplicate role names!"));if(l.some(function(n){return!Boolean(n)}))return Promise.reject(new Error("Contains empty role names!"))}return Promise.resolve()}}],children:(0,e.jsx)(nr,{value:N,urls:m(T)?F.urls:[],onChange:L})})}),(0,e.jsx)(Y.L.StepForm,{initialValues:{},layout:"vertical",title:i.t("user.title","User"),children:J===3&&(0,e.jsx)("div",{style:{height:"calc( 100vh - 400px )"},children:(0,e.jsx)(Ee.ZP,{users:g,roles:N,onChange:$,granting:!0,parentIntl:c})})})]})},ir=sr,lr=function(){var b=(0,G.Z)((0,I.Z)().mark(function a(r){var o;return(0,I.Z)().wrap(function(s){for(;;)switch(s.prev=s.next){case 0:return o=B.default.loading("Adding ..."),s.prev=1,s.next=4,(0,ae.ri)((0,u.Z)({},r));case 4:return o(),B.default.success("Added successfully"),s.abrupt("return",!0);case 9:return s.prev=9,s.t0=s.catch(1),o(),B.default.error("Adding failed, please try again!"),s.abrupt("return",!1);case 14:case"end":return s.stop()}},a,null,[[1,9]])}));return function(r){return b.apply(this,arguments)}}(),ur=function(){var b=(0,G.Z)((0,I.Z)().mark(function a(r){var o;return(0,I.Z)().wrap(function(s){for(;;)switch(s.prev=s.next){case 0:if(o=B.default.loading("Configuring"),s.prev=1,!r.id){s.next=10;break}return s.next=5,(0,ae.KT)({id:r.id},{id:r.id,name:r.name,status:r.status,description:r.description,grantType:r.grantType,grantMode:r.grantMode,avatar:r.avatar,roles:r.roles,users:r.users,proxy:r.proxy,url:r.url,displayName:r.displayName});case 5:return o(),B.default.success("update is successful"),s.abrupt("return",!0);case 10:return B.default.success("update failed, system error"),s.abrupt("return",!1);case 12:s.next=19;break;case 14:return s.prev=14,s.t0=s.catch(1),o(),B.default.error("Configuration failed, please try again!"),s.abrupt("return",!1);case 19:case"end":return s.stop()}},a,null,[[1,14]])}));return function(r){return b.apply(this,arguments)}}(),or=function(a){var r=a.match.params.aid,o=a.history,c=new ne.f("pages.apps",(0,Se.YB)()),s=(0,h.useState)(),y=(0,D.Z)(s,2),i=y[0],R=y[1],w=(0,h.useState)(!1),T=(0,D.Z)(w,2),E=T[0],P=T[1],j=(0,h.useState)(!1),f=(0,D.Z)(j,2),g=f[0],$=f[1],U=function(){var V=(0,G.Z)((0,I.Z)().mark(function N(L){var A,v;return(0,I.Z)().wrap(function(p){for(;;)switch(p.prev=p.next){case 0:return p.prev=0,P(!0),p.next=4,(0,ae.BN)({id:L});case 4:return A=p.sent,A.data&&(v=A.data,R(v)),P(!1),p.abrupt("return",A);case 10:return p.prev=10,p.t0=p.catch(0),$(!0),p.t0.handled||console.error("failed to get app info: ".concat(p.t0)),p.abrupt("return",{success:!1});case 15:case"end":return p.stop()}},N,null,[[0,10]])}));return function(L){return V.apply(this,arguments)}}();return(0,h.useEffect)(function(){r?U(r):R({})},[r]),(0,e.jsx)(be.ZP,{title:r?i==null?void 0:i.name:!1,children:(0,e.jsx)(xe.ZP,{children:i?(0,e.jsx)(ir,{onSubmit:function(){var V=(0,G.Z)((0,I.Z)().mark(function N(L){var A,v;return(0,I.Z)().wrap(function(p){for(;;)switch(p.prev=p.next){case 0:return p.next=2,(i!=null&&i.id?ur:lr)(L);case 2:return A=p.sent,A&&o.push("/apps/".concat((v=i==null?void 0:i.id)!==null&&v!==void 0?v:"")),p.abrupt("return",A);case 5:case"end":return p.stop()}},N)}));return function(N){return V.apply(this,arguments)}}(),values:i!=null&&i.id?i:void 0,parentIntl:c,loading:E,disabled:g}):null})})},dr=or}}]);