function getUserById(id) {
    if (G.AllUsers && G.AllUsers.length > 0) {
        for (var i = 0; i < G.AllUsers.length; i++) {
            if (G.AllUsers[i].Id == id) {
                return G.AllUsers[i];
            }
        }
    }
}

function prependChild(p, e) {
    if (!p || !e) { return; }
    var parentChildren = getElementChildren(p);
    if (parentChildren.length == 0) {
        p.appendChild(e);
    } else {
        var firstChild = parentChildren[0];
        p.insertBefore(e, firstChild);
    }
}

function flintToString(flint, baseOffset, decimalPlaces) {
    if (!baseOffset) { baseOffset = 2; }
    if (!decimalPlaces) { decimalPlaces = 2; }
    var strVal = flint.toString();
    var major  = '';
    var minor  = '';
    if (strVal.length > baseOffset) {
        minor = strVal.substring(strVal.length - baseOffset);
        major = strVal.substring(0, strVal.length - baseOffset);
    } else {
        major = '0';
        minor = strVal;
    }
    minor = trim(padString(minor, decimalPlaces, '0', true), decimalPlaces);
    var ret = major;
    if (decimalPlaces == 0) {
        return ret;
    }
    return ret + '.' + minor;
}

function trim(str, len) {
    if (!str || !len) { return ''; }
    if (str.length < len) { return str; }
    return str.substring(0, len);
}

function padString(str, len, pad, fromBeginning) {
    if (str.length >= len) { return str; }
    if (pad.length == 0) { return ''; }
    var remaining = Math.floor((len - str.length)/pad.length)
    for (var i = 0; i < remaining; i++) {
        if (fromBeginning) {
            str = pad + str;
        } else {
            str = str + pad;
        }
    }
    return str;
}

function stringToFlint(str, baseOffset) {
    if (!baseOffset) { baseOffset = 2; }
    if (!str) { return 0; }
    var dotIdx = str.indexOf('.');
    if (dotIdx > -1) {
        var major = toInt(str.substring(0, dotIdx));
        var minor = str.substring(dotIdx + 1);
        var zeros = takeUntil(minor, function(x) { x != '0'; }).length;
        minor     = toInt(minor) * Math.pow(10, -zeros);
    } else {
        var major = toInt(str);
        var minor = 0;
    }
    return Math.round(Math.pow(10, baseOffset) * (major + minor));
}

function takeUntil(str, f) {
    if (!str || !f) { return ''; }
    var take = '';
    for (var i = 0; i < str.length; i++) {
        if (f(str[i])) { return take; }
        take += str[i];
    }
    return take;
}

function replaceIds(parent, match, replace) {
    if (!parent) { return; }
    var elements = parent.querySelectorAll('*[id^="' + match + '"]')
    for (var i = 0; i < elements.length; i++) {
        var element = elements[i];
        var newId = replace + element.id.substring(match.length);
        element.id = newId;
    }
}

function replaceMatches(parent, q, qp, m, r) {
    if (!parent || !q) { return; }
    var elements = parent.querySelectorAll('*[' + q + qp + ']');
    for (var i = 0; i < elements.length; i++) {
        var element = elements[i];
        var eq = element.getAttribute(q);
        if (eq === null) { continue; }
        var idx = eq.indexOf(m);
        if (idx > 0) {
            var newAtt = eq.substring(0, idx) + r;
            if (eq.length > idx + m.length) {
                newAtt += eq.substring(idx + m.length);
            }
            element.setAttribute(q, newAtt);
        }
    }
}

function removeCharacter(s, c) {
    if (!s) { return ''; }
    return s.replace(c, '');
}

function removeAllChildren(e) {
    if (!e) { return; }
    var children = getAllChildren(e);
    for (var i = 0; i < children.length; i++) {
        removeElement(children[i]);
    }
}

function removeElement(e) {
    if (e && e.parentNode) {
        e.parentNode.removeChild(e);
    }
}

function parseJSON(json, d) {
    var resp;
    try {
        resp = JSON.parse(json);
    } catch (e) {
        return d || e;
    }
    return resp;
}

function http(type, url, data, callback) {
    var isAsync = callback ? true : false;
    var xmlhttp = new XMLHttpRequest();
    if (isAsync) {
        xmlhttp.onreadystatechange = function() {
            var resp = parseJSON(xmlhttp.responseText, {});
            if (xmlhttp.readyState == 4 && isSuccess(xmlhttp.status)) {
                showError();
                callback(resp);
            } else if (xmlhttp.readyState == 4 && isBadRequest(xmlhttp.status)) {
                showError(resp.Message || 'Error');
            }
        }
    }
    xmlhttp.open(type, url, isAsync);
    xmlhttp.send(data || null);
    if (!isAsync) { return xmlhttp.responseText; }
}

function httpGet(urlToGet, callback) {
    return http('GET', urlToGet, null, callback);
}

function httpPost(urlToPost, data, callback) {
    return http('POST', urlToPost, data, callback);
}

function httpDelete(url, callback) {
    return http('DELETE', url, null, callback);
}

function getValueById(id) {
    var e = getElement(id);
    return getValue(e);
}

function getValue(e) {
    if (e) { return e.value || ""; }
    return "";
}

function setValueById(id, v) {
    var e = getElement(id);
    setValue(e, v);
}

function setValue(e, v) {
    if (e) { e.value = v; }
}

function setInnerHTMLById(id, v) {
    var e = getElement(id);
    setInnerHTML(e, v);
}

function setInnerHTML(e, v) {
    if (e) { e.innerHTML = v; }
}

function randomInt(min, max) {
    var rnd = Math.random();
    rnd = Math.floor(rnd * (max - min));
    return rnd + min;
}

function isInList(e, l) {
    if (!l) { return false; }
    for (var i = 0; i < l.length; i++) {
        if (e == l[i]) { return true; }
    }
    return false;
}

function toInt(v, d) {
    if (!d) { d = 0; }
    var n = parseInt(v);
    if (isNaN(n)) { return d || 0; }
    return n;
}

function getElementChildrenById(id) {
    var e = getElement(id);
    if (!e) { return []; }
    return e.children;
}

function getElementChildren(e) {
    if (!e) { return []; }
    return e.children;
}

function getAllUsers() {
    if (G.AllUsers) {
        return G.AllUsers;
    }
    return httpGet('users');
}

function getElement(id) {
    if (!id) {
        return null;
    }
    var e = document.getElementById(id);
    if (!e) {
        return null;
    }
    return e;
}

function getValue(node) {
    if (!node) { return ""; }
    if (node.value) {
        return node.value;
    }
    var e = getElement(node);
    if (!e) { return ""; }
    return e.value || "";
}

function onclickFormSubmit(formId, callback) {
    var form = getElement(formId);
    var location = form.action;
    if (!form) {
        return;
    }
    var data = makeFormData(form);
    if (!data) {
        return;
    }
    var result = httpPost(location, data, callback);
    return result;
}

function makeFormData(form) {
    var f = new FormData();
    formChildren = getAllChildren(form);
    namedChildren = [];
    for (var i = 0; i < formChildren.length; i++) {
        var child = formChildren[i];
        if (child.name) {
            namedChildren.push(child);
        }
    }
    for (var i = 0; i < namedChildren.length; i++) {
        var child = namedChildren[i];
        f.append(child.name, getValue(child));
    }
    return f;
}

function isSuccess(status) {
    return status >= 200 && status < 300;
}

function isBadRequest(status) {
    return status >= 400 && status < 500;
}

function showError(text) {
    if (text) {
        var container = getElement('errormsgcontainer');
        if (container) {
            addClassToElem(container, 'errormsgactive');
        }
        setInnerHtml('errormsg', text);
        show('errormsg');
    } else {
        var container = getElement('errormsgcontainer');
        if (container) {
            removeClassFromElem(container, 'errormsgactive');
        }
        setInnerHtml('errormsg', '');
        hide('errormsg');
    }
}

function show(id) {
    var e = getElement(id);
    removeClassFromElem(e, 'hidden');
}

function showElem(e) {
    removeClassFromElem(e, 'hidden')
}

function hide(id) {
    addClass(id, 'hidden');
}

function hideElem(e) {
    addClassToElem(e, 'hidden');
}

function addClass(id, c) {
    var e = getElement(id);
    addClassToElem(e, c);
}

function addClassToElem(e, c) {
    var cl = [];
    if (e) {
        cl = e.classList;
    } else {
        return;
    }
    if (cl.contains(c)) {
        return;
    }
    cl.add(c);
}

function removeClassFromElem(e, c){
    var cl = [];
    if(e) {
        cl = e.classList
    } else {
        return;
    }
    if (cl.contains(c)) {
        cl.remove(c);
    }
}

function hasClass(e, c) {
    if (!e) { return false; }
    var cl = e.classList;
    return cl.contains(c);
}

function setInnerHtml(id, v){
    var e = getElement(id);
    if (e) {
        e.innerHTML = v;
    }
}

function getAllChildren(node) {
    children = [];
    for (var i = 0; i < node.childNodes.length; i++) {
        children.push(node.childNodes[i]);
        children = appendList(children, getAllChildren(node.childNodes[i]));
    }
    return children;
}

function appendList(a, b) {
    var newList = [];
    if (!a || !b) {
        return [];
    }
    for (var i = 0; i < a.length; i++) {
        newList.push(a[i]);
    }
    for (var i = 0; i < b.length; i++) {
        newList.push(b[i]);
    }
    return newList;
}