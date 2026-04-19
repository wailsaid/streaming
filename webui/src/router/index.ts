import { createRouter, createWebHistory } from 'vue-router'
import { useAuth } from '@/composables/useAuth'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    // Auth routes (public)
    {
      path: '/login',
      name: 'login',
      component: () => import('../views/login.vue'),
      meta: { public: true },
    },
    {
      path: '/register',
      name: 'register',
      component: () => import('../views/register.vue'),
      meta: { public: true },
    },

    // Protected app routes
    {
      path: '/',
      name: 'home',
      component: () => import('../views/home.vue'),
    },
    {
      path: '/watch/:videoId',
      name: 'watch',
      component: () => import('../views/watch.vue'),
    },
    {
      path: '/profile/upload',
      name: 'upload',
      component: () => import('../views/upload.vue'),
    },
    {
      path: '/profile',
      name: 'profile',
      component: () => import('../views/upload.vue'),
    },
    {
      path: '/trending',
      name: 'trending',
      component: () => import('../views/home.vue'),
    },
    {
      path: '/subscriptions',
      name: 'subscriptions',
      component: () => import('../views/home.vue'),
    },
    {
      path: '/profile/history',
      name: 'history',
      component: () => import('../views/home.vue'),
    },
    {
      path: '/profile/watch-later',
      name: 'watch-later',
      component: () => import('../views/home.vue'),
    },
    {
      path: '/profile/liked-videos',
      name: 'liked-videos',
      component: () => import('../views/home.vue'),
    },
  ],
})

// Navigation guard — redirect unauthenticated users to /login
router.beforeEach((to, _from, next) => {
  const { isAuthenticated } = useAuth()
  if (to.meta.public) {
    next()
    return
  }
  if (!isAuthenticated()) {
    next({ name: 'login' })
    return
  }
  next()
})

export default router
