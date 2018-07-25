/**
 * Created by medcl on 2017/1/7.
 */
function safeGetValue(v){
    if(v==undefined){
        return 0;
    }
    return v;
}

function loadData(){

    // 加载统计数据
    $.get('/stats').done(function (data) {
        if(data["queue.check"]){
            $("#checker_task_num").text(data["queue.check"].pop+" / "+data["queue.check"].push+", valid: "+data["checker.url"].valid_seed);
        }

        if(data["crawler.pipeline"]){
            $("#crawler_task_num").text(safeGetValue(data["crawler.pipeline"].finished)+" / "+data["crawler.pipeline"].total+", error: "+safeGetValue(data["crawler.pipeline"].error)+", break: "+safeGetValue(data["crawler.pipeline"].break)+", queue: "+safeGetValue(data["queue.fetch"].pop)+" / "+data["queue.fetch"].push);
        }
    });


    //加载最新任务数据
    $.ajax({
        url: '/task?size=20',
        type: "get",
        dataType: "json",
        success: function(data, textStatus, jqXHR) {
            drawTable(data.result);
        }
    });

}

timeTicket = setInterval(loadData ,2000);

loadData();
