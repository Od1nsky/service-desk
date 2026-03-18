<template>
  <div style="max-width: 680px">
    <div class="page-header"><h1>Выставить оценку</h1></div>
    <div class="card">
      <n-form ref="formRef" :model="form" :rules="rules" label-placement="top">
        <n-form-item label="Студент" path="student_id">
          <n-select
            v-model:value="form.student_id"
            :options="studentOptions"
            placeholder="Выберите студента"
            filterable
          />
        </n-form-item>
        <n-form-item label="Название работы" path="task_name">
          <n-input v-model:value="form.task_name" placeholder="Название задания/работы" />
        </n-form-item>
        <n-form-item label="Примечания" path="notes">
          <n-input
            v-model:value="form.notes"
            type="textarea"
            :rows="5"
            placeholder="Дополнительные примечания..."
          />
        </n-form-item>
        <div class="form-row">
          <n-form-item label="Дисциплина" path="discipline_id" style="flex: 1">
            <n-select
              v-model:value="form.discipline_id"
              :options="disciplineOptions"
              placeholder="Выберите дисциплину"
            />
          </n-form-item>
          <n-form-item label="Тип оценки" path="grade_type_id" style="flex: 1">
            <n-select
              v-model:value="form.grade_type_id"
              :options="gradeTypeOptions"
              placeholder="Выберите тип оценки"
            />
          </n-form-item>
        </div>
        <div class="form-actions">
          <n-button @click="$router.back()">Отмена</n-button>
          <n-button type="primary" :loading="loading" @click="submit">Выставить оценку</n-button>
        </div>
      </n-form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { NForm, NFormItem, NInput, NSelect, NButton, useMessage } from 'naive-ui'
import type { FormRules } from 'naive-ui'
import { disciplinesApi } from '~/api/disciplines'
import { gradeTypesApi } from '~/api/grade-types'
import { gradesApi } from '~/api/grades'
import { usersApi } from '~/api/users'
import type { Discipline, GradeType, User } from '~/types'

definePageMeta({ middleware: ['teacher'] })

const router = useRouter()
const message = useMessage()
const formRef = ref()
const loading = ref(false)
const form = reactive({
  student_id: null as string | null,
  task_name: '',
  notes: '',
  discipline_id: null as number | null,
  grade_type_id: null as number | null,
})

const rules: FormRules = {
  student_id: [{ required: true, message: 'Выберите студента', trigger: 'change' }],
  task_name: [{ required: true, message: 'Введите название работы', trigger: 'blur' }],
  discipline_id: [{ required: true, type: 'number' as const, message: 'Выберите дисциплину', trigger: 'change' }],
  grade_type_id: [{ required: true, type: 'number' as const, message: 'Выберите тип оценки', trigger: 'change' }],
}

const disciplines = ref<Discipline[]>([])
const gradeTypes = ref<GradeType[]>([])
const students = ref<User[]>([])

onMounted(async () => {
  const [d, g, u] = await Promise.all([
    disciplinesApi.list(),
    gradeTypesApi.list(),
    usersApi.list(),
  ])
  disciplines.value = d.data
  gradeTypes.value = g.data
  students.value = (u.data as User[]).filter((u) => u.role === 'student')
})

const disciplineOptions = computed(() =>
  disciplines.value.map((d) => ({ label: d.name, value: d.id }))
)
const gradeTypeOptions = computed(() =>
  gradeTypes.value.map((g) => ({
    label: g.name,
    value: g.id,
    style: `color: ${g.color}`,
  }))
)
const studentOptions = computed(() =>
  students.value.map((s) => ({ label: `${s.full_name} (${s.email})`, value: s.id }))
)

async function submit() {
  try {
    await formRef.value?.validate()
  } catch {
    return
  }
  loading.value = true
  try {
    await gradesApi.create({
      task_name: form.task_name,
      notes: form.notes,
      discipline_id: form.discipline_id!,
      grade_type_id: form.grade_type_id!,
      student_id: form.student_id!,
    })
    message.success('Оценка добавлена')
    router.push('/teacher/queue')
  } catch (e: any) {
    message.error(e?.response?.data?.error || 'Ошибка добавления оценки')
  } finally {
    loading.value = false
  }
}
</script>

<style lang="scss" scoped>
.form-row {
  @include flex(row, flex-start, flex-start, $sp-4);
  @include mobile {
    flex-direction: column;
  }
}
.form-actions {
  @include flex(row, center, flex-end, $sp-3);
  margin-top: $sp-4;
}
</style>
