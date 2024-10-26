function showForm(parent, form) {
    for (const element of document.getElementById(parent).children) {
        if (element.id == form) {
            element.classList.add("show");
            const refresh = element.getAttribute("refresh")
            if (refresh) {
                window[refresh]();
            }
            continue;
        }
        element.classList.remove("show");
    }
}
