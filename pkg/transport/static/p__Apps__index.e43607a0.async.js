(self.webpackChunkidas_ui=self.webpackChunkidas_ui||[]).push([[577],{19273:function(se,$,t){"use strict";t.d($,{Z:function(){return m}});function m(s){if(s==null)throw new TypeError("Cannot destructure undefined")}},42088:function(se){se.exports={ListItem:"ListItem___osY4B"}},99376:function(se,$,t){"use strict";var m=t(71194),s=t(34883),D=t(402),W=t(97272),p=t(71153),l=t(60331),w=t(57663),J=t(71577),ie=t(34792),B=t(48086),S=t(39428),Q=t(11849),X=t(3182),F=t(2824),N=t(93224),b=t(96486),le=t.n(b),L=t(67294),ue=t(2844),z=t(51042),te=t(43471),Y=t(57119),oe=t(8978),ae=t(42088),ce=t.n(ae),O=t(85893),_e=["intl","request","metas","onEdit","onCreate","onClick"],ne=function(K){var g=K.intl,x=K.request,re=K.metas,k=K.onEdit,u=K.onCreate,n=K.onClick,_=(0,N.Z)(K,_e),v=(0,L.useRef)(),r=(0,L.useState)(),i=(0,F.Z)(r,2),e=i[0],a=i[1],P=(0,L.useState)(),y=(0,F.Z)(P,2),M=y[0],Z=y[1],h=(0,L.useCallback)(function(){var c=(0,X.Z)((0,S.Z)().mark(function o(d){return(0,S.Z)().wrap(function(U){for(;;)switch(U.prev=U.next){case 0:return U.abrupt("return",((0,b.isFunction)(x)?x:x.list)((0,Q.Z)((0,Q.Z)({},d),{},{keywords:e??void 0})).then(function(T){return T.total===void 0||T.pageSize===void 0||T.current===void 0?Z(!0):T.total<T.current*T.pageSize?Z(!1):Z(!0),T}));case 1:case"end":return U.stop()}},o)}));return function(o){return c.apply(this,arguments)}}(),[x,e]),f=function(o,d){var j=arguments.length>2&&arguments[2]!==void 0?arguments[2]:{},U=j.processing,T=j.success,I=j.error,q=o??(!(0,b.isFunction)(x)&&d&&x.patch?function(){var ee=(0,X.Z)((0,S.Z)().mark(function G(H){var V;return(0,S.Z)().wrap(function(R){for(;;)switch(R.prev=R.next){case 0:return R.abrupt("return",(V=x.patch)===null||V===void 0?void 0:V.call(x,H.map(d)));case 1:case"end":return R.stop()}},G)}));return function(G){return ee.apply(this,arguments)}}():null);if(!!q)return function(){var ee=(0,X.Z)((0,S.Z)().mark(function G(H){var V;return(0,S.Z)().wrap(function(R){for(;;)switch(R.prev=R.next){case 0:return V=B.default.loading(U??g.t("message.processing","Processing ...")),R.prev=1,R.next=4,q((0,b.isString)(H)?[{id:H}]:H.map(function(me){return{id:me}}));case 4:return V(),B.default.success(T??g.t("message.operationSuccessd","Operation succeeded.")),R.abrupt("return",!0);case 9:return R.prev=9,R.t0=R.catch(1),V(),B.default.error(I??g.t("message.operationFailed","Operation failed, please try again.")),R.abrupt("return",!1);case 14:case"end":return R.stop()}},G,null,[[1,9]])}));return function(G){return ee.apply(this,arguments)}}()},E=f((0,b.isFunction)(x)?void 0:x.delete,function(c){var o=c.id;return{id:o,isDelete:!0}},{processing:g.t("message.removing","Removing ..."),success:g.t("message.removeSuccessd","Remove successfully and will refresh soon"),error:g.t("message.removeFailed","Remove failed, please try again")}),A=f((0,b.isFunction)(x)?void 0:x.disable,function(c){var o=c.id;return{id:o,isDisable:!0}},{processing:g.t("message.disabling","Disabling ..."),success:g.t("message.disableSuccessd","Disabled successfully and will refresh soon"),error:g.t("message.disableFailed","Disable failed, please try again")}),C=f((0,b.isFunction)(x)?void 0:x.enable,function(c){var o=c.id;return{id:o,isDisable:!1}},{processing:g.t("message.enabling","Enabling ..."),success:g.t("message.enableSuccessd","Enabled successfully and will refresh soon"),error:g.t("message.enableFailed","Enable failed, please try again")});return(0,L.useEffect)(function(){var c;(c=v.current)===null||c===void 0||c.reload()},[h]),(0,O.jsx)(oe.ZP,(0,Q.Z)({pagination:M?{defaultPageSize:20,showSizeChanger:!0}:!1,actionRef:v,showActions:"always",grid:{gutter:16,column:3,xxl:4,xl:4,lg:4,md:3,sm:2,xs:2},toolbar:{search:!0,onSearch:function(o){var d;a(o),(d=v.current)===null||d===void 0||d.reload()},actions:[(0,O.jsxs)(J.Z,{hidden:!u,type:"primary",onClick:u,children:[(0,O.jsx)(z.Z,{}),g.t("button.create","Create")]},"create")],settings:[{icon:(0,O.jsx)(te.Z,{onClick:function(){var o;return(o=v.current)===null||o===void 0?void 0:o.reload()}}),key:"reload",onClick:function(){var o;return(o=v.current)===null||o===void 0?void 0:o.reload()}}]},itemCardProps:{className:ce().ListItem},onItem:function(o){return{onClick:function(){n==null||n(o)}}},metas:(0,Q.Z)({title:{render:function(o,d){var j;return(j=d.displayName)!==null&&j!==void 0?j:d.name}},subTitle:{render:function(o,d){return d.isDisable?(0,O.jsx)(l.Z,{color:"red",children:g.t("button.disabled","Disabled")}):(0,O.jsx)(O.Fragment,{})}},content:{render:function(o,d){var j;return(0,O.jsx)(W.Z.Paragraph,{ellipsis:{rows:2,tooltip:d.description},children:(j=d.description)!==null&&j!==void 0?j:""})}},avatar:{dataIndex:"avatar",search:!1,render:function(o,d){return(0,b.isString)(d.avatar)?(0,O.jsx)(ue.ZP,{size:"default",src:d.avatar}):d.avatar}},actions:k||E||C||A?{cardActionProps:"actions",render:function(o,d){return[k?(0,O.jsx)("a",{onClick:function(){k==null||k(d)},children:g.t("button.edit","Edit")},"edit"):null,E?(0,O.jsx)("a",{onClick:function(){s.Z.confirm({title:g.t("confirm.remove","Are you sure you want to remove this item?"),icon:(0,O.jsx)(Y.Z,{}),onOk:function(){E(d.id)},maskClosable:!0})},children:g.t("button.remove","Remove")},"remove"):null,d.isDisable!==void 0&&(d.isDisable?C:A)?(0,O.jsx)("a",{onClick:function(){s.Z.confirm({title:d.isDisable?g.t("confirm.enable","Are you sure you want to enable this item?"):g.t("confirm.disable","Are you sure you want to disable this item?"),icon:(0,O.jsx)(Y.Z,{}),onOk:function(){var T;(T=d.isDisable?C:A)===null||T===void 0||T(d.id)},maskClosable:!0})},children:d.isDisable?g.t("button.enable","Enable"):g.t("button.disable","Disable")},"disable_or_enable"):void 0]}}:void 0},re),request:h},_))};$.Z=ne},37790:function(se,$,t){"use strict";t.r($);var m=t(39428),s=t(3182),D=t(19273),W=t(67294),p=t(21010),l=t(99376),w=t(44728),J=t(93400),ie=t(75362),B=t(85893),S=function(X){(0,D.Z)(X);var F=new J.f("pages.apps",(0,p.YB)());return(0,B.jsx)(ie.ZP,{children:(0,B.jsx)(l.Z,{intl:F,request:{list:w.C6,delete:function(){var N=(0,s.Z)((0,m.Z)().mark(function le(L){return(0,m.Z)().wrap(function(z){for(;;)switch(z.prev=z.next){case 0:L.forEach(function(te){var Y=te.id;(0,w.P2)({id:Y})});case 1:case"end":return z.stop()}},le)}));function b(le){return N.apply(this,arguments)}return b}()},onEdit:function(b){p.m8.push("/apps/".concat(b.id,"/edit"))},onCreate:function(){p.m8.push("/apps/create")},onClick:function(b){p.m8.push("/apps/".concat(b.id))}})})};$.default=S},44728:function(se,$,t){"use strict";t.d($,{C6:function(){return X},ri:function(){return N},BN:function(){return te},KT:function(){return oe},P2:function(){return ce},qy:function(){return de},Wo:function(){return g},Rk:function(){return re}});var m=t(93224),s=t(39428),D=t(11849),W=t(3182),p=t(72709),l=["id"],w=["id"],J=["id"],ie=null,B=["appId"],S=["appId"],Q=["appId"];function X(u,n){return F.apply(this,arguments)}function F(){return F=(0,W.Z)((0,s.Z)().mark(function u(n,_){return(0,s.Z)().wrap(function(r){for(;;)switch(r.prev=r.next){case 0:return r.abrupt("return",(0,p.WY)("/api/v1/apps",(0,D.Z)({method:"GET",params:(0,D.Z)({},n)},_||{})));case 1:case"end":return r.stop()}},u)})),F.apply(this,arguments)}function N(u,n){return b.apply(this,arguments)}function b(){return b=(0,W.Z)((0,s.Z)().mark(function u(n,_){return(0,s.Z)().wrap(function(r){for(;;)switch(r.prev=r.next){case 0:return r.abrupt("return",(0,p.WY)("/api/v1/apps",(0,D.Z)({method:"POST",headers:{"Content-Type":"application/json"},data:n},_||{})));case 1:case"end":return r.stop()}},u)})),b.apply(this,arguments)}function le(u,n){return L.apply(this,arguments)}function L(){return L=_asyncToGenerator(_regeneratorRuntime().mark(function u(n,_){return _regeneratorRuntime().wrap(function(r){for(;;)switch(r.prev=r.next){case 0:return r.abrupt("return",request("/api/v1/apps",_objectSpread({method:"DELETE",headers:{"Content-Type":"application/json"},data:n},_||{})));case 1:case"end":return r.stop()}},u)})),L.apply(this,arguments)}function ue(u,n){return z.apply(this,arguments)}function z(){return z=_asyncToGenerator(_regeneratorRuntime().mark(function u(n,_){return _regeneratorRuntime().wrap(function(r){for(;;)switch(r.prev=r.next){case 0:return r.abrupt("return",request("/api/v1/apps",_objectSpread({method:"PATCH",headers:{"Content-Type":"application/json"},data:n},_||{})));case 1:case"end":return r.stop()}},u)})),z.apply(this,arguments)}function te(u,n){return Y.apply(this,arguments)}function Y(){return Y=(0,W.Z)((0,s.Z)().mark(function u(n,_){var v,r;return(0,s.Z)().wrap(function(e){for(;;)switch(e.prev=e.next){case 0:return v=n.id,r=(0,m.Z)(n,l),e.abrupt("return",(0,p.WY)("/api/v1/apps/".concat(v),(0,D.Z)({method:"GET",params:(0,D.Z)({},r)},_||{})));case 2:case"end":return e.stop()}},u)})),Y.apply(this,arguments)}function oe(u,n,_){return ae.apply(this,arguments)}function ae(){return ae=(0,W.Z)((0,s.Z)().mark(function u(n,_,v){var r,i;return(0,s.Z)().wrap(function(a){for(;;)switch(a.prev=a.next){case 0:return r=n.id,i=(0,m.Z)(n,w),a.abrupt("return",(0,p.WY)("/api/v1/apps/".concat(r),(0,D.Z)({method:"PUT",headers:{"Content-Type":"application/json"},params:(0,D.Z)({},i),data:_},v||{})));case 2:case"end":return a.stop()}},u)})),ae.apply(this,arguments)}function ce(u,n){return O.apply(this,arguments)}function O(){return O=(0,W.Z)((0,s.Z)().mark(function u(n,_){var v,r;return(0,s.Z)().wrap(function(e){for(;;)switch(e.prev=e.next){case 0:return v=n.id,r=(0,m.Z)(n,J),e.abrupt("return",(0,p.WY)("/api/v1/apps/".concat(v),(0,D.Z)({method:"DELETE",params:(0,D.Z)({},r)},_||{})));case 2:case"end":return e.stop()}},u)})),O.apply(this,arguments)}function _e(u,n,_){return ne.apply(this,arguments)}function ne(){return ne=_asyncToGenerator(_regeneratorRuntime().mark(function u(n,_,v){var r,i;return _regeneratorRuntime().wrap(function(a){for(;;)switch(a.prev=a.next){case 0:return r=n.id,i=_objectWithoutProperties(n,ie),a.abrupt("return",request("/api/v1/apps/".concat(r),_objectSpread({method:"PATCH",headers:{"Content-Type":"application/json"},params:_objectSpread({},i),data:_},v||{})));case 2:case"end":return a.stop()}},u)})),ne.apply(this,arguments)}function de(u,n){return K.apply(this,arguments)}function K(){return K=(0,W.Z)((0,s.Z)().mark(function u(n,_){var v,r;return(0,s.Z)().wrap(function(e){for(;;)switch(e.prev=e.next){case 0:return v=n.appId,r=(0,m.Z)(n,B),e.abrupt("return",(0,p.WY)("/api/v1/apps/".concat(v,"/key"),(0,D.Z)({method:"GET",params:(0,D.Z)({},r)},_||{})));case 2:case"end":return e.stop()}},u)})),K.apply(this,arguments)}function g(u,n,_){return x.apply(this,arguments)}function x(){return x=(0,W.Z)((0,s.Z)().mark(function u(n,_,v){var r,i;return(0,s.Z)().wrap(function(a){for(;;)switch(a.prev=a.next){case 0:return r=n.appId,i=(0,m.Z)(n,S),a.abrupt("return",(0,p.WY)("/api/v1/apps/".concat(r,"/key"),(0,D.Z)({method:"POST",headers:{"Content-Type":"application/json"},params:(0,D.Z)({},i),data:_},v||{})));case 2:case"end":return a.stop()}},u)})),x.apply(this,arguments)}function re(u,n,_){return k.apply(this,arguments)}function k(){return k=(0,W.Z)((0,s.Z)().mark(function u(n,_,v){var r,i;return(0,s.Z)().wrap(function(a){for(;;)switch(a.prev=a.next){case 0:return r=n.appId,i=(0,m.Z)(n,Q),a.abrupt("return",(0,p.WY)("/api/v1/apps/".concat(r,"/key"),(0,D.Z)({method:"DELETE",headers:{"Content-Type":"application/json"},params:(0,D.Z)({},i),data:_},v||{})));case 2:case"end":return a.stop()}},u)})),k.apply(this,arguments)}},90860:function(se,$,t){"use strict";t.d($,{Z:function(){return r}});var m=t(96156),s=t(22122),D=t(90484),W=t(94184),p=t.n(W),l=t(67294),w=t(53124),J=t(98423),ie=function(e){var a,P,y=e.prefixCls,M=e.className,Z=e.style,h=e.size,f=e.shape,E=p()((a={},(0,m.Z)(a,"".concat(y,"-lg"),h==="large"),(0,m.Z)(a,"".concat(y,"-sm"),h==="small"),a)),A=p()((P={},(0,m.Z)(P,"".concat(y,"-circle"),f==="circle"),(0,m.Z)(P,"".concat(y,"-square"),f==="square"),(0,m.Z)(P,"".concat(y,"-round"),f==="round"),P)),C=l.useMemo(function(){return typeof h=="number"?{width:h,height:h,lineHeight:"".concat(h,"px")}:{}},[h]);return l.createElement("span",{className:p()(y,E,A,M),style:(0,s.Z)((0,s.Z)({},C),Z)})},B=ie,S=function(e){var a=e.prefixCls,P=e.className,y=e.active,M=e.shape,Z=M===void 0?"circle":M,h=e.size,f=h===void 0?"default":h,E=l.useContext(w.E_),A=E.getPrefixCls,C=A("skeleton",a),c=(0,J.Z)(e,["prefixCls","className"]),o=p()(C,"".concat(C,"-element"),(0,m.Z)({},"".concat(C,"-active"),y),P);return l.createElement("div",{className:o},l.createElement(B,(0,s.Z)({prefixCls:"".concat(C,"-avatar"),shape:Z,size:f},c)))},Q=S,X=function(e){var a,P=e.prefixCls,y=e.className,M=e.active,Z=e.block,h=Z===void 0?!1:Z,f=e.size,E=f===void 0?"default":f,A=l.useContext(w.E_),C=A.getPrefixCls,c=C("skeleton",P),o=(0,J.Z)(e,["prefixCls"]),d=p()(c,"".concat(c,"-element"),(a={},(0,m.Z)(a,"".concat(c,"-active"),M),(0,m.Z)(a,"".concat(c,"-block"),h),a),y);return l.createElement("div",{className:d},l.createElement(B,(0,s.Z)({prefixCls:"".concat(c,"-button"),size:E},o)))},F=X,N=t(28991),b={icon:{tag:"svg",attrs:{viewBox:"64 64 896 896",focusable:"false"},children:[{tag:"path",attrs:{d:"M888 792H200V168c0-4.4-3.6-8-8-8h-56c-4.4 0-8 3.6-8 8v688c0 4.4 3.6 8 8 8h752c4.4 0 8-3.6 8-8v-56c0-4.4-3.6-8-8-8zM288 604a64 64 0 10128 0 64 64 0 10-128 0zm118-224a48 48 0 1096 0 48 48 0 10-96 0zm158 228a96 96 0 10192 0 96 96 0 10-192 0zm148-314a56 56 0 10112 0 56 56 0 10-112 0z"}}]},name:"dot-chart",theme:"outlined"},le=b,L=t(27029),ue=function(e,a){return l.createElement(L.Z,(0,N.Z)((0,N.Z)({},e),{},{ref:a,icon:le}))};ue.displayName="DotChartOutlined";var z=l.forwardRef(ue),te=function(e){var a=e.prefixCls,P=e.className,y=e.style,M=e.active,Z=e.children,h=l.useContext(w.E_),f=h.getPrefixCls,E=f("skeleton",a),A=p()(E,"".concat(E,"-element"),(0,m.Z)({},"".concat(E,"-active"),M),P),C=Z??l.createElement(z,null);return l.createElement("div",{className:A},l.createElement("div",{className:p()("".concat(E,"-image"),P),style:y},C))},Y=te,oe="M365.714286 329.142857q0 45.714286-32.036571 77.677714t-77.677714 32.036571-77.677714-32.036571-32.036571-77.677714 32.036571-77.677714 77.677714-32.036571 77.677714 32.036571 32.036571 77.677714zM950.857143 548.571429l0 256-804.571429 0 0-109.714286 182.857143-182.857143 91.428571 91.428571 292.571429-292.571429zM1005.714286 146.285714l-914.285714 0q-7.460571 0-12.873143 5.412571t-5.412571 12.873143l0 694.857143q0 7.460571 5.412571 12.873143t12.873143 5.412571l914.285714 0q7.460571 0 12.873143-5.412571t5.412571-12.873143l0-694.857143q0-7.460571-5.412571-12.873143t-12.873143-5.412571zM1097.142857 164.571429l0 694.857143q0 37.741714-26.843429 64.585143t-64.585143 26.843429l-914.285714 0q-37.741714 0-64.585143-26.843429t-26.843429-64.585143l0-694.857143q0-37.741714 26.843429-64.585143t64.585143-26.843429l914.285714 0q37.741714 0 64.585143 26.843429t26.843429 64.585143z",ae=function(e){var a=e.prefixCls,P=e.className,y=e.style,M=e.active,Z=l.useContext(w.E_),h=Z.getPrefixCls,f=h("skeleton",a),E=p()(f,"".concat(f,"-element"),(0,m.Z)({},"".concat(f,"-active"),M),P);return l.createElement("div",{className:E},l.createElement("div",{className:p()("".concat(f,"-image"),P),style:y},l.createElement("svg",{viewBox:"0 0 1098 1024",xmlns:"http://www.w3.org/2000/svg",className:"".concat(f,"-image-svg")},l.createElement("path",{d:oe,className:"".concat(f,"-image-path")}))))},ce=ae,O=function(e){var a,P=e.prefixCls,y=e.className,M=e.active,Z=e.block,h=e.size,f=h===void 0?"default":h,E=l.useContext(w.E_),A=E.getPrefixCls,C=A("skeleton",P),c=(0,J.Z)(e,["prefixCls"]),o=p()(C,"".concat(C,"-element"),(a={},(0,m.Z)(a,"".concat(C,"-active"),M),(0,m.Z)(a,"".concat(C,"-block"),Z),a),y);return l.createElement("div",{className:o},l.createElement(B,(0,s.Z)({prefixCls:"".concat(C,"-input"),size:f},c)))},_e=O,ne=t(85061),de=function(e){var a=function(E){var A=e.width,C=e.rows,c=C===void 0?2:C;if(Array.isArray(A))return A[E];if(c-1===E)return A},P=e.prefixCls,y=e.className,M=e.style,Z=e.rows,h=(0,ne.Z)(Array(Z)).map(function(f,E){return l.createElement("li",{key:E,style:{width:a(E)}})});return l.createElement("ul",{className:p()(P,y),style:M},h)},K=de,g=function(e){var a=e.prefixCls,P=e.className,y=e.width,M=e.style;return l.createElement("h3",{className:p()(a,P),style:(0,s.Z)({width:y},M)})},x=g;function re(i){return i&&(0,D.Z)(i)==="object"?i:{}}function k(i,e){return i&&!e?{size:"large",shape:"square"}:{size:"large",shape:"circle"}}function u(i,e){return!i&&e?{width:"38%"}:i&&e?{width:"50%"}:{}}function n(i,e){var a={};return(!i||!e)&&(a.width="61%"),!i&&e?a.rows=3:a.rows=2,a}var _=function(e){var a=e.prefixCls,P=e.loading,y=e.className,M=e.style,Z=e.children,h=e.avatar,f=h===void 0?!1:h,E=e.title,A=E===void 0?!0:E,C=e.paragraph,c=C===void 0?!0:C,o=e.active,d=e.round,j=l.useContext(w.E_),U=j.getPrefixCls,T=j.direction,I=U("skeleton",a);if(P||!("loading"in e)){var q,ee=!!f,G=!!A,H=!!c,V;if(ee){var pe=(0,s.Z)((0,s.Z)({prefixCls:"".concat(I,"-avatar")},k(G,H)),re(f));V=l.createElement("div",{className:"".concat(I,"-header")},l.createElement(B,(0,s.Z)({},pe)))}var R;if(G||H){var me;if(G){var fe=(0,s.Z)((0,s.Z)({prefixCls:"".concat(I,"-title")},u(ee,H)),re(A));me=l.createElement(x,(0,s.Z)({},fe))}var ve;if(H){var he=(0,s.Z)((0,s.Z)({prefixCls:"".concat(I,"-paragraph")},n(ee,G)),re(c));ve=l.createElement(K,(0,s.Z)({},he))}R=l.createElement("div",{className:"".concat(I,"-content")},me,ve)}var Ee=p()(I,(q={},(0,m.Z)(q,"".concat(I,"-with-avatar"),ee),(0,m.Z)(q,"".concat(I,"-active"),o),(0,m.Z)(q,"".concat(I,"-rtl"),T==="rtl"),(0,m.Z)(q,"".concat(I,"-round"),d),q),y);return l.createElement("div",{className:Ee,style:M},V,R)}return typeof Z!="undefined"?Z:null};_.Button=F,_.Avatar=Q,_.Input=_e,_.Image=ce,_.Node=Y;var v=_,r=v}}]);
