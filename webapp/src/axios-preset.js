import axios from 'axios'
import router from '@/router'
import store from '@/store'

// Creds over cookies
axios.defaults.withCredentials = true

axios.interceptors.response.use(undefined, error => {
  const itIsAuthCheck = error.config.url == "/auth/check"
  // Suppress AutoCheck errors
  if (itIsAuthCheck) return
  const itIsNotLogin = error.config.url != "/login" 
  // Handle unsuccessful request to backend
  if (error.response === undefined) {
    // We got an error without response - backend is not available
    // show error message
    store.commit('addError', {
      message: "Сервер (бекенд) " + error.config.baseURL + " не доступен",
      details: "",
      actionDescr: "На главную", 
      actionPath: "/"
    })
    router.push("/error").catch( ()=>{} )
  } else {
    // error with response - backend is available
    const statusUnauthorized = error.response.status == 401 
    if (statusUnauthorized && itIsNotLogin) {
      // Request to backend was commited by unauthorized user
      // and it is not logining
      store.dispatch('ForgetUser')
      const err = {message: "Пользователь не авторизован", details: "", actionDescr: "Залогиньтесь", actionPath: "/login"}
      // show error message
      store.commit("addError", err)
      router.push("/error").catch( ()=>{} )
    }
  }
  // Continue rejection to rest promises chaining
  return Promise.reject(error)
})

export default axios

export function isItSuccessResponse(res) { return res !== undefined && Math.floor(res.status / 100 ) == 2 }