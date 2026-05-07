async function request(path, token, options = {}) {
  const res = await fetch(path, {
    ...options,
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${token}`,
      ...options.headers,
    },
  })
  if (!res.ok) {
    const text = await res.text()
    throw new Error(text.trim() || 'Request failed')
  }
  if (res.status === 204) return null
  return res.json()
}

export const getTasks = (token) =>
  request('/api/tasks', token)

export const createTask = (token, data) =>
  request('/api/tasks', token, { method: 'POST', body: JSON.stringify(data) })

export const updateTask = (token, id, data) =>
  request(`/api/tasks/${id}`, token, { method: 'PUT', body: JSON.stringify(data) })

export const deleteTask = (token, id) =>
  request(`/api/tasks/${id}`, token, { method: 'DELETE' })
