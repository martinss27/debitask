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

export const getHabits = (token) =>
  request('/api/habits', token)

export const createHabit = (token, data) =>
  request('/api/habits', token, { method: 'POST', body: JSON.stringify(data) })

export const deleteHabit = (token, id) =>
  request(`/api/habits/${id}`, token, { method: 'DELETE' })

export const checkInHabit = (token, id) =>
  request(`/api/habits/${id}/checkin`, token, { method: 'POST' })
