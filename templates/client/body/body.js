const BODY_FORM = "body-type-form";
const BODY_FRAGMENTS = "body-parameter";

function showBodyForm(form) {
    for (const element of document.getElementById(BODY_FORM).children) {
        if (element.id == form) {
            element.classList.add("show");
            continue;
        }
        element.classList.remove("show");
    }
}

function syncronizeBodies(event) {
    const content = event.target.value;
    for (const element of document.getElementsByClassName(BODY_FRAGMENTS)) {
        element.value = content;
    }
}