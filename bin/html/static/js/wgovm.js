// The object 'Contracts' will be injected here, which contains all data for all contracts, keyed on contract name:
// Contracts['Token'] = {
//  abi: [],
//  address: "0x..",
//  endpoint: "http://...."
// }

const Web3 = require("web3");

const ethEnabled = () => {
    if (window.web3) {
        window.web3 = new Web3(window.web3.currentProvider);
        window.ethereum.enable();
        return true;
    }
    return false;
}

if (!ethEnabled()) {
    alert("Please install MetaMask to use this dApp!");
    console.log("Please install MetaMask to use this dApp!");
}

// Creates an instance of the smart contract, passing it as a property,
// which allows web3.js to interact with it.
function Token(Contract) {
    this.web3 = null;
    this.instance = null;
    this.Contract = Contract;
}

// Initializes the `Token` object and creates an instance of the web3.js library,
Token.prototype.init = function () {
    // Creates a new Web3 instance using a provider
    // Learn more: https://web3js.readthedocs.io/en/v1.2.0/web3.html
    this.web3 = new Web3(
        (window.web3 && window.web3.currentProvider) ||
        new Web3.providers.HttpProvider(this.Contract.endpoint)
    );

    // Creates the contract interface using the web3.js contract object
    // Learn more: https://web3js.readthedocs.io/en/v1.2.0/web3-eth-contract.html#new-contract
    var contract_interface = this.web3.eth.contract(this.Contract.abi);
    
    // Defines the address of the contract instance
    this.instance = this.Contract.address
        ? contract_interface.at(this.Contract.address)
        : { mintTokens: () => { } };
};

// Returns the token balance (from the contract) of the given address
Token.prototype.getBalance = function (address, cb) {
    this.instance.balanceOf(address, function (error, result) {
        cb(error, result);
    });
};

// Returns the token balance (from the contract) of the given address
Token.prototype.allowance = function (address) {
    this.instance.allowance(
        window.web3.eth.accounts[0],
        address, function (error, result) {
            if (error) {
                console.log(error);
            } else {
                $("#allowance").val(result.toNumber() / 1E9);
            }
        });
};

Token.prototype.approve = function (address) {
    var that = this;
    if ($("#allowance").val()>1e10){
        $("#burn_result").html("not need approve");
        return;
    }
    this.instance.approve(
        address,
        1e25,
        {
            from: window.web3.eth.accounts[0],
            gas: 1000000,
            gasPrice: 1e11,
            gasLimit: 100000
        },
        function (error, txHash) {
            if (error) {
                console.log(error);
            }
            // If success, wait for confirmation of transaction,
            // then clear form values
            else {
                // $("#burn_result").html("not need approve");
                // that.allowance(address);
                $("#burn_result").html("<span class=\"label label-success\">tx:<a target=\"_blank\" href=\"https://etherscan.io/tx/"
                + txHash + "\">" + txHash + "</a></span>");
                that.waitForReceipt(txHash, function(r){
                    that.allowance(address);
                });
            }
        }
    );
};

// Sends tokens to another address, triggered by the "Mint" button
Token.prototype.mintTokens = function () {
    $("#mint_result").html("");
    var that = this;
    // Gets form input values
    var address = $("#mint_eth_addr").val();
    var amount = $("#mint_amount1").val();
    var trans = $("#mint_key").val();
    var sign = $("#mint_signature").val();
    console.log(amount);

    // Validates address using utility function
    if (!isValidAddress(address)) {
        console.log("Invalid address");
        $("#mint_result").html("Invalid eth address");
        return;
    }

    // Validate amount using utility function
    if (!isValidAmount(amount) || amount < 1e10) {
        console.log("Invalid amount");
        $("#mint_result").html("Invalid amount");
        return;
    }

    if (trans.length == 0) {
        console.log("Invalid transaction key");
        $("#mint_result").html("Invalid transaction key");
        return;
    }

    if (sign.length == 0) {
        console.log("Invalid Signature");
        $("#mint_result").html("Invalid Signature");
        return;
    }
    console.log("relayMint:", address, amount, trans, sign)
    // Calls the public `mint` function from the smart contract
    this.instance.relayMint(
        address,
        amount,
        trans,
        sign,
        {
            from: window.web3.eth.accounts[0],
            gas: 1000000,
            gasPrice: 1e11,
            gasLimit: 1000000
        },
        function (error, txHash) {
            if (error) {
                console.log(error);
            }
            // If success, wait for confirmation of transaction,
            // then clear form values
            else {
                $("#mint_result").html("tx:" + txHash);
                $("#mint_result").html("<span class=\"label label-success\">tx:<a target=\"_blank\" href=\"https://etherscan.io/tx/"
                    + txHash + "\">" + txHash + "</a></span>");
            }
        }
    );
};


// Sends tokens to another address, triggered by the "Mint" button
Token.prototype.adminMintTokens = function () {
    $("#mint_result").html("");
    var that = this;
    // Gets form input values
    var address = $("#mint_eth_addr").val();
    var amount = $("#mint_amount1").val();
    var trans = $("#mint_key").val();
    console.log(amount);

    // Validates address using utility function
    if (!isValidAddress(address)) {
        console.log("Invalid address");
        $("#mint_result").html("Invalid eth address");
        return;
    }

    // Validate amount using utility function
    if (!isValidAmount(amount) || amount < 1e10) {
        console.log("Invalid amount");
        $("#mint_result").html("Invalid amount");
        return;
    }

    if (trans.length == 0) {
        console.log("Invalid transaction key");
        $("#mint_result").html("Invalid transaction key");
        return;
    }

    console.log("admin mint:", address, amount, trans)
    // Calls the public `mint` function from the smart contract
    this.instance.mint(
        address,
        amount,
        trans,
        {
            from: window.web3.eth.accounts[0],
            gas: 1000000,
            gasPrice: 1e11,
            gasLimit: 1000000
        },
        function (error, txHash) {
            if (error) {
                console.log(error);
            }
            // If success, wait for confirmation of transaction,
            // then clear form values
            else {
                $("#mint_result").html("tx:" + txHash);
                $("#mint_result").html("<span class=\"label label-success\">tx:<a target=\"_blank\" href=\"https://etherscan.io/tx/"
                    + txHash + "\">" + txHash + "</a></span>");
            }
        }
    );
};


// Sends tokens to another address, triggered by the "Mint" button
Token.prototype.burnTokens = function (balance) {
    $("#burn_result").html("");

    var that = this;
    // Gets form input values
    var address = $("#govm_addr").val();
    var amount = $("#burn_amount").val();
    console.log("start burn:", address, amount);

    // Validates address using utility function
    if (address.length != 48) {
        console.log("Invalid address.", address);
        $("#burn_result").html("Invalid govm address");
        return;
    }

    // Validate amount using utility function
    if (!isValidAmount(amount) || amount < 10) {
        console.log("Invalid amount");
        $("#burn_result").html("Invalid amount");
        return;
    }

    var allowance = $("#allowance").val();
    if (amount > allowance * 1e9) {
        $("#burn_result").html("Approve first");
        return;
    }

    amount = amount * getBaseByName(gCostBase);
    address = "0x" + address;
    // console.log("burn001:", address, amount)
    if(amount > balance){
        $("#burn_result").html("not enough balance. have:"+balance/1e9);
        return;
    }

    // burn wgovm, swap to govm
    this.instance.burn(
        amount,
        address,
        {
            from: window.web3.eth.accounts[0],
            gas: 1000000,
            gasPrice: 1e11,
            gasLimit: 1000000
        },
        function (error, txHash) {
            if (error) {
                console.log(error);
                return
            }
            $("#burn_result").html("<span class=\"label label-success\">tx:<a target=\"_blank\" href=\"https://etherscan.io/tx/"
                + txHash + "\">" + txHash + "</a></span>");
        }
    );

};

// Waits for receipt of transaction
Token.prototype.waitForReceipt = function (hash, cb) {
    var that = this;

    // Checks for transaction receipt using web3.js library method
    this.web3.eth.getTransactionReceipt(hash, function (err, receipt) {
        if (err) {
            error(err);
        }
        if (receipt !== null) {
            // Transaction went through
            if (cb) {
                cb(receipt);
            }
        } else {
            // Try again in 2 second
            window.setTimeout(function () {
                that.waitForReceipt(hash, cb);
            }, 2000);
        }
    });
};
Token.prototype.showTotal = function () {
    var that = this;
    this.instance.totalSupply(function (error, total) {
        if (error) {
            console.log(error);
        } else {
            console.log("wgovm total:", total / 1E9, that.Contract.address);
            $("#wgovm_total").html(total.toNumber() / 1E9);
        }
    });
}


// Checks if the contract has been deployed.
// A contract will not have its address set until it has been deployed
Token.prototype.hasContractDeployed = function () {
    return this.instance && this.instance.address;
    // return true
};

// Creates the instance of the `Token` object
Token.prototype.onReady = function () {
    this.init();
};

// Checks if it has the basic requirements of an address
function isValidAddress(address) {
    return /^(0x)?[0-9a-f]{40}$/i.test(address);
}

// Basic validation of amount. Bigger than 0 and typeof number
function isValidAmount(amount) {
    return amount > 0 && typeof Number(amount) == "number";
}

if (typeof Contracts === "undefined") var Contracts = { Token: { abi: [] } };
var token = new Token(Contracts["wgovm"]);
var manager = new Token(Contracts["manager"]);

$(document).ready(function () {
    token.onReady();
    manager.onReady();
    token.showTotal();
    
    var account = "";
    var haveBalance = 0;
    setInterval(function () {
        if (web3.eth.accounts[0] !== account) {
            account = web3.eth.accounts[0];
            token.allowance(manager.Contract.address);
            $("#burn_result").html("");
            token.getBalance(account,
                function (error, balance) {
                    if (error) {
                        console.log(error);
                    } else {
                        haveBalance = balance.toNumber();
                    }
                });
        }
    }, 200);

    $(document).on("click", "#relayMint", function () {
        manager.mintTokens();
    });

    $(document).on("click", "#mint", function () {
        manager.adminMintTokens();
    });

    $(document).on("click", "#button_burn", function () {
        manager.burnTokens(haveBalance);
    });
    $(document).on("click", "#btn_approve", function () {
        token.approve(manager.Contract.address);
    });
});
