# Planner

*⚠️ For private usage only. Author doesn't guarantee any future maintenance, releases or breaking changes.*

Planner - это небольшое однопользовательское приложение трекер-задач с веб-интерфейсом.

## Идея

Решение следующих проблем (задач):

- Уменьшение прокрастинации (дать возможность фиксировать время для прокрастинации или развлечений, что должно дать управление прокрастинацией, и отключить её неосознанную активацию и влияние)
- Невозможность запомнить все задачи и их последовательность выполнения
- Фиксирование задачи в электронном виде (отчётность о выполнении)
- Возможность строить долгосрочный план работ

## Поддерживаемые действия (Release plan)

### v0.1 List-based with backlog

- Получение списка задач (в приложении один пользователь)
- Добавление(создание) задачи с назначением даты начала задачи
- Редактирование задачи
    - редактирование названия задачи
    - редактирование времени
- Удаление задачи
- Изменение статуса задачи (done/in progress/new)

### v0.2 Criteria-search

- Поиск задач по критериям

### v0.3 Calendar (alpha)

- Календарь для просмотра запланированных задач

### v0.4 Calendar (beta)

- Создание и планирование задач через календарь

### v0.5 Task priority

- Назначение приоритета задаче

### v0.6 Expired-task monitoring

- Механизм авто-маркировки старых задач

### v0.7 Tags

- Тегирование задач
- Поиск по тегам

### vX.X Board-based

## Idiology of the repo/project

- Go mod doesn't contain deps
- Bare minimum CSS library (Just to make it more visually appealing)
- No direct/manual code in JS/TS/Hyperscript/whatever (only HTMX is allowed AND JS dependencies for Missing.css lib).
- Author has no prior Go/HTMX knowledge

## Dependencies

- [Missing.css](https://missing.style/)
- [HTMX](https://htmx.org/)