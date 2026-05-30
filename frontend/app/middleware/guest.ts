export default defineNuxtRouteMiddleware(async () => {
  const { isAuthenticated, checkAuth } = useAuth()

  await checkAuth()

  if (isAuthenticated.value) {
    return navigateTo('/admin')
  }
})
