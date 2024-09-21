function newQueryRow(event) {
    const template = document.getElementById("query-parameter-template")
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
    
    for(const button of template.getElementsByClassName("remove-query-button")) {
        button.classList.add("show")
        button.disabled = false;
    }
}

function removeQueryRow(event) {
    const parent = event.target.parentElement;
    if(parent.id != "query-parameter-template") {
        parent.remove();
    }
}