import { useCallback, useEffect, useRef, useState } from 'react'
import { api } from '../api'
import type { Genre, Movie } from '../types'
import MovieCard from '../components/MovieCard'
import Pagination from '../components/Pagination'
import { useAuth } from '../auth/AuthContext'

export default function Home() {
  const { user } = useAuth()
  const [movies, setMovies] = useState<Movie[]>([])
  const [genres, setGenres] = useState<Genre[]>([])
  const [query, setQuery] = useState('')
  const [genre, setGenre] = useState('')
  const [page, setPage] = useState(1)
  const [totalPages, setTotalPages] = useState(1)
  const [loading, setLoading] = useState(true)
  const [recommended, setRecommended] = useState<Movie[]>([])
  const [error, setError] = useState('')

  const debounceRef = useRef<ReturnType<typeof setTimeout> | undefined>(undefined)

  // Debounce the search box so we don't hammer the API per keystroke.
  useEffect(() => {
    clearTimeout(debounceRef.current)
    debounceRef.current = setTimeout(() => setPage(1), 350)
    return () => clearTimeout(debounceRef.current)
  }, [query, genre])

  const load = useCallback(async () => {
    setLoading(true)
    setError('')
    try {
      const res = await api.movies({ q: query || undefined, genre: genre || undefined, page })
      setMovies(res.data)
      setTotalPages(res.meta.total_pages)
    } catch {
      setError('Could not load movies. Is the API running?')
    } finally {
      setLoading(false)
    }
  }, [query, genre, page])

  useEffect(() => {
    load()
  }, [load])

  useEffect(() => {
    api
      .genres()
      .then(setGenres)
      .catch(() => undefined)
  }, [])

  useEffect(() => {
    if (user) {
      api
        .recommended()
        .then((m) => setRecommended(m.slice(0, 5)))
        .catch(() => undefined)
    }
  }, [user])

  return (
    <div>
      <section className="hero">
        <h1>Stream something good tonight.</h1>
        <p>
          Browse the catalog, watch the trailer, read what the curators actually think.
        </p>
        <div className="search-row">
          <input
            className="search-input"
            type="search"
            placeholder="Search by title..."
            value={query}
            onChange={(e) => setQuery(e.target.value)}
          />
          <select
            className="search-input select"
            value={genre}
            onChange={(e) => setGenre(e.target.value)}
          >
            <option value="">All genres</option>
            {genres.map((g) => (
              <option key={g.genre_id} value={g.genre_name}>
                {g.genre_name}
              </option>
            ))}
          </select>
        </div>
      </section>

      {user && recommended.length > 0 && (
        <section className="section">
          <h2 className="section-title">Picked for you</h2>
          <div className="grid">
            {recommended.map((m) => (
              <MovieCard key={m.imdb_id} movie={m} />
            ))}
          </div>
        </section>
      )}

      <section className="section">
        <h2 className="section-title">
          {query || genre ? 'Search results' : 'In the catalog'}
        </h2>
        {error && <p className="error-text">{error}</p>}
        {loading ? (
          <p className="muted">Loading…</p>
        ) : movies.length === 0 ? (
          <p className="muted">Nothing matches that. Try a different title or genre.</p>
        ) : (
          <>
            <div className="grid">
              {movies.map((m) => (
                <MovieCard key={m.imdb_id} movie={m} />
              ))}
            </div>
            <Pagination page={page} totalPages={totalPages} onPage={setPage} />
          </>
        )}
      </section>
    </div>
  )
}
