# Файлы для итогового задания

В директории `tests` находятся тесты для проверки API, которое должно быть реализовано в веб-сервере.

Директория `web` содержит файлы фронтенда.

#результат выполения итогового задания

В проекте разработан веб-сервер, реализующий интерфейс API. С его помощью возможно:
1. Выполнять функции планировщика задач;
2. Вести контроль планирования, выполнения, поиска задач, их редактирования и т.д. в соответствии с заданием.

3. Пошагово все задания - в пакете "steps"
4. В файле .env находятся переменные окружения, которые используются в коде, такие как PORT с номером порта, DBFILE с названием файла базы данных и TODO_PASSWORD с паролем для шага аутентификации.
5. В tests/settings.go следует использовать: 
	5.1. var Port = 7540 
	5.2. var DBFile = "../scheduler.db" 
	5.3. var FullNextDate = true 
	5.4. var Search = true 
	5.5. var Token : последнее значение, полученное из настроек в инструментах разаботчика, при котором тесты проходят успешно.

6. В репозитории размещен Dockerfile. 
	Локально был создан докер-образ с разработанным приложением. 
	Контейнер запускается с параметрами, планировщик работает в браузере, подключается к SQLite базе данных на хосте в ОС.
	Команды создания и запуска:
		docker build -t todo-app:v1.0.0 . docker run -d -p 7540:7540 todo-app:v1.0.0

