const CLIENT_FORM = "client-form";
const CLIENT_TAG = "client-tag";

window.onloadFuncs["CLIENT_TAG"] = fixClientLabel;

function fixClientLabel() {
    for (const element of document.getElementsByClassName(CLIENT_TAG)) {
        if (element.checked) {
            const label = document.getElementById(`client-label-${element.id}`);
            label.click();
        }
    }
}

function showForm(form) {
    for (const element of document.getElementById(CLIENT_FORM).children) {
        if (element.id == form) {
            element.classList.add("show");
            continue;
        }
        element.classList.remove("show");
    }
}