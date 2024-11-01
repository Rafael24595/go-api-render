function goto(base, {paths, queries, clean}) {
    const urlOrigin = new URL(window.location);
    
    if(!paths) {
        paths = [];
    }

    if(!queries) {
        queries = {};
    }

    paths = [base, ...paths];

    const urlTarget = new URL(paths.join("/"), urlOrigin.origin);

    if(urlOrigin.pathname.startsWith(base) && !clean) {
        urlTarget.search = urlOrigin.searchParams.toString()
    }

    for (const key of Object.keys(queries)) {
        urlTarget.searchParams.set(key, queries[key])
    }

    window.location = urlTarget.toString()
}

function showForm(event, parent, form) {
    for (const element of document.getElementById(parent).children) {
        if (element.id == form) {
            element.classList.add("show");
            
            const refresh = element.getAttribute("refresh")
            if (refresh) {
                window[refresh]();
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

function updateContext(name, value) {
    const url = new URL(window.location);
    url.searchParams.set(name, value);
    history.pushState({}, "", url);
  }