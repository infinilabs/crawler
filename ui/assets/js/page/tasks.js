/**
 * Created by medcl on 2016/12/20.
 */

var totalRow=20
function drawDomainTable(data) {
    $("#domain-records").children().remove();
    for (var i = 0; i < data.length; i++) {
        drawDomainRow(data[i]);
    }
}


function drawDomainRow(rowData) {
    var row = $("<li />")
    $("#domain-records").append(row);

    row.append($("<span class='uk-column-1-1' ><a class=uk-button href=\"javascript:loadData('" + rowData.host + "')\" >" + rowData.host + " ("+rowData.links_count+")</a></span>"));

}


function drawTable(data) {
    //.eq(-1).remove();
    $("#records").children().remove();
    for (var i = 0; i < data.length; i++) {
        drawRow(data[i]);
    }
}

function drawRow(rowData) {
    var row = $("<tr />")
    if(rowData.status==2){
        row = $("<tr class='failed_task' />")
    }
    $("#records").append(row);



    if(rowData.page != undefined ){
        col1= getdata(rowData.domain);
        col2= getdata(rowData.path);

        if(rowData.page.metadata!=undefined){
            col3=getdata(rowData.page.metadata.title);
        }else{
            col3="N/A";
        }
        col4= formatBytes(getdata(rowData.page.size),1);
        row.append($("<td>" + col1 + "</td>"));
        row.append($("<td>" + col2 + "</td>"));
        row.append($("<td>" + col3 + "</td>"));
        row.append($("<td>" + col4 + "</td>"));
    }
    else if(rowData.seed!=undefined){
        row.append($("<td>" + rowData.id + "</td>"));
        row.append($("<td colspan='2'>" + rowData.seed.url + "</td>"));
        row.append($("<td>" + rowData.status + "</td>"));
    }else{
        row.append($("<td>" + rowData.id + "</td>"));
        row.append($("<td>" + rowData.url + "</td>"));
        row.append($("<td colspan='2'>" + rowData.message + "</td>"));
    }

    row.append($("<td class='timeago'>" + timeago(rowData.updated!=undefined?rowData.updated:rowData.created) + "</td>"));
}

function getdata(v){
    try{
        return a=v
    }catch(e){
        return v
    }
}

function timeago(v){
    try{
        return  jQuery.timeago(v)
    }catch(e){
        return v
    }
}


function formatBytes(bytes,decimals) {
    if(bytes == 0) return '0 Byte';
    var k = 1000;
    var dm = decimals + 1 || 3;
    var sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB'];
    var i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + ' ' + sizes[i];
}

function loadDomain(){
    $("#domain-alert").text("Loading");
    //load domain data
    $.ajax({
        url: '/domains?from=0&size=40',
        type: "get",
        dataType: "json",
        success: function(data, textStatus, jqXHR) {
            $("#domain-alert").text("Total: "+data.total);
            drawDomainTable(data.result);
        }
    });
}



function loadData(domain){

    $("#alert").text("Loading");

    para='/tasks?from=0&size=20';
    if(domain!=undefined){
        para=para+"&domain="+domain;
    }

    //load task data
    $.ajax({
        url: para,
        type: "get",
        dataType: "json",
        success: function(data, textStatus, jqXHR) {
            $("#alert").text("Total: "+data.total);
            drawTable(data.result);
        }
    });
}

loadDomain();

loadData();

$('[data-uk-pagination]').on('select.uk.pagination', function(e, pageIndex){
    alert('You have selected page: ' + (pageIndex+1));
});
