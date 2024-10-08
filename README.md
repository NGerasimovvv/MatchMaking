# Matchmaker Service
### Описание
Matchmaker Service — это приложение для создания групп игроков на основе их навыков, задержек и времени ожидания. Оно предоставляет API для добавления игроков и автоматически формирует группы на основе заданных критериев. Приложение поддерживает хранение данных как в памяти, так и в базе данных PostgreSQL.

### Архитектура
1. Основной поток выполнения:
   - Чтение конфигурации: Загружает параметры конфигурации, такие как размер группы и параметры подключения к базе данных.
   - Инициализация хранилища: В зависимости от конфигурации выбирается хранилище данных (в памяти или в базе данных).
   - Создание матчмейкера: Создается объект матчмейкера, который управляет процессом формирования групп.
   - Запуск процесса формирования групп: В отдельной горутине регулярно вызывается метод FormGroups, который формирует группы игроков.
   - Настройка HTTP сервера: Запускается HTTP сервер, обрабатывающий запросы на добавление игроков.

2. Компоненты:
   - Storage: Интерфейс, определяющий методы для работы с хранилищем данных (добавление, получение и удаление игроков).
   - MemoryStorage: Реализация хранилища в памяти.
   - DBStorage: Реализация хранилища в базе данных PostgreSQL.
   - Matchmaker: Основная логика формирования групп, включая расчет разницы по навыкам и задержке, поиск лучших групп и фильтрацию.
   - Player: Структура, представляющая игрока.
   - Group: Структура, представляющая группу игроков.

____


### Подробное описание алгоритма
Этот алгоритм матчмейкинга используется для формирования оптимальных групп игроков на основе их навыков (skill), задержки (latency) и времени ожидания. 

**Шаги работы алгоритма**:
1. Генерация всех возможных комбинаций игроков
2. Оценка групп на основе разницы в навыках и задержке
3. Дополнительная проверка групп на основе времени ожидания
4. Возврат лучшей группы


Ниже приведено описание методов алгоритма:

- ***Метод HandleAddPlayer:***  

  Этот метод обрабатывает HTTP-запросы для добавления новых игроков в систему.
После получения данных игрока через JSON-запрос, игрок добавляется в хранилище с текущим временем присоединения.

- ***Метод FormGroups:***

  Этот метод отвечает за формирование групп из списка доступных игроков.
Он извлекает всех игроков из хранилища и формирует группы, пока количество оставшихся игроков позволяет сформировать хотя бы одну полную группу.
Для каждой группы, которая соответствует условиям, выводится информация о группе, и игроки этой группы удаляются из хранилища.

- ***Метод FindBestGroup:***

  Этот метод находит лучшую возможную группу из всех комбинаций игроков.
Для каждой комбинации вычисляется разница в навыках и задержке, и выбирается та группа, где эти параметры минимальны и соответствуют заданным ограничениям.
Также пото идет проверка: если группа удовлетворяет более строгим условиям по навыкам и задержке - это значит, что в таких группах требования к однородности игроков значительно выше. 
После этой проверки игроки внутри группы сортируются по времени ожидания, чтобы те, кто ждут дольше, получили приоритет при формировании.

- ***FilterGroup:***

  Дополнительная проверка для учета времени ожидания.

- ***CalculateGroupMaxDifferences:***

  Эта функция вычисляет максимальную разницу между значениями в группе (разницу между максимальным и минимальным значениями навыков или задержки).

- ***Сombine:***

  Генерирует все возможные комбинации игроков для заданного размера группы.


____



- В проекте  реализована подгрузка config из .env. В файле env прописаны такие параметры, как: размер группы, переключатель хранилища(в Memory или в Postgres), параметры для Postgres.
  Их можно изменять по своему желанию. Также в main.go в бесконечном цикле в горутине есть time.Sleep на 1 секунду для регулярного формирования групп. Это также можно изменять по своему желанию, в зависимости как часто нужно формировать группы.

- В проекте реализованы несколько тестов(matchmaker_test.go, storage_test.go)

- В проекте есть директория test с файлом "test_matchmaking.go". Этот файл создан исключительно для удобства добавления игроков в пул ожидания группы.
Достаточно запустить проект("go run cmd/main.go") и запустить этот файл("go run test/test_matchmaking.go")

