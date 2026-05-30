import type { SiteSettings } from '~/types'

export function useSiteSettings() {
  const config = useRuntimeConfig()

  const siteSettings = useState<SiteSettings | null>('site-settings', () => null)
  const loading = useState<boolean>('site-settings-loading', () => false)

  const siteName = computed(() => siteSettings.value?.site_name || '在线报名系统')

  async function fetchSiteSettings() {
    if (siteSettings.value) return

    loading.value = true
    try {
      const res = await $fetch<any>(`${config.public.apiBase}/settings`, {
        credentials: 'include',
      })
      if (res.code === 0 && res.data) {
        siteSettings.value = res.data
      }
    } catch (e) {
      console.error('Failed to fetch site settings', e)
    } finally {
      loading.value = false
    }
  }

  return {
    siteSettings,
    siteName,
    loading,
    fetchSiteSettings,
  }
}