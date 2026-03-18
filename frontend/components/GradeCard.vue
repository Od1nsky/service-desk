<template>
  <div class="grade-card">
    <div class="grade-card__header">
      <span class="grade-card__num">#{{ grade.number }}</span>
      <span class="status-badge" :class="`status-badge--${grade.status}`">{{ statusLabel }}</span>
      <span class="priority-tag">
        <span class="dot" :style="{ background: grade.grade_type?.color }"></span>
        {{ grade.grade_type?.name }}
      </span>
    </div>
    <div class="grade-card__title">{{ grade.task_name }}</div>
    <div class="grade-card__meta">
      <span>{{ grade.discipline?.name }}</span>
      <span>·</span>
      <span>{{ formatDate(grade.created_at) }}</span>
      <span v-if="grade.student">· {{ grade.student.full_name }}</span>
      <span v-if="grade.score != null">· Балл: {{ grade.score }}</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Grade } from '~/types'

const { grade } = defineProps<{ grade: Grade }>()

const statusLabels: Record<string, string> = {
  pending: 'Ожидает',
  graded: 'Оценено',
  not_certified: 'Не аттестован',
  certified: 'Аттестован',
  closed: 'Завершено',
}
const statusLabel = computed(() => statusLabels[grade.status] || grade.status)
const formatDate = (d: string) => new Date(d).toLocaleDateString('ru-RU')
</script>

<style lang="scss" scoped>
.grade-card {
  @include card;
  padding: $sp-4 $sp-5;
  cursor: pointer;
  transition: box-shadow 0.15s, transform 0.1s;
  &:hover {
    box-shadow: $shadow;
    transform: translateY(-1px);
  }

  &__header {
    @include flex(row, center, flex-start, $sp-3);
    margin-bottom: $sp-2;
  }
  &__num {
    font-size: 12px;
    color: $text-muted;
    font-weight: 600;
  }
  &__title {
    font-size: 15px;
    font-weight: 500;
    margin-bottom: $sp-2;
    @include truncate;
  }
  &__meta {
    font-size: 12px;
    color: $text-muted;
    @include flex(row, center, flex-start, $sp-2);
  }
}
</style>
