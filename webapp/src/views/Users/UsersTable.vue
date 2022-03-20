<template>
  <v-container fluid max-width="1000">
    <v-data-table
      :items="users"
      :headers="headers"
      :search="search"
      class="elevation-2"
    >
      <!-- Top table -->
      <template v-slot:top>
        <v-toolbar flat>
          <v-toolbar-title>Список пользователей</v-toolbar-title>
          <v-divider class="mx-4" inset vertical></v-divider>
          <!-- Search field -->
          <v-text-field
            v-model="search"
            append-icon="mdi-magnify"
            label="Поиск"
            single-line
            hide-details
            class="blue lighten-5 rounded-lg"
          ></v-text-field>
          <!-- End Search field -->
          <v-divider class="mx-4" inset vertical></v-divider>

          <v-btn @click="initialize" class="mb-2" >
            <v-icon dark> mdi-refresh </v-icon>
          </v-btn>

          <v-divider class="mx-4" inset vertical></v-divider>
          
          <!-- Edit form -->
          <v-btn
            color="primary"
            class="mb-2"
            @click="newItem"
            :disabled="!IsAdmin"
          >
            Добавить
          </v-btn>

          <AppUserForm />
          <!-- End Edit form -->
          

          <!-- Close dialog question -->
          <v-dialog v-model="dialogDelete" max-width="450">
            <v-card>
              <v-card-title class="subtitle-1">Вы уверены что хотите удалить Пользователя?</v-card-title>
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

      <!-- Add column with actions -->
      <template v-slot:[`item.actions`]="{ item }">
        <v-icon
          small
          class="mr-2"
          @click="editItem(item)"
        >
          mdi-pencil
        </v-icon>
        <v-icon
          small
          @click="deleteItem(item)"
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
import AppUserForm from './AppUserForm.vue'
import { mapGetters, mapMutations, mapActions } from 'vuex'

export default {
  name: "UsersTablePage",
  components: {
    AppUserForm
  },
  
  created () {
    this.initialize()
    this.paintTableHeader()
  },

  data () {
    return {
      users: [],
      deletedItem: null,
      deletedIndex: -1,
      dialogDelete: false,
      search: "",
      headers: [
        { text: "ID", align: "start", value: "id", sortable: false },
        { text: "Логин", value: "username", sortable: false },
        { text: "Почта", value: "email", sortable: false },
        { text: "Админ", value: "isAdmin", sortable: false },
        { text: "Действия", value: "actions", sortable: false },
      ],
    }
  },

  computed: {
    ...mapGetters(["AppUserDefault", "AppUserCurrentItem", "AppUserCurrentIndex", "AppUserNew", "AppUserHasBeenSaved", "IsAdmin"]),

    username () { return this.deletedItem ? this.deletedItem.username : "" }
  },

  watch: {
    dialogDelete (val) {
      val || this.closeDelete()
    },
    AppUserHasBeenSaved (val) {
      // Waiting notification when devaccount will be changed
      if (val) {
        this.newOrEditItemConfirm()
        this.setAppUserHasBeenSaved(false)
      }
    }
  },

  methods: {
    ...mapMutations(["setAppUserHasBeenSaved"]),
    ...mapActions(["ShowError", "OpenAppUserDialog"]),

    initialize () {
      this.users = []
      axios.get("auth/users")
        .then( res => { 
          // Loop through all elements, convert to devacc format and pushing into array
          res.data.forEach( el => this.pushAppUserInTable(el) )
        })
        .catch( err => { if (err) { /* do nothing */ }} )
    },

    pushAppUserInTable (data) {
      var newAppUser = Object.assign({}, this.DevAccDefault)
      newAppUser = Object.assign(newAppUser, data)
      this.users.push(newAppUser)
    },

    newItem () {
      this.OpenAppUserDialog({ 
        new: true,
        editedItem: null,
        currentIndex: -1
      })
    },

    editItem (user) {      
      this.OpenAppUserDialog({
          new: false,
          editedItem: Object.assign({}, user),
          currentIndex: this.users.indexOf(user) // remember global index of editing account
      })
    },


    newOrEditItemConfirm () {
      if (this.AppUserNew) {
        // New account
        this.pushAppUserInTable(this.AppUserCurrentItem)
      } else {
        // Existing account
        Object.assign(this.users[this.AppUserCurrentIndex], this.AppUserCurrentItem)
      }
    },
    
    deleteItem (user) {
      this.deletedIndex = this.users.indexOf(user)
      this.deletedItem = Object.assign({}, user)
      this.dialogDelete = true
    },

    deleteItemConfirm () {
      const action = "delete"
      axios.post("auth/users/"+action, this.deletedItem)
        .then( () => {
          this.users.splice(this.deletedIndex, 1) 
        })
        .catch( err => { 
          this.ShowError({
              message: "Ошибка удаления Пользователя",
              details: (err.response.status == 400 ) ? err.response.data.error : "",
              actionDescr: "Вернуться к списку пользователей",
              actionPath: "/users"
          })
        })
        .finally( this.closeDelete() )
    },

    closeDelete () {
      this.dialogDelete = false
      this.$nextTick(() => {
        this.deletedItem = null
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