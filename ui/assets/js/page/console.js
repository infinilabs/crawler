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
               return appendLoggingLog(req)
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


function updateLogging(enabled){

    $("#update_config_response").text("");

    txt=$("input").val();

    var data={};
    data.realtime=enabled;
    data.push_log_level=$("#log_level").val();
    data.func_pattern=$("#func_pattern").val();
    data.file_pattern=$("#file_pattern").val();

        //load domain data
        $.ajax({
            url: '/setting/logger',
            type: "post",
            data: JSON.stringify(data),
            cache : false,
            dataType: "json",
            success: function(data, textStatus, jqXHR) {
                UIkit.notify("Success");

                $(".uk-icon-play").attr("disabled",enabled)
                $(".uk-icon-stop").attr("disabled",!enabled)

            },
            error: function(XMLHttpRequest, textStatus, errorThrown) {
                UIkit.notify("Error: " + errorThrown);
            }

        });


}



window.onload = function(){
    var s = setInterval("scrollDiv()", 10000);
}
function scrollDiv(){
    var logging = document.getElementById("logging");
    logging.scrollTop = logging.scrollHeight;
    var log = document.getElementById("log");
    log.scrollTop = log.scrollHeight;
}
