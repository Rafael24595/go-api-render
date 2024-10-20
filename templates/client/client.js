const CLIENT_FORM_OPTIONS = "client-form-options";
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
    for (const element of document.getElementById(CLIENT_FORM_OPTIONS).children) {
        if (element.id == form) {
            element.classList.add("show");
            const refresh = element.getAttribute("refresh")
            if(refresh) {
                window[refresh]();
            }
            continue;
        }
        element.classList.remove("show");
    }
}

async function removeRequest(id) {
    try {
        const response = await fetch(`/client/${id}`, {
            method: 'DELETE',
        });
        
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        const newHtml = await response.text();

        document.body.innerHTML = newHtml;
    } catch (error) {
        console.error('Error sending DELETE request:', error);
    }
}