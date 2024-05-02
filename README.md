# involta_testtask
Тестовое задание для Involta на вакансию Backend-разработчика

~~1 Развернуть Reindexer в докере~~

~~2 Написать на Go микросервис с постоянным подключением к реиндексеру посредством пакета от разработчика~~
~~2.1 сделать вариативную конфигурацию~~
2.1.1 конфигурация через локальный YAML файл
~~2.1.2* конфигурация через Environment~~

~~3 Проверки при запуске~~
~~3.1 Проверка подключение к Reindexer ~~
~~3.2 Проверка наличия необходимых коллекций~~
~~3.2.1* В случае отсутствия необходимых коллекций их необходимо создать~~

~~4 Микросервис должен создавать, редактировать, выводить информацию о списке имеющихся документов или заданного документа (CRUD)~~
~~4.1* Для вывода списка предусмотреть параметры для пагинации и кол-ва выводимых документов~~

~~5 Выдача одного документа должна кешироваться~~
~~5.1* Кеш устаревает раз в 15 минут~~

~~6 Документ содержит 2 уровня вложенности каждый из которых массив документов~~

~~7* Массив документов первого уровня вложенности должен сортироваться по полю sort (целочисленное) (обратная сортировка)~~

8* Каждый документ перед отправкой ответа обрабатывается в отдельном потоке при этом общая сортировка не должна пострадать (обработка подразумевает исключение одного или нескольких полей из каждого уровня документа)