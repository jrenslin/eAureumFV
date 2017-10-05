/* ---------------------------------
// Search functions:
// - Loads index of all files
// - Binds opening the search menu to the corresponding navigation option
// - Searches through the file index using regular expressions
// + Inserts the results into a list
// + Binds clicking 
// --------------------------------- */

// Load index of all files
// If possible, from session storage (--> https://developer.mozilla.org/en-US/docs/Web/API/Web_Storage_API)
// Else, load data and store it in session storage

var stored = sessionStorage['fileIndex'];
var elements;
if (stored) elements = JSON.parse(stored);
else {
    var requestURL  = '/index';
    var request     = new XMLHttpRequest();
    request.open('GET', requestURL);
    request.setRequestHeader("Cache-Control", "no-cache");
    request.responseType = 'json';
    request.send();
    var localelements;
    request.onload = function() {
        elements = request.response;
        sessionStorage['fileIndex'] = JSON.stringify(elements);
    };
}

// Open search on click on search
if (document.getElementById('navigation_Search') != null) {
    document.getElementById('navigation_Search').addEventListener("click", function(){
        document.getElementById('search').style = "";
    });
}

// Bind functions to search form
if (document.getElementById('search') != null) {
    document.getElementById('search').addEventListener('keydown', complete);
}

// Removes all contents from a given element
function emptyElement (id) {
    var element = document.getElementById(id);
    while (element.firstChild) {
        element.removeChild(element.firstChild);
    }
}

// Add an new option
function addSelectOption (path) {
    var li = document.createElement("li");                           // Create a <p> node
    var t = document.createTextNode(path);                           // Create a text node
    li.addEventListener("click", function(){
        window.location.href = "/file?p="+ path;
    });
    li.appendChild(t);                                               // Append the text to <p>
    document.getElementById("searchSelectors").appendChild(li);      // Append <p> to <div> with id="myDIV"
}

// Search through list of files and add appropriate options to list
function complete (e) {

    emptyElement("searchSelectors");
    var regex = new RegExp(document.getElementById("searchInput").value);

    var hits = 0;
    for (var i = 0; i < elements.length; i++) {

        var match = regex.exec(elements[i]);

        if (match != null) {
            addSelectOption(elements[i]);
            hits++;
            if (hits >= 5) { // Stop after 5 elements
                break;
            }
        }

    }

}
