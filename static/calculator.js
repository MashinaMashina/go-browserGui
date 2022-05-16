function runCalculator() {
    progressInit()
    progressSetPercent(50)
    newQuestion()

    wrapAfter().find('.progress-bar .finish').on('click', function (event) {
        if (event.ctrlKey) {
            let e = wrap().find('[name="answer"]')

            if (e.attr('type') === 'password') {
                e.attr('type', 'number')
            } else {
                e.attr('type', 'password')
            }
        }
    })
}

function newQuestion() {
    $('[name="answer"]').removeClass('border-green').removeClass('border-red');

    showLoader()
    api('calculator/getQuestion', {}, function (json) {
        hideLoader()

        let template = tpl('calculate-question');

        template = template.replace('%question%', json.question_text)

        wrap().html(template)

        $('[name="answer"]').focus();

        wrap().find('form').on('submit', function (event) {
            event.preventDefault()

            api('calculator/checkAnswer', {answer: $('[name="answer"]', this).val()}, function (resp) {
                progressSetPercent(resp.progress)

                if (resp.success) {
                    $('[name="answer"]').addClass('border-green')
                } else {
                    $('[name="answer"]').addClass('border-red')
                }

                setTimeout(function () {
                    if (resp.code === 'you_lost') {
                        alert('Вы проиграли')
                        location.reload()
                    } else if (resp.code === 'you_win') {
                        runOpener()
                    } else {
                        setTimeout(newQuestion, 150)
                    }
                }, 150)
            });
        })
    })
}
