// Переключение темы
const themeToggle = document.getElementById('theme-icon');
const container = document.querySelector('.container');

themeToggle.addEventListener('click', () => {
  document.body.classList.toggle('dark-mode');
  container.classList.toggle('dark-mode');

  if (document.body.classList.contains('dark-mode')) {
    themeToggle.innerHTML = `<path d="M12 8a4 4 0 1 1-8 0 4 4 0 0 1 8 0..."/>`;
  } else {
    themeToggle.innerHTML = `<path d="M6 .278a.77.77 0 0 1 .08.858..."/>`;
  }
});

async function loadDevices() {
  try {
    const response = await fetch('/devices');
    const devices = await response.json();
    populateTable(devices);
  } catch (error) {
    console.error('Ошибка загрузки девайсов:', error);
  }
}

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
      <td>${dev.mtp || ''}</td>
      <td><input type="checkbox" data-id="${dev.id}" /></td>
    `;
    tbody.appendChild(tr);
  }
}
const modal = document.getElementById('modal');
const closeModal = document.getElementById('closeModal');
const recordForm = document.getElementById('recordForm');
const modalTitle = document.getElementById('modalTitle');

const addBtn = document.getElementById('addBtn');
const editBtn = document.getElementById('editBtn');
const deleteBtn = document.getElementById('deleteBtn');

addBtn.addEventListener('click', () => {
  modalTitle.textContent = 'Добавить запись';
  recordForm.reset();
  recordForm.elements.id.value = '';
  modal.style.display = 'block';
});

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

    recordForm.elements.id.value = dev.id;
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
    recordForm.elements.mpt.value = dev.mtp || '';

    modalTitle.textContent = 'Редактировать запись';
    modal.style.display = 'block';
  } catch (err) {
    console.error('Ошибка получения данных для редактирования:', err);
    alert('Ошибка при получении данных для редактирования');
  }
});

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

recordForm.addEventListener('submit', async (e) => {
  e.preventDefault();
  const formData = new FormData(recordForm);
  const record = {};
  formData.forEach((val, key) => {
    record[key] = val.trim();
  });
  console.log('Отправляем record:', record);
  try {
    if (record.id) {
      // Обновление (PUT)
      await fetch(`/devices/${record.id}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(record),
      });
    } else {
      // Добавление (POST)
      await fetch('/devices', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(record),
      });
    }
    modal.style.display = 'none';
    loadDevices();
  } catch (err) {
    console.error('Ошибка при сохранении записи:', err);
    alert('Ошибка при сохранении записи');
  }
});

closeModal.addEventListener('click', () => {
  modal.style.display = 'none';
});

window.addEventListener('click', (e) => {
  if (e.target === modal) {
    modal.style.display = 'none';
  }
});

loadDevices();
