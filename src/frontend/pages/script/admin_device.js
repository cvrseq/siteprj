// Переключение темы
const themeToggle = document.getElementById('theme-icon');
const container = document.querySelector('.container');

themeToggle.addEventListener('click', () => {
  document.body.classList.toggle('dark-mode');
  container.classList.toggle('dark-mode');
  if (document.body.classList.contains('dark-mode')) {
    themeToggle.innerHTML = `<path d="..."/>`;
  } else {
    themeToggle.innerHTML = `<path d="..."/>`;
  }
});

// Загрузка устройств при старте
async function loadDevices() {
  try {
    const response = await fetch('/devices');
    const devices = await response.json();
    populateTable(devices);
  } catch (error) {
    console.error('Ошибка загрузки девайсов:', error);
  }
}

// Заполняем таблицу устройствами
function populateTable(data) {
  const tbody = document.querySelector('#data-table tbody');
  tbody.innerHTML = '';
  for (const dev of data) {
    const tr = document.createElement('tr');
    tr.innerHTML = `
      <td>${dev.id}</td>
      <td>${dev.type || ''}</td>
      <td>${dev.name || ''}</td>
      <td>${dev.model || ''}</td>
      <td>${dev.fuel || ''}</td>
      <td>${dev.pressure || ''}</td>
      <td>${dev.steam_capacity || ''}</td>
      <td>${dev.steam_temperature || ''}</td>
      <td>${dev.efficiency || ''}</td>
      <td>${dev.power || ''}</td>
      <td>${dev.steam_production || ''}</td>
      <td>${dev.gas_cons || ''}</td>
      <td>${dev.diesel_cons || ''}</td>
      <td>${dev.fuel_oil_cons || ''}</td>
      <td>${dev.solid_fuel_cons || ''}</td>
      <td>${dev.weight || ''}</td>
      <td>${dev.burner || ''}</td>
      <td>${dev.mop || ''}</td>
      <td>${dev.mpt || ''}</td>
      <td><input type="checkbox" data-id="${dev.id}" /></td>
    `;
    tbody.appendChild(tr);
  }
}

// Селекторы для элементов
const modal = document.getElementById('modal');
const closeModal = document.getElementById('closeModal');
const recordForm = document.getElementById('recordForm');
const modalTitle = document.getElementById('modalTitle');

const addBtn = document.getElementById('addBtn');
const editBtn = document.getElementById('editBtn');
const deleteBtn = document.getElementById('deleteBtn');

// «Добавить»: открываем форму в режиме добавления
addBtn.addEventListener('click', () => {
  modalTitle.textContent = 'Добавить запись';
  recordForm.reset(); // Очищаем поля
  recordForm.elements.id.value = ''; // Ставим id пустым
  modal.style.display = 'block'; // Показываем модальное окно
});

// «Редактировать»: загружаем данные выбранной записи, заполняем форму
editBtn.addEventListener('click', async () => {
  const selected = document.querySelectorAll(
    '#data-table tbody input[type="checkbox"]:checked'
  );
  if (selected.length !== 1) {
    alert('Пожалуйста, выберите одну запись для редактирования');
    return;
  }
  const id = selected[0].dataset.id;
  try {
    const response = await fetch(`/devices/${id}`);
    if (!response.ok) {
      throw new Error('Не удалось получить данные для редактирования');
    }
    const dev = await response.json();
    console.log('Полученные данные:', dev);

    // Заполняем форму
    recordForm.elements.id.value = dev.id || '';
    recordForm.elements.type.value = dev.type || '';
    recordForm.elements.name.value = dev.name || '';
    recordForm.elements.model.value = dev.model || '';
    recordForm.elements.fuel.value = dev.fuel || '';
    recordForm.elements.pressure.value = dev.pressure || '';
    recordForm.elements.steam_capacity.value = dev.steam_capacity || '';
    recordForm.elements.steam_temperature.value = dev.steam_temperature || '';
    recordForm.elements.efficiency.value = dev.efficiency || '';
    recordForm.elements.power.value = dev.power || '';
    recordForm.elements.steam_production.value = dev.steam_production || '';
    recordForm.elements.gas_cons.value = dev.gas_cons || '';
    recordForm.elements.diesel_cons.value = dev.diesel_cons || '';
    recordForm.elements.fuel_oil_cons.value = dev.fuel_oil_cons || '';
    recordForm.elements.solid_fuel_cons.value = dev.solid_fuel_cons || '';
    recordForm.elements.weight.value = dev.weight || '';
    recordForm.elements.burner.value = dev.burner || '';
    recordForm.elements.mop.value = dev.mop || '';
    recordForm.elements.mpt.value = dev.mpt || '';

    modalTitle.textContent = 'Редактировать запись';
    modal.style.display = 'block';
  } catch (err) {
    console.error('Ошибка получения данных для редактирования:', err);
    alert('Ошибка при получении данных для редактирования');
  }
});

// «Удалить»: удаляем каждую выбранную запись
deleteBtn.addEventListener('click', async () => {
  const selected = document.querySelectorAll(
    '#data-table tbody input[type="checkbox"]:checked'
  );
  if (selected.length === 0) {
    alert('Пожалуйста, выберите записи для удаления');
    return;
  }
  if (!confirm('Вы уверены, что хотите удалить выбранные записи?')) return;

  for (const checkbox of selected) {
    const id = checkbox.dataset.id;
    try {
      await fetch(`/devices/${id}`, { method: 'DELETE' });
    } catch (err) {
      console.error(`Ошибка при удалении записи с id ${id}:`, err);
    }
  }
  loadDevices();
});

// Сохранение (добавление или редактирование)
recordForm.addEventListener('submit', async (e) => {
  e.preventDefault();
  const formData = new FormData(recordForm);
  const record = {};

  formData.forEach((val, key) => {
    record[key] = val.trim(); // Удаляем пробелы
  });

  console.log('Отправляемые данные:', JSON.stringify(record)); // Логируем отправку

  try {
    if (record.id) {
      // Режим редактирования
      await fetch(`/devices/${record.id}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(record),
      });
    } else {
      // Режим добавления
      await fetch('/devices', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(record),
      });
    }

    modal.style.display = 'none';
    loadDevices(); // Обновляем таблицу
  } catch (err) {
    console.error('Ошибка при сохранении записи:', err);
    alert('Ошибка при сохранении записи');
  }
});

// Закрытие модального окна
closeModal.addEventListener('click', () => {
  modal.style.display = 'none';
});
window.addEventListener('click', (e) => {
  if (e.target === modal) {
    modal.style.display = 'none';
  }
});

// При загрузке страницы
loadDevices();
