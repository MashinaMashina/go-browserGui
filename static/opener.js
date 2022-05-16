function runOpener() {
    wrap().html(tpl('opener-select'))

    wrap().find('a[data-game]').on('click', function (event) {
        event.preventDefault()

        showLoader()
        api('opener/open', {game: $(this).data('game')}, function (json) {
            hideLoader()

            if (json.success) {
                setTimeout(function () {
                    window.close()
                }, 500)
            }

            alert(json.message)
        });
    });
}