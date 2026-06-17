<template>
  <header class="header">
    <div class="header-container">
      <div class="header-left">
        <RouterLink to="/" class="logo">
          <span class="logo-text">SelfDev</span>
        </RouterLink>
      </div>

      <div class="header-right">
        <button
          class="theme-toggle"
          @click="themeStore.toggleTheme"
          :title="`Switch to ${theme === 'light' ? 'dark' : 'light'} mode`"
        >
          <span v-if="theme === 'light'">🌙</span>
          <span v-else>☀️</span>
        </button>
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
import { useTheme } from '@/composables/useTheme'

const themeStore = useTheme()
const theme = themeStore.theme
</script>

<style scoped>
.header {
  background: var(--bg-secondary);
  border-bottom: 1px solid var(--border-color);
  padding: 12px 0;
}

.header-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 16px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.header-left {
  display: flex;
  align-items: center;
}

.logo {
  display: flex;
  gap: 8px;
  text-decoration: none;
  font-weight: 700;
  color: var(--text-primary);
}

.logo-icon {
  font-size: 20px;
}

.logo-text {
  font-size: 18px;
}

.header-nav {
  display: flex;
  gap: 20px;
}

.nav-link {
  text-decoration: none;
  color: var(--text-secondary);
  font-weight: 500;
}

.nav-link.router-link-active {
  color: var(--accent-primary);
}

.header-right {
  display: flex;
  align-items: center;
}

.theme-toggle {
  background: var(--bg-tertiary);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  width: 40px;
  height: 40px;
  cursor: pointer;
}
</style>
