export interface Genre {
  genre_id: number
  genre_name: string
}

export interface Ranking {
  ranking_value: number
  ranking_name: string
}

export interface Movie {
  _id?: string
  imdb_id: string
  title: string
  poster_path: string
  youtube_id: string
  genre: Genre[]
  admin_review: string
  ranking: Ranking
}

export interface User {
  user_id: string
  first_name: string
  last_name: string
  email: string
  role: 'ADMIN' | 'USER'
  favourite_genres: Genre[]
}

export interface RegisterPayload {
  first_name: string
  last_name: string
  email: string
  password: string
  favourite_genres: Genre[]
}

export interface Paginated<T> {
  data: T[]
  meta: {
    page: number
    limit: number
    total: number
    total_pages: number
  }
}

export const RANKING_LABELS: Record<number, string> = {
  1: 'Must Watch',
  2: 'Highly Recommended',
  3: 'Recommended',
  4: 'Average',
  5: 'Skip It',
}
