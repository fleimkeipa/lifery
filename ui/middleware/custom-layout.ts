export default defineNuxtRouteMiddleware((to) => {
    const isAuthenticated = () => process.client ? localStorage.getItem('auth_token') : null;

    if (isAuthenticated()) {
        setPageLayout('default')
    } else {
        setPageLayout('not-authenticated')
    }
  })
  