<template>
  <v-dialog v-model="dialog" max-width="700px">
    <!-- Account editor -->
    <v-form
      ref="form"
      v-model="formValid"
      @submit.prevent="saveAndClose"
    >
      <v-card>
        <v-card-title>
          <span class="text-h5">
            {{ formTitle }} <span class="text-h4 primary--text" v-if="modified" title="Есть несохраненные изменения">*</span>
          </span>          
        </v-card-title>
      
        <!-- Editor fields -->
        <v-card-text>
          <v-container>                    
            <v-row>
              <v-col cols="12" sm="4" md="1">
                <v-text-field
                  v-model.number="editedUser.id"
                  label="ID"
                  disabled
                ></v-text-field>
              </v-col>
              <v-col cols="12" sm="8" md="5">
                <v-text-field
                  required
                  :autofocus="AppUserNew"
                  :rules="[rules.required, rules.username]"
                  v-model="editedUser.username"
                  @change="modified=true"
                  :readonly="!AppUserNew"
                >
                  <template #label>
                    Логин<span class="red--text"><strong v-if="AppUserNew"> *</strong></span>
                  </template>
                </v-text-field>
              </v-col>
              <v-col cols="12" sm="12" md="6">
                <v-text-field
                  required
                  :autofocus="!AppUserNew"
                  :rules="[rules.required, rules.email]"
                  type="email"
                  v-model="editedUser.email"
                  @change="modified=true"
                >
                  <template #label>
                    Почта<span class="red--text"><strong> *</strong></span>
                  </template>
                </v-text-field>
              </v-col>
            </v-row>

            <v-row>
              <v-col cols="12" sm="4" md="4">
                <v-checkbox
                  :readonly="!IsAdmin"
                  v-model="editedUser.isAdmin"
                  label="Администратор"
                  color="red"
                  hide-details
                  @change="modified=true"
                ></v-checkbox>
              </v-col>
              <v-col cols="12" sm="8" md="8">
                <v-select
                  v-model="selectDevAcc"
                  :items="devAccounts"
                  label="Привязанный аккаунт разработчика"
                  @click="getDevAccounts"
                  @change="modified=true"
                ></v-select>
              </v-col>
            </v-row>

            <v-row v-if="AppUserNew">
              <v-col cols="12" sm="12" md="6">
                <v-text-field
                  required
                  :rules="[rules.required, rules.password]"
                  type="password"
                  label="Пароль"
                  v-model="newPass.first"
                  autocomplete="new-password"
                ></v-text-field>
              </v-col>
              <v-col cols="12" sm="12" md="6">
                <v-text-field
                  required
                  :rules="[rules.required, rules.password, rules.equalPassConfirm]"
                  type="password"
                  label="Подтверждение"
                  v-model="newPass.second"
                  autocomplete="new-password"
                ></v-text-field>
              </v-col>
            </v-row>

          </v-container>
        </v-card-text>
        <!-- End Editor fields -->

        <!-- Change password dialog -->
        <v-dialog v-model="dialogSetPass" max-width="500">
          <v-card>
            <v-card-title class="deep-orange--text subtitle-1">Укажите новый пароль</v-card-title>

            <v-card-text>
              <v-form
                ref="formNewPass"
              >
                <v-text-field
                  readonly
                  solo
                  dark
                  v-model="newPass.username"
                ></v-text-field>
                <v-text-field
                  autofocus
                  type="password"
                  label="Пароль"
                  v-model="newPass.first"
                  autocomplete="new-password"
                ></v-text-field>
                <v-text-field
                  type="password"
                  label="Подтверждение"
                  v-model="newPass.second"
                  autocomplete="new-password"
                ></v-text-field>
              </v-form>
              <p v-if="!!newPass.error" class="red--text" >{{ newPass.error }}</p>
            </v-card-text>

            <v-card-actions>
              <v-spacer></v-spacer>
              <v-btn color="blue darken-1" text @click="dialogSetPass = false">Отмена</v-btn>
              <v-btn color="blue darken-1" text @click="setPasswordConfirm">Установить</v-btn>
              <v-spacer></v-spacer>
            </v-card-actions>
            
          </v-card>
        </v-dialog>
        <!-- End Change password dialog -->

        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn
            color="deep-orange darken-1"
            text
            @click="showDialogNewPass"
            v-if="!AppUserNew"
          >
            Сменить пароль
          </v-btn>

          <v-spacer></v-spacer>

          <v-btn
            color="blue darken-1"
            text
            @click="close"
          >
            {{ modified ? "Отмена" : "Закрыть"}}
          </v-btn>
          <v-btn
            color="blue darken-1"
            text
            type="submit"
            :disabled="!formValid || !modified"
          >
            Сохранить
          </v-btn>

        </v-card-actions>
      </v-card>
    </v-form>
    <!-- End Account editor -->
    <v-snackbar
      v-model="snackbar"
    >
      {{ snackbarText }}
      <template v-slot:action="{ attrs }">
        <v-btn
          color="white"
          text
          v-bind="attrs"
          @click="snackbar = false"
        >
          <v-icon>mdi-close</v-icon>
        </v-btn>
      </template>
    </v-snackbar>
  </v-dialog>
  <!-- End Edit form -->
</template>

<script>
import axios from 'axios'
import { mapGetters, mapMutations, mapActions } from 'vuex'

export default {
  components: {},
  data() {
    return {
      formValid: true,
      modified: false,
      dialogSave: false,
      dialogSetPass: false,
      rules: {
        required: value => !!value || "Обязательное поле",
        username: value => {
          const pattern = /^[a-zA-Z][a-zA-Z0-9_]*$/
          return pattern.test(value) || "Только английские буквы и цифры, допускается символ подчеркивания"
        },
        email: value => {
          const pattern = /^(([^<>()[\]\\.,;:\s@"]+(\.[^<>()[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/
          return pattern.test(value) || "Неверный формат e-mail."
        },
        password: value => {
          return value.length >= 8 || "Длина пароля должна быть не менее 8 символов"
        },
        equalPassConfirm: value => {
          return value == this.newPass.first || "Пароль и подтверждение должны совпадать"
        }
      },
      newPass: {
        username: "",
        first: "",
        second: "",
        error: ""
      },
      snackbar: false,
      snackbarText: "",
      devAccounts: [],
      selectDevAcc: null
    }
  },
  
  computed: {
    ...mapGetters(["AppUserShowForm", "AppUserCurrentItem", "AppUserNew", "IsAdmin"]),

    dialog: {
      get () { return this.AppUserShowForm },
      set (val) {
        // Mark form as unmodified when closing 
        if (!val) this.modified = false
        // Save global form showing state
        this.showFormAppUser(val) 
      }
    },

    editedUser () { 
      // Access by reference to the object in the global store
      return this.AppUserCurrentItem 
    },

    formTitle () {
      if (this.AppUserNew)
        return "Новый пользователь" 
      else
        return "Редактирование пользователя"
    }
    
  },

  watch: {
    dialog (val) {
      // When open fill selectDevAcc
      if (val && !this.selectDevAcc) {
        var devAccItem = null
        if (this.editedUser.hasDevAccount) {
          devAccItem = {
            text: this.editedUser.devAccountUsername + " (" + this.editedUser.devAccountId + ")",
            value: this.editedUser.devAccountId
          }
        } else {
          devAccItem = {
            text: "<нет аккаунта>",
            value: null
          }
        }
        this.devAccounts = [devAccItem]
        this.selectDevAcc = devAccItem.value
      }

      // Close when dialog == false
      val || this.close()
    },

    dialogSetPass (val) {
      if (!val) {
        this.$refs.formNewPass.reset()
        this.newPass.error = ""
      }
    }
  },

  methods: {
    ...mapMutations(["showFormAppUser", "setAppUserHasBeenSaved"]),
    ...mapActions(["ShowError"]),

    saveAndClose () {
      try {
        this.save()
      } finally {
        this.close()
      }
    },

    save () {
      let action = "edit"
      if (this.AppUserNew) action = "add"

      this.editedUser.hasDevAccount = this.selectDevAcc !== null
      this.editedUser.devAccountId = (this.selectDevAcc !== null) ? this.selectDevAcc : 0

      axios.post("auth/users/" + action, this.editedUser)
        .then( (res) => {
          // Get id, which assigned at the backend
          this.editedUser.id = res.data.id
          this.setAppUserHasBeenSaved(true)
          this.modified = false
        })
        .then( () => {if (this.AppUserNew) this.setNewPassword() } )
        .catch( err => { 
          this.ShowError({
              message: "Ошибка создания/изменения Пользователя",
              details: (err.response !== undefined && err.response.status == 400) ? err.response.data.error : "",
              actionDescr: "Вернуться к списку пользователей",
              actionPath: "/users"
          })
        })
    },

    close () {
      this.dialog = false
      if (this.modified)
        this.modified = false

      this.snackbarText = ""
      this.selectDevAcc = null
    },

    showDialogNewPass () {
      this.newPass.username = this.editedUser.username
      this.dialogSetPass = true
    },

    setPasswordConfirm () {
      if (this.newPass.first !== this.newPass.second) {
        this.newPass.error = "Пароль и подтверждение должны совпадать"
        return
      }

      if (!this.newPass.first || this.newPass.first.length < 8) {
        this.newPass.error = "Длина пароля должна быть не менее 8 символов"
        return
      }

      this.setNewPassword()
    },

    setNewPassword () {
      axios.post("auth/users/setPassword", {
        id: this.editedUser.id,
        password: this.newPass.first
      })
        .then( () => {
          this.showSnackBar("Пароль изменен")
        })
        .catch( err => { 
          this.ShowError({
              message: "Ошибка создания/изменения Пользователя",
              details: (err.response !== undefined && err.response.status == 400) ? err.response.data.error : "",
              actionDescr: "Вернуться к списку пользователей",
              actionPath: "/users"
          })
        })
        .finally( this.dialogSetPass = false )
    },

    showSnackBar (mes) {
      this.snackbarText = mes
      this.snackbar = true
    },

    async getDevAccounts () {
      this.devAccounts = [{
        text: "<нет аккаунта>",
        value: null
      }]
      axios.get("auth/accounts")
      .then( res => {
        res.data.forEach( el => {
          this.devAccounts.push({
            text: el.username + " (" + el.id + ")",
            value: el.id
          })
        })
      })
      .catch( err => { if (err) { /* do nothing */ }} )
    }

  }
}
</script>