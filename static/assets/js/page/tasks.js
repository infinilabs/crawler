/**
 * Created by medcl on 2016/12/20.
 */

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

$(function(){
    $(".timeago").timeago();
});
