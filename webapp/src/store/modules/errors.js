import router from '@/router'

const state = {
    defaultError: {
        message: "",
        details: "",
        actionDescr: "",
        actionPath: ""
    },
    errors: []
}

const getters = {
    PopErrors: state => { 
        let e = state.errors
        state.errors = []
        return e  
    },
    GetErrors: state => state.errors,
    ErrorsIsEmpty: state => state.errors.length == 0
}

const actions = {
  ShowError ({ commit }, error) {
    commit("addError", error)
    router.push("/error").catch( ()=>{} )
  }
}

const mutations = {
    addError: (state, error) => { 
        let newErr = Object.assign({}, state.defaultError)
        newErr = Object.assign(newErr, error)
        state.errors.push(newErr)
    },
    dropErrors: (state) => {
        state.errors = []
    }
}

export default {
    state,
    getters,
    actions,
    mutations
}