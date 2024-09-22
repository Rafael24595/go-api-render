const CLIENT_FORM = "client-form";

function showForm(form) {
    for (const element of document.getElementById(CLIENT_FORM).children) {
        console.log(element.id, form)
        if (element.id == form) {
            element.classList.add("show");
            continue;
        }
        element.classList.remove("show");
    }
}