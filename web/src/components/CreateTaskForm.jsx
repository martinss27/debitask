import { useState } from 'react'

export default function CreateTaskForm({ onSubmit }) {
  const [title, setTitle] = useState('')
  const [description, setDescription] = useState('')
  const [date, setDate] = useState('')
  const [time, setTime] = useState('')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')

  async function handleSubmit(e) {
    e.preventDefault()
    setError('')
    setLoading(true)
    try {
      const deadline = new Date(`${date}T${time || '00:00'}`).toISOString()
      await onSubmit({ title, description, deadline })
      setTitle('')
      setDescription('')
      setDate('')
      setTime('')
    } catch (err) {
      setError(err.message)
    } finally {
      setLoading(false)
    }
  }

  return (
    <form className="create-task-form" onSubmit={handleSubmit}>
      <div className="field">
        <label>Title</label>
        <input
          type="text"
          value={title}
          onChange={e => setTitle(e.target.value)}
          placeholder="What needs to be done?"
          required
        />
      </div>
      <div className="field">
        <label>Description <span className="optional">optional</span></label>
        <input
          type="text"
          value={description}
          onChange={e => setDescription(e.target.value)}
          placeholder="Any details?"
        />
      </div>
      <div className="field">
        <label>Deadline</label>
        <div className="deadline-inputs">
          <input
            type="date"
            value={date}
            onChange={e => setDate(e.target.value)}
            required
          />
          <input
            type="time"
            value={time}
            onChange={e => setTime(e.target.value)}
          />
        </div>
      </div>
      {error && <p className="error">{error}</p>}
      <button type="submit" disabled={loading}>
        {loading ? 'Creating...' : 'Create task'}
      </button>
    </form>
  )
}
