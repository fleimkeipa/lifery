export default defineNuxtRouteMiddleware(() => {
  if (process.client && !localStorage.getItem('auth_token')) {
    return navigateTo('/intro');
  }
}); 