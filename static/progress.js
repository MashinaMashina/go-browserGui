function progressInit() {
    if ($('.progress-bar', wrapAfter()).length) {
        return
    }

    wrapAfter().append(tpl('progress'))
}
function progressGetPercent() {
    let percent = parseInt($('.progress-now').data('percent'))

    if (! percent) {
        return 0
    }

    return percent
}
function progressAdd(percent) {
    percent = progressGetPercent() + parseInt(percent);

    progressSetPercent(percent)
}
function progressRemove(percent) {
    percent = progressGetPercent() - percent;

    progressSetPercent(percent)
}
function progressSetPercent(percent) {
    if (percent < 0) {
        percent = 0
    }
    if (percent > 100) {
        percent = 100
    }

    $('.progress-now').data('percent', percent).animate({'left': percent + '%'}, 100)
}