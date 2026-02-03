const API_BASE = import.meta.env.VITE_API_BASE || ''

export async function request(path, options = {}) {
  const response = await fetch(`${API_BASE}${path}`, {
    headers: {
      'Content-Type': 'application/json',
      ...(options.headers || {})
    },
    ...options
  })

  if (response.status === 204) {
    return null
  }

  const data = await response.json().catch(() => null)
  if (!response.ok) {
    const message = data?.error || 'Request failed'
    throw new Error(message)
  }
  return data
}
