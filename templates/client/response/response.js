const RESPONSE_PAYLOAD = "response-payload";

function showContainer(form) {
    for (const element of document.getElementById(RESPONSE_PAYLOAD).children) {
        if (element.id == form) {
            element.classList.add("show");
            continue;
        }
        element.classList.remove("show");
    }
}