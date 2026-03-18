<template>
  <div>
    <div class="page-header">
      <h1>Дисциплины</h1>
      <n-button type="primary" @click="openModal()">Добавить дисциплину</n-button>
    </div>
    <div class="card">
      <n-data-table :columns="columns" :data="disciplines" :loading="loading" />
    </div>

    <n-modal
      v-model:show="showModal"
      :title="editItem ? 'Редактировать' : 'Новая дисциплина'"
      preset="card"
      style="width: 440px"
    >
      <n-form ref="formRef" :model="form" label-placement="top">
        <n-form-item
          label="Название"
          path="name"
          :rule="{ required: true, message: 'Введите название' }"
        >
          <n-input v-model:value="form.name" />
        </n-form-item>
        <n-form-item label="Описание" path="description">
          <n-input v-model:value="form.description" type="textarea" :rows="2" />
        </n-form-item>
      </n-form>
      <template #footer>
        <div style="display: flex; justify-content: flex-end; gap: 8px">
          <n-button @click="showModal = false">Отмена</n-button>
          <n-button type="primary" :loading="saving" @click="save">Сохранить</n-button>
        </div>
      </template>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { NDataTable, NButton, NModal, NForm, NFormItem, NInput, useDialog, useMessage } from 'naive-ui'
import { h } from 'vue'
import { disciplinesApi } from '~/api/disciplines'
import type { Discipline } from '~/types'

definePageMeta({ middleware: ['auth', 'admin'] })

const disciplines = ref<Discipline[]>([])
const loading = ref(false)
const showModal = ref(false)
const saving = ref(false)
const editItem = ref<Discipline | null>(null)
const formRef = ref()
const form = reactive({ name: '', description: '' })
const message = useMessage()
const dialog = useDialog()

onMounted(fetchDisciplines)

async function fetchDisciplines() {
  loading.value = true
  try {
    const { data } = await disciplinesApi.list()
    disciplines.value = data
  } finally {
    loading.value = false
  }
}

function openModal(item?: Discipline) {
  editItem.value = item || null
  form.name = item?.name || ''
  form.description = item?.description || ''
  showModal.value = true
}

async function save() {
  try {
    await formRef.value?.validate()
  } catch {
    return
  }
  saving.value = true
  try {
    if (editItem.value) {
      await disciplinesApi.update(editItem.value.id, form)
      message.success('Обновлено')
    } else {
      await disciplinesApi.create(form)
      message.success('Создано')
    }
    showModal.value = false
    fetchDisciplines()
  } catch (e: any) {
    message.error(e?.response?.data?.error || 'Ошибка')
  } finally {
    saving.value = false
  }
}

function deleteItem(item: Discipline) {
  dialog.warning({
    title: 'Удалить дисциплину?',
    content: `"${item.name}" будет удалена`,
    positiveText: 'Удалить',
    negativeText: 'Отмена',
    onPositiveClick: async () => {
      try {
        await disciplinesApi.delete(item.id)
        message.success('Удалено')
        fetchDisciplines()
      } catch (e: any) {
        message.error(e?.response?.data?.error || 'Нельзя удалить')
      }
    },
  })
}

const columns = [
  { title: 'Название', key: 'name' },
  { title: 'Описание', key: 'description', ellipsis: { tooltip: true } },
  {
    title: 'Дата',
    key: 'created_at',
    width: 140,
    render: (row: Discipline) => new Date(row.created_at).toLocaleDateString('ru-RU'),
  },
  {
    title: '',
    key: 'actions',
    width: 140,
    render: (row: Discipline) =>
      h('div', { style: 'display:flex;gap:8px' }, [
        h(NButton, { size: 'small', onClick: () => openModal(row) }, { default: () => 'Изменить' }),
        h(
          NButton,
          { size: 'small', type: 'error', ghost: true, onClick: () => deleteItem(row) },
          { default: () => 'Удалить' }
        ),
      ]),
  },
]
</script>
