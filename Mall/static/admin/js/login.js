$(function () {
    loginApp.init();
})
 
var loginApp = {
    init: function () {
        this.getCaptcha()  // 调用获取验证码方法
        this.captchaImgChange()  // 调用验证码改变方法
    },
    getCaptcha: function () { // 获取验证码
        $.get("/admin/captcha?t=" + Math.random(), function (response) { // ? t= 随机数,防止浏览器缓存
            //把验证码赋值给input
            $("#captchaId").val(response.captchaId)
            $("#captchaImg").attr("src", response.captchaImage)
        })
    },
    captchaImgChange: function () { // 验证码改变
        var that = this;
        $("#captchaImg").click(function () {
            that.getCaptcha()
        })
    }
}