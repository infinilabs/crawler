$(function () {

    var conn;
    var msg = $("#msg");
    var log = $("#log");

    function appendRequestLog(msg) {
        req = $("<div class='request'><pre>" + msg + "</pre></div>");
        appendLog(req)
    }

    function appendResponseLog(msg) {
        req = $("<div class='response'><pre>" + msg + "</pre></div>");
        appendLog(req)
    }

    function appendErrorLog(msg) {
        req = $("<div class='error'><b>" + msg + "</b></div>");
        appendLog(req)
    }

    function appendLog(msg) {
        var d = log[0]
        var doScroll = d.scrollTop == d.scrollHeight - d.clientHeight;
        msg.appendTo(log)
        if (doScroll) {
            d.scrollTop = d.scrollHeight - d.clientHeight;
        }
    }

    $("#form").submit(function () {
        if (!conn) {
            return false;
        }
        if (!msg.val()) {
            return false;
        }

        if(msg.val()=="CLS"||msg.val()=="cls"||msg.val()=="clear"||msg.val()=="CLEAR"){
            log.children().remove();
            return
        }

        conn.send(msg.val());
        appendRequestLog(msg.val())
        msg.val("");
        return false
    });

    var refreshFunc = function(){
        console.log("reconnecting");
        connect();
    };

    var tm
    function connect(){

        clearInterval(tm);

        if (window["WebSocket"]) {
            host = location.hostname + (location.port ? ':' + location.port : '');
            conn = new WebSocket("ws://" + host + "/ws");
            conn.onopen = function (evt) {
                var msg="Connection established.";
                log.children().remove();
                $("#connect_status").text(msg);
                $("#connect_status").removeClass("uk-alert-danger").addClass("uk-alert-success");
                clearInterval(tm);
                console.log("remove refresh trigger")
            }
            conn.onclose = function (evt) {
                var msg="Connection closed.";
                $("#connect_status").text(msg);
                $("#connect_status").removeClass("uk-alert-success").addClass("uk-alert-danger");
                tm=setInterval(refreshFunc, 5000);
                console.log("set refresh trigger")
            }
            conn.onmessage = function (evt) {
                appendResponseLog(evt.data);
            }
        } else {
            appendErrorLog("Your browser does not support WebSockets.");
        }
    }

    connect();

});

jQuery(document).ready(function(){
    $(document).bind('keyup.s', function (){
        $("#msg").focus()
    });
});



