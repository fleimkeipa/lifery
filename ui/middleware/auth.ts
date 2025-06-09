import jwtDecode from "jwt-decode";

export default defineNuxtRouteMiddleware((to, from) => {
  const isAuthenticated = () => {
    if (process.server) {
      return null;
    }

    const token = localStorage.getItem("auth_token");
    if (!token) return null;

    try {
      const decoded: { eat: number } = jwtDecode(token);
      if (decoded.eat * 1000 < Date.now()) {
        localStorage.removeItem("auth_token");
        return null;
      }
      return token;
    } catch (error) {
      localStorage.removeItem("auth_token");
      return null;
    }
  };

  if (process.server) {
    return;
  }

  if (!isAuthenticated()) {
    return navigateTo("/login");
  }
});
