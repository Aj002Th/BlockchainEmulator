(self.webpackChunk_N_E=self.webpackChunk_N_E||[]).push([[506],{2027:function(e,t,n){Promise.resolve().then(n.bind(n,2540))},2540:function(e,t,n){"use strict";n.r(t),n.d(t,{default:function(){return z}});var s=n(7437),r=n(5430),l=n(2752);function o(){return(0,s.jsx)("h2",{children:"欢迎来到界面，点击按钮以连接Supervisor"})}var a=n(2067),i=n(7013),c=n(7394),d=n(2943),u=n(305),x=n(9944),h=n(5171),m=n(5771),p=n(7600);class g{async retryTillConnect(){for(console.log("in retryTillConnect");;)try{console.log("creating websocket"),new x.x;let e=(0,p.j)({url:"ws://localhost:7697"});return console.log("now racing for open"),e}catch(e){console.error("Connection failed"),this.fsm.toFailed()}}async receiveMsgLoop(e){let t=e.asObservable();console.log("in receiveMsgLoop"),t.pipe((0,m.h)(e=>"heartbeat"==e.type));let n=t.pipe((0,m.h)(e=>"heartbeat"!=e.type));for(;;){try{console.log("before receive msg");var s=await (0,h.z)(n);console.log("after receive msg,before fsm.next"),this.fsm.next(s)}catch(e){console.log("ERR happened: ",e),this.fsm.toFailed();return}if("bye"==s.type)return}}async start(){console.log("in start");let e=new WebSocket("ws://localhost:7697/api/status");e.onopen=t=>{console.log("opened"),e.onmessage=e=>{this.fsm.next(JSON.parse(e.data))}},e.onerror=e=>{console.log("error happened"),console.log(e),this.fsm.toFailed()}}constructor(){this.fsm=new j}}var v={type:"no_use",content:null};class j{init(e){return(0,u.z)(()=>{l.i7.val=l.qb.Init}),this.connecting}connecting(e){if(console.log("there is connecting. the m is:"),console.log(e),console.log("the m.type is ".concat(e.type)),"hello"==e.type)return(0,u.z)(()=>{l.i7.val=l.qb.Connected}),this.connected;throw Error()}connected(e){if("started"==e.type)return(0,u.z)(()=>{l.i7.val=l.qb.Started}),this.computing;throw Error()}computing(e){if("computing"==e.type)return(0,u.z)(()=>{l.i7.val=l.qb.Running,l.YD.val=e.content}),this.computing;if("completed"==e.type)return(0,u.z)(()=>{l.i7.val=l.qb.Completed,l.Hj.val=e.content}),this.completed;throw Error()}completed(e){if("bye"==e.type)return this.completed;throw Error()}failed(e){throw Error()}toFailed(){(0,u.z)(()=>{l.i7.val=l.qb.Disconnected}),this.state=this.failed}next(e){console.log("FSM got Msg: ".concat(e)),console.log(e),console.log("FSM state before:"),console.log(this.state);try{this.state=this.state(e)}catch(e){console.log("FSM run into error catched:"),console.log(e)}console.log("FSM state after:"),console.log(this.state)}constructor(){this.state=this.init,this.next(v)}}var f=n(5096),b=n(2265),y=n(4129);let w=()=>{console.log("start the receiver"),new g().start().then(()=>{console.log("Receiver Exit.")})},k=(0,r.Pi)(e=>{let{status:t}=e;return(0,b.useEffect)(()=>{w()},[]),(0,s.jsx)("div",{className:"mx-auto",children:(0,s.jsxs)("h1",{className:"pb-3",children:["Emulator运行中 状态为：",(0,s.jsx)("code",{className:"font-bold text-xl",children:l.qb[t.val]})]})})}),N=(0,r.Pi)(f.y),C=(0,r.Pi)(e=>{let{status:t}=e;return(0,s.jsx)("div",{children:(()=>{switch(t.val){case l.qb.Init:return(0,s.jsxs)(s.Fragment,{children:["初始状态 正在连接到Sup服务器。",(0,s.jsx)(a.Z,{size:"large"}),(0,s.jsx)(o,{})]});case l.qb.Connected:return(0,s.jsxs)("div",{className:"flex flex-col items-center",children:[(0,s.jsx)(a.Z,{size:"large"}),"连接成功 Sup还没开始运行。等待Sup开始运行。"]});case l.qb.ConnectionFailed:return(0,s.jsx)(i.ZP,{status:"error",title:"连接失败",subTitle:"请关闭后重新启动"});case l.qb.Started:return(0,s.jsxs)("div",{className:"flex flex-col items-center",children:[(0,s.jsx)(a.Z,{size:"large"}),"Sup已开始运行。等待Sup的返回结果。"]});case l.qb.Running:return(0,s.jsx)("div",{className:"flex",children:(0,s.jsxs)("div",{className:"rounded-lg border-2 border-blue-500 shadow mx-auto  w-6/12",children:[(0,s.jsxs)(c.Z,{gap:"middle",align:"center",justify:"space-between",style:{margin:16},vertical:!0,children:[(0,s.jsx)(d.Z,{type:"dashboard",steps:8,percent:Math.round(100*l.YD.val.count/l.YD.val.total),trailColor:"rgba(0, 0, 0, 0.06)",strokeWidth:20}),(0,s.jsx)(a.Z,{size:"large"})]}),l.YD.val.count!=l.YD.val.total?(0,s.jsxs)(s.Fragment,{children:[(0,s.jsx)("div",{className:"m-8 text-xl",children:"正在运行中:"}),(0,s.jsxs)("div",{className:"m-8 text-xl",children:["有",(0,s.jsx)("b",{children:l.YD.val.count}),"个交易处理完毕"]}),(0,s.jsxs)("div",{className:"m-8 text-xl",children:["这次运行的总量为",(0,s.jsx)("b",{children:l.YD.val.total}),"个交易。"]})]}):(0,s.jsx)(s.Fragment,{children:"处理已结束，Supervisor正在汇总Metrics。等待Supervisor返回Complete消息。"})]})});case l.qb.RunningFailed:return(0,s.jsx)(i.ZP,{status:"error",title:"运行失败",subTitle:"请关闭后重新启动"});case l.qb.Completed:return(0,s.jsxs)("div",{className:"rounded-lg border-2 border-blue-500 shadow",children:[(0,s.jsx)(i.ZP,{icon:(0,s.jsx)(y.Z,{}),title:"运行完毕!"}),(0,s.jsx)("article",{className:"text-wrap rounded-lg border-8 shadow bg-white",children:(0,s.jsx)("div",{className:"rounded-lg border-white border-8",children:(0,s.jsx)(N,{report:l.Hj})})})]});case l.qb.Disconnected:return(0,s.jsx)(s.Fragment,{children:"已退出"});default:return(0,s.jsx)(s.Fragment,{})}})()})});var S=()=>(0,s.jsxs)("div",{className:"absolute left-1/2 top-5 translate-x-[-50%]  w-6/12 flex flex-col",children:[(0,s.jsx)(k,{status:l.i7}),(0,s.jsx)(C,{status:l.i7})]}),F=n(7138),Z=n(9887),q=n(7449),D=n.n(q);function z(){return(0,s.jsxs)(s.Fragment,{children:[(0,s.jsx)(D(),{children:(0,s.jsx)("title",{children:"Blockemulator Supervisor Monitor"})}),(0,s.jsx)(F.default,{href:"https://github.com/Aj002Th/BlockchainEmulator",children:(0,s.jsx)(Z.Z,{})}),(0,s.jsx)(S,{})]})}},5096:function(e,t,n){"use strict";n.d(t,{y:function(){return M}});var s=n(7437),r=n(2752),l=n(9551),o=n(2449),a=n(7817),i=n(8924),c=n(1469),d=n(7129),u=n(2870),x=n(2265),h=n(5211),m=n(7346);h.kL.register(h.uw,h.f$,h.ZL,h.Dx,h.u,h.De);let p={plugins:{title:{display:!0,text:"PBFT交易池统计结果 - 柱状图"}},responsive:!0,interaction:{mode:"index",intersect:!1},scales:{x:{stacked:!0},y:{stacked:!0}}},g=e=>{let{report:t}=e,{pbftShardCsv:n}=t.val,r={labels:n.map((e,t)=>"Round ".concat(t+1)),datasets:[{label:"交易池大小",data:n.map(e=>e.txpool_size),backgroundColor:"rgb(255, 99, 132)",stack:"Stack 0"},{label:"交易计数",data:n.map(e=>e.tx),backgroundColor:"rgb(75, 192, 192)",stack:"Stack 1"}]};return(0,s.jsx)(m.$Q,{options:p,data:r})};var v=n(2249),j=n(9919),f=n(9523),b=n(7840),y=n(1083),w=n(2599),k=n.n(w);let N=e=>e.map((e,t)=>({label:e.name,key:"".concat(t),children:(0,s.jsxs)("div",{children:[(0,s.jsxs)("div",{className:"italic text-sky-600",children:["“",e.name,"”统计量： ",e.desc]}),e.elems.map((e,t)=>{let n=e.desc?e.desc:"（没有描述）";return(0,s.jsxs)("div",{className:"flex flex-row content-between m-4 items-center",children:[(0,s.jsx)("div",{className:"shrink mx-4 font-bold",children:e.name}),(0,s.jsx)(f.Z,{title:n,placement:"right",children:(0,s.jsx)("div",{className:"shrink text-xs font-mono",children:(0,s.jsx)(y.Z,{})})}),(0,s.jsx)("div",{className:"grow"}),(0,s.jsx)("div",{className:"font-bold font-mono",children:e.val})]},e.name)})]})})),C=e=>k().range(e).map(e=>e.toString());var S=e=>{let{mos:t}=e;return(0,s.jsx)(b.Z,{items:N(t),defaultActiveKey:C(t.length),onChange:e=>{console.log(e)}})};let{Title:F,Paragraph:Z,Text:q,Link:D}=l.default,z=()=>(console.log(r.Hj),(0,s.jsx)(s.Fragment,{})),E=[{title:"轮数",dataIndex:"round",key:"round"},{title:"txpool大小",dataIndex:"txpool_size",key:"txpool_size"},{title:"tx计数",dataIndex:"tx",key:"tx"}],P=x.createContext({name:"Default"}),M=e=>{let{report:t}=e;var n=t.val.pbftShardCsv.map((e,t)=>{let{txpool_size:n,tx:s,ctx:r}=e;return{round:t+1,txpool_size:n,tx:s}});let r=[{key:"1",label:"表格视图",children:(0,s.jsx)(o.Z,{dataSource:n,columns:E}),icon:(0,s.jsx)(d.Z,{})},{key:"2",label:"柱状图视图",children:(0,s.jsx)(g,{report:t}),icon:(0,s.jsx)(u.Z,{})}],[l,h]=a.ZP.useNotification(),m=(0,x.useMemo)(()=>({name:"切换成功"}),[]);return console.log("MOS:"),console.log(t.val.measureOutputs),(0,s.jsxs)(s.Fragment,{children:[(0,s.jsxs)("div",{className:"flex flex-row justify-between items-center",children:[(0,s.jsx)("h3",{className:"text-3xl border-white font-bold m-4",children:"报告"}),(0,s.jsx)(i.ZP,{className:"mr-12",type:"primary",onClick:()=>{var e,n;e="data:text/json;charset=utf-8,"+encodeURIComponent(JSON.stringify(t.val)),(n=document.createElement("a")).setAttribute("href",e),n.setAttribute("download","report.json"),document.body.appendChild(n),n.click(),n.remove()},children:"导出并下载json"})]}),(0,s.jsx)("div",{className:"font-bold text-xl my-8 mx-2",children:"PBFT交易池统计结果"}),(0,s.jsxs)(P.Provider,{value:m,children:[h,(0,s.jsx)(c.Z,{defaultActiveKey:"1",items:r,centered:!0})]}),(0,s.jsx)("div",{className:"font-bold text-xl my-8 mx-2",children:"测度输出"}),z(),(0,s.jsx)(S,{mos:t.val.measureOutputs}),(0,s.jsx)("div",{className:"font-bold text-xl my-8 mx-2",children:"原始输出"}),(0,s.jsx)(v.ZP,{value:t.val.measureOutputs,style:j.u})]})}},2752:function(e,t,n){"use strict";n.d(t,{Hj:function(){return i},YD:function(){return a},i7:function(){return o},qb:function(){return r}});var s,r,l=n(305);(s=r||(r={}))[s.Init=0]="Init",s[s.Connected=1]="Connected",s[s.ConnectionFailed=2]="ConnectionFailed",s[s.Started=3]="Started",s[s.Running=4]="Running",s[s.RunningFailed=5]="RunningFailed",s[s.Completed=6]="Completed",s[s.Disconnected=7]="Disconnected";let o=(0,l.ky)({val:0}),a=(0,l.ky)({val:{count:0,total:0}}),i=(0,l.ky)({val:{pbftShardCsv:[],measureOutputs:[]}})}},function(e){e.O(0,[866,674,147,710,317,924,712,632,136,971,23,744],function(){return e(e.s=2027)}),_N_E=e.O()}]);