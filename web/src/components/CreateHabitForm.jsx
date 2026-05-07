import { useState } from 'react'

const DAYS = ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun']

export default function CreateHabitForm({ onSubmit }) {
  const [name, setName] = useState('')
  const [selectedDays, setSelectedDays] = useState(0)
  const [penalty, setPenalty] = useState('')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')

  function toggleDay(index) {
    setSelectedDays(prev => prev ^ (1 << index))
  }

  function isDaySelected(index) {
    return (selectedDays & (1 << index)) !== 0
  }

  async function handleSubmit(e) {
    e.preventDefault()
    if (selectedDays === 0) {
      setError('Select at least one day')
      return
    }
    setError('')
    setLoading(true)
    try {
      await onSubmit({ name, days: selectedDays, penalty: parseFloat(penalty) })
      setName('')
      setSelectedDays(0)
      setPenalty('')
    } catch (err) {
      setError(err.message)
    } finally {
      setLoading(false)
    }
  }

  return (
    <form className="create-task-form" onSubmit={handleSubmit}>
      <div className="field">
        <label>Name</label>
        <input
          type="text"
          value={name}
          onChange={e => setName(e.target.value)}
          placeholder="e.g. Read 10 pages"
          required
        />
      </div>
      <div className="field">
        <label>Days</label>
        <div className="day-selector">
          {DAYS.map((day, i) => (
            <button
              key={day}
              type="button"
              className={`day-btn ${isDaySelected(i) ? 'day-btn--active' : ''}`}
              onClick={() => toggleDay(i)}
            >
              {day}
            </button>
          ))}
        </div>
      </div>
      <div className="field">
        <label>Penalty ($)</label>
        <input
          type="number"
          min="0"
          step="0.01"
          value={penalty}
          onChange={e => setPenalty(e.target.value)}
          placeholder="0.50"
          required
        />
      </div>
      {error && <p className="error">{error}</p>}
      <button type="submit" disabled={loading}>
        {loading ? 'Creating...' : 'Create habit'}
      </button>
    </form>
  )
}
