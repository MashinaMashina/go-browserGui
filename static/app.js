function keepAlive() {
    $.post('/api/keepalive').fail(function () {
        window.close()
    })
}

keepAlive()
setInterval(keepAlive, 5000)

function wrap() {
    return $('.wrapper')
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
