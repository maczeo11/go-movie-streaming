interface PaginationProps {
  page: number
  totalPages: number
  onPage: (page: number) => void
}

export default function Pagination({ page, totalPages, onPage }: PaginationProps) {
  if (totalPages <= 1) return null

  const pages: number[] = []
  const start = Math.max(1, Math.min(page - 2, totalPages - 4))
  for (let i = start; i <= Math.min(totalPages, start + 4); i++) pages.push(i)

  return (
    <div className="pagination">
      <button
        className="btn btn-ghost"
        disabled={page <= 1}
        onClick={() => onPage(page - 1)}
      >
        &larr; Prev
      </button>
      {pages.map((p) => (
        <button
          key={p}
          className={`page-btn${p === page ? ' active' : ''}`}
          onClick={() => onPage(p)}
        >
          {p}
        </button>
      ))}
      <button
        className="btn btn-ghost"
        disabled={page >= totalPages}
        onClick={() => onPage(page + 1)}
      >
        Next &rarr;
      </button>
    </div>
  )
}
