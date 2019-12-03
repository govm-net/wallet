# wallet

web assembly

**if exist function wasmLoaded() in javascript,it will be called**.

about [web assembly](https://github.com/golang/go/wiki/WebAssembly)

function list:

1. loadWallet(password,walletString) return "ok"
2. newWallet(password) return "ok"
3. changePwd(old password,new password) return "ok"
4. getAddress() return AddressString
5. newTransaction(ops,...)
   * OpsTransfer:0,chain(number),peerAddress(hex),cost(number,t0),*energy(number,t0),dataString*
     * newTransaction(0,1,"01ccaf415a3a6dc8964bf935a1f40e55654a4243ae99c709",1000000000)
     * newTransaction(0,2,"01ccaf415a3a6dc8964bf935a1f40e55654a4243ae99c709",1000000000,100000000,"hello")
   * OpsMove:1,chain(number),DstChain(number),cost(number,t0),*energy(number,t0)*
   * OpsRunApp:4,chain(number),appName(hex),cost(number,t0),*data(hex),energy(number,t0)*
   * OpsUpdateAppLife:5,chain(number),appName(hex),life(number,t0),*energy(number,t0)*
