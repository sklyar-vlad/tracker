import { ref } from 'vue'

export type ThemeType = 'light' | 'dark'

const theme = ref<ThemeType>(
  (localStorage.getItem('theme') as ThemeType) || 'dark'
)

const applyTheme = (value: ThemeType) => {
  theme.value = value
  localStorage.setItem('theme', value)
  document.documentElement.setAttribute('data-theme', value)
}

applyTheme(theme.value)

export function useTheme() {
  const toggleTheme = () => {
    applyTheme(theme.value === 'dark' ? 'light' : 'dark')
  }

  return {
    theme,
    toggleTheme,
  }
}