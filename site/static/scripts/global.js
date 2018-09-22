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

function httpGet(urlToGet, callback) {
    var isAsync = callback ? true : false;
    var xmlhttp = new XMLHttpRequest();
    if (isAsync) {
        xmlhttp.onreadystatechange = function() {
            var resp = parseJSON(xmlhttp.responseText, {});
            if (xmlhttp.readyState == 4 && isSuccess(xmlhttp.status)) {
                showError();
                callback(resp);
            } else if (xmlhttp.readyState == 4 && isBadRequest(xmlhttp.status)) {
                showError(resp.Message);
            }
        }
    }
    xmlhttp.open("GET", urlToGet, isAsync);
    xmlhttp.send();
    if (!isAsync) { return xmlhttp.responseText; }
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
    if (isNaN(n)) { return d; }
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

function httpPost(urlToPost, data, callback) {
    var isAsync = callback ? true : false;
    var xmlhttp = new XMLHttpRequest();
    if (isAsync) {
        xmlhttp.onreadystatechange = function() {
            var resp = parseJSON(xmlhttp.responseText, {});
            if (xmlhttp.readyState == 4 && isSuccess(xmlhttp.status)) {
                showError();
                callback(resp);
            } else if (xmlhttp.readyState == 4 && isBadRequest(xmlhttp.status)) {
                showError(resp.Message);
            }
        }
    }
    xmlhttp.open("POST", urlToPost, isAsync);
    xmlhttp.send(data);
    if (!isAsync) { return xmlhttp.responseText; }
}

function isSuccess(status) {
    return status >= 200 && status < 300;
}

function isBadRequest(status) {
    return status >= 400 && status < 500;
}

function showError(text) {
    if(text) {
        setInnerHtml('errormsg', text);
        show('errormsg');
    } else {
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