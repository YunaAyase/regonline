import type { RegistrationRecord, ApiResponse } from '~/types'

export function useRegistrations() {
  const config = useRuntimeConfig()
  const registrations = ref<RegistrationRecord[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)

  async function fetchRegistrations(name?: string, classId?: number) {
    loading.value = true
    error.value = null

    try {
      const params = new URLSearchParams()
      if (name) params.set('name', name)
      if (classId && classId > 0) params.set('class_id', String(classId))

      const qs = params.toString()
      const url = qs ? `${config.public.apiBase}/registrations?${qs}` : `${config.public.apiBase}/registrations`

      const res = await $fetch<ApiResponse<RegistrationRecord[]>>(url, { credentials: 'include' })

      if (res.code === 0) {
        registrations.value = res.data || []
      } else {
        error.value = res.message || '获取报名列表失败'
      }
    } catch (e: any) {
      error.value = e.data?.message || '获取报名列表失败，请检查网络'
    } finally {
      loading.value = false
    }
  }

  async function deleteRegistration(id: number): Promise<boolean> {
    try {
      const url = `${config.public.apiBase}/registrations/${String(id)}`
      const res = await $fetch<ApiResponse>(url, {
        method: 'DELETE',
        credentials: 'include',
      })
      return res.code === 0
    } catch (e) {
      console.error('Delete failed', e)
    }
    return false
  }

  async function deleteAndRefresh(id: number, onSuccess?: () => void): Promise<boolean> {
    const ok = await deleteRegistration(id)
    if (ok) {
      registrations.value = registrations.value.filter(r => r.id !== id)
      onSuccess?.()
    }
    return ok
  }

  async function batchDelete(ids: number[], onSuccess?: () => void): Promise<number> {
    let successCount = 0
    for (const id of ids) {
      if (await deleteRegistration(id)) successCount++
    }
    if (successCount > 0) {
      const deletedSet = new Set(ids.slice(0, successCount))
      registrations.value = registrations.value.filter(r => !deletedSet.has(r.id))
      onSuccess?.()
    }
    return successCount
  }

  function getPhotoUrl(photoPath: string | null): string | null {
    if (!photoPath) return null
    const filename = photoPath.replace(/\\/g, '/').split('/').pop()
    if (!filename) return null
    return `${config.public.apiBase}/photos/${filename}`
  }

  return {
    registrations,
    loading,
    error,
    fetchRegistrations,
    deleteRegistration,
    deleteAndRefresh,
    batchDelete,
    getPhotoUrl,
  }
}