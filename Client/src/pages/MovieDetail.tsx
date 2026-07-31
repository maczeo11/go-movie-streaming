import { useCallback, useEffect, useState } from 'react'
import { Link, useParams } from 'react-router-dom'
import { api, ApiRequestError } from '../api'
import type { Movie } from '../types'
import { useAuth } from '../auth/AuthContext'

const MEDIA_BASE = import.meta.env.VITE_MEDIA_BASE ?? '/api/media'

export default function MovieDetail() {
  const { imdbId } = useParams<{ imdbId: string }>()
  const { user } = useAuth()
  const [movie, setMovie] = useState<Movie | null>(null)
  const [media, setMedia] = useState<string[]>([])
  const [error, setError] = useState('')

  const [review, setReview] = useState('')
  const [updating, setUpdating] = useState(false)
  const [reviewMsg, setReviewMsg] = useState('')

  useEffect(() => {
    if (!imdbId) return
    api
      .movie(imdbId)
      .then(setMovie)
      .catch(() => setError('Movie not found.'))
  }, [imdbId])

  useEffect(() => {
    api
      .media()
      .then((res) => setMedia(res.data))
      .catch(() => undefined)
  }, [])

  const saveReview = useCallback(async () => {
    if (!movie || !imdbId || !review.trim()) return
    setUpdating(true)
    setReviewMsg('')
    try {
      const res = await api.updateReview(imdbId, review.trim())
      setMovie((m) => (m ? { ...m, admin_review: review.trim(), ranking: { ...m.ranking, ranking_name: res.ranking_name } } : m))
      setReviewMsg(`Ranked: ${res.ranking_name}`)
    } catch (err) {
      setReviewMsg(
        err instanceof ApiRequestError && err.status === 500
          ? 'Review ranking needs an OpenAI API key on the server.'
          : 'Could not update the review.',
      )
    } finally {
      setUpdating(false)
    }
  }, [movie, imdbId, review])

  if (error) return <p className="error-text">{error}</p>
  if (!movie) return <p className="muted">Loading…</p>

  const trailerUrl = `https://www.youtube.com/embed/${movie.youtube_id}`
  const clip = media.length > 0 ? `${MEDIA_BASE}/${encodeURIComponent(media[0])}` : null
  const isAdmin = user?.role === 'ADMIN'

  return (
    <div className="detail">
      <div className="detail-hero">
        <img className="detail-poster" src={movie.poster_path} alt={`${movie.title} poster`} />
        <div className="detail-info">
          <h1>{movie.title}</h1>
          <div className="genre-pills">
            {movie.genre.map((g) => (
              <span key={g.genre_id} className="pill">
                {g.genre_name}
              </span>
            ))}
          </div>
          {movie.ranking && movie.ranking.ranking_value <= 5 && (
            <p className="rank-line">
              <span className={`badge badge-${movie.ranking.ranking_value}`}>
                {movie.ranking.ranking_name}
              </span>
              <span className="muted"> by the MagicStream curators</span>
            </p>
          )}
          <p className="imdb-line">
            <a
              href={`https://www.imdb.com/title/${movie.imdb_id}/`}
              target="_blank"
              rel="noreferrer"
            >
              View on IMDb &rarr;
            </a>
          </p>
        </div>
      </div>

      {movie.admin_review && (
        <section className="card review">
          <h2>What the curators say</h2>
          <p>{movie.admin_review}</p>
        </section>
      )}

      {clip && (
        <section className="section">
          <h2 className="section-title">Watch</h2>
          <div className="video-shell">
            <video controls preload="metadata" src={clip}>
              Your browser does not support HTML5 video.
            </video>
            <p className="muted">Demo clip streamed from the Go backend.</p>
          </div>
        </section>
      )}

      <section className="section">
        <h2 className="section-title">Trailer</h2>
        <div className="video-shell">
          <iframe
            src={trailerUrl}
            title={`${movie.title} trailer`}
            allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture"
            allowFullScreen
          />
        </div>
      </section>

      {isAdmin && (
        <section className="card review-editor">
          <h2>Curator review</h2>
          <textarea
            rows={4}
            placeholder="Write your review... the server will rank it via AI."
            value={review}
            onChange={(e) => setReview(e.target.value)}
          />
          <div className="row-between">
            <button
              className="btn btn-primary"
              disabled={updating || !review.trim()}
              onClick={saveReview}
            >
              {updating ? 'Ranking…' : 'Save & rank review'}
            </button>
            {reviewMsg && <span className="muted">{reviewMsg}</span>}
          </div>
        </section>
      )}

      <p className="muted">
        <Link to="/">&larr; Back to catalog</Link>
      </p>
    </div>
  )
}
