<template>
    <nav class="navbar navbar-inverse navbar-fixed-top">
      <div class="container-fluid">
        <div class="navbar-header">
          <a class="navbar-brand" href="#">SweetFish 管理控制台</a>
        </div>
        <div id="navbar" class="navbar-collapse collapse">
          <ul class="nav navbar-nav navbar-right">
            <li><a href="/v1/ajax/auth/logout">Logout</a></li>
          </ul>
        </div>
      </div>
    </nav>
</template>

<script>
require('jquery')

require('@/assets/bootstrap/css/bootstrap.min.css')
require('@/assets/bootstrap/css/dashboard.css')

require('@/assets/bootstrap/js/bootstrap.min.js')


export default {
  name: 'Header',
  data () {
    return {
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>

</style>
