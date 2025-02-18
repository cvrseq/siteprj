import csv
import sqlite3
import glob
import os

# Путь к базе данных и директории с CSV
db_path = 'mydatabase.db'
csv_directory = './devices'  # например, папка csv_files

# Подключаемся к SQLite (файл базы будет создан, если не существует)
conn = sqlite3.connect(db_path)
cursor = conn.cursor()

# Создаем таблицу, если она еще не существует (пример для полной схемы)
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

def convert_value(value, to_type=float):
    value = value.strip()
    if value == "":
        return None
    try:
        # Если число записано с запятой, заменяем на точку
        return to_type(value.replace(',', '.'))
    except ValueError:
        return value

# Собираем список всех CSV-файлов в указанной папке
csv_files = glob.glob(os.path.join(csv_directory, '*.csv'))

# Счетчик импортированных записей
total_imported = 0

# Перебираем файлы
for csv_file in csv_files:
    print(f'Импортируем файл: {csv_file}')
    with open(csv_file, newline='', encoding='utf-8') as f:
        reader = csv.DictReader(f, delimiter=';')
        rows = []
        for row in reader:
            # Формируем кортеж в том же порядке, что и в таблице
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
        try:
            insert_query = '''
            INSERT INTO devices (
                "тип", "название", "Модель", "Топливо", "Давление, атм", 
                "Паропроизводительность, кг/ч", "температура пара", "КПД", 
                "мощность, кВт", "Производство пара, кг/ч", "Расход газа", 
                "Расход дизеля", "Расход мазута", "Расход твердого топлива", 
                "вес, кг", "Расход твердого топлива3", "Расход твердого топлива4"
            ) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?);
            '''
            cursor.executemany(insert_query, rows)
            conn.commit()
            imported = cursor.rowcount
            total_imported += imported
            print(f'Импортировано записей: {imported}')
        except sqlite3.Error as e:
            print(f'Ошибка при импорте из файла {csv_file}: {e}')

print(f'Всего импортировано записей: {total_imported}')

conn.close()
