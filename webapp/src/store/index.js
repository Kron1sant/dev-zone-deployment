import Vue from 'vue'
import Vuex from 'vuex'
import auth from './modules/auth'
import errors from './modules/errors'
import devacc from './modules/devacc'
import appusers from './modules/appusers'
import vm from './modules/vm'

Vue.use(Vuex)

export default new Vuex.Store({
    modules: {
        auth,
        errors,
        devacc,
        appusers,
        vm
    }
})