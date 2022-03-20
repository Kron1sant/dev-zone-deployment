<template>
  <v-overlay
    :value="dialog"
    opacity=0.5
  >
    <v-card
      light
      width="500"
    >
      <v-card-title>
        Ключ OpenVPN
      </v-card-title>

      <v-card-subtitle class="subtitle-1">
        <span class="primary--text">{{ DevAccCurrentItem.username }}</span>
        <span class="caption accent--text "> / {{ DevAccCurrentItem.surname }} {{ DevAccCurrentItem.name }} {{ DevAccCurrentItem.patronomic }}</span>
      </v-card-subtitle>

      <v-card-text>
        Доступ к сети Dev-зоны выполняется по технологии OpenVPN.
        Для подключения нужен набор сертификатов (ключей) упакованных в файл <b>*.ovpn</b>
      </v-card-text>

      <v-card-text v-if="DevAccCurrentItem.hasOVPNCert" class="subtitle-2 green--text">Ранее ключ OpenVPN был выдан</v-card-text>
      <v-card-text v-else class="subtitle-2 red--text">Ранее ключ OpenVPN не выдавался</v-card-text>

      <v-btn
        color="orange darken-1"
        text
        block
        @click="getOpenVPNKey"
      >
        <v-icon>mdi-key</v-icon>
        <span class="mx-3">{{ DevAccCurrentItem.hasOVPNCert ? "Получить ключ заново" : "Получить новый ключ"}}</span>
        <v-icon>mdi-key</v-icon>
      </v-btn>  

      <v-card-actions>
        
        <v-spacer></v-spacer>
        <v-btn
          color="blue darken-1"
          text
          @click="close"
        >
          Отмена
        </v-btn>
      </v-card-actions>
    </v-card> 
  </v-overlay>
</template>

<script>
import axios from 'axios'
import { mapGetters, mapMutations, mapActions } from 'vuex'

export default {
  name: "OpenVPNKeysPage",
  data() {
    return {}
    
  },
  computed: {
    ...mapGetters(["DevAccOpenVPNShowForm", "DevAccCurrentItem"]),

    dialog: {
      get () { return this.DevAccOpenVPNShowForm },
      set (val) {
        // Save global form showing state
        this.showOpenVPNFormAcc(val) 
      }
    },
  },

  methods: {
    ...mapMutations(["showOpenVPNFormAcc", "setDevAccHasBeenSaved"]),
    ...mapActions(["ShowError"]),

    close () {
      this.dialog = false
    },

    getOpenVPNKey () {
      axios.post("auth/openvpnkey", this.DevAccCurrentItem)
      .then( res => {
        var FileSaver = require('file-saver');
        var blob = new Blob([atob(res.data)], {type: "application/json;charset=utf-8"});
        FileSaver.saveAs(blob, this.DevAccCurrentItem.username + ".ovpn");
      })
      .then( () => {
        this.DevAccCurrentItem.hasOVPNCert = true
        this.setDevAccHasBeenSaved(true)
      })
      .catch( err => {
        this.ShowError({
            message: "Ошибка получения ключа OpenVPN для " + this.DevAccCurrentItem.username,
            details: (err.response !== undefined && err.response.status == 400) ? err.response.data.error : err.message,
            actionDescr: "Вернуться к списку аккаунтов",
            actionPath: "/devaccounts"
        })      
      })
      .finally(this.close())
    }
  }
}
</script>