window.onload = function() {

    var anchors = document.getElementsByTagName("a");

    for (var i = 0; i < anchors.length; i++) {
        if((anchors[i].text))
        anchors[i].href = "/snapshot/?url=" + anchors[i].href
    }
}


