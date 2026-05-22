# Organization Structure API

API для управления организационной структурой (подразделения и сотрудники).

## Запуск

```bash
git clone https://github.com/Tim73916/org-structure-api
cd org-structure-api
cp .env.example .env
docker-compose up --build
```
Также можно запустить программу через  Makefile:
```bash
make docker-up   # запустить
make docker-down # остановить
make test        # запустить тесты
make run         # локальный запуск
```

## API

**POST** `/departments/` - создать подразделение

**GET** `/departments/{id}?depth=2&include_employees=true` - получить дерево

**PATCH** `/departments/{id}` - обновить/переместить

**DELETE** `/departments/{id}?mode=cascade|reassign` - удалить

**POST** `/departments/{id}/employees/` - создать сотрудника

**GET** `/employees?department_id={id}` - список сотрудников

**PUT** `/employees/{id}` - обновить сотрудника

**DELETE** `/employees/{id}` - удалить сотрудника


## Пример

```bash
curl -X POST http://localhost:8080/departments/ \
  -H "Content-Type: application/json" \
  -d '{"name":"IT"}'
```

## Стек

Go + PostgreSQL + GORM + Docker
