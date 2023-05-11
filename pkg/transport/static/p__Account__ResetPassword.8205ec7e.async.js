(self.webpackChunkidas_ui=self.webpackChunkidas_ui||[]).push([[310],{2603:function(v,u,e){"use strict";e.d(u,{Z:function(){return M}});var O=e(28991),c=e(67294),p={icon:{tag:"svg",attrs:{viewBox:"64 64 896 896",focusable:"false"},children:[{tag:"path",attrs:{d:"M832 464h-68V240c0-70.7-57.3-128-128-128H388c-70.7 0-128 57.3-128 128v224h-68c-17.7 0-32 14.3-32 32v384c0 17.7 14.3 32 32 32h640c17.7 0 32-14.3 32-32V496c0-17.7-14.3-32-32-32zM332 240c0-30.9 25.1-56 56-56h248c30.9 0 56 25.1 56 56v224H332V240zm460 600H232V536h560v304zM484 701v53c0 4.4 3.6 8 8 8h40c4.4 0 8-3.6 8-8v-53a48.01 48.01 0 10-56 0z"}}]},name:"lock",theme:"outlined"},l=p,C=e(27029),m=function(D,j){return c.createElement(C.Z,(0,O.Z)((0,O.Z)({},D),{},{ref:j,icon:l}))};m.displayName="LockOutlined";var M=c.forwardRef(m)},57112:function(v){v.exports={container:"container___2z00C",lang:"lang___2aIoM",content:"content___2ZfKz",icon:"icon___CO56T"}},57600:function(v,u,e){"use strict";e.r(u);var O=e(57663),c=e(71577),p=e(2824),l=e(39428),C=e(34792),m=e(48086),M=e(11849),h=e(3182),D=e(34442),j=e.n(D),w=e(67294),P=e(21010),R=e(60923),T=e(45953),B=e(71390),A=e(89652),K=e(93400),U=e(89366),g=e(2603),W=e(29464),E=e(5966),Z=e(57112),o=e.n(Z),s=e(85893),x=function(y){var b=(0,P.YB)(),t=new K.f("pages.resetPassword",b),z=function(){var d=(0,h.Z)((0,l.Z)().mark(function n(i){var f,r;return(0,l.Z)().wrap(function(a){for(;;)switch(a.prev=a.next){case 0:return a.prev=0,a.next=3,(0,A.c0)((0,M.Z)({},i));case 3:if(f=a.sent,!f.success){a.next=8;break}return r=t.t("success","The password is reset successfully. Please login again with the new password."),m.default.success(r),a.abrupt("return",!0);case 8:a.next=13;break;case 10:a.prev=10,a.t0=a.catch(0),console.error(a.t0);case 13:return a.abrupt("return",!1);case 14:case"end":return a.stop()}},n,null,[[0,10]])}));return function(i){return d.apply(this,arguments)}}(),_=y.location.query,N=(0,w.useState)(!1),I=(0,p.Z)(N,2),S=I[0],L=I[1];return(0,s.jsxs)("div",{className:o().container,children:[(0,s.jsx)("div",{className:o().lang,"data-lang":!0,children:P.pD&&(0,s.jsx)(P.pD,{})}),(0,s.jsx)("div",{className:o().content,children:(0,s.jsxs)(W.U,{logo:(0,s.jsx)("img",{alt:"logo",src:window.publicPath+"logo.svg"}),title:R.Z.title,subTitle:(0,s.jsx)(s.Fragment,{children:" "}),initialValues:{username:_.username,token:_.token},submitter:{render:function(n){return(0,s.jsx)(c.Z,{loading:S,onClick:n.submit,block:!0,type:"primary",children:"\u91CD\u7F6E\u5BC6\u7801"})}},onFinish:function(){var d=(0,h.Z)((0,l.Z)().mark(function n(i){return(0,l.Z)().wrap(function(r){for(;;)switch(r.prev=r.next){case 0:return L(!0),r.next=3,z({newPassword:i.newPassword,oldPassword:i.oldPassword,userId:_.userId,token:_.token,username:_.username});case 3:if(!r.sent){r.next=5;break}P.m8.push(T.Q);case 5:L(!1);case 6:case"end":return r.stop()}},n)}));return function(n){return d.apply(this,arguments)}}(),children:[(0,s.jsx)(E.Z,{fieldProps:{value:_.username,size:"large",disabled:!0,prefix:(0,s.jsx)(U.Z,{className:o().prefixIcon})}}),(0,s.jsx)(E.Z.Password,{name:"oldPassword",fieldProps:{size:"large",prefix:(0,s.jsx)(g.Z,{className:o().prefixIcon})},hidden:_.token,placeholder:t.t("pages.login.oldPassword.placeholder","Please enter current password"),rules:[{required:!_.token,message:t.t("oldPassword.required","Please enter current password!")}]}),(0,s.jsx)(E.Z.Password,{name:"newPassword",fieldProps:{size:"large",prefix:(0,s.jsx)(g.Z,{className:o().prefixIcon})},placeholder:t.t("pages.login.password.placeholder","Please enter a new password"),rules:[{required:!0,message:t.t("password.required","Please enter a new password!")}]}),(0,s.jsx)(E.Z.Password,{name:"newPasswordConfirm",fieldProps:{size:"large",prefix:(0,s.jsx)(g.Z,{className:o().prefixIcon})},placeholder:t.t("pages.login.confirmPassword.placeholder","Confirm new password."),rules:[{required:!0,message:t.t("confirmPassword.required","Confirm new password!")},function(d){var n=d.getFieldValue;return{validator:function(f,r){return!r||n("newPassword")===r?Promise.resolve():Promise.reject(new Error("The two passwords that you entered do not match!"))}}}]})]})}),(0,s.jsx)(B.Z,{})]})};u.default=x}}]);
