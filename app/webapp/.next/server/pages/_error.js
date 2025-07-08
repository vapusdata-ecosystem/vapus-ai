"use strict";(()=>{var e={};e.id=2731,e.ids=[636,2731,3220],e.modules={1523:(e,t,r)=>{e.exports=r(63885).vendored.contexts.HeadManagerContext},3147:(e,t)=>{function r(e){if("function"!=typeof WeakMap)return null;var t=new WeakMap,n=new WeakMap;return(r=function(e){return e?n:t})(e)}t._=function(e,t){if(!t&&e&&e.__esModule)return e;if(null===e||"object"!=typeof e&&"function"!=typeof e)return{default:e};var n=r(t);if(n&&n.has(e))return n.get(e);var o={__proto__:null},a=Object.defineProperty&&Object.getOwnPropertyDescriptor;for(var i in e)if("default"!==i&&Object.prototype.hasOwnProperty.call(e,i)){var l=a?Object.getOwnPropertyDescriptor(e,i):null;l&&(l.get||l.set)?Object.defineProperty(o,i,l):o[i]=e[i]}return o.default=e,n&&n.set(e,o),o}},8732:e=>{e.exports=require("react/jsx-runtime")},17341:(e,t,r)=>{Object.defineProperty(t,"__esModule",{value:!0}),!function(e,t){for(var r in t)Object.defineProperty(e,r,{enumerable:!0,get:t[r]})}(t,{default:function(){return g},defaultHead:function(){return c}});let n=r(87020),o=r(3147),a=r(8732),i=o._(r(82015)),l=n._(r(95996)),d=r(57043),s=r(1523),u=r(28725);function c(e){void 0===e&&(e=!1);let t=[(0,a.jsx)("meta",{charSet:"utf-8"},"charset")];return e||t.push((0,a.jsx)("meta",{name:"viewport",content:"width=device-width"},"viewport")),t}function p(e,t){return"string"==typeof t||"number"==typeof t?e:t.type===i.default.Fragment?e.concat(i.default.Children.toArray(t.props.children).reduce((e,t)=>"string"==typeof t||"number"==typeof t?e:e.concat(t),[])):e.concat(t)}r(83901);let f=["name","httpEquiv","charSet","itemProp"];function b(e,t){let{inAmpMode:r}=t;return e.reduce(p,[]).reverse().concat(c(r).reverse()).filter(function(){let e=new Set,t=new Set,r=new Set,n={};return o=>{let a=!0,i=!1;if(o.key&&"number"!=typeof o.key&&o.key.indexOf("$")>0){i=!0;let t=o.key.slice(o.key.indexOf("$")+1);e.has(t)?a=!1:e.add(t)}switch(o.type){case"title":case"base":t.has(o.type)?a=!1:t.add(o.type);break;case"meta":for(let e=0,t=f.length;e<t;e++){let t=f[e];if(o.props.hasOwnProperty(t)){if("charSet"===t)r.has(t)?a=!1:r.add(t);else{let e=o.props[t],r=n[t]||new Set;("name"!==t||!i)&&r.has(e)?a=!1:(r.add(e),n[t]=r)}}}}return a}}()).reverse().map((e,t)=>{let n=e.key||t;if(process.env.__NEXT_OPTIMIZE_FONTS&&!r&&"link"===e.type&&e.props.href&&["https://fonts.googleapis.com/css","https://use.typekit.net/"].some(t=>e.props.href.startsWith(t))){let t={...e.props||{}};return t["data-href"]=t.href,t.href=void 0,t["data-optimized-fonts"]=!0,i.default.cloneElement(e,t)}return i.default.cloneElement(e,{key:n})})}let g=function(e){let{children:t}=e,r=(0,i.useContext)(d.AmpStateContext),n=(0,i.useContext)(s.HeadManagerContext);return(0,a.jsx)(l.default,{reduceComponentsToState:b,headManager:n,inAmpMode:(0,u.isInAmpMode)(r),children:t})};("function"==typeof t.default||"object"==typeof t.default&&null!==t.default)&&void 0===t.default.__esModule&&(Object.defineProperty(t.default,"__esModule",{value:!0}),Object.assign(t.default,t),e.exports=t.default)},26931:(e,t,r)=>{r.a(e,async(e,n)=>{try{r.r(t),r.d(t,{config:()=>h,default:()=>p,getServerSideProps:()=>g,getStaticPaths:()=>b,getStaticProps:()=>f,reportWebVitals:()=>m,routeModule:()=>j,unstable_getServerProps:()=>_,unstable_getServerSideProps:()=>P,unstable_getStaticParams:()=>x,unstable_getStaticPaths:()=>v,unstable_getStaticProps:()=>y});var o=r(63885),a=r(80237),i=r(81413),l=r(58548),d=r.n(l),s=r(36411),u=r(66631),c=e([s]);s=(c.then?(await c)():c)[0];let p=(0,i.M)(u,"default"),f=(0,i.M)(u,"getStaticProps"),b=(0,i.M)(u,"getStaticPaths"),g=(0,i.M)(u,"getServerSideProps"),h=(0,i.M)(u,"config"),m=(0,i.M)(u,"reportWebVitals"),y=(0,i.M)(u,"unstable_getStaticProps"),v=(0,i.M)(u,"unstable_getStaticPaths"),x=(0,i.M)(u,"unstable_getStaticParams"),_=(0,i.M)(u,"unstable_getServerProps"),P=(0,i.M)(u,"unstable_getServerSideProps"),j=new o.PagesRouteModule({definition:{kind:a.A.PAGES,page:"/_error",pathname:"/_error",bundlePath:"",filename:""},components:{App:s.default,Document:d()},userland:u});n()}catch(e){n(e)}})},28725:(e,t)=>{function r(e){let{ampFirst:t=!1,hybrid:r=!1,hasQuery:n=!1}=void 0===e?{}:e;return t||r&&n}Object.defineProperty(t,"__esModule",{value:!0}),Object.defineProperty(t,"isInAmpMode",{enumerable:!0,get:function(){return r}})},33873:e=>{e.exports=require("path")},35124:(e,t)=>{Object.defineProperty(t,"__esModule",{value:!0}),!function(e,t){for(var r in t)Object.defineProperty(e,r,{enumerable:!0,get:t[r]})}(t,{NEXT_REQUEST_META:function(){return r},addRequestMeta:function(){return a},getRequestMeta:function(){return n},removeRequestMeta:function(){return i},setRequestMeta:function(){return o}});let r=Symbol.for("NextInternalRequestMeta");function n(e,t){let n=e[r]||{};return"string"==typeof t?n[t]:n}function o(e,t){return e[r]=t,t}function a(e,t,r){let a=n(e);return a[t]=r,o(e,a)}function i(e,t){let r=n(e);return delete r[t],o(e,r)}},36411:(e,t,r)=>{r.a(e,async(e,n)=>{try{r.r(t),r.d(t,{default:()=>l});var o=r(8732);r(82015),r(64148);var a=r(77884),i=e([a]);function l({Component:e,pageProps:t}){return(0,o.jsx)(a.Auth0Provider,{children:(0,o.jsx)(e,{...t})})}a=(i.then?(await i)():i)[0],n()}catch(e){n(e)}})},40361:e=>{e.exports=require("next/dist/compiled/next-server/pages.runtime.prod.js")},50265:(e,t,r)=>{r.d(t,{A:()=>s});var n=r(40093),o=r.n(n),a=r(40866),i=r.n(a),l=r(63685),d=i()(o());d.i(l.A),d.push([e.id,`:root {
  --background: #ffffff;
  --foreground: #171717;

  --font-sans: ui-sans-serif, system-ui, sans-serif, "Apple Color Emoji",
    "Segoe UI Emoji", "Segoe UI Symbol", "Noto Color Emoji";
  --font-mono: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas,
    "Liberation Mono", "Courier New", monospace;
}
.custom-text {
  font-family: var(--font-sans);
}

.code-snippet {
  font-family: var(--font-mono);
}

@theme inline {
  --color-background: var(--background);
  --color-foreground: var(--foreground);
  --font-sans: var(--font-geist-sans);
  --font-mono: var(--font-geist-mono);
}

@media (prefers-color-scheme: dark) {
  :root {
    --background: #0a0a0a;
    --foreground: #ededed;
  }
}

body {
  background: var(--background);
  color: var(--foreground);
  font-family: Arial, Helvetica, sans-serif;
}

table.dataTable {
  @apply w-full border-collapse border border-gray-300;
}

.dataTables_wrapper .dataTables_filter input {
  @apply border border-gray-300 rounded-md px-2 py-1;
}

.scrollbar::-webkit-scrollbar {
  width: 6px;
  height: 8px;
  cursor: pointer;
}

.scrollbar::-webkit-scrollbar-track {
  border-radius: 10vh;
  background: oklch(0.37 0.013 285.805);
}

.scrollbar::-webkit-scrollbar-thumb {
  background: #dd290a;
  border-radius: 10vh;
  /* border: 1px solid #ffffff; */
}

.scrollbar::-webkit-scrollbar-thumb:hover {
  background: #dd290a;
}

.labels {
  display: block;
  font-size: 0.875rem;
  font-weight: 500;
  color: #d2d6de;
}

.form-input-field {
  margin-top: 0.25rem;
  display: block;
  font-size: 0.875rem;
  background-color: #3f3f46;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
  padding: 0.5rem;
  border: none;
  width: 100%;
  color: #f3f4f6;
  height: 35px;
}

#form-textarea {
  margin-top: 0.25rem;
  display: block;
  font-size: 0.875rem;
  background-color: #3f3f46;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
  padding: 0.5rem;
  border: none;
  width: 100%;
  color: #f3f4f6;
}

/* .vapus-dropdown {
  position: relative;
}

.vapus-dropdown-toggle {
  padding: 6px;
  border: 1px solid oklch(0.552 0.016 285.938);
  cursor: pointer;
}

.vapus-dropdown-menu {
  border: 1px solid oklch(0.552 0.016 285.938);
  background-color: oklch(0.37 0.013 285.805);
  color: rgb(235, 235, 235);
}

.vapus-dropdown-item {
  padding: 3px;
  cursor: pointer;
  font-weight: bold;
}

.vapus-dropdown-item.parent {
  font-weight: bold;
  cursor: default;
  color: #010101;
  background-color: #cecfcf;
}

.vapus-dropdown-item.child {
  padding-left: 10px;
}

.vapus-dropdown-item.child:hover {
  background-color: oklch(0.21 0.006 285.885);
  color: #ffffff;
} */
`,""]);let s=d},57043:(e,t,r)=>{e.exports=r(63885).vendored.contexts.AmpContext},64148:(e,t,r)=>{var n=r(46458),o=r.n(n),a=r(74227),i=r.n(a),l=r(48447),d=r.n(l),s=r(44306),u=r.n(s),c=r(65278),p=r.n(c),f=r(20023),b=r.n(f),g=r(50265),h={};h.styleTagTransform=b(),h.setAttributes=u(),h.insert=d().bind(null,"head"),h.domAPI=i(),h.insertStyleElement=p(),o()(g.A,h),g.A&&g.A.locals&&g.A.locals},66631:(e,t,r)=>{Object.defineProperty(t,"__esModule",{value:!0}),Object.defineProperty(t,"default",{enumerable:!0,get:function(){return u}});let n=r(87020),o=r(8732),a=n._(r(82015)),i=n._(r(17341)),l={400:"Bad Request",404:"This page could not be found",405:"Method Not Allowed",500:"Internal Server Error"};function d(e){let t,{req:n,res:o,err:a}=e,i=o&&o.statusCode?o.statusCode:a?a.statusCode:404;if(n){let{getRequestMeta:e}=r(35124),o=e(n,"initURL");o&&(t=new URL(o).hostname)}return{statusCode:i,hostname:t}}let s={error:{fontFamily:'system-ui,"Segoe UI",Roboto,Helvetica,Arial,sans-serif,"Apple Color Emoji","Segoe UI Emoji"',height:"100vh",textAlign:"center",display:"flex",flexDirection:"column",alignItems:"center",justifyContent:"center"},desc:{lineHeight:"48px"},h1:{display:"inline-block",margin:"0 20px 0 0",paddingRight:23,fontSize:24,fontWeight:500,verticalAlign:"top"},h2:{fontSize:14,fontWeight:400,lineHeight:"28px"},wrap:{display:"inline-block"}};class u extends a.default.Component{render(){let{statusCode:e,withDarkMode:t=!0}=this.props,r=this.props.title||l[e]||"An unexpected error has occurred";return(0,o.jsxs)("div",{style:s.error,children:[(0,o.jsx)(i.default,{children:(0,o.jsx)("title",{children:e?e+": "+r:"Application error: a client-side exception has occurred"})}),(0,o.jsxs)("div",{style:s.desc,children:[(0,o.jsx)("style",{dangerouslySetInnerHTML:{__html:"body{color:#000;background:#fff;margin:0}.next-error-h1{border-right:1px solid rgba(0,0,0,.3)}"+(t?"@media (prefers-color-scheme:dark){body{color:#fff;background:#000}.next-error-h1{border-right:1px solid rgba(255,255,255,.3)}}":"")}}),e?(0,o.jsx)("h1",{className:"next-error-h1",style:s.h1,children:e}):null,(0,o.jsx)("div",{style:s.wrap,children:(0,o.jsxs)("h2",{style:s.h2,children:[this.props.title||e?r:(0,o.jsxs)(o.Fragment,{children:["Application error: a client-side exception has occurred"," ",!!this.props.hostname&&(0,o.jsxs)(o.Fragment,{children:["while loading ",this.props.hostname]})," ","(see the browser console for more information)"]}),"."]})})]})]})}}u.displayName="ErrorPage",u.getInitialProps=d,u.origGetInitialProps=d,("function"==typeof t.default||"object"==typeof t.default&&null!==t.default)&&void 0===t.default.__esModule&&(Object.defineProperty(t.default,"__esModule",{value:!0}),Object.assign(t.default,t),e.exports=t.default)},77884:e=>{e.exports=import("@auth0/nextjs-auth0")},80237:(e,t)=>{Object.defineProperty(t,"A",{enumerable:!0,get:function(){return r}});var r=function(e){return e.PAGES="PAGES",e.PAGES_API="PAGES_API",e.APP_PAGE="APP_PAGE",e.APP_ROUTE="APP_ROUTE",e.IMAGE="IMAGE",e}({})},81413:(e,t)=>{Object.defineProperty(t,"M",{enumerable:!0,get:function(){return function e(t,r){return r in t?t[r]:"then"in t&&"function"==typeof t.then?t.then(t=>e(t,r)):"function"==typeof t&&"default"===r?t:void 0}}})},82015:e=>{e.exports=require("react")},83901:(e,t)=>{Object.defineProperty(t,"__esModule",{value:!0}),Object.defineProperty(t,"warnOnce",{enumerable:!0,get:function(){return r}});let r=e=>{}},95996:(e,t,r)=>{Object.defineProperty(t,"__esModule",{value:!0}),Object.defineProperty(t,"default",{enumerable:!0,get:function(){return i}});let n=r(82015),o=()=>{},a=()=>{};function i(e){var t;let{headManager:r,reduceComponentsToState:i}=e;function l(){if(r&&r.mountedInstances){let t=n.Children.toArray(Array.from(r.mountedInstances).filter(Boolean));r.updateHead(i(t,e))}}return null==r||null==(t=r.mountedInstances)||t.add(e.children),l(),o(()=>{var t;return null==r||null==(t=r.mountedInstances)||t.add(e.children),()=>{var t;null==r||null==(t=r.mountedInstances)||t.delete(e.children)}}),o(()=>(r&&(r._pendingUpdate=l),()=>{r&&(r._pendingUpdate=l)})),a(()=>(r&&r._pendingUpdate&&(r._pendingUpdate(),r._pendingUpdate=null),()=>{r&&r._pendingUpdate&&(r._pendingUpdate(),r._pendingUpdate=null)})),null}}};var t=require("../webpack-runtime.js");t.C(e);var r=e=>t(t.s=e),n=t.X(0,[8548,6164],()=>r(26931));module.exports=n})();