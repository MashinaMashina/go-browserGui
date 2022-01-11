## Программа для запуска Minecraft после верного решения математических примеров.

При открытии приложения запускается локальный сервер и автоматически открывается страница в браузере.
Раз в 5 секунд на сервер отправляется ajax запрос, который оповещает сервер о том, что страница еще открыта. Если запрос не приходил 11 секунд, программа закрывается.

Если по какой-то причине отключается сервер, то страница в браузере закрывается.

Прогресс хранится в JWT токене. Общение с сервером идет через REST API.

Команды для сборки:
1. go generate  browserGui/cmd/main
2. go build -ldflags "-s -H windowsgui" browserGui/cmd/main