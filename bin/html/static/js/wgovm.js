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
    // var contract_interface = new this.web3.eth.Contract(this.Contract.abi,this.Contract.address);

    // Defines the address of the contract instance
    this.instance = this.Contract.address
        ? contract_interface.at(this.Contract.address)
        : { mintTokens: () => { } };

};

// Displays the token balance of an address, triggered by the "Check balance" button
Token.prototype.showAddressBalance = function (hash, cb) {
    var that = this;
    $("#message").text("Waiting");
    // Gets form input value
    var address = $("#search_addr").val();

    // Validates address using utility function
    if (!isValidAddress(address)) {
        console.log("Invalid address");
        return;
    }

    // Gets the value stored within the `balances` mapping of the contract
    this.getBalance(address, function (error, balance) {
        if (error) {
            console.log(error);
        } else {
            console.log(balance / 1E9);
            $("#message").text(balance.toNumber() / 1E9);
        }
    });
};

// Returns the token balance (from the contract) of the given address
Token.prototype.getBalance = function (address, cb) {
    this.instance.balanceOf(address, function (error, result) {
        cb(error, result);
    });
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
Token.prototype.burnTokens = function () {
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

    amount = amount * getBaseByName(gCostBase);
    address = "0x" + address;
    // console.log("burn001:", address, amount)

    // burn wgovm, swap to govm
    this.instance.burn(
        amount,
        address,
        {
            from: window.web3.eth.accounts[0],
            gas: 1000000,
            gasPrice: 0,
            gasLimit: 200000
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
    this.instance.totalSupply(function (error, total) {
        if (error) {
            console.log(error);
        } else {
            console.log("wgovm total:", total / 1E9);
            $("#wgovm_total").html(total.toNumber() / 1E9);
        }
    });
}

// Binds functions to the buttons defined in app.html
Token.prototype.bindButtons = function () {
    var that = this;

    $(document).on("click", "#relayMint", function () {
        that.mintTokens();
    });

    $(document).on("click", "#button-check", function () {
        that.showAddressBalance();
    });

    $(document).on("click", "#button_burn", function () {
        that.burnTokens();
    });
};

// Removes the welcome content, and display the main content.
// Called once a contract has been deployed
Token.prototype.updateDisplayContent = function () {
    this.hideWelcomeContent();
    this.showMainContent();
};

// Checks if the contract has been deployed.
// A contract will not have its address set until it has been deployed
Token.prototype.hasContractDeployed = function () {
    return this.instance && this.instance.address;
    // return true
};

Token.prototype.hideWelcomeContent = function () {
    $("#welcome-container").addClass("hidden");
};

Token.prototype.showMainContent = function () {
    $("#main-container").removeClass("hidden");
};

// Creates the instance of the `Token` object
Token.prototype.onReady = function () {
    this.init();
    if (this.hasContractDeployed()) {
        this.updateDisplayContent();
        this.bindButtons();
        this.showTotal();
    }
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

$(document).ready(function () {
    token.onReady();
});
