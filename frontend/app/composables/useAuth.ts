import { computed } from 'vue'

export function useAuth() {
  const config = useRuntimeConfig()
  
  const token = useState<string | null>('auth-token', () => null)
  const username = useState<string | null>('auth-username', () => null)
  const loading = useState<boolean>('auth-loading', () => false)
  const error = useState<string | null>('auth-error', () => null)
  const checked = useState<boolean>('auth-checked', () => false)

  const isAuthenticated = computed(() => !!token.value && !!username.value)

  async function login(usr: string, pwd: string) {
    loading.value = true
    error.value = null

    try {
      const response = await $fetch<any>(`${config.public.apiBase}/auth/login`, {
        method: 'POST',
        body: { username: usr, password: pwd },
        credentials: 'include',
      })

      if (response.code === 0) {
        username.value = response.data.username
        token.value = response.data.username
        checked.value = true
      } else {
        error.value = response.message || '登录失败'
      }
    } catch (e: any) {
      error.value = e.data?.message || '登录失败，请检查网络或后端服务'
    } finally {
      loading.value = false
    }
  }

  async function logout() {
    try {
      await $fetch(`${config.public.apiBase}/auth/logout`, {
        method: 'POST',
        credentials: 'include',
      })
    } catch (e) {
    } finally {
      token.value = null
      username.value = null
      loading.value = false
      error.value = null
      checked.value = false
      await navigateTo('/admin/login')
    }
  }

  async function checkAuth() {
    if (checked.value) {
      return
    }

    try {
      const response = await $fetch<any>(`${config.public.apiBase}/auth/me`, {
        credentials: 'include',
      })

      if (response.code === 0 && response.data?.username) {
        username.value = response.data.username
        token.value = response.data.username
      }
    } catch (e) {
      token.value = null
      username.value = null
    } finally {
      checked.value = true
    }
  }

  return {
    token,
    username,
    loading,
    error,
    checked,
    isAuthenticated,
    login,
    logout,
    checkAuth,
  }
}
