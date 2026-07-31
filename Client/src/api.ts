import type { Genre, Movie, Paginated, RegisterPayload, User } from './types'

const BASE = import.meta.env.VITE_API_BASE ?? '/api'

interface ApiError {
  error: string
  details?: string
}

export class ApiRequestError extends Error {
  status: number
  details?: string

  constructor(status: number, message: string, details?: string) {
    super(message)
    this.status = status
    this.details = details
  }
}

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const res = await fetch(`${BASE}${path}`, {
    credentials: 'include',
    headers: init?.body ? { 'Content-Type': 'application/json' } : undefined,
    ...init,
  })

  if (!res.ok) {
    let err: ApiError = { error: res.statusText }
    try {
      err = (await res.json()) as ApiError
    } catch {
      // non-JSON error body, keep the status text
    }
    throw new ApiRequestError(res.status, err.error, err.details)
  }

  if (res.status === 204) return undefined as T
  return (await res.json()) as T
}

export const api = {
  movies: (params: { q?: string; genre?: string; page?: number; limit?: number; sort?: string }) => {
    const query = new URLSearchParams()
    if (params.q) query.set('q', params.q)
    if (params.genre) query.set('genre', params.genre)
    if (params.page) query.set('page', String(params.page))
    if (params.limit) query.set('limit', String(params.limit))
    if (params.sort) query.set('sort', params.sort)
    const qs = query.toString()
    return request<Paginated<Movie>>(`/movies${qs ? `?${qs}` : ''}`)
  },

  movie: (imdbId: string) => request<Movie>(`/movie/${imdbId}`),
  genres: () => request<Genre[]>('/genres'),
  media: () => request<{ data: string[]; count: number }>('/media'),

  register: (payload: RegisterPayload) =>
    request<{ inserted_id?: string }>('/user', {
      method: 'POST',
      body: JSON.stringify(payload),
    }),

  login: (email: string, password: string) =>
    request<User>('/login', {
      method: 'POST',
      body: JSON.stringify({ email, password }),
    }),

  logout: (userId: string) =>
    request<{ message: string }>('/logout', {
      method: 'POST',
      body: JSON.stringify({ user_id: userId }),
    }),

  profile: () => request<User>('/profile'),
  updateGenres: (favourite_genres: Genre[]) =>
    request<{ message: string }>('/profile/genres', {
      method: 'PATCH',
      body: JSON.stringify({ favourite_genres }),
    }),

  recommended: () => request<Movie[]>('/recommendedmovies'),

  addMovie: (movie: Omit<Movie, '_id'>) =>
    request<{ inserted_id?: string }>('/addmovie', {
      method: 'POST',
      body: JSON.stringify(movie),
    }),

  updateReview: (imdbId: string, adminReview: string) =>
    request<{ ranking_name: string; admin_review: string }>(
      `/updatereview/${imdbId}`,
      {
        method: 'PATCH',
        body: JSON.stringify({ admin_review: adminReview }),
      },
    ),
}
