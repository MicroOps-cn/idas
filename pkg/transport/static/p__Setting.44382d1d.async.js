(self.webpackChunkidas_ui=self.webpackChunkidas_ui||[]).push([[868],{64335:function(T,B,i){"use strict";var p=i(67294),O=(0,p.createContext)({});B.Z=O},21349:function(T,B,i){"use strict";var p=i(84305),O=i(88182),Q=i(85893),H=i(94184),o=i.n(H),L=i(67294),v=i(64335),I=i(53645),R=i.n(I),u=function(S){var t=(0,L.useContext)(v.Z),E=S.children,St=S.contentWidth,vt=S.className,mt=S.style,ht=(0,L.useContext)(O.ZP.ConfigContext),pt=ht.getPrefixCls,ot=S.prefixCls||pt("pro"),lt=St||t.contentWidth,$t="".concat(ot,"-grid-content");return(0,Q.jsx)("div",{className:o()($t,vt,{wide:lt==="Fixed"}),style:mt,children:(0,Q.jsx)("div",{className:"".concat(ot,"-grid-content-children"),children:E})})};B.Z=u},19273:function(T,B,i){"use strict";i.d(B,{Z:function(){return p}});function p(O){if(O==null)throw new TypeError("Cannot destructure undefined")}},32108:function(T){T.exports={main:"main___1jT5e",leftMenu:"leftMenu___1QY4v",right:"right___2a0zh",title:"title___2n0kH"}},53645:function(){},58582:function(T,B,i){"use strict";i.r(B),i.d(B,{default:function(){return Xt}});var p=i(2824),O=i(19273),Q=i(30887),H=i(43084),o=i(67294),L=i(21010),v=i(93400),I=i(21349),R=i(54421),u=i(38272),U=i(22385),S=i(45777),t=i(39428),E=i(3182),St=i(77576),vt=i(12028),mt=i(69610),ht=i(54941),pt=i(81306),ot=i(6298),lt=i(86582),$t=i(49111),gt=i(19650),qt=i(57663),nt=i(71577),_t=i(47673),yt=i(28166),it=i(11849),te=i(34792),jt=i(48086),Lt=i(93224),It=i(28991),Ut={icon:function(s,a){return{tag:"svg",attrs:{viewBox:"64 64 896 896",focusable:"false"},children:[{tag:"path",attrs:{d:"M512 64C264.6 64 64 264.6 64 512s200.6 448 448 448 448-200.6 448-448S759.4 64 512 64zm0 820c-205.4 0-372-166.6-372-372s166.6-372 372-372 372 166.6 372 372-166.6 372-372 372z",fill:s}},{tag:"path",attrs:{d:"M512 140c-205.4 0-372 166.6-372 372s166.6 372 372 372 372-166.6 372-372-166.6-372-372-372zm193.4 225.7l-210.6 292a31.8 31.8 0 01-51.7 0L318.5 484.9c-3.8-5.3 0-12.7 6.5-12.7h46.9c10.3 0 19.9 5 25.9 13.3l71.2 98.8 157.2-218c6-8.4 15.7-13.3 25.9-13.3H699c6.5 0 10.3 7.4 6.4 12.7z",fill:a}},{tag:"path",attrs:{d:"M699 353h-46.9c-10.2 0-19.9 4.9-25.9 13.3L469 584.3l-71.2-98.8c-6-8.3-15.6-13.3-25.9-13.3H325c-6.5 0-10.3 7.4-6.5 12.7l124.6 172.8a31.8 31.8 0 0051.7 0l210.6-292c3.9-5.3.1-12.7-6.4-12.7z",fill:s}}]}},name:"check-circle",theme:"twotone"},Rt=Ut,Wt=i(27029),Et=function(s,a){return o.createElement(Wt.Z,(0,It.Z)((0,It.Z)({},s),{},{ref:a,icon:Rt}))};Et.displayName="CheckCircleTwoTone";var xt=o.forwardRef(Et),wt=i(73218),n=i(85893),Bt=["intl","onSave","autoSave","value","tooltip","prefix","suffix"],ut=function(s){var a=s.intl,d=s.onSave,e=s.autoSave,c=s.value,m=s.tooltip,g=s.prefix,h=s.suffix,x=(0,Lt.Z)(s,Bt),r=(0,o.useState)(c),M=(0,p.Z)(r,2),b=M[0],w=M[1],k=(0,o.useState)(!1),j=(0,p.Z)(k,2),X=j[0],F=j[1],A=function(){var W=(0,E.Z)((0,t.Z)().mark(function V(){return(0,t.Z)().wrap(function($){for(;;)switch($.prev=$.next){case 0:if($.t0=b!==void 0,!$.t0){$.next=5;break}return $.next=4,d(b);case 4:$.t0=$.sent;case 5:if(!$.t0){$.next=8;break}jt.default.info(a.t("finish","Update successful.")),F(!1);case 8:case"end":return $.stop()}},V)}));return function(){return W.apply(this,arguments)}}();return(0,o.useEffect)(function(){w(c)},[c]),(0,n.jsx)("div",{style:{display:"flex",alignItems:"center"},children:X?(0,n.jsxs)(gt.Z.Compact,{size:"small",children:[(0,n.jsx)(S.Z,{title:m?a.t("input.tooltip",m):void 0,children:(0,n.jsx)(yt.Z,(0,it.Z)((0,it.Z)({size:"small",min:1,defaultValue:c,onChange:function(V){x.type==="number"?w(parseInt(V.target.value,10)):w(V.target.value)},onBlur:function(){e&&A()},suffix:h?a.t("input.suffix",h):void 0,prefix:g?a.t("input.prefix",g):void 0},x),{},{style:(0,it.Z)({width:70},x.style)}))}),(0,n.jsx)(nt.Z,{size:"small",hidden:e,onClick:A,children:(0,n.jsx)(xt,{})}),(0,n.jsx)(nt.Z,{size:"small",onClick:function(){F(!1)},hidden:e,children:(0,n.jsx)(wt.Z,{})})]}):(0,n.jsx)("a",{style:{marginLeft:10},onClick:function(){F(!0)},children:a.t("modify","Modify")})})},Mt=function(s){var a=s.intl,d=s.onSave,e=s.suffix,c=s.prefix,m=s.count,g=s.autoSave,h=s.style,x=s.tooltip,r=s.value,M=s.type,b=(0,o.useState)(r??[]),w=(0,p.Z)(b,2),k=w[0],j=w[1],X=(0,o.useState)(!1),F=(0,p.Z)(X,2),A=F[0],W=F[1],V=function(){var Y=(0,E.Z)((0,t.Z)().mark(function K(){return(0,t.Z)().wrap(function(Z){for(;;)switch(Z.prev=Z.next){case 0:if(Z.t0=k,!Z.t0){Z.next=5;break}return Z.next=4,d(k);case 4:Z.t0=Z.sent;case 5:if(!Z.t0){Z.next=8;break}jt.default.info(a.t("finish","Update successful.")),W(!1);case 8:case"end":return Z.stop()}},K)}));return function(){return Y.apply(this,arguments)}}(),G=function(K,J){j([].concat((0,lt.Z)(k.slice(0,K)),[J],(0,lt.Z)(k.slice(K+1))))},$=function(){for(var K=[],J=function(D){var q,_,tt,et,st=(q=x==null?void 0:x[D])!==null&&q!==void 0?q:void 0,at=(_=e==null?void 0:e[D])!==null&&_!==void 0?_:void 0,y=(tt=c==null?void 0:c[D])!==null&&tt!==void 0?tt:void 0,f=(et=h==null?void 0:h[D])!==null&&et!==void 0?et:void 0;K.push((0,n.jsx)(S.Z,{title:st?a.t("input.".concat(D,".tooltip"),st):void 0,children:(0,n.jsx)(yt.Z,{style:(0,it.Z)({width:70},f),min:1,suffix:at?a.t("input.".concat(D,".suffix"),at):void 0,prefix:y?a.t("input.".concat(D,".prefix"),y):void 0,defaultValue:r==null?void 0:r[D],type:M,onChange:function(z){M==="number"?G(D,parseInt(z.target.value,10)):G(D,z.target.value)},onBlur:function(){g&&V()}})},"input-".concat(D)))},Z=0;Z<m;Z++)J(Z);return K};return(0,o.useEffect)(function(){j(r??[])},[r]),(0,n.jsx)("div",{style:{display:"flex",alignItems:"center"},children:A?(0,n.jsxs)(gt.Z.Compact,{size:"small",children:[$(),(0,n.jsx)(nt.Z,{hidden:g,onClick:V,children:(0,n.jsx)(xt,{})}),(0,n.jsx)(nt.Z,{onClick:function(){W(!1)},hidden:g,children:(0,n.jsx)(wt.Z,{})})]}):(0,n.jsx)("a",{style:{marginLeft:10},onClick:function(){W(!0)},children:a.t("modify","Modify")})})},bt=i(72709);function kt(C){return Ct.apply(this,arguments)}function Ct(){return Ct=(0,E.Z)((0,t.Z)().mark(function C(s){return(0,t.Z)().wrap(function(d){for(;;)switch(d.prev=d.next){case 0:return d.abrupt("return",(0,bt.WY)("/api/v1/config/security",(0,it.Z)({method:"GET"},s||{})));case 1:case"end":return d.stop()}},C)})),Ct.apply(this,arguments)}function Ot(C,s){return Zt.apply(this,arguments)}function Zt(){return Zt=(0,E.Z)((0,t.Z)().mark(function C(s,a){return(0,t.Z)().wrap(function(e){for(;;)switch(e.prev=e.next){case 0:return e.abrupt("return",(0,bt.WY)("/api/v1/config/security",(0,it.Z)({method:"PATCH",headers:{"Content-Type":"application/json"},data:s},a||{})));case 1:case"end":return e.stop()}},C)})),Zt.apply(this,arguments)}var Dt=i(1870),ee=i(43358),Vt=i(66074),Nt=i(96486),ct=i(14300),Kt=i(88293),Ht=function(s){var a=s.parentIntl,d=new v.f("password-complexity",a),e=[{label:d.t("option.unsafe","Unsafe"),desc:d.t("option.unsafe-description","Any character.")},{label:d.t("option.general","General"),desc:d.t("option.general-description","Composed of at least any two combinations of uppercase letters, lowercase letters, and numbers.")},{label:d.t("option.safe","Safe"),desc:d.t("option.safe-description","Must include uppercase and lowercase letters and numbers.")},{label:d.t("option.very_safe","Very Safe"),desc:d.t("option.very_safe-description","Must contain uppercase and lowercase letters, numbers, and special characters.")}];return(0,n.jsx)(S.Z,{placement:"bottomLeft",overlayStyle:{maxWidth:"max-content"},overlayInnerStyle:{backgroundColor:"rgba(61, 62, 64, 0.85)"},title:(0,n.jsx)("div",{children:e.map(function(c){return(0,n.jsxs)("li",{children:[c.label,": ",c.desc]},c.label)})}),children:(0,n.jsx)(Dt.Z,{style:{color:"rgba(61, 62, 64, 0.45)",marginInlineStart:4}})})},zt=function(s){return(0,Nt.isNumber)(s)?s:s!==void 0?ct.FI[s]:0},Ft=function(s){return(0,Nt.isNumber)(s)?ct.FI[s]:s!==void 0?s:ct.FI[0]},Yt=function(C){var s=C.parentIntl,a=C.value,d=C.onSave,e=new v.f("password-complexity",s),c=(0,o.useState)(!1),m=(0,p.Z)(c,2),g=m[0],h=m[1],x=(0,o.useState)(0),r=(0,p.Z)(x,2),M=r[0],b=r[1];return(0,o.useEffect)(function(){a!==void 0&&b(zt(a))},[a]),g?(0,n.jsx)(n.Fragment,{children:(0,n.jsxs)(gt.Z.Compact,{size:"small",children:[(0,n.jsx)(Vt.Z,{style:{width:120},defaultValue:zt(a),onChange:function(k){b(k)},options:(0,Kt.MM)(ct.FI,e,"option")}),(0,n.jsx)(nt.Z,{size:"small",onClick:(0,E.Z)((0,t.Z)().mark(function w(){return(0,t.Z)().wrap(function(j){for(;;)switch(j.prev=j.next){case 0:return j.next=2,d(M);case 2:if(!j.sent){j.next=4;break}h(!1);case 4:case"end":return j.stop()}},w)})),children:(0,n.jsx)(xt,{})}),(0,n.jsx)(nt.Z,{size:"small",onClick:function(){h(!1)},children:(0,n.jsx)(wt.Z,{})})]})}):(0,n.jsx)("a",{onClick:function(){h(!0)},children:e.t("modify","Modify")})},Pt=function(C){(0,pt.Z)(a,C);var s=(0,ot.Z)(a);function a(d){var e;return(0,mt.Z)(this,a),kt().then(function(c){e.setState({setting:c.data})}),e=s.call(this,d),e.state={intl:new v.f("security",e.props.parentIntl)},e.t=function(c,m,g,h){var x,r=e.state.intl;return(x=r==null?void 0:r.t(c,m,g,h))!==null&&x!==void 0?x:""},e.getData=function(){var c,m,g,h,x,r,M,b,w,k,j,X,F,A,W,V,G,$,Y,K,J,Z,ft,D,q,_,tt,et,st,at;return[{title:e.t("force-enable-mfa"),description:(0,n.jsx)(n.Fragment,{children:(c=e.state.setting)!==null&&c!==void 0&&c.forceEnableMfa?e.t("force-enable-mfa-enabled"):e.t("force-enable-mfa-disabled")}),actions:[(0,n.jsx)(vt.Z,{size:"small",checked:(m=e.state.setting)===null||m===void 0?void 0:m.forceEnableMfa,onChange:function(f){e.handleUpdateSetting({forceEnableMfa:f})}},"force-enable-mfa")]},{title:(0,n.jsxs)(n.Fragment,{children:[e.t("password-complexity"),(0,n.jsx)(Ht,{parentIntl:e.state.intl})]}),description:(0,n.jsx)(n.Fragment,{children:e.t("password-complexity.option.".concat(Ft((g=e.state.setting)===null||g===void 0?void 0:g.passwordComplexity),"-description"))}),actions:[(0,n.jsx)(Yt,{parentIntl:e.state.intl,onSave:function(){var y=(0,E.Z)((0,t.Z)().mark(function f(N){return(0,t.Z)().wrap(function(l){for(;;)switch(l.prev=l.next){case 0:return l.abrupt("return",e.handleUpdateSetting({passwordComplexity:N}));case 1:case"end":return l.stop()}},f)}));return function(f){return y.apply(this,arguments)}}(),value:(h=e.state.setting)===null||h===void 0?void 0:h.passwordComplexity},"password-complexity")]},{title:e.t("password-min-length"),description:(0,n.jsx)(n.Fragment,{children:(x=e.state.setting)!==null&&x!==void 0&&x.passwordMinLength?e.t("password-min-length-description","","",{minLen:e.state.setting.passwordMinLength}):e.t("password-min-length-unrestricted")}),actions:[(0,n.jsx)(ut,{type:"number",intl:new v.f("password-min-length",e.state.intl),onSave:function(){var y=(0,E.Z)((0,t.Z)().mark(function f(N){return(0,t.Z)().wrap(function(l){for(;;)switch(l.prev=l.next){case 0:return l.abrupt("return",e.handleUpdateSetting({passwordMinLength:N}));case 1:case"end":return l.stop()}},f)}));return function(f){return y.apply(this,arguments)}}(),value:(r=e.state.setting)===null||r===void 0?void 0:r.passwordMinLength},"password-min-length")]},{title:e.t("password-expire-time"),description:(0,n.jsx)(n.Fragment,{children:(M=e.state.setting)!==null&&M!==void 0&&M.passwordExpireTime?e.t("password-expire-time-description","","",{days:(b=e.state.setting)===null||b===void 0?void 0:b.passwordExpireTime}):e.t("password-expire-time-unrestricted")}),actions:[(0,n.jsx)(ut,{type:"number",intl:new v.f("password-expire-time",e.state.intl),onSave:function(){var y=(0,E.Z)((0,t.Z)().mark(function f(N){return(0,t.Z)().wrap(function(l){for(;;)switch(l.prev=l.next){case 0:return l.abrupt("return",e.handleUpdateSetting({passwordExpireTime:N}));case 1:case"end":return l.stop()}},f)}));return function(f){return y.apply(this,arguments)}}(),style:{width:100},suffix:"day",value:(w=e.state.setting)===null||w===void 0?void 0:w.passwordExpireTime},"password-expire-time")]},{title:e.t("password-failed-lock"),description:(0,n.jsx)(n.Fragment,{children:(k=e.state.setting)!==null&&k!==void 0&&k.passwordFailedLockDuration&&(j=e.state.setting)!==null&&j!==void 0&&j.passwordFailedLockThreshold?e.t("password-failed-lock-description","","",{min:(X=e.state.setting)===null||X===void 0?void 0:X.passwordFailedLockDuration,fails:(F=e.state.setting)===null||F===void 0?void 0:F.passwordFailedLockThreshold}):e.t("password-failed-lock-unrestricted")}),actions:[(0,n.jsx)(Mt,{type:"number",intl:new v.f("password-failed-lock",e.state.intl),onSave:function(){var y=(0,E.Z)((0,t.Z)().mark(function f(N){var z,l,rt;return(0,t.Z)().wrap(function(P){for(;;)switch(P.prev=P.next){case 0:return z=(0,p.Z)(N,2),l=z[0],rt=z[1],P.abrupt("return",e.handleUpdateSetting({passwordFailedLockThreshold:l,passwordFailedLockDuration:rt}));case 2:case"end":return P.stop()}},f)}));return function(f){return y.apply(this,arguments)}}(),suffix:["failed","minute"],style:[{width:120},{width:100}],tooltip:["Number of consecutive password input errors.","The duration of account lockout."],count:2,value:[(A=(W=e.state.setting)===null||W===void 0?void 0:W.passwordFailedLockThreshold)!==null&&A!==void 0?A:0,(V=(G=e.state.setting)===null||G===void 0?void 0:G.passwordFailedLockDuration)!==null&&V!==void 0?V:0]},"password-failed-lock")]},{title:e.t("password-history"),description:(0,n.jsx)(n.Fragment,{children:($=e.state.setting)!==null&&$!==void 0&&$.passwordHistory?e.t("password-history-description","","",{count:(Y=e.state.setting)===null||Y===void 0?void 0:Y.passwordHistory}):e.t("password-history-unrestricted")}),actions:[(0,n.jsx)(ut,{type:"number",style:{width:100},suffix:"day",intl:new v.f("password-history",e.state.intl),onSave:function(){var y=(0,E.Z)((0,t.Z)().mark(function f(N){return(0,t.Z)().wrap(function(l){for(;;)switch(l.prev=l.next){case 0:return l.abrupt("return",e.handleUpdateSetting({passwordHistory:N}));case 1:case"end":return l.stop()}},f)}));return function(f){return y.apply(this,arguments)}}(),value:(K=e.state.setting)===null||K===void 0?void 0:K.passwordHistory},"password-history")]},{title:e.t("account-inactive-lock"),description:(0,n.jsx)(n.Fragment,{children:(J=e.state.setting)!==null&&J!==void 0&&J.accountInactiveLock?e.t("account-inactive-lock-description","","",{days:(Z=e.state.setting)===null||Z===void 0?void 0:Z.accountInactiveLock}):e.t("account-inactive-lock-unrestricted")}),actions:[(0,n.jsx)(ut,{type:"number",style:{width:100},suffix:"day",intl:new v.f("account-inactive-lock",e.state.intl),onSave:function(){var y=(0,E.Z)((0,t.Z)().mark(function f(N){return(0,t.Z)().wrap(function(l){for(;;)switch(l.prev=l.next){case 0:return l.abrupt("return",e.handleUpdateSetting({accountInactiveLock:N}));case 1:case"end":return l.stop()}},f)}));return function(f){return y.apply(this,arguments)}}(),value:(ft=e.state.setting)===null||ft===void 0?void 0:ft.accountInactiveLock},"account-inactive-lock")]},{title:(0,n.jsxs)(n.Fragment,{children:[e.t("login-session-expiration-time"),(0,n.jsx)(S.Z,{placement:"bottomLeft",overlayStyle:{maxWidth:"max-content"},overlayInnerStyle:{backgroundColor:"rgba(61, 62, 64, 0.85)"},title:e.t("login-session-expiration-tooltip"),children:(0,n.jsx)(Dt.Z,{style:{color:"rgba(61, 62, 64, 0.45)",marginInlineStart:4}})})]}),description:(0,n.jsx)(n.Fragment,{children:(D=e.state.setting)!==null&&D!==void 0&&D.loginSessionInactivityTime?e.t("login-session-expiration-time-description","Automatically log out after {loginSessionInactivityHours} hours of inactivity, with a maximum session duration of {loginSessionMaxHours} hours.","",{loginSessionInactivityHours:(q=e.state.setting)===null||q===void 0?void 0:q.loginSessionInactivityTime,loginSessionMaxHours:(_=e.state.setting)===null||_===void 0?void 0:_.loginSessionMaxTime}):e.t("login-session-expiration-time-unrestricted")}),actions:[(0,n.jsx)(Mt,{type:"number",intl:new v.f("login-session-expiration-time",e.state.intl),onSave:function(){var y=(0,E.Z)((0,t.Z)().mark(function f(N){var z,l,rt;return(0,t.Z)().wrap(function(P){for(;;)switch(P.prev=P.next){case 0:return z=(0,p.Z)(N,2),l=z[0],rt=z[1],P.abrupt("return",e.handleUpdateSetting({loginSessionInactivityTime:l,loginSessionMaxTime:rt}));case 2:case"end":return P.stop()}},f)}));return function(f){return y.apply(this,arguments)}}(),suffix:["hours","hours"],style:[{width:100},{width:100}],tooltip:["Session inactive automatic logout time.","The maximum duration of the session."],count:2,value:[(tt=(et=e.state.setting)===null||et===void 0?void 0:et.loginSessionInactivityTime)!==null&&tt!==void 0?tt:0,(st=(at=e.state.setting)===null||at===void 0?void 0:at.loginSessionMaxTime)!==null&&st!==void 0?st:0]},"login-session-expiration-time")]}]},e.state={intl:new v.f("security",e.props.parentIntl)},e}return(0,ht.Z)(a,[{key:"handleUpdateSetting",value:function(){var d=(0,E.Z)((0,t.Z)().mark(function c(m){var g=this,h;return(0,t.Z)().wrap(function(r){for(;;)switch(r.prev=r.next){case 0:return r.next=2,Ot(m);case 2:return h=r.sent,kt().then(function(M){g.setState({setting:M.data})}),r.abrupt("return",h.success);case 5:case"end":return r.stop()}},c)}));function e(c){return d.apply(this,arguments)}return e}()},{key:"render",value:function(){var e=this.getData();return(0,n.jsx)(n.Fragment,{children:(0,n.jsx)(u.ZP,{itemLayout:"horizontal",dataSource:e,renderItem:function(m){return(0,n.jsx)(u.ZP.Item,{actions:m.actions,children:(0,n.jsx)(u.ZP.Item.Meta,{title:m.title,description:m.description})})}})})}}]),a}(o.Component),Qt=Pt,At=i(32108),dt=i.n(At),Jt=H.Z.Item,Tt=function(s){(0,O.Z)(s);var a=new v.f("settings",(0,L.YB)()),d=(0,o.useRef)(null),e={base:a.t("menuMap.basic","Basic Settings"),security:a.t("menuMap.security","Security Settings")},c=(0,o.useState)("security"),m=(0,p.Z)(c,2),g=m[0],h=m[1],x=function(){switch(g){case"security":return(0,n.jsx)(Qt,{parentIntl:a});default:break}return null},r=function(){return Object.keys(e).map(function(w){return(0,n.jsx)(Jt,{children:e[w]},w)})},M=function(){return e[g]};return(0,n.jsx)(I.Z,{children:(0,n.jsxs)("div",{className:dt().main,ref:d,children:[(0,n.jsx)("div",{className:dt().leftMenu,children:(0,n.jsx)(H.Z,{mode:"inline",selectedKeys:[g],onClick:function(w){var k=w.key;return h(k)},children:r()})}),(0,n.jsxs)("div",{className:dt().right,children:[(0,n.jsx)("div",{className:dt().title,children:M()}),x()]})]})})},Xt=Tt},14300:function(T,B,i){"use strict";i.d(B,{FI:function(){return p},Bz:function(){return H},w$:function(){return L},qJ:function(){return v},FQ:function(){return I},J0:function(){return u}});var p;(function(t){t[t.unsafe=0]="unsafe",t[t.general=1]="general",t[t.safe=2]="safe",t[t.very_safe=3]="very_safe"})(p||(p={}));var O;(function(t){t[t.refresh_token=0]="refresh_token",t[t.authorization_code=1]="authorization_code",t[t.password=2]="password",t[t.client_credentials=3]="client_credentials"})(O||(O={}));var Q;(function(t){t[t.default=0]="default",t[t.code=1]="code",t[t.token=2]="token"})(Q||(Q={}));var H;(function(t){t[t.oauth2=6]="oauth2",t[t.enable_mfa_sms=12]="enable_mfa_sms",t[t.normal=0]="normal",t[t.mfa_totp=1]="mfa_totp",t[t.mfa_email=2]="mfa_email",t[t.mfa_sms=3]="mfa_sms",t[t.email=4]="email",t[t.sms=5]="sms",t[t.enable_mfa_totp=10]="enable_mfa_totp",t[t.enable_mfa_email=11]="enable_mfa_email"})(H||(H={}));var o;(function(t){t[t.token_signature=3]="token_signature",t[t.basic=0]="basic",t[t.signature=1]="signature",t[t.token=2]="token"})(o||(o={}));var L;(function(t){t[t.unknown=0]="unknown",t[t.normal=1]="normal",t[t.disable=2]="disable"})(L||(L={}));var v;(function(t){t[t.client_credentials=8]="client_credentials",t[t.proxy=16]="proxy",t[t.oidc=32]="oidc",t[t.none=0]="none",t[t.authorization_code=1]="authorization_code",t[t.implicit=2]="implicit",t[t.password=4]="password"})(v||(v={}));var I;(function(t){t[t.manual=0]="manual",t[t.full=1]="full"})(I||(I={}));var R;(function(t){t[t.system=1]="system",t[t.user=0]="user"})(R||(R={}));var u;(function(t){t[t.disabled=1]="disabled",t[t.user_inactive=2]="user_inactive",t[t.password_expired=4]="password_expired",t[t.normal=0]="normal"})(u||(u={}));var U;(function(t){t[t.text=0]="text",t[t.digitRange=4]="digitRange",t[t.select=8]="select",t[t.dateTimeRange=14]="dateTimeRange",t[t.radio=6]="radio",t[t.multiSelect=9]="multiSelect",t[t.date=11]="date",t[t.dateRange=12]="dateRange",t[t.textarea=2]="textarea",t[t.timeRange=10]="timeRange",t[t.dateTime=13]="dateTime",t[t.digit=3]="digit",t[t.checkbox=5]="checkbox",t[t.switch=7]="switch"})(U||(U={}));var S;(function(t){t[t.all=0]="all",t[t.disabled=1]="disabled",t[t.enabled=2]="enabled"})(S||(S={}))},88293:function(T,B,i){"use strict";i.d(B,{MM:function(){return p},fs:function(){return O},GG:function(){return Q}});var p=function(o,L,v,I){var R=[];for(var u in o)if(Object.prototype.propertyIsEnumerable.call(o,u)&&isNaN(Number(u))){var U=o[u];if(I&&!I(u,U))continue;R.push({label:L.formatMessage({id:"".concat(v,".").concat(u),defaultMessage:u}),key:u,value:U})}return R},O=function(o,L,v,I){var R=new Map;for(var u in o)if(Object.prototype.hasOwnProperty.call(o,u)&&isNaN(Number(u))){var U=o[u];if(I&&!I(u,U))continue;R.set(U,L.formatMessage({id:"".concat(v,".").concat(u),defaultMessage:u}))}return R},Q=function(o,L,v,I){var R={};for(var u in o)if(Object.prototype.hasOwnProperty.call(o,u)&&!isNaN(Number(u))){var U,S=o[u];R[u]={text:L.formatMessage({id:"".concat(v,".").concat(S),defaultMessage:S}),status:(U=I[S])!==null&&U!==void 0?U:"Default"}}return R}}}]);