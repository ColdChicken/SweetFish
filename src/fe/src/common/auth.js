require("jquery")

import Ajax from "./ajax"
import Config from "../config"

const Auth = (($) => {

  class Auth {
    constructor() {

    }

    // 获取当前用户信息
    get_current_user_info(callback, current_path="") {
      var that = this
      Ajax.get(
        Config.USER_INFO_URL,
        {},
        true,
        // success
        (data, textStatus, jqXHR) => {
          callback(jqXHR.responseJSON)
        },
        // error
        (jqXHR, textStatus, errorThrown) => {
          console.log('登陆失败')
          console.log(jqXHR)
          // 登陆失败则获取跳转页面地址
          window.location.href = `${Config.LOGIN_PAGE}`
        }
      )
    }
  }

  return Auth

})(jQuery)

export default Auth