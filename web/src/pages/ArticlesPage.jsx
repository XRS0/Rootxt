import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { listPublishedArticles } from '../api/articles'

export default function ArticlesPage() {
  const [articles, setArticles] = useState([])
  const [error, setError] = useState('')

  useEffect(() => {
    let active = true
    listPublishedArticles()
      .then((data) => {
        if (active) {
          setArticles(data)
        }
      })
      .catch((err) => {
        if (active) {
          setError(err.message)
        }
      })

    return () => {
      active = false
    }
  }, [])

  return (
    <section className="page">
      <h1>Articles</h1>
      <p className="muted">Chronological system notes and build logs.</p>

      {error && <p className="muted">{error}</p>}
      {!error && articles.length === 0 && <p className="muted">No published articles yet.</p>}

      <ul className="article-list">
        {articles.map((article) => (
          <li key={article.slug}>
            <Link to={`/articles/${article.slug}`}>{article.title}</Link>
            <p className="excerpt">{article.excerpt}</p>
          </li>
        ))}
      </ul>
    </section>
  )
}
