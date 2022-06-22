<template>
  <div class="login-con">
    <Form class="login-form" @submit="onSubmit">
      <h2>Hotspot Dashboard</h2>
      <br/>
      <CellGroup inset>
        <Field
          v-model="password"
          type="password"
          name="password"
          label="Password"
          placeholder="Password"
          :rules="[{ required: true, message: 'Input Password' }]"
        />
      </CellGroup>
      <div style="margin: 16px;">
        <Button round block type="primary" native-type="submit">
          Enter
        </Button>
      </div>
    </Form>
  </div>
</template>

<script setup>

import { Button, Form, CellGroup, Cell, Field } from 'vant'
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import * as api from './api'
import { AuthToken } from './api/auth';

const router = useRouter()
const password = ref('')

async function onSubmit() {
  const r = await api.login(password.value)
  console.log(r)
  AuthToken.set(r)
  router.replace('/')
}

</script>

<style lang="less" scoped>
.login-con {
  display: flex;
  flex: 1;
//   width: 100vw;
  height: 100vh;
  align-items: center;
  justify-content: center;
}
.login-form {
  flex: 1;
  min-width: 400px;
  max-width: 600px;
  text-align: center;
}
</style>