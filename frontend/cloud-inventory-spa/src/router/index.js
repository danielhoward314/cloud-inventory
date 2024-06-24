import { createRouter, createWebHistory } from 'vue-router'
import LoginView from '../views/LoginView.vue'
import SignupView from '../views/SignupView.vue'
import constants from '@/consts/consts';

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
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
      component: () => import('../views/HomeView.vue')
    }
  ]
})

router.beforeEach((to, from, next) => {
  const requiresAuth = to.matched.some(record => record.meta.requiresAuth);
  const sessionJWT = localStorage.getItem(constants.localStorageKeys.sessionJWT);
  // TODO: validate sessionJWT
  /*
    with only this implementation,
    a user could do `localStorage.setItem('sessionJWT', 'not a real jwt');`
    and gain access without authenticating
  */
  if (requiresAuth && !sessionJWT) {
    next('/login');
  } else {
    next();
  }
});

export default router
