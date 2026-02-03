import { useEffect, useMemo, useState } from 'react'
import { marked } from 'marked'
import {
  loginAdmin,
  listAdminArticles,
  createAdminArticle,
  updateAdminArticle,
  deleteAdminArticle,
  publishAdminArticle,
  unpublishAdminArticle
} from '../api/articles'

const EMPTY_FORM = {
  id: null,
  title: '',
  markdown: '',
  status: 'draft'
}

export default function AdminPage() {
  const [token, setToken] = useState(() => localStorage.getItem('adminToken') || '')
  const [loginError, setLoginError] = useState('')
  const [loginForm, setLoginForm] = useState({ email: '', password: '' })
  const [articles, setArticles] = useState([])
  const [form, setForm] = useState(EMPTY_FORM)
  const [statusMessage, setStatusMessage] = useState('')
  const [error, setError] = useState('')

  useEffect(() => {
    if (!token) {
      return
    }
    listAdminArticles(token)
      .then((data) => setArticles(data))
      .catch((err) => setError(err.message))
  }, [token])

  const previewHTML = useMemo(() => {
    return marked.parse(form.markdown || '')
  }, [form.markdown])

  function handleLogin(event) {
    event.preventDefault()
    setLoginError('')
    loginAdmin(loginForm)
      .then((data) => {
        localStorage.setItem('adminToken', data.token)
        setToken(data.token)
      })
      .catch((err) => setLoginError(err.message))
  }

  function handleSelect(article) {
    setForm({
      id: article.id,
      title: article.title,
      markdown: article.markdown,
      status: article.status
    })
    setStatusMessage('')
    setError('')
  }

  function handleNew() {
    setForm(EMPTY_FORM)
    setStatusMessage('')
    setError('')
  }

  function handleSave() {
    setStatusMessage('')
    setError('')
    const payload = { title: form.title, markdown: form.markdown }
    const action = form.id
      ? updateAdminArticle(token, form.id, payload)
      : createAdminArticle(token, payload)

    action
      .then((data) => {
        setStatusMessage('Saved.')
        return listAdminArticles(token).then((items) => {
          setArticles(items)
          setForm({
            id: data.id,
            title: data.title,
            markdown: data.markdown,
            status: data.status
          })
        })
      })
      .catch((err) => setError(err.message))
  }

  function handleDelete() {
    if (!form.id) {
      return
    }
    setStatusMessage('')
    setError('')
    deleteAdminArticle(token, form.id)
      .then(() => {
        setStatusMessage('Deleted.')
        setForm(EMPTY_FORM)
        return listAdminArticles(token).then((items) => setArticles(items))
      })
      .catch((err) => setError(err.message))
  }

  function handlePublish() {
    if (!form.id) {
      return
    }
    setStatusMessage('')
    setError('')
    publishAdminArticle(token, form.id)
      .then((data) => {
        setForm((prev) => ({ ...prev, status: data.status }))
        setStatusMessage('Published.')
        return listAdminArticles(token).then((items) => setArticles(items))
      })
      .catch((err) => setError(err.message))
  }

  function handleUnpublish() {
    if (!form.id) {
      return
    }
    setStatusMessage('')
    setError('')
    unpublishAdminArticle(token, form.id)
      .then((data) => {
        setForm((prev) => ({ ...prev, status: data.status }))
        setStatusMessage('Unpublished.')
        return listAdminArticles(token).then((items) => setArticles(items))
      })
      .catch((err) => setError(err.message))
  }

  if (!token) {
    return (
      <section className="admin">
        <h1>Admin Access</h1>
        <form className="admin-login" onSubmit={handleLogin}>
          <label>
            Email
            <input
              type="email"
              value={loginForm.email}
              onChange={(event) => setLoginForm({ ...loginForm, email: event.target.value })}
              required
            />
          </label>
          <label>
            Password
            <input
              type="password"
              value={loginForm.password}
              onChange={(event) => setLoginForm({ ...loginForm, password: event.target.value })}
              required
            />
          </label>
          {loginError && <p className="muted">{loginError}</p>}
          <button type="submit">Enter</button>
        </form>
      </section>
    )
  }

  return (
    <section className="admin">
      <div className="admin-header">
        <div>
          <h1>Articles Control Room</h1>
          <p className="muted">Draft, publish, and archive written work.</p>
        </div>
        <button type="button" onClick={() => {
          localStorage.removeItem('adminToken')
          setToken('')
        }}>
          Logout
        </button>
      </div>

      <div className="admin-grid">
        <div className="admin-list">
          <div className="admin-list-header">
            <h2>Inventory</h2>
            <button type="button" onClick={handleNew}>New Draft</button>
          </div>
          <table>
            <thead>
              <tr>
                <th>Title</th>
                <th>Status</th>
              </tr>
            </thead>
            <tbody>
              {articles.map((article) => (
                <tr key={article.id} onClick={() => handleSelect(article)}>
                  <td>{article.title}</td>
                  <td>{article.status}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>

        <div className="admin-editor">
          <div className="admin-actions">
            <button type="button" onClick={handleSave}>Save</button>
            <button type="button" onClick={handlePublish} disabled={!form.id || form.status === 'published'}>
              Publish
            </button>
            <button type="button" onClick={handleUnpublish} disabled={!form.id || form.status === 'draft'}>
              Unpublish
            </button>
            <button type="button" onClick={handleDelete} disabled={!form.id}>
              Delete
            </button>
          </div>

          {statusMessage && <p className="muted">{statusMessage}</p>}
          {error && <p className="muted">{error}</p>}

          <label>
            Title
            <input
              type="text"
              value={form.title}
              onChange={(event) => setForm({ ...form, title: event.target.value })}
            />
          </label>

          <div className="editor-grid">
            <div>
              <label>
                Markdown
                <textarea
                  rows="18"
                  value={form.markdown}
                  onChange={(event) => setForm({ ...form, markdown: event.target.value })}
                />
              </label>
            </div>
            <div className="preview">
              <h3>Preview</h3>
              <div className="preview-body" dangerouslySetInnerHTML={{ __html: previewHTML }} />
            </div>
          </div>
        </div>
      </div>
    </section>
  )
}
