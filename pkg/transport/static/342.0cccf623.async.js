(self.webpackChunkidas_ui=self.webpackChunkidas_ui||[]).push([[342],{44728:function(H,A,b){"use strict";b.d(A,{C6:function(){return e},ri:function(){return I},BN:function(){return O},KT:function(){return $},P2:function(){return K},qy:function(){return s},Wo:function(){return Z},Rk:function(){return M}});var w=b(93224),u=b(39428),i=b(11849),l=b(3182),c=b(72709),D=["id"],R=["id"],E=["id"],j=null,_=["appId"],h=["appId"],C=["appId"];function e(r,t){return x.apply(this,arguments)}function x(){return x=(0,l.Z)((0,u.Z)().mark(function r(t,p){return(0,u.Z)().wrap(function(n){for(;;)switch(n.prev=n.next){case 0:return n.abrupt("return",(0,c.WY)("/api/v1/apps",(0,i.Z)({method:"GET",params:(0,i.Z)({},t)},p||{})));case 1:case"end":return n.stop()}},r)})),x.apply(this,arguments)}function I(r,t){return P.apply(this,arguments)}function P(){return P=(0,l.Z)((0,u.Z)().mark(function r(t,p){return(0,u.Z)().wrap(function(n){for(;;)switch(n.prev=n.next){case 0:return n.abrupt("return",(0,c.WY)("/api/v1/apps",(0,i.Z)({method:"POST",headers:{"Content-Type":"application/json"},data:t},p||{})));case 1:case"end":return n.stop()}},r)})),P.apply(this,arguments)}function y(r,t){return U.apply(this,arguments)}function U(){return U=_asyncToGenerator(_regeneratorRuntime().mark(function r(t,p){return _regeneratorRuntime().wrap(function(n){for(;;)switch(n.prev=n.next){case 0:return n.abrupt("return",request("/api/v1/apps",_objectSpread({method:"DELETE",headers:{"Content-Type":"application/json"},data:t},p||{})));case 1:case"end":return n.stop()}},r)})),U.apply(this,arguments)}function B(r,t){return k.apply(this,arguments)}function k(){return k=_asyncToGenerator(_regeneratorRuntime().mark(function r(t,p){return _regeneratorRuntime().wrap(function(n){for(;;)switch(n.prev=n.next){case 0:return n.abrupt("return",request("/api/v1/apps",_objectSpread({method:"PATCH",headers:{"Content-Type":"application/json"},data:t},p||{})));case 1:case"end":return n.stop()}},r)})),k.apply(this,arguments)}function O(r,t){return q.apply(this,arguments)}function q(){return q=(0,l.Z)((0,u.Z)().mark(function r(t,p){var d,n;return(0,u.Z)().wrap(function(v){for(;;)switch(v.prev=v.next){case 0:return d=t.id,n=(0,w.Z)(t,D),v.abrupt("return",(0,c.WY)("/api/v1/apps/".concat(d),(0,i.Z)({method:"GET",params:(0,i.Z)({},n)},p||{})));case 2:case"end":return v.stop()}},r)})),q.apply(this,arguments)}function $(r,t,p){return N.apply(this,arguments)}function N(){return N=(0,l.Z)((0,u.Z)().mark(function r(t,p,d){var n,W;return(0,u.Z)().wrap(function(f){for(;;)switch(f.prev=f.next){case 0:return n=t.id,W=(0,w.Z)(t,R),f.abrupt("return",(0,c.WY)("/api/v1/apps/".concat(n),(0,i.Z)({method:"PUT",headers:{"Content-Type":"application/json"},params:(0,i.Z)({},W),data:p},d||{})));case 2:case"end":return f.stop()}},r)})),N.apply(this,arguments)}function K(r,t){return Y.apply(this,arguments)}function Y(){return Y=(0,l.Z)((0,u.Z)().mark(function r(t,p){var d,n;return(0,u.Z)().wrap(function(v){for(;;)switch(v.prev=v.next){case 0:return d=t.id,n=(0,w.Z)(t,E),v.abrupt("return",(0,c.WY)("/api/v1/apps/".concat(d),(0,i.Z)({method:"DELETE",params:(0,i.Z)({},n)},p||{})));case 2:case"end":return v.stop()}},r)})),Y.apply(this,arguments)}function z(r,t,p){return o.apply(this,arguments)}function o(){return o=_asyncToGenerator(_regeneratorRuntime().mark(function r(t,p,d){var n,W;return _regeneratorRuntime().wrap(function(f){for(;;)switch(f.prev=f.next){case 0:return n=t.id,W=_objectWithoutProperties(t,j),f.abrupt("return",request("/api/v1/apps/".concat(n),_objectSpread({method:"PATCH",headers:{"Content-Type":"application/json"},params:_objectSpread({},W),data:p},d||{})));case 2:case"end":return f.stop()}},r)})),o.apply(this,arguments)}function s(r,t){return m.apply(this,arguments)}function m(){return m=(0,l.Z)((0,u.Z)().mark(function r(t,p){var d,n;return(0,u.Z)().wrap(function(v){for(;;)switch(v.prev=v.next){case 0:return d=t.appId,n=(0,w.Z)(t,_),v.abrupt("return",(0,c.WY)("/api/v1/apps/".concat(d,"/key"),(0,i.Z)({method:"GET",params:(0,i.Z)({},n)},p||{})));case 2:case"end":return v.stop()}},r)})),m.apply(this,arguments)}function Z(r,t,p){return a.apply(this,arguments)}function a(){return a=(0,l.Z)((0,u.Z)().mark(function r(t,p,d){var n,W;return(0,u.Z)().wrap(function(f){for(;;)switch(f.prev=f.next){case 0:return n=t.appId,W=(0,w.Z)(t,h),f.abrupt("return",(0,c.WY)("/api/v1/apps/".concat(n,"/key"),(0,i.Z)({method:"POST",headers:{"Content-Type":"application/json"},params:(0,i.Z)({},W),data:p},d||{})));case 2:case"end":return f.stop()}},r)})),a.apply(this,arguments)}function M(r,t,p){return g.apply(this,arguments)}function g(){return g=(0,l.Z)((0,u.Z)().mark(function r(t,p,d){var n,W;return(0,u.Z)().wrap(function(f){for(;;)switch(f.prev=f.next){case 0:return n=t.appId,W=(0,w.Z)(t,C),f.abrupt("return",(0,c.WY)("/api/v1/apps/".concat(n,"/key"),(0,i.Z)({method:"DELETE",headers:{"Content-Type":"application/json"},params:(0,i.Z)({},W),data:p},d||{})));case 2:case"end":return f.stop()}},r)})),g.apply(this,arguments)}},14300:function(H,A,b){"use strict";b.d(A,{FI:function(){return w},Bz:function(){return l},w$:function(){return D},qJ:function(){return R},FQ:function(){return E},J0:function(){return _}});var w;(function(e){e[e.unsafe=0]="unsafe",e[e.general=1]="general",e[e.safe=2]="safe",e[e.very_safe=3]="very_safe"})(w||(w={}));var u;(function(e){e[e.client_credentials=3]="client_credentials",e[e.refresh_token=0]="refresh_token",e[e.authorization_code=1]="authorization_code",e[e.password=2]="password"})(u||(u={}));var i;(function(e){e[e.default=0]="default",e[e.code=1]="code",e[e.token=2]="token"})(i||(i={}));var l;(function(e){e[e.mfa_email=2]="mfa_email",e[e.mfa_sms=3]="mfa_sms",e[e.email=4]="email",e[e.enable_mfa_email=11]="enable_mfa_email",e[e.enable_mfa_sms=12]="enable_mfa_sms",e[e.normal=0]="normal",e[e.mfa_totp=1]="mfa_totp",e[e.sms=5]="sms",e[e.oauth2=6]="oauth2",e[e.enable_mfa_totp=10]="enable_mfa_totp"})(l||(l={}));var c;(function(e){e[e.token_signature=3]="token_signature",e[e.basic=0]="basic",e[e.signature=1]="signature",e[e.token=2]="token"})(c||(c={}));var D;(function(e){e[e.unknown=0]="unknown",e[e.normal=1]="normal",e[e.disable=2]="disable"})(D||(D={}));var R;(function(e){e[e.none=0]="none",e[e.authorization_code=1]="authorization_code",e[e.implicit=2]="implicit",e[e.password=4]="password",e[e.client_credentials=8]="client_credentials",e[e.proxy=16]="proxy",e[e.oidc=32]="oidc",e[e.radius=64]="radius"})(R||(R={}));var E;(function(e){e[e.manual=0]="manual",e[e.full=1]="full"})(E||(E={}));var j;(function(e){e[e.user=0]="user",e[e.system=1]="system"})(j||(j={}));var _;(function(e){e[e.user_inactive=2]="user_inactive",e[e.password_expired=4]="password_expired",e[e.normal=0]="normal",e[e.disabled=1]="disabled"})(_||(_={}));var h;(function(e){e[e.dateRange=12]="dateRange",e[e.textarea=2]="textarea",e[e.digit=3]="digit",e[e.select=8]="select",e[e.timeRange=10]="timeRange",e[e.digitRange=4]="digitRange",e[e.radio=6]="radio",e[e.switch=7]="switch",e[e.dateTime=13]="dateTime",e[e.dateTimeRange=14]="dateTimeRange",e[e.text=0]="text",e[e.checkbox=5]="checkbox",e[e.multiSelect=9]="multiSelect",e[e.date=11]="date"})(h||(h={}));var C;(function(e){e[e.enabled=2]="enabled",e[e.all=0]="all",e[e.disabled=1]="disabled"})(C||(C={}))},10861:function(H,A,b){"use strict";b.d(A,{Rf:function(){return _},r4:function(){return C},Vt:function(){return x},Q8:function(){return P},bG:function(){return U},Nq:function(){return k},Pw:function(){return Y}});var w=b(93224),u=b(39428),i=b(11849),l=b(3182),c=b(72709),D=["id"],R=["id"],E=null,j=null;function _(o,s){return h.apply(this,arguments)}function h(){return h=(0,l.Z)((0,u.Z)().mark(function o(s,m){return(0,u.Z)().wrap(function(a){for(;;)switch(a.prev=a.next){case 0:return a.abrupt("return",(0,c.WY)("/api/v1/users",(0,i.Z)({method:"GET",params:(0,i.Z)({},s)},m||{})));case 1:case"end":return a.stop()}},o)})),h.apply(this,arguments)}function C(o,s){return e.apply(this,arguments)}function e(){return e=(0,l.Z)((0,u.Z)().mark(function o(s,m){return(0,u.Z)().wrap(function(a){for(;;)switch(a.prev=a.next){case 0:return a.abrupt("return",(0,c.WY)("/api/v1/users",(0,i.Z)({method:"POST",headers:{"Content-Type":"application/json"},data:s},m||{})));case 1:case"end":return a.stop()}},o)})),e.apply(this,arguments)}function x(o,s){return I.apply(this,arguments)}function I(){return I=(0,l.Z)((0,u.Z)().mark(function o(s,m){return(0,u.Z)().wrap(function(a){for(;;)switch(a.prev=a.next){case 0:return a.abrupt("return",(0,c.WY)("/api/v1/users",(0,i.Z)({method:"DELETE",headers:{"Content-Type":"application/json"},data:s},m||{})));case 1:case"end":return a.stop()}},o)})),I.apply(this,arguments)}function P(o,s){return y.apply(this,arguments)}function y(){return y=(0,l.Z)((0,u.Z)().mark(function o(s,m){return(0,u.Z)().wrap(function(a){for(;;)switch(a.prev=a.next){case 0:return a.abrupt("return",(0,c.WY)("/api/v1/users",(0,i.Z)({method:"PATCH",headers:{"Content-Type":"application/json"},data:s},m||{})));case 1:case"end":return a.stop()}},o)})),y.apply(this,arguments)}function U(o,s){return B.apply(this,arguments)}function B(){return B=(0,l.Z)((0,u.Z)().mark(function o(s,m){var Z,a;return(0,u.Z)().wrap(function(g){for(;;)switch(g.prev=g.next){case 0:return Z=s.id,a=(0,w.Z)(s,D),g.abrupt("return",(0,c.WY)("/api/v1/users/".concat(Z),(0,i.Z)({method:"GET",params:(0,i.Z)({},a)},m||{})));case 2:case"end":return g.stop()}},o)})),B.apply(this,arguments)}function k(o,s,m){return O.apply(this,arguments)}function O(){return O=(0,l.Z)((0,u.Z)().mark(function o(s,m,Z){var a,M;return(0,u.Z)().wrap(function(r){for(;;)switch(r.prev=r.next){case 0:return a=s.id,M=(0,w.Z)(s,R),r.abrupt("return",(0,c.WY)("/api/v1/users/".concat(a),(0,i.Z)({method:"PUT",headers:{"Content-Type":"application/json"},params:(0,i.Z)({},M),data:m},Z||{})));case 2:case"end":return r.stop()}},o)})),O.apply(this,arguments)}function q(o,s){return $.apply(this,arguments)}function $(){return $=_asyncToGenerator(_regeneratorRuntime().mark(function o(s,m){var Z,a;return _regeneratorRuntime().wrap(function(g){for(;;)switch(g.prev=g.next){case 0:return Z=s.id,a=_objectWithoutProperties(s,E),g.abrupt("return",request("/api/v1/users/".concat(Z),_objectSpread({method:"DELETE",params:_objectSpread({},a)},m||{})));case 2:case"end":return g.stop()}},o)})),$.apply(this,arguments)}function N(o,s,m){return K.apply(this,arguments)}function K(){return K=_asyncToGenerator(_regeneratorRuntime().mark(function o(s,m,Z){var a,M;return _regeneratorRuntime().wrap(function(r){for(;;)switch(r.prev=r.next){case 0:return a=s.id,M=_objectWithoutProperties(s,j),r.abrupt("return",request("/api/v1/users/".concat(a),_objectSpread({method:"PATCH",headers:{"Content-Type":"application/json"},params:_objectSpread({},M),data:m},Z||{})));case 2:case"end":return r.stop()}},o)})),K.apply(this,arguments)}function Y(o,s){return z.apply(this,arguments)}function z(){return z=(0,l.Z)((0,u.Z)().mark(function o(s,m){return(0,u.Z)().wrap(function(a){for(;;)switch(a.prev=a.next){case 0:return a.abrupt("return",(0,c.WY)("/api/v1/users/sendActivateMail",(0,i.Z)({method:"POST",headers:{"Content-Type":"application/json"},data:s},m||{})));case 1:case"end":return a.stop()}},o)})),z.apply(this,arguments)}},88293:function(H,A,b){"use strict";b.d(A,{MM:function(){return w},fs:function(){return u},GG:function(){return i}});var w=function(c,D,R,E){var j=[];for(var _ in c)if(Object.prototype.propertyIsEnumerable.call(c,_)&&isNaN(Number(_))){var h=c[_];if(E&&!E(_,h))continue;j.push({label:D.formatMessage({id:"".concat(R,".").concat(_),defaultMessage:_}),key:_,value:h})}return j},u=function(c,D,R,E){var j=new Map;for(var _ in c)if(Object.prototype.hasOwnProperty.call(c,_)&&isNaN(Number(_))){var h=c[_];if(E&&!E(_,h))continue;j.set(h,D.formatMessage({id:"".concat(R,".").concat(_),defaultMessage:_}))}return j},i=function(c,D,R,E){var j={};for(var _ in c)if(Object.prototype.hasOwnProperty.call(c,_)&&!isNaN(Number(_))){var h,C=c[_];j[_]={text:D.formatMessage({id:"".concat(R,".").concat(C),defaultMessage:C}),status:(h=E[C])!==null&&h!==void 0?h:"Default"}}return j}}}]);