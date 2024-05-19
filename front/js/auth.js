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
            localStorage.setItem('user_id', data.user_id);
            window.location.href = 'food';
        } else {
            console.error('Error:', response.statusText);
        }
    } catch (error) {
        console.error('Fetch error:', error);
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
