document.addEventListener('DOMContentLoaded', () => {
  const loginForm = document.getElementById('loginform');

  if (loginForm) {
    loginForm.addEventListener('submit', (e) => {
      e.preventDefault();

      const username = document.getElementById('username').value;
      const password = document.getElementById('password').value;

      fetch('/login', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          username: username,
          password: password,
        }),
      })
        .then((response) => {
          if (!response.ok) {
            throw new Error('Неверное имя пользователя или пароль');
          }
          return response.json();
        })
        .then((data) => {
          console.log(`Успешный вход, данные: ${data}`);
          if (data.redirect) {
            console.log('Перенаправление на:', data.redirect);
            setTimeout(() => {
              window.location.href = data.redirect;
            }, 500);
          } else {
            console.log('Ошибка: отсутствует URL для перенаправления');
          }
        })
        .catch((error) => {
          console.error('Ошибка входа:', error);
          alert(`Ошибка входа: ${error.message}`);
        });
    });
  }
});

/* <---------------------------------------------------> */

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

/* <---------------------------------------------------> */

function showCacheNotification() {
  let notification = document.getElementById('cache-notification');
  if (!notification) {
    notification = document.createElement('div');
    notification.id = 'cache-notification';
    notification.className = 'cache-notification';
    notification.textContent = 'Кэш успешно очищен';
    document.body.appendChild(notification);
  }

  setTimeout(() => {
    notification.classList.add('show');
  }, 100);

  setTimeout(() => {
    notification.classList.remove('show');
  }, 2000);
}

function clearCache() {
  localStorage.clear();
  sessionStorage.clear();

  if ('serviceWorker' in navigator) {
    navigator.serviceWorker.getRegistrations().then((registrations) => {
      for (const registration of registrations) {
        registration.unregister();
        console.log('Service Worker отменен');
      }
    });
  }

  const timestamp = new Date().getTime();

  const styleSheets = document.querySelectorAll('link[rel="stylesheet"]');
  for (const link of styleSheets) {
    const url = link.href.split('?')[0];
    link.href = `${url}?v=${timestamp}`;
  }

  const scripts = document.querySelectorAll('script[src]');
  for (const script of scripts) {
    const url = script.src.split('?')[0];
    script.src = `${url}?v=${timestamp}`;
  }

  const images = document.querySelectorAll('img[src]');
  for (const img of images) {
    const url = img.src.split('?')[0];
    img.src = `${url}?v=${timestamp}`;
  }

  showCacheNotification();

  setTimeout(() => {
    window.location.reload(true);
  }, 1500);
}
