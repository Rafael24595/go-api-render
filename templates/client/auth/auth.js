const AUTH_FORM = "auth-type-form";

function showAuthForm(form) {
    for (const element of document.getElementById(AUTH_FORM).children) {
        if (element.id == form) {
            element.classList.add("show");
            continue;
        }
        element.classList.remove("show");
    }
}