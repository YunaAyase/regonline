import type { ClassItem, ApiResponse } from '~/types'

const classes = ref<ClassItem[]>([])
const loading = ref(false)
const error = ref<string | null>(null)

export function useClasses() {
  const config = useRuntimeConfig()

  async function fetchClasses() {
    loading.value = true
    error.value = null

    try {
      const res = await $fetch<ApiResponse<ClassItem[]>>(`${config.public.apiBase}/classes`, { credentials: 'include' })

      if (res.code === 0) {
        classes.value = res.data || []
      } else {
        error.value = res.message || '获取班级列表失败'
      }
    } catch (e: any) {
      error.value = e.data?.message || '获取班级列表失败，请检查网络'
    } finally {
      loading.value = false
    }
  }

  return {
    classes,
    loading,
    error,
    fetchClasses,
  }
}
