import { ref, watch } from 'vue'

export type ThemeType = 'light' | 'dark'

const theme = ref<ThemeType>('light')

export function useTheme() {
  const applyTheme = (value: ThemeType) => {
    document.documentElement.setAttribute('data-theme', value)
    localStorage.setItem('theme', value)
  }

  const initTheme = () => {
    const stored = localStorage.getItem('theme') as ThemeType | null

    if (stored === 'light' || stored === 'dark') {
      theme.value = stored
    } else {
      theme.value = 'light'
    }

    applyTheme(theme.value)
  }

  const toggleTheme = () => {
    theme.value = theme.value === 'light' ? 'dark' : 'light'
  }

  const setTheme = (value: ThemeType) => {
    theme.value = value
  }

  watch(theme, (newValue) => {
    applyTheme(newValue)
  })

  return {
    theme,
    initTheme,
    toggleTheme,
    setTheme,
  }
}
