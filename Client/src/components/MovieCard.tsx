import { Link } from 'react-router-dom'
import type { Movie } from '../types'

export default function MovieCard({ movie }: { movie: Movie }) {
  const rank = movie.ranking?.ranking_value

  return (
    <Link to={`/movie/${movie.imdb_id}`} className="movie-card">
      <div className="poster-wrap">
        <img
          src={movie.poster_path}
          alt={`${movie.title} poster`}
          loading="lazy"
          onError={(e) => {
            ;(e.currentTarget as HTMLImageElement).src =
              'https://placehold.co/400x600/18181b/e4e4e7?text=No+poster'
          }}
        />
        {rank && rank <= 5 && (
          <span className={`badge badge-${rank}`}>{movie.ranking.ranking_name}</span>
        )}
      </div>
      <div className="movie-card-body">
        <h3>{movie.title}</h3>
        <div className="genre-pills">
          {movie.genre.slice(0, 3).map((g) => (
            <span key={g.genre_id} className="pill">
              {g.genre_name}
            </span>
          ))}
        </div>
      </div>
    </Link>
  )
}
