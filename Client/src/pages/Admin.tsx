import { useEffect, useState } from 'react'
import type { FormEvent } from 'react'
import { api, ApiRequestError } from '../api'
import { RANKING_LABELS } from '../types'
import type { Genre } from '../types'

const empty = {
  imdb_id: '',
  title: '',
  poster_path: '',
  youtube_id: '',
  genre: [] as Genre[],
  admin_review: '',
  ranking: { ranking_value: 1, ranking_name: RANKING_LABELS[1] },
}

export default function Admin() {
  const [genres, setGenres] = useState<Genre[]>([])
  const [form, setForm] = useState(empty)
  const [picked, setPicked] = useState<number[]>([])
  const [msg, setMsg] = useState('')
  const [error, setError] = useState('')
  const [busy, setBusy] = useState(false)

  useEffect(() => {
    api.genres().then(setGenres).catch(() => undefined)
  }, [])

  const set = <K extends keyof typeof empty>(key: K, value: (typeof empty)[K]) => {
    setForm((f) => ({ ...f, [key]: value }))
  }

  const toggleGenre = (id: number) => {
    setPicked((prev) => (prev.includes(id) ? prev.filter((g) => g !== id) : [...prev, id]))
  }

  const submit = async (e: FormEvent) => {
    e.preventDefault()
    setBusy(true)
    setMsg('')
    setError('')
    try {
      const payload = {
        ...form,
        genre: genres.filter((g) => picked.includes(g.genre_id)),
        ranking: {
          ranking_value: form.ranking.ranking_value,
          ranking_name: RANKING_LABELS[form.ranking.ranking_value],
        },
      }
      await api.addMovie(payload)
      setMsg(`Added "${payload.title}" to the catalog.`)
      setForm(empty)
      setPicked([])
    } catch (err) {
      setError(
        err instanceof ApiRequestError
          ? `${err.message}${err.details ? ` — ${err.details}` : ''}`
          : 'Failed to add the movie.',
      )
    } finally {
      setBusy(false)
    }
  }

  return (
    <div className="admin">
      <h1>Admin — add a movie</h1>
      <form className="card admin-form" onSubmit={submit}>
        <label>
          Title
          <input
            required
            value={form.title}
            onChange={(e) => set('title', e.target.value)}
          />
        </label>

        <div className="row-2">
          <label>
            IMDb ID
            <input
              required
              placeholder="tt0111161"
              value={form.imdb_id}
              onChange={(e) => set('imdb_id', e.target.value)}
            />
          </label>
          <label>
            YouTube trailer ID
            <input
              required
              placeholder="dQw4w9WgXcQ"
              value={form.youtube_id}
              onChange={(e) => set('youtube_id', e.target.value)}
            />
          </label>
        </div>

        <label>
          Poster URL
          <input
            required
            type="url"
            placeholder="https://image.tmdb.org/t/p/w500/..."
            value={form.poster_path}
            onChange={(e) => set('poster_path', e.target.value)}
          />
        </label>

        <label>
          Genres
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
        </label>

        <label>
          Ranking
          <select
            value={form.ranking.ranking_value}
            onChange={(e) =>
              set('ranking', {
                ranking_value: Number(e.target.value),
                ranking_name: RANKING_LABELS[Number(e.target.value)],
              })
            }
          >
            {[1, 2, 3, 4, 5].map((v) => (
              <option key={v} value={v}>
                {v} — {RANKING_LABELS[v]}
              </option>
            ))}
          </select>
        </label>

        <label>
          Curator review (optional)
          <textarea
            rows={3}
            value={form.admin_review}
            onChange={(e) => set('admin_review', e.target.value)}
          />
        </label>

        {error && <p className="error-text">{error}</p>}
        {msg && <p className="ok-text">{msg}</p>}

        <button className="btn btn-primary btn-block" disabled={busy}>
          {busy ? 'Adding…' : 'Add to catalog'}
        </button>
      </form>
    </div>
  )
}
