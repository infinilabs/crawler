$(function () {

    var conn;
    var msg = $("#msg");
    var log = $("#log");
    var logging = $("#logging");
    var maxLog=50;

    function appendRequestLog(msg) {
        req = $("<div class='request'><pre>" + msg + "</pre></div>");
        appendLog(req)
    }

    function parsePushedData(msg) {
        op=msg.split(" ",1);
        if(op!=undefined){
            msg=msg.substr(msg.split(" ",1)[0].length+1,msg.length);
            if(op=="PRIVATE"){
                req = $("<div class='response'><pre>" + msg + "</pre></div>");
                return appendLog(req)
            }else if(op=="PUBLIC"){
                req = $("<div class='response'><pre>" + msg + "</pre></div>");
                appendLoggingLog(req)
            }
        }
        console.log("invalid format: ",msg)
    }
    logSize=0;
    function appendErrorLog(msg) {
        req = $("<div class='error'><b>" + msg + "</b></div>");
        appendLog(req)
    }

    function appendLog(msg) {
        if(log.children().length>=maxLog){
            $(log).find('div:first').remove();
            logSize--;
        }
       msg.appendTo(log);
        logSize++;
    }

    loggingSize=0;
    //TODO performance improve
    function appendLoggingLog(msg) {
        if(loggingSize>=maxLog){
            $(logging).find('div:first').remove();
            loggingSize--;
        }
        msg.appendTo(logging);
        loggingSize++;
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

    var tm;
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
            };
            conn.onclose = function (evt) {
                var msg="Connection closed.";
                $("#connect_status").text(msg);
                $("#connect_status").removeClass("uk-alert-success").addClass("uk-alert-danger");
                tm=setInterval(refreshFunc, 5000);
                console.log("set refresh trigger")
            };
            conn.onmessage = function (evt) {
                parsePushedData(evt.data);
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



