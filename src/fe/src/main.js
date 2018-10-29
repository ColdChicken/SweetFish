// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.
import Vue from 'vue'
import App from './App'
import router from './router'

import Auth from './common/auth'

Vue.config.productionTip = false

// 当前用户信息
Vue.prototype._user = ""

// 认证
const auth = new Auth()
auth.get_current_user_info(
  (userinfo) => {
    console.log("当前用户: " + userinfo['username'])
    Vue.prototype._user = userinfo['username']
  }
)

router.beforeEach((to, from, next) => {
  auth.get_current_user_info(
    (userinfo) => {
      console.log("当前用户: " + userinfo['username'])
      Vue.prototype._user = userinfo['username']
      next()
    },
    to.fullpath
  )
})

/* eslint-disable no-new */
new Vue({
  el: '#app',
  router,
  components: { App },
  template: '<App/>'
})


