<template>
  <v-app>
    <AppBar/>
    <Navigation/>

    <v-main>
      <v-container fluid>
        <router-view></router-view>
      </v-container>
    </v-main>

    <v-footer app>
      <!-- -->
    </v-footer>
  </v-app>
</template>

<script>
import Navigation from '@/components/Navigation.vue'
import AppBar from '@/components/AppBar.vue'

export default {
  name: 'App',
  components: {
    Navigation,
    AppBar
  },
  data () {
    return {}
  },
  created: function() {
    // When SPA created we need define wether current user has valid token or not
    this.$store.dispatch("CheckAuth")
  }
}
</script>
