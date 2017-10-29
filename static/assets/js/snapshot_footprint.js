window.onload = function() {

    var anchors = document.getElementsByTagName("a");

    for (var i = 0; i < anchors.length; i++) {
        if((anchors[i].text))
        anchors[i].href = "/redirect/?url=" + encodeURIComponent(anchors[i].href)
    }
}


