const VARIABLE_TEMPLATE_NAME = "collection-variable-template";
const VARIABLE_REMOVE_BUTTON = "key-value-form-remove-button";
const VARIABLE_ID_SEPARATOR = "#";

function newVariableRow() {
    const template = document.getElementById(VARIABLE_TEMPLATE_NAME)
    template.insertAdjacentElement("afterend", template.cloneNode(true))
    template.id = ""

    const uuid = uuidv4()
    for(const label of [...template.getElementsByTagName("label")]) {
        label.setAttribute("for", `${label.getAttribute("for")}${VARIABLE_ID_SEPARATOR}${uuid}`)
        for(const input of label.getElementsByTagName("input")) {
            if(input.type == "checkbox") {
                input.checked = true;
                input.disabled = false;
            }
            input.name = `${input.name}${VARIABLE_ID_SEPARATOR}${uuid}`;
            input.onkeydown = undefined;
        }
    }
    
    for(const button of template.getElementsByClassName(VARIABLE_REMOVE_BUTTON)) {
        button.classList.add("show")
        button.disabled = false;
    }
}

function removeVariableRow(event) {
    const parent = event.target.parentElement;
    if(parent.id != VARIABLE_TEMPLATE_NAME) {
        parent.remove();
    }
}