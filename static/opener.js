function runOpener() {
    wrap().html(tpl('opener-select'))

    wrap().find('a[data-game]').on('click', function (event) {
        event.preventDefault()

        showLoader()
        api('opener/open', {game: $(this).data('game')}, function (json) {
            hideLoader()

            if (json.success) {
                wrap().html('<center><h1>Запускаем...</h1></center>');
                setTimeout(function () {
                    window.close()
                }, 500)
                return
            }

            alert(json.message)
        });
    });
}