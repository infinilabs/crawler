/**
 * Created by medcl on 2016/12/20.
 */

var totalRow=20
function drawTable(data) {
    //.eq(-1).remove();
    $("#records").children().remove();
    for (var i = 0; i < data.length; i++) {
        drawRow(data[i]);
    }
}

function drawRow(rowData) {
    var row = $("<tr />")
    if(!rowData.success){
        row = $("<tr class='failed_task' />")
    }
    $("#records").append(row);

    if(rowData.page != undefined ){
        col1= getdata(rowData.page.domain);
        col2= getdata(rowData.page.path);
        if(rowData.page.metadata!=undefined){
            col3=getdata(rowData.page.metadata.title);
        }else{
            col3="N/A";
        }
        col4= getdata(rowData.page.size);
        row.append($("<td>" + col1 + "</td>"));
        row.append($("<td>" + col2 + "</td>"));
        row.append($("<td>" + col3 + "</td>"));
        row.append($("<td>" + col4 + "</td>"));
    }else{
        row.append($("<td colspan='2'>" + rowData.url + "</td>"));
        row.append($("<td colspan='2'>" + rowData.message + "</td>"));
    }

    row.append($("<td class='timeago'>" + timeago(rowData.updated) + "</td>"));
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
