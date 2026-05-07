import { useEffect, useState } from 'react'
import { useNavigate, Link } from 'react-router-dom'
import { useAuth } from '../context/AuthContext'
import { getHabits, createHabit, deleteHabit, checkInHabit } from '../api/habits'
import CreateHabitForm from '../components/CreateHabitForm'
import HabitGrid from '../components/HabitGrid'

const DAYS = ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun']

function formatDays(bitmask) {
  return DAYS.filter((_, i) => (bitmask & (1 << i)) !== 0).join(', ')
}

export default function Habits() {
  const { user, token, clearAuth } = useAuth()
  const navigate = useNavigate()
  const [habits, setHabits] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  useEffect(() => {
    fetchHabits()
  }, [])

  async function fetchHabits() {
    try {
      const data = await getHabits(token)
      setHabits(data)
    } catch (err) {
      setError(err.message)
    } finally {
      setLoading(false)
    }
  }

  async function handleCreate(data) {
    const habit = await createHabit(token, data)
    setHabits(prev => [...prev, habit])
  }

  async function handleDelete(id) {
    await deleteHabit(token, id)
    setHabits(prev => prev.filter(h => h.id !== id))
  }

  async function handleCheckIn(id) {
    await checkInHabit(token, id)
    setHabits(prev => prev.map(h => h.id === id ? { ...h, checked_in_today: true } : h))
  }

  function handleLogout() {
    clearAuth()
    navigate('/login')
  }

  return (
    <div className="dashboard-container">
      <header className="dashboard-header">
        <h1>debitask</h1>
        <div className="header-right">
          <Link to="/dashboard" className="nav-link">Tasks</Link>
          <Link to="/habits" className="nav-link nav-link--active">Habits</Link>
          <span>{user?.email}</span>
          <button onClick={handleLogout}>Logout</button>
        </div>
      </header>

      <main>
        <div className="dashboard-layout">
          <section className="create-section">
            <h2>New habit</h2>
            <CreateHabitForm onSubmit={handleCreate} />
          </section>

          <section className="tasks-section">
            {error && <p className="error">{error}</p>}
            {loading ? (
              <p className="tasks-empty">Loading...</p>
            ) : habits.length === 0 ? (
              <p className="tasks-empty">No habits yet. Create one to get started.</p>
            ) : (
              habits.map(habit => (
                <div key={habit.id} className="habit-row">
                  <div className="habit-row__header">
                    <div className="habit-row__info">
                      <p className="habit-row__name">{habit.name}</p>
                      <p className="habit-row__meta">{formatDays(habit.days)} · ${habit.penalty.toFixed(2)} penalty</p>
                    </div>
                    <div className="habit-row__actions">
                      {!habit.checked_in_today && (
                        <button className="btn-complete" onClick={() => handleCheckIn(habit.id)}>
                          Done today
                        </button>
                      )}
                      {habit.checked_in_today && (
                        <span className="habit-done-badge">✓ Done</span>
                      )}
                      <button className="btn-delete" onClick={() => handleDelete(habit.id)}>Delete</button>
                    </div>
                  </div>
                  <HabitGrid habit={habit} logs={null} />
                </div>
              ))
            )}
          </section>
        </div>
      </main>
    </div>
  )
}
