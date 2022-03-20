import axios from "axios"

const state = {
  showVMForm: false,
  VMUser: false,
  readOnlyVM: false,
  savedVM: false,
  editedVMIndex: -1,
  editedVM: {
    id: "",
    name: "",
    status: "",
    description: "",
    params: "",
    hasDevAccount: false,
    devAccountId: 0,
    devAccountUsername: ""
  },
  defaultVM: {
    id: "",
    name: "",
    status: "",
    description: "",
    params: "",
    hasDevAccount: false,
    devAccountId: 0,
    devAccountUsername: ""
  }
}

const getters = {
  VMShowForm: state => state.showVMForm,
  VMNew: state => state.newVM,
  VMReadOnly: state => state.readOnlyVM,
  VMCurrentItem: state => state.editedVM,
  VMCurrentIndex: state => state.editedVMIndex,
  VMDefault: state => state.defaultVM,
  VMHasBeenSaved: state => state.savedVM,
}

const actions = {
  async OpenVMDialog ({ commit, VMDefault }, vmParams) {
    if (vmParams.editedItem === null)
      commit("setVMCurrentItem", VMDefault)
    else
      // Obtain VM name 
      if (vmParams.editedItem.hasDevAccount) {
        vmParams.editedItem.devAccountUsername = await getDevAccUsername(vmParams.editedItem.devAccountId)
      }
      commit("setVMCurrentItem", vmParams.editedItem)

    commit("setNewVM", vmParams.new)
    commit("setVMCurrentIndex", vmParams.currentIndex) 
    commit("showFormVM", true)
  }
}

const mutations = {
  setVMCurrentItem (state, vm) {state.editedVM = Object.assign({}, vm) },
  setVMCurrentIndex (state, currentIndex) { state.editedVMIndex = currentIndex },
  setVMUser (state, newAppUser) { state.newAppUser = newAppUser },
  setReadOnlyVM (state, readOnlyAcc) { state.readOnlyAcc = readOnlyAcc },
  showFormVM (state, show) { state.showVMForm = show },
  setVMHasBeenSaved (state, saved) { state.savedVM = saved },
}

async function getDevAccUsername(devAccountId) {
  let devAccUsername = ""
  try {
    var res = await axios.get("auth/accounts?accountid="+devAccountId)
    devAccUsername = res.data[0].username
  } catch (err) { err => { if (err) { /* do nothing */ }} }
  
  return devAccUsername
}

export default {
  state,
  getters,
  actions,
  mutations
}