export default defineNuxtRouteMiddleware(async (to) => {
  if (to.path.startsWith('/admin') && to.path !== '/admin/login') {
    const { isAuthenticated, checkAuth } = useAuth()

    if (!isAuthenticated.value) {
      await checkAuth()
    }

    if (!isAuthenticated.value) {
      return navigateTo({
        path: '/admin/login',
        query: { redirect: to.fullPath },
      })
    }
  }
})
