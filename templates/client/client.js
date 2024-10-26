
async function removeRequest(id, type) {
    try {
        const response = await fetch(`/client/${id}?${type}`, {
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

async function saveRequest(codeInputName, codeInputDoRequest) {
    const nameInput = document.getElementById(codeInputName);
    const doRequestInput = document.getElementById(codeInputDoRequest);

    let name = document.getElementById(codeInputName).value;
    if (!name || name == "") {
        name = prompt("Request name:");
    }

    if (name !== null) {
        nameInput.value = name;
        let auxDoRequest = doRequestInput.value;
        doRequestInput.value = false;
        document.getElementById("client-form").submit();
        doRequestInput.value = auxDoRequest;
    } else {
        alert("Save request canceled.");
    }
}