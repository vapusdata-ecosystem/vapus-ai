"use strict";(()=>{var o={};o.id=636,o.ids=[636],o.modules={8732:o=>{o.exports=require("react/jsx-runtime")},36411:(o,r,e)=>{e.a(o,async(o,a)=>{try{e.r(r),e.d(r,{default:()=>l});var n=e(8732);e(82015),e(64148);var d=e(77884),t=o([d]);function l({Component:o,pageProps:r}){return(0,n.jsx)(d.Auth0Provider,{children:(0,n.jsx)(o,{...r})})}d=(t.then?(await t)():t)[0],a()}catch(o){a(o)}})},50265:(o,r,e)=>{e.d(r,{A:()=>i});var a=e(40093),n=e.n(a),d=e(40866),t=e.n(d),l=e(63685),s=t()(n());s.i(l.A),s.push([o.id,`:root {
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
`,""]);let i=s},64148:(o,r,e)=>{var a=e(46458),n=e.n(a),d=e(74227),t=e.n(d),l=e(48447),s=e.n(l),i=e(44306),c=e.n(i),f=e(65278),p=e.n(f),u=e(20023),b=e.n(u),m=e(50265),g={};g.styleTagTransform=b(),g.setAttributes=c(),g.insert=s().bind(null,"head"),g.domAPI=t(),g.insertStyleElement=p(),n()(m.A,g),m.A&&m.A.locals&&m.A.locals},77884:o=>{o.exports=import("@auth0/nextjs-auth0")},82015:o=>{o.exports=require("react")}};var r=require("../webpack-runtime.js");r.C(o);var e=o=>r(r.s=o),a=r.X(0,[6164],()=>e(36411));module.exports=a})();