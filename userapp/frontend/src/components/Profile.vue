<template>
  <div class="container">
    <div class="col-3"></div>
    <div class="col-6">

    </div>
    <b-alert :show="errMsg.length > 0" variant="danger">{{errMsg}}</b-alert>
    <b-card :title="user.name || 'Your Name'" :img-src="user.avatar" img-alt="Card image" img-left class="mb-3">
      <b-list-group flush>
        <b-list-group-item>Email: {{user.email || 'your email'}}</b-list-group-item>
        <b-list-group-item><b-button v-on:click="goToEdit" variant="primary">编辑</b-button></b-list-group-item>
      </b-list-group>
    </b-card>

  </div>
</template>

<script>
import {BCard, BListGroup, BListGroupItem, BButton} from 'bootstrap-vue'
export default {
  name: "Profile",
  comments: {
    BCard,
    BListGroup,
    BListGroupItem,
    BButton
  },
  created() {
    this.$http.get(process.env.VUE_APP_ENDPOINT + "/profile", {}, {
      withCredential:true
    })
    .then(resp => {
      this.user = resp.body.data
    }, resp => {
      if (resp.body !== undefined && resp.body.msg !== undefined) {
        this.errMsg = resp.body.msg
      } else {
        this.errMsg ='系统异常'
      }
    })
  },
  data() {
    return {
      errMsg: '',
      user: {
        email: "",
        avatar: "",
      }
    }
  },
  methods: {
    reset() {

    },
    goToEdit() {
      this.$router.push("/edit")
    }
  }
}
</script>

<style scoped>

</style>