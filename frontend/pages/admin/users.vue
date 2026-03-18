<template>
  <div>
    <div class="page-header"><h1>Пользователи</h1></div>
    <div class="card">
      <n-data-table :columns="columns" :data="users" :loading="loading" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { NDataTable, NSelect, NSwitch, useMessage } from 'naive-ui'
import { h, ref, onMounted } from 'vue'
import { usersApi } from '~/api/users'
import type { User, Role } from '~/types'

definePageMeta({ middleware: ['auth', 'admin'] })

const users = ref<User[]>([])
const loading = ref(false)
const message = useMessage()
const auth = useAuthStore()

onMounted(async () => {
  loading.value = true
  try {
    const { data } = await usersApi.list()
    users.value = data
  } finally {
    loading.value = false
  }
})

const roleOptions = [
  { label: 'Студент', value: 'student' },
  { label: 'Преподаватель', value: 'teacher' },
  { label: 'Администратор', value: 'admin' },
]

async function changeRole(id: string, role: string) {
  if (id === auth.user?.id) {
    message.warning('Нельзя изменить свою роль')
    return
  }
  try {
    await usersApi.updateRole(id, role)
    message.success('Роль обновлена')
    const u = users.value.find((u) => u.id === id)
    if (u) u.role = role as Role
  } catch (e: any) {
    message.error(e?.response?.data?.error || 'Ошибка')
  }
}

async function toggleBlock(id: string, current: boolean) {
  if (id === auth.user?.id) {
    message.warning('Нельзя заблокировать себя')
    return
  }
  try {
    await usersApi.block(id, !current)
    const u = users.value.find((u) => u.id === id)
    if (u) u.is_active = !current
  } catch (e: any) {
    message.error(e?.response?.data?.error || 'Ошибка')
  }
}

const columns = [
  { title: 'Имя', key: 'full_name' },
  { title: 'Email', key: 'email' },
  {
    title: 'Роль',
    key: 'role',
    width: 200,
    render: (row: User) =>
      h(NSelect, {
        value: row.role,
        size: 'small',
        options: roleOptions,
        onUpdateValue: (v: string) => changeRole(row.id, v),
      }),
  },
  {
    title: 'Активен',
    key: 'is_active',
    width: 100,
    render: (row: User) =>
      h(NSwitch, {
        value: row.is_active,
        onUpdateValue: () => toggleBlock(row.id, row.is_active),
      }),
  },
  {
    title: 'Дата регистрации',
    key: 'created_at',
    width: 150,
    render: (row: User) => new Date(row.created_at).toLocaleDateString('ru-RU'),
  },
]
</script>
