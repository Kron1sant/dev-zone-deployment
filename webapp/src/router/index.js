import Vue from 'vue'
import VueRouter from 'vue-router'

import Home from '@/views/Home.vue'
import Login from '@/views/Login.vue'
import ErrorsPage from '@/views/Errors.vue'
import DevAccounts from '@/views/DevAccounts/AccTable.vue'
import AppUsers from '@/views/Users/UsersTable.vue'
import VM from '@/views/VM/VMTable.vue'
import Test from '@/views/Test.vue'

Vue.use(VueRouter)
export const routes = [
    {
      path: "/",
      name: "Home",
      component: Home,
      meta: { 
        label: "Домой",
        navigation: false, // Don't display in common navigation panel
        requiresAuth: false 
      }
      
    },
    {
      path: "/test",
      name: "Test",
      component: Test,
      meta: { 
        label: "Тест",
        navigation: false,
        requiresAuth: true 
      }
    },
    {
      path: "/login",
      name: "Login",
      component: Login,
      meta: { 
        label: "Вход",
        navigation: false, // Don't display in common navigation panel
        requiresAuth: false 
      }
    },
    {
      path: "/error",
      name: "Error",
      component: ErrorsPage,
      meta: { 
        label: "",
        navigation: false, // Don't display in common navigation panel
        requiresAuth: false 
      }
    },
    {
      path: "/devaccounts",
      name: "Developer accounts",
      component: DevAccounts,
      meta: { 
        label: "Аккаунты разработчиков",
        navigation: true,
        requiresAuth: true 
      }
    },
    {
      path: "/users",
      name: "Application users",
      component: AppUsers,
      meta: { 
        label: "Список пользователей",
        navigation: false,
        requiresAuth: true 
      }
    },
    {
      path: "/virtualmachines",
      name: "VM instances",
      component: VM,
      meta: { 
        label: "Виртуальные машины",
        navigation: true,
        requiresAuth: true 
      }
    },
]

const router = new VueRouter({
    mode: 'history',
    base: process.env.BASE_URL,
    routes
})


export default router