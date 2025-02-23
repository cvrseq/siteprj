import csv
import sqlite3
import glob
import os

db_path = 'mydatabase.db'
csv_directory = './devices'  

conn = sqlite3.connect(db_path)
cursor = conn.cursor()

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
    weight TEXT,
    burner TEXT,
    modification_one_pump TEXT,
    modification_two_pump TEXT
);
'''
cursor.execute(create_table_query)
conn.commit()

def convert_value(value, to_type=float):
    value = value.strip()
    if value == "":
        return None
    try:
        return to_type(value.replace(',', '.'))
    except ValueError:
        return value

csv_files = glob.glob(os.path.join(csv_directory, '*.csv'))

total_imported = 0

for csv_file in csv_files:
    print(f'Импортируем файл: {csv_file}')
    with open(csv_file, newline='', encoding='utf-8') as f:
        reader = csv.DictReader(f, delimiter=';')
        rows = []
        for row in reader:
            data = (
                row.get("type", "").strip(),
                row.get("name", "").strip(),
                row.get("model", "").strip(),
                row.get("fuel", "").strip(),
                row.get("pressure", "").strip(),
                row.get("steam_capacity", "").strip(),
                row.get("steam_temperature", "").strip(),
                row.get("efficiency", "").strip(),
                row.get("power", "").strip(),
                row.get("steam_production", "").strip(),
                row.get("gas_cons", "").strip(),
                row.get("diesel_cons", "").strip(),
                row.get("fuel_oil_cons", "").strip(),
                row.get("solid_fuel_cons", "").strip(),
                row.get("weight", "").strip(),
                row.get("burner", "").strip(),
                row.get("modification_one_pump", "").strip(),
                row.get("modification_two_pump", "").strip()
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
