import { useEffect, useState } from 'react'
import { useAuth } from '../auth/AuthContext'
import { api, ApiRequestError } from '../api'
import type { Genre } from '../types'

export default function Profile() {
  const { user } = useAuth()
  const [allGenres, setAllGenres] = useState<Genre[]>([])
  const [picked, setPicked] = useState<number[]>([])
  const [msg, setMsg] = useState('')
  const [error, setError] = useState('')
  const [busy, setBusy] = useState(false)

  useEffect(() => {
    api.genres().then(setAllGenres).catch(() => undefined)
  }, [])

  useEffect(() => {
    if (user) setPicked(user.favourite_genres.map((g) => g.genre_id))
  }, [user])

  const toggleGenre = (id: number) => {
    setPicked((prev) => (prev.includes(id) ? prev.filter((g) => g !== id) : [...prev, id]))
  }

  const save = async () => {
    setBusy(true)
    setMsg('')
    setError('')
    try {
      const selected = allGenres.filter((g) => picked.includes(g.genre_id))
      await api.updateGenres(selected)
      setMsg('Saved. Recommendations will pick up the new genres.')
    } catch (err) {
      setError(err instanceof ApiRequestError ? err.message : 'Failed to save.')
    } finally {
      setBusy(false)
    }
  }

  return (
    <div className="profile">
      <h1>Your profile</h1>

      <div className="card profile-card">
        <p>
          <strong>
            {user?.first_name} {user?.last_name}
          </strong>
        </p>
        <p className="muted">{user?.email}</p>
      </div>

      <section className="card">
        <h2>Favourite genres</h2>
        <p className="muted">
          The recommendation engine picks movies from these. Fewer, focused
          picks give better results than picking everything.
        </p>
        <div className="genre-picker">
          {allGenres.map((g) => (
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
        {msg && <p className="ok-text">{msg}</p>}
        <button className="btn btn-primary" onClick={save} disabled={busy}>
          {busy ? 'Saving…' : 'Save genres'}
        </button>
      </section>
    </div>
  )
}
