import axios from 'axios'
import router from '@/router'

const state = {
  user: null,
  isAdmin: false,
  devAccId: null
}

const getters = {
  IsAuthenticated: state => !!state.user,
  IsAdmin: state => state.isAdmin,
  CurrentUser: state => state.user,
  CurrentDevAccId: state => state.devAccId
}

const actions = {
  ForgetUser({ commit }) {
    commit("unauthorize")
  },

  CheckAuth({ commit }) {
    axios.post('/auth/check')
      .then( (res) => commit('authorize', res.data) )
      .catch( err => { 
        if (err) { /* do nothing */ }
        if (router.currentRoute.meta.requiresAuth)
        // if user is not auth and page requires auth, then go to login page
          router.push("/login")
      })
  }
}

const mutations = {
  authorize(state, userdata) {
    state.user = userdata.username
    state.isAdmin = userdata.isAdmin
    //state.devAccId = userdata.devAccId add into JWt
  },
  
  unauthorize(state) {
    state.user = null
    state.isAdmin = false
    state.devAccId = null
  }
}

export default {
  state,
  getters,
  actions,
  mutations
}