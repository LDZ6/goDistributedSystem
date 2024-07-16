$(function () {
    baseApp.init();
    //当窗口重置时,重新计算窗口高度
    $(window).resize(function () {
        baseApp.resizeIframe();
    })
})
 
var baseApp = {
    init: function () {
        this.initAside()
        this.resizeIframe()
    },
    initAside: function () { //左侧菜单栏隐藏显示子栏
        $(".aside h4").click(function () {
            $(this).sibling("ul").slideToggle();
        })
    },
    resizeIframe: function() {  // 设置iframe高度
        $("rightMain").height($(window).height()-80)
    },
}