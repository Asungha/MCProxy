(self.webpackChunk_N_E=self.webpackChunk_N_E||[]).push([[923],{1850:function(e,t,r){Promise.resolve().then(r.bind(r,2935)),Promise.resolve().then(r.bind(r,7302)),Promise.resolve().then(r.bind(r,7165)),Promise.resolve().then(r.bind(r,8688)),Promise.resolve().then(r.bind(r,6614)),Promise.resolve().then(r.bind(r,1784)),Promise.resolve().then(r.bind(r,3697))},5884:function(e,t,r){"use strict";r.d(t,{Z:function(){return b}});var n=r(444),i=r(2110),o=r(2265),a=r(3167),s=r(6860),u=r(4e3),c=r(8836),l=r(4198),d=r(7437),f=(0,l.Z)((0,d.jsx)("path",{d:"M12 12c2.21 0 4-1.79 4-4s-1.79-4-4-4-4 1.79-4 4 1.79 4 4 4zm0 2c-2.67 0-8 1.34-8 4v2h16v-2c0-2.66-5.33-4-8-4z"}),"Person"),h=r(5749),p=r(8074);let v=["alt","children","className","component","slots","slotProps","imgProps","sizes","src","srcSet","variant"],g=(0,u.U)("MuiAvatar"),m=e=>{let{classes:t,variant:r,colorDefault:n}=e;return(0,s.Z)({root:["root",r,n&&"colorDefault"],img:["img"],fallback:["fallback"]},h.$,t)},y=(0,c.ZP)("div",{name:"MuiAvatar",slot:"Root",overridesResolver:(e,t)=>{let{ownerState:r}=e;return[t.root,t[r.variant],r.colorDefault&&t.colorDefault]}})(e=>{let{theme:t}=e;return{position:"relative",display:"flex",alignItems:"center",justifyContent:"center",flexShrink:0,width:40,height:40,fontFamily:t.typography.fontFamily,fontSize:t.typography.pxToRem(20),lineHeight:1,borderRadius:"50%",overflow:"hidden",userSelect:"none",variants:[{props:{variant:"rounded"},style:{borderRadius:(t.vars||t).shape.borderRadius}},{props:{variant:"square"},style:{borderRadius:0}},{props:{colorDefault:!0},style:(0,i.Z)({color:(t.vars||t).palette.background.default},t.vars?{backgroundColor:t.vars.palette.Avatar.defaultBg}:(0,i.Z)({backgroundColor:t.palette.grey[400]},t.applyStyles("dark",{backgroundColor:t.palette.grey[600]})))}]}}),Z=(0,c.ZP)("img",{name:"MuiAvatar",slot:"Img",overridesResolver:(e,t)=>t.img})({width:"100%",height:"100%",textAlign:"center",objectFit:"cover",color:"transparent",textIndent:1e4}),$=(0,c.ZP)(f,{name:"MuiAvatar",slot:"Fallback",overridesResolver:(e,t)=>t.fallback})({width:"75%",height:"75%"});var b=o.forwardRef(function(e,t){let r=g({props:e,name:"MuiAvatar"}),{alt:s,children:u,className:c,component:l="div",slots:f={},slotProps:h={},imgProps:b,sizes:S,src:x,srcSet:M,variant:w="circular"}=r,k=(0,n.Z)(r,v),D=null,C=function(e){let{crossOrigin:t,referrerPolicy:r,src:n,srcSet:i}=e,[a,s]=o.useState(!1);return o.useEffect(()=>{if(!n&&!i)return;s(!1);let e=!0,o=new Image;return o.onload=()=>{e&&s("loaded")},o.onerror=()=>{e&&s("error")},o.crossOrigin=t,o.referrerPolicy=r,o.src=n,i&&(o.srcset=i),()=>{e=!1}},[t,r,n,i]),a}((0,i.Z)({},b,{src:x,srcSet:M})),j=x||M,O=j&&"error"!==C,P=(0,i.Z)({},r,{colorDefault:!O,component:l,variant:w}),z=m(P),[_,R]=(0,p.Z)("img",{className:z.img,elementType:Z,externalForwardedProps:{slots:f,slotProps:{img:(0,i.Z)({},b,h.img)}},additionalProps:{alt:s,src:x,srcSet:M,sizes:S},ownerState:P});return D=O?(0,d.jsx)(_,(0,i.Z)({},R)):u||0===u?u:j&&s?s[0]:(0,d.jsx)($,{ownerState:P,className:z.fallback}),(0,d.jsx)(y,(0,i.Z)({as:l,ownerState:P,className:(0,a.Z)(z.root,c),ref:t},k,{children:D}))})},5749:function(e,t,r){"use strict";r.d(t,{$:function(){return o}});var n=r(8399),i=r(7520);function o(e){return(0,i.ZP)("MuiAvatar",e)}let a=(0,n.Z)("MuiAvatar",["root","colorDefault","circular","rounded","square","img","fallback"]);t.Z=a},2935:function(e,t,r){"use strict";r.r(t),r.d(t,{buttonClasses:function(){return i.Z},default:function(){return n.Z},getButtonUtilityClass:function(){return i.F}});var n=r(6718),i=r(4270)},7302:function(e,t,r){"use strict";r.r(t),r.d(t,{cardClasses:function(){return i.Z},default:function(){return n.Z},getCardUtilityClass:function(){return i.y}});var n=r(5092),i=r(8025)},3206:function(e,t,r){"use strict";r.d(t,{Z:function(){return b}});var n=r(444),i=r(2110),o=r(2265),a=r(3167),s=r(6860),u=r(1869),c=r(8836),l=r(3043),d=r(6441),f=r(5135),h=r(8399),p=r(7520);function v(e){return(0,p.ZP)("MuiIconButton",e)}let g=(0,h.Z)("MuiIconButton",["root","disabled","colorInherit","colorPrimary","colorSecondary","colorError","colorInfo","colorSuccess","colorWarning","edgeStart","edgeEnd","sizeSmall","sizeMedium","sizeLarge"]);var m=r(7437);let y=["edge","children","className","color","disabled","disableFocusRipple","size"],Z=e=>{let{classes:t,disabled:r,color:n,edge:i,size:o}=e,a={root:["root",r&&"disabled","default"!==n&&"color".concat((0,f.Z)(n)),i&&"edge".concat((0,f.Z)(i)),"size".concat((0,f.Z)(o))]};return(0,s.Z)(a,v,t)},$=(0,c.ZP)(d.Z,{name:"MuiIconButton",slot:"Root",overridesResolver:(e,t)=>{let{ownerState:r}=e;return[t.root,"default"!==r.color&&t["color".concat((0,f.Z)(r.color))],r.edge&&t["edge".concat((0,f.Z)(r.edge))],t["size".concat((0,f.Z)(r.size))]]}})(e=>{let{theme:t,ownerState:r}=e;return(0,i.Z)({textAlign:"center",flex:"0 0 auto",fontSize:t.typography.pxToRem(24),padding:8,borderRadius:"50%",overflow:"visible",color:(t.vars||t).palette.action.active,transition:t.transitions.create("background-color",{duration:t.transitions.duration.shortest})},!r.disableRipple&&{"&:hover":{backgroundColor:t.vars?"rgba(".concat(t.vars.palette.action.activeChannel," / ").concat(t.vars.palette.action.hoverOpacity,")"):(0,u.Fq)(t.palette.action.active,t.palette.action.hoverOpacity),"@media (hover: none)":{backgroundColor:"transparent"}}},"start"===r.edge&&{marginLeft:"small"===r.size?-3:-12},"end"===r.edge&&{marginRight:"small"===r.size?-3:-12})},e=>{var t;let{theme:r,ownerState:n}=e,o=null==(t=(r.vars||r).palette)?void 0:t[n.color];return(0,i.Z)({},"inherit"===n.color&&{color:"inherit"},"inherit"!==n.color&&"default"!==n.color&&(0,i.Z)({color:null==o?void 0:o.main},!n.disableRipple&&{"&:hover":(0,i.Z)({},o&&{backgroundColor:r.vars?"rgba(".concat(o.mainChannel," / ").concat(r.vars.palette.action.hoverOpacity,")"):(0,u.Fq)(o.main,r.palette.action.hoverOpacity)},{"@media (hover: none)":{backgroundColor:"transparent"}})}),"small"===n.size&&{padding:5,fontSize:r.typography.pxToRem(18)},"large"===n.size&&{padding:12,fontSize:r.typography.pxToRem(28)},{["&.".concat(g.disabled)]:{backgroundColor:"transparent",color:(r.vars||r).palette.action.disabled}})});var b=o.forwardRef(function(e,t){let r=(0,l.Z)({props:e,name:"MuiIconButton"}),{edge:o=!1,children:s,className:u,color:c="default",disabled:d=!1,disableFocusRipple:f=!1,size:h="medium"}=r,p=(0,n.Z)(r,y),v=(0,i.Z)({},r,{edge:o,color:c,disabled:d,disableFocusRipple:f,size:h}),g=Z(v);return(0,m.jsx)($,(0,i.Z)({className:(0,a.Z)(g.root,u),centerRipple:!0,focusRipple:!f,disabled:d,ref:t},p,{ownerState:v,children:s}))})},7165:function(e,t,r){"use strict";r.r(t),r.d(t,{default:function(){return n.Z},getInputAdornmentUtilityClass:function(){return i.w},inputAdornmentClasses:function(){return i.Z}});var n=r(9418),i=r(9755)},8688:function(e,t,r){"use strict";r.r(t),r.d(t,{default:function(){return n.Z},getOutlinedInputUtilityClass:function(){return i.e},outlinedInputClasses:function(){return i.Z}});var n=r(2961),i=r(908)},6614:function(e,t,r){"use strict";r.r(t),r.d(t,{default:function(){return n.Z},stackClasses:function(){return i}});var n=r(895),i=(0,r(8399).Z)("MuiStack",["root"])},1784:function(e,t,r){"use strict";r.r(t),r.d(t,{default:function(){return n.Z},getTypographyUtilityClass:function(){return i.f},typographyClasses:function(){return i.Z}});var n=r(8087),i=r(8935)},8074:function(e,t,r){"use strict";r.d(t,{Z:function(){return f}});var n=r(2110),i=r(444),o=r(4255),a=r(1678),s=r(2880),u=r(1028);let c=["className","elementType","ownerState","externalForwardedProps","getSlotOwnerState","internalForwardedProps"],l=["component","slots","slotProps"],d=["component"];function f(e,t){let{className:r,elementType:f,ownerState:h,externalForwardedProps:p,getSlotOwnerState:v,internalForwardedProps:g}=t,m=(0,i.Z)(t,c),{component:y,slots:Z={[e]:void 0},slotProps:$={[e]:void 0}}=p,b=(0,i.Z)(p,l),S=Z[e]||f,x=(0,a.x)($[e],h),M=(0,s.L)((0,n.Z)({className:r},m,{externalForwardedProps:"root"===e?b:void 0,externalSlotProps:x})),{props:{component:w},internalRef:k}=M,D=(0,i.Z)(M.props,d),C=(0,o.Z)(k,null==x?void 0:x.ref,t.ref),j=v?v(D):{},O=(0,n.Z)({},h,j),P="root"===e?w||y:w,z=(0,u.$)(S,(0,n.Z)({},"root"===e&&!y&&!Z[e]&&g,"root"!==e&&!Z[e]&&g,D,P&&{as:P},{ref:C}),O);return Object.keys(j).forEach(e=>{delete z[e]}),[S,z]}},4e3:function(e,t,r){"use strict";r.d(t,{U:function(){return i}});var n=r(3043);function i(e){return n.Z}},9539:function(e){var t,r,n,i,o,a,s,u,c,l,d,f,h,p,v,g,m,y,Z,$,b,S;e.exports=(t="millisecond",r="second",n="minute",i="hour",o="week",a="month",s="quarter",u="year",c="date",l="Invalid Date",d=/^(\d{4})[-/]?(\d{1,2})?[-/]?(\d{0,2})[Tt\s]*(\d{1,2})?:?(\d{1,2})?:?(\d{1,2})?[.:]?(\d+)?$/,f=/\[([^\]]+)]|Y{1,4}|M{1,4}|D{1,2}|d{1,4}|H{1,2}|h{1,2}|a|A|m{1,2}|s{1,2}|Z{1,2}|SSS/g,h=function(e,t,r){var n=String(e);return!n||n.length>=t?e:""+Array(t+1-n.length).join(r)+e},(v={})[p="en"]={name:"en",weekdays:"Sunday_Monday_Tuesday_Wednesday_Thursday_Friday_Saturday".split("_"),months:"January_February_March_April_May_June_July_August_September_October_November_December".split("_"),ordinal:function(e){var t=["th","st","nd","rd"],r=e%100;return"["+e+(t[(r-20)%10]||t[r]||"th")+"]"}},g="$isDayjsObject",m=function(e){return e instanceof b||!(!e||!e[g])},y=function e(t,r,n){var i;if(!t)return p;if("string"==typeof t){var o=t.toLowerCase();v[o]&&(i=o),r&&(v[o]=r,i=o);var a=t.split("-");if(!i&&a.length>1)return e(a[0])}else{var s=t.name;v[s]=t,i=s}return!n&&i&&(p=i),i||!n&&p},Z=function(e,t){if(m(e))return e.clone();var r="object"==typeof t?t:{};return r.date=e,r.args=arguments,new b(r)},($={s:h,z:function(e){var t=-e.utcOffset(),r=Math.abs(t);return(t<=0?"+":"-")+h(Math.floor(r/60),2,"0")+":"+h(r%60,2,"0")},m:function e(t,r){if(t.date()<r.date())return-e(r,t);var n=12*(r.year()-t.year())+(r.month()-t.month()),i=t.clone().add(n,a),o=r-i<0,s=t.clone().add(n+(o?-1:1),a);return+(-(n+(r-i)/(o?i-s:s-i))||0)},a:function(e){return e<0?Math.ceil(e)||0:Math.floor(e)},p:function(e){return({M:a,y:u,w:o,d:"day",D:c,h:i,m:n,s:r,ms:t,Q:s})[e]||String(e||"").toLowerCase().replace(/s$/,"")},u:function(e){return void 0===e}}).l=y,$.i=m,$.w=function(e,t){return Z(e,{locale:t.$L,utc:t.$u,x:t.$x,$offset:t.$offset})},S=(b=function(){function e(e){this.$L=y(e.locale,null,!0),this.parse(e),this.$x=this.$x||e.x||{},this[g]=!0}var h=e.prototype;return h.parse=function(e){this.$d=function(e){var t=e.date,r=e.utc;if(null===t)return new Date(NaN);if($.u(t))return new Date;if(t instanceof Date)return new Date(t);if("string"==typeof t&&!/Z$/i.test(t)){var n=t.match(d);if(n){var i=n[2]-1||0,o=(n[7]||"0").substring(0,3);return r?new Date(Date.UTC(n[1],i,n[3]||1,n[4]||0,n[5]||0,n[6]||0,o)):new Date(n[1],i,n[3]||1,n[4]||0,n[5]||0,n[6]||0,o)}}return new Date(t)}(e),this.init()},h.init=function(){var e=this.$d;this.$y=e.getFullYear(),this.$M=e.getMonth(),this.$D=e.getDate(),this.$W=e.getDay(),this.$H=e.getHours(),this.$m=e.getMinutes(),this.$s=e.getSeconds(),this.$ms=e.getMilliseconds()},h.$utils=function(){return $},h.isValid=function(){return this.$d.toString()!==l},h.isSame=function(e,t){var r=Z(e);return this.startOf(t)<=r&&r<=this.endOf(t)},h.isAfter=function(e,t){return Z(e)<this.startOf(t)},h.isBefore=function(e,t){return this.endOf(t)<Z(e)},h.$g=function(e,t,r){return $.u(e)?this[t]:this.set(r,e)},h.unix=function(){return Math.floor(this.valueOf()/1e3)},h.valueOf=function(){return this.$d.getTime()},h.startOf=function(e,t){var s=this,l=!!$.u(t)||t,d=$.p(e),f=function(e,t){var r=$.w(s.$u?Date.UTC(s.$y,t,e):new Date(s.$y,t,e),s);return l?r:r.endOf("day")},h=function(e,t){return $.w(s.toDate()[e].apply(s.toDate("s"),(l?[0,0,0,0]:[23,59,59,999]).slice(t)),s)},p=this.$W,v=this.$M,g=this.$D,m="set"+(this.$u?"UTC":"");switch(d){case u:return l?f(1,0):f(31,11);case a:return l?f(1,v):f(0,v+1);case o:var y=this.$locale().weekStart||0,Z=(p<y?p+7:p)-y;return f(l?g-Z:g+(6-Z),v);case"day":case c:return h(m+"Hours",0);case i:return h(m+"Minutes",1);case n:return h(m+"Seconds",2);case r:return h(m+"Milliseconds",3);default:return this.clone()}},h.endOf=function(e){return this.startOf(e,!1)},h.$set=function(e,o){var s,l=$.p(e),d="set"+(this.$u?"UTC":""),f=((s={}).day=d+"Date",s[c]=d+"Date",s[a]=d+"Month",s[u]=d+"FullYear",s[i]=d+"Hours",s[n]=d+"Minutes",s[r]=d+"Seconds",s[t]=d+"Milliseconds",s)[l],h="day"===l?this.$D+(o-this.$W):o;if(l===a||l===u){var p=this.clone().set(c,1);p.$d[f](h),p.init(),this.$d=p.set(c,Math.min(this.$D,p.daysInMonth())).$d}else f&&this.$d[f](h);return this.init(),this},h.set=function(e,t){return this.clone().$set(e,t)},h.get=function(e){return this[$.p(e)]()},h.add=function(e,t){var s,c=this;e=Number(e);var l=$.p(t),d=function(t){var r=Z(c);return $.w(r.date(r.date()+Math.round(t*e)),c)};if(l===a)return this.set(a,this.$M+e);if(l===u)return this.set(u,this.$y+e);if("day"===l)return d(1);if(l===o)return d(7);var f=((s={})[n]=6e4,s[i]=36e5,s[r]=1e3,s)[l]||1,h=this.$d.getTime()+e*f;return $.w(h,this)},h.subtract=function(e,t){return this.add(-1*e,t)},h.format=function(e){var t=this,r=this.$locale();if(!this.isValid())return r.invalidDate||l;var n=e||"YYYY-MM-DDTHH:mm:ssZ",i=$.z(this),o=this.$H,a=this.$m,s=this.$M,u=r.weekdays,c=r.months,d=r.meridiem,h=function(e,r,i,o){return e&&(e[r]||e(t,n))||i[r].slice(0,o)},p=function(e){return $.s(o%12||12,e,"0")},v=d||function(e,t,r){var n=e<12?"AM":"PM";return r?n.toLowerCase():n};return n.replace(f,function(e,n){return n||function(e){switch(e){case"YY":return String(t.$y).slice(-2);case"YYYY":return $.s(t.$y,4,"0");case"M":return s+1;case"MM":return $.s(s+1,2,"0");case"MMM":return h(r.monthsShort,s,c,3);case"MMMM":return h(c,s);case"D":return t.$D;case"DD":return $.s(t.$D,2,"0");case"d":return String(t.$W);case"dd":return h(r.weekdaysMin,t.$W,u,2);case"ddd":return h(r.weekdaysShort,t.$W,u,3);case"dddd":return u[t.$W];case"H":return String(o);case"HH":return $.s(o,2,"0");case"h":return p(1);case"hh":return p(2);case"a":return v(o,a,!0);case"A":return v(o,a,!1);case"m":return String(a);case"mm":return $.s(a,2,"0");case"s":return String(t.$s);case"ss":return $.s(t.$s,2,"0");case"SSS":return $.s(t.$ms,3,"0");case"Z":return i}return null}(e)||i.replace(":","")})},h.utcOffset=function(){return-(15*Math.round(this.$d.getTimezoneOffset()/15))},h.diff=function(e,t,c){var l,d=this,f=$.p(t),h=Z(e),p=(h.utcOffset()-this.utcOffset())*6e4,v=this-h,g=function(){return $.m(d,h)};switch(f){case u:l=g()/12;break;case a:l=g();break;case s:l=g()/3;break;case o:l=(v-p)/6048e5;break;case"day":l=(v-p)/864e5;break;case i:l=v/36e5;break;case n:l=v/6e4;break;case r:l=v/1e3;break;default:l=v}return c?l:$.a(l)},h.daysInMonth=function(){return this.endOf(a).$D},h.$locale=function(){return v[this.$L]},h.locale=function(e,t){if(!e)return this.$L;var r=this.clone(),n=y(e,t,!0);return n&&(r.$L=n),r},h.clone=function(){return $.w(this.$d,this)},h.toDate=function(){return new Date(this.valueOf())},h.toJSON=function(){return this.isValid()?this.toISOString():null},h.toISOString=function(){return this.$d.toISOString()},h.toString=function(){return this.$d.toUTCString()},e}()).prototype,Z.prototype=S,[["$ms",t],["$s",r],["$m",n],["$H",i],["$W","day"],["$M",a],["$y",u],["$D",c]].forEach(function(e){S[e[1]]=function(t){return this.$g(t,e[0],e[1])}}),Z.extend=function(e,t){return e.$i||(e(t,b,Z),e.$i=!0),Z},Z.locale=y,Z.isDayjs=m,Z.unix=function(e){return Z(1e3*e)},Z.en=v[p],Z.Ls=v,Z.p={},Z)},3697:function(e,t,r){"use strict";r.r(t),r.d(t,{CustomersTable:function(){return S}});var n=r(7437),i=r(2265),o=r(5884),a=r(304),s=r(5092),u=r(8563),c=r(8129),l=r(895),d=r(6520),f=r(360),h=r(3407),p=r(412),v=r(8840),g=r(5167),m=r(8087),y=r(9539),Z=r.n(y),$=r(4355);function b(){}function S(e){var t,r;let{count:y=0,rows:S=[],page:x=0,rowsPerPage:M=0}=e,w=i.useMemo(()=>S.map(e=>e.id),[S]),{selectAll:k,deselectAll:D,selectOne:C,deselectOne:j,selected:O}=(0,$.c)(w),P=(null!==(t=null==O?void 0:O.size)&&void 0!==t?t:0)>0&&(null!==(r=null==O?void 0:O.size)&&void 0!==r?r:0)<S.length,z=S.length>0&&(null==O?void 0:O.size)===S.length;return(0,n.jsxs)(s.Z,{children:[(0,n.jsx)(a.Z,{sx:{overflowX:"auto"},children:(0,n.jsxs)(d.Z,{sx:{minWidth:"800px"},children:[(0,n.jsx)(p.Z,{children:(0,n.jsxs)(g.Z,{children:[(0,n.jsx)(h.Z,{padding:"checkbox",children:(0,n.jsx)(u.Z,{checked:z,indeterminate:P,onChange:e=>{e.target.checked?k():D()}})}),(0,n.jsx)(h.Z,{children:"Name"}),(0,n.jsx)(h.Z,{children:"Email"}),(0,n.jsx)(h.Z,{children:"Location"}),(0,n.jsx)(h.Z,{children:"Phone"}),(0,n.jsx)(h.Z,{children:"Signed Up"})]})}),(0,n.jsx)(f.Z,{children:S.map(e=>{let t=null==O?void 0:O.has(e.id);return(0,n.jsxs)(g.Z,{hover:!0,selected:t,children:[(0,n.jsx)(h.Z,{padding:"checkbox",children:(0,n.jsx)(u.Z,{checked:t,onChange:t=>{t.target.checked?C(e.id):j(e.id)}})}),(0,n.jsx)(h.Z,{children:(0,n.jsxs)(l.Z,{sx:{alignItems:"center"},direction:"row",spacing:2,children:[(0,n.jsx)(o.Z,{src:e.avatar}),(0,n.jsx)(m.Z,{variant:"subtitle2",children:e.name})]})}),(0,n.jsx)(h.Z,{children:e.email}),(0,n.jsxs)(h.Z,{children:[e.address.city,", ",e.address.state,", ",e.address.country]}),(0,n.jsx)(h.Z,{children:e.phone}),(0,n.jsx)(h.Z,{children:Z()(e.createdAt).format("MMM D, YYYY")})]},e.id)})})]})}),(0,n.jsx)(c.Z,{}),(0,n.jsx)(v.Z,{component:"div",count:y,onPageChange:b,onRowsPerPageChange:b,page:x,rowsPerPage:M,rowsPerPageOptions:[5,10,25]})]})}},4355:function(e,t,r){"use strict";r.d(t,{c:function(){return i}});var n=r(2265);function i(){let e=arguments.length>0&&void 0!==arguments[0]?arguments[0]:[],[t,r]=n.useState(new Set);n.useEffect(()=>{r(new Set)},[e]);let i=n.useCallback(()=>{r(new Set)},[]),o=n.useCallback(e=>{r(t=>{let r=new Set(t);return r.delete(e),r})},[]),a=n.useCallback(()=>{r(new Set(e))},[e]),s=n.useCallback(e=>{r(t=>{let r=new Set(t);return r.add(e),r})},[]),u=t.size>0,c=t.size===e.length;return{deselectAll:i,deselectOne:o,selectAll:a,selectOne:s,selected:t,selectedAny:u,selectedAll:c}}}},function(e){e.O(0,[626,17,686,961,956,563,669,106,469,971,69,744],function(){return e(e.s=1850)}),_N_E=e.O()}]);