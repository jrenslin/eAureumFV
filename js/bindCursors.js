/* Bind keybindings to body. */
if (document.getElementsByTagName('body')[0] != null) {
    document.getElementsByTagName('body')[0].addEventListener('keydown', hotkeys);
}

/* Function for referring */
function replaceWindow (id) {
    var link = document.getElementById(id).href;
    if (link != null) {
        window.location.replace(link);
    }
}

/* Define usable keybindings. */
function hotkeys (e) {

    // Bind links with ID 0-9 to keypresses on number keys
    switch (e.keyCode) {
    case 48:
        if (e.ctrlKey) {
            replaceWindow("link0ctrl");
            break;
        }
        else {
            replaceWindow("link0");
            break;
        }
    case 49:
        replaceWindow("link1");
        break;
    case 50:
        replaceWindow("link2");
        break;
    case 51:
        replaceWindow("link3");
        break;
    case 52:
        replaceWindow("link4");
        break;
    case 53:
        replaceWindow("link5");
        break;
    case 54:
        replaceWindow("link6");
        break;
    case 55:
        replaceWindow("link7");
        break;
    case 56:
        replaceWindow("link8");
        break;
    case 57:
        replaceWindow("link9");
        break;
    }

    // Bind cursor keys (only with CTRL pressed)
    if (e.ctrlKey) {
        switch (e.keyCode) {
        case 37:
            replaceWindow("prev");
            break;
        case 38:
            replaceWindow("goUp");
            break;
        case 39:
            replaceWindow("next");
            break;
        }
    }
};
