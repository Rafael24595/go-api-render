const QUERY_TEMPLATE_NAME = "query-parameter-template";
const QUERY_REMOVE_BUTTON = "query-remove-button";
const QUERY_ID_SEPARATOR = "#";

function newQueryRow() {
    const template = document.getElementById(QUERY_TEMPLATE_NAME)
    template.insertAdjacentElement("afterend", template.cloneNode(true))
    template.id = ""

    const uuid = uuidv4()
    for(const label of [...template.getElementsByTagName("label")]) {
        label.setAttribute("for", `${label.getAttribute("for")}${QUERY_ID_SEPARATOR}${uuid}`)
        for(const input of label.getElementsByTagName("input")) {
            if(input.type == "checkbox") {
                input.checked = true;
                input.disabled = false;
            }
            input.name = `${input.name}${QUERY_ID_SEPARATOR}${uuid}`;
            input.onkeydown = undefined;
        }
    }
    
    for(const button of template.getElementsByClassName(QUERY_REMOVE_BUTTON)) {
        button.classList.add("show")
        button.disabled = false;
    }
}

function removeQueryRow(event) {
    const parent = event.target.parentElement;
    if(parent.id != QUERY_TEMPLATE_NAME) {
        parent.remove();
    }
}