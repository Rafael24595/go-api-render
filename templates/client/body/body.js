const BODY_FORM = "body-type-form";
const BODY_FRAGMENTS = "body-parameter";

function showBodyForm(form) {
    for (const element of document.getElementById(BODY_FORM).children) {
        if (element.id == form) {
            element.classList.add("show");
            refreshBody(element.getAttribute("refresh"))
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

function refreshBody(type) {
    switch (type.toUpperCase()) {
        case "JSON":
            jsonEditor.refresh()
            break;
        case "TEXT":
            textEditor.refresh()
        default:
            console.log(`Type ${type} not recognized.`)
            break;
    }
}

function refreshBodies() {
    jsonEditor.refresh()
    textEditor.refresh()
}