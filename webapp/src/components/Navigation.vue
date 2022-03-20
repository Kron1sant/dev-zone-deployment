<template>
    <v-navigation-drawer 
      app
      v-model="drawer" 
      clipped
      permanent
    >
      <v-list
        dense
        nav  
      >
        <v-list-item
          v-for="link in links"
          :key="link.id"
          :to="link.path"
        >
          <v-list-item-content>
            <v-list-item-title class="subtitle-2">{{ link.label }}</v-list-item-title>
          </v-list-item-content>
        </v-list-item>  
      </v-list>
      <template v-slot:append>
          <v-list nav>
            <v-list-item to="/users" class="blue-grey darken-1 white--text subtitle-2"> 
              Список пользователей
            </v-list-item>
          </v-list>
        </template>
    </v-navigation-drawer>
</template>

<script>
import { routes } from '@/router'
import { mapGetters } from 'vuex'
export default {
  name: 'NavigationBar',
  components: {},
  data() {
    return {
      permanent: true,
      drawer: true
    }
  },

  computed: {
    ...mapGetters(["IsAuthenticated"]),

    links () {
      // get links from routes
      let l = []
      routes.forEach(el => {
        var meta = Object.assign({ 
          label: "",
          navigation: false,
          requiresAuth: false 
        }, el.meta)

        if (!meta.navigation) return // skip home page and login
        if (meta.requiresAuth && !this.IsAuthenticated) return // skip the restricted access page if the user is not logged in
        l.push({ id: el.name, label: meta.label, path: el.path, needAuth: meta.requiresAuth })})
        return l
      }
    }
}

</script>