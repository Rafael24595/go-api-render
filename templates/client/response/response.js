function refreshResponseBody(element) {
    let target = element.getAttribute("target")
    switch (target.toUpperCase()) {
        case "JSON":
            document.jsonViewer.refresh()
            break;
        case "TEXT":
            document.textViewer.refresh()
            break;
        case "HTML":
            document.htmlViewer.refresh()
            break;
        default:
            console.log(`Type ${target} not recognized.`)
            break;
    }
}