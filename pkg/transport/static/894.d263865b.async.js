(self.webpackChunkidas_ui=self.webpackChunkidas_ui||[]).push([[894],{42489:function(ye,J,r){"use strict";r.d(J,{I:function(){return cr}});var f=r(28481),Z=r(96156),n=r(28991),b=r(55507),M=r(92137),O=r(81253),H=r(17503),t=r(67294),X=r(7902),W=0;function de(e){var o=(0,t.useState)(function(){return e.proFieldKey?e.proFieldKey.toString():(W+=1,W.toString())}),T=(0,f.Z)(o,1),c=T[0],A=(0,t.useRef)(c),I=function(){var y=(0,M.Z)((0,b.Z)().mark(function w(){var l,P;return(0,b.Z)().wrap(function(v){for(;;)switch(v.prev=v.next){case 0:return v.next=2,(l=e.request)===null||l===void 0?void 0:l.call(e,e.params,e);case 2:return P=v.sent,v.abrupt("return",P);case 4:case"end":return v.stop()}},w)}));return function(){return y.apply(this,arguments)}}();(0,t.useEffect)(function(){return function(){W+=1}},[]);var N=(0,X.ZP)([A.current,e.params],I,{revalidateOnFocus:!1,shouldRetryOnError:!1,revalidateOnReconnect:!1}),L=N.data,i=N.error;return[L||i]}var Ze=r(85061),C=r(90484),u=r(88306),a=r(8880),K=r(74763),Qe=r(92210);function ke(e){return(0,C.Z)(e)!=="object"?!1:e===null?!0:!(t.isValidElement(e)||e.constructor===RegExp||e instanceof Map||e instanceof Set||e instanceof HTMLElement||e instanceof Blob||e instanceof File||Array.isArray(e))}var Ue=function(o,T){var c=arguments.length>2&&arguments[2]!==void 0?arguments[2]:!0,A=Object.keys(T).reduce(function(L,i){var y=T[i];return(0,K.k)(y)||(L[i]=y),L},{});if(Object.keys(A).length<1||typeof window=="undefined"||(0,C.Z)(o)!=="object"||(0,K.k)(o)||o instanceof Blob)return o;var I=Array.isArray(o)?[]:{},N=function L(i,y){var w=Array.isArray(i),l=w?[]:{};return i==null||i===void 0?l:(Object.keys(i).forEach(function(P){var ie=y?[y,P].flat(1):[P].flat(1),v=i[P],V=(0,u.Z)(A,ie),z=function E(s){return Array.isArray(s)&&s.forEach(function(F,x){!F||(typeof F=="function"&&(v[x]=F(v,P,i)),(0,C.Z)(F)==="object"&&!Array.isArray(F)&&Object.keys(F).forEach(function(k){if(typeof F[k]=="function"){var ue=F[k](i[P][x][k],P,i);v[x][k]=(0,C.Z)(ue)==="object"?ue[k]:ue}}),(0,C.Z)(F)==="object"&&Array.isArray(F)&&E(F))}),P},Q=function(){var s=typeof V=="function"?V==null?void 0:V(v,P,i):z(V);if(Array.isArray(s)){l=(0,a.Z)(l,s,v);return}(0,C.Z)(s)==="object"&&!Array.isArray(I)?I=(0,n.Z)((0,n.Z)({},I),s):(0,C.Z)(s)==="object"&&Array.isArray(I)?l=(0,n.Z)((0,n.Z)({},l),s):s&&(l=(0,a.Z)(l,[s],v))};if(V&&typeof V=="function"&&Q(),typeof window!="undefined"){if(ke(v)){var U=L(v,ie);if(Object.keys(U).length<1)return;l=(0,a.Z)(l,[P],U);return}Q()}}),c?l:i)};return I=Array.isArray(o)&&Array.isArray(I)?(0,Ze.Z)(N(o)):(0,Qe.T)({},N(o),I),I},qe=r(22270),Ae=r(48171),me=r(26369),Be=r(60249),Ke=r(41036),_e=r(21770),ce=r(75661),Le=r(64385),Ce=r(23312),er=r(45095),fe=r(63202),ne=r(88182),$e=r(11382),Ie=r(94184),We=r.n(Ie),te=r(97435),De=r(80334),be=r(71577),Me=r(19650),p=r(85893),oe=function(o){var T=(0,H.YB)(),c=fe.Z.useFormInstance();if(o.render===!1)return null;var A=o.onSubmit,I=o.render,N=o.onReset,L=o.searchConfig,i=L===void 0?{}:L,y=o.submitButtonProps,w=o.resetButtonProps,l=w===void 0?{}:w,P=function(){c.submit(),A==null||A()},ie=function(){c.resetFields(),N==null||N()},v=i.submitText,V=v===void 0?T.getMessage("tableForm.submit","\u63D0\u4EA4"):v,z=i.resetText,Q=z===void 0?T.getMessage("tableForm.reset","\u91CD\u7F6E"):z,U=[];l!==!1&&U.push((0,t.createElement)(be.Z,(0,n.Z)((0,n.Z)({},(0,te.Z)(l,["preventDefault"])),{},{key:"rest",onClick:function(F){var x;l!=null&&l.preventDefault||ie(),l==null||(x=l.onClick)===null||x===void 0||x.call(l,F)}}),Q)),y!==!1&&U.push((0,t.createElement)(be.Z,(0,n.Z)((0,n.Z)({type:"primary"},(0,te.Z)(y||{},["preventDefault"])),{},{key:"submit",onClick:function(F){var x;y!=null&&y.preventDefault||P(),y==null||(x=y.onClick)===null||x===void 0||x.call(y,F)}}),V));var E=I?I((0,n.Z)((0,n.Z)({},o),{},{form:c,submit:P,reset:ie}),U):U;return E?Array.isArray(E)?(E==null?void 0:E.length)<1?null:(E==null?void 0:E.length)===1?E[0]:(0,p.jsx)(Me.Z,{wrap:!0,children:E}):E:null},le=oe,ae=r(10279),G=r(66758),rr=r(2514),dr=r(9105),fr=["children","contentRender","submitter","fieldProps","formItemProps","groupProps","transformKey","formRef","onInit","form","loading","formComponentType","extraUrlParams","syncToUrl","onUrlSearchChange","onReset","omitNil","isKeyPressSubmit","autoFocusFirstInput","grid","rowProps","colProps"],vr=["extraUrlParams","syncToUrl","isKeyPressSubmit","syncToUrlAsImportant","syncToInitialValues","children","contentRender","submitter","fieldProps","proFieldProps","formItemProps","groupProps","dateFormatter","formRef","onInit","form","formComponentType","onReset","grid","rowProps","colProps","omitNil","request","params","initialValues","formKey","readonly","onLoadingChange","loading"],Ne=function(o,T,c){return o===!0?T:(0,qe.h)(o,T,c)},nr=function(o){return!o||Array.isArray(o)?o:[o]};function mr(e){var o,T=e.children,c=e.contentRender,A=e.submitter,I=e.fieldProps,N=e.formItemProps,L=e.groupProps,i=e.transformKey,y=e.formRef,w=e.onInit,l=e.form,P=e.loading,ie=e.formComponentType,v=e.extraUrlParams,V=v===void 0?{}:v,z=e.syncToUrl,Q=e.onUrlSearchChange,U=e.onReset,E=e.omitNil,s=E===void 0?!0:E,F=e.isKeyPressSubmit,x=e.autoFocusFirstInput,k=x===void 0?!0:x,ue=e.grid,or=e.rowProps,Oe=e.colProps,Pe=(0,O.Z)(e,fr),xe=fe.Z.useFormInstance(),we=(ne.ZP===null||ne.ZP===void 0||(o=ne.ZP.useConfig)===null||o===void 0?void 0:o.call(ne.ZP))||{componentSize:"middle"},Ee=we.componentSize,q=(0,t.useRef)(l||xe),ze=(0,rr.zx)({grid:ue,rowProps:or}),Je=ze.RowWrapper,Fe=(0,Ae.J)(function(){return xe}),Se=(0,t.useMemo)(function(){return{getFieldsFormatValue:function(g){var h;return i((h=Fe())===null||h===void 0?void 0:h.getFieldsValue(g),s)},getFieldFormatValue:function(){var g,h=arguments.length>0&&arguments[0]!==void 0?arguments[0]:[],m=nr(h);if(!m)throw new Error("nameList is require");var D=(g=Fe())===null||g===void 0?void 0:g.getFieldValue(m),_=m?(0,a.Z)({},m,D):D;return(0,u.Z)(i(_,s,m),m)},getFieldFormatValueObject:function(g){var h,m=nr(g),D=(h=Fe())===null||h===void 0?void 0:h.getFieldValue(m),_=m?(0,a.Z)({},m,D):D;return i(_,s,m)},validateFieldsReturnFormatValue:function(){var d=(0,M.Z)((0,b.Z)().mark(function h(m){var D,_,ve;return(0,b.Z)().wrap(function($){for(;;)switch($.prev=$.next){case 0:if(!(!Array.isArray(m)&&m)){$.next=2;break}throw new Error("nameList must be array");case 2:return $.next=4,(D=Fe())===null||D===void 0?void 0:D.validateFields(m);case 4:return _=$.sent,ve=i(_,s),$.abrupt("return",ve||{});case 7:case"end":return $.stop()}},h)}));function g(h){return d.apply(this,arguments)}return g}(),formRef:q}},[s,i]),Y=(0,t.useMemo)(function(){return t.Children.toArray(T).map(function(d,g){return g===0&&t.isValidElement(d)&&k?t.cloneElement(d,(0,n.Z)((0,n.Z)({},d.props),{},{autoFocus:k})):d})},[k,T]),R=(0,t.useMemo)(function(){return typeof A=="boolean"||!A?{}:A},[A]),je=(0,t.useMemo)(function(){if(A!==!1)return(0,p.jsx)(le,(0,n.Z)((0,n.Z)({},R),{},{onReset:function(){var g,h,m=i((g=q.current)===null||g===void 0?void 0:g.getFieldsValue(),s);if(R==null||(h=R.onReset)===null||h===void 0||h.call(R,m),U==null||U(m),z){var D,_=Object.keys(i((D=q.current)===null||D===void 0?void 0:D.getFieldsValue(),!1)).reduce(function(ve,pe){return(0,n.Z)((0,n.Z)({},ve),{},(0,Z.Z)({},pe,m[pe]||void 0))},V);Q(Ne(z,_,"set"))}},submitButtonProps:(0,n.Z)({loading:P},R.submitButtonProps)}),"submitter")},[A,R,P,i,s,U,z,V,Q]),Ve=(0,t.useMemo)(function(){var d=ue?(0,p.jsx)(Je,{children:Y}):Y;return c?c(d,je,q.current):d},[ue,Je,Y,c,je]),Re=(0,me.D)(e.initialValues);return(0,t.useEffect)(function(){if(!(z||!e.initialValues||!Re||Pe.request)){var d=(0,Be.A)(e.initialValues,Re);(0,De.ET)(d,"initialValues \u53EA\u5728 form \u521D\u59CB\u5316\u65F6\u751F\u6548\uFF0C\u5982\u679C\u4F60\u9700\u8981\u5F02\u6B65\u52A0\u8F7D\u63A8\u8350\u4F7F\u7528 request\uFF0C\u6216\u8005 initialValues ? <Form/> : null "),(0,De.ET)(d,"The initialValues only take effect when the form is initialized, if you need to load asynchronously recommended request, or the initialValues ? <Form/> : null ")}},[e.initialValues]),(0,t.useImperativeHandle)(y,function(){return(0,n.Z)((0,n.Z)({},q.current),Se)},[]),(0,t.useEffect)(function(){var d,g,h=i((d=q.current)===null||d===void 0||(g=d.getFieldsValue)===null||g===void 0?void 0:g.call(d,!0),s);w==null||w(h,q.current)},[]),(0,p.jsx)(Ke.J.Provider,{value:Se,children:(0,p.jsx)(ne.ZP,{componentSize:Pe.size||Ee,children:(0,p.jsxs)(rr._p.Provider,{value:{grid:ue,colProps:Oe},children:[Pe.component!==!1&&(0,p.jsx)("input",{type:"text",style:{display:"none"}}),Ve]})})})}var tr=0;function cr(e){var o=e.extraUrlParams,T=o===void 0?{}:o,c=e.syncToUrl,A=e.isKeyPressSubmit,I=e.syncToUrlAsImportant,N=I===void 0?!1:I,L=e.syncToInitialValues,i=L===void 0?!0:L,y=e.children,w=e.contentRender,l=e.submitter,P=e.fieldProps,ie=e.proFieldProps,v=e.formItemProps,V=e.groupProps,z=e.dateFormatter,Q=z===void 0?"string":z,U=e.formRef,E=e.onInit,s=e.form,F=e.formComponentType,x=e.onReset,k=e.grid,ue=e.rowProps,or=e.colProps,Oe=e.omitNil,Pe=Oe===void 0?!0:Oe,xe=e.request,we=e.params,Ee=e.initialValues,q=e.formKey,ze=q===void 0?tr:q,Je=e.readonly,Fe=e.onLoadingChange,Se=e.loading,Y=(0,O.Z)(e,vr),R=(0,t.useRef)({}),je=(0,_e.Z)(!1,{onChange:Fe,value:Se}),Ve=(0,f.Z)(je,2),Re=Ve[0],d=Ve[1],g=(0,er.l)({},{disabled:!c}),h=(0,f.Z)(g,2),m=h[0],D=h[1],_=(0,t.useRef)((0,ce.x)());(0,t.useEffect)(function(){tr+=0},[]);var ve=de({request:xe,params:we,proFieldKey:ze}),pe=(0,f.Z)(ve,1),$=pe[0],Pr=(0,t.useContext)(ne.ZP.ConfigContext),Fr=Pr.getPrefixCls,ar=Fr("pro-form"),ir=(0,Le.Xj)("ProForm",function(ee){return(0,Z.Z)({},".".concat(ar),(0,Z.Z)({},"> div:not(".concat(ee.proComponentsCls,"-form-light-filter)"),{".pro-field":{maxWidth:"100%","&-xs":{width:104},"&-s":{width:216},"&-sm":{width:216},"&-m":{width:328},"&-md":{width:328},"&-l":{width:440},"&-lg":{width:440},"&-xl":{width:552}}}))}),gr=ir.wrapSSR,hr=ir.hashId,yr=(0,t.useState)(function(){return c?Ne(c,m,"get"):{}}),ur=(0,f.Z)(yr,2),lr=ur[0],Zr=ur[1],He=(0,t.useRef)({}),Ge=(0,t.useRef)({}),Ye=(0,t.useCallback)(function(ee,S,j){return Ue((0,Ce.lp)(ee,Q,Ge.current,S,j),He.current,S)},[Q]);(0,t.useEffect)(function(){i||Zr({})},[i]),(0,t.useEffect)(function(){!c||D((0,n.Z)((0,n.Z)({},m),T))},[T,c]);var Cr=(0,t.useMemo)(function(){if(typeof window!="undefined"&&F&&["DrawerForm"].includes(F))return function(ee){return ee.parentNode||document.body}},[F]),Er=(0,Ae.J)((0,M.Z)((0,b.Z)().mark(function ee(){var S,j,re,ge,Te,se;return(0,b.Z)().wrap(function(B){for(;;)switch(B.prev=B.next){case 0:if(Y.onFinish){B.next=2;break}return B.abrupt("return");case 2:if(!Re){B.next=4;break}return B.abrupt("return");case 4:return d(!0),B.prev=5,re=R==null||(S=R.current)===null||S===void 0||(j=S.getFieldsFormatValue)===null||j===void 0?void 0:j.call(S),B.next=9,Y.onFinish(re);case 9:c&&(se=Object.keys(R==null||(ge=R.current)===null||ge===void 0||(Te=ge.getFieldsFormatValue)===null||Te===void 0?void 0:Te.call(ge,void 0,!1)).reduce(function(he,sr){var Xe;return(0,n.Z)((0,n.Z)({},he),{},(0,Z.Z)({},sr,(Xe=re[sr])!==null&&Xe!==void 0?Xe:void 0))},T),Object.keys(m).forEach(function(he){se[he]!==!1&&se[he]!==0&&!se[he]&&(se[he]=void 0)}),D(Ne(c,se,"set"))),d(!1),B.next=17;break;case 13:B.prev=13,B.t0=B.catch(5),console.log(B.t0),d(!1);case 17:case"end":return B.stop()}},ee,null,[[5,13]])})));return(0,t.useImperativeHandle)(U,function(){return R.current},[!$]),!$&&e.request?(0,p.jsx)("div",{style:{paddingTop:50,paddingBottom:50,textAlign:"center"},children:(0,p.jsx)($e.Z,{})}):gr((0,p.jsx)(dr.A.Provider,{value:{mode:e.readonly?"read":"edit"},children:(0,p.jsx)(H._Y,{needDeps:!0,children:(0,p.jsx)(G.Z.Provider,{value:{formRef:R,fieldProps:P,proFieldProps:ie,formItemProps:v,groupProps:V,formComponentType:F,getPopupContainer:Cr,formKey:_.current,setFieldValueType:function(S,j){var re=j.valueType,ge=re===void 0?"text":re,Te=j.dateFormat,se=j.transform;!Array.isArray(S)||(He.current=(0,a.Z)(He.current,S,se),Ge.current=(0,a.Z)(Ge.current,S,{valueType:ge,dateFormat:Te}))}},children:(0,p.jsx)(ae.J.Provider,{value:{},children:(0,p.jsx)(fe.Z,(0,n.Z)((0,n.Z)({onKeyPress:function(S){if(!!A&&S.key==="Enter"){var j;(j=R.current)===null||j===void 0||j.submit()}},autoComplete:"off",form:s},(0,te.Z)(Y,["labelWidth","autoFocusFirstInput"])),{},{initialValues:N?(0,n.Z)((0,n.Z)((0,n.Z)({},Ee),$),lr):(0,n.Z)((0,n.Z)((0,n.Z)({},lr),Ee),$),onValuesChange:function(S,j){var re;Y==null||(re=Y.onValuesChange)===null||re===void 0||re.call(Y,Ye(S,!!Pe),Ye(j,!!Pe))},className:We()(e.className,ar,hr),onFinish:Er,children:(0,p.jsx)(mr,(0,n.Z)((0,n.Z)({transformKey:Ye,autoComplete:"off",loading:Re,onUrlSearchChange:D},e),{},{formRef:R,initialValues:(0,n.Z)((0,n.Z)({},Ee),$)}))}))})})})}))}},9105:function(ye,J,r){"use strict";r.d(J,{A:function(){return Z}});var f=r(67294),Z=f.createContext({mode:"edit"})},31649:function(ye,J,r){"use strict";var f=r(28991),Z=r(81253),n=r(38127),b=r(22270),M=r(60249),O=r(67294),H=r(93189),t=r(9105),X=r(85893),W=["fieldProps","children","labelCol","label","autoFocus","isDefaultDom","render","proFieldProps","renderFormItem","valueType","initialValue","onChange","valueEnum","params","name","dependenciesValues","cacheForSwr","valuePropName"],de=function(u){var a=u.fieldProps,K=u.children,Qe=u.labelCol,ke=u.label,Ue=u.autoFocus,qe=u.isDefaultDom,Ae=u.render,me=u.proFieldProps,Be=u.renderFormItem,Ke=u.valueType,_e=u.initialValue,ce=u.onChange,Le=u.valueEnum,Ce=u.params,er=u.name,fe=u.dependenciesValues,ne=u.cacheForSwr,$e=ne===void 0?!1:ne,Ie=u.valuePropName,We=Ie===void 0?"value":Ie,te=(0,Z.Z)(u,W),De=(0,O.useContext)(t.A),be=(0,O.useMemo)(function(){return fe&&te.request?(0,f.Z)((0,f.Z)({},Ce),fe||{}):Ce},[fe,Ce,te.request]),Me=(0,O.useMemo)(function(){if(K)return O.isValidElement(K)?O.cloneElement(K,(0,f.Z)((0,f.Z)({},te),{},{onChange:function(){for(var oe=arguments.length,le=new Array(oe),ae=0;ae<oe;ae++)le[ae]=arguments[ae];if(a!=null&&a.onChange){var G;a==null||(G=a.onChange)===null||G===void 0||G.call.apply(G,[a].concat(le));return}ce==null||ce.apply(void 0,le)}},K.props)):(0,X.jsx)(X.Fragment,{children:K})},[K,a==null?void 0:a.onChange,ce,te]);return Me||(0,X.jsx)(n.ZP,(0,f.Z)((0,f.Z)((0,f.Z)({text:a==null?void 0:a[We],render:Ae,renderFormItem:Be,valueType:Ke||"text",cacheForSwr:$e,fieldProps:(0,f.Z)((0,f.Z)({autoFocus:Ue},a),{},{onChange:function(){if(a!=null&&a.onChange){for(var oe,le=arguments.length,ae=new Array(le),G=0;G<le;G++)ae[G]=arguments[G];a==null||(oe=a.onChange)===null||oe===void 0||oe.call.apply(oe,[a].concat(ae));return}}}),valueEnum:(0,b.h)(Le)},me),te),{},{mode:(me==null?void 0:me.mode)||De.mode||"edit",params:be}))},Ze=(0,H.G)((0,O.memo)(de,function(C,u){return(0,M.A)(u,C,["onChange","onBlur"])}));J.Z=Ze},5966:function(ye,J,r){"use strict";var f=r(28991),Z=r(81253),n=r(67294),b=r(31649),M=r(85893),O=["fieldProps","proFieldProps"],H=["fieldProps","proFieldProps"],t="text",X=function(C){var u=C.fieldProps,a=C.proFieldProps,K=(0,Z.Z)(C,O);return(0,M.jsx)(b.Z,(0,f.Z)({valueType:t,fieldProps:u,filedConfig:{valueType:t},proFieldProps:a},K))},W=function(C){var u=C.fieldProps,a=C.proFieldProps,K=(0,Z.Z)(C,H);return(0,M.jsx)(b.Z,(0,f.Z)({valueType:"password",fieldProps:u,proFieldProps:a,filedConfig:{valueType:t}},K))},de=X;de.Password=W,de.displayName="ProFormComponent",J.Z=de},48171:function(ye,J,r){"use strict";r.d(J,{J:function(){return n}});var f=r(85061),Z=r(67294),n=function(M){var O=(0,Z.useRef)(null);return O.current=M,(0,Z.useCallback)(function(){for(var H,t=arguments.length,X=new Array(t),W=0;W<t;W++)X[W]=arguments[W];return(H=O.current)===null||H===void 0?void 0:H.call.apply(H,[O].concat((0,f.Z)(X)))},[])}},22270:function(ye,J,r){"use strict";r.d(J,{h:function(){return f}});function f(Z){if(typeof Z=="function"){for(var n=arguments.length,b=new Array(n>1?n-1:0),M=1;M<n;M++)b[M-1]=arguments[M];return Z.apply(void 0,b)}return Z}}}]);