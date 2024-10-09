const RESPONSE_CONTAINER = "response-container";

function showContainer(form) {
    for (const element of document.getElementById(RESPONSE_CONTAINER).children) {
        if (element.id == form) {
            element.classList.add("show");
            continue;
        }
        element.classList.remove("show");
    }
}