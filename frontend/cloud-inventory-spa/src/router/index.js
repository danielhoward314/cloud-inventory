import { createRouter, createWebHistory } from 'vue-router';
import LoginView from '@/views/LoginView.vue';
import SignupView from '@/views/SignupView.vue';
import constants from '@/consts/consts';
import { useAuthStore } from '@/stores/auth';

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'base',
      redirect: '/login',
    },
    {
      path: '/signup',
      name: 'signup',
      component: SignupView,
    },
    {
      path: '/login',
      name: 'login',
      component: LoginView,
    },
    {
      path: '/home',
      name: 'home',
      meta: { requiresAuth: true },
      component: () => import('@/views/HomeView.vue'),
    },
    {
      path: '/account',
      name: 'account',
      meta: { requiresAuth: true },
      component: () => import('@/views/AccountView.vue'),
    },
  ],
});

router.beforeEach(async (to, from, next) => {
  const authStore = useAuthStore();
  const requiresAuth = to.matched.some(record => record.meta.requiresAuth);
  const token = localStorage.getItem(constants.localStorageKeys.adminUiAccessToken);

  // unprotected route, allow navigation
  if (!requiresAuth) {
    next();
    return;
  }

  // no access token in localStorage, user must authenticate
  if (requiresAuth && !token) {
    next('/login');
    return;
  }

  // after one successful call to /session, should fall into this case
  // to avoid further API calls
  if (authStore.isAuthenticated && authStore.token === token) {
    next();
    return;
  }

  // call session API to validate the access token
  // if valid, set in-memory state, which won't persist across app reloads
  // 401s signal a valid, but expired access token
  // call the refresh API, which will return a new access token if valid
  try {
    const formData = {
      jwt: token,
    };
    const fetchOptions = {
      headers: {
        'Content-Type': 'application/json',
      },
      method: 'POST',
      mode: 'cors',
      body: JSON.stringify(formData),
    };
    const response = await fetch('http://localhost:8080/v1/session', fetchOptions);

    if (response.status === 200) {
      const data = await response.json();
      if (!data || !data.jwt || data.jwt !== token) {
        localStorage.removeItem(constants.localStorageKeys.adminUiAccessToken);
        authStore.clearAuthentication();
        next('/login');
        return;
      }
      authStore.setAuthenticated(true, token);
      next();
      return;
    } else if (response.status === 401) {
      const refreshResponse = await fetch('http://localhost:8080/v1/refresh', fetchOptions);
      const refreshData = await refreshResponse.json();
      if (refreshData && refreshData.jwt) {
        localStorage.setItem(constants.localStorageKeys.adminUiAccessToken, refreshData.jwt);
        authStore.setAuthenticated(true, refreshData.jwt);
        next('/home');
        return;
      } else {
        localStorage.removeItem(constants.localStorageKeys.adminUiAccessToken);
        authStore.clearAuthentication();
        next('/login');
        return;
      }
    } else {
      localStorage.removeItem(constants.localStorageKeys.adminUiAccessToken);
      authStore.clearAuthentication();
      next('/login');
      return;
    }
  } catch (error) {
    console.error('Error submitting form:', error);
    localStorage.removeItem(constants.localStorageKeys.adminUiAccessToken);
    authStore.clearAuthentication();
    next('/login');
    return;
  }
});

export default router;
