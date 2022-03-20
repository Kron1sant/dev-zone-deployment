<template>
  <v-app-bar 
    app 
    clipped-left 
    dense dark
    color="accent-2"
  >
    <v-app-bar-nav-icon v-if="!permanent" @click.stop="drawer = !drawer" />
    <v-toolbar-title><v-btn to="/">DevZone</v-btn></v-toolbar-title>

    <v-spacer></v-spacer>

    <a v-if="IsAuthenticated" @click="openAppUser" class="teal--text subtitle-1 mx-2">
       <v-icon color="teal" small>mdi-account</v-icon>
      {{ CurrentUser }}
    </a>
        
    <v-btn v-if="!IsAuthenticated" to="/login">Войти</v-btn>
    <v-btn v-if="IsAuthenticated" @click="logout">Выйти</v-btn>
  </v-app-bar>
</template>

<script>
import axios from 'axios'
import router from '@/router'
import { mapGetters, mapMutations } from 'vuex'

export default {
  name: 'AppBar',
  components: {},
  data () {
    return {
      permanent:true
    }
  },
  computed: {
      ...mapGetters(["IsAuthenticated", "CurrentUser"])
  },
  methods: {
    ...mapMutations(["unauthorize"]),

    logout () {
      axios.post("/logout")
        .then( this.unauthorize() )
        .catch( err => console.error(err) )
      // Go home after logged out
      this.$router.push('/').catch(()=>{})
    },

    openAppUser () { router.push("/users").catch( (err) => {if (err) () => {}} ) }
  },
}
</script>
