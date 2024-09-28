const CLIENT_FORM = "client-form";

function showForm(form) {
    for (const element of document.getElementById(CLIENT_FORM).children) {
        if (element.id == form) {
            element.classList.add("show");
            continue;
        }
        element.classList.remove("show");
    }
}