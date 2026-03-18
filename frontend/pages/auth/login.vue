<template>
  <div class="auth-page">
    <div class="auth-card">
      <div class="auth-card__header">
        <n-icon size="32" color="#0d6efd"><TicketOutline /></n-icon>
        <h1>АИС Успеваемость</h1>
        <p>Войдите в систему</p>
      </div>
      <n-form ref="formRef" :model="form" :rules="rules" @submit.prevent="submit">
        <n-form-item path="email" label="Email">
          <n-input v-model:value="form.email" placeholder="user@university.ru" size="large" />
        </n-form-item>
        <n-form-item path="password" label="Пароль">
          <n-input
            v-model:value="form.password"
            type="password"
            show-password-on="click"
            placeholder="••••••••"
            size="large"
          />
        </n-form-item>
        <n-alert v-if="error" type="error" :title="error" style="margin-bottom: 16px" />
        <n-button type="primary" size="large" block :loading="loading" @click="submit">
          Войти
        </n-button>
      </n-form>
      <div class="auth-card__footer">
        Нет аккаунта? <NuxtLink to="/auth/register">Зарегистрироваться</NuxtLink>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { NForm, NFormItem, NInput, NButton, NAlert, NIcon } from 'naive-ui'
import type { FormRules } from 'naive-ui'
import { TicketOutline } from '@vicons/ionicons5'
import { authApi } from '~/api/auth'

definePageMeta({ middleware: 'guest', layout: false })

const auth = useAuthStore()
const router = useRouter()
const formRef = ref()
const loading = ref(false)
const error = ref('')
const form = reactive({ email: '', password: '' })

const rules: FormRules = {
  email: [{ required: true, message: 'Введите email', trigger: 'blur' }],
  password: [{ required: true, message: 'Введите пароль', trigger: 'blur' }],
}

async function submit() {
  error.value = ''
  try {
    await formRef.value?.validate()
  } catch {
    return
  }
  loading.value = true
  try {
    const { data } = await authApi.login(form)
    auth.setTokens(data.access_token, data.refresh_token, data.user)
    router.push('/')
  } catch (e: any) {
    error.value = e?.response?.data?.error || 'Ошибка входа'
  } finally {
    loading.value = false
  }
}
</script>

<style lang="scss" scoped>
.auth-page {
  min-height: 100vh;
  @include flex(row, center, center);
  background: $bg;
}
.auth-card {
  @include card;
  padding: $sp-8;
  width: 100%;
  max-width: 400px;
  @include mobile {
    margin: $sp-4;
  }

  &__header {
    text-align: center;
    margin-bottom: $sp-8;
    h1 {
      font-size: 22px;
      font-weight: 700;
      margin: $sp-2 0 $sp-1;
    }
    p {
      color: $text-muted;
      font-size: 14px;
    }
  }
  &__footer {
    text-align: center;
    margin-top: $sp-5;
    font-size: 13px;
    color: $text-muted;
    a {
      color: $primary;
      font-weight: 500;
    }
  }
}
</style>
