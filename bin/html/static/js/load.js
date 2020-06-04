var gChainID
var gCostBase
var gLanguage
gChainID = getCookie("chain_id")
if (gChainID == "") {
    gChainID = "1"
    setCookie("chain_id", gChainID)
}
gCostBase = getCookie("cost_base")
if (gCostBase == "") {
    gCostBase = "govm"
    setCookie("cost_base", gCostBase)
}
gLanguage = getCookie("language")

$.get("navbar.page?v="+Math.random(), function (data) {
    $("#navbar").html(data);
    var url = window.location.pathname;
    if (url == "/") {
        url = "/index.html";
    }
    $('ul.nav a[href="' + url + '"]').parent().addClass('active');
    $('ul.nav a').filter(function () {
        return this.href.pathname == url;
    }).parent().addClass('active');
    loadLanguage();
});

function loadLanguage() {
    var fn = window.location.pathname;
    if (fn == "/") {
        fn = "/index.html";
    }
    if (gLanguage != "") {
        $("[data-localize]").localize("i18n" + fn + "/govm", { language: gLanguage }).localize("i18n/navbar/govm", { language: gLanguage })
    } else {
        $("[data-localize]").localize("i18n" + fn + "/govm").localize("i18n/navbar/govm")
    }
}

function base64ToArrayBuffer(base64) {
    var binary_string = window.atob(base64);
    var len = binary_string.length;
    var bytes = new Uint8Array(len);
    for (var i = 0; i < len; i++) {
        bytes[i] = binary_string.charCodeAt(i);
    }
    return bytes;
}

function bytes2Str(arr) {
    var str = "";
    if (typeof arr === 'string') {
        arr = base64ToArrayBuffer(arr);
    }
    for (var i = 0; i < arr.length; i++) {
        var tmp = arr[i].toString(16);
        if (tmp.length == 1) {
            tmp = "0" + tmp;
        }
        str += tmp;
    }
    return str;
}

function getUrlParam(name) {
    var reg = new RegExp("(^|&)" + name + "=([^&]*)(&|$)");
    var r = window.location.search.substr(1).match(reg);
    if (r != null) return unescape(r[2]);
    return "";
}

function fmoney(s) {
    // console.log("money:"+s)
    s = parseFloat((s + "").replace(/[^\d\.-]/g, "")).toFixed(10) + "";
    var l = s.split(".")[0].split("").reverse();
    var t = "";
    for (i = 0; i < l.length; i++) {
        t += l[i] + ((i + 1) % 3 == 0 && (i + 1) != l.length ? "," : "");
    }
    return t.split("").reverse().join("");
}

function getLinkString(path, chain, key) {
    var out = ""
    var k = ""
    if (key === undefined) {
        return "NULL"
    }
    if (typeof key === "string") {
        k = key;
    } else {
        k = bytes2Str(key);
    }
    if (k == "0000000000000000000000000000000000000000000000000000000000000000") {
        return "NULL"
    }
    out += "<a href=\"" + path + "?key=" + k;
    out += "&chain=" + chain + "\">" + k;
    return out
}

function setCookie(cname, cvalue) {
    var d = new Date();
    d.setTime(d.getTime() + (7 * 24 * 60 * 60 * 1000));
    var expires = "expires=" + d.toGMTString();
    document.cookie = cname + "=" + cvalue + "; " + expires + "; path=/";
}

function getCookie(cname) {
    var name = cname + "=";
    var ca = document.cookie.split(';');
    for (var i = 0; i < ca.length; i++) {
        var c = ca[i].trim();
        if (c.indexOf(name) == 0) { return c.substring(name.length, c.length); }
    }
    return "";
}

// t0,t3,t6,t9,tc,govm=t9
function getBaseByName(name) {
    var num = 1;
    var split = name.split("t")
    if (split.length <= 1){
        split = name.split("g")
        split[1]="9"
    }
    var tn = parseInt(split[1],16)
    for (i = 0; i < tn; i++) {
        num = num * 10
    }
    return num
}

function getValueWithBase(v, name) {
    return Math.floor((v / getBaseByName(name)) * 1000) / 1000
}