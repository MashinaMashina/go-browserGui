function keepAlive() {
    $.post('/keepalive/').fail(function () {
        window.close()
    })
}

keepAlive()
setInterval(keepAlive, 5000)

function wrapBefore() {
    return $('.app-before')
}
function wrap() {
    return $('.app')
}
function wrapAfter() {
    return $('.app-after')
}

function tpl(id) {
    return $('#' + id).html()
}
function showLoader() {
    wrap().addClass('load')
}
function hideLoader() {
    wrap().removeClass('load')
}

$(function () {
    hideLoader()
    runCalculator()
})

let authHeader;

function api(method, data, callback) {
    let headers = {};

    if (authHeader) {
        headers['Authorization'] = authHeader;
    }

    data = JSON.stringify(data)

    $.ajax({
        url: '/api/' + method,
        type: 'POST',
        data: data,
        headers: headers,
        dataType: 'json',
        contentType: "application/json"
    }).always(function (response, status, request) {
        if (status === "error") {
            request = response;
            response = response.responseJSON;
        }

        let header = request.getResponseHeader("Authorization");

        if (header) {
            authHeader = header;
        }

        callback(response, status, request)
    })
}