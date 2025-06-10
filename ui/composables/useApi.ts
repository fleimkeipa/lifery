import { createFetch } from "@vueuse/core";


export default createFetch({
  baseUrl: process.env.NUXT_API_BASE_URL || "http://localhost:8080",
  options: {
    beforeFetch({ options }) {
      const token = typeof window !== 'undefined' ? localStorage.getItem('auth_token') : null;
      if (token) {
        options.headers = {
          ...options.headers,
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json',
        };
      } else {
        options.headers = {
          ...options.headers,
          'Content-Type': 'application/json',
        };
      }
      return { options };
    },
  },
});
