import { ref } from 'vue'
import jwtDecode from 'jwt-decode'

export const useAuth = () => {
  const user = ref<{ id: string; username: string } | null>(null)

  if (process.client) {
    const token = localStorage.getItem('auth_token')
    if (token) {
      try {
        const decoded = jwtDecode(token) as { id: string; username: string }
        user.value = decoded
      } catch (error) {
        console.error('Failed to decode token:', error)
      }
    }
  }

  return {
    user
  }
} 