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
   * OpsTransfer:0,chain(number),cost(number,t0),peerAddress(hex),*energy(number,t0),dataString*
     * newTransaction(0,1,1000000000,"01ccaf415a3a6dc8964bf935a1f40e55654a4243ae99c709")
     * newTransaction(0,2,1000000000,"01ccaf415a3a6dc8964bf935a1f40e55654a4243ae99c709",100000000,"hello")
   * OpsMove:1,chain(number),cost(number,t0),DstChain(number),*energy(number,t0)*
   * OpsRunApp:4,chain(number),cost(number,t0),appName(hex),*data(hex),energy(number,t0)*
   * OpsUpdateAppLife:5,chain(number),0,appName(hex),life(number,t0),*energy(number,t0)*
