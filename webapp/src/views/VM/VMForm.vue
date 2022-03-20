<template>
  <v-dialog v-model="dialog" max-width="700px">
    <!-- Account editor -->
    <v-form
      ref="form"
      v-model="formValid"
      @submit.prevent="saveAndClose"
      :readonly="VMReadOnly"
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
              <v-col cols="12" sm="4" md="4">
                <v-text-field
                  v-model="editedVM.id"
                  label="ID"
                  disabled
                ></v-text-field>
              </v-col>
              <v-col cols="12" sm="8" md="8">
                <v-text-field
                  required
                  :autofocus="VMNew"
                  :rules="[rules.required, rules.name]"
                  v-model="editedVM.name"
                  @change="modified=true"
                  :readonly="!VMNew"
                >
                  <template #label>
                    Имя ВМ<span class="red--text"><strong v-if="VMNew"> *</strong></span>
                  </template>
                </v-text-field>
              </v-col>
            </v-row>

            <v-row>
              <v-col cols="12" sm="6" md="6">
                <v-text-field
                  v-model="editedVM.description"
                  @change="modified=true"
                  :readonly="!VMNew"
                  label="Описание"
                ></v-text-field>
              </v-col>
              <v-col cols="12" sm="6" md="6">
                <v-text-field
                  v-model="editedVM.params"
                  @change="modified=true"
                  :readonly="!VMNew"
                  label="Параметры"
                ></v-text-field>
              </v-col>
            </v-row>
            
            <v-row>
              <v-col cols="12" sm="4" md="4">
                <v-text-field
                  readonly
                  type="status"
                  v-model="editedVM.status"
                  label="Статус"
                ></v-text-field>
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
          </v-container>
        </v-card-text>
        <!-- End Editor fields --> 

        <v-card-actions>
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
        name: value => {
          const pattern = /^[a-zA-Z][a-zA-Z0-9-_]*$/
          return pattern.test(value) || "Только английские буквы и цифры, допускается символ подчеркивания и дефис"
        }
      },
      snackbar: false,
      snackbarText: "",
      devAccounts: [],
      selectDevAcc: null
    }
  },
  
  computed: {
    ...mapGetters(["VMShowForm", "VMCurrentItem", "VMNew", "VMReadOnly", "IsAdmin"]),

    dialog: {
      get () { return this.VMShowForm },
      set (val) {
        // Mark form as unmodified when closing 
        if (!val) this.modified = false
        // Save global form showing state
        this.showFormVM(val) 
      }
    },

    editedVM () { 
      // Access by reference to the object in the global store
      return this.VMCurrentItem 
    },

    formTitle () {
      if (this.VMNew)
        return "Новая Виртуальная машина" 
      else
        return "Редактирование Виртуальной машины"
    }
    
  },

  watch: {
    dialog (val) {
      // When open fill selectDevAcc
      if (val && !this.selectDevAcc) {
        var devAccItem = null
        if (this.editedVM.hasDevAccount) {
          devAccItem = {
            text: this.editedVM.devAccountUsername + " (" + this.editedVM.devAccountId + ")",
            value: this.editedVM.devAccountId
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
    }
  },

  methods: {
    ...mapMutations(["showFormVM", "setVMHasBeenSaved"]),
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
      if (this.VMNew) action = "add"

      this.editedVM.hasDevAccount = this.selectDevAcc !== null
      this.editedVM.devAccountId = (this.selectDevAcc !== null) ? this.selectDevAcc : 0

      axios.post("auth/vm/" + action, this.editedVM)
        .then( (res) => {
          // Get id, which assigned at the backend
          this.editedVM.id = res.data.id
          this.setVMHasBeenSaved(true)
          this.modified = false
        })
        .catch( err => { 
          this.ShowError({
              message: "Ошибка создания/изменения ВМ",
              details: (err.response !== undefined && err.response.status == 400) ? err.response.data.error : "",
              actionDescr: "Вернуться к списку ВМ",
              actionPath: "/virtualmachines"
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