import { createRouter, createWebHistory } from 'vue-router'
import LoginView from '@/views/LoginView.vue'
import SignupView from '@/views/SignupView.vue'
import constants from '@/consts/consts';

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path:'/',
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
      component: LoginView
    },
    {
      path: '/home',
      name: 'home',
      meta: { requiresAuth: true },
      // lazy-load when the route is visited
      component: () => import('@/views/HomeView.vue')
    }
  ]
})

router.beforeEach((to, from, next) => {
  const requiresAuth = to.matched.some(record => record.meta.requiresAuth);
  const sessionJWT = localStorage.getItem(constants.localStorageKeys.sessionJWT);
  if (!requiresAuth) {
    next();
    return
  }
  if (requiresAuth && !sessionJWT) {
    next('/login');
    localStorage.removeItem(constants.localStorageKeys.sessionJWT);
    return
  } else {
    const formData = {
      jwt: sessionJWT,
    };
    const fetchOptions = {
        headers: {
        'Content-Type': 'application/json',
        },
        method: 'POST',
        mode: "cors",
        body: JSON.stringify(formData),
    };
    fetch('http://localhost:8080/v1/session', fetchOptions)
    .then(response => response.json())
    .then(data => {
      if (!data || !data.jwt) {
        next('/login');
        localStorage.removeItem(constants.localStorageKeys.sessionJWT);
        return
      }
      if (sessionJWT !== data.jwt) {
        next('/login');
        localStorage.removeItem(constants.localStorageKeys.sessionJWT);
        return
      }
      next();
      return
    })
    .catch(error => {
      console.error('Error submitting form:', error);
      next('/login');
      localStorage.removeItem(constants.localStorageKeys.sessionJWT);
    });
  }
});

export default router
