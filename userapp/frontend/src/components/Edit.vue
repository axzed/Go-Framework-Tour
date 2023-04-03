<template>
  <div class="container">
    <div class="row">
      <div class="col-3"></div>
      <div class="col-6">
        <b-alert :show="errMsg.length > 0" variant="danger">{{errMsg}}</b-alert>
        <b-form>
<!--          <b-form-group id="input-group-1" label-for="input-1">-->
<!--            <b-img :src="form.avatar"></b-img>-->
<!--            <b-form-file-->
<!--                accept="image/*"-->
<!--                v-model="avatarImg"-->
<!--                :state="Boolean(avatarImg)"-->
<!--                placeholder="Choose a file or drop it here..."-->
<!--                drop-placeholder="Drop file here..."-->
<!--            ></b-form-file>-->
<!--          </b-form-group>-->
          <br/>
          <b-form-group id="input-group-2" label-for="input-2">
            <b-form-input
                id="input-2"
                v-model="form.name"
                placeholder="Enter name"
                required
            ></b-form-input>
          </b-form-group>
          <br />
          <b-form-group id="input-group-3" label-for="input-3">
            <b-form-input
                id="input-3"
                v-model="form.email"
                placeholder="Enter Email"
                required
            ></b-form-input>
          </b-form-group>
          <br>
          <b-button variant="primary" v-on:click="save">保存</b-button>
        </b-form>
      </div>

    </div>

  </div>
</template>

<script>
import {BAlert, BForm, BFormInput, BFormFile, BImg} from 'bootstrap-vue'
export default {
  name: "Edit",
  comments: {
    BAlert,
    BForm,
    BFormInput,
    BFormFile,
    BImg
  },
  created() {
    this.$http.get(process.env.VUE_APP_ENDPOINT + "/profile")
        .then(resp => {
          this.form = resp.body.data
        }, resp => {
          if (resp.body !== undefined && this.body.msg !== undefined) {
            this.errMsg = resp.body.msg
          } else {
            this.errMsg ='系统异常'
          }
        })
  },
  data() {
    return {
      errMsg: '',
      avatarImg: null,
      form: {
        name: '',
        email: '',
        avatar: ''
      }
    }
  },
  watch: {
    // avatarImg(val) {
    //   let data = new FormData()
    //   data.append('file', val)
    //   this.$http.post(process.env.VUE_APP_ENDPOINT + "/upload", data, {
    //     headers: { 'Content-Type': 'multipart/form-data' }
    //   }).then(resp => {
    //     if (resp.body !== undefined && resp.body.Data !== undefined) {
    //       this.form.avatar = resp.body.Data
    //       console.log("aaa")
    //     } else {
    //       this.errMsg = 'unknown error'
    //     }
    //   }, resp => {
    //     if (resp.body !== undefined && this.body.Msg !== undefined) {
    //       this.errMsg = resp.body.Msg
    //     } else {
    //       this.errMsg ='system error'
    //     }
    //   })
    // }
  },
  methods: {
    // goToEdit() {
    //   this.$router.push("/edit")
    // },
    save() {
      this.$http.post(process.env.VUE_APP_ENDPOINT + "/update", this.form)
      .then(() => {
        this.$router.push("/profile")
      }, resp => {
        if (resp.body !== undefined && this.body.msg !== undefined) {
          this.errMsg = resp.body.msg
        } else {
          this.errMsg ='系统异常'
        }
      })
    }
  }
}
</script>

<style scoped>

</style>