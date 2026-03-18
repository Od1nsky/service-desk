<template>
  <NConfigProvider :theme-overrides="themeOverrides">
    <NMessageProvider>
      <NDialogProvider>
        <NuxtPage />
      </NDialogProvider>
    </NMessageProvider>
  </NConfigProvider>
</template>

<script setup lang="ts">
import { NConfigProvider, NMessageProvider, NDialogProvider } from 'naive-ui'
import type { GlobalThemeOverrides } from 'naive-ui'

const themeOverrides: GlobalThemeOverrides = {
  common: {
    primaryColor: '#0d6efd',
    primaryColorHover: '#3d8bfd',
    primaryColorPressed: '#0a58ca',
    borderRadius: '8px',
  },
}

const auth = useAuthStore()
onMounted(() => auth.init())
</script>
