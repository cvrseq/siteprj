import csv
import sqlite3

# Путь к базе данных
db_path = 'mydatabase.db'
# Путь к CSV-файлу
csv_path = './devices/Общий-Table 1.csv'

# Подключаемся к SQLite (файл базы будет создан, если не существует)
conn = sqlite3.connect(db_path)
cursor = conn.cursor()

# Создаем таблицу devices, если её еще нет
create_table_query = '''
CREATE TABLE IF NOT EXISTS devices (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    "тип" TEXT,
    "название" TEXT,
    "Модель" TEXT,
    "Топливо" TEXT,
    "Давление, атм" REAL,
    "Паропроизводительность, кг/ч" REAL,
    "температура пара" REAL,
    "КПД" REAL,
    "мощность, кВт" REAL,
    "Производство пара, кг/ч" REAL,
    "Расход газа" REAL,
    "Расход дизеля" REAL,
    "Расход мазута" REAL,
    "Расход твердого топлива" REAL,
    "вес, кг" REAL,
    "Расход твердого топлива3" REAL,
    "Расход твердого топлива4" REAL
);
'''
cursor.execute(create_table_query)
conn.commit()

# Функция для преобразования значений: если пустая строка, возвращаем None
def convert_value(value, to_type=float):
    if value.strip() == "":
        return None
    try:
        return to_type(value.replace(',', '.'))  # заменяем запятую на точку, если нужно
    except ValueError:
        return value  # если не число, возвращаем строку

# Открываем CSV-файл
with open(csv_path, newline='', encoding='utf-8') as csvfile:
    reader = csv.DictReader(csvfile, delimiter=';')
    # Для каждой строки, полученной из CSV, формируем кортеж значений для вставки.
    rows = []
    for row in reader:
        # Поскольку ключи словаря соответствуют заголовкам CSV,
        # извлекаем данные в том же порядке, как в запросе ниже.
        data = (
            row.get("тип", "").strip(),
            row.get("название", "").strip(),
            row.get("Модель", "").strip(),
            row.get("Топливо", "").strip(),
            convert_value(row.get("Давление, атм", "")),
            convert_value(row.get("Паропроизводительность, кг/ч", "")),
            convert_value(row.get("температура пара", "")),
            convert_value(row.get("КПД", "")),
            convert_value(row.get("мощность, кВт", "")),
            convert_value(row.get("Производство пара, кг/ч", "")),
            convert_value(row.get("Расход газа", "")),
            convert_value(row.get("Расход дизеля", "")),
            convert_value(row.get("Расход мазута", "")),
            convert_value(row.get("Расход твердого топлива", "")),
            convert_value(row.get("вес, кг", "")),
            convert_value(row.get("Расход твердого топлива3", "")),
            convert_value(row.get("Расход твердого топлива4", ""))
        )
        rows.append(data)

# Подготавливаем запрос для пакетной вставки
insert_query = '''
INSERT INTO devices (
    "тип",
    "название",
    "Модель",
    "Топливо",
    "Давление, атм",
    "Паропроизводительность, кг/ч",
    "температура пара",
    "КПД",
    "мощность, кВт",
    "Производство пара, кг/ч",
    "Расход газа",
    "Расход дизеля",
    "Расход мазута",
    "Расход твердого топлива",
    "вес, кг",
    "Расход твердого топлива3",
    "Расход твердого топлива4"
)
VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?);
'''

# Выполняем пакетную вставку
cursor.executemany(insert_query, rows)
conn.commit()

print(f"Импортировано {cursor.rowcount} записей.")

# Закрываем соединение
conn.close()
