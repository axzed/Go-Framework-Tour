<template>
  <div>
    <h1>Welcome to simple-user demo</h1>
    <div class="col-6 container">
      <b-form>
        <b-form-group id="input-group-1">
          <b-form-input
              id="input-1"
              v-model="form.email"
              type="email"
              placeholder="Enter email"
              required
              :invalid-feedback="'invalid email'"
              :state="isValidEmail"
          ></b-form-input>
        </b-form-group>
        <br/>

        <b-form-group id="input-group-2">
          <b-form-input
              id="input-2"
              v-model="form.password"
              type="password"
              placeholder="8~16 characters"
              required
              :invalid-feedback="invalidPwd"
              :state="isValidPwd"
          ></b-form-input>
        </b-form-group>

        <br/>
        <b-form-group  v-if="signup" id="input-group-2">
          <b-form-input
              id="input-2"
              v-model="form.confirmPwd"
              placeholder="confirm password"
              type="password"
              required
              :invalid-feedback="notMatchPwd"
              :state="form.confirmPwd.length >= 8 && form.password === form.confirmPwd"
          ></b-form-input>
          <br />
        </b-form-group>

        <b-alert :show="msg.length > 0" variant="success">{{msg}}</b-alert>
        <b-alert :show="errMsg.length > 0" variant="danger">{{errMsg}}</b-alert>

        <b-button v-if="!signup" v-on:click="login" variant="primary">Sign In</b-button>
        &nbsp;&nbsp;
        <a v-if="!signup" v-on:click="toSignUp" href="#">Create Account</a>
        <b-button v-if="signup" v-on:click="signUp" variant="danger">Sign Up</b-button>

<!--        <b-button v-on:click="goToProfile(123)" variant="danger">GO GO GO</b-button>-->
        &nbsp;&nbsp;
        <a v-if="signup" v-on:click="toLogin" href="#">Already have account</a>
        <br />
        <br />
      </b-form>
    </div>
  </div>
</template>

<script>

import { BForm, BButton, BFormGroup, BFormInput, BAlert } from 'bootstrap-vue'
import md5 from 'blueimp-md5'

export default {
  name: "Main",
  comments:{
    BForm,
    BButton,
    BFormGroup,
    BFormInput,
    BAlert
  },
  data() {
    return {
      msg: '',
      errMsg:'',
      form: {
        email: '',
        password: '',
        confirmPwd: '',
      },
      signup: false,
    }
  },
  methods: {
    toSignUp() {
      this.signup = true
      this.resetAlert()
    },

    toLogin() {
      this.signup = false
      this.resetAlert()
    },
    login() {
      this.resetAlert()
      // let v = this
      let data = {
        "email": this.form.email,
        "password": md5(this.form.password),
      }
      this.$http.post(process.env.VUE_APP_ENDPOINT + "/login", data)
      .then(() => {
        this.$router.push("/profile")
      }, resp => {
        if (resp.body !== undefined) {
          this.errMsg = resp.body
        } else {
          this.errMsg = "系统异常"
        }
      })

    },
    // signup
    signUp() {
      this.resetAlert()
      if (this.form.email === "") {
        this.errMsg = "邮箱不能为空"
        return;
      }

      if(this.form.password === "") {
        this.errMsg = "密码不能为空"
        return;
      }
      if (this.form.password !== this.form.confirmPwd) {
        this.errMsg = "两次密码不一致，请重新输入"
        return;
      }
      // let v = this
      let data = {
        "email": this.form.email,
        "password": md5(this.form.password),
        "confirm_pwd": md5(this.form.confirmPwd),
      }

      this.$http.post(process.env.VUE_APP_ENDPOINT + "/signup", data)
      .then(() => {
        this.toLogin()
        this.msg = "注册成功，请登录"
      }, resp => {
        let msg = "系统异常"
        if(resp.body !== undefined && resp.body.msg !== undefined) {
          msg = resp.body.msg
        }
        this.errMsg = msg
      })
    },

    resetAlert() {
      this.msg = ''
      this.errMsg = ''
    }
  },
  computed: {
    notMatchPwd() {
      return "两次密码输入不一致"
    },
    invalidPwd() {
      return "密码不少于八个字符"
    },

    isValidPwd() {
      if (!this.signup) {
        return undefined
      }
      return this.form.password.length >= 8
    },

    isValidEmail : function() {
      if (!this.signup) {
        return undefined
      }
      let re = /(.+)@(.+){2,}\.(.+){2,}/;
      return re.test(this.form.email.toLowerCase());
    }
  }
}
</script>

<style scoped>

</style>