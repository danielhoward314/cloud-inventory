import axios from 'axios'
import constants from '../consts/consts'
import router from '../router'

export const baseClient = axios.create({
  baseURL: 'http://localhost:8080/v1',
  headers: {
    'Content-Type': 'application/json'
  },
  mode: 'cors'
})

baseClient.interceptors.request.use(
  async (config) => {
    let token = localStorage.getItem(constants.localStorageKeys.apiAccessToken)
    if (token) {
      config.headers['Authorization'] = `Bearer ${token}`
    }
    const isMethodWithBody =
      config.method === 'post' ||
      config.method === 'put' ||
      config.method === 'patch' ||
      config.method === 'delete'
    // `!config._retry` is used to avoid re-stringifying a request body
    // for a retry that has already had this done
    if (isMethodWithBody && config.data && !config._retry) {
      config.data = JSON.stringify(config.data)
    }
    return config
  },
  (error) => Promise.reject(error)
)

baseClient.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config
    if (error.response.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true
      try {
        const refreshToken = localStorage.getItem(constants.localStorageKeys.adminUiRefreshToken)
        const refreshFormData = {
          jwt: refreshToken,
          claimsType: 2 // ClaimsType_API_AUTHORIZATION
        }
        const refreshResponse = await axios.post(
          'http://localhost:8080/v1/refresh',
          refreshFormData
        )
        const hasToken = refreshResponse && refreshResponse.data && refreshResponse.data.jwt
        if (!hasToken) {
          console.error('No API access token in refresh response')
          router.push('/login')
          return
        }
        localStorage.setItem(constants.localStorageKeys.apiAccessToken, refreshResponse.data.jwt)
        originalRequest.headers['Authorization'] = `Bearer ${refreshResponse.data.jwt}`
        return baseClient(originalRequest)
      } catch (error) {
        console.error('Error refreshing token:', error)
        router.push('/login')
        return
      }
    }
    return Promise.reject(error)
  }
)
