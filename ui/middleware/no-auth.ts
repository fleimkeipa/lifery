export default defineNuxtRouteMiddleware((to) => {
  const isAuthenticated = () => process.client ? localStorage.getItem('auth_token') : null;
  
  // Eğer kullanıcı giriş yapmışsa ve login sayfasına gitmeye çalışıyorsa
  if (isAuthenticated() && to.path === '/login') {
    return navigateTo('/');
  }
});
