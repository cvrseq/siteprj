const themeToggle = document.getElementById('theme-icon');
const container = document.querySelector('.container');

themeToggle.addEventListener('click', () => {
  document.body.classList.toggle('dark-mode');
  container.classList.toggle('dark-mode');

  if (document.body.classList.contains('dark-mode')) {
    themeToggle.innerHTML = `<path d="M12 8a4 4 0 1 1-8 0 4 4 0 0 1 8 0M8 0a.5.5 0 0 1 .5.5v2a.5.5 0 0 1-1 0v-2A.5.5 0 0 1 8 0m0 13a.5.5 0 0 1 .5.5v2a.5.5 0 0 1-1 0v-2A.5.5 0 0 1 8 13m8-5a.5.5 0 0 1-.5.5h-2a.5.5 0 0 1 0-1h2a.5.5 0 0 1 .5.5M3 8a.5.5 0 0 1-.5.5h-2a.5.5 0 0 1 0-1h2A.5.5 0 0 1 3 8m10.657-5.657a.5.5 0 0 1 0 .707l-1.414 1.415a.5.5 0 1 1-.707-.708l1.414-1.414a.5.5 0 0 1 .707 0m-9.193 9.193a.5.5 0 0 1 0 .707L3.05 13.657a.5.5 0 0 1-.707-.707l1.414-1.414a.5.5 0 0 1 .707 0m9.193 2.121a.5.5 0 0 1 0 .707l-1.414-1.414a.5.5 0 0 1 .707-.707l1.414 1.414a.5.5 0 0 1 0 .707M4.464 4.465a.5.5 0 0 1-.707 0L2.343 3.05a.5.5 0 1 1 .707-.707l1.414 1.414a.5.5 0 0 1 0 .708"/>`;
  } else {
    themeToggle.innerHTML = `<path d="M6 .278a.77.77 0 0 1 .08.858 7.2 7.2 0 0 0-.878 3.46c0 4.021 3.278 7.277 7.318 7.277q.792-.001 1.533-.16a.79.79 0 0 1 .81.316.73.73 0 0 1-.031.893A8.35 8.35 0 0 1 8.344 16C3.734 16 0 12.286 0 7.71 0 4.266 2.114 1.312 5.124.06A.75.75 0 0 1 6 .278"/>`;
  }
});

/* <----------------------------------------------------------------------------------------------------------> */

const modal = document.getElementById('modal');
const closeModal = document.getElementById('closeModal');
const recordForm = document.getElementById('recordForm');
const modalTitle = document.getElementById('modalTitle');
const addBtn = document.getElementById('addBtn');
const editBtn = document.getElementById('editBtn');
const deleteBtn = document.getElementById('deleteBtn');

async function loadEmployees() {
  try {
    const response = await fetch('/employees');

    if (!response.ok) {
      throw new Error(`Ошибка HTTP: ${response.status}`);
    }

    const employees = await response.json();
    populateTable(employees);
  } catch (error) {
    console.error('Ошибка загрузки сотрудников:', error);
    showNotification('Не удалось загрузить данные сотрудников', 'error');
  }
}

/**
 *
 * @param {Array} data
 */
function populateTable(data) {
  const tbody = document.querySelector('#data-table tbody');
  tbody.innerHTML = '';

  if (!data || data.length === 0) {
    const tr = document.createElement('tr');
    tr.innerHTML = '<td colspan="5" class="text-center">Нет данных</td>';
    tbody.appendChild(tr);
    return;
  }

  for (const emp of data) {
    const tr = document.createElement('tr');
    tr.innerHTML = `
      <td>${emp.id || ''}</td>
      <td>${emp.username || ''}</td>
      <td>${emp.password || ''}</td>
      <td>${emp.role || ''}</td>
      <td><input type="checkbox" data-id="${emp.id}" /></td>
    `;
    tbody.appendChild(tr);
  }
}

/**
 *
 * @param {string} message
 * @param {string} type
 */
function showNotification(message, type = 'info') {
  alert(`${type.toUpperCase()}: ${message}`);
}

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
    showNotification(
      'Пожалуйста, выберите одну запись для редактирования',
      'warning'
    );
    return;
  }

  const id = selected[0].dataset.id;

  try {
    const response = await fetch(`/employees/${id}`);

    if (!response.ok) {
      const errorText = await response.text();
      throw new Error(`Ошибка сервера: ${response.status} - ${errorText}`);
    }

    const emp = await response.json();
    console.log('Данные для редактирования:', emp);

    recordForm.elements.id.value = emp.id || '';
    recordForm.elements.username.value = emp.username || '';
    recordForm.elements.password.value = emp.password || '';
    recordForm.elements.role.value = emp.role || '';

    modalTitle.textContent = 'Редактировать запись';
    modal.style.display = 'block';
  } catch (err) {
    console.error('Ошибка при получении данных:', err);
    showNotification(
      `Ошибка при получении данных для редактирования: ${err.message}`,
      'error'
    );
  }
});

deleteBtn.addEventListener('click', async () => {
  const selected = document.querySelectorAll(
    '#data-table tbody input[type="checkbox"]:checked'
  );

  if (selected.length === 0) {
    showNotification('Пожалуйста, выберите записи для удаления', 'warning');
    return;
  }

  if (
    !confirm(`Вы уверены, что хотите удалить ${selected.length} запись(ей)?`)
  ) {
    return;
  }

  let hasErrors = false;

  for (const checkbox of selected) {
    const id = checkbox.dataset.id;
    try {
      const response = await fetch(`/employees/${id}`, { method: 'DELETE' });

      if (!response.ok) {
        const errorText = await response.text();
        throw new Error(
          `Ошибка при удалении записи ${id}: ${response.status} - ${errorText}`
        );
      }
    } catch (err) {
      console.error(`Ошибка при удалении записи с id ${id}:`, err);
      hasErrors = true;
    }
  }

  await loadEmployees();

  if (hasErrors) {
    showNotification('Некоторые записи не удалось удалить', 'warning');
  } else {
    showNotification('Записи успешно удалены', 'success');
  }
});

recordForm.addEventListener('submit', async (e) => {
  e.preventDefault();

  const formData = new FormData(recordForm);
  const formEntries = {};

  formData.forEach((val, key) => {
    if (key === 'id' && val.trim() !== '') {
      formEntries[key] = Number.parseInt(val.trim(), 10);
    } else {
      formEntries[key] = val.trim();
    }
  });

  try {
    let response;

    if (formEntries.id) {
      response = await fetch(`/employees/${formEntries.id}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(formEntries),
      });
    } else {
      const { id, ...recordWithoutId } = formEntries;
      console.log(
        'Отправляемые данные (создание):',
        JSON.stringify(recordWithoutId)
      );

      response = await fetch('/employees', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(recordWithoutId),
      });
    }

    if (!response.ok) {
      const errorText = await response.text();
      throw new Error(`Ошибка сервера: ${response.status} - ${errorText}`);
    }

    modal.style.display = 'none';
    await loadEmployees();
    showNotification('Запись успешно сохранена', 'success');
  } catch (err) {
    console.error('Ошибка при сохранении записи:', err);
    showNotification(`Ошибка при сохранении записи: ${err.message}`, 'error');
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

document.addEventListener('DOMContentLoaded', () => {
  loadEmployees();
});
