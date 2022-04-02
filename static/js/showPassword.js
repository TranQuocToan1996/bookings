function showPassword(id) {
	var inputPassword = document.getElementById(id);
	if (inputPassword.type === "password") {
		inputPassword.type = "text";
	} else {
		inputPassword.type = "password";
	}
}
