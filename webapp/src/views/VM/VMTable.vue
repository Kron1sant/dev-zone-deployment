<template>
  <v-container fluid max-width="1000">
    <v-data-table
      :items="vm"
      :headers="headers"
      :search="search"
      class="elevation-2"
    >
      <!-- Top table -->
      <template v-slot:top>
        <v-toolbar flat>
          <v-toolbar-title>Список виртуальных машин</v-toolbar-title>
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
          <v-btn
            color="primary"
            class="mb-2"
            @click="cloudSync"
            :disabled="!IsAdmin"
          >
            Синхронизировать из облака
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

          <VMForm />
          <!-- End Edit form -->
          
          <!-- Close dialog question -->
          <v-dialog v-model="dialogDelete" max-width="450">
            <v-card>
              <v-card-title class="subtitle-1">Вы уверены что хотите удалить Вирутальную машину?</v-card-title>
              <v-card-subtitle class="subtitle-2 red--text">{{ vmName }}</v-card-subtitle>
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

      <template v-slot:[`item.name`]="{ item }">
        <a @click="showItem(item)">
          {{ item.name }}
        </a>
      </template>

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
import VMForm from './VMForm.vue'
import { mapGetters, mapMutations, mapActions } from 'vuex'

export default {
  name: "VMTablePage",
  components: {
    VMForm
  },
  
  created () {
    this.initialize()
    this.paintTableHeader()
  },

  data () {
    return {
      vm: [],
      deletedItem: null,
      deletedIndex: -1,
      dialogDelete: false,
      search: "",
      headers: [
        { text: "ID", align: "start", value: "id", sortable: false },
        { text: "Имя ВМ", value: "name", sortable: false },
        { text: "Привязана к разработчику", value: "hasDevAccount", sortable: false },
        { text: "ID разработчика", value: "devAccountId", sortable: false },
        { text: "Статус", value: "status", sortable: false },
        { text: "Описание", value: "description", sortable: false },
        { text: "Параметры", value: "params", sortable: false },
        { text: "Действия", value: "actions", sortable: false },
      ],
    }
  },

  computed: {
    ...mapGetters(["VMDefault", "VMCurrentItem", "VMCurrentIndex", "VMNew", "VMHasBeenSaved", "IsAdmin"]),

    vmName () { return this.deletedItem ? this.deletedItem.name : "" }
  },

  watch: {
    dialogDelete (val) {
      val || this.closeDelete()
    },
    VMHasBeenSaved (val) {
      // Waiting notification when devaccount will be changed
      if (val) {
        this.newOrEditItemConfirm()
        this.setVMHasBeenSaved(false)
      }
    }
  },

  methods: {
    ...mapMutations(["setVMHasBeenSaved"]),
    ...mapActions(["ShowError", "OpenVMDialog"]),

    initialize () {
      this.vm = []
      axios.get("auth/vm")
        .then( res => { 
          // Loop through all elements, convert to vm format and pushing into array
          res.data.forEach( el => this.pushVMInTable(el) )
        })
        .catch( err => { if (err) { /* do nothing */ }} )
    },

    cloudSync () {
      axios.post("auth/vm/update")
        .then( () => {
          this.initialize() 
        })
        .catch( err => { if (err) { /* do nothing */ }} )  
    },

    pushVMInTable (data) {
      var newVM = Object.assign({}, this.VMDefault)
      newVM = Object.assign(newVM, data)
      this.vm.push(newVM)
    },

    newItem () {
      this.OpenVMDialog({ 
        new: true,
        editedItem: null,
        currentIndex: -1
      })
    },

    editItem (item) {      
      this.OpenVMDialog({
          new: false,
          editedItem: Object.assign({}, item),
          currentIndex: this.vm.indexOf(item) // remember global index of editing vm
      })
    },
    
    showItem (item) {      
      this.OpenVMDialog({
          new: false,
          readOnly: true,
          editedItem: Object.assign({}, item),
          currentIndex: this.vm.indexOf(item) // remember global index of editing account
      })
    },

    newOrEditItemConfirm () {
      if (this.VMNew) {
        // New account
        this.pushVMInTable(this.VMCurrentItem)
      } else {
        // Existing account
        Object.assign(this.vm[this.VMCurrentIndex], this.VMCurrentItem)
      }
    },
    
    deleteItem (item) {
      this.deletedIndex = this.vm.indexOf(item)
      this.deletedItem = Object.assign({}, item)
      this.dialogDelete = true
    },

    deleteItemConfirm () {
      const action = "delete"
      axios.post("auth/vm/"+action, this.deletedItem)
        .then( () => {
          this.vm.splice(this.deletedIndex, 1) 
        })
        .catch( err => { 
          this.ShowError({
              message: "Ошибка удаления Виртуальной машины",
              details: (err.response.status == 400 ) ? err.response.data.error : "",
              actionDescr: "Вернуться к списку ВМ",
              actionPath: "/virtualmachines"
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