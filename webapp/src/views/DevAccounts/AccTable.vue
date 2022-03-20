<template>
  <v-container fluid max-width="1000">
    <v-data-table
      :items="devaccounts"
      :headers="headers"
      :search="search"
      class="elevation-2"
    >
      <!-- Top table -->
      <template v-slot:top>
        <v-toolbar flat>
          <v-toolbar-title>Аккаунты разработчиков</v-toolbar-title>
          <v-divider
            class="mx-4"
            inset
            vertical
          ></v-divider>

          <v-text-field
            v-model="search"
            append-icon="mdi-magnify"
            label="Поиск"
            single-line
            hide-details
            class="blue lighten-5 rounded-lg"
          ></v-text-field>

          <v-divider class="mx-4" inset vertical></v-divider>

          <v-btn @click="initialize" class="mb-2" >
            <v-icon dark> mdi-refresh </v-icon>
          </v-btn>

          <v-divider class="mx-4" inset vertical></v-divider>

          <!-- Edit form -->
          <v-btn
            color="primary"
            dark
            class="mb-2"
            @click="newItem"
          >
            Добавить
          </v-btn>

          <AccForm />
          <!-- End Edit form -->

          <!-- Close dialog question -->
          <v-dialog v-model="dialogDelete" max-width="450">
            <v-card>
              <v-card-title class="subtitle-1">Вы уверены что хотите удалить аккаунт разработчика?</v-card-title>
              <v-card-subtitle class="subtitle-2 red--text">{{ username }}</v-card-subtitle>
              <v-card-actions>
                <v-spacer></v-spacer>
                <v-btn color="blue darken-1" text @click="closeDelete">Отмена</v-btn>
                <v-btn color="blue darken-1" text @click="deleteItemConfirm">Удалить</v-btn>
                <v-spacer></v-spacer>
              </v-card-actions>
            </v-card>
          </v-dialog>
          <!-- End Close dialog question -->

        </v-toolbar>
      </template> 
      <!-- End Top -->

      <template v-slot:[`item.username`]="{ item }">
        <a @click="showItem(item)">
          {{ item.username }}
        </a>
      </template>

      <!-- Add column with actions -->
      <template v-slot:[`item.actions`]="{ item }">
        <v-icon
          small
          class="mr-2"
          @click="editItem(item)"
          :disabled="!IsAdmin"
        >
          mdi-pencil
        </v-icon>
        <v-icon
          small
          @click="deleteItem(item)"
          :disabled="!IsAdmin"
        >
          mdi-delete
        </v-icon>
      </template>
      <!-- End Add column with actions -->
    </v-data-table>
  </v-container>
</template>

<script>
import axios from 'axios'
import { mapMutations, mapGetters, mapActions } from 'vuex'
import AccForm from './AccForm.vue'

export default {
  name: "DevAccountsPage",
  components: {
    AccForm
  },

  created () {
    this.initialize()
    this.paintTableHeader()
  },
  
  data () {
    return {
      devaccounts: [],
      deletedAccount: null,
      deletedIndex: -1,
      dialogDelete: false,
      search: "",
      headers: [
        { text: "ID", align: "start", value: "id", sortable: false },
        { text: "Логин", value: "username", sortable: false },
        { text: "Фамилия", value: "surname", sortable: false },
        { text: "Имя", value: "name", sortable: false },
        { text: "Отчество", value: "patronomic", sortable: false },
        { text: "Почта", value: "email", sortable: false },
        { text: "Имеет ключ OpenVPN", value: "hasOVPNCert", sortable: false },
        { text: "Комментарий", value: "comment", sortable: false },
        { text: "Действия", value: "actions", sortable: false },
      ],
    }
  },

  computed: {
    ...mapGetters(["DevAccHasBeenSaved", "DevAccDefault", "DevAccCurrentItem", "DevAccCurrentIndex", "DevAccNew", "IsAdmin"]),

    username () { return this.deletedAccount ? this.deletedAccount.username : "" }
  },

  watch: {
    dialogDelete (val) {
      val || this.closeDelete()
    },
    DevAccHasBeenSaved (val) {
      // Waiting notification when devaccount will be changed
      if (val) {
        this.newOrEditItemConfirm()
        this.setDevAccHasBeenSaved(false)
      }
    }
  },

  methods: {
    ...mapMutations(["setDevAccHasBeenSaved"]),
    ...mapActions(["OpenDevAccDialog", "ShowError", "ClearDevAccCurrentItem"]),

    initialize () {
      this.devaccounts = []
      axios.get("auth/accounts")
        .then( res => { 
          // Loop through all elements, convert to devacc format and pushing into array
          res.data.forEach( el => this.pushDevAccountInTable(el) )
        })
        .catch( err => { if (err) { /* do nothing */ }} )
    },

    pushDevAccountInTable (data) {
      var newAcc = Object.assign({}, this.DevAccDefault)
      newAcc = Object.assign(newAcc, data)
      this.devaccounts.push(newAcc)
    },

    newItem () {
      this.OpenDevAccDialog({ 
        new: true,
        readOnly: false,
        editedItem: null,
        currentIndex: -1
      })
    },

    editItem (devaccount) {      
      this.OpenDevAccDialog({
          new: false,
          readOnly: false,
          editedItem: Object.assign({}, devaccount),
          currentIndex: this.devaccounts.indexOf(devaccount) // remember global index of editing account
      })
    },

    showItem (devaccount) {      
      this.OpenDevAccDialog({
          new: false,
          readOnly: true,
          editedItem: Object.assign({}, devaccount),
          currentIndex: this.devaccounts.indexOf(devaccount) // remember global index of editing account
      })
    },

    newOrEditItemConfirm () {
      if (this.DevAccNew) {
        // New account
        this.pushDevAccountInTable(this.DevAccCurrentItem)
      } else {
        // Existing account
        Object.assign(this.devaccounts[this.DevAccCurrentIndex], this.DevAccCurrentItem)
      }
      //this.ClearDevAccCurrentItem() // ToDo - перенести в момент "закрытия" формы редактирования
    },
    
    deleteItem (devaccount) {
      this.deletedIndex = this.devaccounts.indexOf(devaccount)
      this.deletedAccount = Object.assign({}, devaccount)
      this.dialogDelete = true
    },

    deleteItemConfirm () {
      const action = "delete"
      axios.post("auth/accounts/"+action, this.deletedAccount)
        .then( () => {
          this.devaccounts.splice(this.deletedIndex, 1) 
          this.closeDelete()
        })
        .catch( err => { 
          this.ShowError({
              message: "Ошибка удаления Аккаунта разработчика",
              details: (err.response.status == 400 ) ? err.response.data.error : "",
              actionDescr: "Вернуться к списку аккаунтов",
              actionPath: "/devaccounts"
          })
          this.closeDelete()
        })
    },

    closeDelete () {
      this.dialogDelete = false
      this.$nextTick(() => {
        this.deletedAccount = null
        this.deletedIndex = -1
      })
    },
    
    paintTableHeader () {
      this.headers.forEach(header => {
        header.class = "blue-grey darken-3 white--text"
      });
    }
  }
}
</script>