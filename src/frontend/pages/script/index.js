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
