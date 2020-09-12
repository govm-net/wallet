

/* str: string
    hex: hex string
    bytes: byte array
    array: string/number array
*/
function dataEncode(input, type) {
    function _bytes2hex(arr) {
        if (arr === undefined || arr == "") {
            return ""
        }
        var str = "";
        for (var i = 0; i < arr.length; i++) {
            var tmp = arr[i].toString(16);
            if (tmp.length == 1) {
                tmp = "0" + tmp;
            }
            str += tmp;
        }
        return str;
    };


    function stringToBytes(str) {
        if (str === undefined || str == "") {
            return []
        }
        var bytes = new Array();
        var len, c;
        len = str.length;
        for (var i = 0; i < len; i++) {
            c = str.charCodeAt(i);
            if (c >= 0x010000 && c <= 0x10FFFF) {
                bytes.push(((c >> 18) & 0x07) | 0xF0);
                bytes.push(((c >> 12) & 0x3F) | 0x80);
                bytes.push(((c >> 6) & 0x3F) | 0x80);
                bytes.push((c & 0x3F) | 0x80);
            } else if (c >= 0x000800 && c <= 0x00FFFF) {
                bytes.push(((c >> 12) & 0x0F) | 0xE0);
                bytes.push(((c >> 6) & 0x3F) | 0x80);
                bytes.push((c & 0x3F) | 0x80);
            } else if (c >= 0x000080 && c <= 0x0007FF) {
                bytes.push(((c >> 6) & 0x1F) | 0xC0);
                bytes.push((c & 0x3F) | 0x80);
            } else {
                bytes.push(c & 0xFF);
            }
        }
        return bytes;
    }


    function _bytesToString(arr) {
        if (arr === undefined || arr == "") {
            return ""
        }
        if (typeof arr === 'string') {
            return arr;
        }
        var str = '',
            _arr = arr;
        for (var i = 0; i < _arr.length; i++) {
            var one = _arr[i].toString(2);
            var v = one.match(/^1+?(?=0)/);
            if (v && one.length == 8) {
                var bytesLength = v[0].length;
                var store = _arr[i].toString(2).slice(7 - bytesLength);
                for (var st = 1; st < bytesLength; st++) {
                    store += _arr[st + i].toString(2).slice(2);
                }
                str += String.fromCharCode(parseInt(store, 2));
                i += bytesLength - 1;
            } else {
                str += String.fromCharCode(_arr[i]);
            }
        }
        return str;
    }

    function _intToHex(str, bit) {
        if (str == "" || str === undefined) {
            return ""
        }
        bit = bit || 64;

        var val = parseInt(str).toString(16);
        var n = bit / 4 - val.length;
        for (var i = 0; i < n; i++) {
            val = "0" + val;
        }
        return val
    }

    switch (type) {
        case "bytes2str":
            return _bytesToString(input)
        case "array2str":
            return input.toString()
        case "bytes2hex":
            return _bytes2hex(input)
        case "bytes2int":
            return dataEncode(dataEncode(input, "bytes2hex"), "hex2int")
        case "hex2int":
            return parseInt(input, 16)
        case "hex2bytes":
            var myUint8Array = new Uint8Array(input.match(/[\da-f]{2}/gi).map(function (h) {
                return parseInt(h, 16)
            }))
            return Array.from(myUint8Array);
        case "hex2str":
            return dataEncode(dataEncode(input, "hex2bytes"), "bytes2str")
        case "str2bytes":
            return stringToBytes(input)
        case "str2hex":
            return dataEncode(dataEncode(input, "str2bytes"), "bytes2hex")
        case "str2int":
            return parseInt(input)
        case "int322hex":
            return _intToHex(input, 32)
        case "int642hex":
            return _intToHex(input, 64)
        case "int162hex":
            return _intToHex(input, 16)
        case "int82hex":
            return _intToHex(input, 8)
        case "uristr":
            return decodeURI(escape(input))
        case "str2json":
            return JSON.parse(input)
        case "json2str":
            return JSON.stringify(input, null, 2)
        case "hex2json":
            return dataEncode(dataEncode(input, "hex2str"), "str2json")
        case "time2str":
            var myDate = new Date()
            myDate.setTime(input)
            return myDate.toString()
        case "base2bytes":
            var binary_string = window.atob(input);
            var len = binary_string.length;
            var bytes = new Uint8Array(len);
            for (var i = 0; i < len; i++) {
                bytes[i] = binary_string.charCodeAt(i);
            }
            return bytes;
        case "base2hex":
            return dataEncode(dataEncode(input, "base2bytes"), "bytes2hex")
        case "cost2num":
            if (input === undefined || input == "" || input == "0" || typeof arr !== 'string') {
                return 0;
            }
            // t0,t3,t6,t9,tc,govm=t9
            var split = input.split("t")
            if (split.length <= 1){
                split = input.split("g")
                split[1]="9"
            }
            var num = parseFloat(split[0])
            if(split.length > 1){
                var tn = parseInt(split[1],16)
                for(i=0;i<tn;i++){
                    num=num*10
                }
            }
            return num
        default:
            // console.log("unknow type:", type, input)
            return input
    }
}
