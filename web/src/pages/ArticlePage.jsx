import { useEffect, useState } from 'react'
import { useParams } from 'react-router-dom'
import { getPublishedArticle } from '../api/articles'

export default function ArticlePage() {
  const { slug } = useParams()
  const [article, setArticle] = useState(null)
  const [error, setError] = useState('')

  useEffect(() => {
    let active = true
    getPublishedArticle(slug)
      .then((data) => {
        if (active) {
          setArticle(data)
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
  }, [slug])

  if (error) {
    return (
      <section className="page">
        <h1>Article not found</h1>
        <p className="muted">{error}</p>
      </section>
    )
  }

  if (!article) {
    return (
      <section className="page">
        <p className="muted">Loading...</p>
      </section>
    )
  }

  return (
    <article className="page article">
      <header>
        <h1>{article.title}</h1>
        <p className="meta">Published {new Date(article.createdAt).toLocaleDateString()}</p>
      </header>
      <div
        className="article-body"
        dangerouslySetInnerHTML={{ __html: article.html }}
      />
    </article>
  )
}
