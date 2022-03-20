<template>
  <div id="errors">
    <v-alert
      prominent
      type="success"
      v-if="noErrors"
    >
      <v-row align="center">
        <v-col class="grow">
          Нет новых ошибок
        </v-col>
      </v-row>
    </v-alert>
    
    <template v-for="(err, index) in errors">
      <v-alert
        prominent
        type="error"
        v-bind:key="index"
      >
        <v-row align="center">
          <v-col class="grow">
            <p class="h5">{{ err.message }}</p>
            <p class="caption">{{ err.details }}</p>
          </v-col>
          <v-col class="shrink">
            <v-btn :to="err.actionPath">{{ err.actionDescr }}</v-btn>
          </v-col>
        </v-row>
      </v-alert>
    </template>
    
  </div>
</template>

<script>
export default {
  name: "ErrorsPage",
  components: {},
  data () {
    return {}
  },
  beforeRouteLeave (to, from, next) {
    // Clear errors if leaving Error page
   if (to !== from) 
      this.$store.commit("dropErrors")
    next()
  },
  computed: {
    errors () { return this.$store.getters.GetErrors },
    noErrors () { return this.$store.getters.ErrorsIsEmpty }
  }
}
</script>