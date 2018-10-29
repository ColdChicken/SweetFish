<template>
  <div id="app">
    <Header/>
    <Content/>
  </div>
</template>

<script>

import Header from '@/components/framework/Header'
import Content from '@/components/framework/Content'

export default {
  name: 'App',
  components: {
    Header,
    Content
  }
}
</script>

<style>

</style>
