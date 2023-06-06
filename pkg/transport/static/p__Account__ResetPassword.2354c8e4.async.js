(self.webpackChunkidas_ui=self.webpackChunkidas_ui||[]).push([[310],{2603:function(g,m,e){"use strict";e.d(m,{Z:function(){return D}});var O=e(28991),P=e(67294),M={icon:{tag:"svg",attrs:{viewBox:"64 64 896 896",focusable:"false"},children:[{tag:"path",attrs:{d:"M832 464h-68V240c0-70.7-57.3-128-128-128H388c-70.7 0-128 57.3-128 128v224h-68c-17.7 0-32 14.3-32 32v384c0 17.7 14.3 32 32 32h640c17.7 0 32-14.3 32-32V496c0-17.7-14.3-32-32-32zM332 240c0-30.9 25.1-56 56-56h248c30.9 0 56 25.1 56 56v224H332V240zm460 600H232V536h560v304zM484 701v53c0 4.4 3.6 8 8 8h40c4.4 0 8-3.6 8-8v-53a48.01 48.01 0 10-56 0z"}}]},name:"lock",theme:"outlined"},_=M,T=e(27029),E=function(h,w){return P.createElement(T.Z,(0,O.Z)((0,O.Z)({},h),{},{ref:w,icon:_}))};E.displayName="LockOutlined";var D=P.forwardRef(E)},57112:function(g){g.exports={container:"container___2z00C",lang:"lang___2aIoM",content:"content___2ZfKz"}},57600:function(g,m,e){"use strict";e.r(m);var O=e(57663),P=e(71577),M=e(2824),_=e(39428),T=e(34792),E=e(48086),D=e(11849),p=e(3182),h=e(34442),w=e.n(h),K=e(67294),c=e(21010),U=e(60923),W=e(45953),Z=e(29791),x=e(89652),y=e(93400),b=e(72709),S=e(89366),C=e(2603),z=e(29464),v=e(5966),N=e(57112),o=e.n(N),s=e(85893),V=function(F){var I,L,j,$=(0,c.YB)(),t=new y.f("pages.resetPassword",$),H=function(){var d=(0,p.Z)((0,_.Z)().mark(function r(u){var f,n;return(0,_.Z)().wrap(function(a){for(;;)switch(a.prev=a.next){case 0:return a.prev=0,a.next=3,(0,x.c0)((0,D.Z)({},u));case 3:if(f=a.sent,!f.success){a.next=8;break}return n=t.t("success","The password is reset successfully. Please login again with the new password."),E.default.success(n),a.abrupt("return",!0);case 8:a.next=13;break;case 10:a.prev=10,a.t0=a.catch(0),console.error(a.t0);case 13:return a.abrupt("return",!1);case 14:case"end":return a.stop()}},r,null,[[0,10]])}));return function(u){return d.apply(this,arguments)}}(),l=F.location.query,G=(0,K.useState)(!1),A=(0,M.Z)(G,2),Q=A[0],B=A[1],Y=(0,c.tT)("@@initialState"),R=Y.initialState,i=(I=R==null?void 0:R.globalConfig)!==null&&I!==void 0?I:null;return(0,s.jsxs)("div",{className:o().container,children:[(0,s.jsx)("div",{className:o().lang,"data-lang":!0,children:c.pD&&(0,s.jsx)(c.pD,{})}),(0,s.jsx)("div",{className:o().content,children:(0,s.jsxs)(z.U,{logo:(L=i==null?void 0:i.logo)!==null&&L!==void 0?L:(0,b.Ak)("logo.svg"),title:(j=i==null?void 0:i.title)!==null&&j!==void 0?j:U.Z.title,subTitle:(0,s.jsx)(s.Fragment,{children:" "}),initialValues:{username:l.username,token:l.token},submitter:{render:function(r){return(0,s.jsx)(P.Z,{loading:Q,onClick:r.submit,block:!0,type:"primary",children:"\u91CD\u7F6E\u5BC6\u7801"})}},onFinish:function(){var d=(0,p.Z)((0,_.Z)().mark(function r(u){return(0,_.Z)().wrap(function(n){for(;;)switch(n.prev=n.next){case 0:return B(!0),n.next=3,H({newPassword:u.newPassword,oldPassword:u.oldPassword,userId:l.userId,token:l.token,username:l.username});case 3:if(!n.sent){n.next=5;break}c.m8.push(W.Q);case 5:B(!1);case 6:case"end":return n.stop()}},r)}));return function(r){return d.apply(this,arguments)}}(),children:[(0,s.jsx)(v.Z,{fieldProps:{value:l.username,size:"large",disabled:!0,prefix:(0,s.jsx)(S.Z,{className:o().prefixIcon})}}),(0,s.jsx)(v.Z.Password,{name:"oldPassword",fieldProps:{size:"large",prefix:(0,s.jsx)(C.Z,{className:o().prefixIcon})},hidden:l.token,placeholder:t.t("pages.login.oldPassword.placeholder","Please enter current password"),rules:[{required:!l.token,message:t.t("oldPassword.required","Please enter current password!")}]}),(0,s.jsx)(v.Z.Password,{name:"newPassword",fieldProps:{size:"large",prefix:(0,s.jsx)(C.Z,{className:o().prefixIcon})},placeholder:t.t("pages.login.password.placeholder","Please enter a new password"),rules:[{required:!0,message:t.t("password.required","Please enter a new password!")}]}),(0,s.jsx)(v.Z.Password,{name:"newPasswordConfirm",fieldProps:{size:"large",prefix:(0,s.jsx)(C.Z,{className:o().prefixIcon})},placeholder:t.t("pages.login.confirmPassword.placeholder","Confirm new password."),rules:[{required:!0,message:t.t("confirmPassword.required","Confirm new password!")},function(d){var r=d.getFieldValue;return{validator:function(f,n){return!n||r("newPassword")===n?Promise.resolve():Promise.reject(new Error("The two passwords that you entered do not match!"))}}}]})]})}),(0,s.jsx)(Z.Z,{})]})};m.default=V}}]);
