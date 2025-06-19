export default defineNuxtRouteMiddleware((to) => {
  const isAuthenticated = () => process.client ? localStorage.getItem('auth_token') : null;
  
  // Eğer kullanıcı giriş yapmışsa ve login veya ana sayfaya gitmeye çalışıyorsa
  if (isAuthenticated() && (to.path === '/login' || to.path === '/')) {
    return navigateTo('/home');
  }
});
