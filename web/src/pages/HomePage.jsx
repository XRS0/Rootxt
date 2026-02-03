import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { listPublishedArticles } from '../api/articles'

export default function HomePage() {
  const [articles, setArticles] = useState([])
  const [error, setError] = useState('')

  useEffect(() => {
    let active = true
    listPublishedArticles()
      .then((data) => {
        if (!active) return
        const items = Array.isArray(data) ? data : []
        setArticles(items.slice(0, 3))
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
      <div className="hero">
        <h1>Rootix / Backend Engineer</h1>
        <p>
          Systems are narratives. I write them in code, verify them in logs, and
          refine them in public.
        </p>
      </div>

      <section className="manifesto">
        <h2>Beliefs</h2>
        <ul>
          <li>Protocols over presentations.</li>
          <li>Latency is a design decision.</li>
          <li>Every interface is a contract.</li>
          <li>Operational truth beats polished fiction.</li>
        </ul>
      </section>

      <section className="article-preview">
        <div className="section-header">
          <h2>Recent Articles</h2>
          <Link to="/articles">View all</Link>
        </div>
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
    </section>
  )
}
