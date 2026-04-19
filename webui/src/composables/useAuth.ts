// Auth composable — manages JWT token in localStorage
export const TOKEN_KEY = 'auth_token'

export function useAuth() {
  function getToken(): string | null {
    return localStorage.getItem(TOKEN_KEY)
  }

  function setToken(token: string): void {
    localStorage.setItem(TOKEN_KEY, token)
  }

  function clearToken(): void {
    localStorage.removeItem(TOKEN_KEY)
  }

  function isAuthenticated(): boolean {
    const token = getToken()
    if (!token) return false
    // Quick expiry check by decoding payload (no verification — just for UI gating)
    try {
      const payload = JSON.parse(atob(token.split('.')[1]))
      return payload.exp > Math.floor(Date.now() / 1000)
    } catch {
      return false
    }
  }

  function getUser(): { id: number; email: string; username: string } | null {
    const token = getToken()
    if (!token) return null
    try {
      const payload = JSON.parse(atob(token.split('.')[1]))
      return { id: payload.uid, email: payload.sub, username: payload.sub }
    } catch {
      return null
    }
  }

  return { getToken, setToken, clearToken, isAuthenticated, getUser }
}
