import csv
import sqlite3
import glob
import os

db_path = 'mydatabase.db'
csv_directory = './devices'

# Открываем соединение с базой
conn = sqlite3.connect(db_path)
cursor = conn.cursor()

# Создаем таблицу, если её ещё нет.
create_table_query = '''
CREATE TABLE IF NOT EXISTS devices (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    type TEXT,
    name TEXT,
    model TEXT,
    fuel TEXT,
    pressure TEXT,
    steam_capacity TEXT,
    steam_temperature TEXT,
    efficiency TEXT,
    power TEXT,
    steam_production TEXT,
    gas_cons TEXT,
    diesel_cons TEXT, 
    fuel_oil_cons TEXT,
    solid_fuel_cons TEXT,
    weight INTEGER,
    burner TEXT,
    modification_one_pump TEXT,
    modification_two_pump TEXT
);
'''
cursor.execute(create_table_query)
conn.commit()

# Функция для конвертации значений (опционально, если хотите преобразовывать числовые поля)
def convert_value(value, to_type=float):
    value = value.strip()
    if value == "":
        return None
    try:
        # Заменяем запятую на точку, чтобы можно было парсить как число
        return to_type(value.replace(',', '.'))
    except ValueError:
        return value

# Получаем список CSV-файлов в директории
csv_files = glob.glob(os.path.join(csv_directory, '*.csv'))

total_imported = 0

for csv_file in csv_files:
    print(f'Импортируем файл: {csv_file}')
    with open(csv_file, newline='', encoding='utf-8') as f:
        # Убедитесь, что разделитель соответствует вашему CSV (здесь используется точка с запятой)
        reader = csv.DictReader(f, delimiter=';')
        rows = []
        for row in reader:
            # Проверьте, что названия ключей точно совпадают с заголовками в CSV-файле.
            data = (
                row.get("тип", "").strip(),
                row.get("название", "").strip(),
                row.get("Модель", "").strip(),
                row.get("Топливо", "").strip(),
                row.get("Давление, атм", "").strip(),  # Если здесь ожидается число, можно вызвать convert_value
                row.get("Паропроизводительность, кг/ч", "").strip(),
                row.get("температура пара", "").strip(),
                row.get("КПД", "").strip(),
                row.get("мощность, кВт", "").strip(),
                row.get("Производство пара, кг/ч", "").strip(),
                row.get("Расход газа", "").strip(),
                row.get("Расход дизеля", "").strip(),
                row.get("Расход мазута", "").strip(),
                row.get("Расход твердого топлива", "").strip(),
                row.get("вес, кг", "").strip(),
                row.get("Горелка", "").strip(),
                row.get("модификация с одним питательным насосом", "").strip(),
                row.get("модификация с двумя питательными насосами", "").strip()
            )
            rows.append(data)
        try:
            insert_query = '''
            INSERT INTO devices (
                type, name, model, fuel, pressure, steam_capacity, 
                steam_temperature, efficiency, power, steam_production, 
                gas_cons, diesel_cons, fuel_oil_cons, solid_fuel_cons,
                weight, burner, modification_one_pump, modification_two_pump
            ) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?);
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
