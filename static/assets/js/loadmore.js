/**
 * Created by JetBrains PhpStorm.
 * User: liaosy
 * Date: 17-06-18
 * Time: 下午9:10
 */

//把URL参数解析为一个对象
function parseQueryString(url) {
    var str = url.split("?")[1];
    var items = str.split("&");
    var result = {};
    var arr = [];

    for (var i = 0; i < items.length; i++) {
        if (items[i]) {
            arr = items[i].split('=');
            if (arr.length <= 1) {
                continue;
            }
            //decode 参数值,避免ajax请求参数被自动二次编码
            result[arr[0]] = decodeURIComponent(arr[1]);
        }
    }

    return result;
}

/**
 ----------加载更多模块-----------
 */
$(function () {
    var $e = $(".pnnext");
    //手动点击加载更多方式获取数据
    $e.on("click", function () {
        var $this = $(this);
        var $tips = $(".load-tips");
        var originText = $this.text();
        var total = parseInt($this.attr("data-total"));
        var size = parseInt($this.attr("data-size"));
        var from = parseInt($this.attr("data-from"));
        var queryObj = parseQueryString(location.href);
        queryObj.from = from + size;
        var url = "ajax_more_item/";

        if (from < total) {
            if ($this.hasClass("disabled")) {
                return false;
            }
            $this.addClass("disabled");
            $this.children(".load-icon").addClass("loading");
            $this.children(".load-text").text($this.attr("data-load-text"));

            $.ajax({
                url: url,
                type: 'GET',
                data: queryObj,
                dataType: 'html',
                error: function () {
                    $tips.text("Network error, please try again!").show();
                },
                success: function (res) {
                    if (res) {
                        $this.attr("data-from", queryObj.from);
                        $(".item-view").append(res);
                    } else {
                        $tips.text("Error with data loading, please try again!").show();
                    }
                },
                complete: function () {
                    $this.children(".load-icon").removeClass("loading");
                    $this.children(".load-text").text(originText);
                    $this.removeClass("disabled");
                }
            });
        } else {
            $tips.text("No more data!").show();
            $this.addClass("disabled");
            setTimeout(function () {
                $this.fadeOut();
            }, 3000);
        }
    });

    //目标元素在首屏时,直接隐藏(内容少,不需要加载更多)
    //if($e.offset().top <= $(window).height()) {
    //    $e.hide();
    //}

    //监听滚动条,模拟触发点击加载更多
    $(window).scroll(function (event) {
        var winHeight = $(window).height();//可视区域高度
        var winPos = $(window).scrollTop();//屏幕距离页面顶部偏移量
        var pnnextPos = $e.offset().top;//目标监控元素
        var heightDiff = winHeight + winPos - pnnextPos;//差值大于0 表示目标元素已经出现在可视区域
        //偏移量超过
        if (heightDiff > 50) {
            if (typeof($e.attr("disabled")) == "undefined") {
                $e.click();
            }
        }
    });
});


/**
 ----------回到顶部模块-----------
 */
$(function () {
    $(".c-back").on("click", function () {
        var _this = $(this);
        $('html,body').animate({scrollTop: 0}, 500, function () {
            _this.fadeOut();
        });
    });

    $(window).scroll(function () {
        var winHeight = $(window).height();//可视区域高度
        var htmlTop = $(window).scrollTop();

        if (htmlTop > (winHeight * 1.5)) {
            $(".c-back").fadeIn();
        } else {
            $(".c-back").fadeOut();
        }
    });
});
