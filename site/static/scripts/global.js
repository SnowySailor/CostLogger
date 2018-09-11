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
            var resp = {};
            try { resp = JSON.parse(xmlhttp.responseText); }
            catch { }
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
    removeClass(id, 'hidden');
}

function hide(id) {
    addClass(id, 'hidden');
}

function addClass(id, c) {
    var e = getElement(id);
    if(e) {
        e = e.classList;
    } else {
        return;
    }
    if (e.contains(c)) {
        return;
    }
    e.add(c);
}

function removeClass(id, c){
    var e = getElement(id);
    if(e) {
        e = e.classList
    } else {
        return;
    }
    if (e.contains(c)) {
        e.remove(c);
    }
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