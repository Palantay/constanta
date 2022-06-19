# Constanta API
___

## Общая информация
___
#### Реализован эмулятор платежного сервиса:
- Создание платежа (В 20% случаев статус платежа записывается, как «ОШИБКА»)
- Изменение статуса платежа платежной системой (Доступно только авторизированным пользователям)
- Проверка статуса платежа по ID
- Получение списка всех платежей пользователя по его ID
- Отмена платежа по его ID (Статус отмены платежа становится true)

#### Сервис имеет HTTP API и принимает/отдаёт запросы/ответы в формате JSON.

## Использование
___
Создайте в корне проекта файл `.env`, скопируйте в него содержимое файла `.env-example`

Для запуска приложения выполните в терминале команду: 
```
make api-build
```

У Вас должна быть установлена программа [make](https://www.gnu.org/software/make/)

Если у вас нет программы [make](https://www.gnu.org/software/make/). Тогда выполните, по очереди, следующие команды:
````
docker-compose up -d
````
````
migrate -path migrations -database "postgres://localhost:5432/db?sslmode=disable&user=admin&password=admin" up
````
````
go build -v ./cmd/api/
`````
`````
./api
`````

При выполнении команд запустится докер контейнер с базой данных (PostgreSQL) **на порту 5432**,
скомпилируется приложение и **запустится на порту 8080**