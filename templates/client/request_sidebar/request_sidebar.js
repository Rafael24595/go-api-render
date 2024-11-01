async function updateRequest(id, type, name) {
    name = prompt(`New request name for ${name}:`);

    if (name == null) {
        alert("Save request canceled.");
        return;
    }
    
    try {

        const formData = new FormData();
        formData.append('name', name);

        const response = await fetch(`/client/${id}?${type}`, {
            method: 'PUT',
            body: formData
        });

        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        const newHtml = await response.text();

        document.body.innerHTML = newHtml;
    } catch (error) {
        console.error('Error sending request:', error);
    }
}