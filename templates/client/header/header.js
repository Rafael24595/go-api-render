const HEADER_TEMPLATE_NAME = "header-parameter-template";
const HEADER_REMOVE_BUTTON = "key-value-form-remove-button";
const HEADER_ID_SEPARATOR = "#";

function newHeaderRow() {
    const template = document.getElementById(HEADER_TEMPLATE_NAME)
    template.insertAdjacentElement("afterend", template.cloneNode(true))
    template.id = ""

    const uuid = uuidv4()
    for(const label of [...template.getElementsByTagName("label")]) {
        label.setAttribute("for", `${label.getAttribute("for")}${HEADER_ID_SEPARATOR}${uuid}`)
        for(const input of label.getElementsByTagName("input")) {
            if(input.type == "checkbox") {
                input.checked = true;
                input.disabled = false;
            }
            input.name = `${input.name}${HEADER_ID_SEPARATOR}${uuid}`;
            input.onkeydown = undefined;
        }
    }
    
    for(const button of template.getElementsByClassName(HEADER_REMOVE_BUTTON)) {
        button.classList.add("show")
        button.disabled = false;
    }
}

function removeHeaderRow(event) {
    const parent = event.target.parentElement;
    if(parent.id != HEADER_TEMPLATE_NAME) {
        parent.remove();
    }
}