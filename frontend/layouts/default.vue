<template>
  <div class="page-wrapper">
    <!-- Mobile overlay -->
    <div v-if="sidebarOpen" class="overlay" @click="sidebarOpen = false" />

    <!-- Sidebar -->
    <aside class="sidebar" :class="{ 'sidebar--open': sidebarOpen }">
      <div class="sidebar__logo">
        <n-icon size="24" color="#0d6efd"><TicketOutline /></n-icon>
        <span>АИС Успеваемость</span>
      </div>
      <nav class="sidebar__nav">
        <NuxtLink
          v-for="item in navItems"
          :key="item.to"
          :to="item.to"
          class="sidebar__item"
          active-class="sidebar__item--active"
          @click="sidebarOpen = false"
        >
          <n-icon size="18"><component :is="item.icon" /></n-icon>
          <span>{{ item.label }}</span>
        </NuxtLink>
      </nav>
      <div class="sidebar__user">
        <div class="sidebar__user-info">
          <div class="sidebar__user-name">{{ auth.user?.full_name }}</div>
          <div class="sidebar__user-role">{{ roleLabel }}</div>
        </div>
        <n-button text @click="auth.logout()">
          <n-icon size="18"><LogOutOutline /></n-icon>
        </n-button>
      </div>
    </aside>

    <!-- Main -->
    <div class="main-content">
      <!-- Header -->
      <header class="app-header">
        <button class="burger" data-testid="burger-btn" @click="sidebarOpen = !sidebarOpen">
          <n-icon size="22"><MenuOutline /></n-icon>
        </button>
        <div class="app-header__title">{{ pageTitle }}</div>
      </header>

      <div class="content-area">
        <slot />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { NButton, NIcon } from 'naive-ui'
import {
  TicketOutline,
  HomeOutline,
  ListOutline,
  AddOutline,
  PeopleOutline,
  FolderOutline,
  BarChartOutline,
  MenuOutline,
  LogOutOutline,
} from '@vicons/ionicons5'

const auth = useAuthStore()
const route = useRoute()
const sidebarOpen = ref(false)

const allNavItems = [
  { to: '/', label: 'Дашборд', icon: HomeOutline, roles: ['admin', 'teacher'] },
  { to: '/grades', label: 'Мои оценки', icon: ListOutline, roles: ['student'] },
  { to: '/grades/create', label: 'Выставить оценку', icon: AddOutline, roles: ['teacher'] },
  { to: '/teacher/queue', label: 'Очередь на оценивание', icon: TicketOutline, roles: ['teacher', 'admin'] },
  { to: '/admin/users', label: 'Пользователи', icon: PeopleOutline, roles: ['admin'] },
  { to: '/admin/disciplines', label: 'Дисциплины', icon: FolderOutline, roles: ['admin'] },
  { to: '/', label: 'Статистика', icon: BarChartOutline, roles: ['admin', 'teacher'] },
]

const navItems = computed(() =>
  allNavItems.filter((i) => i.roles.includes(auth.user?.role || ''))
)

const roleLabel = computed(
  () =>
    ({
      student: 'Студент',
      teacher: 'Преподаватель',
      admin: 'Администратор',
    }[auth.user?.role || ''] || '')
)

const pageTitle = computed(() => {
  const map: Record<string, string> = {
    '/': 'Дашборд',
    '/grades': 'Мои оценки',
    '/grades/create': 'Выставить оценку',
    '/teacher/queue': 'Очередь на оценивание',
    '/admin/users': 'Пользователи',
    '/admin/disciplines': 'Дисциплины',
  }
  return map[route.path] || 'АИС Успеваемость'
})
</script>

<style lang="scss" scoped>
.sidebar {
  position: fixed;
  top: 0;
  left: 0;
  bottom: 0;
  width: $sidebar-width;
  background: $surface;
  border-right: 1px solid $border;
  display: flex;
  flex-direction: column;
  z-index: 100;
  transition: transform 0.25s ease;

  @include mobile {
    transform: translateX(-100%);
    &--open {
      transform: translateX(0);
      box-shadow: $shadow-lg;
    }
  }

  &__logo {
    @include flex(row, center, flex-start, $sp-3);
    padding: $sp-5 $sp-6;
    font-size: 17px;
    font-weight: 700;
    color: $text;
    border-bottom: 1px solid $border;
  }

  &__nav {
    flex: 1;
    overflow-y: auto;
    padding: $sp-3;
  }

  &__item {
    @include flex(row, center, flex-start, $sp-3);
    padding: $sp-3 $sp-4;
    border-radius: $radius;
    color: $text-muted;
    font-size: 14px;
    transition: all 0.15s;
    margin-bottom: 2px;
    &:hover {
      background: $bg;
      color: $text;
    }
    &--active {
      background: rgba($primary, 0.08);
      color: $primary;
      font-weight: 500;
    }
  }

  &__user {
    @include flex(row, center, space-between, $sp-3);
    padding: $sp-4 $sp-5;
    border-top: 1px solid $border;
  }
  &__user-info {
    flex: 1;
    min-width: 0;
  }
  &__user-name {
    font-size: 13px;
    font-weight: 600;
    @include truncate;
  }
  &__user-role {
    font-size: 11px;
    color: $text-muted;
  }
}

.overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.4);
  z-index: 99;
  @include desktop {
    display: none;
  }
}

.app-header {
  @include flex(row, center, flex-start, $sp-3);
  height: $header-height;
  padding: 0 $sp-6;
  background: $surface;
  border-bottom: 1px solid $border;
  flex-shrink: 0;
  &__title {
    font-size: 16px;
    font-weight: 600;
  }
  .burger {
    display: flex;
    align-items: center;
    justify-content: center;
    background: none;
    border: none;
    cursor: pointer;
    padding: 4px;
    color: $text;
    border-radius: $radius-sm;
    &:hover { background: $bg; }
    @include desktop {
      display: none;
    }
  }
}
</style>
