<template>
  <div>
    <div class="page-header">
      <h1>Дисциплины</h1>
      <n-button v-if="auth.isAdmin" type="primary" @click="openModal()">
        <template #icon><n-icon><AddOutline /></n-icon></template>
        Добавить дисциплину
      </n-button>
    </div>

    <div v-if="loading" class="center"><n-spin /></div>
    <div v-else-if="!disciplines.length" class="empty">
      <n-empty description="Дисциплины не найдены" />
    </div>
    <div v-else class="disc-grid">
      <div
        v-for="d in disciplines"
        :key="d.id"
        class="disc-card card"
        @click="$router.push(`/disciplines/${d.id}`)"
      >
        <div class="disc-card__icon">
          <n-icon size="28" color="#0d6efd"><BookOutline /></n-icon>
        </div>
        <div class="disc-card__body">
          <div class="disc-card__name">{{ d.name }}</div>
          <div class="disc-card__desc">{{ d.description || 'Без описания' }}</div>
        </div>
        <div class="disc-card__actions" @click.stop>
          <template v-if="auth.isAdmin">
            <n-button size="small" @click="openModal(d)">Изменить</n-button>
            <n-button size="small" type="error" ghost @click="deleteItem(d)">Удалить</n-button>
          </template>
          <n-icon v-else size="18" color="#adb5bd"><ChevronForwardOutline /></n-icon>
        </div>
      </div>
    </div>

    <!-- Add / Edit modal -->
    <n-modal
      v-model:show="showModal"
      :title="editItem ? 'Редактировать дисциплину' : 'Новая дисциплина'"
      preset="card"
      style="width: 440px"
    >
      <n-form ref="formRef" :model="form" label-placement="top">
        <n-form-item
          label="Название"
          path="name"
          :rule="{ required: true, message: 'Введите название' }"
        >
          <n-input v-model:value="form.name" placeholder="Например: Математический анализ" />
        </n-form-item>
        <n-form-item label="Описание" path="description">
          <n-input v-model:value="form.description" type="textarea" :rows="2" placeholder="Краткое описание курса" />
        </n-form-item>
      </n-form>
      <template #footer>
        <div style="display:flex;justify-content:flex-end;gap:8px">
          <n-button @click="showModal = false">Отмена</n-button>
          <n-button type="primary" :loading="saving" @click="save">Сохранить</n-button>
        </div>
      </template>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { NSpin, NEmpty, NIcon, NButton, NModal, NForm, NFormItem, NInput, useMessage, useDialog } from 'naive-ui'
import { BookOutline, ChevronForwardOutline, AddOutline } from '@vicons/ionicons5'
import { disciplinesApi } from '~/api/disciplines'
import type { Discipline } from '~/types'

definePageMeta({ middleware: ['auth', 'teacher'] })

const auth = useAuthStore()
const message = useMessage()
const dialog = useDialog()

const disciplines = ref<Discipline[]>([])
const loading = ref(true)
const showModal = ref(false)
const saving = ref(false)
const editItem = ref<Discipline | null>(null)
const formRef = ref()
const form = reactive({ name: '', description: '' })

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
    content: `«${item.name}» будет удалена безвозвратно`,
    positiveText: 'Удалить',
    negativeText: 'Отмена',
    onPositiveClick: async () => {
      try {
        await disciplinesApi.delete(item.id)
        message.success('Удалено')
        fetchDisciplines()
      } catch (e: any) {
        message.error(e?.response?.data?.error || 'Нельзя удалить: есть связанные оценки')
      }
    },
  })
}
</script>

<style lang="scss" scoped>
.disc-grid {
  display: flex;
  flex-direction: column;
  gap: $sp-3;
}
.disc-card {
  @include flex(row, center, flex-start, $sp-4);
  cursor: pointer;
  transition: box-shadow 0.15s, border-color 0.15s;
  &:hover {
    border-color: $primary;
    box-shadow: 0 2px 12px rgba($primary, 0.12);
  }
  &__icon {
    flex-shrink: 0;
    width: 48px;
    height: 48px;
    border-radius: $radius;
    background: rgba($primary, 0.06);
    @include flex(row, center, center);
  }
  &__body {
    flex: 1;
    min-width: 0;
  }
  &__name {
    font-size: 15px;
    font-weight: 600;
    color: $text;
    margin-bottom: 2px;
  }
  &__desc {
    font-size: 13px;
    color: $text-muted;
    @include truncate;
  }
  &__actions {
    @include flex(row, center, flex-end, $sp-2);
    flex-shrink: 0;
  }
}
.empty {
  @include flex(column, center, center);
  padding: 60px;
}
.center {
  @include flex(row, center, center);
  padding: 60px;
}
</style>
