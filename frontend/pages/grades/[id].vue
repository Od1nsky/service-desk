<template>
  <div v-if="store.loading" class="center"><n-spin /></div>
  <div v-else-if="store.current" class="grade-detail">
    <!-- Header -->
    <div class="page-header">
      <div>
        <h1>{{ store.current.task_name }}</h1>
        <span class="grade-num">#{{ store.current.number }}</span>
      </div>
      <div class="header-actions">
        <!-- Teacher: assign -->
        <n-button v-if="canAssign" size="small" :loading="actionLoading" @click="doAssign">
          Взять в работу
        </n-button>
        <!-- Teacher: change status -->
        <n-select
          v-if="canChangeStatus"
          v-model:value="newStatus"
          :options="availableStatuses"
          placeholder="Сменить статус"
          style="width: 200px"
          size="small"
          @update:value="doChangeStatus"
        />
        <!-- Student: acknowledge -->
        <n-button
          v-if="canAcknowledge"
          size="small"
          type="primary"
          :loading="actionLoading"
          @click="doAcknowledge"
        >
          Подтвердить ознакомление
        </n-button>
      </div>
    </div>

    <div class="grade-layout">
      <!-- Left: notes + comments -->
      <div class="grade-main">
        <div class="card" style="margin-bottom: 16px">
          <div class="grade-notes">{{ store.current.notes }}</div>
        </div>

        <!-- Comments -->
        <div class="card">
          <h3 class="section-title">Комментарии</h3>
          <div v-if="!comments.length" class="empty-comments">Нет комментариев</div>
          <div class="comment-list">
            <div
              v-for="c in comments"
              :key="c.id"
              class="comment"
              :class="{ 'comment--internal': c.is_internal }"
            >
              <div class="comment__header">
                <strong>{{ c.author?.full_name }}</strong>
                <span class="comment__date">{{ formatDate(c.created_at) }}</span>
                <n-tag v-if="c.is_internal" size="small" type="warning">Внутренний</n-tag>
              </div>
              <div class="comment__body">{{ c.body }}</div>
            </div>
          </div>
          <!-- Comment form -->
          <div class="comment-form">
            <n-input
              v-model:value="newComment"
              type="textarea"
              :rows="3"
              placeholder="Написать комментарий..."
            />
            <div class="comment-form__footer">
              <n-checkbox v-if="auth.isTeacher" v-model:checked="isInternal"
                >Внутренняя заметка</n-checkbox
              >
              <n-button
                type="primary"
                :loading="commentLoading"
                :disabled="!newComment.trim()"
                @click="addComment"
              >
                Отправить
              </n-button>
            </div>
          </div>
        </div>
      </div>

      <!-- Right: meta -->
      <aside class="grade-aside">
        <div class="card">
          <div class="meta-row">
            <span class="meta-label">Статус оценки</span>
            <span class="status-badge" :class="`status-badge--${store.current.status}`">{{
              statusLabel(store.current.status)
            }}</span>
          </div>
          <div class="meta-row">
            <span class="meta-label">Тип оценки</span>
            <span class="priority-tag">
              <span class="dot" :style="{ background: store.current.grade_type?.color }"></span>
              {{ store.current.grade_type?.name }}
            </span>
          </div>
          <div class="meta-row">
            <span class="meta-label">Дисциплина</span>
            <span>{{ store.current.discipline?.name }}</span>
          </div>
          <div class="meta-row">
            <span class="meta-label">Студент</span>
            <span>{{ store.current.student?.full_name }}</span>
          </div>
          <div class="meta-row">
            <span class="meta-label">Преподаватель</span>
            <span>{{ store.current.teacher?.full_name || '—' }}</span>
          </div>
          <div class="meta-row">
            <span class="meta-label">Добавлено</span>
            <span>{{ formatDate(store.current.created_at) }}</span>
          </div>
          <div v-if="store.current.score != null" class="meta-row">
            <span class="meta-label">Балл</span>
            <span class="score-value">{{ store.current.score }}</span>
          </div>
        </div>

        <!-- History -->
        <div v-if="store.history.length" class="card history-card">
          <h4 class="section-title">История</h4>
          <n-timeline>
            <n-timeline-item
              v-for="h in store.history"
              :key="h.id"
              :content="`${h.old_value || '—'} → ${h.new_value}`"
              :time="formatDate(h.changed_at)"
            />
          </n-timeline>
        </div>
      </aside>
    </div>
  </div>
</template>

<script setup lang="ts">
import {
  NSpin,
  NButton,
  NSelect,
  NInput,
  NCheckbox,
  NTag,
  NTimeline,
  NTimelineItem,
  useMessage,
} from 'naive-ui'
import { commentsApi } from '~/api/comments'
import { gradesApi } from '~/api/grades'
import type { Comment, GradeStatus } from '~/types'

definePageMeta({ middleware: 'auth' })

const route = useRoute()
const store = useGradesStore()
const auth = useAuthStore()
const message = useMessage()
const id = route.params.id as string

const comments = ref<Comment[]>([])
const newComment = ref('')
const isInternal = ref(false)
const commentLoading = ref(false)
const actionLoading = ref(false)
const newStatus = ref<string | null>(null)

onMounted(async () => {
  await store.fetchById(id)
  const { data } = await commentsApi.list(id)
  comments.value = data
})

const canAssign = computed(
  () => auth.isTeacher && !store.current?.teacher_id && store.current?.status === 'pending'
)
const canChangeStatus = computed(
  () => auth.isTeacher && store.current?.status !== 'closed'
)
const canAcknowledge = computed(
  () =>
    auth.isStudent &&
    store.current?.student_id === auth.user?.id &&
    store.current?.status === 'graded'
)

const transitions: Record<string, GradeStatus[]> = {
  pending: ['graded'],
  graded: ['certified', 'not_certified', 'pending'],
  not_certified: ['graded'],
  certified: ['closed'],
  closed: [],
}

const statusOptions: Record<string, string> = {
  pending: 'Ожидает',
  graded: 'Оценено',
  not_certified: 'Не аттестован',
  certified: 'Аттестован',
  closed: 'Завершено',
}

const availableStatuses = computed(() => {
  const cur = store.current?.status || ''
  return (transitions[cur] || []).map((s) => ({ label: statusOptions[s], value: s }))
})

async function doAssign() {
  actionLoading.value = true
  try {
    await gradesApi.assign(id)
    await store.fetchById(id)
    message.success('Оценка назначена на вас')
  } catch (e: any) {
    message.error(e?.response?.data?.error || 'Ошибка')
  } finally {
    actionLoading.value = false
  }
}

async function doChangeStatus(status: string) {
  try {
    await gradesApi.updateStatus(id, status)
    await store.fetchById(id)
    newStatus.value = null
    message.success('Статус обновлён')
  } catch (e: any) {
    message.error(e?.response?.data?.error || 'Ошибка')
  }
}

async function doAcknowledge() {
  actionLoading.value = true
  try {
    await gradesApi.acknowledge(id)
    await store.fetchById(id)
    message.success('Ознакомление подтверждено')
  } catch (e: any) {
    message.error(e?.response?.data?.error || 'Ошибка')
  } finally {
    actionLoading.value = false
  }
}

async function addComment() {
  if (!newComment.value.trim()) return
  commentLoading.value = true
  try {
    const { data } = await commentsApi.create(id, {
      body: newComment.value,
      is_internal: isInternal.value,
    })
    comments.value.push(data)
    newComment.value = ''
    isInternal.value = false
  } catch (e: any) {
    message.error(e?.response?.data?.error || 'Ошибка')
  } finally {
    commentLoading.value = false
  }
}

const statusLabel = (s: string) => statusOptions[s] || s
const formatDate = (d: string) =>
  new Date(d).toLocaleString('ru-RU', {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  })
</script>

<style lang="scss" scoped>
.grade-layout {
  display: grid;
  grid-template-columns: 1fr 300px;
  gap: $sp-4;
  align-items: flex-start;
  @include mobile {
    grid-template-columns: 1fr;
  }
}
.grade-main {
}
.grade-aside {
}
.grade-num {
  font-size: 13px;
  color: $text-muted;
}
.grade-notes {
  line-height: 1.7;
  white-space: pre-wrap;
}
.meta-row {
  @include flex(row, center, space-between, $sp-2);
  padding: $sp-2 0;
  border-bottom: 1px solid $bg;
  font-size: 13px;
  &:last-child {
    border-bottom: none;
  }
}
.meta-label {
  color: $text-muted;
}
.score-value {
  font-weight: 700;
  font-size: 15px;
  color: $primary;
}
.section-title {
  font-size: 15px;
  font-weight: 600;
  margin-bottom: $sp-4;
}
.comment-list {
  display: flex;
  flex-direction: column;
  gap: $sp-3;
  margin-bottom: $sp-5;
}
.comment {
  padding: $sp-3 $sp-4;
  border-radius: $radius;
  background: $bg;
  &--internal {
    background: rgba($warning, 0.07);
    border-left: 3px solid $warning;
  }
  &__header {
    @include flex(row, center, flex-start, $sp-3);
    margin-bottom: $sp-2;
  }
  &__date {
    color: $text-muted;
    font-size: 12px;
  }
  &__body {
    font-size: 13px;
    line-height: 1.6;
    white-space: pre-wrap;
  }
}
.empty-comments {
  color: $text-muted;
  font-size: 13px;
  margin-bottom: $sp-4;
}
.comment-form {
  border-top: 1px solid $border;
  padding-top: $sp-4;
}
.comment-form__footer {
  @include flex(row, center, space-between);
  margin-top: $sp-3;
}
.header-actions {
  @include flex(row, center, flex-start, $sp-3);
  flex-wrap: wrap;
}
.history-card {
  margin-top: $sp-4;
}
.center {
  @include flex(row, center, center);
  padding: 60px;
}
</style>
