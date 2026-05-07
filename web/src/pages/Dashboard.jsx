import { useAuth } from '../context/AuthContext'
import { useNavigate } from 'react-router-dom'

export default function Dashboard() {
  const { user, clearAuth } = useAuth()
  const navigate = useNavigate()

  function handleLogout() {
    clearAuth()
    navigate('/login')
  }

  return (
    <div className="dashboard-container">
      <header className="dashboard-header">
        <h1>debitask</h1>
        <div className="header-right">
          <span>{user?.email}</span>
          <button onClick={handleLogout}>Logout</button>
        </div>
      </header>
      <main>
        <p>Tasks coming soon.</p>
      </main>
    </div>
  )
}
