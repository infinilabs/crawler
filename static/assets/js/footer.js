/**
 * Created by JetBrains PhpStorm.
 * User: Medcl
 * Date: 12-12-4
 * Time: 下午8:10
 */


//<![CDATA[

if ("undefined" != typeof maxpage) {
    var userAgent = navigator.userAgent.toLowerCase();
    var is_opera = userAgent.indexOf('opera') != -1 && opera.version();
    var is_moz = (navigator.product == 'Gecko') && userAgent.substr(userAgent.indexOf('firefox') + 8, 3);
    var is_ie = (userAgent.indexOf('msie') != -1 && !is_opera) && userAgent.substr(userAgent.indexOf('msie') + 5, 3);
} else {
    var maxpage = 0;
}
document.onkeyup = function (e) {
    e = e ? e : window.event;
    var tagname = is_ie ? e.srcElement.tagName : e.target.tagName;
    if (tagname == 'INPUT' || tagname == 'TEXTAREA') return;
    actualCode = e.keyCode ? e.keyCode : e.charCode;
    if (actualCode == 83) {
        $('#query').focus();
        $('#query').select();
    }

    if (maxpage > 1) {
        if (actualCode == 39) {
            if ("undefined" != typeof next_page) {
                window.location = next_page;
            }
        }
        if (actualCode == 37) {
            if ("undefined" != typeof prev_page) {
                window.location = prev_page;
            }
        }
    }

};
//]]>


//<![CDATA[

jQuery(function () {
    var a1;
    var onAutocompleteSelect = function (value, data) {
        $("#f").submit();
    };

    var autocompleteSelectWidth = document.getElementsByClassName("sbx-custom__wrapper").item(0).offsetWidth - 2;
    var suggest_url = $("#search-context a[class='active']").attr("suggest_url");
    var context = $("#search-context a[class='active']").attr("context");
    if (suggest_url == undefined || suggest_url == "") {
        suggest_url = "/suggest/"
    }
    var options = {
        serviceUrl: suggest_url,
        width: autocompleteSelectWidth,
        delimiter: /(,|;)\s*/,
        onSelect: onAutocompleteSelect,
        deferRequestBy: 50, //miliseconds
        params: {version: '1.0', type: "input", context: context},
        noCache: true //set to true, to disable caching
    };
    a1 = $('#query').autocomplete(options);

});
//]]>


//<![CDATA[
document.querySelector('.searchbox [type="reset"]').addEventListener('click', function () {
    this.parentNode.querySelector('input').setAttribute("value", "");
    this.parentNode.querySelector('input').focus();
});
//]]>


//<![CDATA[
function submitSearch() {
    var kw = document.getElementById('query').value;
    console.log(kw)
//mpq.track(kw);
}
//]]>
