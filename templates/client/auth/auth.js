function showAuthForm(form) {
    for (const element of document.getElementById("auth-type-form").children) {
        if (element.id == form) {
            element.classList.add("show");
            continue;
        }
        element.classList.remove("show");
    }
}