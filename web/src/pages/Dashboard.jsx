import { useEffect, useState } from 'react'
import { useNavigate, Link } from 'react-router-dom'
import { useAuth } from '../context/AuthContext'
import { getTasks, createTask, updateTask, deleteTask } from '../api/tasks'
import CreateTaskForm from '../components/CreateTaskForm'
import TaskCard from '../components/TaskCard'

export default function Dashboard() {
  const { user, token, clearAuth } = useAuth()
  const navigate = useNavigate()
  const [tasks, setTasks] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  useEffect(() => {
    fetchTasks()
  }, [])

  async function fetchTasks() {
    try {
      const data = await getTasks(token)
      setTasks(data)
    } catch (err) {
      setError(err.message)
    } finally {
      setLoading(false)
    }
  }

  async function handleCreate(data) {
    const task = await createTask(token, data)
    setTasks(prev => [task, ...prev])
  }

  async function handleComplete(id) {
    const task = await updateTask(token, id, { status: 'completed' })
    setTasks(prev => prev.map(t => t.id === id ? task : t))
  }

  async function handleDelete(id) {
    await deleteTask(token, id)
    setTasks(prev => prev.filter(t => t.id !== id))
  }

  function handleLogout() {
    clearAuth()
    navigate('/login')
  }

  const pending = tasks.filter(t => t.status === 'pending')
  const overdue = tasks.filter(t => t.status === 'overdue')
  const completed = tasks.filter(t => t.status === 'completed')

  return (
    <div className="dashboard-container">
      <header className="dashboard-header">
        <h1>debitask</h1>
        <div className="header-right">
          <Link to="/dashboard" className="nav-link nav-link--active">Tasks</Link>
          <Link to="/habits" className="nav-link">Habits</Link>
          <span>{user?.email}</span>
          <button onClick={handleLogout}>Logout</button>
        </div>
      </header>

      <main>
        <div className="dashboard-layout">
          <section className="create-section">
            <h2>New task</h2>
            <CreateTaskForm onSubmit={handleCreate} />
          </section>

          <section className="tasks-section">
            {error && <p className="error">{error}</p>}
            {loading ? (
              <p className="tasks-empty">Loading...</p>
            ) : (
              <>
                {overdue.length > 0 && (
                  <div className="task-group">
                    <h3 className="task-group__title task-group__title--overdue">
                      Overdue <span>{overdue.length}</span>
                    </h3>
                    {overdue.map(t => (
                      <TaskCard key={t.id} task={t} onComplete={handleComplete} onDelete={handleDelete} />
                    ))}
                  </div>
                )}

                {pending.length > 0 && (
                  <div className="task-group">
                    <h3 className="task-group__title task-group__title--pending">
                      Pending <span>{pending.length}</span>
                    </h3>
                    {pending.map(t => (
                      <TaskCard key={t.id} task={t} onComplete={handleComplete} onDelete={handleDelete} />
                    ))}
                  </div>
                )}

                {completed.length > 0 && (
                  <div className="task-group">
                    <h3 className="task-group__title task-group__title--completed">
                      Completed <span>{completed.length}</span>
                    </h3>
                    {completed.map(t => (
                      <TaskCard key={t.id} task={t} onComplete={handleComplete} onDelete={handleDelete} />
                    ))}
                  </div>
                )}

                {tasks.length === 0 && (
                  <p className="tasks-empty">No tasks yet. Create one to get started.</p>
                )}
              </>
            )}
          </section>
        </div>
      </main>
    </div>
  )
}
