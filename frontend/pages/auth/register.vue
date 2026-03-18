<template>
  <div class="auth-page">
    <div class="auth-card">
      <div class="auth-card__header">
        <n-icon size="32" color="#0d6efd"><TicketOutline /></n-icon>
        <h1>АИС Успеваемость</h1>
        <p>Создайте аккаунт</p>
      </div>
      <n-form ref="formRef" :model="form" :rules="rules" @submit.prevent="submit">
        <n-form-item path="full_name" label="Полное имя">
          <n-input v-model:value="form.full_name" placeholder="Иван Иванов" size="large" />
        </n-form-item>
        <n-form-item path="email" label="Email">
          <n-input v-model:value="form.email" placeholder="user@university.ru" size="large" />
        </n-form-item>
        <n-form-item path="password" label="Пароль">
          <n-input
            v-model:value="form.password"
            type="password"
            show-password-on="click"
            placeholder="Минимум 6 символов"
            size="large"
          />
        </n-form-item>
        <n-alert v-if="error" type="error" :title="error" style="margin-bottom: 16px" />
        <n-alert v-if="successMsg" type="success" :title="successMsg" style="margin-bottom: 16px" />
        <n-button type="primary" size="large" block :loading="loading" @click="submit">
          Зарегистрироваться
        </n-button>
      </n-form>
      <div class="auth-card__footer">
        Уже есть аккаунт? <NuxtLink to="/auth/login">Войти</NuxtLink>
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

const router = useRouter()
const formRef = ref()
const loading = ref(false)
const error = ref('')
const successMsg = ref('')
const form = reactive({ full_name: '', email: '', password: '' })

const rules: FormRules = {
  full_name: [{ required: true, message: 'Введите полное имя', trigger: 'blur' }],
  email: [
    { required: true, message: 'Введите email', trigger: 'blur' },
    { type: 'email' as const, message: 'Некорректный email', trigger: 'blur' },
  ],
  password: [
    { required: true, message: 'Введите пароль', trigger: 'blur' },
    { min: 6, message: 'Пароль должен быть не менее 6 символов', trigger: 'blur' },
  ],
}

async function submit() {
  error.value = ''
  successMsg.value = ''
  try {
    await formRef.value?.validate()
  } catch {
    return
  }
  loading.value = true
  try {
    await authApi.register(form)
    successMsg.value = 'Аккаунт создан! Перенаправляем на вход...'
    setTimeout(() => router.push('/auth/login'), 1500)
  } catch (e: any) {
    error.value = e?.response?.data?.error || 'Ошибка регистрации'
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
