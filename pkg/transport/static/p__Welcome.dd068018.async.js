(self.webpackChunkidas_ui=self.webpackChunkidas_ui||[]).push([[185],{42088:function(D){D.exports={ListItem:"ListItem___osY4B"}},8668:function(D){D.exports={AppList:"AppList___1Uhdo"}},99376:function(D,p,e){"use strict";var V=e(71194),R=e(34883),w=e(402),I=e(97272),F=e(71153),T=e(60331),$=e(57663),W=e(71577),z=e(34792),O=e(48086),c=e(39428),P=e(11849),E=e(3182),L=e(2824),N=e(93224),u=e(96486),H=e.n(u),m=e(67294),A=e(2844),q=e(51042),ee=e(43471),Y=e(57119),ne=e(8978),se=e(42088),te=e.n(se),l=e(85893),ae=["intl","request","metas","onEdit","onCreate","onClick"],ie=function(v){var t=v.intl,_=v.request,re=v.metas,g=v.onEdit,G=v.onCreate,B=v.onClick,oe=(0,N.Z)(v,ae),M=(0,m.useRef)(),le=(0,m.useState)(),X=(0,L.Z)(le,2),b=X[0],_e=X[1],ue=(0,m.useState)(),k=(0,L.Z)(ue,2),de=k[0],U=k[1],J=(0,m.useCallback)(function(){var a=(0,E.Z)((0,c.Z)().mark(function s(n){return(0,c.Z)().wrap(function(d){for(;;)switch(d.prev=d.next){case 0:return d.abrupt("return",((0,u.isFunction)(_)?_:_.list)((0,P.Z)((0,P.Z)({},n),{},{keywords:b??void 0})).then(function(r){return r.total===void 0||r.pageSize===void 0||r.current===void 0?U(!0):r.total<r.current*r.pageSize?U(!1):U(!0),r}));case 1:case"end":return d.stop()}},s)}));return function(s){return a.apply(this,arguments)}}(),[_,b]),K=function(s,n){var i=arguments.length>2&&arguments[2]!==void 0?arguments[2]:{},d=i.processing,r=i.success,S=i.error,Q=s??(!(0,u.isFunction)(_)&&n&&_.patch?function(){var x=(0,E.Z)((0,c.Z)().mark(function C(h){var f;return(0,c.Z)().wrap(function(o){for(;;)switch(o.prev=o.next){case 0:return o.abrupt("return",(f=_.patch)===null||f===void 0?void 0:f.call(_,h.map(n)));case 1:case"end":return o.stop()}},C)}));return function(C){return x.apply(this,arguments)}}():null);if(!!Q)return function(){var x=(0,E.Z)((0,c.Z)().mark(function C(h){var f;return(0,c.Z)().wrap(function(o){for(;;)switch(o.prev=o.next){case 0:return f=O.default.loading(d??t.t("message.processing","Processing ...")),o.prev=1,o.next=4,Q((0,u.isString)(h)?[{id:h}]:h.map(function(me){return{id:me}}));case 4:return f(),O.default.success(r??t.t("message.operationSuccessd","Operation succeeded.")),o.abrupt("return",!0);case 9:return o.prev=9,o.t0=o.catch(1),f(),O.default.error(S??t.t("message.operationFailed","Operation failed, please try again.")),o.abrupt("return",!1);case 14:case"end":return o.stop()}},C,null,[[1,9]])}));return function(C){return x.apply(this,arguments)}}()},j=K((0,u.isFunction)(_)?void 0:_.delete,function(a){var s=a.id;return{id:s,isDelete:!0}},{processing:t.t("message.removing","Removing ..."),success:t.t("message.removeSuccessd","Remove successfully and will refresh soon"),error:t.t("message.removeFailed","Remove failed, please try again")}),y=K((0,u.isFunction)(_)?void 0:_.disable,function(a){var s=a.id;return{id:s,isDisable:!0}},{processing:t.t("message.disabling","Disabling ..."),success:t.t("message.disableSuccessd","Disabled successfully and will refresh soon"),error:t.t("message.disableFailed","Disable failed, please try again")}),Z=K((0,u.isFunction)(_)?void 0:_.enable,function(a){var s=a.id;return{id:s,isDisable:!1}},{processing:t.t("message.enabling","Enabling ..."),success:t.t("message.enableSuccessd","Enabled successfully and will refresh soon"),error:t.t("message.enableFailed","Enable failed, please try again")});return(0,m.useEffect)(function(){var a;(a=M.current)===null||a===void 0||a.reload()},[J]),(0,l.jsx)(ne.ZP,(0,P.Z)({pagination:de?{defaultPageSize:20,showSizeChanger:!0}:!1,actionRef:M,showActions:"always",grid:{gutter:16,column:3,xxl:4,xl:4,lg:4,md:3,sm:2,xs:2},toolbar:{search:!0,onSearch:function(s){var n;_e(s),(n=M.current)===null||n===void 0||n.reload()},actions:[(0,l.jsxs)(W.Z,{hidden:!G,type:"primary",onClick:G,children:[(0,l.jsx)(q.Z,{}),t.t("button.create","Create")]},"create")],settings:[{icon:(0,l.jsx)(ee.Z,{onClick:function(){var s;return(s=M.current)===null||s===void 0?void 0:s.reload()}}),key:"reload",onClick:function(){var s;return(s=M.current)===null||s===void 0?void 0:s.reload()}}]},itemCardProps:{className:te().ListItem},onItem:function(s){return{onClick:function(){B==null||B(s)}}},metas:(0,P.Z)({title:{render:function(s,n){var i;return(i=n.displayName)!==null&&i!==void 0?i:n.name}},subTitle:{render:function(s,n){return n.isDisable?(0,l.jsx)(T.Z,{color:"red",children:t.t("button.disabled","Disabled")}):(0,l.jsx)(l.Fragment,{})}},content:{render:function(s,n){var i;return(0,l.jsx)(I.Z.Paragraph,{ellipsis:{rows:2,tooltip:n.description},children:(i=n.description)!==null&&i!==void 0?i:""})}},avatar:{dataIndex:"avatar",search:!1,render:function(s,n){return(0,u.isString)(n.avatar)?(0,l.jsx)(A.ZP,{size:"default",src:n.avatar}):n.avatar}},actions:g||j||Z||y?{cardActionProps:"actions",render:function(s,n){return[g?(0,l.jsx)("a",{onClick:function(){g==null||g(n)},children:t.t("button.edit","Edit")},"edit"):null,j?(0,l.jsx)("a",{onClick:function(){R.Z.confirm({title:t.t("confirm.remove","Are you sure you want to remove this item?"),icon:(0,l.jsx)(Y.Z,{}),onOk:function(){j(n.id)},maskClosable:!0})},children:t.t("button.remove","Remove")},"remove"):null,n.isDisable!==void 0&&(n.isDisable?Z:y)?(0,l.jsx)("a",{onClick:function(){R.Z.confirm({title:n.isDisable?t.t("confirm.enable","Are you sure you want to enable this item?"):t.t("confirm.disable","Are you sure you want to disable this item?"),icon:(0,l.jsx)(Y.Z,{}),onOk:function(){var r;(r=n.isDisable?Z:y)===null||r===void 0||r(n.id)},maskClosable:!0})},children:n.isDisable?t.t("button.enable","Enable"):t.t("button.disable","Disable")},"disable_or_enable"):void 0]}}:void 0},re),request:J},oe))};p.Z=ie},53922:function(D,p,e){"use strict";e.r(p);var V=e(34792),R=e(48086),w=e(67294),I=e(21010),F=e(60923),T=e(99376),$=e(37959),W=e(89652),z=e(93400),O=e(43037),c=e(8668),P=e.n(c),E=e(85893),L=function(){var u=new z.f("pages.welcome",(0,I.YB)());return(0,E.jsx)(O.ZP,{logo:window.publicPath+"logo.svg",title:F.Z.title,layout:"top",navTheme:"dark",rightContentRender:function(){return(0,E.jsx)($.Z,{})},children:(0,E.jsx)(T.Z,{cardProps:{className:P().AppList},intl:u,request:{list:W.hX},onClick:function(m){if(m.url){var A=window.open("about:blank");A&&(A.location.href=m.url)}else R.default.warn("The URL for this application is not configured. Please contact the administrator.")}})})};p.default=L}}]);
