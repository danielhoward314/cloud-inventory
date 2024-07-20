import { createRouter, createWebHistory } from 'vue-router'
import LoginView from '@/views/LoginView.vue'
import SignupView from '@/views/SignupView.vue'
import constants from '@/consts/consts'
import { useAuthStore } from '@/stores/auth'
import { parseJwt } from '@/util/jwt'
import { useOrganizationStore } from '@/stores/organization'
import { useAdministratorStore } from '@/stores/administrator'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'base',
      redirect: '/login'
    },
    {
      path: '/signup',
      name: 'signup',
      component: SignupView
    },
    {
      path: '/login',
      name: 'login',
      component: LoginView
    },
    {
      path: '/home',
      name: 'home',
      meta: { requiresAuth: true },
      component: () => import('@/views/HomeView.vue')
    },
    {
      path: '/account',
      name: 'account',
      meta: { requiresAuth: true },
      component: () => import('@/views/AccountView.vue')
    }
  ]
})

router.beforeEach(async (to, from, next) => {
  const authStore = useAuthStore()
  const orgStore = useOrganizationStore()
  const adminStore = useAdministratorStore()
  const requiresAuth = to.matched.some((record) => record.meta.requiresAuth)
  const accessToken = localStorage.getItem(constants.localStorageKeys.adminUiAccessToken)

  // Allow navigation for unprotected routes.
  if (!requiresAuth) {
    next()
    return
  }

  // No access token in localStorage? The user must authenticate.
  if (requiresAuth && !accessToken) {
    next('/login')
    return
  }

  // After one successful call to /session, we should fall into this case.
  // This avoids superfluous /session API calls as the user navigates
  // to different routes.
  if (authStore.isAuthenticated && authStore.token === accessToken) {
    next()
    return
  }

  // Call session API to validate the access token.
  // If valid, set in-memory pinia state,
  // which will spare more session API calls as user navigates,
  // but won't persist across app reloads.
  // A 401 status code signals a valid-but-expired access token.
  // Call the refresh API to validate the refresh token.
  // If valid, return a new access token to set in localStorage and pinia.
  try {
    const formData = {
      jwt: accessToken
    }
    const fetchOptions = {
      headers: {
        'Content-Type': 'application/json'
      },
      method: 'POST',
      mode: 'cors',
      body: JSON.stringify(formData)
    }
    const response = await fetch('http://localhost:8080/v1/session', fetchOptions)

    if (response.status === 200) {
      const data = await response.json()
      if (!data || !data.jwt || data.jwt !== accessToken) {
        localStorage.removeItem(constants.localStorageKeys.adminUiAccessToken)
        authStore.clearAuthentication()
        next('/login')
        return
      }
      const parsedToken = parseJwt(data.jwt)
      orgStore.setOrganization({
        id: parsedToken.organization_id
      })
      adminStore.setAdministrator({
        id: parsedToken.sub,
        authorizationRole: parsedToken.authorization_role
      })
      authStore.setAuthenticated(true, accessToken)
      next()
      return
    } else if (response.status === 401) {
      const refreshToken = localStorage.getItem(constants.localStorageKeys.adminUiRefreshToken)
      const refreshFormData = {
        jwt: refreshToken,
        claimsType: 1 // ClaimsType_ADMIN_UI_SESSION
      }
      const refreshFetchOptions = {
        headers: {
          'Content-Type': 'application/json'
        },
        method: 'POST',
        mode: 'cors',
        body: JSON.stringify(refreshFormData)
      }
      const refreshResponse = await fetch('http://localhost:8080/v1/refresh', refreshFetchOptions)
      const refreshData = await refreshResponse.json()
      if (refreshData && refreshData.jwt) {
        localStorage.setItem(constants.localStorageKeys.adminUiAccessToken, refreshData.jwt)
        authStore.setAuthenticated(true, refreshData.jwt)
        next('/home')
        return
      } else {
        localStorage.removeItem(constants.localStorageKeys.adminUiAccessToken)
        authStore.clearAuthentication()
        next('/login')
        return
      }
    } else {
      localStorage.removeItem(constants.localStorageKeys.adminUiAccessToken)
      authStore.clearAuthentication()
      next('/login')
      return
    }
  } catch (error) {
    console.error('Error submitting form:', error)
    localStorage.removeItem(constants.localStorageKeys.adminUiAccessToken)
    authStore.clearAuthentication()
    next('/login')
    return
  }
})

export default router
