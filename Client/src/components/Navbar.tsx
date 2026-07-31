import { Link, NavLink, useNavigate } from 'react-router-dom'
import { useAuth } from '../auth/AuthContext'

export default function Navbar() {
  const { user, logout } = useAuth()
  const navigate = useNavigate()

  const handleLogout = async () => {
    await logout()
    navigate('/')
  }

  return (
    <header className="navbar">
      <Link to="/" className="brand">
        <svg viewBox="0 0 64 64" aria-hidden="true">
          <rect width="64" height="64" rx="14" fill="var(--accent)" />
          <polygon points="24,18 50,32 24,46" fill="#fff" />
        </svg>
        MagicStream
      </Link>

      <nav className="nav-links">
        <NavLink to="/" end>
          Browse
        </NavLink>
        {user && (
          <NavLink to="/profile">My genres</NavLink>
        )}
        {user?.role === 'ADMIN' && <NavLink to="/admin">Admin</NavLink>}
      </nav>

      <div className="nav-user">
        {user ? (
          <>
            <span className="nav-name">
              {user.first_name} {user.last_name}
              {user.role === 'ADMIN' && <span className="admin-tag">admin</span>}
            </span>
            <button className="btn btn-ghost" onClick={handleLogout}>
              Sign out
            </button>
          </>
        ) : (
          <>
            <Link className="btn btn-ghost" to="/login">
              Sign in
            </Link>
            <Link className="btn btn-primary" to="/register">
              Join free
            </Link>
          </>
        )}
      </div>
    </header>
  )
}
