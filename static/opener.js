function runOpener() {
    api("opener/available", {}, function (json) {
        let str = ""
        json.forEach(function (program) {
            str += '<a href="#" data-game="'+program.name+'" class="form-control btn btn-minecraft">'+program.name+'</a>'
        })

        let template = tpl('opener-select')
        template = template.replace('%programs%', str)

        wrap().html(template)

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
    })

}