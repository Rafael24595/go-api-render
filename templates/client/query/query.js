const QUERY_TEMPLATE_NAME = "query-parameter-template";
const QUERY_REMOVE_BUTTON = "query-remove-button";

function newQueryRow() {
    const template = document.getElementById(QUERY_TEMPLATE_NAME)
    template.insertAdjacentElement("afterend", template.cloneNode(true))
    template.id = ""

    const uuid = uuidv4()
    for(const label of [...template.getElementsByTagName("label")]) {
        label.setAttribute("for", `${label.getAttribute("for")}-${uuid}`)
        for(const input of label.getElementsByTagName("input")) {
            input.name = `${input.name}-${uuid}`;
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