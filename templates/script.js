window.onloadFuncs = !window.onloadFuncs ? {} : window.onloadFuncs;

window.onload = () => {
    for (const key of Object.keys(window.onloadFuncs)) {
        window.onloadFuncs[key]()
    }
}
