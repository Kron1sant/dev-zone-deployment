<template>
  <v-container>
    <v-card
      max-width="600"
      elevation=8
      class="mx-auto mt-12"
    >
      <v-card-title class="justify-center">Вход в систему</v-card-title>

      <v-form
        ref="form"
        v-model="valid"
        lazy-validation
        class="pa-6"
        @submit.prevent="login"
      >
        <v-text-field
          v-model="username"
          :rules="nameRules"
          label="Имя пользователя"
          autocomplete="username"
          required
        ></v-text-field>

        <v-text-field
          v-model="password"
          :rules="passRules"
          type="password"
          autocomplete="current-password"
          label="Пароль"
          required
        ></v-text-field>

        <v-flex>
        <v-btn
          :disabled="!valid"
          color="success"
          class="mr-4"
          type="submit"
        >
          Войти
        </v-btn>
          <p v-if="!!authError" class="red--text mt-5 d-flex justify-center">{{ authError }}</p>
        </v-flex>
      </v-form>
    </v-card>
  </v-container>
</template>

<script>
import { mapMutations } from 'vuex'
import axios from 'axios'

export default {
  name: "LoginPage",
  components: {
  },
  data() {
    return {
      valid: true,
      username: "",
      nameRules: [
        v => !!v || 'Укажите имя пользователя'
      ],
      password: '',
      passRules: [
        v => !!v || 'Укажите пароль',
      ],
      authError: ""
    }
  },
  methods: {
    ...mapMutations(["authorize"]),
    login() {
      if (!this.$refs.form.validate()) return
      this.valid = false

      axios.post("/login", {
        username: this.username,
        password: this.password 
        })
      .then( (res) => {
        // ToDo wait handler to expire and reclaim new cred
        this.authorize(parseJwt(res.data.token)) 
        this.authError = ""
        this.$router.push("/devaccounts")
      })
      .catch( () => {
        this.authError = "Некорректное имя пользователя и/или пароль"
        this.valid = true
      })
    }
  },
}

function parseJwt (token) {
    var base64Url = token.split('.')[1];
    var base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
    var jsonPayload = decodeURIComponent(atob(base64).split('').map(function(c) {
        return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2);
    }).join(''));

    return JSON.parse(jsonPayload);
}
</script>