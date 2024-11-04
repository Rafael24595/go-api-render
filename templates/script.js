function goto(base, { paths, queries, clean }) {
    const urlTarget = applyContext(base, {paths, queries, clean});
    window.location = urlTarget.toString()
}

function showForm(event, parent, form) {
    for (const element of document.getElementById(parent).children) {
        if (element.id == form) {
            element.classList.add("show");

            const refresh = element.getAttribute("refresh")
            if (refresh) {
                window[refresh](element);
            }

            const input = event.target.control
            const name = input.name
            const value = input.value
            updateContext(name, value)
            continue;
        }
        element.classList.remove("show");
    }
}

async function submitForm(event, { paths, queries, clean }) {
    event.preventDefault();
    
    try {
        const form = event.target;
        const method = form.getAttribute("method");
        const url = new URL(form.action);
        const action = applyContext(url.pathname, {paths, queries, clean});

        history.pushState({}, "", action);

        const formData = new FormData(form);

        const response = await fetch(action, {
            method: method.toUpperCase(),
            body: formData
        });

        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        const newHtml = await response.text();

        document.body.innerHTML = newHtml;
 
        reloadScripts();

        //document.dispatchEvent(new Event('DOMContentLoaded'));
    } catch (error) {
        console.error('Error sending request:', error);
    }
}

function reloadScripts() {
    document.body.querySelectorAll("script").forEach(script => {
        const newScript = document.createElement("script");
        if (script.src) {
            newScript.src = script.src;
        } else {
            newScript.textContent = script.textContent;
        }
        document.body.appendChild(newScript);
    });
}

function applyContext(base, { paths, queries, clean }) {
    const urlOrigin = new URL(window.location);

    if (!paths) {
        paths = [];
    }

    if (!queries) {
        queries = {};
    }

    paths = [base, ...paths];

    const urlTarget = new URL(paths.join("/"), urlOrigin.origin);

    if (urlOrigin.pathname.startsWith(base) && !clean) {
        urlTarget.search = urlOrigin.searchParams.toString()
    }

    for (const key of Object.keys(queries)) {
        const value = queries[key];
        if(!value || value == "") {
            continue;
        }
        urlTarget.searchParams.set(key, value)
    }

    return urlTarget.toString()
}

function updateContext(name, value) {
    const url = new URL(window.location);
    url.searchParams.set(name, value);
    history.pushState({}, "", url);
}
