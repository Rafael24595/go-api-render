async function updateRequest(id, name) {
    name = prompt(`New request name for ${name}:`);

    if (name !== null) {
        // TODO: Request.
    } else {
        alert("Save request canceled.");
    }
}