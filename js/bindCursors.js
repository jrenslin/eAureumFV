/* Bind keys to body. */
if (document.getElementsByTagName('body')[0] != null) {
    document.getElementsByTagName('body')[0].addEventListener('keydown', hotkeys);
}

/* Function for referring based on a link in the page */
function replaceWindow (id) {
    var link = document.getElementById(id);
    if (link != null) {
        window.location.href = link.href;
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

    if (document.getElementById("file") != null) { // Only for use on the file / preview page
        if (e.ctrlKey) {
            switch (e.keyCode) {
            case 40:
                if (String(window.location.search).search("fullPreview") < 1) {
                    target = String(window.location).replace(window.location.search, window.location.search + "&fullPreview=yes");
                } else {
                    target = String(window.location).replace("&fullPreview=yes", "");
                }
                window.location.href = target;
            }
        }
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
    // Bind cursor keys (without CTRL) for CBZ viewer
    else if (document.getElementById("page0") != null) { // this only applies if we have a paged site
        // Get number of the current page from anchor
		    var no = parseInt(window.location.hash.slice(-1));
        var offset = window.location.search.split("&")[1].split("=")[1];

        switch (e.keyCode) {
        case 37:
        case 38:
            if (no > 0) {
                window.location.hash = "#page" + (no - 1);
            } else if (document.getElementById("prevBatch") != null) {
                // Currently offset is always the second parameter.
                var target = "";
                target = String(window.location).replace("offset=" + offset, "offset=" + (parseInt(offset) - 10)).split("#")[0];
                target = target + "#page9";
                window.location.href = target;
            }
            break;

        case 39:
        case 40:
            if (parseInt(offset) + parseInt(no) == parseInt(document.getElementById("max").innerHTML) - 1) break;
            if (no < 9) {
                window.location.hash = "#page" + (no + 1);
            } else if (document.getElementById("nextBatch") != null) {
                // Currently offset is always the second parameter.
                var target = String(window.location).replace("offset=" + offset, "offset=" + (parseInt(offset) + 10)).split("#")[0];
                target = target + "#page0";
                window.location.href = target;
            }
            break;
        }

    }
};

