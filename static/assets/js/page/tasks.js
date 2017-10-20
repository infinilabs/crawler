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



$('[data-uk-pagination]').on('select.uk.pagination', function(e, pageIndex){
    alert('You have selected page: ' + (pageIndex+1));
});

$(function(){
    $(".timeago").timeago();
});
