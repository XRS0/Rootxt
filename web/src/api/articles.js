import { request } from './client'

export function listPublishedArticles() {
  return request('/api/articles')
}

export function getPublishedArticle(slug) {
  return request(`/api/articles/${slug}`)
}

export function loginAdmin(payload) {
  return request('/api/admin/login', {
    method: 'POST',
    body: JSON.stringify(payload)
  })
}

function authHeaders(token) {
  return {
    Authorization: `Bearer ${token}`
  }
}

export function listAdminArticles(token) {
  return request('/api/admin/articles', {
    headers: authHeaders(token)
  })
}

export function createAdminArticle(token, payload) {
  return request('/api/admin/articles', {
    method: 'POST',
    headers: authHeaders(token),
    body: JSON.stringify(payload)
  })
}

export function updateAdminArticle(token, id, payload) {
  return request(`/api/admin/articles/${id}`, {
    method: 'PUT',
    headers: authHeaders(token),
    body: JSON.stringify(payload)
  })
}

export function deleteAdminArticle(token, id) {
  return request(`/api/admin/articles/${id}`, {
    method: 'DELETE',
    headers: authHeaders(token)
  })
}

export function publishAdminArticle(token, id) {
  return request(`/api/admin/articles/${id}/publish`, {
    method: 'POST',
    headers: authHeaders(token)
  })
}

export function unpublishAdminArticle(token, id) {
  return request(`/api/admin/articles/${id}/unpublish`, {
    method: 'POST',
    headers: authHeaders(token)
  })
}
