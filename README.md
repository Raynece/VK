Task 1:
Subscribe: Создает подписку и запускает горутину для обработки сообщений.

Publish: Рассылает сообщения всем подписчикам асинхронно.

Close: Останавливает все горутины с учетом контекста.

Используется sync.Map для потокобезопасного хранения подписчиков.

Task2:
Graceful Shutdown: Сервер корректно завершает работу при получении сигналов.

Dependency Injection: Экземпляр subpub передается в сервер через конструктор.

Логирование: Используется стандартный логгер с настройкой уровня через конфиг.

Тестирование: Написаны unit-тесты для subpub, проверяющие FIFO и обработку ошибок.
