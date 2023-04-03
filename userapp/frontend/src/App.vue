<template>
  <div id="app">
    <img alt="Shopee logo" style="width: 400px; border: none" src="./assets/img.png">
    <!-- 对应的组件内容渲染到router-view中 -->
    <router-view></router-view>

<!--    <Main />-->
<!--    <HelloWorld msg="Welcome to Your Vue.js App"/>-->
  </div>
</template>

<script>
// import HelloWorld from './components/HelloWorld.vue'
// import Main from './components/Main'

export default {
  name: 'App',
  components: {
    // HelloWorld
    // Main
  }
}
</script>

<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
}
</style>
