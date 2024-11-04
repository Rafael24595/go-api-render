function syncronizeBodies(event) {
    const content = event.target.value;
    for (const element of document.getElementsByClassName("body-parameter")) {
        element.value = content;
    }
}

function refreshRequestBody(element) {
    let target = element.getAttribute("target")
    switch (target.toUpperCase()) {
        case "JSON":
            document.jsonEditor.refresh()
            break;
        case "TEXT":
            document.textEditor.refresh()
            break;
        default:
            console.log(`Type ${target} not recognized.`)
            break;
    }
}

function refreshBodies() {
    document.jsonEditor.refresh()
    document.textEditor.refresh()
}