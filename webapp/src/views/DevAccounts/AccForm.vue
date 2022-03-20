<template>
  <v-dialog v-model="dialog" max-width="700px">
    <!-- Account editor -->
    <v-form
      ref="form"
      v-model="formValid"
      :readonly="DevReadOnlyAcc"
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
                  v-model.number="editedAccount.id"
                  label="ID"
                  disabled
                ></v-text-field>
              </v-col>
              <v-col cols="12" sm="8" md="5">
                <v-text-field
                  required
                  autofocus
                  :rules="[rules.required, rules.username]"
                  v-model="editedAccount.username"
                  @change="modified=true"
                >
                  <template #label>
                    Логин<span class="red--text"><strong> *</strong></span>
                  </template>
                </v-text-field>
              </v-col>
              <v-col cols="12" sm="12" md="6">
                <v-text-field
                  required
                  :rules="[rules.required, rules.email]"
                  type="email"
                  v-model="editedAccount.email"
                  @change="modified=true"
                >
                  <template #label>
                    Почта<span class="red--text"><strong> *</strong></span>
                  </template>
                </v-text-field>
              </v-col>
            </v-row>
            
            <v-row>
              <v-col cols="12" sm="6" md="4">
                <v-text-field
                  v-model="editedAccount.surname"
                  label="Фамилия"
                  @change="modified=true"
                ></v-text-field>
              </v-col>
              <v-col cols="12" sm="6" md="4">
                <v-text-field
                  v-model="editedAccount.name"
                  label="Имя"
                  @change="modified=true"
                ></v-text-field>
              </v-col>
              <v-col cols="12" sm="6" md="4">
                <v-text-field
                  v-model="editedAccount.patronomic"
                  label="Отчество"
                  @change="modified=true"
                ></v-text-field>
              </v-col> 
            </v-row>
            
            <v-row>
              <v-col cols="12" sm="12" md="12">
                <v-textarea
                  outlined
                  clearable
                  v-model="editedAccount.comment"
                  label="Комментарий"
                  @change="modified=true"
                ></v-textarea>
              </v-col>
            </v-row>

          </v-container>
        </v-card-text>
        <!-- End Editor fields -->

        <v-card-actions>
          <v-btn 
            class="white--text"
            :color="editedAccount.hasOVPNCert ? 'green' : 'red darken-1' "
            :disabled="DevAccNew || (!IsAdmin && CurrentDevAccId != editedAccount.id)"
            @click="showOpenVPNForm"
          >
            {{ OpenVPNButtonText }}
          </v-btn>         

          <OpenVPNKeys />

          <v-spacer></v-spacer>
          <v-btn
            color="blue darken-1"
            text
            @click="setReadOnlyAcc(false)"
            v-if="DevReadOnlyAcc && IsAdmin"
          >
            Изменить
          </v-btn>
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
            :disabled="!formValid || DevReadOnlyAcc || !modified"
          >
            Сохранить
          </v-btn>

          <!-- Save dialog question -->
          <v-dialog v-model="dialogSave" max-width="300px">
            <v-card>
              <v-card-title class="subtitle-1"><v-spacer></v-spacer>Сначала сохраните изменения<v-spacer></v-spacer></v-card-title>
              <v-card-actions>
                <v-spacer></v-spacer>
                <v-btn color="blue darken-1" text @click="dialogSave=false">Отмена</v-btn>
                <v-btn color="blue darken-1" text @click="saveItemConfirm">Сохранить</v-btn>
                <v-spacer></v-spacer>
              </v-card-actions>
            </v-card>
          </v-dialog>
          <!-- End Close dialog question -->
        </v-card-actions>
      </v-card>
    </v-form>
    <!-- End Account editor -->

  </v-dialog>
  <!-- En Edit form -->
</template>

<script>
import axios from 'axios'
import { mapGetters, mapMutations, mapActions } from 'vuex'
import OpenVPNKeys from './AccOpenVPNKeys.vue'

export default {
  components: {
    OpenVPNKeys
  },
  data() {
    return {
      formValid: true,
      modified: false,
      dialogSave: false,
      rules: {
        required: value => !!value || 'Обязательное поле',
        username: value => {
          const pattern = /^[a-zA-Z][a-zA-Z0-9_]*$/
          return pattern.test(value) || 'Только английские буквы и цифры, допускается символ подчеркивания'
        },
        email: value => {
          const pattern = /^(([^<>()[\]\\.,;:\s@"]+(\.[^<>()[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/
          return pattern.test(value) || 'Неверный формат e-mail.'
        },
      },
    }
  },
  
  computed: {
    ...mapGetters(["DevAccShowForm", "DevAccCurrentItem", "DevAccNew", "DevReadOnlyAcc", "IsAdmin", "CurrentDevAccId"]),

    dialog: {
      get () { return this.DevAccShowForm },
      set (val) {
        // Mark form as unmodified when closing 
        if (!val) this.modified = false
        // Save global form showing state
        this.showFormAcc(val) 
      }
    },

    editedAccount () { 
      // Access by reference to the object in the global store
      return this.DevAccCurrentItem 
    },

    formTitle () {
      if (this.DevReadOnlyAcc) 
        return "Аккаунт разработчика"
      else if (this.DevAccNew)
        return "Новый аккаунт" 
      else
        return "Редактирование аккаунта"
    },

    OpenVPNButtonText () {
      if (this.editedAccount.hasOVPNCert) 
        return "Ключ OpenVPN выпущен"
      else 
        return "Ключ OpenVPN не выпускался"
    }
  },

  watch: {
    dialog (val) {
      val || this.close()
    }
  },

  methods: {
    ...mapMutations(["showFormAcc", "setDevAccHasBeenSaved", "showOpenVPNFormAcc", "setReadOnlyAcc"]),
    ...mapActions(["ShowError"]),

    validate () {
      this.$refs.form.validate()
    },

    saveAndClose () {
      try {
        this.save()
      } finally {
        this.close()
      }
    },

    save () {
      let action = "edit"
      if (this.DevAccNew) action = "add"
      axios.post("auth/accounts/"+action, this.editedAccount)
        .then( (res) => {
          // Get id, which assigned at the backend
          this.editedAccount.id = res.data.id
          this.setDevAccHasBeenSaved(true)
          this.modified = false
        })
        .catch( err => { 
          this.ShowError({
              message: "Ошибка создания/изменения Аккаунта разработчика",
              details: (err.response !== undefined && err.response.status == 400) ? err.response.data.error : "",
              actionDescr: "Вернуться к списку аккаунтов",
              actionPath: "/devaccounts"
          })
        })
    },

    close () {
      this.dialog = false
      if (this.modified)
        this.modified = false
    },

    showOpenVPNForm () {
      if (this.modified) {
        this.dialogSave = true
        return
      }
      this.showOpenVPNFormAcc(true)
    },

    saveItemConfirm () {
      this.dialogSave = false
      this.save()
      this.showOpenVPNForm()
    }

  }
}
</script>