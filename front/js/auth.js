async function submitUsername() {
    const usernameInput = document.getElementById('usernameInput').value;
    const passwordInput = document.getElementById('passwordInput').value;
    try {
        const response = await fetch("http://localhost:3000/users", {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ username: usernameInput, password: passwordInput })
        });

        if (response.ok) {
            const data = await response.json();
            console.log('User ID:', data.user_id);
            sessionStorage.setItem('user_id', data.user_id);
            window.location.href = 'food'; 
        } else {
            const errorData = await response.json();
            console.error('Error:', response.statusText);
            alert(`Ошибка: ${errorData.message || response.statusText}`);
        }
    } catch (error) {
        console.error('Fetch error:', error);
        alert('Произошла ошибка при отправке данных. Пожалуйста, попробуйте снова.');
    }
}

document.addEventListener("DOMContentLoaded", function() {
    const loginButton = document.querySelector(".submit__btn");

    loginButton.addEventListener("click", function() {
        const usernameInput = document.getElementById("usernameInput");
        const passwordInput = document.getElementById("passwordInput");

        const username = usernameInput.value;
        const password = passwordInput.value;
        if (username && password) {
            submitUsername();
        } else {
            alert("Пожалуйста, введите имя пользователя и пароль");
        }
    });
});
