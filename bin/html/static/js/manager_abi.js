/* Autogenerated - do not fiddle */
if (typeof Contracts === "undefined") var Contracts = { Token: { abi: [] } };

Contracts["manager"] = {
    // address: "0x25B5ab5f67a783fb9334c6A14577a7F7AAD9c55e",
    address: "0xEb12a835C2FD57A912A76590C527fBB8e5E28f91",
    network: "browser",
    endpoint: "http://localhost:8545",
    abi: [
        {
            "constant": false,
            "inputs": [
                {
                    "name": "_to",
                    "type": "address"
                },
                {
                    "name": "_amount",
                    "type": "uint256"
                },
                {
                    "name": "_trans",
                    "type": "bytes32"
                },
                {
                    "name": "approvalData",
                    "type": "bytes"
                }
            ],
            "name": "relayMint",
            "outputs": [
                {
                    "name": "",
                    "type": "bool"
                }
            ],
            "payable": false,
            "stateMutability": "nonpayable",
            "type": "function"
        },
        {
            "constant": false,
            "inputs": [
                {
                    "name": "_to",
                    "type": "address"
                },
                {
                    "name": "_amount",
                    "type": "uint256"
                },
                {
                    "name": "_trans",
                    "type": "bytes32"
                }
            ],
            "name": "mint",
            "outputs": [
                {
                    "name": "",
                    "type": "bool"
                }
            ],
            "payable": false,
            "stateMutability": "nonpayable",
            "type": "function"
        },
        {
            "constant": false,
            "inputs": [],
            "name": "transferAppOwnership",
            "outputs": [],
            "payable": false,
            "stateMutability": "nonpayable",
            "type": "function"
        },
        {
            "constant": true,
            "inputs": [],
            "name": "owner",
            "outputs": [
                {
                    "name": "",
                    "type": "address"
                }
            ],
            "payable": false,
            "stateMutability": "view",
            "type": "function"
        },
        {
            "constant": true,
            "inputs": [],
            "name": "app",
            "outputs": [
                {
                    "name": "",
                    "type": "address"
                }
            ],
            "payable": false,
            "stateMutability": "view",
            "type": "function"
        },
        {
            "constant": false,
            "inputs": [
                {
                    "name": "_value",
                    "type": "uint256"
                },
                {
                    "name": "_addr",
                    "type": "bytes"
                }
            ],
            "name": "burn",
            "outputs": [
                {
                    "name": "",
                    "type": "bool"
                }
            ],
            "payable": false,
            "stateMutability": "nonpayable",
            "type": "function"
        }
    ]
};
// console.log(data);