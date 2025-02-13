var token = "";

const submitFormCallback = async (event) => {
    event.preventDefault(); // Prevent default form submission

    const formData = new FormData(event.target);
    const response = await fetch(event.target.action, {
        method: event.target.method,
        headers: {
            'Authorization': `${token}`
        },
        body: formData
    });

    try {
        const result = await response.json();
        console.log(result);
        localStorage.setItem("token", result.token);
        localStorage.setItem("authority", result.authority);
        token = result.token;

        await Dashboard(result.authority);
    } catch (error) {
        console.error(error);
    }
    // Handle the response (e.g., navigate to another page, show a message, etc.)
};

function HomePage() {
    console.log("Home page");

    const mainElement = document.getElementById("main");
    mainElement.className = "";
    mainElement.className += "centered-items-in-column";

    mainElement.innerHTML = "";

    const titleElement = document.createElement("h1");
    titleElement.textContent = "Home";

    mainElement.appendChild(titleElement);

	const wrapperElement = document.createElement("div")
	wrapperElement.className = "home-btn-wrapper";

    const loginButton = document.createElement("button");
    loginButton.textContent = "Login";
    loginButton.addEventListener("click", () => {
        LoginPage();
    });

    wrapperElement.appendChild(loginButton);

    const registerButton = document.createElement("button");
    registerButton.textContent = "Register";
    registerButton.addEventListener("click", () => {
        RegisterPage();
    });

    wrapperElement.appendChild(registerButton);
	
	mainElement.appendChild(wrapperElement)
}

async function Dashboard(userAuthority) {
    const mainElement = document.getElementById("main");
    mainElement.className = "";
    mainElement.className += "centered-items-in-column";

    mainElement.innerHTML = "";

    const titleWrapper = document.createElement("div");
    titleWrapper.className = "admin-title-wrapper";

    const titleElement = document.createElement("h1");
    titleElement.textContent = "Dashboard";

    titleWrapper.appendChild(titleElement);

    const userElement = document.createElement("p");
    userElement.textContent = `Welcome to dashboard!`;

    titleWrapper.appendChild(userElement);

    const logoutButton = document.createElement("button");
    logoutButton.textContent = "Logout";
    logoutButton.addEventListener("click", async () => {
        try {
            const response = await fetch("/api/logout", {
                headers: {
                    "Authorization": token
                }
            });

            const result = await response.json();
            console.log(result);

            localStorage.removeItem("token");
            localStorage.removeItem("authority");
            token = "";
            HomePage();

        } catch (error) {
            console.error(error);
        }
    });

    titleWrapper.appendChild(logoutButton);

    mainElement.appendChild(titleWrapper);

    try {
        console.log("User authority fetching...");
        const response = await fetch("/api/authority", {
            headers: {
                "Authority": userAuthority,
                "Authorization": token
            }
        });

        const result = await response.json();
        console.log(result);

        if (result.authority !== "admin") {
            return;
        }

        console.log("Users list fetching...");
        const usersData = await fetch("/api/users", {
            headers: {
                "Authority": userAuthority,
                "Authorization": token
            }
        });

        const users = await usersData.json();
        console.log(users);

        const usersListContainer = document.createElement("div");
        usersListContainer.className = "users-list";

        const usersShortListContainer = document.createElement("div");
        usersShortListContainer.className = "users-short-list";

        users.forEach(user => {
            const userElement = document.createElement("div");
            userElement.innerHTML = `<p>${user.username}</p>`;
            userElement.addEventListener("click", async () => {
                const userDetailedInfoContainer = document.getElementById("user-detailed-info");
                userDetailedInfoContainer.innerHTML = "";

                const editForm = document.createElement("form");
                editForm.innerHTML = `
                    <label>Username: <input type="text" name="username" value="${user.username}"></label>
                    <label>Email: <input type="email" name="email" value="${user.email}"></label>
                    <label>Phone: <input type="tel" name="phone" value="${user.phone}"></label>
                    <label>Country: <input type="text" name="country" value="${user.contry}"></label>
                    <label>Birthdate: <input type="date" name="bday" value="${user.bday}"></label>
                    <label>Gender: <select name="gender">
                        <option value="male">Male</option>
                        <option value="female">Female</option>
                        <option value="other">Other</option>
                    </select></label>
                    <button type="submit">Save</button>
                `;

                editForm.querySelector(`select[name="gender"]`).value = user.gender;

                editForm.addEventListener("submit", async (event) => {
                    event.preventDefault();
                    const formData = new FormData(editForm);
                    const updatedUser = {
                        id: user.id,
                        username: formData.get("username"),
                        email: formData.get("email"),
                        phone: formData.get("phone"),
                        country: formData.get("country"),
                        bday: formData.get("bday")
                    };

                    await updateUser(updatedUser);
                    userDetailedInfoContainer.innerHTML = "";
                    Dashboard("admin");
                });

                userDetailedInfoContainer.appendChild(editForm);
            });

            usersShortListContainer.appendChild(userElement);
        });

        usersListContainer.appendChild(usersShortListContainer);

        const userDetailedInfoContainer = document.createElement("div");
        userDetailedInfoContainer.id = "user-detailed-info";
        userDetailedInfoContainer.className = "centered-items-in-column user-detailed-info";

        usersListContainer.appendChild(userDetailedInfoContainer)

        mainElement.appendChild(usersListContainer);
    } catch (error) {
        console.error(error);
    }
}

async function updateUser(user) {
    try {
        const response = await fetch(`/api/users/${user.id}`, {
            method: "put",
            headers: {
                "Authorization": token,
                "Content-Type": "application/json"
            },
            body: JSON.stringify(user)
        });

        if (!response.ok) {
            throw new Error("Failed to update user");
        }

        console.log("User updated successfully");
    } catch (error) {
        console.error(error);
    }
}

function LoginPage() {
    const mainElement = document.getElementById("main");
    mainElement.className = "";
    mainElement.className += "centered-items-in-column";

    mainElement.innerHTML = "";

    const loginForm = document.createElement("form");
    loginForm.setAttribute("id", "login-form");
    loginForm.setAttribute("method", "post");
    loginForm.setAttribute("action", "/api/login");

    const titleElement = document.createElement("h1");
    titleElement.textContent = "Login";

    loginForm.appendChild(titleElement);

    const emailInput = document.createElement("input");
    emailInput.setAttribute("type", "email");
    emailInput.setAttribute("name", "email");
    emailInput.setAttribute("placeholder", "Email");

    loginForm.appendChild(emailInput);

    const passwordInput = document.createElement("input");
    passwordInput.setAttribute("type", "password");
    passwordInput.setAttribute("name", "password");
    passwordInput.setAttribute("placeholder", "Password");

    loginForm.appendChild(passwordInput);

    const loginButton = document.createElement("button");
    loginButton.setAttribute("type", "submit");
    loginButton.textContent = "Login";

    loginForm.appendChild(loginButton);

    const backButton = document.createElement("button");
    backButton.setAttribute("type", "button");
    backButton.textContent = "Back";
    backButton.addEventListener("click", () => {
        HomePage();
    });

    loginForm.appendChild(backButton);

    mainElement.appendChild(loginForm);

    // Handle form submission
    loginForm.addEventListener("submit", submitFormCallback);
}

function RegisterPage() {
    const mainElement = document.getElementById("main");
    mainElement.className = "";
    mainElement.className += "centered-items-in-column";

    mainElement.innerHTML = "";

    const registerForm = document.createElement("form");
    registerForm.setAttribute("id", "register-form");
    registerForm.setAttribute("method", "post");
    registerForm.setAttribute("action", "/api/register");

    const titleElement = document.createElement("h1");
    titleElement.textContent = "Register";

    registerForm.appendChild(titleElement);

    const emailInput = document.createElement("input");
    emailInput.setAttribute("type", "email");
    emailInput.setAttribute("name", "email");
    emailInput.setAttribute("placeholder", "Email");

    registerForm.appendChild(emailInput);

    const usernameInput = document.createElement("input");
    usernameInput.setAttribute("type", "text");
    usernameInput.setAttribute("name", "username");
    usernameInput.setAttribute("placeholder", "Username");

    registerForm.appendChild(usernameInput);

    const phoneInput = document.createElement("input");
    phoneInput.setAttribute("type", "tel");
    phoneInput.setAttribute("name", "phone");
    phoneInput.setAttribute("placeholder", "Phone");

    registerForm.appendChild(phoneInput);

    const passwordInput = document.createElement("input");
    passwordInput.setAttribute("type", "password");
    passwordInput.setAttribute("name", "password");
    passwordInput.setAttribute("placeholder", "Password");

    registerForm.appendChild(passwordInput);

    const confirmPasswordInput = document.createElement("input");
    confirmPasswordInput.setAttribute("type", "password");
    confirmPasswordInput.setAttribute("name", "confirmPassword");
    confirmPasswordInput.setAttribute("placeholder", "Confirm Password");
    confirmPasswordInput.onchange = () => {
        if (confirmPasswordInput.value !== passwordInput.value) {
            confirmPasswordInput.setCustomValidity("Passwords don't match");
        } else {
            confirmPasswordInput.setCustomValidity("");
        }
    }

    registerForm.appendChild(confirmPasswordInput);

    const countryInput = document.createElement("input");
    countryInput.setAttribute("type", "text");
    countryInput.setAttribute("name", "contry");
    countryInput.setAttribute("placeholder", "Country");

    registerForm.appendChild(countryInput);

    const birthdateInput = document.createElement("input");
    birthdateInput.setAttribute("type", "date");
    birthdateInput.setAttribute("name", "bday");
    birthdateInput.setAttribute("placeholder", "Birthdate");

    registerForm.appendChild(birthdateInput);

    const genderInput = document.createElement("select");
    genderInput.setAttribute("name", "gender");
    genderInput.setAttribute("required", "required");

    const optionNone = document.createElement("option");
    optionNone.setAttribute("value", "");
    optionNone.textContent = "Select Gender";
    genderInput.appendChild(optionNone);

    const optionMale = document.createElement("option");
    optionMale.setAttribute("value", "male");
    optionMale.textContent = "Male";
    genderInput.appendChild(optionMale);

    const optionFemale = document.createElement("option");
    optionFemale.setAttribute("value", "female");
    optionFemale.textContent = "Female";
    genderInput.appendChild(optionFemale);

    const optionOther = document.createElement("option");
    optionOther.setAttribute("value", "other");
    optionOther.textContent = "Other";
    genderInput.appendChild(optionOther);

    registerForm.appendChild(genderInput);

    const photoInput = document.createElement("input");
    photoInput.setAttribute("type", "file");
    photoInput.setAttribute("name", "photo");

    registerForm.appendChild(photoInput);

    const dataTheftLabel = document.createElement("label");
    dataTheftLabel.textContent = "I agree to have my data stolen";
    dataTheftLabel.setAttribute("for", "agreement");

    registerForm.appendChild(dataTheftLabel);

    const dataTheftCheckbox = document.createElement("input");
    dataTheftCheckbox.setAttribute("type", "checkbox");
    dataTheftCheckbox.setAttribute("name", "agreement");
    dataTheftCheckbox.value = "false";

    registerForm.appendChild(dataTheftCheckbox);

    const registerButton = document.createElement("button");
    registerButton.setAttribute("type", "submit");
    registerButton.textContent = "Register";

    registerForm.appendChild(registerButton);

    const backButton = document.createElement("button");
    backButton.setAttribute("type", "button");
    backButton.textContent = "Back";
    backButton.addEventListener("click", async () => {
        HomePage();
    });

    registerForm.appendChild(backButton);

    mainElement.appendChild(registerForm);

    // Handle form submission
    registerForm.addEventListener("submit", submitFormCallback);
}

document.addEventListener("DOMContentLoaded", function () {
    HomePage();
});
