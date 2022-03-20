const state = {
  showAccForm: false,
  showAccOpenVPNForm: false,
  newAcc: false,
  readOnlyAcc: false,
  savedAcc: false,
  editedAccountIndex: -1,
  editedAccount: {
    id: 0,
    username: "",
    name: "",
    surname: "",
    patronomic: "",
    email: "",
    hasOVPNCert: false,
    comment: "",
  },
  defaultAccount: {
    id: 0,
    username: "",
    name: "",
    surname: "",
    patronomic: "",
    email: "",
    hasOVPNCert: false,
    comment: "",
  }
}

const getters = {
  DevAccShowForm: state => state.showAccForm,
  DevAccNew: state => state.newAcc,
  DevReadOnlyAcc: state => state.readOnlyAcc,
  DevAccCurrentItem: state => state.editedAccount,
  DevAccCurrentIndex: state => state.editedAccountIndex,
  DevAccDefault: state => state.defaultAccount,
  DevAccHasBeenSaved: state => state.savedAcc,

  DevAccOpenVPNShowForm: state => state.showAccOpenVPNForm,
}

const actions = {
  OpenDevAccDialog ({ commit, DevAccDefault }, devAccParams) {
    if (devAccParams.editedItem === null)
      commit("setDevAccCurrentItem", DevAccDefault)
    else
      commit("setDevAccCurrentItem", devAccParams.editedItem)

    commit("setNewAcc", devAccParams.new)
    commit("setReadOnlyAcc", devAccParams.readOnly)
    commit("setDevAccCurrentIndex", devAccParams.currentIndex) 
    commit("showFormAcc", true)
  },
  ClearDevAccCurrentItem ({ commit }) {
    commit("setDevAccCurrentItem", {})
    commit("setDevAccCurrentIndex", -1)
  }
}

const mutations = {
  setDevAccCurrentItem (state, devAcc) { state.editedAccount = Object.assign({}, devAcc) },
  setDevAccCurrentIndex (state, currentIndex) { state.editedAccountIndex = currentIndex },
  setNewAcc (state, newAcc) { state.newAcc = newAcc },
  setReadOnlyAcc (state, readOnlyAcc) { state.readOnlyAcc = readOnlyAcc },
  showFormAcc (state, show) { state.showAccForm = show },
  setDevAccHasBeenSaved (state, saved) { state.savedAcc = saved },

  showOpenVPNFormAcc (state, show) { state.showAccOpenVPNForm = show },
}

export default {
  state,
  getters,
  actions,
  mutations
}