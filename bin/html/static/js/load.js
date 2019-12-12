var gChainID
var gCostBase
gChainID = getCookie("chain_id")
if (gChainID == "") {
    gChainID = "1"
    setCookie("chain_id", gChainID)
}
gCostBase = getCookie("cost_base")
if (gCostBase == "") {
    gCostBase = "t9"
    setCookie("cost_base", gCostBase)
}

$.get("navbar.page", function (data) {
    $("#navbar").html(data);
    var url = window.location.pathname;
    if (url == "/") {
        url = "/index.html";
    }
    $('ul.nav a[href="' + url + '"]').parent().addClass('active');
    $('ul.nav a').filter(function () {
        return this.href.pathname == url;
    }).parent().addClass('active');
});


function bytes2Str(arr) {
    var str = "";
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
    var k = bytes2Str(key)
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

function getBaseByName(name) {
    switch (name) {
        case "t3":
            return 1000;
        case "t6":
            return 1000000;
        case "t9":
            return 1000000000;
    }
    return 1
}

function getValueWithBase(v,name){
    return Math.floor((v/getBaseByName(name))*100)/100
}