import axios from "axios"

const state = {
  showAppUserForm: false,
  newAppUser: false,
  // readOnlyAcc: false,
  savedAppUser: false,
  editedAppUserIndex: -1,
  editedAppUser: {
    id: 0,
    username: "",
    email: "",
    password: "",
    isAdmin: false,
    hasDevAccount: false,
    devAccountId: 0,
    devAccountUsername: ""
  },
  defaultAppUser: {
    id: 0,
    username: "",
    email: "",
    password: "",
    isAdmin: false,
    hasDevAccount: false,
    devAccountId: 0,
    devAccountUsername: ""
  }
}

const getters = {
  AppUserShowForm: state => state.showAppUserForm,
  AppUserNew: state => state.newAppUser,
  // DevReadOnlyAcc: state => state.readOnlyAcc,
  AppUserCurrentItem: state => state.editedAppUser,
  AppUserCurrentIndex: state => state.editedAppUserIndex,
  AppUserDefault: state => state.defaultAppUser,
  AppUserHasBeenSaved: state => state.savedAppUser,
}

const actions = {
  async OpenAppUserDialog ({ commit, AppUserDefault }, appUserParams) {
    if (appUserParams.editedItem === null)
      commit("setAppUserCurrentItem", AppUserDefault)
    else
      // Obtain devacc name 
      if (appUserParams.editedItem.hasDevAccount) {
        appUserParams.editedItem.devAccountUsername = await getDevAccUsername(appUserParams.editedItem.devAccountId)
      }
      commit("setAppUserCurrentItem", appUserParams.editedItem)

    commit("setNewAppUser", appUserParams.new)
    commit("setAppUserCurrentIndex", appUserParams.currentIndex) 
    commit("showFormAppUser", true)
  }
}

const mutations = {
  setAppUserCurrentItem (state, appUser) {state.editedAppUser = Object.assign({}, appUser) },
  setAppUserCurrentIndex (state, currentIndex) { state.editedAppUserIndex = currentIndex },
  setNewAppUser (state, newAppUser) { state.newAppUser = newAppUser },
  // setReadOnlyAcc (state, readOnlyAcc) { state.readOnlyAcc = readOnlyAcc },
  showFormAppUser (state, show) { state.showAppUserForm = show },
  setAppUserHasBeenSaved (state, saved) { state.savedAppUser = saved },
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