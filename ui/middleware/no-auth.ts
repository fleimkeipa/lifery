export default defineNuxtRouteMiddleware((to, from) => {
  const isAuthenticated = () => process.client ? localStorage.getItem('auth_token') : null;
  if (!!isAuthenticated()) {
    return navigateTo("/");
  }
});
