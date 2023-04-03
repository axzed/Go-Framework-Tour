import Vue from "vue";
import VueRouter from "vue-router";

// 引入组件
import Main from "./components/Main.vue";
import Edit from "./components/Edit.vue";
import Profile from "./components/Profile";

// 要告诉 vue 使用 vueRouter
Vue.use(VueRouter);

const routes = [
    {
        path:"*",
        component: Main
    },
    {
        path: "/edit",
        component: Edit
    },
    {
        path: "/profile",
        component: Profile
    },
]

let router =  new VueRouter({
    routes
})
export default router;