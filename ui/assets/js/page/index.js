/**
 * Created by medcl on 2017/1/7.
 */
function safeGetValue(v){
    if(v==undefined){
        return 0;
    }
    return v;
}

var myChart = echarts.init(document.getElementById('gopa_stats'));

option = {
    tooltip : {
        formatter: "{a} <br/>{b} : {c}%"
    },
    toolbox: {
        feature: {
            restore: {show: true},
            saveAsImage: {show: true}
        }
    },
    series: [
        {
            name: 'Checker',
            center: ['20%', '55%'],
            type: 'gauge',
            detail: {formatter:'{value}%'},
            data: [{value: 0, name: 'Checker'}]
        },{
            name: 'Crawler',
            center: ['77%', '50%'],
            type: 'gauge',
            detail: {formatter:'{value}%'},
            data: [{value: 0, name: 'Crawler'}]
        }
    ]
};

myChart.setOption(option);

function loadData(){

    // 加载统计数据
    $.get('/stats').done(function (data) {
        if(data["queue.check"]){
            $("#checker_task_num").text(data["queue.check"].pop+" / "+data["queue.check"].push+", valid: "+data["checker.url"].valid_seed);
            option.series[0].data[0].value = ((data["queue.check"].pop/data["queue.check"].push)*100).toFixed(2) - 0;
        }

        if(data["crawler.pipeline"]){
            $("#crawler_task_num").text(safeGetValue(data["crawler.pipeline"].finished)+" / "+data["crawler.pipeline"].total+", error: "+safeGetValue(data["crawler.pipeline"].error)+", break: "+safeGetValue(data["crawler.pipeline"].break)+", queue: "+safeGetValue(data["queue.fetch"].pop)+" / "+data["queue.fetch"].push);
            option.series[1].data[0].value = (((parseInt(safeGetValue(data["queue.fetch"].pop)))/parseInt(safeGetValue(data["queue.fetch"].push)))*100).toFixed(2) - 0;
        }

        myChart.setOption(option, true);
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
