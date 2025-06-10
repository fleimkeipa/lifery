import { createFetch } from "@vueuse/core";

export default createFetch({
  baseUrl: useRuntimeConfig().public.apiBase,
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
