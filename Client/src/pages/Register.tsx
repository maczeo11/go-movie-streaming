import { useEffect, useState } from 'react'
import type { FormEvent } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { useAuth } from '../auth/AuthContext'
import { api, ApiRequestError } from '../api'
import type { Genre } from '../types'

export default function Register() {
  const { register } = useAuth()
  const navigate = useNavigate()

  const [firstName, setFirstName] = useState('')
  const [lastName, setLastName] = useState('')
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [genres, setGenres] = useState<Genre[]>([])
  const [picked, setPicked] = useState<number[]>([])
  const [error, setError] = useState('')
  const [busy, setBusy] = useState(false)

  useEffect(() => {
    api
      .genres()
      .then(setGenres)
      .catch(() => undefined)
  }, [])

  const toggleGenre = (id: number) => {
    setPicked((prev) =>
      prev.includes(id) ? prev.filter((g) => g !== id) : [...prev, id],
    )
  }

  const submit = async (e: FormEvent) => {
    e.preventDefault()
    setBusy(true)
    setError('')
    try {
      await register({
        first_name: firstName,
        last_name: lastName,
        email,
        password,
        favourite_genres: genres.filter((g) => picked.includes(g.genre_id)),
      })
      navigate('/')
    } catch (err) {
      setError(
        err instanceof ApiRequestError
          ? err.message
          : 'Something went wrong. Try again.',
      )
    } finally {
      setBusy(false)
    }
  }

  return (
    <div className="auth-wrap">
      <form className="card auth-card" onSubmit={submit}>
        <h1>Create your account</h1>
        <p className="muted">
          Pick a few genres you love — we'll use them for recommendations.
        </p>

        <div className="row-2">
          <label>
            First name
            <input
              required
              minLength={2}
              value={firstName}
              onChange={(e) => setFirstName(e.target.value)}
            />
          </label>
          <label>
            Last name
            <input
              required
              minLength={2}
              value={lastName}
              onChange={(e) => setLastName(e.target.value)}
            />
          </label>
        </div>

        <label>
          Email
          <input
            type="email"
            required
            value={email}
            onChange={(e) => setEmail(e.target.value)}
          />
        </label>

        <label>
          Password
          <input
            type="password"
            required
            minLength={6}
            value={password}
            onChange={(e) => setPassword(e.target.value)}
          />
        </label>

        <div className="genre-picker">
          {genres.map((g) => (
            <button
              type="button"
              key={g.genre_id}
              className={`pill pickable${picked.includes(g.genre_id) ? ' on' : ''}`}
              onClick={() => toggleGenre(g.genre_id)}
            >
              {g.genre_name}
            </button>
          ))}
        </div>

        {error && <p className="error-text">{error}</p>}

        <button className="btn btn-primary btn-block" disabled={busy}>
          {busy ? 'Creating…' : 'Create account'}
        </button>

        <p className="muted">
          Already have an account? <Link to="/login">Sign in</Link>
        </p>
      </form>
    </div>
  )
}
