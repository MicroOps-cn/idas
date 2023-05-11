(self.webpackChunkidas_ui=self.webpackChunkidas_ui||[]).push([[193],{2603:function(v,u,e){"use strict";e.d(u,{Z:function(){return M}});var f=e(28991),c=e(67294),O={icon:{tag:"svg",attrs:{viewBox:"64 64 896 896",focusable:"false"},children:[{tag:"path",attrs:{d:"M832 464h-68V240c0-70.7-57.3-128-128-128H388c-70.7 0-128 57.3-128 128v224h-68c-17.7 0-32 14.3-32 32v384c0 17.7 14.3 32 32 32h640c17.7 0 32-14.3 32-32V496c0-17.7-14.3-32-32-32zM332 240c0-30.9 25.1-56 56-56h248c30.9 0 56 25.1 56 56v224H332V240zm460 600H232V536h560v304zM484 701v53c0 4.4 3.6 8 8 8h40c4.4 0 8-3.6 8-8v-53a48.01 48.01 0 10-56 0z"}}]},name:"lock",theme:"outlined"},r=O,g=e(27029),m=function(D,A){return c.createElement(g.Z,(0,f.Z)((0,f.Z)({},D),{},{ref:A,icon:r}))};m.displayName="LockOutlined";var M=c.forwardRef(m)},45672:function(v){v.exports={container:"container___3ZjIt",lang:"lang___iAuLM",content:"content___3AaeS",icon:"icon___28cRT"}},35365:function(v,u,e){"use strict";e.r(u);var f=e(57663),c=e(71577),O=e(2824),r=e(39428),g=e(34792),m=e(48086),M=e(11849),h=e(3182),D=e(34442),A=e.n(D),T=e(67294),E=e(21010),C=e(60923),R=e(45953),B=e(71390),U=e(89652),K=e(93400),W=e(89366),L=e(2603),Z=e(29464),p=e(5966),x=e(45672),_=e.n(x),s=e(85893),y=function(w){var o=new K.f("pages.activateAccount",(0,E.YB)()),b=function(){var i=(0,h.Z)((0,r.Z)().mark(function t(d){var P,n;return(0,r.Z)().wrap(function(a){for(;;)switch(a.prev=a.next){case 0:return a.prev=0,a.next=3,(0,U.NG)((0,M.Z)({},d));case 3:if(P=a.sent,!P.success){a.next=8;break}return n=o.t("success","The account activation succeeded. Please login again with the new password."),m.default.success(n),a.abrupt("return",!0);case 8:a.next=13;break;case 10:a.prev=10,a.t0=a.catch(0),console.error(a.t0);case 13:return a.abrupt("return",!1);case 14:case"end":return a.stop()}},t,null,[[0,10]])}));return function(d){return i.apply(this,arguments)}}(),l=w.location.query,N=(0,T.useState)(!1),j=(0,O.Z)(N,2),S=j[0],I=j[1];return(0,s.jsxs)("div",{className:_().container,children:[(0,s.jsx)("div",{className:_().lang,"data-lang":!0,children:E.pD&&(0,s.jsx)(E.pD,{})}),(0,s.jsx)("div",{className:_().content,children:(0,s.jsxs)(Z.U,{logo:(0,s.jsx)("img",{alt:"logo",src:window.publicPath+"logo.svg"}),title:C.Z.title,subTitle:(0,s.jsx)(s.Fragment,{children:" "}),initialValues:{username:l.username,token:l.token},submitter:{render:function(t){return(0,s.jsx)(c.Z,{loading:S,onClick:t.submit,block:!0,type:"primary",children:o.t("button.activation","Activation")})}},onFinish:function(){var i=(0,h.Z)((0,r.Z)().mark(function t(d){return(0,r.Z)().wrap(function(n){for(;;)switch(n.prev=n.next){case 0:return I(!0),n.next=3,b({newPassword:d.newPassword,oldPassword:d.oldPassword,userId:l.userId,token:l.token,storage:l.storage});case 3:if(!n.sent){n.next=5;break}E.m8.push(R.Q);case 5:I(!1);case 6:case"end":return n.stop()}},t)}));return function(t){return i.apply(this,arguments)}}(),children:[(0,s.jsx)(p.Z,{fieldProps:{value:l.username,size:"large",disabled:!0,prefix:(0,s.jsx)(W.Z,{className:_().prefixIcon})}}),(0,s.jsx)(p.Z.Password,{name:"newPassword",fieldProps:{size:"large",prefix:(0,s.jsx)(L.Z,{className:_().prefixIcon})},placeholder:o.t("password.placeholder","Please enter a new password"),rules:[{required:!0,message:o.t("password.required","Please enter a new password!")}]}),(0,s.jsx)(p.Z.Password,{name:"newPasswordConfirm",fieldProps:{size:"large",prefix:(0,s.jsx)(L.Z,{className:_().prefixIcon})},placeholder:o.t("password.placeholder","Please enter the password again to confirm it is correct."),rules:[{required:!0,message:o.t("password.required","Please enter the confirmation password!")},function(i){var t=i.getFieldValue;return{validator:function(P,n){return!n||t("newPassword")===n?Promise.resolve():Promise.reject(new Error("The two passwords that you entered do not match!"))}}}]})]})}),(0,s.jsx)(B.Z,{})]})};u.default=y}}]);
