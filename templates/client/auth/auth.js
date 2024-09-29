const AUTH_FORM = "auth-type-form";
const AUTH_TAG = "auth-tag";

window.onloadFuncs["AUTH_TAG"] = fixAuthLabel;

function fixAuthLabel() {
    for (const element of document.getElementsByClassName(AUTH_TAG)) {
        if (element.checked) {
            const label = document.getElementById(`auth-label-${element.id}`);
            label.click();
        }
    }
}

function showAuthForm(form) {
    for (const element of document.getElementById(AUTH_FORM).children) {
        if (element.id == form) {
            element.classList.add("show");
            continue;
        }
        element.classList.remove("show");
    }
}