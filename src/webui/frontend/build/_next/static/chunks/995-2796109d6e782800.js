"use strict";(self.webpackChunk_N_E=self.webpackChunk_N_E||[]).push([[995],{41028:function(e,n,o){o.d(n,{$:function(){return i}});var t=o(14749),l=o(9413);function i(e,n,o){return void 0===e||(0,l.X)(e)?n:(0,t.Z)({},n,{ownerState:(0,t.Z)({},n.ownerState,o)})}},68508:function(e,n,o){o.d(n,{_:function(){return t}});function t(e,n=[]){if(void 0===e)return{};let o={};return Object.keys(e).filter(o=>o.match(/^on[A-Z]/)&&"function"==typeof e[o]&&!n.includes(o)).forEach(n=>{o[n]=e[n]}),o}},9413:function(e,n,o){o.d(n,{X:function(){return t}});function t(e){return"string"==typeof e}},72880:function(e,n,o){o.d(n,{L:function(){return a}});var t=o(14749),l=o(75504),i=o(68508);function r(e){if(void 0===e)return{};let n={};return Object.keys(e).filter(n=>!(n.match(/^on[A-Z]/)&&"function"==typeof e[n])).forEach(o=>{n[o]=e[o]}),n}function a(e){let{getSlotProps:n,additionalProps:o,externalSlotProps:a,externalForwardedProps:u,className:c}=e;if(!n){let e=(0,l.Z)(null==o?void 0:o.className,c,null==u?void 0:u.className,null==a?void 0:a.className),n=(0,t.Z)({},null==o?void 0:o.style,null==u?void 0:u.style,null==a?void 0:a.style),i=(0,t.Z)({},o,u,a);return e.length>0&&(i.className=e),Object.keys(n).length>0&&(i.style=n),{props:i,internalRef:void 0}}let v=(0,i._)((0,t.Z)({},u,a)),s=r(a),d=r(u),f=n(v),m=(0,l.Z)(null==f?void 0:f.className,null==o?void 0:o.className,c,null==u?void 0:u.className,null==a?void 0:a.className),p=(0,t.Z)({},null==f?void 0:f.style,null==o?void 0:o.style,null==u?void 0:u.style,null==a?void 0:a.style),Z=(0,t.Z)({},f,o,d,s);return m.length>0&&(Z.className=m),Object.keys(p).length>0&&(Z.style=p),{props:Z,internalRef:f.ref}}},21678:function(e,n,o){o.d(n,{x:function(){return t}});function t(e,n,o){return"function"==typeof e?e(n,o):e}},37630:function(e,n,o){var t=o(70444),l=o(14749),i=o(2265),r=o(75504),a=o(76860),u=o(41869),c=o(58836),v=o(39497),s=o(93043),d=o(13499),f=o(57437);let m=["className","component","elevation","square","variant"],p=e=>{let{square:n,elevation:o,variant:t,classes:l}=e;return(0,a.Z)({root:["root",t,!n&&"rounded","elevation"===t&&"elevation".concat(o)]},d.J,l)},Z=(0,c.ZP)("div",{name:"MuiPaper",slot:"Root",overridesResolver:(e,n)=>{let{ownerState:o}=e;return[n.root,n[o.variant],!o.square&&n.rounded,"elevation"===o.variant&&n["elevation".concat(o.elevation)]]}})(e=>{var n;let{theme:o,ownerState:t}=e;return(0,l.Z)({backgroundColor:(o.vars||o).palette.background.paper,color:(o.vars||o).palette.text.primary,transition:o.transitions.create("box-shadow")},!t.square&&{borderRadius:o.shape.borderRadius},"outlined"===t.variant&&{border:"1px solid ".concat((o.vars||o).palette.divider)},"elevation"===t.variant&&(0,l.Z)({boxShadow:(o.vars||o).shadows[t.elevation]},!o.vars&&"dark"===o.palette.mode&&{backgroundImage:"linear-gradient(".concat((0,u.Fq)("#fff",(0,v.Z)(t.elevation)),", ").concat((0,u.Fq)("#fff",(0,v.Z)(t.elevation)),")")},o.vars&&{backgroundImage:null==(n=o.vars.overlays)?void 0:n[t.elevation]}))}),h=i.forwardRef(function(e,n){let o=(0,s.Z)({props:e,name:"MuiPaper"}),{className:i,component:a="div",elevation:u=1,square:c=!1,variant:v="elevation"}=o,d=(0,t.Z)(o,m),h=(0,l.Z)({},o,{component:a,elevation:u,square:c,variant:v}),g=p(h);return(0,f.jsx)(Z,(0,l.Z)({as:a,ownerState:h,className:(0,r.Z)(g.root,i),ref:n},d))});n.Z=h},13499:function(e,n,o){o.d(n,{J:function(){return i}});var t=o(28399),l=o(37520);function i(e){return(0,l.ZP)("MuiPaper",e)}let r=(0,t.Z)("MuiPaper",["root","rounded","outlined","elevation","elevation0","elevation1","elevation2","elevation3","elevation4","elevation5","elevation6","elevation7","elevation8","elevation9","elevation10","elevation11","elevation12","elevation13","elevation14","elevation15","elevation16","elevation17","elevation18","elevation19","elevation20","elevation21","elevation22","elevation23","elevation24"]);n.Z=r},39497:function(e,n){n.Z=e=>((e<1?5.11916*e**2:4.5*Math.log(e+1)+2)/100).toFixed(2)},34198:function(e,n,o){o.d(n,{Z:function(){return y}});var t=o(14749),l=o(2265),i=o(70444),r=o(75504),a=o(76860),u=o(95135),c=o(93043),v=o(58836),s=o(28399),d=o(37520);function f(e){return(0,d.ZP)("MuiSvgIcon",e)}(0,s.Z)("MuiSvgIcon",["root","colorPrimary","colorSecondary","colorAction","colorError","colorDisabled","fontSizeInherit","fontSizeSmall","fontSizeMedium","fontSizeLarge"]);var m=o(57437);let p=["children","className","color","component","fontSize","htmlColor","inheritViewBox","titleAccess","viewBox"],Z=e=>{let{color:n,fontSize:o,classes:t}=e,l={root:["root","inherit"!==n&&"color".concat((0,u.Z)(n)),"fontSize".concat((0,u.Z)(o))]};return(0,a.Z)(l,f,t)},h=(0,v.ZP)("svg",{name:"MuiSvgIcon",slot:"Root",overridesResolver:(e,n)=>{let{ownerState:o}=e;return[n.root,"inherit"!==o.color&&n["color".concat((0,u.Z)(o.color))],n["fontSize".concat((0,u.Z)(o.fontSize))]]}})(e=>{var n,o,t,l,i,r,a,u,c,v,s,d,f;let{theme:m,ownerState:p}=e;return{userSelect:"none",width:"1em",height:"1em",display:"inline-block",fill:p.hasSvgAsChild?void 0:"currentColor",flexShrink:0,transition:null==(n=m.transitions)||null==(o=n.create)?void 0:o.call(n,"fill",{duration:null==(t=m.transitions)||null==(t=t.duration)?void 0:t.shorter}),fontSize:({inherit:"inherit",small:(null==(l=m.typography)||null==(i=l.pxToRem)?void 0:i.call(l,20))||"1.25rem",medium:(null==(r=m.typography)||null==(a=r.pxToRem)?void 0:a.call(r,24))||"1.5rem",large:(null==(u=m.typography)||null==(c=u.pxToRem)?void 0:c.call(u,35))||"2.1875rem"})[p.fontSize],color:null!=(v=null==(s=(m.vars||m).palette)||null==(s=s[p.color])?void 0:s.main)?v:({action:null==(d=(m.vars||m).palette)||null==(d=d.action)?void 0:d.active,disabled:null==(f=(m.vars||m).palette)||null==(f=f.action)?void 0:f.disabled,inherit:void 0})[p.color]}}),g=l.forwardRef(function(e,n){let o=(0,c.Z)({props:e,name:"MuiSvgIcon"}),{children:a,className:u,color:v="inherit",component:s="svg",fontSize:d="medium",htmlColor:f,inheritViewBox:g=!1,titleAccess:y,viewBox:S="0 0 24 24"}=o,b=(0,i.Z)(o,p),N=l.isValidElement(a)&&"svg"===a.type,x=(0,t.Z)({},o,{color:v,component:s,fontSize:d,instanceFontSize:e.fontSize,inheritViewBox:g,viewBox:S,hasSvgAsChild:N}),w={};g||(w.viewBox=S);let k=Z(x);return(0,m.jsxs)(h,(0,t.Z)({as:s,className:(0,r.Z)(k.root,u),focusable:"false",color:f,"aria-hidden":!y||void 0,role:y?"img":void 0,ref:n},w,b,N&&a.props,{ownerState:x,children:[N?a.props.children:a,y?(0,m.jsx)("title",{children:y}):null]}))});function y(e,n){function o(o,l){return(0,m.jsx)(g,(0,t.Z)({"data-testid":"".concat(n,"Icon"),ref:l},o,{children:e}))}return o.muiName=g.muiName,l.memo(l.forwardRef(o))}g.muiName="SvgIcon"},65735:function(e,n,o){var t=o(66871);n.Z=t.Z}}]);