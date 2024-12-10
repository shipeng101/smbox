import { ref, watchEffect } from 'vue'
import { darkTheme, useOsTheme } from 'naive-ui'
import type { GlobalTheme } from 'naive-ui'

export function useTheme() {
  const isDarkTheme = ref(false)
  const osTheme = useOsTheme()
  
  const theme = ref<GlobalTheme | null>(null)
  
  const toggleTheme = () => {
    isDarkTheme.value = !isDarkTheme.value
  }
  
  // 监听主题变化
  watchEffect(() => {
    // 根据开关状态或系统主题设置
    if (isDarkTheme.value || osTheme.value === 'dark') {
      theme.value = darkTheme
    } else {
      theme.value = null
    }
  })
  
  return {
    isDarkTheme,
    theme,
    toggleTheme
  }
} 