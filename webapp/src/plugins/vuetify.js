import Vue from 'vue'
import Vuetify from 'vuetify'
import 'vuetify/dist/vuetify.min.css'
import colors from 'vuetify/lib/util/colors'

Vue.use(Vuetify, {
  theme: {
    themes: {
      light: {
        primarytext: colors.purple,
        secondary: colors.grey.darken1,
        accent: colors.shades.black,
        error: colors.red.accent3,
      },
      dark: {
        primarytext: colors.blue.lighten3,
      },
    },
  }
})

const opts = {}

export default new Vuetify(opts)