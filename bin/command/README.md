# wallet

1. operation:
   * wallet:create new wallet
   * show:show account info of the wallet
   * transaction: create new transaction,and write to trans_file
   * sendTrans: send transaction(trans_file) to network
2. trans_ops:
   * OpsTransfer:0
   * OpsMove:1
   * OpsRunApp:4
   * OpsUpdateAppLife:5

**how to use**.

1. build project:
     * go build
2. new wallet:
     * make sure the operation value of conf.json equal wallet
     * ./cmd
     * input password of wallet
     * it will create "wallet.key"
3. get account info
    * set operation=show
    * ./cmd
    * will get the account info
4. transfer
    * set operation=transaction
    * set transfer.peer
    * set cost of transfer
    * ./cmd
    * input password
    * data will write to transaction.data
    * set operation=sendTrans
    * ./cmd