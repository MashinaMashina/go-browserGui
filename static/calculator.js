function runCalculator() {
    setQuestion()
}

function setQuestion() {
    showLoader()
    $.post('/api/calculator/getQuestion').done(function (resp) {
        hideLoader()

        let json = JSON.parse(resp)
        let template = tpl('calculate-question');

        template = template.replace('%question%', json.question_text)

        wrap().html(template)

        wrap().find('form').on('submit', function (event) {
            event.preventDefault()

            $.post('/api/calculator/checkAnswer').done(function (resp) {
                console.log(resp)
            });
        })
    })
}