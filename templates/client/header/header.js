
function newHeaderRow() {
    const HEADER_ID_SEPARATOR = "#";
    
    const template = document.getElementById("header-parameter-template")
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
    
    for(const button of template.getElementsByClassName("key-value-form-remove-button")) {
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